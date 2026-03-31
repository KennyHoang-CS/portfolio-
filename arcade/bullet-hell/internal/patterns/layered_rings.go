package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type LayeredRingsPattern struct {
	timer    int
	duration int
}

// Draw implements [game.BulletPattern].
func (p *LayeredRingsPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewLayeredRingsPattern() *LayeredRingsPattern {
	return &LayeredRingsPattern{
		duration: 260,
	}
}

func (p *LayeredRingsPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	if p.timer%28 != 0 {
		return
	}

	cx := float64(g.ScreenWidth()) / 2
	cy := float64(g.ScreenHeight()) / 3

	layers := 2 + difficulty/2
	num := 12 + difficulty
	baseSpeed := 1.4 + float64(difficulty)*0.05

	for l := 0; l < layers; l++ {
		speed := baseSpeed + float64(l)*0.2
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

func (p *LayeredRingsPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *LayeredRingsPattern) Reset() {
	p.timer = 0
}

func (p *LayeredRingsPattern) Name() string { return "LayeredRings" }
