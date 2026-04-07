package input

import (
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/editor"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/terminal"
)

func Update(e *editor.Editor, t *terminal.Terminal) {
	// For now, all keyboard input goes to the editor.
	// Later: focus switching, hotkeys, etc.
	e.Update()
	t.Update()
}