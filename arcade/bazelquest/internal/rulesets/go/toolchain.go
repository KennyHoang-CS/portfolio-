package gorules

// GoToolchain represents a minimal toolchain descriptor.
type GoToolchain struct{}

func DefaultToolchain() *GoToolchain { return &GoToolchain{} }
