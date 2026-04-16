package dynamic_prog

// ============================================================
// PROBLEM 11: Maximum Product Subarray (LeetCode #152) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums, find a contiguous subarray that has
//   the largest product, and return that product.
//
// PARAMETERS:
//   nums []int — array of integers (may contain negatives and zeros)
//
// RETURN:
//   int — the maximum product of any contiguous subarray
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 2 × 10⁴
//   • -10 ≤ nums[i] ≤ 10
//   • The product of any subarray fits in a 32-bit integer
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [2, 3, -2, 4]
//   Output: 6
//   Why:    [2, 3] has the largest product
//
// Example 2:
//   Input:  nums = [-2, 0, -1]
//   Output: 0
//   Why:    Subarray [0] gives the max product
//
// Example 3:
//   Input:  nums = [-2, 3, -4]
//   Output: 24
//   Why:    [-2, 3, -4] — two negatives make a positive: (-2)×3×(-4) = 24
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Track both curMax and curMin at each position (negatives flip sign)
// • curMax = max(nums[i], curMax*nums[i], curMin*nums[i])
// • curMin = min(nums[i], curMax*nums[i], curMin*nums[i]) — compute before updating curMax
// • Update global max at each step
// • Target: O(n) time, O(1) space
func MaxProduct(nums []int) int {
	return 0
}
