package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type RotatingEmittersPattern struct {
	timer     int
	duration  int
	baseAngle float64
}

// Draw implements [game.BulletPattern].
func (p *RotatingEmittersPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewRotatingEmittersPattern() *RotatingEmittersPattern {
	return &RotatingEmittersPattern{
		duration: 300,
	}
}

func (p *RotatingEmittersPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	cx := float64(g.ScreenWidth()) / 2
	cy := float64(g.ScreenHeight()) / 3

	emitters := 4 + difficulty/2
	radius := 60.0
	p.baseAngle += 0.02 + float64(difficulty)*0.005

	if p.timer%6 == 0 {
		for i := 0; i < emitters; i++ {
			angle := p.baseAngle + 2*math.Pi*float64(i)/float64(emitters)
			ex := cx + math.Cos(angle)*radius
			ey := cy + math.Sin(angle)*radius

			speed := 1.6 + float64(difficulty)*0.05
			bAngle := angle + math.Pi/2

			g.NewBullet(
				ex, ey,
				math.Cos(bAngle)*speed,
				math.Sin(bAngle)*speed,
				game.BulletStar,
			)
		}
	}
}

func (p *RotatingEmittersPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *RotatingEmittersPattern) Reset() {
	p.timer = 0
	p.baseAngle = 0
}

func (p *RotatingEmittersPattern) Name() string { return "RotatingEmitters" }
