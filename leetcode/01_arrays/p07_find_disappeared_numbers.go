package arrays

import "gointerviewprep/utils"

// ! ============================================================
// ! PROBLEM 7: Find All Numbers Disappeared in an Array
// ! LeetCode #448 — EASY
// ! ============================================================

// * PROBLEM STATEMENT:
//   Given an array nums of n integers where nums[i] is in the range [1, n],
//   return an array of all the integers in the range [1, n] that do NOT
//   appear in nums.
//
//   You must do it WITHOUT extra space (O(1) extra, not counting the output)
//   and in O(n) time.

// * PARAMETERS:
//   nums []int — an array of n integers, each in the range [1, n].
//                Some values may appear twice, causing others to be missing.
//
// * RETURN:
//   []int — a list of all missing numbers from the range [1, n].

// * CONSTRAINTS:
//   n == nums.length
//   1 <= n <= 10⁵
//   1 <= nums[i] <= n

// ! ─── EXAMPLES ───────────────────────────────────────────────

// * Example 1 — Two missing:
//   Input:  nums = [4, 3, 2, 7, 8, 2, 3, 1]     (n = 8)
//   Output: [5, 6]
//   Why:    The range [1..8] is missing 5 and 6.
//           2 and 3 each appear twice, "stealing" slots from 5 and 6.

// * Example 2 — One missing:
//   Input:  nums = [1, 1]     (n = 2)
//   Output: [2]
//   Why:    1 appears twice. 2 is missing.

// * Example 3 — None missing:
//   Input:  nums = [1, 2, 3, 4]     (n = 4)
//   Output: []
//   Why:    Every number in [1..4] is present exactly once.

// * Example 4 — All same:
//   Input:  nums = [1, 1, 1, 1]     (n = 4)
//   Output: [2, 3, 4]

// * Example 5 — Single element:
//   Input:  nums = [1]     (n = 1)
//   Output: []

// ! ─── KEY INSIGHT: USING THE ARRAY AS A HASH MAP ─────────────

//   Since values are in [1, n] and indices are in [0, n-1],
//   there is a natural mapping:  value v → index (v - 1).
//
//   If value v is present, "mark" index (v-1) by negating the value there.
//   Use absolute value when reading, because a slot may already be negated.
//   After marking, any index i still positive means (i+1) is missing.

// ? ─── THINGS TO THINK ABOUT ────────────────────────────────
// ?  1. The hash set approach works but uses O(n) space. Can you avoid that?
// ?  2. How can you use the ARRAY ITSELF as a marking structure?
// ?  3. Why do we use absolute values when reading during the mark pass?
// ?  4. After marking, how do you identify which numbers are missing?
// ?  5. Target: O(n) time, O(1) extra space.

// FindDisappearedNumbers returns missing numbers using O(1) extra space.
// Time: O(n)  Space: O(1) extra
func FindDisappearedNumbers(nums []int) []int {
	for _, v := range nums {
		idx := utils.Abs(v) - 1 // get the target index
		if nums[idx] > 0 {
			nums[idx] = -nums[idx] // mark as visited
		}
	}

	var result []int
	for i, v := range nums {
		if v > 0 {
			result = append(result, i+1)
		}
	}
	return result
}
