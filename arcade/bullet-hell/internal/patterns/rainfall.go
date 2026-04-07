package patterns

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/KennyHoang-CS/portfolio/bullet-hell/internal/game"
)

type RainfallPattern struct {
	timer    int
	duration int

	// wind + gusts
	windOffset float64
	windTimer  int

	// lightning flash
	flashTimer int
}

func NewRainfallPattern() *RainfallPattern {
	return &RainfallPattern{
		duration: 1200, // ~20 seconds
	}
}

func (p *RainfallPattern) Update(g *game.Game, difficulty int) {
	p.timer++
	p.windTimer++

	// -----------------------------
	// WIND GUSTS
	// -----------------------------
	if p.windTimer%180 == 0 { // every 3 seconds
		p.windOffset = (rand.Float64()*2 - 1) * 0.8 // -0.8 to +0.8
	}

	// -----------------------------
	// LIGHTNING FLASH
	// -----------------------------
	if p.timer%240 == 0 { // every 4 seconds
		p.flashTimer = 8 // flash lasts 8 frames
	}

	// -----------------------------
	// SPAWN RATE
	// -----------------------------
	spawnRate := 12 - difficulty
	if spawnRate < 3 {
		spawnRate = 3
	}

	if p.timer%spawnRate == 0 {
		x := rand.Float64() * float64(g.ScreenWidth())
		y := float64(-10)

		// base fall speed
		speed := 2.0 + float64(difficulty)*0.3

		// slight wind drift
		dx := p.windOffset + (rand.Float64()*0.4 - 0.2)

		// HEAVY DROP (rare)
		if rand.Float64() < 0.1 {
			speed *= 1.6
			g.NewBullet(x, y, dx, speed, game.BulletStar)
		} else {
			g.NewBullet(x, y, dx, speed, game.BulletOrb)
		}
	}

	// -----------------------------
	// SPLASH SPARKS (bottom hits)
	// -----------------------------
	for _, e := range g.World().Entities() {
		if e.Bullet == nil {
			continue
		}

		if e.Position.Y > float64(g.ScreenHeight()+20) {
			// splash sparks
			for i := 0; i < 4; i++ {
				s := g.World().NewEntity()
				s.Position = &game.Position{
					X: e.Position.X,
					Y: float64(g.ScreenHeight()) - 4,
				}
				s.Spark = &game.SparkTag{Life: 12}
			}

			e.Destroy = true
		}
	}
}

func (p *RainfallPattern) Draw(screen *ebiten.Image, g *game.Game) {
	// -------------------------------------
	// CLOUD BAND AT TOP (dark gradient)
	// -------------------------------------
	for i := 0; i < 80; i++ {
		alpha := uint8(80 - i)
		vector.FillRect(
			screen,
			0, float32(i),
			float32(g.ScreenWidth()), 1,
			color.RGBA{0, 0, 0, alpha},
			false,
		)
	}

	// -------------------------------------
	// SOFT, VISIBLE, SAFE LIGHTNING TINT
	// -------------------------------------
	if p.flashTimer > 0 {
		p.flashTimer--
		alpha := uint8(2 + p.flashTimer*2)
		// a pale gray-blue that shows up on dark backgrounds
		overlay := ebiten.NewImage(g.ScreenWidth(), g.ScreenHeight())
		overlay.Fill(color.RGBA{160, 180, 220, alpha})

		screen.DrawImage(overlay, nil)
	}
}
func (p *RainfallPattern) IsFinished() bool { return p.timer > p.duration }
func (p *RainfallPattern) Reset()           { p.timer = 0 }
func (p *RainfallPattern) Name() string     { return "RainfallPattern" }
