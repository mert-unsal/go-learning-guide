package arrays_slices
import (
"reflect"
"testing"
)
func TestReverseSlice(t *testing.T) {
tests := []struct {
input []int
want  []int
}{
{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
{[]int{1, 2}, []int{2, 1}},
{[]int{42}, []int{42}},
{[]int{}, []int{}},
}
for _, tt := range tests {
s := make([]int, len(tt.input))
copy(s, tt.input)
ReverseSliceSolution(s)
if !reflect.DeepEqual(s, tt.want) {
t.Errorf("ReverseSlice(%v) = %v, want %v", tt.input, s, tt.want)
}
}
}
func TestRemoveDuplicates(t *testing.T) {
tests := []struct {
input []int
want  []int
}{
{[]int{1, 1, 2, 3, 3, 4}, []int{1, 2, 3, 4}},
{[]int{1, 1, 1}, []int{1}},
{[]int{1, 2, 3}, []int{1, 2, 3}},
{[]int{}, []int{}},
}
for _, tt := range tests {
s := make([]int, len(tt.input))
copy(s, tt.input)
got := RemoveDuplicatesSolution(s)
if !reflect.DeepEqual(got, tt.want) {
t.Errorf("RemoveDuplicates(%v) = %v, want %v", tt.input, got, tt.want)
}
}
}
func TestMake2D(t *testing.T) {
m := Make2DSolution(3, 4)
if len(m) != 3 {
t.Fatalf("Make2D rows = %d, want 3", len(m))
}
for i, row := range m {
if len(row) != 4 {
t.Errorf("row[%d] len = %d, want 4", i, len(row))
}
}
// Verify rows are independent: modifying m[0] should not affect m[1]
m[0][0] = 99
if m[1][0] == 99 {
t.Error("rows share underlying array — they should be independent")
}
}
func TestRotateLeft(t *testing.T) {
tests := []struct {
input []int
k     int
want  []int
}{
{[]int{1, 2, 3, 4, 5}, 2, []int{3, 4, 5, 1, 2}},
{[]int{1, 2, 3, 4, 5}, 0, []int{1, 2, 3, 4, 5}},
{[]int{1, 2, 3, 4, 5}, 5, []int{1, 2, 3, 4, 5}},
{[]int{1, 2, 3}, 1, []int{2, 3, 1}},
}
for _, tt := range tests {
s := make([]int, len(tt.input))
copy(s, tt.input)
RotateLeftSolution(s, tt.k)
if !reflect.DeepEqual(s, tt.want) {
t.Errorf("RotateLeft(%v, %d) = %v, want %v", tt.input, tt.k, s, tt.want)
}
}
}
func TestFilter(t *testing.T) {
isEven := func(n int) bool { return n%2 == 0 }
got := FilterSolution([]int{1, 2, 3, 4, 5, 6}, isEven)
want := []int{2, 4, 6}
if !reflect.DeepEqual(got, want) {
t.Errorf("Filter(isEven) = %v, want %v", got, want)
}
isPositive := func(n int) bool { return n > 0 }
got2 := FilterSolution([]int{-2, -1, 0, 1, 2}, isPositive)
want2 := []int{1, 2}
if !reflect.DeepEqual(got2, want2) {
t.Errorf("Filter(isPositive) = %v, want %v", got2, want2)
}
}
func TestMergeSorted(t *testing.T) {
tests := []struct {
a, b []int
want []int
}{
{[]int{1, 3, 5}, []int{2, 4, 6}, []int{1, 2, 3, 4, 5, 6}},
{[]int{1, 2, 3}, []int{}, []int{1, 2, 3}},
{[]int{}, []int{4, 5}, []int{4, 5}},
{[]int{1}, []int{1}, []int{1, 1}},
}
for _, tt := range tests {
got := MergeSortedSolution(tt.a, tt.b)
if !reflect.DeepEqual(got, tt.want) {
t.Errorf("MergeSorted(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
}
}
}