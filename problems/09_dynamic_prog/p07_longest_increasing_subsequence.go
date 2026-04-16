package dynamic_prog

// ============================================================
// PROBLEM 7: Longest Increasing Subsequence (LeetCode #300) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums, return the length of the longest
//   strictly increasing subsequence.
//
// PARAMETERS:
//   nums []int — array of integers
//
// RETURN:
//   int — length of the longest strictly increasing subsequence
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 2500
//   • -10⁴ ≤ nums[i] ≤ 10⁴
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [10, 9, 2, 5, 3, 7, 101, 18]
//   Output: 4
//   Why:    [2, 3, 7, 101] is the longest increasing subsequence
//
// Example 2:
//   Input:  nums = [0, 1, 0, 3, 2, 3]
//   Output: 4
//   Why:    [0, 1, 2, 3]
//
// Example 3:
//   Input:  nums = [7, 7, 7, 7, 7]
//   Output: 1
//   Why:    All elements are equal — strictly increasing requires length 1
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • O(n²) DP: dp[i] = length of LIS ending at index i
// • Recurrence: dp[i] = max(dp[j]+1) for all j < i where nums[j] < nums[i]
// • O(n log n) patience sorting: maintain a tails array, binary search for position
// • Target: O(n²) time, O(n) space (or O(n log n) with binary search)
func LengthOfLIS(nums []int) int {
	return 0
}
