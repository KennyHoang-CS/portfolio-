package patterns

import (
	"math"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type CurvedAimedPattern struct {
	timer    int
	duration int
}

// Draw implements [game.BulletPattern].
func (p *CurvedAimedPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewCurvedAimedPattern() *CurvedAimedPattern {
	return &CurvedAimedPattern{
		duration: 280,
	}
}

func (p *CurvedAimedPattern) Update(g *game.Game, difficulty int) {
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
	baseAngle := math.Atan2(dy, dx)
	speed := 1.7 + float64(difficulty)*0.05

	// fire a small fan that will curve
	count := 3 + difficulty/2
	curve := 0.04 + float64(difficulty)*0.005

	for i := -count; i <= count; i++ {
		angle := baseAngle + float64(i)*0.08
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed

		e := g.NewBullet(cx, cy, vx, vy, game.BulletAmulet)
		e.Bullet.Curve = curve * float64(i)
	}
}

func (p *CurvedAimedPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *CurvedAimedPattern) Reset() {
	p.timer = 0
}

func (p *CurvedAimedPattern) Name() string { return "CurvedAimed" }
