package heap_priority_queue

// ============================================================
// PROBLEM 5: Sliding Window Maximum (LeetCode #239) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   You are given an array of integers nums and a sliding window
//   of size k which moves from the very left to the very right.
//   You can only see the k numbers in the window. Each time the
//   window moves right by one position, return the max in each window.
//
// PARAMETERS:
//   nums []int — integer array
//   k    int   — size of the sliding window
//
// RETURN:
//   []int — array of maximum values for each window position
//
// CONSTRAINTS:
//   • 1 <= len(nums) <= 10^5
//   • -10^4 <= nums[i] <= 10^4
//   • 1 <= k <= len(nums)
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [1,3,-1,-3,5,3,6,7], k = 3
//   Output: [3,3,5,5,6,7]
//   Why:    Windows: [1,3,-1]→3, [3,-1,-3]→3, [-1,-3,5]→5, [-3,5,3]→5, [5,3,6]→6, [3,6,7]→7.
//
// Example 2:
//   Input:  nums = [1], k = 1
//   Output: [1]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Monotonic deque: maintain a decreasing deque of indices
// • Front of deque = index of current window max
// • Pop from back when new element >= deque back; pop from front when out of window
// • Target: O(n) time, O(k) space
func MaxSlidingWindow(nums []int, k int) []int {
	return nil
}
