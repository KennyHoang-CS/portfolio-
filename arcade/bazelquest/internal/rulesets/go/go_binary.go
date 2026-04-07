package gorules

// GoBinary simulates a rules_go go_binary target.
type GoBinary struct {
	Name string
}

func NewGoBinary() *GoBinary { return &GoBinary{} }
