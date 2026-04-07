package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/bullet-hell/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpiralAimedHybridPattern struct {
	timer    int
	duration int
	angle    float64
}

// Draw implements [game.BulletPattern].
func (p *SpiralAimedHybridPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewSpiralAimedHybridPattern() *SpiralAimedHybridPattern {
	return &SpiralAimedHybridPattern{
		duration: 300,
	}
}

func (p *SpiralAimedHybridPattern) Update(g *game.Game, difficulty int) {
	p.timer++
	p.angle += 0.07 + float64(difficulty)*0.01

	cx := float64(g.ScreenWidth()) / 2
	cy := float64(g.ScreenHeight()) / 3

	if p.timer%3 == 0 {
		speed := 1.7 + float64(difficulty)*0.05
		g.NewBullet(
			cx, cy,
			math.Cos(p.angle)*speed,
			math.Sin(p.angle)*speed,
			game.BulletStar,
		)
	}

	if p.timer%40 == 0 {
		player := g.GetPlayer()
		if player != nil {
			dx := player.Position.X - cx
			dy := player.Position.Y - cy
			dist := math.Hypot(dx, dy)
			speed := 1.8 + float64(difficulty)*0.05
			g.NewBullet(
				cx, cy,
				dx/dist*speed,
				dy/dist*speed,
				game.BulletArrow,
			)
		}
	}
}
func (p *SpiralAimedHybridPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *SpiralAimedHybridPattern) Reset() {
	p.timer = 0
	p.angle = 0
}

func (p *SpiralAimedHybridPattern) Name() string { return "SpiralAimedHybrid" }
