package sliding_window

// ============================================================
// PROBLEM 3: Permutation in String (LeetCode #567) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given two strings s1 and s2, return true if s2 contains a
//   permutation of s1 (i.e., one of s1's permutations is a substring
//   of s2).
//
// PARAMETERS:
//   s1 string — the pattern string whose permutation to find
//   s2 string — the string to search within
//
// RETURN:
//   bool — true if any permutation of s1 is a substring of s2
//
// CONSTRAINTS:
//   • 1 <= len(s1), len(s2) <= 10^4
//   • s1 and s2 consist of lowercase English letters
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  s1 = "ab", s2 = "eidbaooo"
//   Output: true
//   Why:    s2 contains "ba" which is a permutation of "ab".
//
// Example 2:
//   Input:  s1 = "ab", s2 = "eidboaoo"
//   Output: false
//   Why:    No substring of length 2 in s2 is a permutation of "ab".
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Use a fixed-size sliding window of length len(s1) over s2.
// • Maintain a frequency array (26 lowercase letters) for the window
//   and compare against s1's frequency array at each position.
// • Alternatively, track a "matches" count of how many of the 26
//   character counts are equal between the two arrays.
// • Target: O(|s1| + |s2|) time, O(1) space (fixed 26-char array)

func CheckInclusion(s1 string, s2 string) bool {
	return false
}
