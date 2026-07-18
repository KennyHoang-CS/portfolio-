package game

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type Ability struct {
	Name            string
	Description     string
	Level           int
	Icon            *ebiten.Image
	Enabled         bool
	Apply           func(g *Game) // runs every frame if Enabled
	DescScroll      float64
	DescExpanded    bool
	DescHover       bool
	DescTouchY      int
	DescTouchActive bool
}

type DamageNumber struct {
	X, Y  float64
	Value int
	Life  float64
}

type Crystal struct {
	Pos    Vec
	Alive  bool
	Sprite *ebiten.Image
}

type Projectile struct {
	Pos    Vec
	Vel    Vec
	Damage float64
	Alive  bool
	Sprite *ebiten.Image

	// Behavior flags
	Lifetime int     // frames until despawn
	Curve    float64 // boomerang arc
	Return   bool    // flips direction mid-flight
	Homing   bool    // seeks nearest enemy
	Orbit    bool    // circles around player
	Radius   float64 // orbit radius
	Angle    float64 // orbit angle
	Speed    float64 // movement speed

	// Optional callbacks
	OnHit    func(*Enemy)
	OnUpdate func(*Projectile, *Game)
}

type FireTile struct {
	Pos    Vec
	Life   float64
	Damage float64
}

type StunWave struct {
	Pos    Vec
	Radius float64
	Life   float64
}

type Game struct {
	player        *Player
	lastTime      time.Time
	Enemies       []*Enemy
	DamageNumbers []*DamageNumber
	Crystals      []*Crystal

	waveTimer  float64
	waveNumber int

	bg             *ebiten.Image
	trainingRoomBg *ebiten.Image
	grassBg        *ebiten.Image
	crystalSprite  *ebiten.Image

	// Abilities
	AvailableAbilities []*Ability
	LevelUpMenuOpen    bool
	LevelUpChoices     []*Ability

	// Unified projectile system
	Projectiles []*Projectile

	// Fire tiles & waves remain separate
	FireTiles []*FireTile
	StunWaves []*StunWave

	// Global spell stats
	CrystalMagnetRadius float64
	InfiniteLoopChance  float64

	// Dagger (now projectile)
	DaggerCooldown float64
	DaggerRate     float64

	GameState     GameState
	TitleBG       *ebiten.Image
	FadeAlpha     float64
	BlinkTimer    float64
	TitleFont     font.Face
	TitleFontSize float64
	TitleYOffset  float64

	inTrainingRoom bool

	abilityEnabled map[string]bool

	ToggleUIOpen    bool
	ScrollOffset    float64
	LastTouchY      int
	LastTouchActive bool

	contentArea        image.Rectangle
	ToggleDragging     bool
	TogglePrevY        int
	ToggleStartY       int
	ToggleOffsetY      int
	ToggleOffsetStartY int
	ToggleVelocityY    int
}