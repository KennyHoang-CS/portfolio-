package editor

type Editor struct {
	lines  []string
	cursor Cursor
}

type Cursor struct {
	Line int
	Col  int
}

func New() *Editor {
	return &Editor{
		lines: []string{""},
		cursor: Cursor{
			Line: 0,
			Col:  0,
		},
	}
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