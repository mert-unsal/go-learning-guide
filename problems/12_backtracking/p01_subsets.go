package backtracking

// ============================================================
// PROBLEM 1: Subsets (LeetCode #78) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums of unique elements, return all
//   possible subsets (the power set). The solution set must not
//   contain duplicate subsets. Return the subsets in any order.
//
// PARAMETERS:
//   nums []int — array of unique integers
//
// RETURN:
//   [][]int — all possible subsets of nums
//
// CONSTRAINTS:
//   • 1 <= len(nums) <= 10
//   • -10 <= nums[i] <= 10
//   • All elements of nums are unique
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [1,2,3]
//   Output: [[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
//   Why:    All 2^3 = 8 subsets of {1,2,3}.
//
// Example 2:
//   Input:  nums = [0]
//   Output: [[],[0]]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Backtracking: at each index, choose to include or exclude the element
// • Iterative: start with [[]], for each num append it to all existing subsets
// • Bitmask: iterate 0..2^n-1, each bit decides inclusion
// • Target: O(n * 2^n) time, O(n * 2^n) space
func Subsets(nums []int) [][]int {
	return nil
}
