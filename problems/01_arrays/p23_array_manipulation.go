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
	diff := make([]int64, n+2) // 1-indexed, extra slot for n+1

	for _, q := range queries {
		a, b, v := q[0], q[1], int64(q[2])
		diff[a] += v
		if b+1 <= n {
			diff[b+1] -= v
		}
	}

	var maxVal, running int64
	for i := 1; i <= n; i++ {
		running += diff[i]
		if running > maxVal {
			maxVal = running
		}
	}
	return maxVal
}
