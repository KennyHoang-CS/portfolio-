package ui

import (
	"image/color"
	"log"
	"bytes"
	"fmt"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type ToolbarUI struct {
    UI *ebitenui.UI
}

func NewToolbarUI(onRun, onBuild, onFormat func()) *ToolbarUI {
    // Load EbitenUI-compatible font
	face, _ := loadFont(18)

    // Button images
    idle := image.NewNineSliceColor(color.NRGBA{70, 70, 70, 255})
    hover := image.NewNineSliceColor(color.NRGBA{90, 90, 90, 255})
    pressed := image.NewNineSliceColor(color.NRGBA{50, 50, 50, 255})

    root := widget.NewContainer(
        widget.ContainerOpts.Layout(widget.NewRowLayout(
            widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
            widget.RowLayoutOpts.Spacing(12),
        )),
        widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 0, 0})),
    )

    // Helper to add a button
    addBtn := func(label string, handler func()) {
        btn := widget.NewButton(
            widget.ButtonOpts.Image(&widget.ButtonImage{
                Idle:    idle,
                Hover:   hover,
                Pressed: pressed,
            }),
            widget.ButtonOpts.Text(label, &face, &widget.ButtonTextColor{
                Idle: color.White,
            }),
            widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
                handler()
            }),
        )
        root.AddChild(btn)
    }

    addBtn("Run", onRun)
    addBtn("Build", onBuild)
    addBtn("Format", onFormat)

    ui := &ebitenui.UI{Container: root}

    return &ToolbarUI{UI: ui}
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Error loading font: %w", err)
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}