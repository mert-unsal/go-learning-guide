package hard

// ============================================================
// PROBLEM 11: Edit Distance (LeetCode #72) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   Given two strings word1 and word2, return the minimum number
//   of operations required to convert word1 into word2. You have
//   three operations: insert a character, delete a character, or
//   replace a character.
//
// PARAMETERS:
//   word1 string — the source string
//   word2 string — the target string
//
// RETURN:
//   int — minimum number of edit operations (insert, delete, replace)
//
// CONSTRAINTS:
//   • 0 <= len(word1), len(word2) <= 500
//   • word1 and word2 consist of lowercase English letters
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  word1 = "horse", word2 = "ros"
//   Output: 3
//   Why:    horse → rorse (replace h→r) → rose (remove r) → ros (remove e)
//
// Example 2:
//   Input:  word1 = "intention", word2 = "execution"
//   Output: 5
//   Why:    intention → exention → exection → execition → execution (5 ops)
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Classic DP: dp[i][j] = edit distance between word1[:i] and word2[:j]
// • If chars match: dp[i][j] = dp[i-1][j-1]
// • Else: 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1]) for delete, insert, replace
// • Target: O(m*n) time, O(m*n) space (optimizable to O(min(m,n)) with rolling array)
func MinDistance(word1 string, word2 string) int {
	return 0
}
func min3(a, b, c int) int {
	return 0
}
