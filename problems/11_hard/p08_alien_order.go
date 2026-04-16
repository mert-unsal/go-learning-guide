package hard

// ============================================================
// PROBLEM 9: Alien Dictionary (LeetCode #269) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   There is a new alien language that uses the English alphabet,
//   but the order of the letters is unknown. You are given a list
//   of strings words from the alien dictionary, sorted
//   lexicographically by the rules of this new language. Derive
//   the order of letters in this language. If the order is invalid,
//   return "". If there are multiple valid orderings, return any.
//
// PARAMETERS:
//   words []string — list of words sorted in alien lexicographic order
//
// RETURN:
//   string — a string of unique letters in alien-sorted order, or "" if invalid
//
// CONSTRAINTS:
//   • 1 <= len(words) <= 100
//   • 1 <= len(words[i]) <= 100
//   • words[i] consists of only lowercase English letters
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  words = ["wrt","wrf","er","ett","rftt"]
//   Output: "wertf"
//   Why:    From adjacent comparisons: w<e, e<r, t<f, r<t → "wertf".
//
// Example 2:
//   Input:  words = ["z","x"]
//   Output: "zx"
//
// Example 3:
//   Input:  words = ["z","x","z"]
//   Output: ""
//   Why:    Invalid ordering (cycle: z before x, then x before z).
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Build a directed graph from adjacent word comparisons
// • Topological sort (BFS with Kahn's algorithm or DFS)
// • Detect invalid input: prefix case ("abc" before "ab") and cycles
// • Target: O(C) time, O(1) space (fixed 26-letter alphabet) where C = total chars
func AlienOrder(words []string) string {
	return ""
}
