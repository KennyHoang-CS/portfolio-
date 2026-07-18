package game

import (
    "image/color"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
    "github.com/hajimehoshi/ebiten/v2/vector"
)

const (
    GridTop    = 160
    CellWidth  = 260
    CellHeight = 100
    Columns    = 3
)

func (g *Game) updateToggleUI() {
    // Desktop: TAB opens UI
    if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
        g.ToggleUIOpen = true
    }

    // Desktop: ESC closes UI
    if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
        g.ToggleUIOpen = false
    }

    // ---------------------------------------------------------
    // Mobile: tap the Abilities button
    // ---------------------------------------------------------
    var ids []ebiten.TouchID
    ids = ebiten.AppendTouchIDs(ids)

    for _, id := range ids {
        tx, ty := ebiten.TouchPosition(id)
        if tx >= 10 && tx <= 150 && ty >= 540 && ty <= 590 {
            g.ToggleUIOpen = true
        }
    }

    // ---------------------------------------------------------
    // Desktop: single-click toggle (smooth)
    // ---------------------------------------------------------
    if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
        mx, my := ebiten.CursorPosition()
        g.handleToggleTouch(mx, my)
    }

    if !g.ToggleUIOpen {
        return
    }

    // ---------------------------------------------------------
    // SCROLLING (desktop wheel)
    // ---------------------------------------------------------
    _, wheelY := ebiten.Wheel()
    g.ScrollOffset += wheelY * -20

    // ---------------------------------------------------------
    // SCROLLING (mobile drag)
    // ---------------------------------------------------------
    ids = ids[:0]
    ids = ebiten.AppendTouchIDs(ids)

    if len(ids) > 0 {
        id := ids[0]
        _, ty := ebiten.TouchPosition(id)

        if g.LastTouchActive {
            dy := float64(ty - g.LastTouchY)
            g.ScrollOffset += dy
        }

        g.LastTouchY = ty
        g.LastTouchActive = true
    } else {
        g.LastTouchActive = false
    }

    // ---------------------------------------------------------
    // Clamp scroll offset (with bottom padding)
    // ---------------------------------------------------------
    rows := (len(g.AvailableAbilities) + Columns - 1) / Columns
    maxScroll := float64(rows*CellHeight - (520 - GridTop) + 20)

    if g.ScrollOffset < 0 {
        g.ScrollOffset = 0
    }
    if g.ScrollOffset > maxScroll {
        g.ScrollOffset = maxScroll
    }

    // ---------------------------------------------------------
    // Mobile: handle toggles (touch)
    // ---------------------------------------------------------
    for _, id := range ids {
        tx, ty := ebiten.TouchPosition(id)
        g.handleToggleTouch(tx, ty)
    }
}

// ---------------------------------------------------------
// Handle checkbox clicks / taps
// ---------------------------------------------------------

func (g *Game) handleToggleTouch(x, y int) {
    // Close button
    if x >= 300 && x <= 500 && y >= 520 && y <= 560 {
        g.ToggleUIOpen = false
        return
    }

    startY := GridTop - int(g.ScrollOffset)

    for i, a := range g.AvailableAbilities {
        row := i / Columns
        col := i % Columns

        boxX := 40 + col*CellWidth + 10
        boxY := startY + row*CellHeight + 10

        if x >= boxX && x <= boxX+40 && y >= boxY && y <= boxY+40 {
            a.Enabled = !a.Enabled
            if a.Enabled {
                a.Apply(g)
            }
            return
        }
    }
}

// ---------------------------------------------------------
// DRAW: Desktop + Mobile Unified UI
// ---------------------------------------------------------

func (g *Game) drawToggleUI(screen *ebiten.Image) {
    // Mobile button (always visible)
    vector.FillRect(screen, 10, 540, 140, 50, color.RGBA{20, 20, 40, 200}, false)
    ebitenutil.DebugPrintAt(screen, "Abilities", 20, 555)

    if !g.ToggleUIOpen {
        return
    }

    // Background panel
    vector.FillRect(screen, 0, 80, 800, 520, color.RGBA{10, 10, 30, 230}, false)
    vector.StrokeRect(screen, 0, 80, 800, 520, 3, color.RGBA{200, 200, 255, 255}, false)

    // Header (fixed)
    ebitenutil.DebugPrintAt(screen, "TRAINING ROOM — Ability Toggles", 20, 100)

    // Draw abilities in a 3-column grid
    startY := GridTop - int(g.ScrollOffset)

    for i, a := range g.AvailableAbilities {
        row := i / Columns
        col := i % Columns

        y := startY + row*CellHeight
        x := 40 + col*CellWidth

        // Skip off-screen rows
        if y < GridTop || y > 520 {
            continue
        }

        // Checkbox
        if a.Enabled {
            vector.FillRect(screen, float32(x+10), float32(y+10), 40, 40, color.RGBA{80, 200, 255, 255}, false)
        } else {
            vector.FillRect(screen, float32(x+10), float32(y+10), 40, 40, color.RGBA{40, 40, 80, 255}, false)
        }

        // Icon
        if a.Icon != nil {
            bounds := a.Icon.Bounds()
            h := bounds.Dy()
            scale := 40.0 / float64(h)

            op := &ebiten.DrawImageOptions{}
            op.GeoM.Scale(scale, scale)
            op.GeoM.Translate(float64(x+60), float64(y+10))
            screen.DrawImage(a.Icon, op)
        }

        // Name
        ebitenutil.DebugPrintAt(screen, a.Name, x+120, y+20)
    }

    // Close button
    vector.FillRect(screen, 300, 520, 200, 40, color.RGBA{80, 80, 255, 255}, false)
    ebitenutil.DebugPrintAt(screen, "Close", 360, 530)
}
