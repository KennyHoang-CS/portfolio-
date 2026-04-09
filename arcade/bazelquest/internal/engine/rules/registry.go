package rules

import "github.com/KennyHoang-CS/portfolio/bazelquest/internal/engine"

// registry maps rule names to engine.Rule implementations.
var registry = map[string]engine.Rule{}

// Register adds a rule implementation to the global registry.
func Register(name string, r engine.Rule) {
	registry[name] = r
}

// Get retrieves a rule implementation by name.
func Get(name string) engine.Rule {
	return registry[name]
}

// All returns the entire registry.
func All() map[string]engine.Rule {
	return registry
}