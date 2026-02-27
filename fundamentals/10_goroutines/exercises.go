package goroutines
import "sync"
// ============================================================
// EXERCISES — 10 Goroutines
// ============================================================
// Exercise 1:
// Run fn concurrently n times using goroutines + WaitGroup.
// Wait for all goroutines to finish before returning.
func RunConcurrently(n int, fn func(id int)) {
// TODO: launch n goroutines each calling fn(i), use WaitGroup
}
// Exercise 2:
// ExCounter is a thread-safe counter (uses a different name to avoid
// conflict with concepts.go). Implement Inc() and Value() using a Mutex.
type ExCounter struct {
mu    sync.Mutex
value int
}
func (c *ExCounter) Inc() {
// TODO: c.mu.Lock(); defer c.mu.Unlock(); c.value++
}
func (c *ExCounter) Value() int {
// TODO: c.mu.Lock(); defer c.mu.Unlock(); return c.value
return 0
}
// Exercise 3:
// SumConcurrent splits nums into two halves, sums each half
// in a separate goroutine, then returns the total.
func SumConcurrent(nums []int) int {
// TODO: WaitGroup + Mutex, two goroutines, combine results
return 0
}
// Exercise 4:
// RunOnce calls setup exactly once using sync.Once.
func RunOnce(setup func()) {
// TODO: var once sync.Once; once.Do(setup)
}