package binary_search

// ============================================================
// PROBLEM 9: Find First and Last Position (LeetCode #34) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array of integers nums sorted in non-decreasing order,
//   find the starting and ending position of a given target value.
//   If target is not found, return [-1, -1]. You must achieve
//   O(log n) runtime complexity.
//
// PARAMETERS:
//   nums   []int — a sorted (non-decreasing) array of integers
//   target int   — the value to find the range for
//
// RETURN:
//   []int — a two-element slice [first, last] position of target, or [-1, -1]
//
// CONSTRAINTS:
//   • 0 <= len(nums) <= 10^5
//   • -10^9 <= nums[i] <= 10^9
//   • nums is sorted in non-decreasing order
//   • -10^9 <= target <= 10^9
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [5,7,7,8,8,10], target = 8
//   Output: [3,4]
//   Why:    8 first appears at index 3 and last at index 4.
//
// Example 2:
//   Input:  nums = [5,7,7,8,8,10], target = 6
//   Output: [-1,-1]
//   Why:    6 is not present in the array.
//
// Example 3:
//   Input:  nums = [], target = 0
//   Output: [-1,-1]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Run two binary searches: one biased left (find first occurrence)
//   and one biased right (find last occurrence).
// • For leftmost: when nums[mid] == target, move right = mid.
//   For rightmost: when nums[mid] == target, move left = mid.
// • Target: O(log n) time, O(1) space

func SearchRange(nums []int, target int) []int {
	return nil
}
