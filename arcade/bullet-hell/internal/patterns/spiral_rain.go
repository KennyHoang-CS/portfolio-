package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpiralRainPattern struct {
	timer    int
	duration int
	angle    float64
}

// Draw implements [game.BulletPattern].
func (p *SpiralRainPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewSpiralRainPattern() *SpiralRainPattern {
	return &SpiralRainPattern{
		duration: 260,
	}
}

func (p *SpiralRainPattern) Update(g *game.Game, difficulty int) {
	p.timer++
	p.angle += 0.06 + float64(difficulty)*0.008

	if p.timer%4 == 0 {
		cx := float64(g.ScreenWidth()) / 2
		y := float64(g.ScreenHeight()) / 6

		speed := 1.5 + float64(difficulty)*0.05
		x := cx + math.Cos(p.angle)*120

		g.NewBullet(
			x, y,
			0,
			speed,
			game.BulletOrb,
		)
	}
}
func (p *SpiralRainPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *SpiralRainPattern) Reset() {
	p.timer = 0
	p.angle = 0
}

func (p *SpiralRainPattern) Name() string { return "SpiralRain" }
