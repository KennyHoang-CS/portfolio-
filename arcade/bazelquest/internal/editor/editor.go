package editor

import (
	"image/color"
	"strings"

	fonts "github.com/KennyHoang-CS/portfolio/bazelquest/internal/assets"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/workspace"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Editor struct {
	lines      []string
	cursor     Cursor
	offsetX    int
	offsetY    int
	lineHeight int

	currentPath string
	workspace   *workspace.Workspace
}

type Cursor struct {
	Line int
	Col  int
}

func New(ws *workspace.Workspace) *Editor {
	metrics := fonts.EditorFont.Metrics()
	lineSpacing := metrics.HAscent + metrics.HDescent + metrics.HLineGap

	return &Editor{
		lines:      []string{""},
		cursor:     Cursor{Line: 0, Col: 0},
		offsetX:    0,
		offsetY:    0,
		lineHeight: int(lineSpacing),

		workspace: ws,
	}
}

func (e *Editor) Open(path string) {
    e.currentPath = path

    if content, ok := e.workspace.Read(path); ok {
        e.SetText(content)
    } else {
        // If file doesn't exist, create it empty
        e.workspace.Write(path, "")
        e.SetText("")
    }
}

func (e *Editor) Save() {
    if e.currentPath == "" {
        return
    }
    e.workspace.Write(e.currentPath, e.Buffer())
}

// Allow UI to reposition editor dynamically
func (e *Editor) SetOffset(x, y int) {
	e.offsetX = x
	e.offsetY = y
}

func (e *Editor) InsertLine(line string) {
	e.lines = append(e.lines, line)
}

func (e *Editor) Buffer() string {
	out := ""
	for i, l := range e.lines {
		if i > 0 {
			out += "\n"
		}
		out += l
	}
	return out
}

func (e *Editor) moveCursorLine(delta int) {
	newLine := e.cursor.Line + delta
	if newLine < 0 {
		newLine = 0
	}
	if newLine >= len(e.lines) {
		newLine = len(e.lines) - 1
	}
	e.cursor.Line = newLine

	if e.cursor.Col > len(e.lines[e.cursor.Line]) {
		e.cursor.Col = len(e.lines[e.cursor.Line])
	}
}

func (e *Editor) moveCursorCol(delta int) {
	newCol := e.cursor.Col + delta
	if newCol < 0 {
		newCol = 0
	}
	if newCol > len(e.lines[e.cursor.Line]) {
		newCol = len(e.lines[e.cursor.Line])
	}
	e.cursor.Col = newCol
}

// -----------------------------
// Mouse click → cursor position
// -----------------------------
func (e *Editor) Click(x, y int) {
	editorX := x - e.offsetX
	editorY := y - e.offsetY

	lineIdx := editorY / e.lineHeight
	if lineIdx < 0 {
		lineIdx = 0
	}
	if lineIdx >= len(e.lines) {
		lineIdx = len(e.lines) - 1
	}
	line := e.lines[lineIdx]

	var prevAdvance float64

	for i := range line {
		curAdvance := text.Advance(line[:i], fonts.EditorFont)

		var x0, x1 int
		if i == 0 {
			x0 = 0
		} else {
			x0 = int((prevAdvance + curAdvance) / 2)
		}

		var nextAdvance float64
		if i+1 <= len(line) {
			nextAdvance = text.Advance(line[:i+1], fonts.EditorFont)
		} else {
			nextAdvance = curAdvance
		}

		x1 = int((curAdvance + nextAdvance) / 2)

		if editorX >= x0 && editorX < x1 {
			e.cursor.Line = lineIdx
			e.cursor.Col = i
			return
		}

		prevAdvance = curAdvance
	}

	// Clicked past end
	e.cursor.Line = lineIdx
	e.cursor.Col = len(line)
}

// -----------------------------
// Rendering
// -----------------------------
func (e *Editor) Draw(screen *ebiten.Image) {
	y := e.offsetY

	for _, line := range e.lines {
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(e.offsetX), float64(y))
		text.Draw(screen, line, fonts.EditorFont, opts)
		y += e.lineHeight
	}

	cursorX := e.offsetX + e.cursorPixelX()
	cursorY := e.offsetY + e.cursor.Line*e.lineHeight

	cursorImg := ebiten.NewImage(2, e.lineHeight)
	cursorImg.Fill(color.White)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(cursorX), float64(cursorY))
	screen.DrawImage(cursorImg, op)
}

// Compute pixel X using proportional font advance
func (e *Editor) cursorPixelX() int {
	line := e.lines[e.cursor.Line]
	if e.cursor.Col > len(line) {
		return 0
	}
	sub := line[:e.cursor.Col]
	advance := text.Advance(sub, fonts.EditorFont)
	return int(advance)
}

func (e *Editor) SetText(s string) {
	e.lines = strings.Split(s, "\n")
	e.cursor.Line = 0
	e.cursor.Col = 0
}
