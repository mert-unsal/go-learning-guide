package strings_problems

// ============================================================
// PROBLEM 11: Longest Palindromic Substring (LeetCode #5) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a string s, return the longest palindromic substring in s.
//   A palindrome reads the same forward and backward.
//
// CONSTRAINTS:
//   • 1 <= s.length <= 1000
//   • s consists of only digits and English letters.
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: s="babad" → "bab" (or "aba")
// Example 2: s="cbbd"  → "bb"
// Example 3: s="a"     → "a"
// Example 4: s="ac"    → "a" (or "c")
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Expand around center: for each position, try expanding outward.
//   • 2n-1 centers: n single chars (odd-length) + n-1 gaps (even-length).
//   • Target: O(n²) time, O(1) space.

// LongestPalindrome returns the longest palindromic substring.
// Time: O(n²)  Space: O(1)
func LongestPalindrome(s string) string {
	// TODO: implement
	return ""
}
