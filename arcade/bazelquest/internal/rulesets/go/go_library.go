package gorules

// GoLibrary simulates a rules_go go_library target.
type GoLibrary struct {
	Name string
}

func NewGoLibrary() *GoLibrary { return &GoLibrary{} }
