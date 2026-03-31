package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

type DrawablePattern interface {
	Draw(screen *ebiten.Image, g *Game)
}

// PatternManager interface to avoid circular imports
type PatternManager interface {
	Update(g *Game)
	Reset()
	CurrentName() string
	Difficulty() int
	Current() BulletPattern
	Shuffle()
	Done() bool
}

type Game struct {
	world          *World
	Score          int
	patternManager PatternManager
	player         *Entity
	menuBGM        *audio.Player
	gameBGM        *audio.Player
	retryBGM       *audio.Player
	bgmStarted     bool
	fadeOut        bool
	fadeTimer      float64
	state          GameState
	fadeIn         float64
	frameCount     int
}

func NewGame(pm PatternManager) *Game {
	g := &Game{
		world:          NewWorld(),
		Score:          0,
		patternManager: pm,
		state:          StateMenu,
		fadeIn:         1.0,
	}

	g.initPlayer()

	menuTrack, err := loadBGM()
	if err == nil {
		g.menuBGM = menuTrack
		g.menuBGM.SetVolume(0.05)
		g.menuBGM.Play() // starts immediately (desktop OK)
	}

	gameTrack, err := loadRandomBGM()
	if err == nil {
		g.gameBGM = gameTrack
		g.gameBGM.SetVolume(0.05)
	}

	retryTrack, err := loadRetryBGM()
	if err == nil {
		g.retryBGM = retryTrack
		g.retryBGM.SetVolume(0.05)
	}

	return g
}

func (g *Game) initPlayer() {
	e := g.world.NewEntity()
	e.Position = &Position{X: ScreenWidth / 2, Y: ScreenHeight - 80}
	e.Velocity = &Velocity{}
	e.Hitbox = &Hitbox{Radius: 8}
	e.Player = &PlayerTag{}
}

func (g *Game) SpawnBullet(x, y, dx, dy float64, t BulletType) {
	e := g.world.NewEntity()
	e.Position = &Position{X: x, Y: y}
	e.Velocity = &Velocity{VX: dx, VY: dy}
	e.Hitbox = &Hitbox{Radius: 4}
	e.Bullet = &BulletTag{
		Speed: 4,
		Curve: 0,
		Type:  t,
	}
}

func (g *Game) Update() error {
	g.frameCount++
	switch g.state {
	case StateMenu:
		g.updateMenu()
	case StateControls:
		g.updateControls()
	case StatePlaying:
		g.updatePlaying()
	case StateRetry:
		g.updateRetry()
	case StateVictory:
		g.updateVictory()
	}
	return nil
}

func (g *Game) updateMenu() {
	if g.fadeIn > 0 {
		g.fadeIn -= 1.0 / 60.0
	}

	// ⭐ SAFETY: If menuBGM is nil, load it
	if g.menuBGM == nil {
		mbgm, err := loadBGM()
		if err == nil {
			g.menuBGM = mbgm
			g.menuBGM.SetVolume(0.1)
			g.menuBGM.Play()
		}
	}

	// ⭐ SAFETY: Only call IsPlaying/Rewind/Play if menuBGM exists
	if g.menuBGM != nil && !g.menuBGM.IsPlaying() {
		g.menuBGM.Rewind()
		g.menuBGM.Play()
	}

	// ENTER → Start game
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if g.menuBGM != nil {
			g.menuBGM.Pause()
		}

		g.reset(true, true) // load gameplay BGM + shuffle patterns

		if g.gameBGM != nil {
			g.gameBGM.Rewind()
			g.gameBGM.Play()
		}

		g.state = StatePlaying
		return
	}

	// C → Controls
	if ebiten.IsKeyPressed(ebiten.KeyC) {
		g.state = StateControls
		return
	}
}

