package graphs

// Shared types for graph problems.

// GraphNode is a graph node for the clone problem.
type GraphNode struct {
	Val       int
	Neighbors []*GraphNode
}
