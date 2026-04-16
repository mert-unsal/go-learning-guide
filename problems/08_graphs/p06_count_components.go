package graphs

// ============================================================
// PROBLEM 6: Number of Connected Components (LeetCode #323) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given n nodes labeled from 0 to n-1 and a list of undirected
//   edges, return the number of connected components in the graph.
//
// PARAMETERS:
//   n     int     — number of nodes
//   edges [][]int — list of undirected edges [u, v]
//
// RETURN:
//   int — number of connected components
//
// CONSTRAINTS:
//   • 1 <= n <= 2000
//   • 0 <= edges.length <= 5000
//   • edges[i].length == 2
//   • 0 <= ui, vi < n
//   • ui != vi
//   • No duplicate edges
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  n = 5, edges = [[0,1],[1,2],[3,4]]
//   Output: 2
//   Why:    Component {0,1,2} and component {3,4}
//
// Example 2:
//   Input:  n = 5, edges = [[0,1],[1,2],[2,3],[3,4]]
//   Output: 1
//   Why:    All nodes are connected in a single component
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • DFS/BFS: iterate all nodes; for each unvisited node, run traversal and increment count
// • Union-Find: start with n components, each union reduces count by 1
// • Target: O(V+E) time, O(V+E) space
func CountComponents(n int, edges [][]int) int {
	return 0
}
