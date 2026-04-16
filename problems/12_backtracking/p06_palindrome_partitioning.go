package backtracking

// ============================================================
// PROBLEM 7: Palindrome Partitioning (LeetCode #131) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a string s, partition s such that every substring of
//   the partition is a palindrome. Return all possible palindrome
//   partitionings of s.
//
// PARAMETERS:
//   s string — input string of lowercase English letters
//
// RETURN:
//   [][]string — all possible partitions where each part is a palindrome
//
// CONSTRAINTS:
//   • 1 <= len(s) <= 16
//   • s consists of only lowercase English letters
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  s = "aab"
//   Output: [["a","a","b"],["aa","b"]]
//   Why:    Both partitions consist entirely of palindromes.
//
// Example 2:
//   Input:  s = "a"
//   Output: [["a"]]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Backtracking: try every prefix that is a palindrome, then recurse on suffix
// • Precompute a palindrome table dp[i][j] for O(1) palindrome checks
// • Target: O(n * 2^n) time, O(n) space for recursion depth
func Partition(s string) [][]string {
	return nil
}
