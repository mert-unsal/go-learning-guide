package binary_search

// ============================================================
// PROBLEM 1: Search in Rotated Sorted Array (LeetCode #33) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   An integer array nums sorted in ascending order (with distinct
//   values) is rotated at an unknown pivot. Given the rotated array
//   and a target value, return the index of target, or -1 if not found.
//   You must achieve O(log n) runtime complexity.
//
// PARAMETERS:
//   nums   []int — a rotated sorted array of distinct integers
//   target int   — the value to search for
//
// RETURN:
//   int — the index of target in nums, or -1 if not present
//
// CONSTRAINTS:
//   • 1 <= len(nums) <= 5000
//   • -10^4 <= nums[i] <= 10^4
//   • All values in nums are unique
//   • nums was sorted ascending then rotated 1 to n times
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [4,5,6,7,0,1,2], target = 0
//   Output: 4
//   Why:    0 is at index 4 in the rotated array.
//
// Example 2:
//   Input:  nums = [4,5,6,7,0,1,2], target = 3
//   Output: -1
//   Why:    3 is not present in the array.
//
// Example 3:
//   Input:  nums = [1], target = 0
//   Output: -1
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • At each step, one half of the array is always sorted. Determine
//   which half is sorted, then check if target lies within that range.
// • If target is in the sorted half, narrow to it; otherwise, search
//   the other half.
// • Target: O(log n) time, O(1) space

func Search(nums []int, target int) int {
	return 0
}
