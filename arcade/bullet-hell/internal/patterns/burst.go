package patterns

import (
	"math"
	"math/rand"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type BurstPattern struct {
	timer    int
	duration int
}

// Draw implements [game.BulletPattern].
func (p *BurstPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewBurstPattern() *BurstPattern {
	return &BurstPattern{
		duration: 220,
	}
}

func (p *BurstPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	if p.timer%12 == 0 {
		cx := float64(g.ScreenWidth()) / 2
		cy := float64(g.ScreenHeight()) / 3

		count := 10 + difficulty // density increases
		speed := 1.6 + float64(difficulty)*0.05

		for i := 0; i < count; i++ {
			angle := rand.Float64() * (2 * math.Pi)
			g.NewBullet(
				cx, cy,
				math.Cos(angle)*speed,
				math.Sin(angle)*speed,
				game.BulletOrb,
			)
		}
	}
}

func (p *BurstPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *BurstPattern) Reset() {
	p.timer = 0
}

func (p *BurstPattern) Name() string { return "Burst" }
