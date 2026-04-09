package terminal

import (
	"image/color"
	"strings"

	fonts "github.com/KennyHoang-CS/portfolio/bazelquest/internal/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Line struct {
	Text  string
	Color Color
}

type Color int

const (
	ColorInfo Color = iota
	ColorError
	ColorSuccess
)

type Terminal struct {
	lines      []Line
	offsetX    int
	offsetY    int
	lineHeight int
}

func New() *Terminal {
	metrics := fonts.EditorFont.Metrics()
	lineSpacing := metrics.HAscent + metrics.HDescent + metrics.HLineGap

	return &Terminal{
		lines:      make([]Line, 0, 128),
		offsetX:    0,
		offsetY:    0,
		lineHeight: int(lineSpacing),
	}
}

func (t *Terminal) SetOffset(x, y int) {
	t.offsetX = x
	t.offsetY = y
}

func (t *Terminal) LogInfo(msg string) {
    for _, line := range strings.Split(msg, "\n") {
        t.lines = append(t.lines, Line{Text: line, Color: ColorInfo})
    }
}

func (t *Terminal) LogError(msg string) {
    for _, line := range strings.Split(msg, "\n") {
        t.lines = append(t.lines, Line{Text: line, Color: ColorError})
    }
}

func (t *Terminal) LogSuccess(msg string) {
    for _, line := range strings.Split(msg, "\n") {
        t.lines = append(t.lines, Line{Text: line, Color: ColorSuccess})
    }
}

func (t *Terminal) Update() {
	// later: scrolling, animations, etc.
}

func (t *Terminal) Lines() []Line {
	return t.lines
}

// -----------------------------
// Rendering
// -----------------------------
func (t *Terminal) Draw(screen *ebiten.Image) {
	y := t.offsetY

	for _, line := range t.lines {
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(t.offsetX), float64(y))
		opts.ColorScale.ScaleWithColor(t.colorFor(line.Color))

		text.Draw(screen, line.Text, fonts.EditorFont, opts)
		y += t.lineHeight
	}
}

func (t *Terminal) colorFor(c Color) color.RGBA {
	switch c {
	case ColorError:
		return color.RGBA{255, 80, 80, 255}
	case ColorSuccess:
		return color.RGBA{120, 255, 120, 255}
	default:
		return color.RGBA{200, 200, 255, 255}
	}
}