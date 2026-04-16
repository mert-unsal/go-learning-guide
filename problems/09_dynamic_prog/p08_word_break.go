package dynamic_prog

// ============================================================
// PROBLEM 8: Word Break (LeetCode #139) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a string s and a dictionary of strings wordDict, return true
//   if s can be segmented into a space-separated sequence of one or
//   more dictionary words. The same dictionary word may be reused
//   multiple times in the segmentation.
//
// PARAMETERS:
//   s        string   — the input string to segment
//   wordDict []string — list of valid dictionary words
//
// RETURN:
//   bool — true if s can be segmented into dictionary words
//
// CONSTRAINTS:
//   • 1 ≤ len(s) ≤ 300
//   • 1 ≤ len(wordDict) ≤ 1000
//   • 1 ≤ len(wordDict[i]) ≤ 20
//   • s and wordDict[i] consist of only lowercase English letters
//   • All strings in wordDict are unique
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  s = "leetcode", wordDict = ["leet", "code"]
//   Output: true
//   Why:    "leet" + "code" = "leetcode"
//
// Example 2:
//   Input:  s = "applepenapple", wordDict = ["apple", "pen"]
//   Output: true
//   Why:    "apple" + "pen" + "apple"
//
// Example 3:
//   Input:  s = "catsandog", wordDict = ["cats", "dog", "sand", "and", "cat"]
//   Output: false
//   Why:    No valid segmentation exists
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • DP: dp[i] = true if s[:i] can be segmented
// • For each position i, check all words: if dp[i-len(w)] && s[i-len(w):i] == w
// • Put wordDict into a map/set for O(1) lookups
// • Target: O(n² × m) time, O(n) space where m = max word length
func WordBreak(s string, wordDict []string) bool {
	return false
}
