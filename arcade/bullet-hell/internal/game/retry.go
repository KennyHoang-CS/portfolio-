package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *Game) updateRetry() {
	// Retry (R)
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		if g.retryBGM != nil {
			g.retryBGM.Pause()
		}

		g.reset(true, true)

		// Only operate on gameBGM if it exists
		if g.gameBGM != nil {
			g.gameBGM.Rewind()
			g.gameBGM.SetVolume(0.1)

			// WASM CAN play now because user pressed R
			g.gameBGM.Play()
		}

		g.state = StatePlaying
		return
	}
	
    // Main Menu (ESC)
    if ebiten.IsKeyPressed(ebiten.KeyEscape) {
        if g.retryBGM != nil {
            g.retryBGM.Pause()
        }

        g.reset(false, true)

        // Just resume existing menuBGM if we have it
        if g.menuBGM != nil {
            g.menuBGM.Rewind()
            g.menuBGM.SetVolume(0.05)
            g.menuBGM.Play()
        }

        g.state = StateMenu
        return
    }


}

func (g *Game) drawRetry(screen *ebiten.Image) {
	g.drawBackground(screen)
	g.drawMenuParticles(screen)

	// Dim overlay
	dim := ebiten.NewImage(g.ScreenWidth(), g.ScreenHeight())
	dim.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(dim, nil)

	// Common metrics
	m := ScoreFont.Metrics()
	lineSpacing := m.HAscent + m.HDescent + m.HLineGap

	screenW := float64(g.ScreenWidth())
	screenH := float64(g.ScreenHeight())

	// -------------------------
	// YOU DIED — dynamic scale
	// -------------------------
	title := "YOU DIED"
	titleW, _ := text.Measure(title, ScoreFont, lineSpacing)

	titleTarget := screenW * 0.55 // 55% of screen width
	titleScale := titleTarget / float64(titleW)

	if titleScale > 2.5 {
		titleScale = 2.5
	}
	if titleScale < 0.8 {
		titleScale = 0.8
	}

	scaledTitleW := float64(titleW) * titleScale
	titleX := int((screenW - scaledTitleW) / 2)
	titleY := int(screenH * 0.22)

	neonText(screen, title, titleX, titleY, color.RGBA{255, 120, 120, 255}, titleScale)

	// -------------------------
	// Press R to Retry — dynamic scale
	// -------------------------
	staticColor := color.RGBA{200, 200, 255, 255}

	line1 := "Press R to Retry"
	line1W, _ := text.Measure(line1, ScoreFont, lineSpacing)

	line1Target := screenW * 0.45
	line1Scale := line1Target / float64(line1W)

	if line1Scale > 1.6 {
		line1Scale = 1.6
	}
	if line1Scale < 0.7 {
		line1Scale = 0.7
	}

	scaledLine1W := float64(line1W) * line1Scale
	line1X := int((screenW - scaledLine1W) / 2)
	y1 := int(screenH * 0.50)

	// Split into parts
	left := "Press "
	key := "R"
	right := " to Retry"

	leftW, _ := text.Measure(left, ScoreFont, lineSpacing)
	keyW, _ := text.Measure(key, ScoreFont, lineSpacing)

	leftPx := int(float64(leftW) * line1Scale)
	keyPx := int(float64(keyW) * line1Scale)

	// Static left
	neonText(screen, left, line1X, y1, staticColor, line1Scale)

	// Pulse R
	p := pulse(g.frameCount, 0.08)
	keyScale := line1Scale * (1.0 + 0.10*p)
	keyAlpha := uint8(180 + 75*p)
	keyColor := color.RGBA{200, 200, 255, keyAlpha}

	neonTextChar(screen, key, line1X+leftPx, y1, keyColor, keyScale)

	// Static right
	neonText(screen, right, line1X+leftPx+keyPx, y1, staticColor, line1Scale)

	// -------------------------
	// Press ESC for main menu — dynamic scale
	// -------------------------
	line2 := "Press ESC for main menu"
	line2W, _ := text.Measure(line2, ScoreFont, lineSpacing)

	line2Target := screenW * 0.45
	line2Scale := line2Target / float64(line2W)

	if line2Scale > 1.6 {
		line2Scale = 1.6
	}
	if line2Scale < 0.7 {
		line2Scale = 0.7
	}

	scaledLine2W := float64(line2W) * line2Scale
	line2X := int((screenW - scaledLine2W) / 2)
	y2 := int(screenH * 0.58)

	left2 := "Press "
	key2 := "ESC"
	right2 := " for main menu"

	left2W, _ := text.Measure(left2, ScoreFont, lineSpacing)
	key2W, _ := text.Measure(key2, ScoreFont, lineSpacing)

	left2Px := int(float64(left2W) * line2Scale)
	key2Px := int(float64(key2W) * line2Scale)

	// Static left
	neonText(screen, left2, line2X, y2, staticColor, line2Scale)

	// Pulse ESC
	p2 := pulse(g.frameCount+30, 0.08)
	key2Scale := line2Scale * (1.0 + 0.10*p2)
	key2Alpha := uint8(180 + 75*p2)
	key2Color := color.RGBA{200, 200, 255, key2Alpha}

	neonTextChar(screen, key2, line2X+left2Px, y2, key2Color, key2Scale)

	// Static right
	neonText(screen, right2, line2X+left2Px+key2Px, y2, staticColor, line2Scale)
}

func (g *Game) reset(loadGameplayBGM bool, shufflePatterns bool) {
	g.world = NewWorld()
	g.Score = 0
	g.initPlayer()

	if g.patternManager != nil {
		g.patternManager.Reset()
		if shufflePatterns {
			g.patternManager.Shuffle()
		}
	}

	if g.gameBGM != nil {
		g.gameBGM.Pause()
	}

	if loadGameplayBGM {
		var next *audio.Player

		if g.preloader != nil && g.preloader.Stage2Done() {
			next = LoadRandomGameplayFull()
		} else {
			next = LoadRandomGameplayStage1()
		}

		if next != nil {
			next.Rewind()
			next.SetVolume(0.1)
		}

		g.gameBGM = next
	}
}
