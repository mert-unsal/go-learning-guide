package strings_problems

// ============================================================
// PROBLEM 10: Encode and Decode Strings (LeetCode #271) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Design an algorithm to encode a list of strings to a single string
//   and decode it back. The encoded string must handle ANY characters
//   including delimiters and special characters.
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: ["hello","world"] → encode → decode → ["hello","world"]
// Example 2: ["",""] → encode → decode → ["",""]  (empty strings)
// Example 3: ["contains#delimiter"] → must handle '#' in content
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Length-prefix encoding: "5#hello5#world"
//   • Each string is prefixed by its length and a delimiter '#'.
//   • On decode, read length, skip '#', read that many chars.
//   • This handles ANY characters because length tells you exactly
//     how many bytes to read — no ambiguity.

// Encode encodes a list of strings to a single string.
func Encode(strs []string) string {
	// TODO: implement
	return ""
}

// Decode decodes the encoded string back to a list of strings.
func Decode(s string) []string {
	// TODO: implement
	return nil
}
