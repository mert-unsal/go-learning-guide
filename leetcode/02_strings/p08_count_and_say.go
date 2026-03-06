package strings_problems

// ============================================================
// PROBLEM 8: Count and Say (LeetCode #38) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   The count-and-say sequence is a sequence of digit strings defined
//   by recursive description:
//     Term 1: "1"
//     Term n: describe term (n-1) by "saying" each group of consecutive
//             identical digits. E.g., "1" → "one 1" → "11".
//
//   "1"    → "11"   (one 1)
//   "11"   → "21"   (two 1s)
//   "21"   → "1211" (one 2, then one 1)
//   "1211" → "111221"
//
// CONSTRAINTS:
//   • 1 <= n <= 30
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: n=1 → "1"
// Example 2: n=4 → "1211"
// Example 3: n=5 → "111221"
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Build each term from the previous one.
//   • Walk through the previous string, counting consecutive identical chars.
//   • Use a strings.Builder for efficiency.
//   • Target: O(n × length of each term) time.

// CountAndSay returns the nth term of the count-and-say sequence.
// Time: O(n * len of each term)  Space: O(n)
func CountAndSay(n int) string {
	// TODO: implement
	return ""
}
