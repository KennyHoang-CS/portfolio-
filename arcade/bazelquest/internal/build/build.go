package build

import (
	"strings"

	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/terminal"
)

func Run(buffer string, term *terminal.Terminal) {
	term.LogInfo("bazel build //app:server")
	if strings.Contains(buffer, "go_library(") {
		term.LogSuccess("BUILD SUCCESS")
	} else {
		term.LogError("ERROR: expected go_library() in BUILD file")
	}
}