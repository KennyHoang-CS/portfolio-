package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DamageNumber struct {
	X, Y  float64
	Value int
	Life  float64
}

type Crystal struct {
	Pos    Vec
	Alive  bool
	Sprite *ebiten.Image
}

type Ability struct {
	Name        string
	Description string
	Level       int
	Icon        *ebiten.Image
	Apply       func(g *Game)
}

type Orb struct {
	Angle  float64
	Radius float64
	Speed  float64
	Damage float64
	Sprite *ebiten.Image
}

type Dagger struct {
	Pos    Vec
	Vel    Vec
	Damage float64
	Alive  bool
	Sprite *ebiten.Image
}

type FireTile struct {
	Pos    Vec
	Life   float64
	Damage float64
}

type StunWave struct {
	Pos    Vec
	Radius float64
	Life   float64
}

type Game struct {
	player        *Player
	lastTime      time.Time
	Enemies       []*Enemy
	DamageNumbers []*DamageNumber
	Crystals      []*Crystal

	waveTimer  float64
	waveNumber int

	bg            *ebiten.Image
	crystalSprite *ebiten.Image

	// Abilities
	AvailableAbilities []*Ability
	LevelUpMenuOpen    bool
	LevelUpChoices     []*Ability

	// Spell systems
	Orbs      []*Orb
	Daggers   []*Dagger
	FireTiles []*FireTile
	StunWaves []*StunWave

	// Global spell stats
	CrystalMagnetRadius float64
	InfiniteLoopChance  float64

	// Dagger
	DaggerCooldown float64
	DaggerRate     float64
}

func PlaceholderIcon(r, g, b uint8) *ebiten.Image {
	img := ebiten.NewImage(32, 32)
	img.Fill(color.RGBA{r, g, b, 255})
	return img
}

func NewGame() *Game {
	InitMonsters()

	g := &Game{
		player:              NewPlayer(),
		lastTime:            time.Now(),
		waveTimer:           3.0,
		waveNumber:          1,
		bg:                  LoadImage("assets/grass_tile.png"),
		crystalSprite:       LoadImage("assets/crystal.png"),
		Crystals:            []*Crystal{},
		CrystalMagnetRadius: 0,
		InfiniteLoopChance:  0.0,
		DaggerCooldown:      0,
		DaggerRate:          3.2, // fires every 3.2 seconds by default
	}

	g.initAbilities()

	for i := 0; i < 10; i++ {
		mt := MonsterPool[rand.Intn(len(MonsterPool))]
		g.Enemies = append(g.Enemies, NewEnemy(
			mt,
			float64(200+i*50),
			float64(200+i*30),
		))
	}

	return g
}

