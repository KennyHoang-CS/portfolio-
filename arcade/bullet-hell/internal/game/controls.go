package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func pulse(frame int, speed float64) float64 {
	return 0.5 + 0.5*math.Sin(float64(frame)*speed)
}

func neonTextChar(screen *ebiten.Image, s string, x, y int, col color.Color, scale float64) {
	op := &text.DrawOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	text.Draw(screen, s, ScoreFont, op)
}

func (g *Game) updateControls() {
    if ebiten.IsKeyPressed(ebiten.KeyEscape) {
        g.state = StateMenu

        if g.menuBGM != nil && !g.menuBGM.IsPlaying() {
            g.menuBGM.Rewind()
            g.menuBGM.Play()
        }
    }
}

func (g *Game) drawControls(screen *ebiten.Image) {
    // Background (same as menu)
    g.drawBackground(screen)
    g.drawMenuParticles(screen)

    // Dim overlay
    dim := ebiten.NewImage(g.ScreenWidth(), g.ScreenHeight())
    dim.Fill(color.RGBA{0, 0, 0, 160})
    screen.DrawImage(dim, nil)

    // Common metrics
    m := ScoreFont.Metrics()
    lineSpacing := m.HAscent + m.HDescent + m.HLineGap

    // -------------------------
    // TITLE
    // -------------------------
    title := "CONTROLS"
    titleW, _ := text.Measure(title, ScoreFont, lineSpacing)

    titleTarget := float64(g.ScreenWidth()) * 0.50
    titleScale := titleTarget / float64(titleW)

    if titleScale > 2.0 {
        titleScale = 2.0
    }
    if titleScale < 0.6 {
        titleScale = 0.6
    }

    scaledTitleW := float64(titleW) * titleScale
    titleX := int((float64(g.ScreenWidth()) - scaledTitleW) / 2)

    neonText(screen, title, titleX, 120, color.RGBA{180, 180, 255, 255}, titleScale)

    // -------------------------
    // CONTROL LINE 1 — Move: W A S D
    // -------------------------
    line1 := "Move: W A S D"
    line1W, _ := text.Measure(line1, ScoreFont, lineSpacing)

    line1Target := float64(g.ScreenWidth()) * 0.40
    line1Scale := line1Target / float64(line1W)

    if line1Scale > 1.4 {
        line1Scale = 1.4
    }
    if line1Scale < 0.7 {
        line1Scale = 0.7
    }

    scaledLine1W := float64(line1W) * line1Scale
    line1X := int((float64(g.ScreenWidth()) - scaledLine1W) / 2)

    // Split into static + pulsing keys
    keys := []string{"Move:", " ", "W", " ", "A", " ", "S", " ", "D"}

    cursorX := line1X
    y := 220

    staticColor := color.RGBA{180, 180, 255, 255}

    for i, k := range keys {
        if k == "W" || k == "A" || k == "S" || k == "D" {
            // Pulse
            p := pulse(g.frameCount+i*12, 0.10)
            scale := line1Scale * (1.0 + 0.12*p)
            alpha := uint8(180 + 75*p)
            col := color.RGBA{200, 200, 255, alpha}

            advance, _ := text.Measure(k, ScoreFont, lineSpacing)
            pixelWidth := int(float64(advance) * scale)

            neonTextChar(screen, k, cursorX, y, col, scale)
            cursorX += pixelWidth
        } else {
            // STATIC — use neonText to match original look
            advance, _ := text.Measure(k, ScoreFont, lineSpacing)
            pixelWidth := int(float64(advance) * line1Scale)

            neonText(screen, k, cursorX, y, staticColor, line1Scale)
            cursorX += pixelWidth
        }
    }

    // -------------------------
    // CONTROL LINE 1B — static
    // -------------------------
    line1b := "(up, left, down, right)"
    line1bW, _ := text.Measure(line1b, ScoreFont, lineSpacing)

    line1bTarget := float64(g.ScreenWidth()) * 0.40
    line1bScale := line1bTarget / float64(line1bW)

    if line1bScale > 1.4 {
        line1bScale = 1.4
    }
    if line1bScale < 0.7 {
        line1bScale = 0.7
    }

    scaledLine1bW := float64(line1bW) * line1bScale
    line1bX := int((float64(g.ScreenWidth()) - scaledLine1bW) / 2)

    line1bY := 260
    neonText(screen, line1b, line1bX, line1bY, staticColor, line1bScale)

    // Dynamic spacing
    line1bHeight := (m.HAscent + m.HDescent + m.HLineGap) * line1bScale
    padding := 20.0

    // -------------------------
    // CONTROL LINE 2 — Back: ESC
    // -------------------------
    line2 := "Back: ESC"
    line2W, _ := text.Measure(line2, ScoreFont, lineSpacing)

    line2Target := float64(g.ScreenWidth()) * 0.40
    line2Scale := line2Target / float64(line2W)

    if line2Scale > 1.4 {
        line2Scale = 1.4
    }
    if line2Scale < 0.7 {
        line2Scale = 0.7
    }

    scaledLine2W := float64(line2W) * line2Scale
    line2X := int((float64(g.ScreenWidth()) - scaledLine2W) / 2)

    line2Y := int(float64(line1bY) + line1bHeight + padding)

    // Split into static + pulsing ESC
    backLabel := "Back: "
    backW, _ := text.Measure(backLabel, ScoreFont, lineSpacing)
    backPx := int(float64(backW) * line2Scale)

    // Draw static "Back:"
    neonText(screen, backLabel, line2X, line2Y, staticColor, line2Scale)

    // Pulse ESC
    p := pulse(g.frameCount, 0.08)
    escScale := line2Scale * (1.0 + 0.10*p)
    escAlpha := uint8(180 + 75*p)
    escColor := color.RGBA{200, 200, 255, escAlpha}

    neonTextChar(screen, "ESC", line2X+backPx, line2Y, escColor, escScale)
}