package strings_problems

import "strings"

// ============================================================
// Pangrams — [E]
// ============================================================
// A pangram contains every letter of the alphabet at least once.
// Determine if the sentence is a pangram.
//
// Example: "We promptly judged antique ivory buckles for the next prize" → "pangram"

// Pangram returns "pangram" if every letter appears, else "not pangram".
// Time: O(n)  Space: O(1)
func Pangram(sentence string) string {
	var seen [26]bool
	for _, ch := range strings.ToLower(sentence) {
		if ch >= 'a' && ch <= 'z' {
			seen[ch-'a'] = true
		}
	}
	for _, v := range seen {
		if !v {
			return "not pangram"
		}
	}
	return "pangram"
}
