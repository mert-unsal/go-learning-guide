package dynamic_prog

// ============================================================
// PROBLEM 15: Interleaving String (LeetCode #97) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given strings s1, s2, and s3, determine whether s3 is formed by
//   an interleaving of s1 and s2. An interleaving of two strings is a
//   configuration where s1 and s2 are divided into substrings such
//   that s3 = s1_1 + s2_1 + s1_2 + s2_2 + ... (preserving order
//   within each string, but interleaving between them).
//
// PARAMETERS:
//   s1 string — first source string
//   s2 string — second source string
//   s3 string — target string to check
//
// RETURN:
//   bool — true if s3 is formed by interleaving s1 and s2
//
// CONSTRAINTS:
//   • 0 ≤ len(s1), len(s2) ≤ 100
//   • 0 ≤ len(s3) ≤ 200
//   • s1, s2, and s3 consist of lowercase English letters
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  s1 = "aabcc", s2 = "dbbca", s3 = "aadbbcbcac"
//   Output: true
//   Why:    s3 interleaves s1="aa|bc|c" with s2="dbbc|a"
//
// Example 2:
//   Input:  s1 = "aabcc", s2 = "dbbca", s3 = "aadbbbaccc"
//   Output: false
//
// Example 3:
//   Input:  s1 = "", s2 = "", s3 = ""
//   Output: true
//   Why:    All empty strings trivially interleave
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • First check: len(s1) + len(s2) must equal len(s3)
// • 2D DP: dp[i][j] = true if s3[:i+j] is an interleaving of s1[:i] and s2[:j]
// • Transition: dp[i][j] = (dp[i-1][j] && s1[i-1]==s3[i+j-1]) || (dp[i][j-1] && s2[j-1]==s3[i+j-1])
// • Space optimization: use 1D array of length len(s2)+1
// • Target: O(m×n) time, O(n) space
func IsInterleave(s1, s2, s3 string) bool {
	return false
}
