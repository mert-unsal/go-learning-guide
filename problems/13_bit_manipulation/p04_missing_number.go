package bit_manipulation

// ============================================================
// PROBLEM 4: Missing Number (LeetCode #268) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array nums containing n distinct numbers in the
//   range [0, n], return the only number in the range that is
//   missing from the array.
//
// PARAMETERS:
//   nums []int — array of n distinct integers from [0, n]
//
// RETURN:
//   int — the missing number from the range [0, n]
//
// CONSTRAINTS:
//   • n == len(nums)
//   • 1 <= n <= 10^4
//   • 0 <= nums[i] <= n
//   • All values of nums are unique
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [3,0,1]
//   Output: 2
//   Why:    Range is [0,3]; 2 is missing.
//
// Example 2:
//   Input:  nums = [0,1]
//   Output: 2
//
// Example 3:
//   Input:  nums = [9,6,4,2,3,5,7,0,1]
//   Output: 8
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • XOR approach: XOR all numbers with all indices (0..n) — the missing one remains
// • Math approach: n*(n+1)/2 - sum(nums) = missing number
// • Target: O(n) time, O(1) space
func MissingNumber(nums []int) int {
	return 0
}
