package graphs

// ============================================================
// PROBLEM 8: Pacific Atlantic Water Flow (LeetCode #417) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an m x n island height map, water can flow from a cell
//   to its 4-directional neighbor if the neighbor's height is less
//   than or equal to the current cell's height. The Pacific ocean
//   touches the left and top edges; the Atlantic touches the right
//   and bottom edges. Return all cells from which water can flow
//   to both oceans.
//
// PARAMETERS:
//   heights [][]int — m x n matrix of non-negative heights
//
// RETURN:
//   [][]int — list of [row, col] coordinates that can reach both oceans
//
// CONSTRAINTS:
//   • m == heights.length
//   • n == heights[i].length
//   • 1 <= m, n <= 200
//   • 0 <= heights[i][j] <= 10^5
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  heights = [[1,2,2,3,5],[3,2,3,4,4],[2,4,5,3,1],[6,7,1,4,5],[5,1,1,2,4]]
//   Output: [[0,4],[1,3],[1,4],[2,2],[3,0],[3,1],[4,0]]
//   Why:    These cells can reach both Pacific (top/left) and Atlantic (bottom/right)
//
// Example 2:
//   Input:  heights = [[1]]
//   Output: [[0,0]]
//   Why:    Single cell touches both oceans
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Reverse flow: BFS/DFS from ocean borders inward (go to cells with height >= current)
// • Two visited matrices: one for Pacific-reachable, one for Atlantic-reachable
// • Result = intersection of both reachable sets
// • Target: O(m*n) time, O(m*n) space
func PacificAtlantic(heights [][]int) [][]int {
	return nil
}
