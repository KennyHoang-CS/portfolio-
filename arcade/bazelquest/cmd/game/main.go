package main

import (
	"log"

	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/build"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/editor"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/input"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/terminal"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	editor   *editor.Editor
	terminal *terminal.Terminal
}

func NewGame() *Game {
	e := editor.New()
	t := terminal.New()
	t.LogInfo("BazelQuest prototype started.")
	t.LogInfo("Type in the editor. Press Enter+Ctrl (simulated later) or just F5 to build (for now: Enter triggers build).")
	return &Game{
		editor:   e,
		terminal: t,
	}
}

func (g *Game) Update() error {
	input.Update(g.editor, g.terminal)

	// TEMP: pressing F5 or Enter triggers a fake build
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		build.Run(g.editor.Buffer(), g.terminal)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	w, h := screen.Size()
	ui.DrawPanels(screen, w, h)
	g.editor.Draw(screen)
	g.terminal.Draw(screen)
}

func (g *Game) Layout(outW, outH int) (int, int) { return 1280, 720 }

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("BazelQuest - Prototype")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}