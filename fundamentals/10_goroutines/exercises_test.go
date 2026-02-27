package goroutines
import (
"sync"
"sync/atomic"
"testing"
)
func TestRunConcurrently(t *testing.T) {
var counter int64
RunConcurrentlySolution(100, func(id int) {
atomic.AddInt64(&counter, 1)
})
if counter != 100 {
t.Errorf("RunConcurrently: counter = %d, want 100", counter)
}
}
func TestExCounter(t *testing.T) {
c := &ExCounter{}
var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
wg.Add(1)
go func() {
defer wg.Done()
c.IncSolution()
}()
}
wg.Wait()
if got := c.ValueSolution(); got != 1000 {
t.Errorf("ExCounter = %d, want 1000 (possible race condition)", got)
}
}
func TestSumConcurrent(t *testing.T) {
tests := []struct {
nums []int
want int
}{
{[]int{1, 2, 3, 4, 5}, 15},
{[]int{10, 20}, 30},
{[]int{}, 0},
{[]int{100}, 100},
}
for _, tt := range tests {
got := SumConcurrentSolution(tt.nums)
if got != tt.want {
t.Errorf("SumConcurrent(%v) = %d, want %d", tt.nums, got, tt.want)
}
}
}
func TestRunOnce(t *testing.T) {
calls := 0
RunOnceSolution(func() { calls++ })
if calls != 1 {
t.Errorf("setup ran %d times, want exactly 1", calls)
}
}