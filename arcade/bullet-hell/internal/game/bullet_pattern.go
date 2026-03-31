package game

import "github.com/hajimehoshi/ebiten/v2"

type BulletPattern interface {
    Update(g *Game, difficulty int)
    Reset()
    IsFinished() bool
    Name() string

    // Optional: patterns that draw visuals implement this
    Draw(screen *ebiten.Image, g *Game)
}