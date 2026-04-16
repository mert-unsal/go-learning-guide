package backtracking

// ============================================================
// PROBLEM 5: Word Search (LeetCode #79) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an m×n grid of characters board and a string word,
//   return true if word exists in the grid. The word can be
//   constructed from letters of sequentially adjacent cells
//   (horizontally or vertically neighboring). The same cell
//   may not be used more than once.
//
// PARAMETERS:
//   board [][]byte — m×n grid of characters
//   word  string   — the target word to search for
//
// RETURN:
//   bool — true if the word can be found in the grid
//
// CONSTRAINTS:
//   • m == len(board), n == len(board[i])
//   • 1 <= m, n <= 6
//   • 1 <= len(word) <= 15
//   • board and word consist of only lowercase and uppercase English letters
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word = "ABCCED"
//   Output: true
//   Why:    Path: A(0,0)→B(0,1)→C(0,2)→C(1,2)→E(2,2)→D(2,1).
//
// Example 2:
//   Input:  board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word = "SEE"
//   Output: true
//
// Example 3:
//   Input:  board = [["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word = "ABCB"
//   Output: false
//   Why:    Cannot reuse the cell at (0,1) for the final 'B'.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • DFS/backtracking from every cell that matches word[0]
// • Mark cells as visited (e.g., flip to '#'), restore on backtrack
// • Prune early if remaining word length > remaining unvisited cells
// • Target: O(m * n * 3^L) time where L=len(word), O(L) space for recursion
func Exist(board [][]byte, word string) bool {
	return false
}
