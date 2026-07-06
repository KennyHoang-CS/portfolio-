package ui

import (
    "image/color"
    "sort"
    "strings"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/text/v2"

    fonts "github.com/KennyHoang-CS/portfolio/bazelquest/internal/assets"
    "github.com/KennyHoang-CS/portfolio/bazelquest/internal/workspace"
)

type FileSystemPanel struct {
    ws       *workspace.Workspace
    x, y     int
    width    int
    height   int
    onOpen   func(path string)
}

func NewFileSystemPanel(ws *workspace.Workspace, x, y, w, h int, onOpen func(string)) *FileSystemPanel {
    return &FileSystemPanel{
        ws:     ws,
        x:      x,
        y:      y,
        width:  w,
        height: h,
        onOpen: onOpen,
    }
}

func (p *FileSystemPanel) Draw(screen *ebiten.Image) {
    files := p.ws.List()
    sort.Strings(files)

    y := p.y + 10
    lineHeight := int(fonts.EditorFont.Metrics().HAscent +
        fonts.EditorFont.Metrics().HDescent +
        fonts.EditorFont.Metrics().HLineGap)

    for _, path := range files {
        indent := strings.Count(path, "/") * 12

        opts := &text.DrawOptions{}
        opts.GeoM.Translate(float64(p.x+indent), float64(y))
        opts.ColorScale.ScaleWithColor(color.RGBA{200, 200, 200, 255})

        text.Draw(screen, path, fonts.EditorFont, opts)
        y += lineHeight
    }
}

func (p *FileSystemPanel) Click(x, y int) {
    if x < p.x || x > p.x+p.width {
        return
    }
    if y < p.y || y > p.y+p.height {
        return
    }

    files := p.ws.List()
    sort.Strings(files)

    lineHeight := int(fonts.EditorFont.Metrics().HAscent +
        fonts.EditorFont.Metrics().HDescent +
        fonts.EditorFont.Metrics().HLineGap)

    index := (y - p.y - 10) / lineHeight
    if index < 0 || index >= len(files) {
        return
    }

    p.onOpen(files[index])
}