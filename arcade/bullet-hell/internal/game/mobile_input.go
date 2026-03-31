//go:build js && wasm

package game

import (
    "math"
    "strings"
    "syscall/js"

    "github.com/hajimehoshi/ebiten/v2"
)

type MobileInput struct {
    IsMobile bool

    // Joystick state
    JoyActive  bool
    JoyCenterX float64
    JoyCenterY float64
    JoyX       float64
    JoyY       float64

    // Buttons
    ShootPressed bool
    FocusPressed bool
    BombPressed  bool
}

func NewMobileInput() *MobileInput {
    return &MobileInput{
        IsMobile: detectMobile(),
    }
}

// --- Mobile Detection (WASM-safe) ---
func detectMobile() bool {
    ua := js.Global().Get("navigator").Get("userAgent").String()
    ua = strings.ToLower(ua)

    mobile := []string{
        "mobile", "android", "iphone", "ipad", "ipod",
        "webos", "opera mini", "iemobile",
    }

    for _, m := range mobile {
        if strings.Contains(ua, m) {
            return true
        }
    }
    return false
}

// --- Update Touch Input ---
func (m *MobileInput) Update(screenW, screenH int) {
    if !m.IsMobile {
        return
    }

    // Reset each frame
    m.ShootPressed = false
    m.FocusPressed = false
    m.BombPressed = false
    m.JoyActive = false
    m.JoyX, m.JoyY = 0, 0

    touches := ebiten.TouchIDs()
    for _, t := range touches {
        x, y := ebiten.TouchPosition(t)

        // LEFT SIDE → JOYSTICK
        if x < screenW/2 {
            if !m.JoyActive {
                m.JoyActive = true
                m.JoyCenterX = float64(x)
                m.JoyCenterY = float64(y)
            }

            dx := float64(x) - m.JoyCenterX
            dy := float64(y) - m.JoyCenterY
            dist := math.Hypot(dx, dy)

            maxDist := 60.0
            if dist > maxDist {
                scale := maxDist / dist
                dx *= scale
                dy *= scale
            }

            m.JoyX = dx / maxDist
            m.JoyY = dy / maxDist
        }

        // RIGHT SIDE → BUTTONS
        if x > screenW/2 {

            // SHOOT (bottom right)
            if x > screenW-160 && y > screenH-160 {
                m.ShootPressed = true
            }

            // FOCUS (middle right)
            if x > screenW-160 && y > screenH-300 && y < screenH-160 {
                m.FocusPressed = true
            }

            // BOMB (bottom right-left)
            if x > screenW-300 && x < screenW-160 && y > screenH-160 {
                m.BombPressed = true
            }
        }
    }
}