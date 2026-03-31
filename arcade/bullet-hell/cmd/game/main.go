package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/KennyHoang-CS/portfolio/internal/game"
	"github.com/KennyHoang-CS/portfolio/internal/patterns"
)

func main() {
	pm := patterns.NewPatternManager()

	// Warm-up
	pm.Add(patterns.NewRotatingFlowerPattern())
	pm.Add(patterns.NewExpandingRingsPattern())
	pm.Add(patterns.NewRadialPattern())

	// Mid-boss
	pm.Add(patterns.NewDoubleSpiralPattern())
	pm.Add(patterns.NewWavesPattern())
	pm.Add(patterns.NewSpiralRainPattern())

	// Pressure
	pm.Add(patterns.NewAimedPattern())
	pm.Add(patterns.NewCurvedAimedPattern())
	pm.Add(patterns.NewRandomSpreadBiasPattern())

	// Spell Card 1
	pm.Add(patterns.NewButterflyPattern())
	pm.Add(patterns.NewLayeredRingsPattern())
	pm.Add(patterns.NewRotatingEmittersPattern())

	// Spell Card 2
	pm.Add(patterns.NewChaoticWavesPattern())
	pm.Add(patterns.NewRandomBurstsSafeZonesPattern())
	pm.Add(patterns.NewWaveBurstHybridPattern())

	// Final Spell
	pm.Add(patterns.NewSpiralAimedHybridPattern())
	pm.Add(patterns.NewBurstPattern())
	pm.Add(patterns.NewRainfallPattern())
	pm.Add(patterns.NewUltimateDetonation())
	pm.Add(patterns.NewUltimateDetonation())

	pm.Shuffle()

	g := game.NewGame(pm)

	// ⭐ Load all bullet sprites, player sprites, etc.
	if err := game.LoadAssets(); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Touhou-Style Bullet Hell")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
