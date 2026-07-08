package game

import "github.com/hajimehoshi/ebiten/v2"

type MonsterType struct {
    Name   string
    Sprite *ebiten.Image
    HP     float64
    Speed  float64
    Radius float64
    Scale  float64 // NEW: per‑monster sprite size multiplier
}

var (
    MonsterJavascript *MonsterType
	MonsterTypescript *MonsterType
	MonsterCsharp *MonsterType

    MonsterPool []*MonsterType
)

func InitMonsters() {
    MonsterJavascript = &MonsterType{
        Name:   "Javascript",
        Sprite: LoadImage("assets/monsters/javascript_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,   // hitbox size
        Scale:  2.2,  // sprite size multiplier (make it bigger)
    }

	MonsterTypescript = &MonsterType{
        Name:   "Typescript",
        Sprite: LoadImage("assets/monsters/typescript_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,   // hitbox size
        Scale:  2.2,  // sprite size multiplier (make it bigger)
    }

	
	MonsterCsharp = &MonsterType{
        Name:   "Csharp",
        Sprite: LoadImage("assets/monsters/csharp_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,   // hitbox size
        Scale:  2.2,  // sprite size multiplier (make it bigger)
    }



    MonsterPool = []*MonsterType{
        MonsterJavascript,
		MonsterTypescript,
		MonsterCsharp,
    }
}
