package strings_problems

// ============================================================
// PROBLEM 13: String to Integer (atoi) (LeetCode #8) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Implement the myAtoi(string s) function, which converts a string
//   to a 32-bit signed integer.
//
//   Algorithm:
//   1. Skip leading whitespace.
//   2. Determine the sign (+ or -). Default is positive.
//   3. Read digits until a non-digit or end of string.
//   4. Clamp the result to the 32-bit signed integer range [−2³¹, 2³¹−1].
//
// CONSTRAINTS:
//   • 0 <= s.length <= 200
//   • s consists of English letters, digits, ' ', '+', '-', '.'.
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: "42"              → 42
// Example 2: "   -42"          → -42
// Example 3: "4193 with words" → 4193
// Example 4: "words and 987"   → 0    (first non-space is not digit/sign)
// Example 5: "-91283472332"    → -2147483648  (clamped to INT_MIN)
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Handle overflow BEFORE it happens — check before multiplying by 10.
//   • Edge cases: empty string, only whitespace, only sign, leading zeros.
//   • Target: O(n) time, O(1) space.

// MyAtoi converts string to 32-bit integer following LeetCode rules.
// Time: O(n)  Space: O(1)
func MyAtoi(s string) int {
	// TODO: implement
	return 0
}
