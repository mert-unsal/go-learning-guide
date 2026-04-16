package graphs

// ============================================================
// PROBLEM 11: Surrounded Regions (LeetCode #130) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an m x n board containing 'X' and 'O', capture all
//   regions that are 4-directionally surrounded by 'X'. A region
//   is captured by flipping all 'O's into 'X's. 'O's on the
//   border (and any 'O' connected to a border 'O') are NOT captured.
//
// PARAMETERS:
//   board [][]byte — m x n board with 'X' and 'O' (modified in-place)
//
// RETURN:
//   (none) — the board is modified in-place
//
// CONSTRAINTS:
//   • m == board.length
//   • n == board[i].length
//   • 1 <= m, n <= 200
//   • board[i][j] is 'X' or 'O'
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  board = [["X","X","X","X"],
//                     ["X","O","O","X"],
//                     ["X","X","O","X"],
//                     ["X","O","X","X"]]
//   Output: [["X","X","X","X"],
//             ["X","X","X","X"],
//             ["X","X","X","X"],
//             ["X","O","X","X"]]
//   Why:    Interior O's at (1,1),(1,2),(2,2) are surrounded → captured.
//           Border-connected O at (3,1) is NOT captured.
//
// Example 2:
//   Input:  board = [["X"]]
//   Output: [["X"]]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Reverse thinking: instead of finding surrounded O's, find UN-surrounded O's
// • DFS/BFS from all border O's, mark them as safe (e.g., temporarily change to 'S')
// • Then scan entire board: remaining O → X (captured), S → O (restored)
// • Target: O(m*n) time, O(m*n) space
func SurroundedRegions(board [][]byte) {
}
