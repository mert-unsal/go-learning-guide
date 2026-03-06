package arrays

// ============================================================
// PROBLEM 13: Valid Sudoku (LeetCode #36) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Determine if a 9×9 Sudoku board is valid. Only the FILLED cells
//   need to be validated according to the following rules:
//
//   1. Each ROW must contain the digits 1-9 without repetition.
//   2. Each COLUMN must contain the digits 1-9 without repetition.
//   3. Each of the nine 3×3 SUB-BOXES must contain the digits 1-9
//      without repetition.
//
//   Note: A Sudoku board (partially filled) could be valid but is
//   NOT necessarily solvable. Only the filled cells need to be
//   validated. Empty cells are represented by '.'.
//
// PARAMETERS:
//   board [][]byte — a 9×9 grid. Each cell is either a digit '1'-'9'
//                     or '.' (empty).
//
// RETURN:
//   bool — true if the board is valid according to the rules above.
//
// CONSTRAINTS:
//   • board.length == 9
//   • board[i].length == 9
//   • board[i][j] is a digit '1'-'9' or '.'.
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Valid board:
//   Input: (standard valid partially-filled Sudoku)
//   Output: true
//
// Example 2 — Invalid: duplicate in row:
//   Row 0: ['5', '3', '.', '.', '7', '.', '.', '.', '5']
//   Output: false
//   Why:    '5' appears twice in row 0.
//
// Example 3 — Invalid: duplicate in column:
//   Column 0 has '8' at row 0 and '8' at row 4.
//   Output: false
//
// Example 4 — Invalid: duplicate in 3×3 box:
//   Top-left 3×3 box contains '1' twice.
//   Output: false
//
// ─── 3×3 BOX INDEX MAPPING ─────────────────────────────────
//
//   The board has 9 sub-boxes (3 rows of boxes × 3 columns of boxes).
//   For cell at (row, col), its box index is:
//
//     boxIndex = (row / 3) * 3 + (col / 3)
//
//   Box layout:
//     ┌───────┬───────┬───────┐
//     │ Box 0 │ Box 1 │ Box 2 │   rows 0-2
//     ├───────┼───────┼───────┤
//     │ Box 3 │ Box 4 │ Box 5 │   rows 3-5
//     ├───────┼───────┼───────┤
//     │ Box 6 │ Box 7 │ Box 8 │   rows 6-8
//     └───────┴───────┴───────┘
//      cols    cols    cols
//      0-2     3-5     6-8
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • You need to check 3 things per digit: row, column, box.
//   • How many "seen" sets do you need? (Hint: 9 rows + 9 cols + 9 boxes)
//   • You only iterate the board once — O(81) = O(1).
//   • What data structure tracks "have I seen this digit in this row/col/box"?
//   • Target: O(1) time (fixed 9×9), O(1) space (fixed-size sets).

// IsValidSudoku returns true if the board is a valid Sudoku configuration.
// Time: O(81) = O(1)  Space: O(81) = O(1)
func IsValidSudoku(board [][]byte) bool {
	// TODO: implement
	return false
}
