package graphs

// ============================================================
// PROBLEM 4: Clone Graph (LeetCode #133) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a reference to a node in a connected undirected graph,
//   return a deep copy (clone) of the graph. Each node contains a
//   value and a list of its neighbors.
//
// PARAMETERS:
//   node *GraphNode — reference to any node in the graph (may be nil)
//
// RETURN:
//   *GraphNode — reference to the corresponding node in the cloned graph
//
// CONSTRAINTS:
//   • 0 <= number of nodes <= 100
//   • 1 <= Node.Val <= 100
//   • All Node.Val are unique
//   • No repeated edges or self-loops
//   • The graph is connected and undirected
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  adjList = [[2,4],[1,3],[2,4],[1,3]]
//   Output: [[2,4],[1,3],[2,4],[1,3]]
//   Why:    4 nodes: 1→{2,4}, 2→{1,3}, 3→{2,4}, 4→{1,3}. Deep copy preserves structure.
//
// Example 2:
//   Input:  adjList = [[]]
//   Output: [[]]
//   Why:    Single node with no neighbors
//
// Example 3:
//   Input:  adjList = []
//   Output: []
//   Why:    Empty graph (nil node)
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Use a hashmap (old node → new node) to track already-cloned nodes
// • DFS or BFS: for each node, clone it, then recursively clone its neighbors
// • The map also serves as a "visited" set to avoid infinite loops in cycles
// • Target: O(V+E) time, O(V) space
func CloneGraph(node *GraphNode) *GraphNode {
	return nil
}
