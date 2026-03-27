package strings_problems

// ============================================================
// PROBLEM 1: Valid Anagram (LeetCode #242) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given two strings s and t, return true if t is an ANAGRAM of s,
//   and false otherwise.
//
//   An anagram is a word formed by rearranging the letters of another
//   word, using all the original letters exactly once.
//
// PARAMETERS:
//   s string — first string.
//   t string — second string.
//
// RETURN:
//   bool — true if t is an anagram of s.
//
// CONSTRAINTS:
//   • 1 <= s.length, t.length <= 5 × 10⁴
//   • s and t consist of lowercase English letters.
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1: s="anagram", t="nagaram" → true
// Example 2: s="rat",     t="car"     → false
// Example 3: s="a",       t="a"       → true
// Example 4: s="ab",      t="a"       → false (different lengths)
// Example 5: s="aacc",    t="ccac"    → false (same length, different frequencies)
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Two strings are anagrams if and only if they have the same
//     character frequencies.
//   • With only lowercase English letters, a [26]int array suffices.
//   • Increment for s, decrement for t. If all zero → anagram.
//   • Target: O(n) time, O(1) space (26-letter array).

// IsAnagram returns true if s and t are anagrams.
// Time: O(n)  Space: O(1) — only 26 lowercase letters
func IsAnagram(s string, t string) bool {
	// TODO: implement
	return false
}
