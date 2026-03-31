package patterns

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type UltimateDetonation struct {
	timer    int
	duration int

	// orb state
	x, y      float64
	vx, vy    float64
	phase     int // 0 = spawn, 1 = drift, 2 = charge, 3 = explode
	phaseTime int
}

func NewUltimateDetonation() *UltimateDetonation {
	return &UltimateDetonation{
		duration: 600,
	}
}

func (p *UltimateDetonation) Update(g *game.Game, difficulty int) {
	p.timer++
	p.phaseTime++

	switch p.phase {

	// ---------------------------------------------------
	// 0. Spawn at random position
	// ---------------------------------------------------
	case 0:
		p.x = rand.Float64() * float64(g.ScreenWidth())
		p.y = rand.Float64() * float64(g.ScreenHeight()/2)

		// random drift velocity
		p.vx = (rand.Float64()*2 - 1) * 0.8
		p.vy = (rand.Float64()*2 - 1) * 0.8

		p.phase = 1
		p.phaseTime = 0

	// ---------------------------------------------------
	// 1. Drift around for a bit
	// ---------------------------------------------------
	case 1:
		p.x += p.vx
		p.y += p.vy

		// bounce off screen edges
		if p.x < 40 || p.x > float64(g.ScreenWidth()-40) {
			p.vx *= -1
		}
		if p.y < 40 || p.y > float64(g.ScreenHeight()/2) {
			p.vy *= -1
		}

		// drift for ~60 frames
		if p.phaseTime > 60 {
			p.phase = 2
			p.phaseTime = 0
		}

	// ---------------------------------------------------
	// 2. Charge-up (pulsing + shrinking warning ring)
	// ---------------------------------------------------
	case 2:
		// You can draw a pulsing circle in your Draw() using p.phaseTime
		// You can also draw a shrinking ring based on (40 - p.phaseTime)

		if p.phaseTime > 40 {
			p.phase = 3
			p.phaseTime = 0
		}

	// ---------------------------------------------------
	// 3. Multi-stage explosion
	// ---------------------------------------------------
	case 3:
		// Stage 1: inner ring
		innerCount := 12 + difficulty*2
		innerSpeed := 1.8 + float64(difficulty)*0.05

		for i := 0; i < innerCount; i++ {
			angle := 2 * math.Pi * float64(i) / float64(innerCount)
			g.NewBullet(
				p.x, p.y,
				math.Cos(angle)*innerSpeed,
				math.Sin(angle)*innerSpeed,
				game.BulletOrb,
			)
		}

		// Stage 2: outer ring (offset by half-angle)
		outerCount := 16 + difficulty*2
		outerSpeed := 2.4 + float64(difficulty)*0.05

		for i := 0; i < outerCount; i++ {
			angle := 2*math.Pi*float64(i)/float64(outerCount) + math.Pi/float64(outerCount)
			g.NewBullet(
				p.x, p.y,
				math.Cos(angle)*outerSpeed,
				math.Sin(angle)*outerSpeed,
				game.BulletStar, // different sprite for visual flair
			)
		}

		// reset for next orb
		p.phase = 0
		p.phaseTime = 0
	}
}

func (p *UltimateDetonation) IsFinished() bool {
	return p.timer > p.duration
}

func (p *UltimateDetonation) Reset() {
	p.timer = 0
	p.phase = 0
	p.phaseTime = 0
}

func (p *UltimateDetonation) Name() string { return "UltimateDetonation" }

func (p *UltimateDetonation) X() float64     { return p.x }
func (p *UltimateDetonation) Y() float64     { return p.y }
func (p *UltimateDetonation) Phase() int     { return p.phase }
func (p *UltimateDetonation) PhaseTime() int { return p.phaseTime }

func (p *UltimateDetonation) Draw(screen *ebiten.Image, g *game.Game) {
    // -----------------------------------------
    // Blinking red orb (phase 1 and 2)
    // -----------------------------------------
    if p.phase == 1 || p.phase == 2 {
        // blink every 6 frames
        blink := (p.phaseTime/6)%2 == 0

        var orbColor color.RGBA
        if blink {
            orbColor = color.RGBA{255, 40, 40, 255} // bright danger red
        } else {
            orbColor = color.RGBA{255, 120, 120, 255} // softer red
        }

        // base orb glow
        g.DrawGlowCircle(screen, p.x, p.y, 10, orbColor)
    }

    // -----------------------------------------
    // Charge-up extras (phase 2 only)
    // -----------------------------------------
    if p.phase == 2 {
        // pulsing radius
        pulse := 10 + 4*math.Sin(float64(p.phaseTime)/4)
        g.DrawGlowCircle(screen, p.x, p.y, pulse, color.RGBA{255, 160, 200, 255})

        // shrinking warning ring
        ring := 40 - float64(p.phaseTime)
        if ring > 0 {
            g.DrawRing(screen, p.x, p.y, ring, 2, color.RGBA{255, 255, 255, 200})
        }
    }

    // -----------------------------------------
    // Explosion flash (phase 3)
    // -----------------------------------------
    if p.phase == 3 && p.phaseTime < 6 {
        flash := 20 + float64(p.phaseTime)*4
        g.DrawGlowCircle(screen, p.x, p.y, flash, color.RGBA{255, 200, 200, 255})
    }
}