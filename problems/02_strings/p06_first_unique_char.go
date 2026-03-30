package strings_problems

// ============================================================
// PROBLEM 6: First Unique Character in a String (LeetCode #387) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a string s, find the first non-repeating character and return
//   its index. Return -1 if it does not exist.
//
// CONSTRAINTS:
//   • 1 <= s.length <= 10⁵
//   • s consists of only lowercase English letters.
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: s="leetcode"  → 0  ('l' is first unique)
// Example 2: s="loveleetcode" → 2  ('v')
// Example 3: s="aabb"      → -1  (no unique character)
// Example 4: s="z"         → 0
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Two-pass approach: count frequencies, then find first with count 1.
//   • A [26]int frequency array is enough for lowercase letters.
//   • Target: O(n) time, O(1) space.

// FirstUniqChar returns the index of the first non-repeating character.
// Time: O(n)  Space: O(1)
func FirstUniqChar(s string) int {
	return 0
}
