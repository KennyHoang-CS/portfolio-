package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/bullet-hell/internal/game"
)

type SpiralPattern struct {
	timer    int
	angle    float64
	duration int
}

func NewSpiralPattern() *SpiralPattern {
	return &SpiralPattern{
		duration: 240,
	}
}

func (p *SpiralPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	p.angle += 0.08 + float64(difficulty)*0.01

	if p.timer%3 == 0 {
		cx := float64(g.ScreenWidth()) / 2
		cy := float64(g.ScreenHeight()) / 3

		speed := 1.8 + float64(difficulty)*0.05

		g.NewBullet(
			cx, cy,
			math.Cos(p.angle)*speed,
			math.Sin(p.angle)*speed,
			game.BulletStar,
		)
	}
}

func (p *SpiralPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *SpiralPattern) Reset() {
	p.timer = 0
	p.angle = 0
}

func (p *SpiralPattern) Name() string { return "Spiral" }