func (g *Game) initAbilities() {
	g.AvailableAbilities = []*Ability{
		// {
		// 	Name:        "CPU Overclock",
		// 	Description: "Increase slash damage and attack speed.",
		// 	Icon:        PlaceholderIcon(255, 80, 80),
		// 	Apply: func(g *Game) {
		// 		g.player.SlashDamage += 1
		// 		g.player.SlashCooldownBase *= 0.9
		// 		g.player.SlashRadius += 5
		// 		g.player.AttackSpeedMultiplier += 0.05
		// 	},
		// },
		// {
		// 	Name:        "Garbage Collection",
		// 	Description: "Magnet for crystals and tiny DoT to weak enemies.",
		// 	Icon:        PlaceholderIcon(80, 80, 255),
		// 	Apply: func(g *Game) {
		// 		g.CrystalMagnetRadius += 40
		// 	},
		// },
		// {
		// 	Name:        "Type Error",
		// 	Description: "Slow aura around the player. (placeholder)",
		// 	Icon:        PlaceholderIcon(160, 80, 255),
		// 	Apply: func(g *Game) {
		// 		// Implement slow aura in enemy update/contact later
		// 	},
		// },
		{
			Name:        "Index Out of Range",
			Description: "Shoot a dagger in facing direction. (placeholder buff)",
			Icon:        PlaceholderIcon(255, 255, 80),
			Apply: func(g *Game) {
				// Fire a dagger immediately
				g.player.HasDagger = true
				dir := float64(g.player.lastDir)
				d := &Dagger{
					Pos:    Vec{X: g.player.Pos.X, Y: g.player.Pos.Y},
					Vel:    Vec{X: dir * 400, Y: 0},
					Damage: 10 + float64(g.player.Level)*2,
					Alive:  true,
					Sprite: PlaceholderIcon(255, 255, 80),
				}
				// Increase dagger fire rate + damage
				g.DaggerRate *= 0.60 // 10% faster firing
				g.Daggers = append(g.Daggers, d)
			},
		},
		// {
		// 	Name:        "Memory Leak",
		// 	Description: "Leave fire tiles behind you. (placeholder buff)",
		// 	Icon:        PlaceholderIcon(255, 120, 0),
		// 	Apply: func(g *Game) {
		// 		// Spawn fire tile at player's feet
		// 		ft := &FireTile{
		// 			Pos:    Vec{X: g.player.Pos.X, Y: g.player.Pos.Y},
		// 			Life:   3.0,
		// 			Damage: 5 + float64(g.player.Level),
		// 		}
		// 		g.FireTiles = append(g.FireTiles, ft)
		// 	},
		// },
		// {
		// 	Name:        "Segmentation Fault",
		// 	Description: "Periodic stun wave. (placeholder buff)",
		// 	Icon:        PlaceholderIcon(255, 255, 255),
		// 	Apply: func(g *Game) {
		// 		// Increase stun radius or duration later
		// 	},
		// },
		// {
		// 	Name:        "Infinite Loop",
		// 	Description: "Chance to repeat spells.",
		// 	Icon:        PlaceholderIcon(200, 200, 200),
		// 	Apply: func(g *Game) {
		// 		g.InfiniteLoopChance += 0.05
		// 	},
		// },
		// {
		// 	Name:        "Race Condition",
		// 	Description: "Increase movement and attack speed.",
		// 	Icon:        PlaceholderIcon(255, 200, 0),
		// 	Apply: func(g *Game) {
		// 		g.player.MoveSpeedMultiplier += 0.1
		// 		g.player.AttackSpeedMultiplier += 0.05
		// 	},
		// },
		// {
		// 	Name:        "Deadlock",
		// 	Description: "Freeze aura. (placeholder)",
		// 	Icon:        PlaceholderIcon(0, 200, 255),
		// 	Apply: func(g *Game) {
		// 		// Implement freeze aura later
		// 	},
		// },
		// {
		// 	Name:        "Stack Overflow",
		// 	Description: "Increase HP and regen.",
		// 	Icon:        PlaceholderIcon(255, 0, 0),
		// 	Apply: func(g *Game) {
		// 		g.player.MaxHP += 20
		// 		g.player.HP += 20
		// 		g.player.RegenPerSecond += 1
		// 	},
		// },
		// {
		// 	Name:        "Heap Allocation",
		// 	Description: "Orbiting orbs around the player.",
		// 	Icon:        PlaceholderIcon(0, 255, 0),
		// 	Apply: func(g *Game) {
		// 		orb := &Orb{
		// 			Angle:  rand.Float64() * math.Pi * 2,
		// 			Radius: 60,
		// 			Speed:  1.5,
		// 			Damage: 5,
		// 			Sprite: PlaceholderIcon(0, 255, 0),
		// 		}
		// 		g.Orbs = append(g.Orbs, orb)
		// 	},
		// },
		// {
		// 	Name:        "Compiler Optimization",
		// 	Description: "Global cooldown and speed buffs.",
		// 	Icon:        PlaceholderIcon(0, 255, 120),
		// 	Apply: func(g *Game) {
		// 		g.player.MoveSpeedMultiplier += 0.05
		// 		g.player.AttackSpeedMultiplier += 0.05
		// 	},
		// },
		// {
		// 	Name:        "Thread Pool",
		// 	Description: "Summoned minions (placeholder: extra orbs).",
		// 	Icon:        PlaceholderIcon(120, 120, 255),
		// 	Apply: func(g *Game) {
		// 		orb := &Orb{
		// 			Angle:  rand.Float64() * math.Pi * 2,
		// 			Radius: 80,
		// 			Speed:  2.0,
		// 			Damage: 4,
		// 			Sprite: PlaceholderIcon(120, 120, 255),
		// 		}
		// 		g.Orbs = append(g.Orbs, orb)
		// 	},
		// },
		// {
		// 	Name:        "Dependency Injection",
		// 	Description: "Poison/DoT on hit (placeholder: extra damage).",
		// 	Icon:        PlaceholderIcon(0, 200, 0),
		// 	Apply: func(g *Game) {
		// 		g.player.SlashDamage += 0.5
		// 	},
		// },
		// {
		// 	Name:        "API Rate Limit",
		// 	Description: "Reduce spell cooldowns.",
		// 	Icon:        PlaceholderIcon(200, 200, 0),
		// 	Apply: func(g *Game) {
		// 		g.player.AttackSpeedMultiplier += 0.05
		// 	},
		// },
	}
}

