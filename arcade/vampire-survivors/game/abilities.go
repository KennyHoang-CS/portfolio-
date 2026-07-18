package game

import (
    "math"
    "math/rand"

    "github.com/hajimehoshi/ebiten/v2"
)

type Ability struct {
    Name        string
    Description string
    Level       int
    Icon        *ebiten.Image
    Enabled     bool
    Apply       func(g *Game) // runs every frame if Enabled
}

// ---------------------------------------------------------
// Ability Toggle Helpers
// ---------------------------------------------------------

func (g *Game) AbilityEnabled(name string) bool {
    for _, a := range g.AvailableAbilities {
        if a.Name == name {
            return a.Enabled
        }
    }
    return false
}

func (g *Game) ToggleAbility(name string) {
    for _, a := range g.AvailableAbilities {
        if a.Name == name {
            a.Enabled = !a.Enabled
            return
        }
    }
}

// ---------------------------------------------------------
// Level-Up System (unchanged)
// ---------------------------------------------------------

func (g *Game) pickAbility(idx int) {
    if idx < 0 || idx >= len(g.LevelUpChoices) {
        return
    }
    ability := g.LevelUpChoices[idx]
    ability.Level++
    ability.Enabled = true // level-up enables the ability
    g.player.Abilities = append(g.player.Abilities, ability)
}

// ---------------------------------------------------------
// Initialize all abilities (no toggling logic here)
// ---------------------------------------------------------

