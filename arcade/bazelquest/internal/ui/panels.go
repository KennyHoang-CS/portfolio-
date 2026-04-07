package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// DrawPanels renders a two-panel layout: editor (left) and terminal (right).
func DrawPanels(screen *ebiten.Image, w, h int) {
	leftW := w * 3 / 5
	rightW := w - leftW

	// Left panel image (Editor)
	leftImg := ebiten.NewImage(leftW, h)
	leftImg.Fill(EditorBg)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, 0)
	screen.DrawImage(leftImg, opts)

	// Right panel image (Terminal)
	rightImg := ebiten.NewImage(rightW, h)
	rightImg.Fill(TerminalBg)
	opts2 := &ebiten.DrawImageOptions{}
	opts2.GeoM.Translate(float64(leftW), 0)
	screen.DrawImage(rightImg, opts2)

	// Labels
	text.Draw(screen, "Editor Panel", BasicFontFace, 20, 28, LabelColor)
	text.Draw(screen, "Terminal Panel", BasicFontFace, leftW+20, 28, LabelColor)
}
