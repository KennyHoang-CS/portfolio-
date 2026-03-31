package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *Game) drawCircle(screen *ebiten.Image, x, y, r float64, clr color.Color) {
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy <= r*r {
				screen.Set(int(x+dx), int(y+dy), clr)
			}
		}
	}
}

func (g *Game) DrawGlowCircle(screen *ebiten.Image, x, y, r float64, clr color.RGBA) {
	// Outer glow layers
	for i := 0; i < 4; i++ {
		alpha := uint8(40 - i*8)
		glow := color.RGBA{clr.R, clr.G, clr.B, alpha}
		g.drawCircle(screen, x, y, r+float64(i*3), glow)
	}

	// Core
	g.drawCircle(screen, x, y, r, clr)
}

func (g *Game) DrawRing(screen *ebiten.Image, x, y, r float64, thickness float64, clr color.Color) {
	outer := r
	inner := r - thickness

	for dy := -outer; dy <= outer; dy++ {
		for dx := -outer; dx <= outer; dx++ {
			dist := dx*dx + dy*dy
			if dist <= outer*outer && dist >= inner*inner {
				screen.Set(int(x+dx), int(y+dy), clr)
			}
		}
	}
}

func (g *Game) drawVictory(screen *ebiten.Image) {
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
    // YOU WIN! — dynamic scale
    // -------------------------
    title := "YOU WIN!"
    titleW, _ := text.Measure(title, ScoreFont, lineSpacing)

    titleTarget := screenW * 0.55
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

    neonText(screen, title, titleX, titleY, color.RGBA{120, 255, 120, 255}, titleScale)

    // -------------------------
    // Press ESC for Main Menu — dynamic scale + pulsing ESC
    // -------------------------
    staticColor := color.RGBA{200, 200, 255, 255}

    line := "Press ESC for Main Menu"
    lineW, _ := text.Measure(line, ScoreFont, lineSpacing)

    lineTarget := screenW * 0.45
    lineScale := lineTarget / float64(lineW)

    if lineScale > 1.6 {
        lineScale = 1.6
    }
    if lineScale < 0.7 {
        lineScale = 0.7
    }

    scaledLineW := float64(lineW) * lineScale
    lineX := int((screenW - scaledLineW) / 2)
    y := int(screenH * 0.55)

    left := "Press "
    key := "ESC"
    right := " for Main Menu"

    leftW, _ := text.Measure(left, ScoreFont, lineSpacing)
    keyW, _ := text.Measure(key, ScoreFont, lineSpacing)

    leftPx := int(float64(leftW) * lineScale)
    keyPx := int(float64(keyW) * lineScale)

    // Static left
    neonText(screen, left, lineX, y, staticColor, lineScale)

    // Pulse ESC
    p := pulse(g.frameCount+30, 0.08)
    keyScale := lineScale * (1.0 + 0.10*p)
    keyAlpha := uint8(180 + 75*p)
    keyColor := color.RGBA{200, 200, 255, keyAlpha}

    neonTextChar(screen, key, lineX+leftPx, y, keyColor, keyScale)

    // Static right
    neonText(screen, right, lineX+leftPx+keyPx, y, staticColor, lineScale)
}

func (g *Game) updateVictory() {
    if ebiten.IsKeyPressed(ebiten.KeyEscape) {
        g.reset(false, false)
        g.state = StateMenu
    }
}