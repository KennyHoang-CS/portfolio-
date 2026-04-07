package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

func EditorBounds(w, h int) image.Rectangle {
	return image.Rect(0, 0, w/2, h)
}

func TerminalBounds(w, h int) image.Rectangle {
	return image.Rect(w/2, 0, w, h)
}

func DrawPanels(screen *ebiten.Image, w, h int) {
	// Editor panel
	editorRect := EditorBounds(w, h)
	editorImg := ebiten.NewImage(editorRect.Dx(), editorRect.Dy())
	editorImg.Fill(color.RGBA{30, 30, 60, 255})
	screen.DrawImage(editorImg, nil)

	// Terminal panel
	termRect := TerminalBounds(w, h)
	termImg := ebiten.NewImage(termRect.Dx(), termRect.Dy())
	termImg.Fill(color.RGBA{10, 10, 20, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(termRect.Min.X), float64(termRect.Min.Y))
	screen.DrawImage(termImg, op)

	// Labels
	text.Draw(screen, "Editor Panel", basicfont.Face7x13, 16, 24, color.White)
	text.Draw(screen, "Terminal Panel", basicfont.Face7x13, termRect.Min.X+16, 24, color.White)
}