package engine

import (
	"strings"
)

type ParsedRule struct {
	Kind  string
	Attrs map[string]Value
}

// VERY simple parser for early prototype.
// Later: replace with a real Starlark subset parser.
func ParseBUILD(src string) ([]ParsedRule, error) {
	lines := strings.Split(src, "\n")
	var rules []ParsedRule

	var current ParsedRule
	inRule := false

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Start of rule: go_library(
		if strings.HasSuffix(line, "(") {
			inRule = true
			kind := strings.TrimSuffix(line, "(")
			current = ParsedRule{
				Kind:  strings.TrimSpace(kind),
				Attrs: map[string]Value{},
			}
			continue
		}

		// End of rule: )
		if line == ")" {
			inRule = false
			rules = append(rules, current)
			continue
		}

		if inRule {
			// attr = "value"
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.TrimSuffix(val, ",")
			val = strings.Trim(val, "\"")

			// list support: ["a", "b"]
			if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
				val = strings.Trim(val, "[]")
				items := strings.Split(val, ",")
				var list []string
				for _, item := range items {
					item = strings.TrimSpace(item)
					item = strings.Trim(item, "\"")
					if item != "" {
						list = append(list, item)
					}
				}
				current.Attrs[key] = list
			} else {
				current.Attrs[key] = val
			}
		}
	}

	return rules, nil
}