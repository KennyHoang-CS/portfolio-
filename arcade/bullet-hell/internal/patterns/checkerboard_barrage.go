package patterns

import (
	"github.com/KennyHoang-CS/portfolio/bullet-hell/internal/game"
)

type CheckerboardBarragePattern struct {
	timer    int
	duration int
	offset   int
}

func NewCheckerboardBarragePattern() *CheckerboardBarragePattern {
	return &CheckerboardBarragePattern{
		duration: 260,
	}
}

func (p *CheckerboardBarragePattern) Update(g *game.Game, difficulty int) {
    p.timer++

    if p.timer%20 == 0 {
        cols := 10 + difficulty
        rows := 6 + difficulty/2

        cellW := float64(g.ScreenWidth()) / float64(cols)
        cellH := float64(g.ScreenHeight()) / float64(rows)

        p.offset++
        for r := 0; r < rows; r++ {
            for c := 0; c < cols; c++ {
                if (r+c+p.offset)%2 == 0 {
                    x := cellW*float64(c) + cellW/2
                    y := cellH*float64(r) + cellH/2

                    g.NewBullet(
                        x, y,
                        0,
                        1.4+float64(difficulty)*0.04,
                        game.BulletStar,   
                    )
                }
            }
        }
    }
}

func (p *CheckerboardBarragePattern) IsFinished() bool {
	return p.timer > p.duration
}

func (p *CheckerboardBarragePattern) Reset() {
	p.timer = 0
	p.offset = 0
}

func (p *CheckerboardBarragePattern) Name() string { return "CheckerboardBarrage" }