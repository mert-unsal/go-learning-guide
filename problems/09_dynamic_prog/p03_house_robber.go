package dynamic_prog

// ============================================================
// PROBLEM 3: House Robber (LeetCode #198) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   You are a robber planning to rob houses along a street. Each house
//   has a certain amount of money stashed. Adjacent houses have connected
//   security systems — if two adjacent houses are broken into on the same
//   night, the police are alerted. Return the maximum amount of money
//   you can rob without alerting the police.
//
// PARAMETERS:
//   nums []int — amount of money at each house (nums[i] ≥ 0)
//
// RETURN:
//   int — maximum amount you can rob without robbing adjacent houses
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 100
//   • 0 ≤ nums[i] ≤ 400
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [1, 2, 3, 1]
//   Output: 4
//   Why:    Rob house 0 (1) + house 2 (3) = 4
//
// Example 2:
//   Input:  nums = [2, 7, 9, 3, 1]
//   Output: 12
//   Why:    Rob house 0 (2) + house 2 (9) + house 4 (1) = 12
//
// Example 3:
//   Input:  nums = [2, 1, 1, 2]
//   Output: 4
//   Why:    Rob house 0 (2) + house 3 (2) = 4
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Recurrence: rob(i) = max(rob(i-1), rob(i-2) + nums[i])
// • Only need previous two values — no array required
// • Think of it as "take or skip" at each house
// • Target: O(n) time, O(1) space

func Rob(nums []int) int {
	return 0
}

func max(a, b int) int {
	return 0
}
