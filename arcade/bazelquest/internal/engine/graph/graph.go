package graph

// Node represents a dependency node in the build graph.
type Node struct {
	Name string
}

// Graph is a placeholder for dependency graph logic.
type Graph struct{}

// New creates an empty graph.
func New() *Graph { return &Graph{} }
