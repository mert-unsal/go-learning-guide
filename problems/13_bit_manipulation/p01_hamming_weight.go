package bit_manipulation

// ============================================================
// PROBLEM 1: Number of 1 Bits / Hamming Weight (LeetCode #191) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a positive integer n, return the number of set bits
//   (1-bits) in its binary representation (also known as the
//   Hamming weight).
//
// PARAMETERS:
//   n uint32 — a 32-bit unsigned integer
//
// RETURN:
//   int — count of 1-bits in the binary representation of n
//
// CONSTRAINTS:
//   • 1 <= n <= 2^31 - 1
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  n = 11 (binary: 1011)
//   Output: 3
//   Why:    Three bits are set: positions 0, 1, and 3.
//
// Example 2:
//   Input:  n = 128 (binary: 10000000)
//   Output: 1
//
// Example 3:
//   Input:  n = 2147483645 (binary: 1111111111111111111111111111101)
//   Output: 30
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Brian Kernighan's trick: n & (n-1) clears the lowest set bit
// • Count how many times you can clear a bit until n == 0
// • Alternatively, shift and mask: check each bit with n & 1
// • Target: O(k) time where k = number of set bits, O(1) space
func HammingWeight(n uint32) int {
	return 0
}
