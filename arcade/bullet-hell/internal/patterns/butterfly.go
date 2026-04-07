package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/bullet-hell/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type ButterflyPattern struct {
	timer    int
	duration int
	angle    float64
}

// Draw implements [game.BulletPattern].
func (p *ButterflyPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewButterflyPattern() *ButterflyPattern {
	return &ButterflyPattern{
		duration: 280,
	}
}

func (p *ButterflyPattern) Update(g *game.Game, difficulty int) {
	p.timer++
	p.angle += 0.06 + float64(difficulty)*0.008

	if p.timer%3 == 0 {
		cx := float64(g.ScreenWidth()) / 2
		cy := float64(g.ScreenHeight()) / 3

		speed := 1.7 + float64(difficulty)*0.05
		offset := 0.6 + float64(difficulty)*0.03

		a1 := p.angle + offset
		a2 := p.angle - offset

		g.NewBullet(cx, cy, math.Cos(a1)*speed, math.Sin(a1)*speed, game.BulletPetal)
		g.NewBullet(cx, cy, math.Cos(a2)*speed, math.Sin(a2)*speed, game.BulletPetal)
	}
}

func (p *ButterflyPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *ButterflyPattern) Reset() {
	p.timer = 0
	p.angle = 0
}

func (p *ButterflyPattern) Name() string { return "Butterfly" }