func (g *Game) initAbilities() {
    g.AvailableAbilities = []*Ability{

        // -------------------------------------------------
        // CPU Overclock — slash buffs
        // -------------------------------------------------
        {
            Name:        "CPU Overclock",
            Description: "Increase slash damage, radius, and attack speed.",
            Icon:        LoadImage("assets/skills/cpu_overclock_skill.png"),
            Apply: func(g *Game) {
                g.player.SlashDamage += 1
                g.player.SlashCooldownBase *= 0.9
                g.player.SlashRadius += 5
                g.player.AttackSpeedMultiplier += 0.05
            },
        },

        // -------------------------------------------------
        // Garbage Collection — crystal magnet
        // -------------------------------------------------
        {
            Name:        "Garbage Collection",
            Description: "Increase crystal magnet radius.",
            Icon:        LoadImage("assets/skills/garbage_collection_skill.png"),
            Apply: func(g *Game) {
                g.CrystalMagnetRadius += 40
            },
        },

        // -------------------------------------------------
        // Type Error — placeholder aura
        // -------------------------------------------------
        {
            Name:        "Type Error",
            Description: "Slow aura around the player. (placeholder)",
            Icon:        LoadImage("assets/skills/type_error_skill.png"),
            Apply: func(g *Game) {
                // Implement later
            },
        },

        // -------------------------------------------------
        // Callback Hell — boomerang projectile
        // -------------------------------------------------
        {
            Name:        "Callback Hell",
            Description: "Throw a chaotic boomerang.",
            Icon:        LoadImage("assets/skills/callback_hell_skill.png"),
            Apply: func(g *Game) {
                dir := float64(g.player.lastDir)

                p := &Projectile{
                    Pos:      g.player.Pos,
                    Vel:      Vec{X: dir * 350, Y: 0},
                    Damage:   8 + float64(g.player.Level)*2,
                    Alive:    true,
                    Sprite:   PlaceholderIcon(255, 150, 0),
                    Lifetime: 45,
                    Curve:    0.05,
                    Return:   true,
                    Speed:    350,
                }

                g.Projectiles = append(g.Projectiles, p)
            },
        },

        // -------------------------------------------------
        // Index Out of Range — unlock dagger auto-fire
        // -------------------------------------------------
        {
            Name:        "Index Out of Range",
            Description: "Unlock auto‑firing dagger and increase fire rate.",
            Icon:        LoadImage("assets/skills/index_out_of_error_skill.png"),
            Apply: func(g *Game) {
                g.player.HasDagger = true
                g.DaggerRate *= 0.60
            },
        },

        // -------------------------------------------------
        // Memory Leak — fire tiles
        // -------------------------------------------------
        {
            Name:        "Memory Leak",
            Description: "Leave fire tiles behind you.",
            Icon:        LoadImage("assets/skills/memory_leak_skill.png"),
            Apply: func(g *Game) {
                ft := &FireTile{
                    Pos:    g.player.Pos,
                    Life:   3.0,
                    Damage: 5 + float64(g.player.Level),
                }
                g.FireTiles = append(g.FireTiles, ft)
            },
        },

        // -------------------------------------------------
        // Race Condition — movement + attack speed
        // -------------------------------------------------
        {
            Name:        "Race Condition",
            Description: "Increase movement and attack speed.",
            Icon:        LoadImage("assets/skills/race_condition_skill.png"),
            Apply: func(g *Game) {
                g.player.MoveSpeedMultiplier += 0.1
                g.player.AttackSpeedMultiplier += 0.05
            },
        },

        // -------------------------------------------------
        // Deadlock — freeze aura placeholder
        // -------------------------------------------------
        {
            Name:        "Deadlock",
            Description: "Freeze aura. (placeholder)",
            Icon:        LoadImage("assets/skills/deadlock_skill.png"),
            Apply: func(g *Game) {
                // Implement later
            },
        },

        // -------------------------------------------------
        // Stack Overflow — HP + regen
        // -------------------------------------------------
        {
            Name:        "Stack Overflow",
            Description: "Increase HP and regen.",
            Icon:        LoadImage("assets/skills/stack_overflow_skill.png"),
            Apply: func(g *Game) {
                g.player.MaxHP += 20
                g.player.HP += 20
                g.player.RegenPerSecond += 1
            },
        },

        // -------------------------------------------------
        // Heap Allocation — orbiting orb
        // -------------------------------------------------
        {
            Name:        "Heap Allocation",
            Description: "Summon an orbiting orb.",
            Icon:        LoadImage("assets/skills/heap_allocation_skill.png"),
            Apply: func(g *Game) {
                orb := &Projectile{
                    Pos:    g.player.Pos,
                    Alive:  true,
                    Damage: 5,
                    Sprite: PlaceholderIcon(0, 255, 0),
                    Orbit:  true,
                    Angle:  rand.Float64() * math.Pi * 2,
                    Radius: 60,
                    Speed:  1.5,
                }
                g.Projectiles = append(g.Projectiles, orb)
            },
        },

        // -------------------------------------------------
        // Compiler Optimization — global speed buffs
        // -------------------------------------------------
        {
            Name:        "Compiler Optimization",
            Description: "Increase movement and attack speed.",
            Icon:        LoadImage("assets/skills/compiler_optimization_skill.png"),
            Apply: func(g *Game) {
                g.player.MoveSpeedMultiplier += 0.05
                g.player.AttackSpeedMultiplier += 0.05
            },
        },

        // -------------------------------------------------
        // AI Agents — homing minions
        // -------------------------------------------------
        {
            Name:        "AI Agents",
            Description: "Summon autonomous homing agents.",
            Icon:        LoadImage("assets/skills/ai_agents_skill.png"),
            Apply: func(g *Game) {
                for i := 0; i < 2; i++ {
                    agent := &Projectile{
                        Pos:    g.player.Pos,
                        Alive:  true,
                        Damage: 6 + float64(g.player.Level),
                        Sprite: PlaceholderIcon(120, 255, 255),
                        Homing: true,
                        Speed:  250,
                        Vel:    Vec{X: rand.Float64()*2 - 1, Y: rand.Float64()*2 - 1},
                    }
                    g.Projectiles = append(g.Projectiles, agent)
                }
            },
        },

        // -------------------------------------------------
        // Dependency Injection — forward projectile
        // -------------------------------------------------
        {
            Name:        "Dependency Injection",
            Description: "Fire a green projectile forward.",
            Icon:        LoadImage("assets/skills/dependency_injection_skill.png"),
            Apply: func(g *Game) {
                dir := float64(g.player.lastDir)

                p := &Projectile{
                    Pos:      g.player.Pos,
                    Vel:      Vec{X: dir * 300, Y: 0},
                    Damage:   12 + float64(g.player.Level)*2.5,
                    Alive:    true,
                    Sprite:   PlaceholderIcon(0, 200, 0),
                    Speed:    300,
                    Lifetime: 60,
                }

                g.Projectiles = append(g.Projectiles, p)
                g.DaggerRate *= 0.85
            },
        },

        // -------------------------------------------------
        // Merge Conflict — 360° cleave
        // -------------------------------------------------
        {
            Name:        "Merge Conflict",
            Description: "Perform a 360° cleave.",
            Icon:        LoadImage("assets/skills/merge_conflict_skill.png"),
            Apply: func(g *Game) {
                radius := 90.0
                damage := 12 + float64(g.player.Level)*2.5

                for _, e := range g.Enemies {
                    if !e.Alive {
                        continue
                    }

                    dx := e.Pos.X - g.player.Pos.X
                    dy := e.Pos.Y - g.player.Pos.Y
                    dist := math.Hypot(dx, dy)

                    if dist <= radius {
                        e.HP -= damage

                        g.DamageNumbers = append(g.DamageNumbers, &DamageNumber{
                            X:     e.Pos.X,
                            Y:     e.Pos.Y,
                            Value: int(damage),
                            Life:  0.8,
                        })

                        if e.HP <= 0 {
                            e.Alive = false
                            g.Crystals = append(g.Crystals, &Crystal{
                                Pos:    e.Pos,
                                Alive:  true,
                                Sprite: g.crystalSprite,
                            })
                        }
                    }
                }
            },
        },
    }

    // Mark all abilities disabled initially
    for _, a := range g.AvailableAbilities {
        a.Enabled = false
    }
}

func (g *Game) fireDagger() {
	dir := float64(g.player.lastDir)

	p := &Projectile{
		Pos:    Vec{X: g.player.Pos.X, Y: g.player.Pos.Y},
		Vel:    Vec{X: dir * 400, Y: 0},
		Damage: 10 + float64(g.player.Level)*2,
		Alive:  true,
		Sprite: PlaceholderIcon(255, 255, 80),

		// projectile behavior
		Speed:    400, // optional, used by homing logic
		Lifetime: 60,  // optional, despawns after ~1 sec
	}

	g.Projectiles = append(g.Projectiles, p)
}
