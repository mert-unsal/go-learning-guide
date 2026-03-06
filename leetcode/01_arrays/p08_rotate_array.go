package arrays

// ============================================================
// PROBLEM 8: Rotate Array (LeetCode #189) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums, rotate the array to the RIGHT by k steps,
//   where k is non-negative. Do this IN-PLACE with O(1) extra space.
//
//   Rotating to the right by 1 step means the last element moves to
//   the front, and every other element shifts one position to the right.
//
// PARAMETERS:
//   nums []int — the array to rotate (modified in-place).
//   k    int   — number of positions to rotate right.
//
// RETURN:
//   (none — modify nums in-place)
//
// CONSTRAINTS:
//   • 1 <= nums.length <= 10⁵
//   • -2³¹ <= nums[i] <= 2³¹ - 1
//   • 0 <= k <= 10⁵
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Rotate right by 3:
//   Input:  nums = [1, 2, 3, 4, 5, 6, 7], k = 3
//   Output: [5, 6, 7, 1, 2, 3, 4]
//   Step-by-step:
//     After 1 rotation: [7, 1, 2, 3, 4, 5, 6]
//     After 2 rotations: [6, 7, 1, 2, 3, 4, 5]
//     After 3 rotations: [5, 6, 7, 1, 2, 3, 4]
//
// Example 2 — Rotate right by 2:
//   Input:  nums = [-1, -100, 3, 99], k = 2
//   Output: [3, 99, -1, -100]
//
// Example 3 — k equals array length (full rotation = no change):
//   Input:  nums = [1, 2, 3], k = 3
//   Output: [1, 2, 3]
//
// Example 4 — k larger than array length:
//   Input:  nums = [1, 2, 3], k = 5
//   Output: [2, 3, 1]
//   Why:    k=5 with len=3 is equivalent to k = 5 % 3 = 2.
//
// Example 5 — Single element:
//   Input:  nums = [42], k = 7
//   Output: [42]
//   Why:    Any rotation of a single element gives the same array.
//
// Example 6 — k = 0:
//   Input:  nums = [1, 2, 3], k = 0
//   Output: [1, 2, 3]
//   Why:    No rotation needed.
//
// ─── THE REVERSE TRICK ─────────────────────────────────────
//
//   There's an elegant O(n) time, O(1) space approach using reversals:
//
//   Given nums = [1, 2, 3, 4, 5, 6, 7], k = 3:
//
//   Step 1: Reverse the entire array:
//           [7, 6, 5, 4, 3, 2, 1]
//
//   Step 2: Reverse the first k elements [0..k-1]:
//           [5, 6, 7, 4, 3, 2, 1]
//
//   Step 3: Reverse the remaining elements [k..n-1]:
//           [5, 6, 7, 1, 2, 3, 4]  ✓
//
//   Why does this work? Reversing twice restores order, but the two
//   halves are now in swapped positions.
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • What if k >= len(nums)? What does k % len(nums) give you?
//   • The naive approach (shift one by one, k times) is O(n*k). Too slow.
//   • Using a temp array is O(n) space. Can you do O(1)?
//   • You may want a helper: func reverse(nums []int, left, right int)
//   • Target: O(n) time, O(1) space.

// Rotate rotates the array to the right by k positions in-place.
// Time: O(n)  Space: O(1)
func Rotate(nums []int, k int) {
	// TODO: implement
}
