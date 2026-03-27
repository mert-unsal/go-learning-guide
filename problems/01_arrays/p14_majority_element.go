package arrays

// MajorityElement ============================================================
// PROBLEM 14: Majority Element (LeetCode #169) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//
//	Given an array nums of size n, return the MAJORITY ELEMENT.
//
//	The majority element is the element that appears MORE THAN ⌊n/2⌋
//	times. You may assume that the majority element ALWAYS exists
//	in the array.
//
// PARAMETERS:
//
//	nums []int — an array of integers, guaranteed to have a majority element.
//
// RETURN:
//
//	int — the majority element.
//
// CONSTRAINTS:
//   - n == nums.length
//   - 1 <= n <= 5 × 10⁴
//   - -10⁹ <= nums[i] <= 10⁹
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Simple majority:
//
//	Input:  nums = [3, 2, 3]
//	Output: 3
//	Why:    3 appears 2 times out of 3 (> 3/2 = 1.5). ✓
//
// Example 2 — Larger array:
//
//	Input:  nums = [2, 2, 1, 1, 1, 2, 2]
//	Output: 2
//	Why:    2 appears 4 times out of 7 (> 7/2 = 3.5). ✓
//
// Example 3 — Single element:
//
//	Input:  nums = [1]
//	Output: 1
//
// Example 4 — All same:
//
//	Input:  nums = [5, 5, 5, 5, 5]
//	Output: 5
//
// Example 5 — Majority just barely:
//
//	Input:  nums = [1, 1, 2, 2, 1]
//	Output: 1
//	Why:    1 appears 3 times out of 5 (> 5/2 = 2.5). ✓
//
// MajorityElement returns the element that appears more than n/2 times.
// Time: O(n)  Space: O(1)
func MajorityElement(nums []int) int {
	// TODO: implement
	return 0
}
