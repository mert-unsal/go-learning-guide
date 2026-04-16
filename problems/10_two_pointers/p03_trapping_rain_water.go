package two_pointers

// ============================================================
// PROBLEM 3: Trapping Rain Water (LeetCode #42) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   Given n non-negative integers representing an elevation map where
//   the width of each bar is 1, compute how much water it can trap
//   after raining.
//
// PARAMETERS:
//   height []int — array of non-negative integers representing bar heights
//
// RETURN:
//   int — total units of trapped water
//
// CONSTRAINTS:
//   • 1 ≤ len(height) ≤ 2 × 10⁴
//   • 0 ≤ height[i] ≤ 10⁵
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  height = [0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1]
//   Output: 6
//   Why:    Water fills the valleys between bars
//
// Example 2:
//   Input:  height = [4, 2, 0, 3, 2, 5]
//   Output: 9
//
// Example 3:
//   Input:  height = [4, 2, 3]
//   Output: 1
//   Why:    1 unit trapped between index 0 (h=4) and index 2 (h=3) above index 1 (h=2)
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Water at position i = min(maxLeft, maxRight) - height[i]
// • Two-pointer approach: track leftMax and rightMax, advance the smaller side
// • If leftMax < rightMax, water at left is bounded by leftMax — safe to compute
// • Also solvable with prefix/suffix max arrays or a monotonic stack
// • Target: O(n) time, O(1) space (two-pointer approach)
func Trap(height []int) int {
	return 0
}
