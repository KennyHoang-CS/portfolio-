package terminal

import (
	"image/color"

	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

func (t *Terminal) Draw(screen *ebiten.Image) {
	w, h := screen.Size()
	bounds := ui.TerminalBounds(w, h)

	x := bounds.Min.X + 16
	y := bounds.Min.Y + 48
	lineHeight := 16

	start := 0
	if len(t.lines) > (h-64)/lineHeight {
		start = len(t.lines) - (h-64)/lineHeight
	}

	for i, line := range t.lines[start:] {
		col := color.RGBA{180, 180, 180, 255}
		switch line.Color {
		case ColorInfo:
			col = color.RGBA{120, 200, 255, 255}
		case ColorError:
			col = color.RGBA{255, 80, 80, 255}
		case ColorSuccess:
			col = color.RGBA{120, 255, 120, 255}
		}
		text.Draw(screen, line.Text, basicfont.Face7x13, x, y+i*lineHeight, col)
	}
}