package arrays

// FindMinRotated ============================================================
// PROBLEM 9: Find Minimum in Rotated Sorted Array (LeetCode #153) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//
//	Suppose an array of length n sorted in ascending order is ROTATED
//	between 1 and n times. For example, [0,1,2,4,5,6,7] might become
//	[4,5,6,7,0,1,2] after rotating 4 times.
//
//	Given the sorted rotated array nums of UNIQUE elements, return
//	the minimum element.
//
//	You must write an algorithm that runs in O(log n) time.
//
// PARAMETERS:
//
//	nums []int — a sorted array that has been rotated. All elements are unique.
//
// RETURN:
//
//	int — the minimum element in the array.
//
// CONSTRAINTS:
//   - n == nums.length
//   - 1 <= n <= 5000
//   - -5000 <= nums[i] <= 5000
//   - All values are UNIQUE.
//   - nums was sorted ascending, then rotated 1 to n times.
//
// ─── WHAT DOES "ROTATED" MEAN? ──────────────────────────────
//
//	Original sorted: [0, 1, 2, 3, 4, 5, 6, 7]
//	Rotated 4 times: [4, 5, 6, 7, 0, 1, 2, 3]
//	                  ← sorted → | ← sorted →
//
//	A rotated sorted array has TWO sorted halves.
//	The minimum is the first element of the second sorted half
//	(the "inflection point").
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Rotated in the middle:
//
//	Input:  nums = [3, 4, 5, 1, 2]
//	Output: 1
//	Why:    Original was [1,2,3,4,5], rotated 3 times.
//
// Example 2 — Rotated near the end:
//
//	Input:  nums = [4, 5, 6, 7, 0, 1, 2]
//	Output: 0
//
// Example 3 — Not actually rotated (rotated n times = original):
//
//	Input:  nums = [11, 13, 15, 17]
//	Output: 11
//	Why:    Already sorted. The minimum is just the first element.
//
// Example 4 — Two elements:
//
//	Input:  nums = [2, 1]
//	Output: 1
//
// Example 5 — Single element:
//
//	Input:  nums = [42]
//	Output: 42
//
// Example 6 — Rotated by 1:
//
//	Input:  nums = [7, 1, 2, 3, 4, 5, 6]
//	Output: 1
//	Why:    Only the largest element wrapped to the front.
//
// FindMinRotated returns the minimum of a rotated sorted array.
// Time: O(log n)  Space: O(1)
func FindMinRotated(nums []int) int {
	// TODO: implement
	return 0
}
