package arrays

// LongestConsecutive ============================================================
// PROBLEM 11: Longest Consecutive Sequence (LeetCode #128) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//
//	Given an unsorted array of integers nums, return the length of the
//	longest CONSECUTIVE elements sequence.
//
//	Consecutive means values differ by exactly 1: e.g., [1, 2, 3, 4].
//	The sequence does NOT need to be contiguous in the array — elements
//	can be scattered anywhere.
//
//	You must write an algorithm that runs in O(n) time.
//
// PARAMETERS:
//
//	nums []int — an unsorted array of integers.
//
// RETURN:
//
//	int — the length of the longest consecutive sequence.
//
// CONSTRAINTS:
//   - 0 <= nums.length <= 10⁵
//   - -10⁹ <= nums[i] <= 10⁹
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Basic:
//
//	Input:  nums = [100, 4, 200, 1, 3, 2]
//	Output: 4
//	Why:    The longest consecutive sequence is [1, 2, 3, 4]. Length = 4.
//	        100 and 200 are isolated — their sequences are length 1.
//
// Example 2 — Already sorted:
//
//	Input:  nums = [0, 3, 7, 2, 5, 8, 4, 6, 0, 1]
//	Output: 9
//	Why:    The sequence [0, 1, 2, 3, 4, 5, 6, 7, 8] has length 9.
//
// Example 3 — Empty array:
//
//	Input:  nums = []
//	Output: 0
//
// Example 4 — Single element:
//
//	Input:  nums = [42]
//	Output: 1
//
// Example 5 — Duplicates:
//
//	Input:  nums = [1, 2, 0, 1]
//	Output: 3
//	Why:    Sequence [0, 1, 2]. Duplicates don't extend the sequence.
//
// Example 6 — Negative numbers:
//
//	Input:  nums = [-5, -4, -3, 10, 11]
//	Output: 3
//	Why:    [-5, -4, -3] has length 3. [10, 11] has length 2.
//
// Example 7 — All same:
//
//	Input:  nums = [5, 5, 5, 5]
//	Output: 1
//	Why:    Only one unique value → sequence length 1.
//
// LongestConsecutive returns the length of the longest consecutive sequence.
// Time: O(n)  Space: O(n)
func LongestConsecutive(nums []int) int {
	// TODO: implement
	return 0
}
