package editor

import (
	"image/color"

	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

func (e *Editor) Draw(screen *ebiten.Image) {
	w, h := screen.Size()
	bounds := ui.EditorBounds(w, h)

	x := bounds.Min.X + 16
	y := bounds.Min.Y + 48
	lineHeight := 16

	for i, line := range e.lines {
		text.Draw(screen, line, basicfont.Face7x13, x, y+i*lineHeight, color.White)
	}

	// simple cursor
	cursorX := x + e.cursor.Col*7
	cursorY := y + e.cursor.Line*lineHeight
	text.Draw(screen, "|", basicfont.Face7x13, cursorX, cursorY, color.RGBA{200, 200, 0, 255})
}