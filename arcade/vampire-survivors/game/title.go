package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
)

func (g *Game) updateTitle(dt float64) error {
	g.BlinkTimer += dt

	if g.FadeAlpha > 0 {
		g.FadeAlpha -= dt * 0.5
		if g.FadeAlpha < 0 {
			g.FadeAlpha = 0
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.startGameplay()
		g.GameState = StateGameplay
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyT) {
		g.GameState = StateTraining
		g.loadTrainingRoom()
	}

	return nil
}

func (g *Game) updateTraining(dt float64) error {
    g.player.Update(dt)

    // Update dummy (even though Speed = 0)
    for _, e := range g.Enemies {
        e.Update(dt, g.player.Pos)
    }

    // All normal gameplay interactions
    g.resolveEnemyCollisions(dt)
    g.resolvePlayerCollision(dt)
    g.handleCombat(dt)
    g.handleEnemyContact(dt)
    g.handleCrystalPickup()
    g.updateDamageNumbers(dt)
    g.updateSpells(dt)

    if ebiten.IsKeyPressed(ebiten.KeyEscape) {
        g.GameState = StateTitle
    }
    return nil
}


func (g *Game) drawTitle(screen *ebiten.Image) {
	// Draw background image
	if g.TitleBG != nil {
		w, h := g.TitleBG.Size()
		scaleX := 800.0 / float64(w)
		scaleY := 600.0 / float64(h)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		screen.DrawImage(g.TitleBG, op)
	} else {
		vector.FillRect(screen, 0, 0, 800, 600, color.RGBA{20, 20, 40, 255}, false)
	}

	// Blink effects
	keyBlink := (math.Sin(g.BlinkTimer*6) + 1) / 2
	textBlink := (math.Sin(g.BlinkTimer*3) + 1) / 2

	// Base Y positions
	base := g.TitleYOffset

	// --- LINE 1: ENTER + Start Game ---
	drawVSKeyPrompt(screen, "[ENTER]", 480+base, g.TitleFont, keyBlink)
	if textBlink > 0.4 {
		drawVSStyledText(screen, "Start Game", 520+base, g.TitleFont)
	}

	// --- LINE 2: T + Training Room ---
	drawVSKeyPrompt(screen, "[T]", 560+base, g.TitleFont, keyBlink)
	if textBlink > 0.4 {
		drawVSStyledText(screen, "Training Room", 600+base, g.TitleFont)
	}

	// Fade overlay
	if g.FadeAlpha > 0 {
		a := uint8(g.FadeAlpha * 255)
		vector.FillRect(screen, 0, 0, 800, 600, color.RGBA{0, 0, 0, a}, false)
	}
}

func (g *Game) drawTrainingRoom(screen *ebiten.Image) {
	g.drawGameplay(screen)
	ebitenutil.DebugPrintAt(screen, "TRAINING ROOM - Press ESC to return", 10, 10)
}

func (g *Game) loadTrainingRoom() {
	// Reset world
	g.Enemies = []*Enemy{}
	g.Projectiles = []*Projectile{}
	g.FireTiles = []*FireTile{}
	g.Crystals = []*Crystal{}
	g.DamageNumbers = []*DamageNumber{}

	// Reset player state
	g.player.Pos = Vec{X: 400, Y: 300}
	g.player.HP = g.player.MaxHP
	g.player.XP = 0
	g.player.Level = 1
	g.player.Abilities = []*Ability{}
	g.player.HasDagger = false

	// Reset global stats
	g.CrystalMagnetRadius = 0
	g.InfiniteLoopChance = 0
	g.DaggerCooldown = 0
	g.DaggerRate = 3.2

	// Spawn a training dummy
	g.spawnTrainingDummy(Vec{X: 600, Y: 300})
}

func (g *Game) spawnTrainingDummy(pos Vec) {
	dummy := &Enemy{
		Type:     MonsterJavascript,
		Pos:      pos,
		HP:       999999,
		Speed:    0,
		Radius:   MonsterJavascript.Radius,
		Alive:    true,
		HitFlash: 0,
	}

	g.Enemies = append(g.Enemies, dummy)
}
