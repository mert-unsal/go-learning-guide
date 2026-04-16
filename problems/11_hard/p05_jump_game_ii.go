package hard

// ============================================================
// PROBLEM 5: Jump Game II (LeetCode #45) — MEDIUM/HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   You are given a 0-indexed array of integers nums of length n.
//   You are initially positioned at nums[0]. Each element nums[i]
//   represents the maximum length of a forward jump from index i.
//   Return the minimum number of jumps to reach nums[n-1].
//   The test cases are generated such that you can always reach the last index.
//
// PARAMETERS:
//   nums []int — array where nums[i] is the max jump length from index i
//
// RETURN:
//   int — minimum number of jumps to reach the last index
//
// CONSTRAINTS:
//   • 1 <= len(nums) <= 10^4
//   • 0 <= nums[i] <= 1000
//   • It is guaranteed that you can reach nums[n-1]
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [2,3,1,1,4]
//   Output: 2
//   Why:    Jump 1 step from index 0 to 1, then 3 steps to index 4.
//
// Example 2:
//   Input:  nums = [2,3,0,1,4]
//   Output: 2
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Greedy BFS: track the farthest reachable index in the current "level"
// • Increment jumps when you pass the current level boundary
// • Target: O(n) time, O(1) space
func Jump(nums []int) int {
	return 0
}
