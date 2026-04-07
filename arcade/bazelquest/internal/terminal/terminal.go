package terminal

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
	lines []Line
}

func New() *Terminal {
	return &Terminal{
		lines: make([]Line, 0, 128),
	}
}

func (t *Terminal) LogInfo(msg string) {
	t.lines = append(t.lines, Line{Text: msg, Color: ColorInfo})
}

func (t *Terminal) LogError(msg string) {
	t.lines = append(t.lines, Line{Text: msg, Color: ColorError})
}

func (t *Terminal) LogSuccess(msg string) {
	t.lines = append(t.lines, Line{Text: msg, Color: ColorSuccess})
}

func (t *Terminal) Update() {
	// later: scrolling, animations, etc.
}

func (t *Terminal) Lines() []Line {
	return t.lines
}