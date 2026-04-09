package engine

import "fmt"

func Evaluate(parsed []ParsedRule, lookup func(kind string) Rule) *Result {
	result := NewResult()

	for _, pr := range parsed {
		rule := lookup(pr.Kind)
		if rule == nil {
			result.Fail(fmt.Sprintf("Unknown rule: %s", pr.Kind))
			continue
		}

		// Validate attributes
		if err := rule.Validate(pr.Attrs); err != nil {
			result.Fail(fmt.Sprintf("%s: %v", pr.Kind, err))
			continue
		}

		// TODO: resolve deps later
		ctx := NewContext(pr.Attrs, nil)

		// Build
		r, err := rule.Build(ctx)
		if err != nil {
			result.Fail(fmt.Sprintf("%s: %v", pr.Kind, err))
			continue
		}

		result.Logs = append(result.Logs, r.Logs...)
		if !r.Success {
			result.Success = false
		}
	}

	return result
}