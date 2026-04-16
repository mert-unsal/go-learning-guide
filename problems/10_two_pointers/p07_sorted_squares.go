package two_pointers

// ============================================================
// PROBLEM 7: Squares of a Sorted Array (LeetCode #977) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums sorted in non-decreasing order, return
//   an array of the squares of each number sorted in non-decreasing
//   order.
//
// PARAMETERS:
//   nums []int — sorted array of integers (may include negatives)
//
// RETURN:
//   []int — array of squared values, sorted in non-decreasing order
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 10⁴
//   • -10⁴ ≤ nums[i] ≤ 10⁴
//   • nums is sorted in non-decreasing order
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [-4, -1, 0, 3, 10]
//   Output: [0, 1, 9, 16, 100]
//
// Example 2:
//   Input:  nums = [-7, -3, 2, 3, 11]
//   Output: [4, 9, 9, 49, 121]
//
// Example 3:
//   Input:  nums = [1, 2, 3]
//   Output: [1, 4, 9]
//   Why:    All positive — squares are already sorted
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Two pointers at both ends: largest square is at either end (not middle)
// • Compare abs(nums[left]) vs abs(nums[right]), place the larger square at the back
// • Fill result array from right to left
// • Target: O(n) time, O(n) space (for the result array)
func SortedSquares(nums []int) []int {
	return nil
}
