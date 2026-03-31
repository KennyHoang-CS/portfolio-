package patterns

import (
	"math"
	"math/rand"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type RandomBurstsSafeZonesPattern struct {
	timer    int
	duration int
}

// Draw implements [game.BulletPattern].
func (p *RandomBurstsSafeZonesPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewRandomBurstsSafeZonesPattern() *RandomBurstsSafeZonesPattern {
	return &RandomBurstsSafeZonesPattern{
		duration: 260,
	}
}

func (p *RandomBurstsSafeZonesPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	if p.timer%18 != 0 {
		return
	}

	cx := float64(g.ScreenWidth()) / 2
	cy := float64(g.ScreenHeight()) / 3

	count := 24 + difficulty*2
	speed := 1.6 + float64(difficulty)*0.05

	safeAngle := rand.Float64() * 2 * math.Pi
	safeWidth := 0.5 - float64(difficulty)*0.02
	if safeWidth < 0.15 {
		safeWidth = 0.15
	}

	for i := 0; i < count; i++ {
		angle := 2 * math.Pi * float64(i) / float64(count)
		if math.Abs(angle-safeAngle) < safeWidth {
			continue
		}
		g.NewBullet(
			cx, cy,
			math.Cos(angle)*speed,
			math.Sin(angle)*speed,
			game.BulletOrb,
		)
	}
}

func (p *RandomBurstsSafeZonesPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *RandomBurstsSafeZonesPattern) Reset() {
	p.timer = 0
}

func (p *RandomBurstsSafeZonesPattern) Name() string { return "RandomBurstsSafeZones" }
