package stacks_queues

// ============================================================
// PROBLEM 6: Next Greater Element I (LeetCode #496) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   You are given two distinct 0-indexed integer arrays nums1 and
//   nums2, where nums1 is a subset of nums2. For each element in
//   nums1, find the next greater element in nums2. The next greater
//   element of nums1[i] is the first element to the right of
//   nums1[i]'s position in nums2 that is larger. Return -1 if none.
//
// PARAMETERS:
//   nums1 []int — subset array to query
//   nums2 []int — reference array to search in
//
// RETURN:
//   []int — for each element in nums1, the next greater element in nums2 or -1
//
// CONSTRAINTS:
//   • 1 <= len(nums1) <= len(nums2) <= 1000
//   • 0 <= nums1[i], nums2[i] <= 10^4
//   • All values in nums1 and nums2 are unique
//   • All values in nums1 also appear in nums2
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums1 = [4,1,2], nums2 = [1,3,4,2]
//   Output: [-1,3,-1]
//   Why:    4 has no greater to its right in nums2. 1→3 is next greater. 2 has none.
//
// Example 2:
//   Input:  nums1 = [2,4], nums2 = [1,2,3,4]
//   Output: [3,-1]
//   Why:    2→3 is the next greater. 4 has no greater element to its right.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Precompute the next greater element for every value in nums2
//   using a monotonic decreasing stack, storing results in a hash map.
// • Then look up each nums1 element in the map.
// • Target: O(n + m) time, O(n) space

func NextGreaterElement(nums1 []int, nums2 []int) []int {
	return nil
}
