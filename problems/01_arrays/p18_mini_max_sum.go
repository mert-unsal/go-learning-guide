package arrays

import "sort"

// ============================================================
// Mini-Max Sum — [E]
// ============================================================
// Given 5 integers, find the minimum and maximum sums of 4 out of 5 numbers.
//
// Example: arr=[1,2,3,4,5]
// minSum = 1+2+3+4 = 10
// maxSum = 2+3+4+5 = 14
// Output: "10 14"
//
// Key insight: sort → minSum = sum of first 4, maxSum = sum of last 4.
// Or: total - max = minSum, total - min = maxSum

// MiniMaxSum returns the minimum and maximum 4-element sums.
// Time: O(n log n)  Space: O(1)
func MiniMaxSum(arr []int) (minSum, maxSum int) {
	sort.Ints(arr)
	for i := 0; i < 4; i++ {
		minSum += arr[i]
	}
	for i := 1; i < 5; i++ {
		maxSum += arr[i]
	}
	return
}
