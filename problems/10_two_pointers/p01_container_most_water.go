package two_pointers

// ============================================================
// PROBLEM 1: Container With Most Water (LeetCode #11) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given n non-negative integers height[0..n-1] where each represents
//   a vertical line at position i with height height[i], find two lines
//   that together with the x-axis form a container that holds the most
//   water. Return the maximum amount of water a container can store.
//   The container may not be slanted.
//
// PARAMETERS:
//   height []int — array of non-negative integers representing line heights
//
// RETURN:
//   int — maximum area of water the container can hold
//
// CONSTRAINTS:
//   • 2 ≤ len(height) ≤ 10⁵
//   • 0 ≤ height[i] ≤ 10⁴
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  height = [1, 8, 6, 2, 5, 4, 8, 3, 7]
//   Output: 49
//   Why:    Lines at index 1 (h=8) and index 8 (h=7), area = 7 × 7 = 49
//
// Example 2:
//   Input:  height = [1, 1]
//   Output: 1
//   Why:    Only two lines, area = min(1,1) × 1 = 1
//
// Example 3:
//   Input:  height = [4, 3, 2, 1, 4]
//   Output: 16
//   Why:    Lines at index 0 and 4, area = min(4,4) × 4 = 16
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Two pointers: start at both ends, move the shorter line inward
// • Moving the shorter line is the only way to potentially increase area
// • area = min(height[l], height[r]) × (r - l)
// • Target: O(n) time, O(1) space
func MaxArea(height []int) int {
	return 0
}
