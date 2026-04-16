package graphs

// ============================================================
// PROBLEM 10: Max Area of Island (LeetCode #695) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an m x n binary grid where 1 = land and 0 = water,
//   return the area of the largest island. An island's area is
//   the number of connected 1-cells (4-directional adjacency).
//   If there is no island, return 0.
//
// PARAMETERS:
//   grid [][]int — m x n grid with values 0 or 1
//
// RETURN:
//   int — area of the largest island, or 0 if no islands
//
// CONSTRAINTS:
//   • m == grid.length
//   • n == grid[i].length
//   • 1 <= m, n <= 50
//   • grid[i][j] is 0 or 1
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  grid = [[0,0,1,0,0,0,0,1,0,0,0,0,0],
//                    [0,0,0,0,0,0,0,1,1,1,0,0,0],
//                    [0,1,1,0,1,0,0,0,0,0,0,0,0],
//                    [0,1,0,0,1,1,0,0,1,0,1,0,0],
//                    [0,1,0,0,1,1,0,0,1,1,1,0,0],
//                    [0,0,0,0,0,0,0,0,0,0,1,0,0],
//                    [0,0,0,0,0,0,0,1,1,1,0,0,0],
//                    [0,0,0,0,0,0,0,1,1,0,0,0,0]]
//   Output: 6
//   Why:    The largest island (bottom-right cluster) has 6 connected cells
//
// Example 2:
//   Input:  grid = [[0,0,0,0,0,0,0,0]]
//   Output: 0
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • DFS/BFS flood-fill: for each unvisited 1, compute area and track the maximum
// • Mark visited by setting cell to 0 (in-place) or use a visited matrix
// • Target: O(m*n) time, O(m*n) space
func MaxAreaOfIsland(grid [][]int) int {
	return 0
}
