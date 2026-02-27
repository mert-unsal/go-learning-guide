package sort_pkg
import (
"reflect"
"testing"
)
func TestSortByLength(t *testing.T) {
got := SortByLengthSolution([]string{"banana", "fig", "apple", "kiwi"})
want := []string{"fig", "kiwi", "apple", "banana"}
if !reflect.DeepEqual(got, want) {
t.Errorf("SortByLength = %v, want %v", got, want)
}
}
func TestSortByAbsValue(t *testing.T) {
got := SortByAbsValueSolution([]int{-3, 1, -2, 4})
want := []int{1, -2, -3, 4}
if !reflect.DeepEqual(got, want) {
t.Errorf("SortByAbsValue = %v, want %v", got, want)
}
}
func TestRank(t *testing.T) {
got := RankSolution([]int{40, 10, 20, 30})
want := []int{4, 1, 2, 3}
if !reflect.DeepEqual(got, want) {
t.Errorf("Rank = %v, want %v", got, want)
}
}
func TestMedianSorted(t *testing.T) {
tests := []struct{ nums []int; want float64 }{
{[]int{1, 2, 3, 4, 5}, 3.0},
{[]int{1, 2, 3, 4}, 2.5},
{[]int{7}, 7.0},
}
for _, tt := range tests {
if got := MedianSortedSolution(tt.nums); got != tt.want {
t.Errorf("Median(%v) = %v, want %v", tt.nums, got, tt.want)
}
}
}
func TestBinarySearch(t *testing.T) {
s := []int{1, 3, 5, 7, 9, 11}
if got := BinarySearchSolution(s, 7); got != 3 {
t.Errorf("BinarySearch(7) = %d, want 3", got)
}
if got := BinarySearchSolution(s, 4); got != -1 {
t.Errorf("BinarySearch(4) = %d, want -1", got)
}
}