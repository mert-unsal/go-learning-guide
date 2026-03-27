package strings_problems

// ============================================================
// PROBLEM 4: Longest Common Prefix (LeetCode #14) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Write a function to find the longest common prefix string amongst
//   an array of strings. If there is no common prefix, return "".
//
// PARAMETERS:
//   strs []string — an array of strings.
//
// RETURN:
//   string — the longest common prefix, or "" if none.
//
// CONSTRAINTS:
//   • 1 <= strs.length <= 200
//   • 0 <= strs[i].length <= 200
//   • strs[i] consists of only lowercase English letters.
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1: ["flower","flow","flight"] → "fl"
// Example 2: ["dog","racecar","car"]    → ""  (no common prefix)
// Example 3: ["abc","abc","abc"]        → "abc"
// Example 4: ["", "b"]                  → ""
// Example 5: ["a"]                      → "a"
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Use the first string as a reference. Compare character by character.
//   • Stop as soon as a mismatch is found or any string runs out.
//   • Target: O(S) where S = total characters across all strings, O(1) space.

// LongestCommonPrefix returns the longest common prefix of strs.
// Time: O(S) where S = total characters  Space: O(1)
func LongestCommonPrefix(strs []string) string {
	// TODO: implement
	return ""
}
