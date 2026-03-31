package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type AimedPattern struct {
	timer    int
	duration int
}

// Draw implements [game.BulletPattern].
func (p *AimedPattern) Draw(screen *ebiten.Image, g *game.Game) {
	
}

func NewAimedPattern() *AimedPattern {
	return &AimedPattern{
		duration: 300,
	}
}

func (p *AimedPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	if p.timer%35 != 0 {
		return
	}

	player := g.GetPlayer()
	if player == nil {
		return
	}

	cx := float64(g.ScreenWidth()) / 2
	cy := float64(g.ScreenHeight()) / 3

	dx := player.Position.X - cx
	dy := player.Position.Y - cy
	dist := math.Hypot(dx, dy)

	speed := 1.7 + float64(difficulty)*0.05

	g.NewBullet(
		cx, cy,
		dx/dist*speed,
		dy/dist*speed,
		game.BulletArrow,
	)
}

func (p *AimedPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *AimedPattern) Reset() {
	p.timer = 0
}

func (p *AimedPattern) Name() string { return "Aimed" }
