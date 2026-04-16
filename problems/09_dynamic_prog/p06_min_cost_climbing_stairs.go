package dynamic_prog

// ============================================================
// PROBLEM 6: Min Cost Climbing Stairs (LeetCode #746) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   You are given an integer array cost where cost[i] is the cost of
//   the i-th step on a staircase. Once you pay the cost, you can climb
//   one or two steps. You can start from step 0 or step 1. Return the
//   minimum cost to reach the top of the floor (past the last step).
//
// PARAMETERS:
//   cost []int — cost of each step
//
// RETURN:
//   int — minimum cost to reach the top (beyond the last index)
//
// CONSTRAINTS:
//   • 2 ≤ len(cost) ≤ 1000
//   • 0 ≤ cost[i] ≤ 999
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  cost = [10, 15, 20]
//   Output: 15
//   Why:    Start at step 1 (cost 15), jump two steps to the top
//
// Example 2:
//   Input:  cost = [1, 100, 1, 1, 1, 100, 1, 1, 100, 1]
//   Output: 6
//   Why:    Steps 0→2→3→4→6→7→9→top, paying 1+1+1+1+1+1 = 6
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Recurrence: dp[i] = cost[i] + min(dp[i-1], dp[i-2])
// • Answer is min(dp[n-1], dp[n-2]) — you can step from either to the top
// • Only need previous two values — iterate in-place or with two vars
// • Target: O(n) time, O(1) space
func MinCostClimbingStairs(cost []int) int {
	return 0
}
func min2(a, b int) int {
	return 0
}
