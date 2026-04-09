package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func EditorBounds(w, h int) image.Rectangle {
    top := 80 // space for toolbar
    left := 0
    right := w/2 - 10 // small padding
    bottom := h - 20  // bottom padding
    return image.Rect(left, top, right, bottom)
}

func TerminalBounds(w, h int) image.Rectangle {
    top := 80
    left := w/2 + 10
    right := w
    bottom := h - 20
    return image.Rect(left, top, right, bottom)
}


func DrawPanels(screen *ebiten.Image, w, h int, font *text.GoTextFace) {
    editorRect := EditorBounds(w, h)
    terminalRect := TerminalBounds(w, h)

    // Editor panel background
    editorImg := ebiten.NewImage(editorRect.Dx(), editorRect.Dy())
    editorImg.Fill(color.RGBA{20, 26, 31, 255}) // #141A1F
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(float64(editorRect.Min.X), float64(editorRect.Min.Y))
    screen.DrawImage(editorImg, op)

    // Terminal panel background
    termImg := ebiten.NewImage(terminalRect.Dx(), terminalRect.Dy())
    termImg.Fill(color.RGBA{20, 20, 20, 255})
    op2 := &ebiten.DrawImageOptions{}
    op2.GeoM.Translate(float64(terminalRect.Min.X), float64(terminalRect.Min.Y))
    screen.DrawImage(termImg, op2)

    // Titles
    drawTitle(screen, "Editor Panel", editorRect, font)
    drawTitle(screen, "Terminal Panel", terminalRect, font)
}

func drawTitle(screen *ebiten.Image, label string, rect image.Rectangle, font *text.GoTextFace) {
    m := font.Metrics()
    lineHeight := m.HAscent + m.HDescent + m.HLineGap

    labelW, _ := text.Measure(label, font, lineHeight)

    // Center horizontally
    x := float64(rect.Min.X) + (float64(rect.Dx())-labelW)/2

    // Place near the top of the panel
    y := float64(rect.Min.Y) + float64(m.HAscent) + 8

    opts := &text.DrawOptions{}
    opts.GeoM.Translate(x, y)
    opts.ColorScale.ScaleWithColor(color.White)

    text.Draw(screen, label, font, opts)
}

func drawCenteredLabel(
    screen *ebiten.Image,
    textStr string,
    font *text.GoTextFace,
    x, y, w, h int,
    lineHeight float64,
) {
    labelW, _ := text.Measure(textStr, font, lineHeight)

    // Center horizontally
    textX := float64(x) + (float64(w)-labelW)/2

    // Rough vertical centering using ascent
    m := font.Metrics()
    textY := float64(y) + float64(h)/2 + m.HAscent/2

    opts := &text.DrawOptions{}
    opts.GeoM.Translate(textX, textY)
    opts.ColorScale.ScaleWithColor(color.White)

    text.Draw(screen, textStr, font, opts)
}