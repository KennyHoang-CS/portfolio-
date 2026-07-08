package game

import (
    "fmt"
    "math"
    "math/rand"
    "time"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type DamageNumber struct {
    X, Y  float64
    Value int
    Life  float64
}

type Game struct {
    player        *Player
    lastTime      time.Time
    Enemies       []*Enemy
    DamageNumbers []*DamageNumber

    waveTimer   float64
    waveNumber  int

    bg *ebiten.Image
}

func NewGame() *Game {
    InitMonsters()

    g := &Game{
        player:     NewPlayer(),
        lastTime:   time.Now(),
        waveTimer:  3.0,
        waveNumber: 1,
        bg:         LoadImage("assets/grass_tile.png"),
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

    g.handleCombat(dt)
    g.updateDamageNumbers(dt)
    g.updateWaves(dt)

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

    // Draw enemies with outline + scale
    for _, e := range g.Enemies {
        if !e.Alive || e.Type == nil || e.Type.Sprite == nil {
            continue
        }

        w, h := e.Type.Sprite.Size()

        // NEW: per‑monster scale
        scale := e.Type.Scale * (e.Radius * 2) / float64(h)

        // Outline pass
        outline := &ebiten.DrawImageOptions{}
        outline.GeoM.Scale(scale, scale)
        outline.GeoM.Translate(
            e.Pos.X - camX - float64(w)*scale/2,
            e.Pos.Y - camY - float64(h)*scale/2,
        )
        outline.ColorM.Scale(0, 0, 0, 0.6)
        screen.DrawImage(e.Type.Sprite, outline)

        // Main sprite
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Scale(scale, scale)
        op.GeoM.Translate(
            e.Pos.X - camX - float64(w)*scale/2,
            e.Pos.Y - camY - float64(h)*scale/2,
        )

        if e.HitFlash > 0 {
            op.ColorM.Scale(2, 2, 2, 1)
        }

        screen.DrawImage(e.Type.Sprite, op)
    }

    for _, dn := range g.DamageNumbers {
        ebitenutil.DebugPrintAt(
            screen,
            fmt.Sprintf("%d", dn.Value),
            int(dn.X - camX),
            int(dn.Y - camY),
        )
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
            }
        }
    }
}

func (g *Game) updateDamageNumbers(dt float64) {
    for i := len(g.DamageNumbers)-1; i >= 0; i-- {
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
