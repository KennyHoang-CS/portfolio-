package game

type Position struct {
    X, Y float64
}

type Velocity struct {
    VX, VY float64
}

type Hitbox struct {
    Radius float64
}

type PlayerTag struct{}

type BulletType int

const (
    BulletPetal BulletType = iota
    BulletStar
    BulletAmulet
    BulletOrb
    BulletKunai
    BulletArrow
)

type BulletTag struct {
    Speed float64
    Curve float64
    Type  BulletType
}

type Entity struct {
    Position *Position
    Velocity *Velocity
    Hitbox   *Hitbox
    Player   *PlayerTag
    Bullet   *BulletTag
    Spark    *SparkTag   
    Destroy  bool        
}

type SparkTag struct {
    Life int
}