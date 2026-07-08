package main

import (
    "log"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/KennyHoang-CS/portfolio/vampire-survivors/game"
)

func main() {
    g := game.NewGame()

    ebiten.SetWindowSize(800, 600)
    ebiten.SetWindowTitle("Go Gopher Survivors")
    ebiten.SetFullscreen(false) // optional

    if err := ebiten.RunGame(g); err != nil {
        log.Fatal(err)
    }
}
