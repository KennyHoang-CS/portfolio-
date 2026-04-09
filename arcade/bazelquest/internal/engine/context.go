package engine

// Forward declaration — DO NOT import rules here.
type Provider interface {
    Kind() string
}

type Context struct {
    Attrs map[string]Value
    Deps  map[string][]Provider
}

func NewContext(attrs map[string]Value, deps map[string][]Provider) *Context {
    return &Context{
        Attrs: attrs,
        Deps:  deps,
    }
}