func (g *Game) drawMenuParticles(screen *ebiten.Image) {
	for i := 0; i < 30; i++ {
		x := float32((i*53 + g.Score/4) % g.ScreenWidth())
		y := float32((i*97 + g.Score/5) % g.ScreenHeight())
		vector.FillCircle(screen, x, y, 1.2, color.RGBA{200, 200, 255, 60}, false)
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	// Background + particles
	g.drawBackground(screen)
	g.drawMenuParticles(screen)

	// Dim overlay
	dim := ebiten.NewImage(g.ScreenWidth(), g.ScreenHeight())
	dim.Fill(color.RGBA{0, 0, 0, 140})
	screen.DrawImage(dim, nil)

	// Common metrics
	m := ScoreFont.Metrics()
	lineSpacing := m.HAscent + m.HDescent + m.HLineGap

	// -------------------------
	// TITLE (Enhanced Neon FX)
	// -------------------------
	title := "BULLET HELL DEMO"
	titleW, _ := text.Measure(title, ScoreFont, lineSpacing)

	titleTarget := float64(g.ScreenWidth()) * 0.50
	titleScale := titleTarget / float64(titleW)

	if titleScale > 2.0 {
		titleScale = 2.0
	}
	if titleScale < 0.5 {
		titleScale = 0.5
	}

	scaledTitleW := float64(titleW) * titleScale
	titleX := int((float64(g.ScreenWidth()) - scaledTitleW) / 2)
	titleY := 150

	// --- Subtle breathing pulse ---
	p := pulse(g.frameCount, 0.02)
	titleScaleP := titleScale * (1.0 + 0.03*p)
	alpha := uint8(200 + 40*p)

	// --- RGB chromatic aberration ---
	neonText(screen, title, titleX-3, titleY, color.RGBA{255, 80, 80, 120}, titleScaleP) // red shadow
	neonText(screen, title, titleX+3, titleY, color.RGBA{80, 80, 255, 120}, titleScaleP) // blue shadow

	// --- Soft glow aura ---
	glow := ebiten.NewImage(int(scaledTitleW)+60, 120)
	glow.Fill(color.RGBA{80, 80, 140, 40})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(titleX)-30, float64(titleY)-40)
	screen.DrawImage(glow, op)

	// --- Main title ---
	neonText(screen, title, titleX, titleY, color.RGBA{180, 180, 255, alpha}, titleScaleP)

	// -------------------------
	// MENU OPTION 1 — Press ENTER to Start
	// -------------------------
	staticColor := color.RGBA{200, 200, 255, 255}

	start := "Press ENTER to Start"
	startW, _ := text.Measure(start, ScoreFont, lineSpacing)

	startTarget := float64(g.ScreenWidth()) * 0.40
	startScale := startTarget / float64(startW)

	if startScale > 1.5 {
		startScale = 1.5
	}
	if startScale < 0.7 {
		startScale = 0.7
	}

	scaledStartW := float64(startW) * startScale
	startX := int((float64(g.ScreenWidth()) - scaledStartW) / 2)
	yStart := 300

	left := "Press "
	key := "ENTER"
	right := " to Start"

	leftW, _ := text.Measure(left, ScoreFont, lineSpacing)
	leftPx := int(float64(leftW) * startScale)

	keyW, _ := text.Measure(key, ScoreFont, lineSpacing)
	keyPx := int(float64(keyW) * startScale)

	// Static left
	neonText(screen, left, startX, yStart, staticColor, startScale)

	// Pulse ENTER
	pEnter := pulse(g.frameCount, 0.08)
	keyScale := startScale * (1.0 + 0.10*pEnter)
	keyAlpha := uint8(180 + 75*pEnter)
	keyColor := color.RGBA{200, 200, 255, keyAlpha}

	neonTextChar(screen, key, startX+leftPx, yStart, keyColor, keyScale)

	// Static right
	neonText(screen, right, startX+leftPx+keyPx, yStart, staticColor, startScale)

	// -------------------------
	// MENU OPTION 2 — Press C for Controls
	// -------------------------
	controls := "Press C for Controls"
	controlsW, _ := text.Measure(controls, ScoreFont, lineSpacing)

	controlsTarget := float64(g.ScreenWidth()) * 0.40
	controlsScale := controlsTarget / float64(controlsW)

	if controlsScale > 1.5 {
		controlsScale = 1.5
	}
	if controlsScale < 0.7 {
		controlsScale = 0.7
	}

	scaledControlsW := float64(controlsW) * controlsScale
	controlsX := int((float64(g.ScreenWidth()) - scaledControlsW) / 2)
	yControls := 340

	left2 := "Press "
	key2 := "C"
	right2 := " for Controls"

	left2W, _ := text.Measure(left2, ScoreFont, lineSpacing)
	left2Px := int(float64(left2W) * controlsScale)

	key2W, _ := text.Measure(key2, ScoreFont, lineSpacing)
	key2Px := int(float64(key2W) * controlsScale)

	// Static left
	neonText(screen, left2, controlsX, yControls, staticColor, controlsScale)

	// Pulse C
	pC := pulse(g.frameCount+30, 0.08)
	key2Scale := controlsScale * (1.0 + 0.10*pC)
	key2Alpha := uint8(180 + 75*pC)
	key2Color := color.RGBA{200, 200, 255, key2Alpha}

	neonTextChar(screen, key2, controlsX+left2Px, yControls, key2Color, key2Scale)

	// Static right
	neonText(screen, right2, controlsX+left2Px+key2Px, yControls, staticColor, controlsScale)

	// Fade-in overlay
	if g.fadeIn > 0 {
		alpha := uint8(g.fadeIn * 255)
		fade := ebiten.NewImage(g.ScreenWidth(), g.ScreenHeight())
		fade.Fill(color.RGBA{0, 0, 0, alpha})
		screen.DrawImage(fade, nil)
	}

	// -------------------------
	// CREDITS — static line
	// -------------------------
	credits := "OST by SHIFT UP Corporation"
	creditsW, _ := text.Measure(credits, ScoreFont, lineSpacing)

	creditsTarget := float64(g.ScreenWidth()) * 0.38
	creditsScale := creditsTarget / float64(creditsW)

	if creditsScale > 1.3 {
		creditsScale = 1.3
	}
	if creditsScale < 0.6 {
		creditsScale = 0.6
	}

	scaledCreditsW := float64(creditsW) * creditsScale
	creditsX := int((float64(g.ScreenWidth()) - scaledCreditsW) / 2)
	yCredits := 420 // adjust if needed

	neonText(
		screen,
		credits,
		creditsX,
		yCredits,
		color.RGBA{160, 160, 220, 200}, // softer than menu options
		creditsScale,
	)
}

