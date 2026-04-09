package rules

import (
	"errors"
	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/engine"
)

type Rule interface {
    Name() string
    Validate(attrs map[string]engine.Value) error
    Build(ctx *engine.Context) (*engine.Result, error)
}

// Helper for common validation errors
var ErrInvalidAttributes = errors.New("invalid attributes")