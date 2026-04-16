package bit_manipulation

// ============================================================
// PROBLEM 7: Power of Two (LeetCode #231) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer n, return true if it is a power of two.
//   An integer n is a power of two if there exists an integer x
//   such that n == 2^x.
//
// PARAMETERS:
//   n int — integer to check
//
// RETURN:
//   bool — true if n is a power of two
//
// CONSTRAINTS:
//   • -2^31 <= n <= 2^31 - 1
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  n = 1
//   Output: true
//   Why:    2^0 = 1.
//
// Example 2:
//   Input:  n = 16
//   Output: true
//   Why:    2^4 = 16.
//
// Example 3:
//   Input:  n = 3
//   Output: false
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • A power of two has exactly one set bit: n & (n-1) == 0 (for n > 0)
// • Edge case: n must be positive (0 and negatives are not powers of two)
// • Target: O(1) time, O(1) space
func IsPowerOfTwo(n int) bool {
	return false
}
