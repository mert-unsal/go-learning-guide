package arrays

// SubarraySum ============================================================
// PROBLEM 10: Subarray Sum Equals K (LeetCode #560) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//
//	Given an array of integers nums and an integer k, return the total
//	number of CONTINUOUS subarrays whose sum equals k.
//
//	A subarray is a contiguous non-empty sequence of elements within
//	the array.
//
// PARAMETERS:
//
//	nums []int — an array of integers (may contain negatives and zeros).
//	k    int   — the target sum.
//
// RETURN:
//
//	int — the count of contiguous subarrays that sum to exactly k.
//
// CONSTRAINTS:
//   - 1 <= nums.length <= 2 × 10⁴
//   - -1000 <= nums[i] <= 1000
//   - -10⁷ <= k <= 10⁷
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Basic:
//
//	Input:  nums = [1, 1, 1], k = 2
//	Output: 2
//	Why:    Subarrays [1,1] (indices 0-1) and [1,1] (indices 1-2) both sum to 2.
//
// Example 2 — Negative numbers:
//
//	Input:  nums = [1, -1, 0], k = 0
//	Output: 3
//	Why:    [1,-1] = 0, [-1,0] = wait.. no. Let's check:
//	        [1,-1] = 0 ✓, [0] = 0 ✓, [1,-1,0] = 0 ✓ → 3.
//
// Example 3 — Single element equals k:
//
//	Input:  nums = [5], k = 5
//	Output: 1
//
// Example 4 — No subarray sums to k:
//
//	Input:  nums = [1, 2, 3], k = 7
//	Output: 0
//	Why:    Total sum is 6. No contiguous subarray sums to 7.
//
// Example 5 — Entire array sums to k:
//
//	Input:  nums = [1, 2, 3], k = 6
//	Output: 1
//	Why:    Only the full array [1,2,3] sums to 6.
//
// Example 6 — Multiple overlapping subarrays:
//
//	Input:  nums = [1, 2, 1, 2, 1], k = 3
//	Output: 4
//	Why:    [1,2] at indices 0-1, [2,1] at 1-2, [1,2] at 2-3, [2,1] at 3-4.
//
// SubarraySum returns the number of subarrays with sum equal to k.
// Time: O(n)  Space: O(n)
func SubarraySum(nums []int, k int) int {
	return 0
}
