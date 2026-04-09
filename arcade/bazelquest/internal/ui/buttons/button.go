package buttons

import (
    "image/color"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
    "github.com/hajimehoshi/ebiten/v2/text/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"
)

type UIButton struct {
    X, Y int
    W, H int

    Label   string
    OnClick func()

    Hovered bool
}

// Update handles hover + click
func (b *UIButton) Update() {
    x, y := ebiten.CursorPosition()

    inside := x >= b.X && x <= b.X+b.W &&
        y >= b.Y && y <= b.Y+b.H

    b.Hovered = inside

    if inside && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
        if b.OnClick != nil {
            b.OnClick()
        }
    }
}

// Draw renders the button using the provided font
func (b *UIButton) Draw(dst *ebiten.Image, face *text.GoTextFace) {
    // Background
    bg := color.RGBA{60, 60, 60, 255}
    if b.Hovered {
        bg = color.RGBA{80, 80, 80, 255}
    }

    vector.DrawFilledRect(
        dst,
        float32(b.X),
        float32(b.Y),
        float32(b.W),
        float32(b.H),
        bg,
        false,
    )

    // Border
    vector.StrokeRect(
        dst,
        float32(b.X),
        float32(b.Y),
        float32(b.W),
        float32(b.H),
        1,
        color.White,
        false,
    )

    // Text metrics
    m := face.Metrics()
    lineHeight := m.HAscent + m.HDescent + m.HLineGap

    labelW, _ := text.Measure(b.Label, face, lineHeight)

    // Horizontal center
    textX := float64(b.X) + (float64(b.W)-labelW)/2

    // Vertical center (baseline math)
    textY := float64(b.Y) + float64(b.H)/2 + float64(m.HAscent)/2

    // Draw text
    op := &text.DrawOptions{}
    op.GeoM.Translate(textX, textY)
    op.ColorScale.ScaleWithColor(color.White)

    text.Draw(dst, b.Label, face, op)
}