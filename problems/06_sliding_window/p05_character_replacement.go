package sliding_window

// ============================================================
// PROBLEM 5: Longest Repeating Character Replacement (LeetCode #424) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a string s and an integer k, you can choose any character
//   of the string and change it to any other uppercase English letter.
//   You can perform this operation at most k times. Return the length
//   of the longest substring containing the same letter you can get
//   after performing the above operations.
//
// PARAMETERS:
//   s string — a string of uppercase English letters
//   k int    — the maximum number of character replacements allowed
//
// RETURN:
//   int — the length of the longest valid substring after at most k replacements
//
// CONSTRAINTS:
//   • 1 <= len(s) <= 10^5
//   • s consists of only uppercase English letters
//   • 0 <= k <= len(s)
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  s = "ABAB", k = 2
//   Output: 4
//   Why:    Replace both 'A's with 'B' (or vice versa) → "BBBB", length 4.
//
// Example 2:
//   Input:  s = "AABABBA", k = 1
//   Output: 4
//   Why:    Replace the 'B' at index 3 → "AAAAABA" → substring "AAAA", length 4.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Sliding window: a window is valid when (windowSize - maxFreq) <= k,
//   meaning the characters that need replacement don't exceed k.
// • Track the frequency of each character in the window and the
//   running maximum frequency. Shrink left when the window is invalid.
// • Target: O(n) time, O(1) space (26-letter frequency array)

func CharacterReplacement(s string, k int) int {
	return 0
}
