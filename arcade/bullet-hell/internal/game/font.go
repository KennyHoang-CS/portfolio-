package game

import (
    _ "embed"
    "bytes"

    text "github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed assets/fonts/ARCADE_N.ttf
var arcadeTTF []byte

var ScoreFont *text.GoTextFace

func init() {
    // Create a font source from the embedded TTF
    src, err := text.NewGoTextFaceSource(bytes.NewReader(arcadeTTF))
    if err != nil {
        panic(err)
    }

    // Create a face with a given size
    ScoreFont = &text.GoTextFace{
        Source: src,
        Size:   32,
    }
}