package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
)

func (g *Game) updateMovement() {
	for _, e := range g.world.Entities() {
		if e.Position == nil || e.Velocity == nil {
			continue
		}

		// Apply velocity
		e.Position.X += e.Velocity.VX
		e.Position.Y += e.Velocity.VY

		// Clamp player inside screen
		if e.Player != nil {
			r := e.Hitbox.Radius

			if e.Position.X < r {
				e.Position.X = r
			}
			if e.Position.X > float64(ScreenWidth)-r {
				e.Position.X = float64(ScreenWidth) - r
			}
			if e.Position.Y < r {
				e.Position.Y = r
			}
			if e.Position.Y > float64(ScreenHeight)-r {
				e.Position.Y = float64(ScreenHeight) - r
			}
		}
	}
}

func (g *Game) renderEntities(screen *ebiten.Image) {
	for _, e := range g.world.Entities() {
		if e.Position == nil {
			continue
		}

		x := float64(e.Position.X)
		y := float64(e.Position.Y)

		// -------------------------
		// Bullet rendering (glow + sprite)
		// -------------------------
		if e.Bullet != nil {

			// Determine bullet glow color
			var glow color.RGBA

			if e.Bullet.Curve != 0 {
				glow = color.RGBA{0xcc, 0x33, 0xff, 60} // purple glow
			} else {
				glow = color.RGBA{0xff, 0x55, 0x88, 60} // pink-red glow
			}

			// Glow bloom
			vector.FillCircle(screen, float32(x), float32(y), 7, glow, false)

			// Select sprite
			var img *ebiten.Image
			var targetSize float64

			switch e.Bullet.Type {
			case BulletPetal:
				img = ImgBulletPetal
				targetSize = 24
			case BulletStar:
				img = ImgBulletStar
				targetSize = 24
			case BulletAmulet:
				img = ImgBulletAmulet
				targetSize = 32
			case BulletOrb:
				img = ImgBulletOrb
				targetSize = 20
			case BulletKunai:
				img = ImgBulletKunai
				targetSize = 32
			case BulletArrow:
				img = ImgBulletArrow
				targetSize = 24
			}

			if img != nil {
				op := &ebiten.DrawImageOptions{}
				w := float64(img.Bounds().Dx())
				scale := targetSize / w
				op.GeoM.Scale(scale, scale)
				op.GeoM.Translate(x-targetSize/2, y-targetSize/2)
				screen.DrawImage(img, op)
			}

			continue
		}

		// -------------------------
		// Player rendering (Touhou-style, no sprite)
		// -------------------------
		if e.Player != nil {
			// 1. Aura ring (visible + harmonious)
			t := float64(g.Score) * 0.12
			ax := x + math.Sin(t)*2
			ay := y + math.Cos(t)*2

			vector.FillCircle(
				screen,
				float32(ax), float32(ay),
				13,                               // larger aura so it shows
				color.RGBA{0x66, 0xFF, 0xAA, 40}, // visible mint-green glow
				false,
			)

			// 2. Large diamond silhouette (30px tall)
			var path vector.Path
			path.MoveTo(float32(x), float32(y-15))
			path.LineTo(float32(x+10), float32(y))
			path.LineTo(float32(x), float32(y+15))
			path.LineTo(float32(x-10), float32(y))
			path.Close()

			drawOp := &vector.DrawPathOptions{}
			drawOp.AntiAlias = true
			drawOp.ColorScale.ScaleWithColor(color.RGBA{0x00, 0xAA, 0x33, 0xFF}) // deep emerald green

			vector.FillPath(screen, &path, nil, drawOp)

			// 3. Hitbox dot
			vector.FillCircle(screen, float32(x), float32(y), 2, color.White, false)

			continue
		}

		// -------------------------
		// Graze sparks
		// -------------------------
		if e.Spark != nil {
			vector.FillCircle(screen, float32(x), float32(y), 3, color.RGBA{255, 255, 255, 180}, false)
			e.Spark.Life--
			if e.Spark.Life <= 0 {
				e.Destroy = true
			}
		}
	}
}
