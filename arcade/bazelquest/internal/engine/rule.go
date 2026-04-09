package engine

type Rule interface {
	Name() string
	Validate(attrs map[string]Value) error
	Build(ctx *Context) (*Result, error)
}