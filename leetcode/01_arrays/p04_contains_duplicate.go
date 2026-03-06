package arrays

// ============================================================
// PROBLEM 4: Contains Duplicate (LeetCode #217) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums, return true if any value appears
//   AT LEAST TWICE in the array, and return false if every element
//   is distinct.
//
// PARAMETERS:
//   nums []int — an array of integers.
//
// RETURN:
//   bool — true if any duplicate exists, false otherwise.
//
// CONSTRAINTS:
//   • 1 <= nums.length <= 10⁵
//   • -10⁹ <= nums[i] <= 10⁹
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Has duplicate:
//   Input:  nums = [1, 2, 3, 1]
//   Output: true
//   Why:    The value 1 appears at index 0 and index 3.
//
// Example 2 — All unique:
//   Input:  nums = [1, 2, 3, 4]
//   Output: false
//   Why:    Every element is different.
//
// Example 3 — Multiple duplicates:
//   Input:  nums = [1, 1, 1, 3, 3, 4, 3, 2, 4, 2]
//   Output: true
//
// Example 4 — Single element:
//   Input:  nums = [7]
//   Output: false
//   Why:    Can't have a duplicate with only one element.
//
// Example 5 — Two identical elements:
//   Input:  nums = [5, 5]
//   Output: true
//
// Example 6 — Negative numbers:
//   Input:  nums = [-1, -2, -3, -1]
//   Output: true
//   Why:    -1 appears twice.
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Brute force: compare every pair → O(n²). Can you do better?
//   • Sorting first lets you check adjacent elements → O(n log n).
//   • A hash set lets you check "have I seen this before?" in O(1).
//   • What data structure in Go acts as a set? (Hint: map[int]bool)
//   • Target: O(n) time, O(n) space.

// ContainsDuplicate returns true if the slice has any repeated element.
// Time: O(n)  Space: O(n)
func ContainsDuplicate(nums []int) bool {
	// TODO: implement
	return false
}
