package strings_problems

// ============================================================
// PROBLEM 5: Reverse Words in a String (LeetCode #151) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an input string s, reverse the ORDER of the words.
//   A word is a sequence of non-space characters. Words are separated
//   by at least one space. Return the reversed string with single spaces
//   between words and no leading/trailing spaces.
//
// CONSTRAINTS:
//   • 1 <= s.length <= 10⁴
//   • s contains English letters (upper/lower), digits, and spaces.
//   • There is at least one word.
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: "the sky is blue"    → "blue is sky the"
// Example 2: "  hello world  "    → "world hello"
// Example 3: "a good   example"   → "example good a"
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Go's strings.Fields splits on any whitespace and ignores extras.
//   • After splitting, reverse the slice and join with single spaces.
//   • Target: O(n) time, O(n) space.

// ReverseWords reverses the word order in a string.
// Time: O(n)  Space: O(n)
func ReverseWords(s string) string {
	return ""
}
