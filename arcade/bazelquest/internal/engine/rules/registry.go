package rules

// RuleFactory constructs rule instances.
type RuleFactory func() interface{}

// Registry maps rule names to their factory functions.
var Registry = map[string]RuleFactory{}

// Register registers a rule factory under the given name.
func Register(name string, f RuleFactory) {
	Registry[name] = f
}
