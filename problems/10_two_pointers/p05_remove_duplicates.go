package two_pointers

// ============================================================
// PROBLEM 5: Remove Duplicates from Sorted Array (LeetCode #26) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums sorted in non-decreasing order, remove
//   the duplicates in-place such that each unique element appears only
//   once. The relative order must be preserved. Return the number of
//   unique elements k. The first k elements of nums must hold the
//   unique values.
//
// PARAMETERS:
//   nums []int — sorted array of integers (modified in-place)
//
// RETURN:
//   int — k, the number of unique elements
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 3 × 10⁴
//   • -100 ≤ nums[i] ≤ 100
//   • nums is sorted in non-decreasing order
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [1, 1, 2]
//   Output: 2, nums = [1, 2, _]
//   Why:    Two unique elements; values beyond k don't matter
//
// Example 2:
//   Input:  nums = [0, 0, 1, 1, 1, 2, 2, 3, 3, 4]
//   Output: 5, nums = [0, 1, 2, 3, 4, _, _, _, _, _]
//
// Example 3:
//   Input:  nums = [1]
//   Output: 1
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Two pointers: slow marks the last unique position, fast scans ahead
// • When nums[fast] ≠ nums[slow], increment slow and copy
// • Since array is sorted, duplicates are always adjacent
// • Target: O(n) time, O(1) space
func RemoveDuplicates(nums []int) int {
	return 0
}
