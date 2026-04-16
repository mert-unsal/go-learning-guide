package two_pointers

// ============================================================
// PROBLEM 9: Minimum Difference Between Highest and Lowest of K Scores (LeetCode #1984) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array of integers nums representing test scores and an
//   integer k, return the minimum possible difference between the
//   highest and lowest of k scores chosen from the array.
//
// PARAMETERS:
//   nums []int — array of test scores
//   k    int   — number of scores to select
//
// RETURN:
//   int — minimum difference between max and min of the k chosen scores
//
// CONSTRAINTS:
//   • 1 ≤ k ≤ len(nums) ≤ 1000
//   • 0 ≤ nums[i] ≤ 10⁵
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [90], k = 1
//   Output: 0
//   Why:    Only one score selected — difference is 0
//
// Example 2:
//   Input:  nums = [9, 4, 1, 7], k = 2
//   Output: 2
//   Why:    Sort → [1, 4, 7, 9]; pick [7, 9] → diff = 2
//
// Example 3:
//   Input:  nums = [87, 68, 91, 86, 58, 63, 43, 98, 6, 40], k = 6
//   Output: 31
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Sort the array — the optimal k elements will be a contiguous window
// • Slide a window of size k across the sorted array
// • Minimum difference = min(nums[i+k-1] - nums[i]) for all valid i
// • Target: O(n log n) time, O(1) extra space
func MinimumDifference(nums []int, k int) int {
	return 0
}
