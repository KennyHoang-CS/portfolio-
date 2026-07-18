package game

import (
    "image"
    "image/color"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
    "github.com/hajimehoshi/ebiten/v2/vector"
)

const (
    headerHeight = 80
    cellWidth    = 720
    cellHeight   = 120
)

func (g *Game) pointingDevicePosition() (x, y int) {
    touchIDs := ebiten.AppendTouchIDs(nil)
    if len(touchIDs) > 0 {
        return ebiten.TouchPosition(touchIDs[0])
    }
    return ebiten.CursorPosition()
}

func (g *Game) isPointingDevicePressed() bool {
    touchIDs := ebiten.AppendTouchIDs(nil)
    if len(touchIDs) > 0 {
        return true
    }
    return ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

func (g *Game) isPointingDeviceJustReleased() bool {
    touchIDs := inpututil.AppendJustReleasedTouchIDs(nil)
    if len(touchIDs) > 0 {
        return true
    }
    return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

func (g *Game) updateToggleUI() {
    // Open/close
    if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
        g.ToggleUIOpen = true
    }
    if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
        g.ToggleUIOpen = false
    }

    // Mobile button
    touchIDs := ebiten.AppendTouchIDs(nil)
    for _, id := range touchIDs {
        tx, ty := ebiten.TouchPosition(id)
        if tx >= 10 && tx <= 150 && ty >= 540 && ty <= 590 {
            g.ToggleUIOpen = true
        }
    }

    if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
        mx, my := ebiten.CursorPosition()
        g.handleToggleTouch(mx, my)
    }

    if !g.ToggleUIOpen {
        return
    }

    // Define scrollable content area
    g.contentArea = image.Rect(
        40,
        80+headerHeight,
        40+cellWidth,
        520,
    )

    x, y := g.pointingDevicePosition()
    hovering := image.Pt(x, y).In(g.contentArea)

    // Wheel scroll
    if _, wheelY := ebiten.Wheel(); wheelY != 0 && hovering {
        g.ToggleVelocityY = int(wheelY * -15)
    }

    // Release → inertia
    if g.isPointingDeviceJustReleased() && g.ToggleDragging {
        g.ToggleDragging = false
        g.ToggleVelocityY = y - g.TogglePrevY
        g.TogglePrevY = y
        return
    }

    // Inertia
    if !g.isPointingDevicePressed() {
        g.ToggleDragging = false

        g.setInfiniteOffsetY(g.ToggleOffsetY + g.ToggleVelocityY)
        if g.ToggleVelocityY != 0 {
            g.ToggleVelocityY = int(float64(g.ToggleVelocityY) * 15.0 / 16.0)
        }
        g.TogglePrevY = y
        return
    }

    // Stop inertia
    g.ToggleVelocityY = 0

    // Start dragging
    if !g.ToggleDragging && hovering {
        g.ToggleDragging = true
        g.ToggleOffsetStartY = g.ToggleOffsetY
        g.ToggleStartY = y
    }

    // Dragging
    if g.ToggleDragging {
        g.setInfiniteOffsetY(g.ToggleOffsetStartY + y - g.ToggleStartY)
    }

    g.TogglePrevY = y
}

func (g *Game) setInfiniteOffsetY(offsetY int) {
    contentHeight := len(g.AvailableAbilities) * cellHeight
    if contentHeight == 0 {
        return
    }

    // Wrap offset like infinite scroll demo
    offsetY %= contentHeight
    if offsetY < 0 {
        offsetY += contentHeight
    }

    g.ToggleOffsetY = offsetY
}

func (g *Game) handleToggleTouch(x, y int) {
    if x >= 300 && x <= 500 && y >= 520 && y <= 560 {
        g.ToggleUIOpen = false
        return
    }

    if !g.ToggleUIOpen {
        return
    }

    if !image.Pt(x, y).In(g.contentArea) {
        return
    }

    count := len(g.AvailableAbilities)
    if count == 0 {
        return
    }

    for drawIndex := 0; drawIndex < count; drawIndex++ {
        abilityIndex := (drawIndex + (g.ToggleOffsetY / cellHeight)) % count
        a := g.AvailableAbilities[abilityIndex]

        itemY := g.contentArea.Min.Y + drawIndex*cellHeight
        itemX := g.contentArea.Min.X

        itemRect := image.Rect(
            itemX,
            itemY,
            itemX+cellWidth,
            itemY+cellHeight,
        )

        if itemRect.Intersect(g.contentArea).Empty() {
            continue
        }

        boxX := itemRect.Min.X + 10
        boxY := itemRect.Min.Y + 10

        if x >= boxX && x <= boxX+40 && y >= boxY && y <= boxY+40 {
            a.Enabled = !a.Enabled
            if a.Enabled {
                a.Apply(g)
            }
            return
        }
    }
}

