package stacks_queues

// ============================================================
// PROBLEM 1: Valid Parentheses (LeetCode #20) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a string s containing only the characters '(', ')', '{', '}',
//   '[' and ']', determine if the input string is valid. A string is
//   valid if every open bracket is closed by the same type in the
//   correct order, and every close bracket has a corresponding open.
//
// PARAMETERS:
//   s string — a string consisting of bracket characters only
//
// RETURN:
//   bool — true if the bracket string is valid, false otherwise
//
// CONSTRAINTS:
//   • 1 <= len(s) <= 10^4
//   • s consists of parentheses only: '()[]{}' characters
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  s = "()"
//   Output: true
//   Why:    single pair of matching parentheses
//
// Example 2:
//   Input:  s = "()[]{}"
//   Output: true
//   Why:    three pairs, each correctly matched and ordered
//
// Example 3:
//   Input:  s = "(]"
//   Output: false
//   Why:    opening '(' is closed by ']' — type mismatch
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Push opening brackets onto a stack. On closing bracket, check
//   that the stack is non-empty and the top matches the expected type.
// • A map from closing → opening bracket simplifies the match check.
// • Target: O(n) time, O(n) space

func IsValid(s string) bool {
	return false
}
