package hard

// ============================================================
// PROBLEM 4: Longest Valid Parentheses (LeetCode #32) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a string containing just the characters '(' and ')',
//   return the length of the longest valid (well-formed)
//   parentheses substring.
//
// PARAMETERS:
//   s string — input string consisting of '(' and ')' characters only
//
// RETURN:
//   int — length of the longest valid parentheses substring
//
// CONSTRAINTS:
//   • 0 <= len(s) <= 3 * 10^4
//   • s[i] is '(' or ')'
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  s = "(()"
//   Output: 2
//   Why:    The longest valid parentheses substring is "()".
//
// Example 2:
//   Input:  s = ")()())"
//   Output: 4
//   Why:    The longest valid parentheses substring is "()()".
//
// Example 3:
//   Input:  s = ""
//   Output: 0
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Stack approach: push indices, pop on match, length = i - stack.top
// • DP approach: dp[i] = length of longest valid ending at i
// • Two-pass counter approach: scan left-to-right then right-to-left
// • Target: O(n) time, O(n) space (stack) or O(1) space (counter)
func LongestValidParentheses(s string) int {
	return 0
}
