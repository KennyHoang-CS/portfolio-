package patterns

import (
	"math"
	"math/rand"

	"github.com/KennyHoang-CS/portfolio/bullet-hell/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type ChaoticWavesPattern struct {
	timer     int
	duration  int
	amplitude float64
	frequency float64
	speed     float64
}

// Draw implements [game.BulletPattern].
func (p *ChaoticWavesPattern) Draw(screen *ebiten.Image, g *game.Game) {
}

func NewChaoticWavesPattern() *ChaoticWavesPattern {
	return &ChaoticWavesPattern{
		duration:  280,
		amplitude: 40,
		frequency: 0.06,
		speed:     1.6,
	}
}

func (p *ChaoticWavesPattern) Update(g *game.Game, difficulty int) {
	p.timer++

	if p.timer%6 == 0 {
		cx := float64(g.ScreenWidth()) / 2
		cy := float64(g.ScreenHeight()) / 3

		ampJitter := (rand.Float64() - 0.5) * 10
		freqJitter := (rand.Float64() - 0.5) * 0.01

		amp := p.amplitude + float64(difficulty)*3 + ampJitter
		freq := p.frequency + float64(difficulty)*0.002 + freqJitter
		speed := p.speed + float64(difficulty)*0.05

		angle := math.Sin(float64(p.timer) * freq)
		vx := angle * amp * 0.05
		vy := speed

		g.NewBullet(cx, cy, vx, vy, game.BulletOrb)
	}
}

func (p *ChaoticWavesPattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *ChaoticWavesPattern) Reset() {
	p.timer = 0
}

func (p *ChaoticWavesPattern) Name() string { return "ChaoticWaves" }
