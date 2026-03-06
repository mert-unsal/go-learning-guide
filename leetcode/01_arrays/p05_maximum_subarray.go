package arrays

// ============================================================
// PROBLEM 5: Maximum Subarray (LeetCode #53) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums, find the CONTIGUOUS subarray
//   (containing at least one number) which has the largest sum
//   and return its sum.
//
//   A subarray is a contiguous part of an array — elements must be
//   adjacent and in order. For example, [4, -1, 2] is a subarray
//   of [-2, 1, -3, 4, -1, 2, 1, -5, 4], but [4, 2, 1] is NOT
//   (not contiguous).
//
// PARAMETERS:
//   nums []int — an array of integers (may contain negatives).
//
// RETURN:
//   int — the sum of the contiguous subarray with the largest sum.
//
// CONSTRAINTS:
//   • 1 <= nums.length <= 10⁵
//   • -10⁴ <= nums[i] <= 10⁴
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Mixed positive and negative:
//   Input:  nums = [-2, 1, -3, 4, -1, 2, 1, -5, 4]
//   Output: 6
//   Why:    The subarray [4, -1, 2, 1] has the largest sum = 6.
//           Starting at index 3, ending at index 6.
//
// Example 2 — Single element:
//   Input:  nums = [1]
//   Output: 1
//   Why:    Only one element; the subarray is [1].
//
// Example 3 — All negative:
//   Input:  nums = [-3, -2, -5, -1]
//   Output: -1
//   Why:    The "least negative" single element is -1.
//           You must include at least one element.
//
// Example 4 — All positive:
//   Input:  nums = [1, 2, 3, 4]
//   Output: 10
//   Why:    The entire array is the maximum subarray.
//
// Example 5 — Negative sandwiched between positives:
//   Input:  nums = [5, -3, 5]
//   Output: 7
//   Why:    [5, -3, 5] = 7. The negative in the middle is worth absorbing.
//
// Example 6 — Large negative splits the array:
//   Input:  nums = [5, -100, 5]
//   Output: 5
//   Why:    The -100 makes [5, -100, 5] = -90. Better to take just [5].
//
// Example 7 — Zero elements:
//   Input:  nums = [0, -1, 0, -2, 0]
//   Output: 0
//   Why:    [0] or [0, -1, 0] → best is just [0] = 0.
//
// ─── KEY DECISION AT EACH POSITION ─────────────────────────
//
//   At each index i, you face a choice:
//     Option A: EXTEND the previous subarray by including nums[i]
//     Option B: START FRESH at nums[i] (abandon everything before)
//
//   You should start fresh when the accumulated sum so far is
//   negative — because adding a negative prefix only hurts.
//
//   This greedy choice at every step is called KADANE'S ALGORITHM.
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • What two variables do you need to track?
//   • When is it better to "start fresh" vs "extend"?
//   • Why does this work even when all numbers are negative?
//   • Target: O(n) time, O(1) space.

// MaxSubArray returns the largest sum of any contiguous subarray.
// Time: O(n)  Space: O(1)
func MaxSubArray(nums []int) int {
	// TODO: implement
	return 0
}
