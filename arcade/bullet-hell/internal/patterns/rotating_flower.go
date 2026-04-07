package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/bullet-hell/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type RotatingFlowerPattern struct {
	timer     int
	duration  int
	baseAngle float64
}

// Draw implements [game.BulletPattern].
func (p *RotatingFlowerPattern) Draw(screen *ebiten.Image, g *game.Game) {

}

func NewRotatingFlowerPattern() *RotatingFlowerPattern {
	return &RotatingFlowerPattern{
		duration: 240,
	}
}

func (p *RotatingFlowerPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	if p.timer%20 == 0 {
		cx := float64(g.ScreenWidth()) / 2
		cy := float64(g.ScreenHeight()) / 3

		petals := 8 + difficulty
		speed := 1.7 + float64(difficulty)*0.05
		p.baseAngle += 0.1 + float64(difficulty)*0.01

		for i := 0; i < petals; i++ {
			angle := p.baseAngle + 2*math.Pi*float64(i)/float64(petals)
			g.NewBullet(
				cx, cy,
				math.Cos(angle)*speed,
				math.Sin(angle)*speed,
				game.BulletPetal,
			)
		}
	}
}

func (p *RotatingFlowerPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *RotatingFlowerPattern) Reset() {
	p.timer = 0
	p.baseAngle = 0
}

func (p *RotatingFlowerPattern) Name() string { return "RotatingFlower" }
