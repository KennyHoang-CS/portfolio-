package gorules

import (
	"fmt"

	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/engine"
)

type GoBinaryRule struct{}

func NewGoBinaryRule() *GoBinaryRule { return &GoBinaryRule{} }

func (r *GoBinaryRule) Name() string { return "go_binary" }

func (r *GoBinaryRule) Validate(attrs map[string]engine.Value) error {
	srcs := engine.AsStringList(attrs["srcs"])
	if len(srcs) == 0 {
		return ErrMissingBinarySrcs
	}
	return nil
}

func (r *GoBinaryRule) Build(ctx *engine.Context) (*engine.Result, error) {
	res := engine.NewResult()

	if err := DefaultToolchain.Check(); err != nil {
		res.Fail(err.Error())
		return res, nil
	}

	name := engine.AsString(ctx.Attrs["name"])
	if name == "" {
		name = "app"
	}

	res.Log(fmt.Sprintf("Linking go_binary %q", name))

	_ = &GoBinaryProvider{
		Name:   name,
		Binary: name,
	}

	return res, nil
}