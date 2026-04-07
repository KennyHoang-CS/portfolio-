package patterns

import (
	"math"
	"math/rand"

	"github.com/KennyHoang-CS/portfolio/bullet-hell/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type RandomSpreadBiasPattern struct {
	timer    int
	duration int
}

// Draw implements [game.BulletPattern].
func (p *RandomSpreadBiasPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewRandomSpreadBiasPattern() *RandomSpreadBiasPattern {
	return &RandomSpreadBiasPattern{
		duration: 260,
	}
}

func (p *RandomSpreadBiasPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	if p.timer%8 != 0 {
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
	baseAngle := math.Atan2(dy, dx)

	count := 6 + difficulty
	speed := 1.6 + float64(difficulty)*0.05
	bias := 0.4 + float64(difficulty)*0.05

	for i := 0; i < count; i++ {
		angle := baseAngle + (rand.Float64()-0.5)*bias
		g.NewBullet(
			cx, cy,
			math.Cos(angle)*speed,
			math.Sin(angle)*speed,
			game.BulletKunai,
		)
	}
}

func (p *RandomSpreadBiasPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *RandomSpreadBiasPattern) Reset() {
	p.timer = 0
}

func (p *RandomSpreadBiasPattern) Name() string { return "RandomSpreadBias" }
