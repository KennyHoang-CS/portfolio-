package game

import (
    "log"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func LoadImage(path string) *ebiten.Image {
    imgFile, _, err := ebitenutil.NewImageFromFile(path)
    if err != nil {
        log.Fatal(err)
    }

    // No options allowed here — Ebiten only accepts the image itself
    return ebiten.NewImageFromImage(imgFile)
}

type Vec struct {
    X, Y float64
}
