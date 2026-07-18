package game

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font/opentype"
)

func NewGame() *Game {
	InitMonsters()

	g := &Game{
		player:              NewPlayer(),
		lastTime:            time.Now(),
		waveTimer:           3.0,
		waveNumber:          1,
		bg:                  nil,
		trainingRoomBg:      LoadImage("assets/training_room_tile.png"),
		grassBg:             LoadImage("assets/grass_tile.png"),
		crystalSprite:       LoadImage("assets/crystal.png"),
		Crystals:            []*Crystal{},
		CrystalMagnetRadius: 0,
		InfiniteLoopChance:  0.0,
		DaggerCooldown:      0,
		DaggerRate:          3.2, // fires every 3.2 seconds by default
		GameState:           StateTitle,
		TitleBG:             LoadImage("assets/title_screen.png"),
		FadeAlpha:           1.0,
		TitleYOffset:        -40, // move everything up 40px
	}

	g.initAbilities()

	tt, err := opentype.Parse(LoadBytes("assets/fonts/title.ttf"))
	if err != nil {
		panic(err)
	}

	g.TitleFontSize = 48 // big chunky text
	g.TitleFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: g.TitleFontSize,
		DPI:  72,
	})
	if err != nil {
		panic(err)
	}

	return g
}

func (g *Game) startGameplay() {
	g.Enemies = []*Enemy{}
	g.Projectiles = []*Projectile{}
	g.FireTiles = []*FireTile{}
	g.Crystals = []*Crystal{}
	g.DamageNumbers = []*DamageNumber{}

	g.waveTimer = 3.0
	g.waveNumber = 1
}

func (g *Game) CameraOffset() (float64, float64) {
	return g.player.Pos.X - 400, g.player.Pos.Y - 300
}

