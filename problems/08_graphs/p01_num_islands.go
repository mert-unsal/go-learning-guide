package graphs

// ============================================================
// PROBLEM 1: Number of Islands (LeetCode #200) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an m x n 2D grid of '1's (land) and '0's (water),
//   return the number of islands. An island is surrounded by water
//   and is formed by connecting adjacent land cells horizontally
//   or vertically.
//
// PARAMETERS:
//   grid [][]byte — m x n grid where each cell is '1' or '0'
//
// RETURN:
//   int — the number of distinct islands
//
// CONSTRAINTS:
//   • m == grid.length
//   • n == grid[i].length
//   • 1 <= m, n <= 300
//   • grid[i][j] is '0' or '1'
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  grid = [["1","1","1","1","0"],
//                    ["1","1","0","1","0"],
//                    ["1","1","0","0","0"],
//                    ["0","0","0","0","0"]]
//   Output: 1
//   Why:    All land cells are connected into a single island
//
// Example 2:
//   Input:  grid = [["1","1","0","0","0"],
//                    ["1","1","0","0","0"],
//                    ["0","0","1","0","0"],
//                    ["0","0","0","1","1"]]
//   Output: 3
//   Why:    Three separate groups of connected land cells
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • DFS/BFS flood-fill: for each unvisited '1', increment count and mark all connected '1's
// • Mark visited cells by changing '1' to '0' (in-place) to avoid extra visited array
// • Union-Find is an alternative approach
// • Target: O(m*n) time, O(m*n) space (recursion stack in worst case)

func NumIslands(grid [][]byte) int {
	return 0
}
func dfsIsland(grid [][]byte, r, c, rows, cols int) {
}
