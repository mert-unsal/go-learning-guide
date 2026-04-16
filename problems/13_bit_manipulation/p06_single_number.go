package bit_manipulation

// ============================================================
// PROBLEM 6: Single Number (LeetCode #136) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a non-empty array of integers nums, every element
//   appears twice except for one. Find that single one. You must
//   implement a solution with linear runtime and constant extra space.
//
// PARAMETERS:
//   nums []int — array where every element appears twice except one
//
// RETURN:
//   int — the element that appears only once
//
// CONSTRAINTS:
//   • 1 <= len(nums) <= 3 * 10^4
//   • -3 * 10^4 <= nums[i] <= 3 * 10^4
//   • Each element appears exactly twice except one element which appears once
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [2,2,1]
//   Output: 1
//
// Example 2:
//   Input:  nums = [4,1,2,1,2]
//   Output: 4
//   Why:    XOR all elements: pairs cancel out, leaving 4.
//
// Example 3:
//   Input:  nums = [1]
//   Output: 1
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • XOR all elements: a ^ a = 0 and a ^ 0 = a, so pairs cancel
// • One pass through the array XOR-ing into an accumulator
// • Target: O(n) time, O(1) space
func SingleNumber(nums []int) int {
	return 0
}
