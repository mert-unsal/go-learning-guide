package sort_pkg
import "sort"
// SOLUTIONS — 02 sort
func SortByLengthSolution(words []string) []string {
result := make([]string, len(words))
copy(result, words)
sort.SliceStable(result, func(i, j int) bool {
return len(result[i]) < len(result[j])
})
return result
}
func abs(x int) int {
if x < 0 { return -x }
return x
}
func SortByAbsValueSolution(nums []int) []int {
result := make([]int, len(nums))
copy(result, nums)
sort.Slice(result, func(i, j int) bool {
return abs(result[i]) < abs(result[j])
})
return result
}
func RankSolution(nums []int) []int {
n := len(nums)
// Create index slice and sort it by value
indices := make([]int, n)
for i := range indices { indices[i] = i }
sort.Slice(indices, func(a, b int) bool {
return nums[indices[a]] < nums[indices[b]]
})
ranks := make([]int, n)
for rank, originalIdx := range indices {
ranks[originalIdx] = rank + 1 // 1-based rank
}
return ranks
}
func MedianSortedSolution(sorted []int) float64 {
n := len(sorted)
if n == 0 { return 0 }
if n%2 == 1 {
return float64(sorted[n/2])
}
return float64(sorted[n/2-1]+sorted[n/2]) / 2.0
}
func BinarySearchSolution(sorted []int, target int) int {
n := len(sorted)
i := sort.Search(n, func(i int) bool { return sorted[i] >= target })
if i < n && sorted[i] == target {
return i
}
return -1
}