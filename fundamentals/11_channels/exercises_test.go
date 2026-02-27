package channels
import (
"sort"
"testing"
)
func TestSumAsync(t *testing.T) {
ch := make(chan int, 1)
go SumAsyncSolution([]int{1, 2, 3, 4, 5}, ch)
got := <-ch
if got != 15 {
t.Errorf("SumAsync = %d, want 15", got)
}
}
func TestGenerate(t *testing.T) {
ch := GenerateSolution(5)
result := []int{}
for v := range ch {
result = append(result, v)
}
want := []int{1, 2, 3, 4, 5}
for i, v := range want {
if result[i] != v {
t.Errorf("Generate[%d] = %d, want %d", i, result[i], v)
}
}
}
func TestSquare(t *testing.T) {
in := GenerateSolution(4) // 1,2,3,4
out := SquareSolution(in) // 1,4,9,16
want := []int{1, 4, 9, 16}
for _, w := range want {
got := <-out
if got != w {
t.Errorf("Square: got %d, want %d", got, w)
}
}
}
func TestMerge(t *testing.T) {
make123 := func() <-chan int {
ch := make(chan int, 3)
ch <- 1; ch <- 2; ch <- 3
close(ch)
return ch
}
a := make123()
b := make123()
merged := MergeSolution(a, b)
result := []int{}
for v := range merged {
result = append(result, v)
}
sort.Ints(result)
want := []int{1, 1, 2, 2, 3, 3}
for i, v := range want {
if result[i] != v {
t.Errorf("Merge[%d] = %d, want %d", i, result[i], v)
}
}
}
func TestWithTimeout(t *testing.T) {
// Channel that sends immediately
fast := make(chan int, 1)
fast <- 42
v, ok := WithTimeoutSolution(fast, 100)
if !ok || v != 42 {
t.Errorf("WithTimeout fast: got (%d,%v), want (42,true)", v, ok)
}
// Channel that never sends — should timeout
slow := make(chan int)
v, ok = WithTimeoutSolution(slow, 10) // 10ms timeout
if ok {
t.Errorf("WithTimeout slow: expected timeout, got value %d", v)
}
}