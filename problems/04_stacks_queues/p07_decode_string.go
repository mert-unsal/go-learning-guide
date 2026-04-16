package stacks_queues

// ============================================================
// PROBLEM 7: Decode String (LeetCode #394) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an encoded string, return its decoded string. The encoding
//   rule is: k[encoded_string], where the encoded_string inside the
//   brackets is repeated exactly k times. Nesting is allowed.
//
// PARAMETERS:
//   s string — an encoded string following the k[...] pattern
//
// RETURN:
//   string — the fully decoded string
//
// CONSTRAINTS:
//   • 1 <= len(s) <= 30
//   • s consists of lowercase English letters, digits, and '[]'
//   • s is guaranteed to be a valid encoding (balanced brackets, valid k)
//   • 1 <= k <= 300
//   • The decoded string length will not exceed 10^5
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  s = "3[a]2[bc]"
//   Output: "aaabcbc"
//   Why:    "a" repeated 3 times + "bc" repeated 2 times
//
// Example 2:
//   Input:  s = "3[a2[c]]"
//   Output: "accaccacc"
//   Why:    Inner: 2[c]→"cc", so 3[acc]→"accaccacc" (nested decoding)
//
// Example 3:
//   Input:  s = "2[abc]3[cd]ef"
//   Output: "abcabccdcdcdef"
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Use two stacks: one for repeat counts, one for partial strings
//   built so far. On '[', push current state. On ']', pop and repeat.
// • Alternatively, use recursion treating '[' as entering a sub-problem.
// • Target: O(output length) time and space

func DecodeString(s string) string {
	return ""
}
