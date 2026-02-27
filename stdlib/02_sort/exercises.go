package sort_pkg
// ============================================================
// EXERCISES — 02 sort
// ============================================================
// Exercise 1:
// SortByLength sorts a slice of strings by length (shortest first).
// For equal lengths, preserve original relative order (stable sort).
// Example: ["banana","fig","apple","kiwi"] → ["fig","kiwi","apple","banana"]
func SortByLength(words []string) []string {
// TODO: make a copy, use sort.SliceStable with len comparison
return nil
}
// Exercise 2:
// SortByAbsValue sorts []int by absolute value, ascending.
// Example: [-3, 1, -2, 4] → [1, -2, -3, 4]
func SortByAbsValue(nums []int) []int {
// TODO: make a copy, sort.Slice with abs comparison
return nil
}
// Exercise 3:
// Rank returns the rank (1-based position in sorted order) of each element.
// Example: [40, 10, 20, 30] → [4, 1, 2, 3]
func Rank(nums []int) []int {
// TODO: sort a copy with indices, assign ranks back to original positions
return nil
}
// Exercise 4:
// MedianSorted returns the median of a sorted slice.
// For even length, return the average of the two middle elements (as float64).
// Example: [1,2,3,4,5] → 3.0   [1,2,3,4] → 2.5
func MedianSorted(sorted []int) float64 {
// TODO: handle odd/even cases
return 0
}
// Exercise 5:
// BinarySearch returns the index of target in a sorted slice, or -1.
// Implement it manually using sort.Search.
// sort.Search(n, f) returns the smallest index i in [0,n) where f(i) is true.
func BinarySearch(sorted []int, target int) int {
// TODO: use sort.Search, verify sorted[i] == target
return -1
}