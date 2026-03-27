package strings_problems

import "strings"

// ============================================================
// Repeated String — [E]
// ============================================================
// Infinite string formed by repeating s. Count occurrences of 'a' in first n chars.
//
// Example: s="aba", n=10 → 7  ("abaabaabaaba"[:10] = "abaabaabaab" → 7 a's)

// RepeatedString counts 'a' in the first n characters of the infinite repeated string.
// Time: O(|s|)  Space: O(1)
func RepeatedString(s string, n int) int {
	countInS := strings.Count(s, "a")
	fullRepeats := n / len(s)
	remainder := n % len(s)
	countInRemainder := strings.Count(s[:remainder], "a")
	return fullRepeats*countInS + countInRemainder
}
