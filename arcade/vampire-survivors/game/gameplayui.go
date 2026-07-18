package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawHPBar(screen *ebiten.Image) {
	camX, camY := g.CameraOffset()

	barWidth := float32(40)
	barHeight := float32(6)

	x := float32(g.player.Pos.X - camX - float64(barWidth)/2)
	y := float32(g.player.Pos.Y - camY + 45)

	hpPercent := float32(g.player.HP / g.player.MaxHP)

	vector.FillRect(screen, x, y, barWidth, barHeight, color.RGBA{80, 0, 0, 255}, false)
	vector.FillRect(screen, x, y, barWidth*hpPercent, barHeight, color.RGBA{255, 40, 40, 255}, false)
}

func (g *Game) drawXPBar(screen *ebiten.Image) {
	barWidth := float32(800)
	barHeight := float32(10)
	x := float32(0)
	y := float32(0)

	xpNeeded := float32(g.player.Level * 10)
	xpPercent := float32(g.player.XP) / xpNeeded

	vector.FillRect(screen, x, y, barWidth, barHeight, color.RGBA{30, 30, 60, 255}, false)
	vector.FillRect(screen, x, y, barWidth*xpPercent, barHeight, color.RGBA{80, 80, 255, 255}, false)
}
