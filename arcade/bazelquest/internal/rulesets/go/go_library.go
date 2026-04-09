package gorules

import (
	"fmt"

	"github.com/KennyHoang-CS/portfolio/bazelquest/internal/engine"
)

type GoLibraryRule struct{}

func NewGoLibraryRule() *GoLibraryRule { return &GoLibraryRule{} }

func (r *GoLibraryRule) Name() string { return "go_library" }

func (r *GoLibraryRule) Validate(attrs map[string]engine.Value) error {
	srcs := engine.AsStringList(attrs["srcs"])
	if len(srcs) == 0 {
		return ErrMissingSrcs
	}
	return nil
}

func (r *GoLibraryRule) Build(ctx *engine.Context) (*engine.Result, error) {
	res := engine.NewResult()

	if err := DefaultToolchain.Check(); err != nil {
		res.Fail(err.Error())
		return res, nil
	}

	srcs := engine.AsStringList(ctx.Attrs["srcs"])
	name := engine.AsString(ctx.Attrs["name"])
	if name == "" {
		name = "unnamed"
	}

	res.Log(fmt.Sprintf("Compiling go_library %q with %d srcs", name, len(srcs)))

	objects := make([]string, 0, len(srcs))
	for _, s := range srcs {
		objects = append(objects, s+".o")
	}

	_ = &GoCompileProvider{
		PackageName: name,
		Objects:     objects,
	}

	// For now we just log; provider graph comes later.
	return res, nil
}