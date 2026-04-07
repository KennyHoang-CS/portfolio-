package main

import (
	"log"

	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Update() error { return nil }

func (g *Game) Draw(screen *ebiten.Image) {
	w, h := screen.Size()
	ui.DrawPanels(screen, w, h)
}

func (g *Game) Layout(outW, outH int) (int, int) { return 1280, 720 }

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("BazelQuest - Prototype")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
