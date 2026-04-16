package binary_search

// ============================================================
// PROBLEM 3: Binary Search (LeetCode #704) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a sorted array of integers nums and an integer target,
//   return the index of target if it exists, otherwise return -1.
//   You must write an algorithm with O(log n) runtime complexity.
//
// PARAMETERS:
//   nums   []int — a sorted (ascending) array of integers
//   target int   — the value to search for
//
// RETURN:
//   int — index of target in nums, or -1 if not found
//
// CONSTRAINTS:
//   • 1 <= len(nums) <= 10^4
//   • -10^4 < nums[i], target < 10^4
//   • All integers in nums are unique
//   • nums is sorted in ascending order
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [-1,0,3,5,9,12], target = 9
//   Output: 4
//   Why:    9 exists at index 4.
//
// Example 2:
//   Input:  nums = [-1,0,3,5,9,12], target = 2
//   Output: -1
//   Why:    2 does not exist in the array.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Classic binary search: maintain left and right pointers, compare
//   nums[mid] with target to halve the search space.
// • Use mid = left + (right - left) / 2 to avoid integer overflow.
// • Target: O(log n) time, O(1) space

func BinarySearch(nums []int, target int) int {
	return 0
}
