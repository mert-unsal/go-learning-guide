package builtins
// ============================================================
// EXERCISES â€” 03 Builtins
// ============================================================
// Exercise 1:
// DeepCopySlice returns a completely independent copy of src.
// Modifying the copy must NOT affect src.
func DeepCopySlice(src []int) []int {
// TODO: make + copy
return nil
}
// Exercise 2:
// DeepCopyMap returns an independent copy of src map.
func DeepCopyMap(src map[string]int) map[string]int {
// TODO: make + range copy
return nil
}
// Exercise 3:
// SafeDivide uses recover to catch a panic from integer division by zero.
// Return the result and nil, or 0 and the recovered error message as an error.
func SafeDivideEx(a, b int) (result int, err error) {
// TODO: defer a recover() that converts the panic to an error
return a / b, nil
}
// Exercise 4:
// Flatten takes a [][]int (2D slice) and returns a single []int with all values.
// Example: [[1,2],[3],[4,5]] â†’ [1,2,3,4,5]
func Flatten(matrix [][]int) []int {
// TODO: append each row into result
return nil
}
// Exercise 5:
// UniqueInts removes duplicate integers from nums and returns
// a new slice preserving the FIRST occurrence order.
// Example: [3,1,4,1,5,9,2,6,5,3] â†’ [3,1,4,5,9,2,6]
func UniqueInts(nums []int) []int {
// TODO: use a map[int]bool to track seen, append unseen to result
return nil
}
// Exercise 6:
// ChunkSlice splits s into chunks of size n.
// The last chunk may be smaller than n.
// Example: [1,2,3,4,5], n=2 â†’ [[1,2],[3,4],[5]]
func ChunkSlice(s []int, n int) [][]int {
// TODO: iterate with step n, append s[i:min(i+n, len(s))]
return nil
}