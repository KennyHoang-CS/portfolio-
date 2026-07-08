package game

import (
    "log"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func LoadImage(path string) *ebiten.Image {
    img, _, err := ebitenutil.NewImageFromFile(path)
    if err != nil {
        log.Fatalf("failed to load image %s: %v", path, err)
    }
    return img
}

type Vec struct {
    X, Y float64
}
