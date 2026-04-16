package heap_priority_queue

import (
	"container/heap"
	"sort"
)

var _ = heap.Init
var _ = sort.Ints

// ============================================================
// PROBLEM 1: Kth Largest Element in an Array (LeetCode #215) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums and an integer k, return the
//   kth largest element in the array. Note that it is the kth
//   largest element in sorted order, not the kth distinct element.
//   Can you solve it without sorting?
//
// PARAMETERS:
//   nums []int — integer array
//   k    int   — the rank of the desired largest element (1-indexed)
//
// RETURN:
//   int — the kth largest element
//
// CONSTRAINTS:
//   • 1 <= k <= len(nums) <= 10^5
//   • -10^4 <= nums[i] <= 10^4
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [3,2,1,5,6,4], k = 2
//   Output: 5
//   Why:    Sorted descending: [6,5,4,3,2,1]; 2nd element is 5.
//
// Example 2:
//   Input:  nums = [3,2,3,1,2,4,5,5,6], k = 4
//   Output: 4
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Min-heap of size k: push all elements, pop when size > k; top is answer
// • Quickselect (Hoare's algorithm): average O(n), worst O(n^2)
// • Target: O(n log k) time with heap, O(n) average with quickselect, O(k) space
func FindKthLargest(nums []int, k int) int {
	return 0
}
