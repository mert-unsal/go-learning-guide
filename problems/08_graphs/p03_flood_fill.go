package graphs

// ============================================================
// PROBLEM 3: Flood Fill (LeetCode #733) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an m x n image represented by a 2D array of integers,
//   a starting pixel (sr, sc), and a new color, perform a flood
//   fill. Starting from the pixel at (sr, sc), change the color of
//   every pixel connected 4-directionally that shares the same
//   original color to the new color.
//
// PARAMETERS:
//   image [][]int — m x n grid of pixel colors
//   sr    int     — starting row index
//   sc    int     — starting column index
//   color int     — the new color to fill with
//
// RETURN:
//   [][]int — the modified image after flood fill
//
// CONSTRAINTS:
//   • m == image.length
//   • n == image[i].length
//   • 1 <= m, n <= 50
//   • 0 <= image[i][j], color < 2^16
//   • 0 <= sr < m
//   • 0 <= sc < n
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  image = [[1,1,1],[1,1,0],[1,0,1]], sr = 1, sc = 1, color = 2
//   Output: [[2,2,2],[2,2,0],[2,0,1]]
//   Why:    Starting pixel (1,1) has color 1; all connected 1's become 2
//
// Example 2:
//   Input:  image = [[0,0,0],[0,0,0]], sr = 0, sc = 0, color = 0
//   Output: [[0,0,0],[0,0,0]]
//   Why:    New color equals original color — no change needed
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • DFS/BFS from (sr, sc), changing all matching-color neighbors
// • Early return if new color == original color (avoids infinite loop)
// • No visited array needed — changed color acts as the visited marker
// • Target: O(m*n) time, O(m*n) space
func FloodFill(image [][]int, sr int, sc int, color int) [][]int {
	return nil
}
func fill(image [][]int, r, c, originalColor, newColor, rows, cols int) {
}
