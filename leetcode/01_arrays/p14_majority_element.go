package arrays

// ============================================================
// PROBLEM 14: Majority Element (LeetCode #169) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array nums of size n, return the MAJORITY ELEMENT.
//
//   The majority element is the element that appears MORE THAN ⌊n/2⌋
//   times. You may assume that the majority element ALWAYS exists
//   in the array.
//
// PARAMETERS:
//   nums []int — an array of integers, guaranteed to have a majority element.
//
// RETURN:
//   int — the majority element.
//
// CONSTRAINTS:
//   • n == nums.length
//   • 1 <= n <= 5 × 10⁴
//   • -10⁹ <= nums[i] <= 10⁹
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Simple majority:
//   Input:  nums = [3, 2, 3]
//   Output: 3
//   Why:    3 appears 2 times out of 3 (> 3/2 = 1.5). ✓
//
// Example 2 — Larger array:
//   Input:  nums = [2, 2, 1, 1, 1, 2, 2]
//   Output: 2
//   Why:    2 appears 4 times out of 7 (> 7/2 = 3.5). ✓
//
// Example 3 — Single element:
//   Input:  nums = [1]
//   Output: 1
//
// Example 4 — All same:
//   Input:  nums = [5, 5, 5, 5, 5]
//   Output: 5
//
// Example 5 — Majority just barely:
//   Input:  nums = [1, 1, 2, 2, 1]
//   Output: 1
//   Why:    1 appears 3 times out of 5 (> 5/2 = 2.5). ✓
//
// ─── APPROACHES ─────────────────────────────────────────────
//
//   Approach 1 — Hash map: count frequencies, return the one > n/2.
//     Time: O(n), Space: O(n).
//
//   Approach 2 — Sort: the majority element will always be at index n/2.
//     Time: O(n log n), Space: O(1).
//
//   Approach 3 — Boyer-Moore Voting Algorithm: O(n) time, O(1) space.
//     Maintain a "candidate" and a "count".
//     - If count == 0, pick current element as new candidate.
//     - If current element == candidate, increment count.
//     - Otherwise, decrement count.
//     The majority element always survives because it appears more than
//     all others combined — their "votes" cancel each other out.
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Why does Boyer-Moore work? Think of it as a "voting" game.
//   • What happens when two different elements cancel each other?
//   • The majority element has MORE than n/2 votes — it can never
//     be fully cancelled.
//   • Target: O(n) time, O(1) space.

// MajorityElement returns the element that appears more than n/2 times.
// Time: O(n)  Space: O(1)
func MajorityElement(nums []int) int {
	// TODO: implement
	return 0
}
