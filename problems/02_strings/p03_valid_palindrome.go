package strings_problems

// ============================================================
// PROBLEM 3: Valid Palindrome (LeetCode #125) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   A phrase is a palindrome if, after converting all uppercase letters
//   to lowercase and removing all non-alphanumeric characters, it reads
//   the same forward and backward. Alphanumeric = letters and numbers.
//
// PARAMETERS:
//   s string — the input string (may contain spaces, punctuation, mixed case).
//
// RETURN:
//   bool — true if s is a valid palindrome after cleaning.
//
// CONSTRAINTS:
//   • 1 <= s.length <= 2 × 10⁵
//   • s consists only of printable ASCII characters.
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1: "A man, a plan, a canal: Panama" → true  ("amanaplanacanalpanama")
// Example 2: "race a car" → false  ("raceacar")
// Example 3: " " → true  (empty after removing non-alphanumeric)
// Example 4: "0P" → false  ("0p" is not a palindrome)
// Example 5: "aa" → true
// Example 6: ".,," → true  (empty after cleaning)
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Two pointers converging from both ends.
//   • Skip non-alphanumeric characters.
//   • Compare lowercase versions of characters.
//   • Target: O(n) time, O(1) space.

// IsPalindrome returns true if s is a valid palindrome (ignoring case/non-alnum).
// Time: O(n)  Space: O(1)
func IsPalindrome(s string) bool {
	return false
}
