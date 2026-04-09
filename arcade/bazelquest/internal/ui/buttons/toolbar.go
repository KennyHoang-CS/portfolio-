package buttons

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Toolbar struct {
	Buttons []*UIButton
	X, Y    int
	Spacing int
}

func NewToolbar(x, y int) *Toolbar {
	return &Toolbar{
		X:       x,
		Y:       y,
		Spacing: 12,
	}
}

func (t *Toolbar) AddButton(label string, w, h int, onClick func()) {
	btn := &UIButton{
		X:      t.X,
		Y:      t.Y,
		W:      w,
		H:      h,
		Label:  label,
		OnClick: onClick,
	}

	t.Buttons = append(t.Buttons, btn)

	// Move next button position
	t.X += w + t.Spacing
}

func (t *Toolbar) Update() {
	for _, b := range t.Buttons {
		b.Update()
	}
}

func (t *Toolbar) Draw(screen *ebiten.Image, font *text.GoTextFace) {
	for _, b := range t.Buttons {
		b.Draw(screen, font)
	}
}