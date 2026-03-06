package arrays

// ============================================================
// PROBLEM 12: Top K Frequent Elements (LeetCode #347) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums and an integer k, return the k MOST
//   FREQUENT elements. You may return the answer in any order.
//
//   It is guaranteed that the answer is unique (no ties for the kth
//   most frequent).
//
// PARAMETERS:
//   nums []int — an array of integers.
//   k    int   — number of most frequent elements to return.
//
// RETURN:
//   []int — the k most frequent elements (any order).
//
// CONSTRAINTS:
//   • 1 <= nums.length <= 10⁵
//   • -10⁴ <= nums[i] <= 10⁴
//   • k is in the range [1, number of unique elements]
//   • The answer is guaranteed to be unique.
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Basic:
//   Input:  nums = [1, 1, 1, 2, 2, 3], k = 2
//   Output: [1, 2]  (any order)
//   Why:    1 appears 3 times, 2 appears 2 times, 3 appears 1 time.
//           The top 2 frequent elements are 1 and 2.
//
// Example 2 — k equals number of unique elements:
//   Input:  nums = [1], k = 1
//   Output: [1]
//
// Example 3 — All same frequency except one:
//   Input:  nums = [4, 4, 4, 1, 2, 3], k = 1
//   Output: [4]
//   Why:    4 appears 3 times; the rest appear once each.
//
// Example 4 — Negative numbers:
//   Input:  nums = [-1, -1, 2, 2, 2, 3], k = 2
//   Output: [2, -1]  (any order)
//
// Example 5 — Large k:
//   Input:  nums = [1, 2, 3, 4, 5], k = 5
//   Output: [1, 2, 3, 4, 5]  (any order)
//   Why:    All appear once; all 5 are the top 5 most frequent.
//
// ─── APPROACHES TO CONSIDER ─────────────────────────────────
//
//   Approach 1 — Sort by frequency: O(n log n)
//     Count frequencies, sort by count, take top k.
//
//   Approach 2 — Min-heap of size k: O(n log k)
//     Maintain a heap of the k most frequent so far.
//
//   Approach 3 — Bucket sort: O(n)
//     Create buckets where index = frequency.
//     Bucket[freq] = list of numbers with that frequency.
//     Walk from highest frequency bucket down, collect k elements.
//     This works because max possible frequency = len(nums).
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • First step for any approach: count frequencies. What do you use?
//   • For bucket sort: what is the maximum possible frequency?
//   • Why is bucket sort O(n) while sorting is O(n log n)?
//   • Target: O(n) time (with bucket sort), O(n) space.

// TopKFrequent returns the k most frequent elements.
// Time: O(n)  Space: O(n)
func TopKFrequent(nums []int, k int) []int {
	// TODO: implement
	return nil
}
