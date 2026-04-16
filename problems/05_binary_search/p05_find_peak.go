package binary_search

// ============================================================
// PROBLEM 5: Find Peak Element (LeetCode #162) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   A peak element is an element that is strictly greater than its
//   neighbors. Given a 0-indexed integer array nums, find a peak
//   element and return its index. If the array contains multiple
//   peaks, return the index to any of the peaks. You may imagine
//   that nums[-1] = nums[n] = -∞. You must achieve O(log n) time.
//
// PARAMETERS:
//   nums []int — an integer array where no two adjacent elements are equal
//
// RETURN:
//   int — the index of any peak element
//
// CONSTRAINTS:
//   • 1 <= len(nums) <= 1000
//   • -2^31 <= nums[i] <= 2^31 - 1
//   • nums[i] != nums[i+1] for all valid i
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [1,2,3,1]
//   Output: 2
//   Why:    nums[2] = 3 is a peak (3 > 2 and 3 > 1).
//
// Example 2:
//   Input:  nums = [1,2,1,3,5,6,4]
//   Output: 5 (or 1)
//   Why:    nums[5] = 6 is a peak. nums[1] = 2 is also a valid answer.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Binary search: if nums[mid] < nums[mid+1], a peak must exist to
//   the right (move left = mid+1). Otherwise, a peak is at mid or left.
// • The key insight is that ascending slope guarantees a peak ahead.
// • Target: O(log n) time, O(1) space

func FindPeakElement(nums []int) int {
	return 0
}
