package stacks_queues

// ============================================================
// PROBLEM 9: Generate Parentheses (LeetCode #22) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given n pairs of parentheses, write a function to generate all
//   combinations of well-formed parentheses.
//
// PARAMETERS:
//   n int — the number of pairs of parentheses
//
// RETURN:
//   []string — all valid combinations of n pairs of parentheses
//
// CONSTRAINTS:
//   • 1 <= n <= 8
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  n = 3
//   Output: ["((()))","(()())","(())()","()(())","()()()"]
//   Why:    All 5 valid orderings of 3 pairs (Catalan number C(3) = 5)
//
// Example 2:
//   Input:  n = 1
//   Output: ["()"]
//   Why:    Only one way to form a single pair
//
// Example 3:
//   Input:  n = 2
//   Output: ["(())","()()"]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Backtracking: add '(' if open count < n, add ')' if close count < open count.
// • The total number of valid combinations is the n-th Catalan number.
// • Target: O(4^n / √n) time (Catalan number), O(n) recursion space

func GenerateParenthesis(n int) []string {
	return nil
}
