package binary_search

// ============================================================
// PROBLEM 6: Guess Number Higher or Lower (LeetCode #374) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   I pick a number from 1 to n. You have to guess which number I
//   picked. Every time you guess wrong, I tell you whether the number
//   is higher or lower than your guess. You call a pre-defined API
//   guess(num) which returns -1 (your guess is higher), 1 (your guess
//   is lower), or 0 (correct).
//
// PARAMETERS:
//   n       int           — the upper bound of the range [1, n]
//   guessFn func(int) int — the guess API: returns -1, 0, or 1
//
// RETURN:
//   int — the number that was picked
//
// CONSTRAINTS:
//   • 1 <= n <= 2^31 - 1
//   • 1 <= pick <= n
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  n = 10, pick = 6
//   Output: 6
//   Why:    Binary search converges: guess(5)→1, guess(8)→-1, guess(6)→0.
//
// Example 2:
//   Input:  n = 1, pick = 1
//   Output: 1
//   Why:    Only one possible number.
//
// Example 3:
//   Input:  n = 2, pick = 1
//   Output: 1
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Standard binary search between 1 and n, using the guess API to
//   decide which half to narrow to.
// • Use mid = left + (right - left) / 2 to avoid overflow.
// • Target: O(log n) time, O(1) space

func GuessNumber(n int, guessFn func(int) int) int {
	return 0
}
