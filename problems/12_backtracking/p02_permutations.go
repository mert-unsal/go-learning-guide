package backtracking

// ============================================================
// PROBLEM 2: Permutations (LeetCode #46) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array nums of distinct integers, return all possible
//   permutations in any order.
//
// PARAMETERS:
//   nums []int — array of distinct integers
//
// RETURN:
//   [][]int — all possible permutations of nums
//
// CONSTRAINTS:
//   • 1 <= len(nums) <= 6
//   • -10 <= nums[i] <= 10
//   • All integers of nums are unique
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [1,2,3]
//   Output: [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
//   Why:    All 3! = 6 permutations of {1,2,3}.
//
// Example 2:
//   Input:  nums = [0,1]
//   Output: [[0,1],[1,0]]
//
// Example 3:
//   Input:  nums = [1]
//   Output: [[1]]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Backtracking: swap elements into the current position, recurse, swap back
// • Alternatively, use a "used" boolean array to track which elements are placed
// • Target: O(n * n!) time, O(n) space (excluding output)
func Permute(nums []int) [][]int {
	return nil
}
