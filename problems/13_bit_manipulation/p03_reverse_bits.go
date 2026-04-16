package bit_manipulation

// ============================================================
// PROBLEM 3: Reverse Bits (LeetCode #190) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Reverse bits of a given 32-bit unsigned integer. For example,
//   the input binary string 00000010100101000001111010011100
//   returns 00111001011110000010100101000000.
//
// PARAMETERS:
//   n uint32 — a 32-bit unsigned integer
//
// RETURN:
//   uint32 — the integer formed by reversing all 32 bits of n
//
// CONSTRAINTS:
//   • The input must be a 32-bit unsigned integer
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  n = 43261596  (00000010100101000001111010011100)
//   Output: 964176192    (00111001011110000010100101000000)
//   Why:    All 32 bits reversed.
//
// Example 2:
//   Input:  n = 4294967293 (11111111111111111111111111111101)
//   Output: 3221225471    (10111111111111111111111111111111)
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Iterate 32 times: shift result left, OR with the lowest bit of n, shift n right
// • Divide-and-conquer: swap adjacent bits, then pairs, then nibbles, etc.
// • Target: O(1) time (fixed 32 iterations), O(1) space
func ReverseBits(n uint32) uint32 {
	return 0
}
