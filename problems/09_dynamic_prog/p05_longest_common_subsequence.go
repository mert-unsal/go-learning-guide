package dynamic_prog

// ============================================================
// PROBLEM 5: Longest Common Subsequence (LeetCode #1143) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given two strings, return the length of their longest common
//   subsequence. A subsequence is a sequence that can be derived
//   from another sequence by deleting some or no elements without
//   changing the order of the remaining elements.
//
// PARAMETERS:
//   text1 string — first string
//   text2 string — second string
//
// RETURN:
//   int — length of the longest common subsequence (0 if none)
//
// CONSTRAINTS:
//   • 1 ≤ len(text1), len(text2) ≤ 1000
//   • text1 and text2 consist of only lowercase English letters
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  text1 = "abcde", text2 = "ace"
//   Output: 3
//   Why:    LCS is "ace"
//
// Example 2:
//   Input:  text1 = "abc", text2 = "abc"
//   Output: 3
//   Why:    Entire string is the common subsequence
//
// Example 3:
//   Input:  text1 = "abc", text2 = "def"
//   Output: 0
//   Why:    No common subsequence
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Classic 2D DP: dp[i][j] = LCS length of text1[:i] and text2[:j]
// • If characters match: dp[i][j] = dp[i-1][j-1] + 1
// • If mismatch: dp[i][j] = max(dp[i-1][j], dp[i][j-1])
// • Space can be optimized to O(min(m,n)) with rolling row
// • Target: O(m×n) time, O(m×n) space

func LongestCommonSubsequence(text1 string, text2 string) int {
	return 0
}
