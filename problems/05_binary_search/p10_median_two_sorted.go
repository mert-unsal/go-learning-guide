package binary_search

// ============================================================
// PROBLEM 10: Median of Two Sorted Arrays (LeetCode #4) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   Given two sorted arrays nums1 and nums2 of size m and n
//   respectively, return the median of the two sorted arrays.
//   The overall run time complexity should be O(log(min(m, n))).
//
// PARAMETERS:
//   nums1 []int — first sorted array
//   nums2 []int — second sorted array
//
// RETURN:
//   float64 — the median of the combined sorted arrays
//
// CONSTRAINTS:
//   • nums1.length == m, nums2.length == n
//   • 0 <= m <= 1000, 0 <= n <= 1000
//   • 1 <= m + n <= 2000
//   • -10^6 <= nums1[i], nums2[i] <= 10^6
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums1 = [1,3], nums2 = [2]
//   Output: 2.0
//   Why:    Merged = [1,2,3], median = 2.
//
// Example 2:
//   Input:  nums1 = [1,2], nums2 = [3,4]
//   Output: 2.5
//   Why:    Merged = [1,2,3,4], median = (2+3)/2 = 2.5.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Binary search on the partition of the smaller array. For a
//   partition i in nums1, the corresponding partition j in nums2 is
//   (m+n+1)/2 - i.
// • Valid partition: maxLeft1 <= minRight2 and maxLeft2 <= minRight1.
// • Target: O(log(min(m, n))) time, O(1) space

func FindMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	return 0
}
