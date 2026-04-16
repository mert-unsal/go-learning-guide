package bit_manipulation

// ============================================================
// PROBLEM 5: Sum of Two Integers (LeetCode #371) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given two integers a and b, return the sum of the two
//   integers without using the operators + and -.
//
// PARAMETERS:
//   a int — first integer
//   b int — second integer
//
// RETURN:
//   int — the sum a + b computed without arithmetic operators
//
// CONSTRAINTS:
//   • -1000 <= a, b <= 1000
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  a = 1, b = 2
//   Output: 3
//
// Example 2:
//   Input:  a = 2, b = 3
//   Output: 5
//
// Example 3:
//   Input:  a = -1, b = 1
//   Output: 0
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • XOR gives sum without carry: a ^ b
// • AND + left shift gives carry: (a & b) << 1
// • Repeat until carry is 0: a = a ^ b, b = (a & b) << 1
// • Be careful with negative numbers in languages with fixed-width integers
// • Target: O(1) time (bounded by bit width), O(1) space
func GetSum(a int, b int) int {
	return 0
}