func (g *Game) Update() error {
	now := time.Now()
	dt := now.Sub(g.lastTime).Seconds()
	g.lastTime = now

	switch g.GameState {

	// ----------------------------
	// TITLE SCREEN UPDATE
	// ----------------------------
	case StateTitle:
		return g.updateTitle(dt)

	// ----------------------------
	// TRAINING ROOM UPDATE
	// ----------------------------
	case StateTraining:
		g.updateToggleUI()
		return g.updateTraining(dt)

	// ----------------------------
	// GAMEPLAY UPDATE
	// ----------------------------
	case StateGameplay:
		g.bg = g.grassBg
		// If level-up menu is open, pause gameplay and handle input
		if g.LevelUpMenuOpen {
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

	return nil
}

func (g *Game) drawGameplay(screen *ebiten.Image) {
	camX, camY := g.CameraOffset()

	// ----------------------------
	// Draw Background (Tiled)
	// ----------------------------
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

	// Draw Player
	g.player.DrawWithCamera(screen, camX, camY)

	// Draw Enemies
	for _, e := range g.Enemies {
		if !e.Alive || e.Type == nil || e.Type.Sprite == nil {
			continue
		}

		w, h := e.Type.Sprite.Size()
		scale := e.Type.Scale * (e.Radius * 2) / float64(h)

		px := e.Pos.X - camX - float64(w)*scale/2
		py := e.Pos.Y - camY - float64(h)*scale/2

		px = math.Round(px)
		py = math.Round(py)

		outline := &ebiten.DrawImageOptions{}
		outline.Filter = ebiten.FilterNearest
		outline.GeoM.Scale(scale, scale)
		outline.GeoM.Translate(px, py)
		outline.ColorM.Scale(0, 0, 0, 0.6)
		screen.DrawImage(e.Type.Sprite, outline)

		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterNearest
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(px, py)

		if e.HitFlash > 0 {
			op.ColorM.Scale(2, 2, 2, 1)
		}

		screen.DrawImage(e.Type.Sprite, op)
	}

	// Draw Crystals
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

	// Draw Projectiles
	for _, p := range g.Projectiles {
		if !p.Alive || p.Sprite == nil {
			continue
		}

		op := &ebiten.DrawImageOptions{}

		if p.Orbit {
			x := g.player.Pos.X + math.Cos(p.Angle)*p.Radius - camX
			y := g.player.Pos.Y + math.Sin(p.Angle)*p.Radius - camY
			op.GeoM.Translate(x, y)
			screen.DrawImage(p.Sprite, op)
			continue
		}

		px := p.Pos.X - camX
		py := p.Pos.Y - camY

		px = math.Round(px)
		py = math.Round(py)

		op.GeoM.Translate(px, py)
		screen.DrawImage(p.Sprite, op)
	}

	// Draw Fire Tiles
	for _, ft := range g.FireTiles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(ft.Pos.X-camX, ft.Pos.Y-camY)
		screen.DrawImage(PlaceholderIcon(255, 120, 0), op)
	}

	// Draw Damage Numbers
	for _, dn := range g.DamageNumbers {
		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("%d", dn.Value),
			int(dn.X-camX),
			int(dn.Y-camY),
		)
	}

	// UI Bars
	g.drawHPBar(screen)
	g.drawXPBar(screen)

	// Level Up Menu
	if g.LevelUpMenuOpen {
		g.drawLevelUpMenu(screen)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.GameState {

	// ----------------------------
	// TITLE SCREEN
	// ----------------------------
	case StateTitle:
		g.drawTitle(screen)
		return

	// ----------------------------
	// TRAINING ROOM
	// ----------------------------
	case StateTraining:
		g.drawTrainingRoom(screen)
		return

	// ----------------------------
	// GAMEPLAY
	// ----------------------------
	case StateGameplay:
		g.drawGameplay(screen)
		return
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

func (g *Game) checkLevelUp() {
	xpNeeded := float64(g.player.Level * 10)

	if g.player.XP >= xpNeeded {
		g.player.XP -= xpNeeded
		g.player.Level++
		g.openLevelUpMenu()
	}
}

func (g *Game) updateSpells(dt float64) {

	// ----------------------------------------
	// 1. DAGGER AUTO-FIRE (still valid)
	// ----------------------------------------
	if g.player.HasDagger {
		g.DaggerCooldown -= dt
		if g.DaggerCooldown <= 0 {
			g.fireDagger() // now spawns a Projectile
			g.DaggerCooldown = g.DaggerRate
		}
	}

	// ----------------------------------------
	// 2. PROJECTILE UPDATE LOOP (Unified)
	// ----------------------------------------
	for _, p := range g.Projectiles {
		if !p.Alive {
			continue
		}

		// ----------------------------
		// Orbiting projectiles
		// ----------------------------
		if p.Orbit {
			p.Angle += p.Speed * dt
			p.Pos.X = g.player.Pos.X + math.Cos(p.Angle)*p.Radius
			p.Pos.Y = g.player.Pos.Y + math.Sin(p.Angle)*p.Radius
		} else {

			// ----------------------------
			// Homing projectiles
			// ----------------------------
			if p.Homing {
				var target *Enemy
				minDist := 999999.0

				for _, e := range g.Enemies {
					if !e.Alive {
						continue
					}
					dx := e.Pos.X - p.Pos.X
					dy := e.Pos.Y - p.Pos.Y
					dist := math.Hypot(dx, dy)
					if dist < minDist {
						minDist = dist
						target = e
					}
				}

				if target != nil {
					dx := target.Pos.X - p.Pos.X
					dy := target.Pos.Y - p.Pos.Y
					d := math.Hypot(dx, dy)
					if d > 0 {
						p.Vel.X = (dx / d) * p.Speed
						p.Vel.Y = (dy / d) * p.Speed
					}
				}
			}

			// ----------------------------
			// Boomerang curve
			// ----------------------------
			if p.Curve != 0 {
				p.Vel.Y += p.Curve
			}

			// ----------------------------
			// Boomerang return
			// ----------------------------
			if p.Return && p.Lifetime == 20 { // halfway
				p.Vel.X *= -1
				p.Vel.Y *= -1
			}

			// ----------------------------
			// Movement
			// ----------------------------
			p.Pos.X += p.Vel.X * dt
			p.Pos.Y += p.Vel.Y * dt
		}

		// ----------------------------
		// Lifetime expiration
		// ----------------------------
		if p.Lifetime > 0 {
			p.Lifetime--
			if p.Lifetime <= 0 {
				p.Alive = false
				continue
			}
		}

		// ----------------------------
		// Collision with enemies
		// ----------------------------
		for _, e := range g.Enemies {
			if !e.Alive {
				continue
			}

			dx := e.Pos.X - p.Pos.X
			dy := e.Pos.Y - p.Pos.Y
			dist := math.Hypot(dx, dy)

			if dist < e.Radius+10 {
				e.HP -= p.Damage

				if p.OnHit != nil {
					p.OnHit(e)
				}

				// Non-orbiting projectiles disappear on hit
				if !p.Orbit {
					p.Alive = false
				}

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

	// ----------------------------------------
	// 3. Remove dead projectiles
	// ----------------------------------------
	for i := len(g.Projectiles) - 1; i >= 0; i-- {
		if !g.Projectiles[i].Alive {
			g.Projectiles = append(g.Projectiles[:i], g.Projectiles[i+1:]...)
		}
	}

	// ----------------------------------------
	// 4. FIRE TILES (unchanged)
	// ----------------------------------------
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
