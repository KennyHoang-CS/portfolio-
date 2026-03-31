package game

import (
	"image/color"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2"
)

func textWidthScaled(str string, scale float64) float64 {
    return float64(len(str) * 6) * scale
}

func centerTextX(text string, screenWidth int) int {
	charWidth := 6 // Ebiten debug font width
	textWidth := len(text) * charWidth
	return (screenWidth - textWidth) / 2
}

func neonText(screen *ebiten.Image, str string, x, y int, clr color.Color, scale float64) {
    op := &text.DrawOptions{}
    op.GeoM.Scale(scale, scale)
    op.GeoM.Translate(float64(x), float64(y))

    op.ColorScale.ScaleWithColor(clr)

    text.Draw(screen, str, ScoreFont, op)
}

func measureTextWidth(str string, face text.Face, scale float64) float64 {
    m := face.Metrics()
    lineSpacing := m.HAscent + m.HDescent + m.HLineGap
    w, _ := text.Measure(str, face, lineSpacing)
    return float64(w) * scale
}

func centerTextXScaled(str string, screenWidth int, face text.Face, scale float64) int {
    w := measureTextWidth(str, face, scale)
    return int((float64(screenWidth) - w) / 2)
}