func (g *Game) updatePlaying() error {
	// If gameplay BGM finished, load a new random track
	if g.gameBGM != nil && !g.gameBGM.IsPlaying() {
		g.gameBGM.Pause()
		g.gameBGM = nil

		next, err := loadRandomBGM()
		if err == nil {
			g.gameBGM = next
			g.gameBGM.SetVolume(0.1)
			g.gameBGM.Play()
		}
	}

	g.updateInput()
	g.updateMovement()
	g.checkCollisions()

	g.Score += 1

	if g.patternManager != nil {
		g.patternManager.Update(g)
		if g.patternManager.Done() {
			g.state = StateVictory
			return nil
		}
	}

	// Fade-out logic (for transitions)
	if g.fadeOut && g.gameBGM != nil {
		dt := 1.0 / 60.0
		g.fadeTimer -= dt

		newVol := g.gameBGM.Volume() - (dt / 1.0)

		if newVol <= 0 {
			g.gameBGM.SetVolume(0)
			g.gameBGM.Pause()
			g.fadeOut = false
		} else {
			g.gameBGM.SetVolume(newVol)
		}
	}

	g.world.Cleanup()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateMenu:
		g.drawMenu(screen)

	case StatePlaying:
		g.drawPlaying(screen)

	case StateRetry:
		g.drawRetry(screen)

	case StateControls:
		g.drawControls(screen)

	case StateVictory:
		g.drawVictory(screen)
	}
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Base dark blue background
	screen.Fill(color.RGBA{10, 10, 30, 255})

	// Layer 1 — slow stars
	for i := 0; i < 60; i++ {
		x := float32((i*37 + g.Score/6) % g.ScreenWidth())
		y := float32((i*91 + g.Score/8) % g.ScreenHeight())
		vector.FillCircle(screen, x, y, 1, color.RGBA{180, 180, 255, 80}, false)
	}

	// Layer 2 — faster stars
	for i := 0; i < 40; i++ {
		x := float32((i*53 + g.Score/3) % g.ScreenWidth())
		y := float32((i*71 + g.Score/4) % g.ScreenHeight())
		vector.FillCircle(screen, x, y, 1.5, color.RGBA{220, 220, 255, 120}, false)
	}

	// Pattern visuals
	if dp, ok := g.patternManager.Current().(DrawablePattern); ok {
		dp.Draw(screen, g)
	}

	// Entities
	g.renderEntities(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) updateInput() {
	player := g.GetPlayer()
	if player == nil {
		return
	}

	speed := 3.0
	vx, vy := 0.0, 0.0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		vy -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		vy += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		vx -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		vx += speed
	}

	player.Velocity.VX = vx
	player.Velocity.VY = vy
}

