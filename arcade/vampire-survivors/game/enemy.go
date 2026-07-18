package game

import (
	"math"
	"math/rand"
)

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

    // ⭐ Training dummy check: if Speed == 0, do NOT move
    if e.Speed <= 0 {
        return
    }

    // Normal monster movement
    dx := playerPos.X - e.Pos.X
    dy := playerPos.Y - e.Pos.Y
    dist := math.Hypot(dx, dy)

    if dist > 1 {
        e.Pos.X += (dx / dist) * e.Speed * dt
        e.Pos.Y += (dy / dist) * e.Speed * dt
    }
}

func (g *Game) updateWaves(dt float64) {
	g.waveTimer -= dt
	if g.waveTimer <= 0 {
		g.spawnWave()
		g.waveNumber++
		g.waveTimer = 5.0
	}
}

func (g *Game) spawnWave() {
	count := 5 + g.waveNumber*2

	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2
		dist := 400 + rand.Float64()*200

		x := g.player.Pos.X + math.Cos(angle)*dist
		y := g.player.Pos.Y + math.Sin(angle)*dist

		mt := MonsterPool[rand.Intn(len(MonsterPool))]
		g.Enemies = append(g.Enemies, NewEnemy(mt, x, y))
	}
}