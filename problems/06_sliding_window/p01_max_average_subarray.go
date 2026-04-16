package sliding_window

// ============================================================
// PROBLEM 1: Maximum Average Subarray I (LeetCode #643) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums and an integer k, find a contiguous
//   subarray of length k that has the maximum average value and
//   return that value.
//
// PARAMETERS:
//   nums []int — the input array of integers
//   k    int   — the length of the subarray
//
// RETURN:
//   float64 — the maximum average value of any subarray of length k
//
// CONSTRAINTS:
//   • n == len(nums)
//   • 1 <= k <= n <= 10^5
//   • -10^4 <= nums[i] <= 10^4
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [1,12,-5,-6,50,3], k = 4
//   Output: 12.75
//   Why:    Subarray [12,-5,-6,50] has sum 51, average 51/4 = 12.75.
//
// Example 2:
//   Input:  nums = [5], k = 1
//   Output: 5.0
//   Why:    Only one element, average is the element itself.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Fixed-size sliding window of length k. Compute the initial window
//   sum, then slide by adding the new right element and subtracting
//   the element leaving the window.
// • Track the maximum sum seen and divide by k at the end.
// • Target: O(n) time, O(1) space

func FindMaxAverage(nums []int, k int) float64 {
	return 0
}
