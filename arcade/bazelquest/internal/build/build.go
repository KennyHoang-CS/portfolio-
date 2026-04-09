package build

import (
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/engine"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/engine/rules"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/terminal"
)

func Run(buffer string, term *terminal.Terminal) bool {
    term.LogInfo("bazel build //:target")

    parsed, err := engine.ParseBUILD(buffer)
    if err != nil {
        term.LogError("Parse error: " + err.Error())
        return false
    }

    result := engine.Evaluate(parsed, rules.Get)

    for _, log := range result.Logs {
        if result.Success {
            term.LogSuccess(log)
        } else {
            term.LogError(log)
        }
    }

    if result.Success {
        term.LogSuccess("BUILD SUCCESS")
        return true
    } else {
        term.LogError("BUILD FAILED")
        return false
    }
}