package strings_problems

// ============================================================
// PROBLEM 2: Longest Substring Without Repeating Characters (LeetCode #3) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a string s, find the length of the LONGEST SUBSTRING
//   without repeating characters.
//
//   A substring is a contiguous sequence of characters within the string.
//
// PARAMETERS:
//   s string — the input string (may contain letters, digits, symbols, spaces).
//
// RETURN:
//   int — the length of the longest substring with all unique characters.
//
// CONSTRAINTS:
//   • 0 <= s.length <= 5 × 10⁴
//   • s consists of English letters, digits, symbols, and spaces.
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1: s = "abcabcbb" → 3  (substring "abc")
// Example 2: s = "bbbbb"    → 1  (substring "b")
// Example 3: s = "pwwkew"   → 3  (substring "wke", NOT "pwke" — must be contiguous)
// Example 4: s = ""         → 0
// Example 5: s = " "        → 1  (space is a valid character)
// Example 6: s = "dvdf"     → 3  (substring "vdf")
// Example 7: s = "abba"     → 2  ("ab" or "ba")
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • This is a classic SLIDING WINDOW problem.
//   • Two pointers define a window [left, right].
//   • Expand right. When a duplicate is found, shrink from left.
//   • A map from character → last-seen index lets you jump the left
//     pointer directly instead of moving one-by-one.
//   • Target: O(n) time, O(min(n, alphabet_size)) space.

// LengthOfLongestSubstring returns the length of the longest unique-char substring.
// Time: O(n)  Space: O(min(n, alphabet_size))
func LengthOfLongestSubstring(s string) int {
	return 0
}
