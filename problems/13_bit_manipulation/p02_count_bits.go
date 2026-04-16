package bit_manipulation

// ============================================================
// PROBLEM 2: Counting Bits (LeetCode #338) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer n, return an array ans of length n+1 such
//   that for each i (0 <= i <= n), ans[i] is the number of 1-bits
//   in the binary representation of i.
//
// PARAMETERS:
//   n int — non-negative integer upper bound
//
// RETURN:
//   []int — array where result[i] = popcount(i) for i in [0..n]
//
// CONSTRAINTS:
//   • 0 <= n <= 10^5
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  n = 2
//   Output: [0,1,1]
//   Why:    0→0 bits, 1→1 bit, 2→1 bit.
//
// Example 2:
//   Input:  n = 5
//   Output: [0,1,1,2,1,2]
//   Why:    0→0, 1→1, 2→1, 3→2, 4→1, 5→2.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • DP relation: ans[i] = ans[i >> 1] + (i & 1)
// • Alternatively: ans[i] = ans[i & (i-1)] + 1 (Brian Kernighan relation)
// • Target: O(n) time, O(n) space — do it without popcount per number
func CountBits(n int) []int {
	return nil
}
