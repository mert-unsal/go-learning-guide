package binary_search

// ============================================================
// PROBLEM 2: Find Minimum in Rotated Sorted Array (LeetCode #153) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a sorted array of unique elements that has been rotated
//   between 1 and n times, find the minimum element. You must
//   achieve O(log n) runtime complexity.
//
// PARAMETERS:
//   nums []int — a rotated sorted array of unique integers
//
// RETURN:
//   int — the minimum element in the array
//
// CONSTRAINTS:
//   • n == len(nums)
//   • 1 <= n <= 5000
//   • -5000 <= nums[i] <= 5000
//   • All values are unique
//   • nums was sorted ascending then rotated 1 to n times
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [3,4,5,1,2]
//   Output: 1
//   Why:    Original sorted array [1,2,3,4,5] was rotated 3 times.
//
// Example 2:
//   Input:  nums = [4,5,6,7,0,1,2]
//   Output: 0
//   Why:    The minimum element 0 is at the rotation pivot.
//
// Example 3:
//   Input:  nums = [11,13,15,17]
//   Output: 11
//   Why:    Array was rotated 4 times (full cycle), so it's still sorted.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • If nums[mid] > nums[right], the minimum is in the right half.
//   Otherwise, the minimum is in the left half (including mid).
// • Compare against the rightmost element to decide which half to search.
// • Target: O(log n) time, O(1) space

func FindMin(nums []int) int {
	return 0
}
