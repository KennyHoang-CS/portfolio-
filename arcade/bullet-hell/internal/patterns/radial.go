package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type RadialPattern struct {
	timer    int
	duration int
}

// Draw implements [game.BulletPattern].
func (p *RadialPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewRadialPattern() *RadialPattern {
	return &RadialPattern{
		duration: 180,
	}
}

func (p *RadialPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	if p.timer%25 == 0 {
		cx := float64(g.ScreenWidth()) / 2
		cy := float64(g.ScreenHeight()) / 3

		num := 12 + difficulty
		speed := 1.6 + float64(difficulty)*0.05

		for i := 0; i < num; i++ {
			angle := 2 * math.Pi * float64(i) / float64(num)
			g.NewBullet(
				cx, cy,
				math.Cos(angle)*speed,
				math.Sin(angle)*speed,
				game.BulletOrb,
			)
		}
	}
}

func (p *RadialPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *RadialPattern) Reset() {
	p.timer = 0
}

func (p *RadialPattern) Name() string { return "Radial" }
