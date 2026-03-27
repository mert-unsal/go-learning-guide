package arrays

import "strconv"

// ============================================================
// Almost Sorted — [M]
// ============================================================
// Determine if a permutation can be sorted by:
// - swapping exactly two elements (output "swap a b"), or
// - reversing exactly one contiguous subarray (output "reverse a b")
// Otherwise output "no".
// All indices are 1-based.

// AlmostSorted determines what single operation sorts the array.
func AlmostSorted(arr []int) string {
	n := len(arr)
	// Find the leftmost position that's out of order
	left := -1
	for i := 0; i < n-1; i++ {
		if arr[i] > arr[i+1] {
			left = i
			break
		}
	}
	if left == -1 {
		return "yes" // already sorted
	}
	// Find the rightmost position that's out of order
	right := -1
	for i := n - 1; i > 0; i-- {
		if arr[i] < arr[i-1] {
			right = i
			break
		}
	}
	// Try reversing arr[left..right] and check if fully sorted
	reversed := make([]int, n)
	copy(reversed, arr)
	for i, j := left, right; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	for i := 0; i < n-1; i++ {
		if reversed[i] > reversed[i+1] {
			return "no"
		}
	}
	// Sorted after reversal
	if left+1 == right {
		return "swap " + strconv.Itoa(left+1) + " " + strconv.Itoa(right+1)
	}
	return "reverse " + strconv.Itoa(left+1) + " " + strconv.Itoa(right+1)
}


