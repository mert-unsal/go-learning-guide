package two_pointers

// ============================================================
// PROBLEM 8: 4Sum (LeetCode #18) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array nums of n integers and an integer target, return
//   all unique quadruplets [nums[a], nums[b], nums[c], nums[d]] such
//   that a, b, c, d are distinct indices and the four values sum to
//   target. The solution set must not contain duplicate quadruplets.
//
// PARAMETERS:
//   nums   []int — array of integers
//   target int   — target sum for the quadruplet
//
// RETURN:
//   [][]int — list of unique quadruplets that sum to target
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 200
//   • -10⁹ ≤ nums[i] ≤ 10⁹
//   • -10⁹ ≤ target ≤ 10⁹
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [1, 0, -1, 0, -2, 2], target = 0
//   Output: [[-2, -1, 1, 2], [-2, 0, 0, 2], [-1, 0, 0, 1]]
//
// Example 2:
//   Input:  nums = [2, 2, 2, 2, 2], target = 8
//   Output: [[2, 2, 2, 2]]
//
// Example 3:
//   Input:  nums = [1, 2, 3, 4], target = 21
//   Output: []
//   Why:    No quadruplet sums to 21
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Sort, then extend 3Sum: fix two elements, use two pointers for the remaining pair
// • Skip duplicates at each level to avoid repeated quadruplets
// • Watch for integer overflow when summing four values — use int (64-bit in Go)
// • Target: O(n³) time, O(1) extra space (excluding output)
func FourSum(nums []int, target int) [][]int {
	return nil
}
