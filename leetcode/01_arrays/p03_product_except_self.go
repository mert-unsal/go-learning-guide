package arrays

// ============================================================
// PROBLEM 3: Product of Array Except Self (LeetCode #238) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an integer array nums, return an array answer such that
//   answer[i] is equal to the product of ALL the elements of nums
//   EXCEPT nums[i].
//
//   The product of any prefix or suffix of nums is GUARANTEED to fit
//   in a 32-bit integer.
//
//   You must write an algorithm that runs in O(n) time and WITHOUT
//   using the division operator.
//
// PARAMETERS:
//   nums []int — an array of integers (may include zeros and negatives).
//
// RETURN:
//   []int — an array where output[i] = product of all elements except nums[i].
//
// CONSTRAINTS:
//   • 2 <= nums.length <= 10⁵
//   • -30 <= nums[i] <= 30
//   • The product of any prefix or suffix fits in a 32-bit integer.
//
// WHY NO DIVISION?
//   The naive approach would be: compute the total product, then divide by
//   each element. But the problem explicitly forbids division. Also,
//   division breaks when any element is zero (division by zero).
//   You must find another way.
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Basic:
//   Input:  nums = [1, 2, 3, 4]
//   Output: [24, 12, 8, 6]
//   Why:    output[0] = 2*3*4 = 24
//           output[1] = 1*3*4 = 12
//           output[2] = 1*2*4 = 8
//           output[3] = 1*2*3 = 6
//
//
// Example 2 — Contains a zero:
//   Input:  nums = [0, 1, 2, 3]
//   Output: [6, 0, 0, 0]
//   Why:    output[0] = 1*2*3 = 6 (the zero is excluded)
//           For all other positions, the product includes the zero → 0.
//
// Example 3 — Two zeros:
//   Input:  nums = [0, 0, 2, 3]
//   Output: [0, 0, 0, 0]
//   Why:    Every product includes at least one zero.
//
// Example 4 — Negative numbers:
//   Input:  nums = [-1, 1, 0, -3, 3]
//   Output: [0, 0, 9, 0, 0]
//
// Example 5 — Two elements:
//   Input:  nums = [4, 5]
//   Output: [5, 4]
//   Why:    output[0] = 5, output[1] = 4.
//
// Example 6 — All ones:
//   Input:  nums = [1, 1, 1, 1]
//   Output: [1, 1, 1, 1]
//
// ─── VISUALIZATION ─────────────────────────────────────────
//
//   For nums = [a, b, c, d]:
//
//   output[i] = (everything LEFT of i) × (everything RIGHT of i)
//
//   Index:     0        1        2        3
//   Left:      1        a       a*b     a*b*c    ← "prefix products"
//   Right:   b*c*d     c*d       d        1      ← "suffix products"
//   Answer:  b*c*d    a*c*d    a*b*d    a*b*c
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Can you split each output[i] into a LEFT part and a RIGHT part?
//   • Can you compute all left prefix products in one pass?
//   • Can you then compute all right suffix products in a second pass?
//   • Can you do the second pass using just a single running variable
//     instead of a whole array?
//   • Target: O(n) time, O(1) extra space (output array doesn't count).

// ProductExceptSelf returns the product array without division.
// Time: O(n)  Space: O(1) extra (output array doesn't count)
func ProductExceptSelf(nums []int) []int {
	// TODO: implement
	output := make([]int, len(nums))
	output[0] = 1
	for i := 1; i < len(nums); i++ {
		output[i] = output[i-1] * nums[i-1]
	}
	right := 1
	for i := len(nums) - 1; i >= 0; i-- {
		output[i] *= right
		right *= nums[i]
	}
	return output
}
