package fonts

import (
    _ "embed"
    "bytes"

    text "github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed fonts/FiraMono-Medium.ttf
var editorTTF []byte

var EditorFont *text.GoTextFace

func init() {
    src, err := text.NewGoTextFaceSource(bytes.NewReader(editorTTF))
    if err != nil {
        panic(err)
    }

    EditorFont = &text.GoTextFace{
        Source: src,
        Size:   18,
        
    }
}