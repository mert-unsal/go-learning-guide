package hard

// ============================================================
// PROBLEM 10: Regular Expression Matching (LeetCode #10) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an input string s and a pattern p, implement regular
//   expression matching with support for '.' (matches any single
//   character) and '*' (matches zero or more of the preceding
//   element). The matching should cover the entire input string.
//
// PARAMETERS:
//   s string — input string (may be empty, contains only lowercase a-z)
//   p string — pattern (may be empty, contains lowercase a-z, '.', and '*')
//
// RETURN:
//   bool — true if the pattern matches the entire input string
//
// CONSTRAINTS:
//   • 1 <= len(s) <= 20
//   • 1 <= len(p) <= 20
//   • s contains only lowercase English letters
//   • p contains only lowercase English letters, '.', and '*'
//   • For each '*', there is a valid preceding character to match
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  s = "aa", p = "a"
//   Output: false
//   Why:    "a" does not match the entire string "aa".
//
// Example 2:
//   Input:  s = "aa", p = "a*"
//   Output: true
//   Why:    '*' means zero or more of 'a', matching "aa".
//
// Example 3:
//   Input:  s = "ab", p = ".*"
//   Output: true
//   Why:    ".*" means zero or more of any character.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • DP: dp[i][j] = whether s[:i] matches p[:j]
// • Handle '*' by either skipping the "x*" pair (zero match) or consuming one char
// • Base case: dp[0][0] = true; handle patterns like "a*b*" matching empty string
// • Target: O(m*n) time, O(m*n) space where m=len(s), n=len(p)
func IsMatch(s string, p string) bool {
	return false
}
