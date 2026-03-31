package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type DoubleSpiralPattern struct {
	timer    int
	duration int
	angle1   float64
	angle2   float64
}

// Draw implements [game.BulletPattern].
func (p *DoubleSpiralPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewDoubleSpiralPattern() *DoubleSpiralPattern {
	return &DoubleSpiralPattern{
		duration: 260,
	}
}

func (p *DoubleSpiralPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	rotSpeed := 0.08 + float64(difficulty)*0.01
	p.angle1 += rotSpeed
	p.angle2 -= rotSpeed

	if p.timer%3 == 0 {
		cx := float64(g.ScreenWidth()) / 2
		cy := float64(g.ScreenHeight()) / 3
		speed := 1.8 + float64(difficulty)*0.05

		g.NewBullet(cx, cy, math.Cos(p.angle1)*speed, math.Sin(p.angle1)*speed, game.BulletStar)
		g.NewBullet(cx, cy, math.Cos(p.angle2)*speed, math.Sin(p.angle2)*speed, game.BulletStar)
	}
}

func (p *DoubleSpiralPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *DoubleSpiralPattern) Reset() {
	p.timer = 0
	p.angle1 = 0
	p.angle2 = math.Pi
}

func (p *DoubleSpiralPattern) Name() string { return "DoubleSpiral" }