func (g *Game) CameraOffset() (float64, float64) {
	return g.player.Pos.X - 400, g.player.Pos.Y - 300
}

func (g *Game) Update() error {
	now := time.Now()
	dt := now.Sub(g.lastTime).Seconds()
	g.lastTime = now

	// If level-up menu is open, pause gameplay and handle input
	if g.LevelUpMenuOpen {
		// Press 1, 2, or 3 to pick an ability
		if ebiten.IsKeyPressed(ebiten.Key1) {
			g.pickAbility(0)
			g.LevelUpMenuOpen = false
		}
		if ebiten.IsKeyPressed(ebiten.Key2) && len(g.LevelUpChoices) > 1 {
			g.pickAbility(1)
			g.LevelUpMenuOpen = false
		}
		if ebiten.IsKeyPressed(ebiten.Key3) && len(g.LevelUpChoices) > 2 {
			g.pickAbility(2)
			g.LevelUpMenuOpen = false
		}
		return nil
	}

	g.player.Update(dt)

	for _, e := range g.Enemies {
		e.Update(dt, g.player.Pos)
	}

	g.resolveEnemyCollisions(dt)
	g.resolvePlayerCollision(dt)

	g.handleCombat(dt)
	g.handleEnemyContact(dt)
	g.handleCrystalPickup()
	g.updateDamageNumbers(dt)
	g.updateWaves(dt)
	g.updateSpells(dt)
	g.checkLevelUp()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	camX, camY := g.CameraOffset()

	tileW, tileH := g.bg.Size()
	for x := -tileW; x < 800+tileW; x += tileW {
		for y := -tileH; y < 600+tileH; y += tileH {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(
				float64(x)-math.Mod(camX, float64(tileW)),
				float64(y)-math.Mod(camY, float64(tileH)),
			)
			screen.DrawImage(g.bg, op)
		}
	}

	g.player.DrawWithCamera(screen, camX, camY)

	// Draw enemies
	for _, e := range g.Enemies {
		if !e.Alive || e.Type == nil || e.Type.Sprite == nil {
			continue
		}

		w, h := e.Type.Sprite.Size()

		// Keep your original scale logic
		scale := e.Type.Scale * (e.Radius * 2) / float64(h)

		// Compute top-left BEFORE rounding
		px := e.Pos.X - camX - float64(w)*scale/2
		py := e.Pos.Y - camY - float64(h)*scale/2

		// Snap ONLY translation to pixel grid
		px = math.Round(px)
		py = math.Round(py)

		// Outline
		outline := &ebiten.DrawImageOptions{}
		outline.Filter = ebiten.FilterNearest
		outline.GeoM.Scale(scale, scale)
		outline.GeoM.Translate(px, py)
		outline.ColorM.Scale(0, 0, 0, 0.6)
		screen.DrawImage(e.Type.Sprite, outline)

		// Main sprite
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterNearest
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(px, py)

		if e.HitFlash > 0 {
			op.ColorM.Scale(2, 2, 2, 1)
		}

		screen.DrawImage(e.Type.Sprite, op)
	}

	// Draw crystals (scaled + bobbing)
	for _, c := range g.Crystals {
		if !c.Alive {
			continue
		}

		w, h := c.Sprite.Size()
		op := &ebiten.DrawImageOptions{}

		scale := 0.1
		op.GeoM.Scale(scale, scale)

		bob := math.Sin(float64(time.Now().UnixNano())*0.000000005) * 2

		op.GeoM.Translate(
			c.Pos.X-camX-float64(w)*scale/2,
			c.Pos.Y-camY-float64(h)*scale/2+bob,
		)

		screen.DrawImage(c.Sprite, op)
	}

	// Draw orbs
	for _, orb := range g.Orbs {
		op := &ebiten.DrawImageOptions{}
		x := g.player.Pos.X + math.Cos(orb.Angle)*orb.Radius - camX
		y := g.player.Pos.Y + math.Sin(orb.Angle)*orb.Radius - camY
		op.GeoM.Translate(x, y)
		screen.DrawImage(orb.Sprite, op)
	}

	// Draw daggers
	for _, d := range g.Daggers {
		if !d.Alive {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(d.Pos.X-camX, d.Pos.Y-camY)
		screen.DrawImage(d.Sprite, op)
	}

	// Draw fire tiles
	for _, ft := range g.FireTiles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(ft.Pos.X-camX, ft.Pos.Y-camY)
		screen.DrawImage(PlaceholderIcon(255, 120, 0), op)
	}

	// Damage numbers
	for _, dn := range g.DamageNumbers {
		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("%d", dn.Value),
			int(dn.X-camX),
			int(dn.Y-camY),
		)
	}

	g.drawHPBar(screen)
	g.drawXPBar(screen)

	if g.LevelUpMenuOpen {
		g.drawLevelUpMenu(screen)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 800, 600
}

