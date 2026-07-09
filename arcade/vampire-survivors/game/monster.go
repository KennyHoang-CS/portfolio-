package game

import "github.com/hajimehoshi/ebiten/v2"

type MonsterType struct {
    Name   string
    Sprite *ebiten.Image
    HP     float64
    Speed  float64
    Radius float64
    Scale  float64 // sprite size multiplier
}

var (
    MonsterJavascript *MonsterType
    MonsterTypescript *MonsterType
    MonsterCsharp     *MonsterType
    MonsterPhp        *MonsterType
    MonsterPython     *MonsterType
    MonsterRuby       *MonsterType
    MonsterSwift      *MonsterType
    MonsterJava       *MonsterType
    MonsterRust       *MonsterType

    MonsterPool []*MonsterType
)

func InitMonsters() {

    MonsterJavascript = &MonsterType{
        Name:   "Javascript",
        Sprite: LoadImage("assets/monsters/javascript_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,
        Scale:  2.2,
    }

    MonsterTypescript = &MonsterType{
        Name:   "Typescript",
        Sprite: LoadImage("assets/monsters/typescript_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,
        Scale:  2.2,
    }

    MonsterPhp = &MonsterType{
        Name:   "Php",
        Sprite: LoadImage("assets/monsters/php_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,
        Scale:  2.2,
    }

    MonsterPython = &MonsterType{
        Name:   "Python",
        Sprite: LoadImage("assets/monsters/python_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,
        Scale:  2.2,
    }

    MonsterRuby = &MonsterType{
        Name:   "Ruby",
        Sprite: LoadImage("assets/monsters/ruby_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,
        Scale:  2.2,
    }

    MonsterRust = &MonsterType{
        Name:   "Rust",
        Sprite: LoadImage("assets/monsters/rust_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,
        Scale:  2.2,
    }

    MonsterSwift = &MonsterType{
        Name:   "Swift",
        Sprite: LoadImage("assets/monsters/swift_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,
        Scale:  2.2,
    }

    MonsterCsharp = &MonsterType{
        Name:   "Csharp",
        Sprite: LoadImage("assets/monsters/csharp_monster.png"),
        HP:     12,
        Speed:  55,
        Radius: 22,
        Scale:  2.2,
    }

    // ⭐ FIXED: Java was missing, causing nil pointer crash
    MonsterJava = &MonsterType{
        Name:   "Java",
        Sprite: LoadImage("assets/monsters/java_monster.png"),
        HP:     14,
        Speed:  50,
        Radius: 24,
        Scale:  2.3,
    }

    MonsterPool = []*MonsterType{
        MonsterJavascript,
        MonsterTypescript,
        MonsterCsharp,
        MonsterJava,     // now valid
        MonsterPhp,
        MonsterRust,
        MonsterRuby,
        MonsterSwift,
        MonsterPython,
    }
}
