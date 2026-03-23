package arrays

// ============================================================
// PROBLEM 1: Two Sum (LeetCode #1) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array of integers nums and an integer target, return the
//   INDICES of the two numbers such that they add up to target.
//
//   You may assume that each input would have EXACTLY ONE solution,
//   and you may not use the same element twice.
//
//   You can return the answer in any order.
//
// PARAMETERS:
//   nums   []int — an unsorted array of integers (may contain negatives).
//   target int   — the sum we are looking for.
//
// RETURN:
//   []int — a slice of exactly two indices [i, j] where nums[i] + nums[j] == target.
//
// CONSTRAINTS:
//   • 2 <= nums.length <= 10⁴
//   • -10⁹ <= nums[i] <= 10⁹
//   • -10⁹ <= target <= 10⁹
//   • Only one valid answer exists.
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Basic:
//   Input:  nums = [2, 7, 11, 15], target = 9
//   Output: [0, 1]
//   Why:    nums[0] + nums[1] = 2 + 7 = 9.
//
// Example 2 — Not the first two elements:
//   Input:  nums = [3, 2, 4], target = 6
//   Output: [1, 2]
//   Why:    nums[1] + nums[2] = 2 + 4 = 6.
//
// Example 3 — Same value used twice (different indices):
//   Input:  nums = [3, 3], target = 6
//   Output: [0, 1]
//   Why:    Both elements are 3. They are at different indices, so this is valid.
//
// Example 4 — Negative numbers:
//   Input:  nums = [-3, 4, 3, 90], target = 0
//   Output: [0, 2]
//   Why:    -3 + 3 = 0.
//
// Example 5 — Large array, answer near the end:
//   Input:  nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], target = 19
//   Output: [8, 9]
//   Why:    9 + 10 = 19.
//
// Example 6 — Negative target:
//   Input:  nums = [-10, -20, 5, 15], target = -30
//   Output: [0, 1]
//   Why:    -10 + (-20) = -30.
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Brute force checks every pair — what is the time complexity?
//   • Can you do it in a single pass through the array?
//   • For each number, you need its "complement" (target - num).
//     How can you check if the complement was already seen?
//   • What data structure gives you O(1) lookups?
//   • Target: O(n) time, O(n) space.

// TwoSum returns the indices of two numbers that sum to target.
// Time: O(n)  Space: O(n)
func TwoSum(nums []int, target int) []int {
	// TODO: implement
	// value to store the indices of the two numbers
	var indices = make([]int, 2)
	var indicesMap = make(map[int]int)
	for i, num := range nums {
		complement := target - num
		if complementIndex, ok := indicesMap[complement]; ok {
			indices[0] = complementIndex
			indices[1] = i
		} else {
			indicesMap[num] = i
		}
	}
	return indices
}
