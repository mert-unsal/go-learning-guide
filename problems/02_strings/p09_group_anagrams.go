package strings_problems

// ============================================================
// PROBLEM 9: Group Anagrams (LeetCode #49) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array of strings strs, group the anagrams together.
//   You can return the answer in any order.
//
// CONSTRAINTS:
//   • 1 <= strs.length <= 10⁴
//   • 0 <= strs[i].length <= 100
//   • strs[i] consists of lowercase English letters.
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: ["eat","tea","tan","ate","nat","bat"]
//          → [["bat"],["nat","tan"],["ate","eat","tea"]]
// Example 2: [""] → [[""]]
// Example 3: ["a"] → [["a"]]
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Two strings are anagrams if they have the same character frequencies.
//   • Use a [26]int frequency array as a map key to group anagrams.
//   • In Go, arrays (not slices) are comparable and can be map keys.
//   • Target: O(n × k) time where k is max string length, O(n×k) space.

// GroupAnagrams groups strings that are anagrams of each other.
// Time: O(n * k) where k is max string length  Space: O(n*k)
func GroupAnagrams(strs []string) [][]string {
	return nil
}