func (g *Game) GetPlayer() *Entity {
	for _, e := range g.world.Entities() {
		if e.Player != nil {
			return e
		}
	}
	return nil
}

func (g *Game) ScreenWidth() int  { return ScreenWidth }
func (g *Game) ScreenHeight() int { return ScreenHeight }

func (g *Game) World() *World {
	return g.world
}

func (g *Game) NewBullet(x, y, vx, vy float64, t BulletType) *Entity {
	e := g.world.NewEntity()
	e.Position = &Position{X: x, Y: y}
	e.Velocity = &Velocity{VX: vx, VY: vy}
	e.Hitbox = &Hitbox{Radius: 4}
	e.Bullet = &BulletTag{
		Type:  t,
		Speed: 4,
		Curve: 0,
	}
	return e
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	// Base dark blue background
	screen.Fill(color.RGBA{10, 10, 30, 255})

	// Layer 1 — slow stars
	for i := 0; i < 60; i++ {
		x := float32((i*37 + g.Score/6) % g.ScreenWidth())
		y := float32((i*91 + g.Score/8) % g.ScreenHeight())
		vector.FillCircle(screen, x, y, 1, color.RGBA{180, 180, 255, 80}, false)
	}

	// Layer 2 — faster stars
	for i := 0; i < 40; i++ {
		x := float32((i*53 + g.Score/3) % g.ScreenWidth())
		y := float32((i*71 + g.Score/4) % g.ScreenHeight())
		vector.FillCircle(screen, x, y, 1.5, color.RGBA{220, 220, 255, 120}, false)
	}
}

func (g *Game) checkCollisions() {
	player := g.GetPlayer()
	if player == nil {
		return
	}

	for _, e := range g.world.Entities() {
		if e.Bullet == nil {
			continue
		}

		dx := e.Position.X - player.Position.X
		dy := e.Position.Y - player.Position.Y
		dist := math.Hypot(dx, dy)

		// Graze scoring
		grazeDist := e.Hitbox.Radius + player.Hitbox.Radius + 12
		if dist < grazeDist && dist >= e.Hitbox.Radius+player.Hitbox.Radius {
			g.Score += 5

			spark := g.world.NewEntity()
			spark.Position = &Position{
				X: player.Position.X,
				Y: player.Position.Y,
			}
			spark.Spark = &SparkTag{Life: 10}
		}

		// Collision
		if dist < e.Hitbox.Radius+player.Hitbox.Radius {
			if g.gameBGM != nil && !g.fadeOut {
				g.fadeOut = true
				g.fadeTimer = 1.0
			}

			g.fadeIn = 1.0
			g.state = StateRetry

			// stop gameplay music
			if g.gameBGM != nil {
				g.gameBGM.Pause()
			}

			// start retry music
			if g.retryBGM != nil {
				g.retryBGM.Rewind()
				g.retryBGM.Play()
			}
			return
		}
	}
}
