package two_pointers

// ============================================================
// PROBLEM 4: Move Zeroes (LeetCode #283) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums, move all 0s to the end of it while
//   maintaining the relative order of the non-zero elements.
//   You must do this in-place without making a copy of the array.
//
// PARAMETERS:
//   nums []int — array of integers (modified in-place)
//
// RETURN:
//   (none) — the array is modified in-place
//
// CONSTRAINTS:
//   • 1 ≤ len(nums) ≤ 10⁴
//   • -2³¹ ≤ nums[i] ≤ 2³¹ - 1
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  nums = [0, 1, 0, 3, 12]
//   Output: [1, 3, 12, 0, 0]
//
// Example 2:
//   Input:  nums = [0]
//   Output: [0]
//
// Example 3:
//   Input:  nums = [1, 2, 3]
//   Output: [1, 2, 3]
//   Why:    No zeroes — array unchanged
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Two pointers: slow pointer tracks the insert position for non-zero elements
// • Fast pointer scans every element; when non-zero, swap with slow position
// • After one pass, all zeroes are naturally at the end
// • Minimize writes: only swap when slow ≠ fast
// • Target: O(n) time, O(1) space
func MoveZeroes(nums []int) {
}
