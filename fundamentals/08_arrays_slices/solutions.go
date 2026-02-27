package arrays_slices
// ============================================================
// SOLUTIONS — 08 Arrays & Slices
// ============================================================
// helper: reverse a slice between indices lo..hi inclusive
func reverse(s []int, lo, hi int) {
for lo < hi {
s[lo], s[hi] = s[hi], s[lo]
lo++
hi--
}
}
func ReverseSliceSolution(s []int) {
reverse(s, 0, len(s)-1)
}
func RemoveDuplicatesSolution(s []int) []int {
if len(s) == 0 {
return s
}
// write pointer: next position to write a unique value
write := 1
for i := 1; i < len(s); i++ {
if s[i] != s[i-1] { // new unique value
s[write] = s[i]
write++
}
}
return s[:write]
}
func Make2DSolution(rows, cols int) [][]int {
matrix := make([][]int, rows)
for i := range matrix {
matrix[i] = make([]int, cols) // each row is an independent slice
}
return matrix
}
func RotateLeftSolution(s []int, k int) {
n := len(s)
if n == 0 {
return
}
k = k % n // handle k >= n
// Three-reversal trick for left rotation by k:
// reverse(0, k-1), reverse(k, n-1), reverse(0, n-1)
reverse(s, 0, k-1)
reverse(s, k, n-1)
reverse(s, 0, n-1)
}
func FilterSolution(s []int, fn func(int) bool) []int {
result := make([]int, 0)
for _, v := range s {
if fn(v) {
result = append(result, v)
}
}
return result
}
func MergeSortedSolution(a, b []int) []int {
result := make([]int, 0, len(a)+len(b))
i, j := 0, 0
for i < len(a) && j < len(b) {
if a[i] <= b[j] {
result = append(result, a[i])
i++
} else {
result = append(result, b[j])
j++
}
}
// append remaining elements
result = append(result, a[i:]...)
result = append(result, b[j:]...)
return result
}