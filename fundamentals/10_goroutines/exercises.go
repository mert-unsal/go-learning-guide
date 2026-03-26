package goroutines

import "sync"

// ============================================================
// EXERCISES — 10 Goroutines
// ============================================================

// Exercise 1:
// Run fn concurrently n times using goroutines + WaitGroup.
// Wait for all goroutines to finish before returning.
func RunConcurrently(n int, fn func(id int)) {
	// TODO: use sync.WaitGroup to launch n goroutines and wait
	panic("not implemented")
}

// Exercise 2:
// ExCounter is a thread-safe counter (uses a different name to avoid
// conflict with concepts.go). Implement Inc() and Value() using a Mutex.
type ExCounter struct {
	mu    sync.Mutex
	value int
}

func (c *ExCounter) Inc() {
	// TODO: lock, increment, unlock
	panic("not implemented")
}

func (c *ExCounter) Value() int {
	// TODO: lock, read, unlock, return
	panic("not implemented")
}

// Exercise 3:
// SumConcurrent splits nums into two halves, sums each half
// in a separate goroutine, then returns the total.
func SumConcurrent(nums []int) int {
	// TODO: split at midpoint, sum each half in a goroutine, combine
	panic("not implemented")
}

// Exercise 4:
// RunOnce calls setup exactly once using sync.Once.
var runOnce sync.Once

func RunOnce(setup func()) {
	// TODO: use runOnce.Do(...)
	panic("not implemented")
}
