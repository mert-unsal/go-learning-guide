package arrays

// ============================================================
// Diagonal Difference — [E]
// ============================================================
// Given a square matrix, compute the absolute difference
// between the sums of its diagonals.
//
// Example:
//   1 2 3       Primary:   1+5+9 = 15
//   4 5 6       Secondary: 3+5+7 = 15
//   7 8 9       |15 - 15| = 0

// DiagonalDifference returns the absolute diagonal difference.
// Time: O(n)  Space: O(1)
func DiagonalDifference(matrix [][]int) int {
	n := len(matrix)
	primary, secondary := 0, 0
	for i := 0; i < n; i++ {
		primary += matrix[i][i]
		secondary += matrix[i][n-1-i]
	}
	diff := primary - secondary
	if diff < 0 {
		return -diff
	}
	return diff
}
