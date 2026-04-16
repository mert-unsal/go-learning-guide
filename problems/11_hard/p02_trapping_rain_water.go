package hard

// ============================================================
// PROBLEM 2: Trapping Rain Water (LeetCode #42) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   Given n non-negative integers representing an elevation map
//   where the width of each bar is 1, compute how much water it
//   can trap after raining.
//
// PARAMETERS:
//   height []int — elevation map where height[i] is the height of bar i
//
// RETURN:
//   int — total units of water trapped between the bars
//
// CONSTRAINTS:
//   • n == len(height)
//   • 1 <= n <= 2 * 10^4
//   • 0 <= height[i] <= 10^5
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  height = [0,1,0,2,1,0,1,3,2,1,2,1]
//   Output: 6
//   Why:    Water fills gaps between bars totaling 6 units.
//
// Example 2:
//   Input:  height = [4,2,0,3,2,5]
//   Output: 9
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Two-pointer approach: advance the pointer with the smaller max height
// • Water at position i = min(maxLeft, maxRight) - height[i]
// • Target: O(n) time, O(1) space with two pointers
func Trap(height []int) int {
	return 0
}
