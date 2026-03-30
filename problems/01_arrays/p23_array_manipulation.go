package arrays

// ============================================================
// Array Manipulation — [H]
// ============================================================
// Given an n-sized array of zeros and m operations, each operation adds value v
// to all elements between indices a and b (1-indexed inclusive).
// Find the maximum value after all operations.
//
// Example: n=5, queries=[[1,2,100],[2,5,100],[3,4,100]] → 200
//
// Key insight: difference array technique.
// Instead of updating all elements in range (O(n) per op),
// add v at index a subtract v at index b+1 (O(1) per op).
// Then compute prefix sum → original array values.

// ArrayManipulation returns the max value after all range-add operations.
// Time: O(n + m)  Space: O(n)
func ArrayManipulation(n int, queries [][]int) int64 {
	return 0
}
