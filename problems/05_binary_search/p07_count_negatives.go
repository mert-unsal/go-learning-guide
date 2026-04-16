package binary_search

// ============================================================
// PROBLEM 7: Count Negative Numbers in Sorted Matrix (LeetCode #1351) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an m x n matrix grid which is sorted in non-increasing
//   order both row-wise and column-wise, return the number of
//   negative numbers in grid.
//
// PARAMETERS:
//   grid [][]int — m x n matrix sorted in non-increasing order per row and column
//
// RETURN:
//   int — the total count of negative numbers in the matrix
//
// CONSTRAINTS:
//   • m == len(grid)
//   • n == len(grid[i])
//   • 1 <= m, n <= 100
//   • -100 <= grid[i][j] <= 100
//   • Each row and column is sorted in non-increasing order
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  grid = [[4,3,2,-1],[3,2,1,-1],[1,1,-1,-2],[-1,-1,-2,-3]]
//   Output: 8
//   Why:    There are 8 negative numbers in the matrix.
//
// Example 2:
//   Input:  grid = [[3,2],[1,0]]
//   Output: 0
//   Why:    No negative numbers exist.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Start at the top-right corner. If the value is negative, all values
//   below it in that column are also negative — add (m - row) and move left.
//   If non-negative, move down.
// • Alternatively, binary search each row for the first negative.
// • Target: O(m + n) time, O(1) space

func CountNegatives(grid [][]int) int {
	return 0
}
