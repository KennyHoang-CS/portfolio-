package game

import (
    "bytes"
    "image"
	_ "embed"
    _ "image/png"

    "github.com/hajimehoshi/ebiten/v2"
)

// --------------------------------------
// Embedded bullet sprite PNGs
// --------------------------------------

//go:embed assets/bullets/bullet_petal.png
var bulletPetalPNG []byte

//go:embed assets/bullets/bullet_star.png
var bulletStarPNG []byte

//go:embed assets/bullets/bullet_amulet.png
var bulletAmuletPNG []byte

//go:embed assets/bullets/bullet_orb.png
var bulletOrbPNG []byte

//go:embed assets/bullets/bullet_kunai.png
var bulletKunaiPNG []byte

//go:embed assets/bullets/bullet_arrowhead.png
var bulletArrowPNG []byte

// --------------------------------------
// Loaded Ebiten images
// --------------------------------------

var (
    ImgBulletPetal  *ebiten.Image
    ImgBulletStar   *ebiten.Image
    ImgBulletAmulet *ebiten.Image
    ImgBulletOrb    *ebiten.Image
    ImgBulletKunai  *ebiten.Image
    ImgBulletArrow  *ebiten.Image
)

// --------------------------------------
// LoadAssets loads all embedded sprites
// --------------------------------------

func LoadAssets() error {
    decode := func(data []byte) (*ebiten.Image, error) {
        img, _, err := image.Decode(bytes.NewReader(data))
        if err != nil {
            return nil, err
        }
        return ebiten.NewImageFromImage(img), nil
    }

    var err error

    ImgBulletPetal, err = decode(bulletPetalPNG)
    if err != nil { return err }

    ImgBulletStar, err = decode(bulletStarPNG)
    if err != nil { return err }

    ImgBulletAmulet, err = decode(bulletAmuletPNG)
    if err != nil { return err }

    ImgBulletOrb, err = decode(bulletOrbPNG)
    if err != nil { return err }

    ImgBulletKunai, err = decode(bulletKunaiPNG)
    if err != nil { return err }

    ImgBulletArrow, err = decode(bulletArrowPNG)
    if err != nil { return err }

    return nil
}