func (g *Game) handleCombat(dt float64) {
	if !g.player.attacking {
		return
	}

	slashDir := float64(g.player.lastDir)
	slashX := g.player.Pos.X + slashDir*150
	slashY := g.player.Pos.Y - 5
	slashRadius := g.player.SlashRadius

	for _, e := range g.Enemies {
		if !e.Alive {
			continue
		}

		dx := e.Pos.X - slashX
		dy := e.Pos.Y - slashY
		dist := math.Hypot(dx, dy)

		if dist < slashRadius+e.Radius {
			e.HP -= g.player.SlashDamage
			e.HitFlash = 0.1

			kb := 120.0
			e.Pos.X += slashDir * kb * dt

			g.DamageNumbers = append(g.DamageNumbers, &DamageNumber{
				X:     e.Pos.X,
				Y:     e.Pos.Y - 20,
				Value: int(g.player.SlashDamage),
				Life:  1.0,
			})

			if e.HP <= 0 {
				e.Alive = false

				g.Crystals = append(g.Crystals, &Crystal{
					Pos:    Vec{X: e.Pos.X, Y: e.Pos.Y},
					Alive:  true,
					Sprite: g.crystalSprite,
				})
			}
		}
	}
}

func (g *Game) handleEnemyContact(dt float64) {
	for _, e := range g.Enemies {
		if !e.Alive {
			continue
		}

		dx := e.Pos.X - g.player.Pos.X
		dy := e.Pos.Y - g.player.Pos.Y
		dist := math.Hypot(dx, dy)

		if dist < e.Radius+20 {
			g.player.HP -= 10 * dt
			if g.player.HP < 0 {
				g.player.HP = 0
			}
		}
	}
}

func (g *Game) handleCrystalPickup() {
	for _, c := range g.Crystals {
		if !c.Alive {
			continue
		}

		dx := c.Pos.X - g.player.Pos.X
		dy := c.Pos.Y - g.player.Pos.Y
		dist := math.Hypot(dx, dy)

		// Magnet
		if g.CrystalMagnetRadius > 0 && dist < g.CrystalMagnetRadius && dist > 0 {
			dirX := dx / dist
			dirY := dy / dist
			c.Pos.X -= dirX * 200 * (1.0 / 60.0)
			c.Pos.Y -= dirY * 200 * (1.0 / 60.0)
		}

		if dist < 40 {
			c.Alive = false
			g.player.XP += 1
		}
	}
}

func (g *Game) updateDamageNumbers(dt float64) {
	for i := len(g.DamageNumbers) - 1; i >= 0; i-- {
		dn := g.DamageNumbers[i]
		dn.Y -= 30 * dt
		dn.Life -= dt

		if dn.Life <= 0 {
			g.DamageNumbers = append(g.DamageNumbers[:i], g.DamageNumbers[i+1:]...)
		}
	}
}

func (g *Game) updateWaves(dt float64) {
	g.waveTimer -= dt
	if g.waveTimer <= 0 {
		g.spawnWave()
		g.waveNumber++
		g.waveTimer = 5.0
	}
}

func (g *Game) spawnWave() {
	count := 5 + g.waveNumber*2

	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2
		dist := 400 + rand.Float64()*200

		x := g.player.Pos.X + math.Cos(angle)*dist
		y := g.player.Pos.Y + math.Sin(angle)*dist

		mt := MonsterPool[rand.Intn(len(MonsterPool))]
		g.Enemies = append(g.Enemies, NewEnemy(mt, x, y))
	}
}