func (g *Game) drawToggleUI(screen *ebiten.Image) {
    // Mobile button
    vector.FillRect(screen, 10, 540, 140, 50, color.RGBA{20, 20, 40, 200}, false)
    ebitenutil.DebugPrintAt(screen, "Abilities", 20, 555)

    if !g.ToggleUIOpen {
        return
    }

    // Neon scanlines
    for y := 80; y < 600; y += 4 {
        vector.FillRect(screen, 0, float32(y), 800, 1, color.RGBA{0, 255, 255, 20}, false)
    }

    // Background panel
    vector.FillRect(screen, 0, 80, 800, 520, color.RGBA{4, 4, 12, 240}, false)
    vector.StrokeRect(screen, 0, 80, 800, 520, 3, color.RGBA{0, 255, 255, 255}, false)

    // Header box
    vector.FillRect(screen, 0, 80, 800, headerHeight, color.RGBA{6, 6, 20, 240}, false)
    vector.StrokeRect(screen, 0, 80, 800, headerHeight, 3, color.RGBA{0, 255, 255, 255}, false)
    ebitenutil.DebugPrintAt(screen, "TRAINING ROOM — Ability Toggles", 20, 110)

    // Content mask
    screenContent := screen.SubImage(g.contentArea).(*ebiten.Image)

    count := len(g.AvailableAbilities)
    if count == 0 {
        return
    }

    for drawIndex := 0; drawIndex < count; drawIndex++ {
        abilityIndex := (drawIndex + (g.ToggleOffsetY / cellHeight)) % count
        a := g.AvailableAbilities[abilityIndex]

        itemY := g.contentArea.Min.Y + drawIndex*cellHeight
        itemX := g.contentArea.Min.X

        itemRect := image.Rect(
            itemX,
            itemY,
            itemX+cellWidth,
            itemY+cellHeight,
        )

        if itemRect.Intersect(g.contentArea).Empty() {
            continue
        }

        // Neon separator
        vector.StrokeLine(screenContent,
            float32(itemRect.Min.X), float32(itemRect.Min.Y),
            float32(itemRect.Max.X), float32(itemRect.Min.Y),
            1, color.RGBA{0, 200, 255, 255}, false)

        // Checkbox
        vector.FillRect(screenContent, float32(itemRect.Min.X+10), float32(itemRect.Min.Y+10), 40, 40,
            func() color.RGBA {
                if a.Enabled {
                    return color.RGBA{0, 255, 255, 255}
                }
                return color.RGBA{20, 40, 60, 255}
            }(), false)

        if a.Enabled {
            vector.StrokeRect(screenContent, float32(itemRect.Min.X+10), float32(itemRect.Min.Y+10), 40, 40,
                2, color.RGBA{0, 255, 255, 255}, false)
        }

        // Icon
        if a.Icon != nil {
            bounds := a.Icon.Bounds()
            h := bounds.Dy()
            scale := 40.0 / float64(h)

            op := &ebiten.DrawImageOptions{}
            op.GeoM.Scale(scale, scale)
            op.GeoM.Translate(float64(itemRect.Min.X+70), float64(itemRect.Min.Y+10))
            screenContent.DrawImage(a.Icon, op)
        }

        // Name
        ebitenutil.DebugPrintAt(screenContent, a.Name, itemRect.Min.X+130, itemRect.Min.Y+15)

        // Description panel
        descX := itemRect.Min.X + 130
        descY := itemRect.Min.Y + 45
        descW := cellWidth - 160
        descH := cellHeight - 55

        vector.FillRect(screenContent, float32(descX), float32(descY), float32(descW), float32(descH),
            color.RGBA{6, 6, 20, 220}, false)

        lines := wrapTextWords(a.Description, 60)
        lineHeight := 16

        maxDescScroll := float64(len(lines)*lineHeight - descH)
        if maxDescScroll < 0 {
            maxDescScroll = 0
        }
        if a.DescScroll < 0 {
            a.DescScroll = 0
        }
        if a.DescScroll > maxDescScroll {
            a.DescScroll = maxDescScroll
        }

        startLine := int(a.DescScroll / float64(lineHeight))
        if startLine < 0 {
            startLine = 0
        }
        if startLine >= len(lines) {
            startLine = len(lines) - 1
        }

        for li := startLine; li < len(lines); li++ {
            lineY := descY + (li-startLine)*lineHeight
            if lineY > descY+descH {
                break
            }
            ebitenutil.DebugPrintAt(screenContent, lines[li], descX+8, lineY)
        }

        // Scrollbar
        if maxDescScroll > 0 {
            scrollRatio := a.DescScroll / maxDescScroll
            barHeight := float32(descH) * 0.3
            barY := float32(descY) + (float32(descH)-barHeight)*float32(scrollRatio)
            barX := float32(descX + descW - 6)

            vector.FillRect(screenContent, barX, float32(descY), 6, float32(descH),
                color.RGBA{0, 80, 120, 255}, false)
            vector.FillRect(screenContent, barX, barY, 6, barHeight,
                color.RGBA{0, 255, 255, 255}, false)
        }
    }

    // Close button
    vector.FillRect(screen, 300, 520, 200, 40, color.RGBA{0, 80, 120, 255}, false)
    vector.StrokeRect(screen, 300, 520, 200, 40, 2, color.RGBA{0, 255, 255, 255}, false)
    ebitenutil.DebugPrintAt(screen, "Close", 360, 530)
}
