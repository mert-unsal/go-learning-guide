package dynamic_prog

// ============================================================
// PROBLEM 12: Decode Ways (LeetCode #91) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   A message containing letters A-Z can be encoded as numbers using
//   the mapping 'A' → "1", 'B' → "2", ..., 'Z' → "26". Given a
//   string s containing only digits, return the number of ways to
//   decode it. The answer is guaranteed to fit in a 32-bit integer.
//
// PARAMETERS:
//   s string — a string of digits ('0'-'9')
//
// RETURN:
//   int — number of ways to decode the string
//
// CONSTRAINTS:
//   • 1 ≤ len(s) ≤ 100
//   • s contains only digits and may contain leading zeros
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  s = "12"
//   Output: 2
//   Why:    "AB" (1,2) or "L" (12)
//
// Example 2:
//   Input:  s = "226"
//   Output: 3
//   Why:    "BZ" (2,26), "VF" (22,6), or "BBF" (2,2,6)
//
// Example 3:
//   Input:  s = "06"
//   Output: 0
//   Why:    "06" is not a valid encoding — leading zero has no mapping
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • dp[i] = number of ways to decode s[:i]
// • Single digit: if s[i-1] != '0', dp[i] += dp[i-1]
// • Two digits: if 10 ≤ s[i-2..i-1] ≤ 26, dp[i] += dp[i-2]
// • Only need previous two values — similar to Fibonacci with conditions
// • Target: O(n) time, O(1) space
func NumDecodings(s string) int {
	return 0
}
