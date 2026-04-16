package dynamic_prog

// ============================================================
// PROBLEM 10: Partition Equal Subset Sum (LeetCode #416) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums, return true if you can partition the
//   array into two subsets such that the sum of elements in both
//   subsets is equal.
//
// PARAMETERS:
//   nums []int — array of positive integers
//
// RETURN:
//   bool — true if the array can be partitioned into two equal-sum subsets
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 200
//   • 1 ≤ nums[i] ≤ 100
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [1, 5, 11, 5]
//   Output: true
//   Why:    [1, 5, 5] and [11] both sum to 11
//
// Example 2:
//   Input:  nums = [1, 2, 3, 5]
//   Output: false
//   Why:    Total sum is 11 (odd) — cannot split evenly
//
// Example 3:
//   Input:  nums = [1, 2, 5]
//   Output: false
//   Why:    Total sum is 8, target = 4, but no subset sums to 4
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • If total sum is odd, return false immediately
// • Reduce to 0/1 knapsack: can we pick a subset that sums to totalSum/2?
// • DP: dp[j] = true if a subset sums to j, iterate nums in outer loop
// • Traverse j from target down to nums[i] to avoid reusing elements
// • Target: O(n × sum/2) time, O(sum/2) space
func CanPartition(nums []int) bool {
	return false
}
