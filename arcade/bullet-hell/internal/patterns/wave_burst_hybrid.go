package patterns

import (
	"math"
	"math/rand"

	"github.com/KennyHoang-CS/portfolio/bullet-hell/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type WaveBurstHybridPattern struct {
	timer     int
	duration  int
	amplitude float64
	frequency float64
	speed     float64
}

// Draw implements [game.BulletPattern].
func (p *WaveBurstHybridPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewWaveBurstHybridPattern() *WaveBurstHybridPattern {
	return &WaveBurstHybridPattern{
		duration:  300,
		amplitude: 40,
		frequency: 0.06,
		speed:     1.6,
	}
}

func (p *WaveBurstHybridPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	cx := float64(g.ScreenWidth()) / 2
	cy := float64(g.ScreenHeight()) / 3

	if p.timer%6 == 0 {
		amp := p.amplitude + float64(difficulty)*3
		freq := p.frequency + float64(difficulty)*0.002
		speed := p.speed + float64(difficulty)*0.05

		angle := math.Sin(float64(p.timer) * freq)
		vx := angle * amp * 0.05
		vy := speed

		g.NewBullet(
			cx, cy,
			vx, vy,
			game.BulletOrb,
		)
	}

	if p.timer%45 == 0 {
		count := 10 + difficulty
		speed := 1.6 + float64(difficulty)*0.05

		for i := 0; i < count; i++ {
			a := rand.Float64() * 2 * math.Pi
			g.NewBullet(
				cx, cy,
				math.Cos(a)*speed,
				math.Sin(a)*speed,
				game.BulletStar,
			)
		}
	}
}

func (p *WaveBurstHybridPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *WaveBurstHybridPattern) Reset() {
	p.timer = 0
}

func (p *WaveBurstHybridPattern) Name() string { return "WaveBurstHybrid" }
