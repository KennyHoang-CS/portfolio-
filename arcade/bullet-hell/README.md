🎮 Bullet Hell (Go + Ebiten)
A fast‑paced, Touhou‑inspired bullet‑hell shooter built in Go using the Ebiten game engine.
This project focuses on clean architecture, an ECS‑style entity system, modular bullet patterns, and polished visual effects.

✨ Features
🌀 Bullet Patterns
- Modular PatternManager interface
- Curved bullets, spreads, spirals, waves
- Difficulty scaling
- Graze detection + scoring
🎨 Visual Effects
- Parallax starfield background
- Neon bullet sprites
- Player aura + hitbox rendering
- Graze sparks
- Smooth vector rendering
🎮 Player System
- WASD movement
- Precise hitbox
- Aura glow
- Collision + graze scoring
🔊 Audio
- MP3 BGM playback
- Automatic looping
- Volume control
🧩 ECS‑Style Architecture
- Entities composed of components:
- Position
- Velocity
- Hitbox
- BulletTag
- PlayerTag
- SparkTag
- World manages entity lifecycle + cleanup


🚀 Running the Game
1. Install Go
https://go.dev/dl/
2. Install Ebiten
go get github.com/hajimehoshi/ebiten/v2


3. Run the game
From the project root:
go run ./cmd/game



🔧 Controls
Player movement: WASD 



🧠 Technical Highlights
✔ Custom ECS‑style architecture
Lightweight, fast, and easy to extend.
✔ PatternManager
Allows you to plug in new bullet patterns without touching core game logic.
✔ Graze System
Rewards near‑misses with score and visual sparks.
✔ Vector Rendering
Player, aura, and sparks use Ebiten’s vector API for crisp shapes.
✔ Audio System
MP3 decoding + looping BGM via Ebiten’s audio context.

🛠️ Future Improvements
- Focus mode (slow movement + stronger aura)
- Player shooting
- Boss phases + health bars
- Replay system
- WASM export
- Gamepad support
- High‑visibility accessibility mode


🙌 Author
Kenny Hoang
https://www.linkedin.com/in/kennyhoang-cs/
