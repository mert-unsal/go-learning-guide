package graphs

// ============================================================
// PROBLEM 12: Redundant Connection (LeetCode #684) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a graph that started as a tree with n nodes (1 to n)
//   and had one additional edge added, return the edge that can be
//   removed so the result is a tree. If there are multiple answers,
//   return the one that occurs last in the input.
//
// PARAMETERS:
//   edges [][]int — list of n edges [u, v] (1-indexed nodes)
//
// RETURN:
//   []int — the redundant edge [u, v]
//
// CONSTRAINTS:
//   • n == edges.length
//   • 3 <= n <= 1000
//   • edges[i].length == 2
//   • 1 <= ui, vi <= n
//   • ui != vi
//   • No duplicate edges
//   • The input graph is connected
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  edges = [[1,2],[1,3],[2,3]]
//   Output: [2,3]
//   Why:    Removing [2,3] leaves a valid tree: 1-2, 1-3
//
// Example 2:
//   Input:  edges = [[1,2],[2,3],[3,4],[1,4],[1,5]]
//   Output: [1,4]
//   Why:    Removing [1,4] breaks the cycle 1-2-3-4-1, leaving a tree
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Union-Find: process edges in order; the first edge connecting already-connected nodes is redundant
// • Use union by rank + path compression for near O(1) amortized per operation
// • The last such edge in input order is the answer (process all, return last found)
// • Target: O(n * α(n)) ≈ O(n) time, O(n) space
func FindRedundantConnection(edges [][]int) []int {
	return nil
}
