package main

import (
	"log"

	fonts "github.com/KennyHoang-CS/portfolio/bazelquest/internal/assets"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/build"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/editor"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/input"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/levels"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/terminal"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/ui"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/ui/buttons"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/workspace"

	_ "github.com/KennyHoang-CS/portfolio/bazelquest/internal/levels/go"
	gorules "github.com/KennyHoang-CS/portfolio/bazelquest/internal/rulesets/go"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
    editor    *editor.Editor
    terminal  *terminal.Terminal
    toolbar   *buttons.Toolbar
    ToolbarUI *ui.ToolbarUI
    workspace *workspace.Workspace
    level     levels.Level

    filePanel *ui.FileSystemPanel   // ⭐ ADD THIS
}

func NewGame() *Game {
    // 1. Create workspace
    ws := workspace.New()

    // 2. Load the first level
    lvl := levels.First()

    // 3. Populate workspace with level files
    for path, contents := range lvl.InitialFiles {
        ws.Write(path, contents)
    }

    // 4. Create editor bound to workspace
    e := editor.New(ws)

    // 5. Open the BUILD file by default
    if ws.Exists("BUILD") {
        e.Open("BUILD")
    }

    // 6. Create terminal
    t := terminal.New()
    t.LogInfo("Level: " + lvl.Title)
    t.LogInfo(lvl.Description)

    // 7. Create game struct
    g := &Game{
        editor:    e,
        terminal:  t,
        workspace: ws,   // NEW FIELD
        level:     lvl,  // NEW FIELD
    }

    // 8. Initialize UI
    g.initUI()

    return g
}

func (g *Game) Update() error {
    // Update editor + terminal input
    input.Update(g.editor, g.terminal)

    // Update EbitenUI toolbar
    g.ToolbarUI.UI.Update()

    if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
        g.BuildOnly()
    }

    if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()
        g.editor.Click(x, y)
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    w, h := screen.Size()

    // 1. Draw your custom panels
    ui.DrawPanels(screen, w, h, fonts.EditorFont)

    // ⭐ Draw the File System panel
    g.filePanel.Draw(screen)

    // 2. Position editor + terminal inside their panels
    editorRect := ui.EditorBounds(w, h)
    g.editor.SetOffset(editorRect.Min.X+16, editorRect.Min.Y+48)

    termRect := ui.TerminalBounds(w, h)
    g.terminal.SetOffset(termRect.Min.X+16, termRect.Min.Y+48)

    // 3. Draw editor + terminal
    g.editor.Draw(screen)
    g.terminal.Draw(screen)

    // 4. Draw EbitenUI toolbar (layout handles positioning)
    g.ToolbarUI.UI.Draw(screen)
}

func (g *Game) Layout(outW, outH int) (int, int) { return 1280, 720 }

func main() {
	gorules.Register()

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("BazelQuest - Prototype")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) initUI() {
    g.ToolbarUI = ui.NewToolbarUI(
        func() { g.RunProgram() },
        func() { g.BuildOnly() },
        func() { g.FormatCode() },
    )

    // ⭐ Add File System Panel
    g.filePanel = ui.NewFileSystemPanel(
        g.workspace,     // workspace reference
        16, 80,          // x, y position
        200, 600,        // width, height
        func(path string) {
            g.editor.Open(path) // open file in editor when clicked
        },
    )
}

func (g *Game) RunBuild() {
	build.Run(g.editor.Buffer(), g.terminal)
}

func (g *Game) FormatCode() {
	g.terminal.LogInfo("Format not implemented yet")
}

func (g *Game) BuildOnly() {
	build.Run(g.editor.Buffer(), g.terminal)
}

func (g *Game) RunProgram() {
	ok := build.Run(g.editor.Buffer(), g.terminal)
	if !ok {
		return
	}

	g.terminal.LogInfo("Running program...")
}
