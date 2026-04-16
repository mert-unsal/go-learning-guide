package dynamic_prog

// ============================================================
// PROBLEM 13: House Robber II (LeetCode #213) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   All houses are arranged in a circle. Adjacent houses have connected
//   security systems — you cannot rob two adjacent houses. Given an
//   array representing the amount of money at each house, return the
//   maximum amount you can rob without alerting the police. Since the
//   first and last house are also adjacent, you cannot rob both.
//
// PARAMETERS:
//   nums []int — amount of money at each house (arranged in a circle)
//
// RETURN:
//   int — maximum amount you can rob without robbing two adjacent houses
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 100
//   • 0 ≤ nums[i] ≤ 1000
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [2, 3, 2]
//   Output: 3
//   Why:    Rob house 1 only (houses 0 and 2 are adjacent in a circle)
//
// Example 2:
//   Input:  nums = [1, 2, 3, 1]
//   Output: 4
//   Why:    Rob house 0 (1) + house 2 (3) = 4
//
// Example 3:
//   Input:  nums = [1, 2, 3]
//   Output: 3
//   Why:    Rob house 2 only
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Break the circle: run House Robber I on nums[0..n-2] and nums[1..n-1]
// • Answer is max of both runs — this avoids robbing both first and last
// • Reuse the linear House Robber logic as a helper (robRange)
// • Target: O(n) time, O(1) space
func RobII(nums []int) int {
	return 0
}
func robRange(nums []int, start, end int) int {
	return 0
}
