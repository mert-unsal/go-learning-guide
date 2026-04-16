package dynamic_prog

// ============================================================
// PROBLEM 9: Jump Game (LeetCode #55) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   You are given an integer array nums. You are initially positioned
//   at the first index, and each element represents your maximum jump
//   length from that position. Return true if you can reach the last index.
//
// PARAMETERS:
//   nums []int — array where nums[i] is the max jump length from index i
//
// RETURN:
//   bool — true if you can reach the last index
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 10⁴
//   • 0 ≤ nums[i] ≤ 10⁵
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [2, 3, 1, 1, 4]
//   Output: true
//   Why:    Jump 1 step to index 1, then 3 steps to the last index
//
// Example 2:
//   Input:  nums = [3, 2, 1, 0, 4]
//   Output: false
//   Why:    You always arrive at index 3 where nums[3] = 0, stuck
//
// Example 3:
//   Input:  nums = [0]
//   Output: true
//   Why:    Already at the last index
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Greedy: track the farthest reachable index as you scan left to right
// • If current index > farthest reachable, return false
// • Update farthest = max(farthest, i + nums[i])
// • Target: O(n) time, O(1) space
func CanJump(nums []int) bool {
	return false
}
