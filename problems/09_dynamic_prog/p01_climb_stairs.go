package dynamic_prog

// ============================================================
// PROBLEM 1: Climbing Stairs (LeetCode #70) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   You are climbing a staircase. It takes n steps to reach the top.
//   Each time you can either climb 1 or 2 steps. Return the number
//   of distinct ways you can climb to the top.
//
// PARAMETERS:
//   n int — the total number of steps to reach the top (1 ≤ n ≤ 45)
//
// RETURN:
//   int — the number of distinct ways to climb to the top
//
// CONSTRAINTS:
//   • 1 ≤ n ≤ 45
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  n = 2
//   Output: 2
//   Why:    (1+1) or (2) — two distinct ways
//
// Example 2:
//   Input:  n = 3
//   Output: 3
//   Why:    (1+1+1), (1+2), (2+1) — three distinct ways
//
// Example 3:
//   Input:  n = 5
//   Output: 8
//   Why:    Fibonacci pattern — ways(n) = ways(n-1) + ways(n-2)
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • This is the Fibonacci sequence: ways(n) = ways(n-1) + ways(n-2)
// • Base cases: ways(1) = 1, ways(2) = 2
// • You only need the previous two values — no array required
// • Target: O(n) time, O(1) space

func ClimbStairs(n int) int {
	return 0
}