func (g *Game) checkLevelUp() {
	xpNeeded := float64(g.player.Level * 10)

	if g.player.XP >= xpNeeded {
		g.player.XP -= xpNeeded
		g.player.Level++
		g.openLevelUpMenu()
	}
}

func (g *Game) openLevelUpMenu() {
	g.LevelUpMenuOpen = true
	g.LevelUpChoices = []*Ability{}

	indices := rand.Perm(len(g.AvailableAbilities))
	for i := 0; i < 3 && i < len(indices); i++ {
		g.LevelUpChoices = append(g.LevelUpChoices, g.AvailableAbilities[indices[i]])
	}
}

func (g *Game) pickAbility(idx int) {
	if idx < 0 || idx >= len(g.LevelUpChoices) {
		return
	}
	ability := g.LevelUpChoices[idx]
	ability.Level++
	ability.Apply(g)
	g.player.Abilities = append(g.player.Abilities, ability)
}

func (g *Game) updateSpells(dt float64) {
	// Orbs orbit and damage enemies
	for _, orb := range g.Orbs {
		orb.Angle += orb.Speed * dt
		ox := g.player.Pos.X + math.Cos(orb.Angle)*orb.Radius
		oy := g.player.Pos.Y + math.Sin(orb.Angle)*orb.Radius

		for _, e := range g.Enemies {
			if !e.Alive {
				continue
			}
			dx := e.Pos.X - ox
			dy := e.Pos.Y - oy
			dist := math.Hypot(dx, dy)
			if dist < e.Radius+10 {
				e.HP -= orb.Damage
				if e.HP <= 0 {
					e.Alive = false
					g.Crystals = append(g.Crystals, &Crystal{
						Pos:    Vec{X: e.Pos.X, Y: e.Pos.Y},
						Alive:  true,
						Sprite: g.crystalSprite,
					})
				}
			}
		}
	}

	// --- DAGGER FIRING (only if unlocked) ---
	if g.player.HasDagger {
		g.DaggerCooldown -= dt
		if g.DaggerCooldown <= 0 {
			g.fireDagger()
			g.DaggerCooldown = g.DaggerRate
		}
	}

	// --- DAGGER MOVEMENT + COLLISION ---
	for _, d := range g.Daggers {
		if !d.Alive {
			continue
		}

		d.Pos.X += d.Vel.X * dt
		d.Pos.Y += d.Vel.Y * dt

		// Hit enemies
		for _, e := range g.Enemies {
			if !e.Alive {
				continue
			}
			dx := e.Pos.X - d.Pos.X
			dy := e.Pos.Y - d.Pos.Y
			dist := math.Hypot(dx, dy)
			if dist < e.Radius+8 {
				e.HP -= d.Damage
				d.Alive = false
				if e.HP <= 0 {
					e.Alive = false
					g.Crystals = append(g.Crystals, &Crystal{
						Pos:    Vec{X: e.Pos.X, Y: e.Pos.Y},
						Alive:  true,
						Sprite: g.crystalSprite,
					})
				}
			}
		}
	}

	// Remove dead daggers
	for i := len(g.Daggers) - 1; i >= 0; i-- {
		if !g.Daggers[i].Alive {
			g.Daggers = append(g.Daggers[:i], g.Daggers[i+1:]...)
		}
	}

	// --- FIRE TRAIL ---
	for _, ft := range g.FireTiles {
		ft.Life -= dt
		if ft.Life <= 0 {
			ft.Life = 0
		}

		for _, e := range g.Enemies {
			if !e.Alive {
				continue
			}
			dx := e.Pos.X - ft.Pos.X
			dy := e.Pos.Y - ft.Pos.Y
			dist := math.Hypot(dx, dy)
			if dist < e.Radius+20 {
				e.HP -= ft.Damage * dt
				if e.HP <= 0 {
					e.Alive = false
					g.Crystals = append(g.Crystals, &Crystal{
						Pos:    Vec{X: e.Pos.X, Y: e.Pos.Y},
						Alive:  true,
						Sprite: g.crystalSprite,
					})
				}
			}
		}
	}

	// Remove expired fire tiles
	for i := len(g.FireTiles) - 1; i >= 0; i-- {
		if g.FireTiles[i].Life <= 0 {
			g.FireTiles = append(g.FireTiles[:i], g.FireTiles[i+1:]...)
		}
	}

}

