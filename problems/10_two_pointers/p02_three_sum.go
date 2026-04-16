package two_pointers

// ============================================================
// PROBLEM 2: 3Sum (LeetCode #15) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums, return all the triplets
//   [nums[i], nums[j], nums[k]] such that i ≠ j ≠ k and
//   nums[i] + nums[j] + nums[k] == 0. The solution set must
//   not contain duplicate triplets.
//
// PARAMETERS:
//   nums []int — array of integers
//
// RETURN:
//   [][]int — list of unique triplets that sum to zero
//
// CONSTRAINTS:
//   • 3 ≤ len(nums) ≤ 3000
//   • -10⁵ ≤ nums[i] ≤ 10⁵
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [-1, 0, 1, 2, -1, -4]
//   Output: [[-1, -1, 2], [-1, 0, 1]]
//
// Example 2:
//   Input:  nums = [0, 1, 1]
//   Output: []
//   Why:    No triplet sums to zero
//
// Example 3:
//   Input:  nums = [0, 0, 0]
//   Output: [[0, 0, 0]]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Sort the array first, then fix one element and use two pointers for the rest
// • Skip duplicate values for the fixed element and both pointers
// • For fixed nums[i], find j,k such that nums[j]+nums[k] == -nums[i]
// • Early termination: if nums[i] > 0, no valid triplet possible
// • Target: O(n²) time, O(1) extra space (excluding output)
func ThreeSum(nums []int) [][]int {
	return nil
}
