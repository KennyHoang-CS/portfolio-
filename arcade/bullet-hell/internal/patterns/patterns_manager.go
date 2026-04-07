package patterns

import (
	"math/rand"

	"github.com/KennyHoang-CS/portfolio/bullet-hell/internal/game"
)

type PatternManager struct {
	patterns   []game.BulletPattern
	index      int
	difficulty int
}

func NewPatternManager() *PatternManager {
	return &PatternManager{
		patterns:   []game.BulletPattern{},
		index:      0,
		difficulty: 1,
	}
}

func (pm *PatternManager) Done() bool {
	return pm.index >= len(pm.patterns)
}

func (pm *PatternManager) Add(p game.BulletPattern) {
	pm.patterns = append(pm.patterns, p)
}

func (pm *PatternManager) Update(g *game.Game) {
	// No patterns? Nothing to do.
	if len(pm.patterns) == 0 {
		return
	}

	// If we've already used all patterns, stop updating.
	// This allows pm.Done() to return true.
	if pm.index >= len(pm.patterns) {
		return
	}

	// Get the current pattern
	current := pm.patterns[pm.index]

	// Update the pattern
	current.Update(g, pm.difficulty)

	// If the pattern is finished, move to the next one
	if current.IsFinished() {
		pm.index++

		// Reward for clearing a pattern
		g.Score += 500 * pm.difficulty

		// ⭐ NO LOOPING HERE — we simply stop at the end.
		// pm.index will eventually reach len(pm.patterns),
		// and pm.Done() will return true.
	}
}

func (pm *PatternManager) Reset() {
	pm.index = 0
	pm.difficulty = 1

	for _, p := range pm.patterns {
		p.Reset()
	}
}

func (pm *PatternManager) CurrentName() string {
	if len(pm.patterns) == 0 {
		return ""
	}
	return pm.patterns[pm.index].Name()
}

func (pm *PatternManager) Difficulty() int {
	return pm.difficulty
}

func (pm *PatternManager) Current() game.BulletPattern {
	if len(pm.patterns) == 0 {
		return nil
	}
	return pm.patterns[pm.index]
}

func (pm *PatternManager) Shuffle() {
	rand.Shuffle(len(pm.patterns), func(i, j int) {
		pm.patterns[i], pm.patterns[j] = pm.patterns[j], pm.patterns[i]
	})
}
