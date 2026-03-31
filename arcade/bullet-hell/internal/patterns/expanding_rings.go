package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type ExpandingRingsPattern struct {
	timer    int
	duration int
}

// Draw implements [game.BulletPattern].
func (p *ExpandingRingsPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewExpandingRingsPattern() *ExpandingRingsPattern {
	return &ExpandingRingsPattern{
		duration: 260,
	}
}

func (p *ExpandingRingsPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	if p.timer%30 == 0 {
		cx := float64(g.ScreenWidth()) / 2
		cy := float64(g.ScreenHeight()) / 3

		rings := 1 + difficulty/2
		baseSpeed := 1.4 + float64(difficulty)*0.05
		num := 14 + difficulty

		for r := 0; r < rings; r++ {
			speed := baseSpeed + float64(r)*0.2
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
}
func (p *ExpandingRingsPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *ExpandingRingsPattern) Reset() {
	p.timer = 0
}

func (p *ExpandingRingsPattern) Name() string { return "ExpandingRings" }
