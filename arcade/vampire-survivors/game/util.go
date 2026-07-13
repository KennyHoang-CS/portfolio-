package game

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
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

func LoadBytes(path string) []byte {
    b, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }
    return b
}

func drawCenteredText(screen *ebiten.Image, msg string, y float64, face font.Face, xOffset, yOffset int, col color.Color) {
    bounds := text.BoundString(face, msg)
    w := bounds.Dx()
    x := float64(800-w)/2 + float64(xOffset)
    text.Draw(screen, msg, face, int(x), int(y)+yOffset, col)
}

func drawVSStyledText(screen *ebiten.Image, msg string, y float64, face font.Face) {
    // Outline (thick black)
    drawCenteredText(screen, msg, y, face, -2, 0, color.Black)
    drawCenteredText(screen, msg, y, face, 2, 0, color.Black)
    drawCenteredText(screen, msg, y, face, 0, -2, color.Black)
    drawCenteredText(screen, msg, y, face, 0, 2, color.Black)

    // Glow (neon blue)
    glow := color.RGBA{80, 80, 255, 180}
    drawCenteredText(screen, msg, y, face, -1, -1, glow)
    drawCenteredText(screen, msg, y, face, 1, 1, glow)

    // Main text (white)
    drawCenteredText(screen, msg, y, face, 0, 0, color.White)
}

func drawVSKeyPrompt(screen *ebiten.Image, msg string, y float64, face font.Face, blink float64) {
    // Gold glow for key prompts
    glow := color.RGBA{255, 220, 80, 200}

    // Outline
    drawCenteredText(screen, msg, y, face, -2, 0, color.Black)
    drawCenteredText(screen, msg, y, face, 2, 0, color.Black)
    drawCenteredText(screen, msg, y, face, 0, -2, color.Black)
    drawCenteredText(screen, msg, y, face, 0, 2, color.Black)

    // Glow (stronger)
    drawCenteredText(screen, msg, y, face, -1, -1, glow)
    drawCenteredText(screen, msg, y, face, 1, 1, glow)

    // Main text (white)
    if blink > 0.3 {
        drawCenteredText(screen, msg, y, face, 0, 0, color.White)
    }
}
