package gorules

import "github.com/KennyHoang-CS/portfolio/bazelquest/internal/engine/rules"

func Register() {
	rules.Register("go_library", NewGoLibraryRule())
	rules.Register("go_binary", NewGoBinaryRule())
}