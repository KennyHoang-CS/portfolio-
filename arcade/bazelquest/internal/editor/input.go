package editor

import (
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (e *Editor) Update() {
	// Arrow keys
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		e.moveCursorLine(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		e.moveCursorLine(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		e.moveCursorCol(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		e.moveCursorCol(1)
	}

	// Backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		e.backspace()
	}

	// Enter
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		e.insertNewline()
	}

	// Text input
	for _, r := range ebiten.InputChars() {
		if r == '\n' || r == '\r' {
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		e.insertRune(r)
	}
}

func (e *Editor) insertRune(r rune) {
	line := e.lines[e.cursor.Line]
	if e.cursor.Col < 0 || e.cursor.Col > len(line) {
		e.cursor.Col = len(line)
	}
	newLine := line[:e.cursor.Col] + string(r) + line[e.cursor.Col:]
	e.lines[e.cursor.Line] = newLine
	e.cursor.Col++
}

func (e *Editor) backspace() {
	if e.cursor.Col > 0 {
		line := e.lines[e.cursor.Line]
		e.lines[e.cursor.Line] = line[:e.cursor.Col-1] + line[e.cursor.Col:]
		e.cursor.Col--
		return
	}

	// merge with previous line
	if e.cursor.Line > 0 {
		prev := e.lines[e.cursor.Line-1]
		cur := e.lines[e.cursor.Line]
		e.lines = append(e.lines[:e.cursor.Line], e.lines[e.cursor.Line+1:]...)
		e.cursor.Line--
		e.cursor.Col = len(prev)
		e.lines[e.cursor.Line] = prev + cur
	}
}

func (e *Editor) insertNewline() {
	line := e.lines[e.cursor.Line]
	left := line[:e.cursor.Col]
	right := line[e.cursor.Col:]
	e.lines[e.cursor.Line] = left
	e.lines = append(e.lines[:e.cursor.Line+1], append([]string{right}, e.lines[e.cursor.Line+1:]...)...)
	e.cursor.Line++
	e.cursor.Col = 0
}