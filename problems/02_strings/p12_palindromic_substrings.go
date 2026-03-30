package strings_problems

// ============================================================
// PROBLEM 12: Palindromic Substrings (LeetCode #647) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a string s, return the NUMBER of palindromic substrings in it.
//   A single character is always a palindrome.
//   Substrings with different start/end positions are counted separately
//   even if they have the same content.
//
// CONSTRAINTS:
//   • 1 <= s.length <= 1000
//   • s consists of lowercase English letters.
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: s="abc" → 3  ("a","b","c")
// Example 2: s="aaa" → 6  ("a"×3, "aa"×2, "aaa"×1)
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Same expand-around-center technique as Longest Palindromic Substring.
//   • Instead of tracking the longest, just COUNT each valid expansion.
//   • Target: O(n²) time, O(1) space.

// CountSubstrings counts palindromic substrings.
// Time: O(n²)  Space: O(1)
func CountSubstrings(s string) int {
	return 0
}
