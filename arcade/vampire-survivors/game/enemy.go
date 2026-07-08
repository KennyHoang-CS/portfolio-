package game

import "math"

type Enemy struct {
    Type     *MonsterType
    Pos      Vec
    HP       float64
    Speed    float64
    Radius   float64
    Alive    bool
    HitFlash float64
}

func NewEnemy(mt *MonsterType, x, y float64) *Enemy {
    return &Enemy{
        Type:     mt,
        Pos:      Vec{X: x, Y: y},
        HP:       mt.HP,
        Speed:    mt.Speed,
        Radius:   mt.Radius,
        Alive:    true,
        HitFlash: 0,
    }
}

func (e *Enemy) Update(dt float64, playerPos Vec) {
    if !e.Alive {
        return
    }

    if e.HitFlash > 0 {
        e.HitFlash -= dt
        if e.HitFlash < 0 {
            e.HitFlash = 0
        }
    }

    dx := playerPos.X - e.Pos.X
    dy := playerPos.Y - e.Pos.Y
    dist := math.Hypot(dx, dy)

    if dist > 1 {
        e.Pos.X += (dx / dist) * e.Speed * dt
        e.Pos.Y += (dy / dist) * e.Speed * dt
    }
}
