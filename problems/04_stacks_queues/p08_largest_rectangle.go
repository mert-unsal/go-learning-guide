package stacks_queues

// ============================================================
// PROBLEM 8: Largest Rectangle in Histogram (LeetCode #84) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array of integers heights representing the histogram's
//   bar heights where the width of each bar is 1, return the area
//   of the largest rectangle that can be formed in the histogram.
//
// PARAMETERS:
//   heights []int — bar heights of the histogram (each bar has width 1)
//
// RETURN:
//   int — area of the largest rectangle in the histogram
//
// CONSTRAINTS:
//   • 1 <= len(heights) <= 10^5
//   • 0 <= heights[i] <= 10^4
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  heights = [2,1,5,6,2,3]
//   Output: 10
//   Why:    The largest rectangle spans indices 2-3 (heights 5,6) with
//           width 2 and height 5 → area = 10.
//
// Example 2:
//   Input:  heights = [2,4]
//   Output: 4
//   Why:    The single bar of height 4 gives area 4 (wider rectangle
//           of height 2, width 2 also gives 4).
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Use a monotonic increasing stack of indices. When a bar shorter
//   than the stack top is encountered, pop and compute the area using
//   the popped bar's height and the width between current index and
//   new stack top.
// • Add a sentinel height of 0 at the end to flush remaining bars.
// • Target: O(n) time, O(n) space

func LargestRectangleArea(heights []int) int {
	return 0
}