func (g *Game) resolveEnemyCollisions(dt float64) {
	for i := 0; i < len(g.Enemies); i++ {
		e1 := g.Enemies[i]
		if !e1.Alive {
			continue
		}

		for j := i + 1; j < len(g.Enemies); j++ {
			e2 := g.Enemies[j]
			if !e2.Alive {
				continue
			}

			dx := e2.Pos.X - e1.Pos.X
			dy := e2.Pos.Y - e1.Pos.Y
			dist := math.Hypot(dx, dy)
			minDist := e1.Radius + e2.Radius

			if dist < minDist && dist > 0 {
				overlap := minDist - dist
				pushX := dx / dist * overlap * 0.5
				pushY := dy / dist * overlap * 0.5

				e1.Pos.X -= pushX
				e1.Pos.Y -= pushY
				e2.Pos.X += pushX
				e2.Pos.Y += pushY
			}
		}
	}
}

func (g *Game) resolvePlayerCollision(dt float64) {
	for _, e := range g.Enemies {
		if !e.Alive {
			continue
		}

		dx := g.player.Pos.X - e.Pos.X
		dy := g.player.Pos.Y - e.Pos.Y
		dist := math.Hypot(dx, dy)
		minDist := 20 + e.Radius

		if dist < minDist && dist > 0 {
			overlap := minDist - dist
			pushX := dx / dist * overlap
			pushY := dy / dist * overlap

			g.player.Pos.X += pushX
			g.player.Pos.Y += pushY
		}
	}
}

func (g *Game) drawHPBar(screen *ebiten.Image) {
	camX, camY := g.CameraOffset()

	barWidth := float32(40)
	barHeight := float32(6)

	x := float32(g.player.Pos.X - camX - float64(barWidth)/2)
	y := float32(g.player.Pos.Y - camY + 45)

	hpPercent := float32(g.player.HP / g.player.MaxHP)

	vector.FillRect(screen, x, y, barWidth, barHeight, color.RGBA{80, 0, 0, 255}, false)
	vector.FillRect(screen, x, y, barWidth*hpPercent, barHeight, color.RGBA{255, 40, 40, 255}, false)
}

func (g *Game) drawXPBar(screen *ebiten.Image) {
	barWidth := float32(800)
	barHeight := float32(10)
	x := float32(0)
	y := float32(0)

	xpNeeded := float32(g.player.Level * 10)
	xpPercent := float32(g.player.XP) / xpNeeded

	vector.FillRect(screen, x, y, barWidth, barHeight, color.RGBA{30, 30, 60, 255}, false)
	vector.FillRect(screen, x, y, barWidth*xpPercent, barHeight, color.RGBA{80, 80, 255, 255}, false)
}

func (g *Game) drawLevelUpMenu(screen *ebiten.Image) {
	// Simple VS-style centered box with 3 choices
	boxX := 100
	boxY := 150
	boxW := 600
	boxH := 300

	// Background box
	vector.FillRect(screen,
		float32(boxX), float32(boxY),
		float32(boxW), float32(boxH),
		color.RGBA{10, 10, 30, 230},
		false,
	)

	// Border
	vector.StrokeRect(screen,
		float32(boxX), float32(boxY),
		float32(boxW), float32(boxH),
		2,
		color.RGBA{200, 200, 255, 255},
		false,
	)

	ebitenutil.DebugPrintAt(screen, "LEVEL UP! Choose an ability (1/2/3)", boxX+20, boxY+20)

	for i, a := range g.LevelUpChoices {
		y := boxY + 60 + i*70
		// Icon placeholder
		if a.Icon != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(boxX+30), float64(y))
			screen.DrawImage(a.Icon, op)
		}

		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("%d: %s (Lv %d)", i+1, a.Name, a.Level+1),
			boxX+80, y,
		)
		ebitenutil.DebugPrintAt(screen,
			a.Description,
			boxX+80, y+20,
		)
	}
}

func (g *Game) fireDagger() {
	dir := float64(g.player.lastDir)

	d := &Dagger{
		Pos:    Vec{X: g.player.Pos.X, Y: g.player.Pos.Y},
		Vel:    Vec{X: dir * 400, Y: 0},
		Damage: 10 + float64(g.player.Level)*2,
		Alive:  true,
		Sprite: PlaceholderIcon(255, 255, 80),
	}

	g.Daggers = append(g.Daggers, d)
}
