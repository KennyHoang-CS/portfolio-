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
}

func NewGame() *Game {
    InitMonsters()

    g := &Game{
        player:        NewPlayer(),
        lastTime:      time.Now(),
        waveTimer:     3.0,
        waveNumber:    1,
        bg:            LoadImage("assets/grass_tile.png"),
        crystalSprite: LoadImage("assets/crystal.png"),
        Crystals:      []*Crystal{},
    }

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

func (g *Game) CameraOffset() (float64, float64) {
    return g.player.Pos.X - 400, g.player.Pos.Y - 300
}

func (g *Game) Update() error {
    now := time.Now()
    dt := now.Sub(g.lastTime).Seconds()
    g.lastTime = now

    g.player.Update(dt)

    for _, e := range g.Enemies {
        e.Update(dt, g.player.Pos)
    }

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
        scale := e.Type.Scale * (e.Radius * 2) / float64(h)

        outline := &ebiten.DrawImageOptions{}
        outline.GeoM.Scale(scale, scale)
        outline.GeoM.Translate(
            e.Pos.X-camX-float64(w)*scale/2,
            e.Pos.Y-camY-float64(h)*scale/2,
        )
        outline.ColorM.Scale(0, 0, 0, 0.6)
        screen.DrawImage(e.Type.Sprite, outline)

        op := &ebiten.DrawImageOptions{}
        op.GeoM.Scale(scale, scale)
        op.GeoM.Translate(
            e.Pos.X-camX-float64(w)*scale/2,
            e.Pos.Y-camY-float64(h)*scale/2,
        )

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
            c.Pos.X - camX - float64(w)*scale/2,
            c.Pos.Y - camY - float64(h)*scale/2 + bob,
        )

        screen.DrawImage(c.Sprite, op)
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
    slashRadius := 40.0

    for _, e := range g.Enemies {
        if !e.Alive {
            continue
        }

        dx := e.Pos.X - slashX
        dy := e.Pos.Y - slashY
        dist := math.Hypot(dx, dy)

        if dist < slashRadius+e.Radius {
            e.HP -= 5
            e.HitFlash = 0.1

            kb := 120.0
            e.Pos.X += slashDir * kb * dt

            g.DamageNumbers = append(g.DamageNumbers, &DamageNumber{
                X:     e.Pos.X,
                Y:     e.Pos.Y - 20,
                Value: 5,
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
    }
}

func (g *Game) drawHPBar(screen *ebiten.Image) {
    camX, camY := g.CameraOffset()

    barWidth := float32(50)
    barHeight := float32(9)

    // Position the bar slightly below the gopher
    x := float32(g.player.Pos.X - camX - float64(barWidth)/2)
    y := float32(g.player.Pos.Y - camY + 45) // adjust 45 to match your sprite height

    hpPercent := float32(g.player.HP / g.player.MaxHP)

    // Background (dark red)
    vector.FillRect(screen, x, y, barWidth, barHeight, color.RGBA{80, 0, 0, 255}, false)

    // HP fill (bright red)
    vector.FillRect(screen, x, y, barWidth*hpPercent, barHeight, color.RGBA{255, 40, 40, 255}, false)
}

func (g *Game) drawXPBar(screen *ebiten.Image) {
    barWidth := float32(800)
    barHeight := float32(15)
    x := float32(0)
    y := float32(0)

    xpNeeded := float32(g.player.Level * 10)
    xpPercent := float32(g.player.XP) / xpNeeded

    vector.FillRect(screen, x, y, barWidth, barHeight, color.RGBA{30, 30, 60, 255}, false)
    vector.FillRect(screen, x, y, barWidth*xpPercent, barHeight, color.RGBA{80, 80, 255, 255}, false)
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

                // Push both enemies away from each other
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
        minDist := 20 + e.Radius // player radius ~20

        if dist < minDist && dist > 0 {
            overlap := minDist - dist
            pushX := dx / dist * overlap
            pushY := dy / dist * overlap

            // Push player away from enemy
            g.player.Pos.X += pushX
            g.player.Pos.Y += pushY
        }
    }
}
