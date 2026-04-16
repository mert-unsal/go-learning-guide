package dynamic_prog

// ============================================================
// PROBLEM 4: Unique Paths (LeetCode #62) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   A robot is located at the top-left corner of an m × n grid.
//   It can only move either down or right at any point in time.
//   The robot is trying to reach the bottom-right corner.
//   Return the number of unique paths the robot can take.
//
// PARAMETERS:
//   m int — number of rows in the grid
//   n int — number of columns in the grid
//
// RETURN:
//   int — number of unique paths from top-left to bottom-right
//
// CONSTRAINTS:
//   • 1 ≤ m, n ≤ 100
//   • The answer is guaranteed to fit in a 32-bit integer
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  m = 3, n = 7
//   Output: 28
//
// Example 2:
//   Input:  m = 3, n = 2
//   Output: 3
//   Why:    Right→Down→Down, Down→Down→Right, Down→Right→Down
//
// Example 3:
//   Input:  m = 1, n = 1
//   Output: 1
//   Why:    Already at the destination
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • DP recurrence: dp[r][c] = dp[r-1][c] + dp[r][c-1]
// • First row and first column are all 1s (only one way to reach them)
// • Space optimization: use a single 1D array of length n
// • Also solvable via combinatorics: C(m+n-2, m-1)
// • Target: O(m×n) time, O(n) space

func UniquePaths(m int, n int) int {
	return 0
}
