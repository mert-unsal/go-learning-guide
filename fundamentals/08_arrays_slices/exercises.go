package arrays_slices
// ============================================================
// EXERCISES — 08 Arrays & Slices
// ============================================================
// Exercise 1:
// Reverse a slice IN-PLACE (modify the original, return nothing).
// Example: [1,2,3,4,5] → [5,4,3,2,1]
func ReverseSlice(s []int) {
// TODO: use two pointers (left, right) swapping toward center
}
// Exercise 2:
// Remove duplicates from a SORTED slice and return the new slice.
// Must do it WITHOUT creating a new slice (in-place).
// Example: [1,1,2,3,3,4] → [1,2,3,4]
func RemoveDuplicates(s []int) []int {
// TODO: use a write pointer pattern
return nil
}
// Exercise 3:
// Return a 2D slice (matrix) of size rows×cols filled with zeros.
// Example: Make2D(2, 3) → [[0,0,0],[0,0,0]]
func Make2D(rows, cols int) [][]int {
// TODO: allocate rows, then allocate each row
return nil
}
// Exercise 4:
// Rotate a slice LEFT by k positions IN-PLACE.
// Example: [1,2,3,4,5], k=2 → [3,4,5,1,2]
// Hint: three-reversal trick — reverse all, reverse first n-k, reverse last k
func RotateLeft(s []int, k int) {
// TODO: implement three-reversal trick
}
// Exercise 5:
// Given a slice of integers, return a new slice containing only
// the elements that satisfy the predicate fn.
// Example: Filter([1,2,3,4,5], isEven) → [2,4]
func Filter(s []int, fn func(int) bool) []int {
// TODO: implement
return nil
}
// Exercise 6:
// Merge two SORTED slices into one sorted slice.
// Example: [1,3,5], [2,4,6] → [1,2,3,4,5,6]
func MergeSorted(a, b []int) []int {
// TODO: two-pointer merge
return nil
}