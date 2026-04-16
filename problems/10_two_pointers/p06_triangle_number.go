package two_pointers

// ============================================================
// PROBLEM 6: Valid Triangle Number (LeetCode #611) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums, return the number of triplets chosen
//   from the array that can make triangles if we take them as side
//   lengths of a triangle. Three sides a ≤ b ≤ c form a valid triangle
//   if a + b > c.
//
// PARAMETERS:
//   nums []int — array of non-negative integers representing side lengths
//
// RETURN:
//   int — number of valid triangle triplets
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 1000
//   • 0 ≤ nums[i] ≤ 1000
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [2, 2, 3, 4]
//   Output: 3
//   Why:    Valid triplets: (2,3,4), (2,3,4), (2,2,3)
//
// Example 2:
//   Input:  nums = [4, 2, 3, 4]
//   Output: 4
//
// Example 3:
//   Input:  nums = [0, 1, 0]
//   Output: 0
//   Why:    0 + 0 is not greater than 1
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Sort the array first, then fix the largest side (c) and use two pointers
// • For fixed c = nums[k], find pairs (i, j) where nums[i] + nums[j] > nums[k]
// • If nums[i]+nums[j] > nums[k], all indices between i and j pair with j → count += j-i
// • Target: O(n²) time, O(1) extra space (or O(log n) for sort)
func TriangleNumber(nums []int) int {
	return 0
}
