package goroutines

import "sync"

// ============================================================
// EXERCISES — 10 Goroutines
// ============================================================
// Exercise 1:
// Run fn concurrently n times using goroutines + WaitGroup.
// Wait for all goroutines to finish before returning.
func RunConcurrently(n int, fn func(id int)) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fn(id)
		}(i)
	}
	wg.Wait()
}

// Exercise 2:
// ExCounter is a thread-safe counter (uses a different name to avoid
// conflict with concepts.go). Implement Inc() and Value() using a Mutex.
type ExCounter struct {
	mu    sync.Mutex
	value int
}

func (c *ExCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *ExCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Exercise 3:
// SumConcurrent splits nums into two halves, sums each half
// in a separate goroutine, then returns the total.
func SumConcurrent(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	mid := len(nums) / 2
	var sum1, sum2 int
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for _, v := range nums[:mid] {
			sum1 += v
		}
	}()
	go func() {
		defer wg.Done()
		for _, v := range nums[mid:] {
			sum2 += v
		}
	}()
	wg.Wait()
	return sum1 + sum2
}

var runOnce sync.Once

// Exercise 4:
// RunOnce calls setup exactly once using sync.Once.
func RunOnce(setup func()) {
	runOnce.Do(setup)
}
