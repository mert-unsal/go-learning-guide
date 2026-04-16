package backtracking

// ============================================================
// PROBLEM 8: Subsets II (LeetCode #90) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums that may contain duplicates,
//   return all possible subsets (the power set). The solution
//   set must not contain duplicate subsets. Return in any order.
//
// PARAMETERS:
//   nums []int — array of integers (may contain duplicates)
//
// RETURN:
//   [][]int — all unique subsets of nums
//
// CONSTRAINTS:
//   • 1 <= len(nums) <= 10
//   • -10 <= nums[i] <= 10
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [1,2,2]
//   Output: [[],[1],[1,2],[1,2,2],[2],[2,2]]
//   Why:    Duplicate subsets like [2] appearing twice are excluded.
//
// Example 2:
//   Input:  nums = [0]
//   Output: [[],[0]]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Sort nums first to group duplicates
// • Backtracking: skip nums[i] if nums[i] == nums[i-1] and i > start
// • Same dedup pattern as Combination Sum II
// • Target: O(n * 2^n) time, O(n) space for recursion depth
func SubsetsWithDup(nums []int) [][]int {
	return nil
}
