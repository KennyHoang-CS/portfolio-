package game

import (
    "image"
    "math"

    "github.com/hajimehoshi/ebiten/v2"
)

var slashFrames []*ebiten.Image

type Player struct {
    Pos         Vec
    Speed       float64
    Idle        *ebiten.Image

    attacking     bool
    anticipation  float64
    recovery      float64

    attackTimer   float64
    attackFrame   int

    lastDir        int // -1 = left, +1 = right
    attackCooldown float64

    HP    float64
    MaxHP float64
    XP    float64
    Level int
}

func NewPlayer() *Player {
    idle := LoadImage("assets/gopher_idle.png")

    slashSheet := LoadImage("assets/slash.png")
    w, h := slashSheet.Size()

    frameWidth := w / 4
    frameHeight := h

    for i := 0; i < 4; i++ {
        x0 := i * frameWidth
        x1 := x0 + frameWidth
        frame := slashSheet.SubImage(image.Rect(x0, 0, x1, frameHeight)).(*ebiten.Image)
        slashFrames = append(slashFrames, frame)
    }

    return &Player{
        Pos:   Vec{X: 400, Y: 300},
        Speed: 200,
        Idle:  idle,

        lastDir: 1,

        HP:    100,
        MaxHP: 100,
        XP:    0,
        Level: 1,
    }
}

func (p *Player) Update(dt float64) {
    vx, vy := 0.0, 0.0

    if ebiten.IsKeyPressed(ebiten.KeyW) { vy -= 1 }
    if ebiten.IsKeyPressed(ebiten.KeyS) { vy += 1 }
    if ebiten.IsKeyPressed(ebiten.KeyA) { vx -= 1 }
    if ebiten.IsKeyPressed(ebiten.KeyD) { vx += 1 }

    if vx != 0 || vy != 0 {
        l := math.Hypot(vx, vy)
        vx /= l
        vy /= l
    }

    if vx < 0 {
        p.lastDir = -1
    } else if vx > 0 {
        p.lastDir = 1
    }

    p.Pos.X += vx * p.Speed * dt
    p.Pos.Y += vy * p.Speed * dt

    // --- AUTO ATTACK ---
    p.attackCooldown -= dt
    if p.attackCooldown <= 0 && !p.attacking {
        p.attacking = true
        p.anticipation = 0.05
        p.recovery = 0
        p.attackTimer = 0
        p.attackFrame = 0

        p.attackCooldown = 0.6
    }

    if p.anticipation > 0 {
        p.anticipation -= dt
        return
    }

    if p.attacking {
        p.attackTimer += dt

        if p.attackTimer > 0.08 {
            p.attackTimer = 0
            p.attackFrame++

            if p.attackFrame >= len(slashFrames) {
                p.attacking = false
                p.attackFrame = 0
                p.recovery = 0.05
            }
        }
    }

    if p.recovery > 0 {
        p.recovery -= dt
    }
}

func (p *Player) DrawWithCamera(screen *ebiten.Image, camX, camY float64) {

    // --- Draw gopher (flip left/right) ---
    {
        op := &ebiten.DrawImageOptions{}
        src := p.Idle

        srcW, srcH := src.Size()
        targetH := 80.0
        scale := targetH / float64(srcH)

        if p.lastDir == -1 {
            // Flip horizontally
            op.GeoM.Scale(-scale, scale)
            op.GeoM.Translate(float64(srcW)*scale/2, -float64(srcH)*scale/2)
        } else {
            op.GeoM.Scale(scale, scale)
            op.GeoM.Translate(-float64(srcW)*scale/2, -float64(srcH)*scale/2)
        }

        op.GeoM.Translate(p.Pos.X-camX, p.Pos.Y-camY)
        screen.DrawImage(src, op)
    }

    // --- Draw slash (flip + taller) ---
    if p.attacking {
        sop := &ebiten.DrawImageOptions{}
        sop.Filter = ebiten.FilterNearest

        src := slashFrames[p.attackFrame]
        srcW, srcH := src.Size()

        targetW := 150.0
        targetH := 80.0 // <-- increased height

        scaleX := targetW / float64(srcW)
        scaleY := targetH / float64(srcH)

        // Flip slash horizontally when facing left
        if p.lastDir == -1 {
            sop.GeoM.Scale(-scaleX, scaleY)
            sop.GeoM.Translate(float64(srcW)*scaleX/2, -float64(srcH)*scaleY/2)
        } else {
            sop.GeoM.Scale(scaleX, scaleY)
            sop.GeoM.Translate(-float64(srcW)*scaleX/2, -float64(srcH)*scaleY/2)
        }

        if p.anticipation > 0 {
            sop.GeoM.Translate(0, -10)
        }

        if p.recovery > 0 {
            sop.GeoM.Translate(0, 6)
        }

        slashOffsetX := float64(p.lastDir) * 150.0
        slashOffsetY := -5.0

        sop.GeoM.Translate(p.Pos.X+slashOffsetX-camX, p.Pos.Y+slashOffsetY-camY)
        screen.DrawImage(src, sop)
    }
}
