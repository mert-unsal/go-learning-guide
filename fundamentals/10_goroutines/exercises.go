package goroutines

import (
	"context"
	"sync"
)

// ============================================================
// EXERCISES — 10 Goroutines
// ============================================================

// Exercise 1:
// Run fn concurrently n times using goroutines + WaitGroup.
// Wait for all goroutines to finish before returning.
func RunConcurrently(n int, fn func(id int)) {
	// TODO: use sync.WaitGroup to launch n goroutines and wait
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			fn(i)
		}()
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
	// TODO: lock, increment, unlock
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *ExCounter) Value() int {
	// TODO: lock, read, unlock, return
	c.mu.Lock()
	read := c.value
	c.mu.Unlock()
	return read
}

// Exercise 3:
// SumConcurrent splits nums into two halves, sums each half
// in a separate goroutine, then returns the total.
func SumConcurrent(nums []int) int {
	// TODO: split at midpoint, sum each half in a goroutine, combine
	var wg sync.WaitGroup
	var sumCh = make(chan int, 2)

	var subArrays = make([][]int, 2)
	subArrays[0] = nums[:len(nums)/2]
	subArrays[1] = nums[len(nums)/2:]
	wg.Add(len(subArrays))
	for _, subArray := range subArrays {
		go func(array []int) {
			defer wg.Done()
			sum := 0
			for _, v := range array {
				sum += v
			}
			sumCh <- sum
		}(subArray)
	}
	wg.Wait()
	close(sumCh)
	sum := 0
	for v := range sumCh {
		sum += v
	}

	return sum
}

// Exercise 4:
// RunOnce calls setup exactly once using sync.Once.
var runOnce sync.Once

func RunOnce(setup func()) {
	// TODO: use runOnce.Do(...)
	runOnce.Do(setup)
}

// ============================================================
// EXERCISES 5-7 — Concurrent Error Collection Patterns
// ============================================================
//
// Three progressively better patterns for collecting errors
// from multiple goroutines. Each solves the same problem
// differently — understand the tradeoffs.

// Exercise 5: WaitGroup + Mutex (Collect All Errors)
//
// PATTERN: Goroutines append errors to a shared slice protected
// by a mutex. WaitGroup coordinates completion.
//
// Requirements:
//   - Launch one goroutine per operation
//   - Each goroutine calls its operation func, which returns an error or nil
//   - If error is non-nil, append to a shared []error (protected by mutex)
//   - Wait for all goroutines to finish
//   - Return all collected errors (nil if none)
//   - Use defer wg.Done() for safety
//
// When to use: You need ALL errors, not just the first one.
// Tradeoff: Mutex contention if many goroutines error simultaneously.
func CollectAllErrors(operations []func() error) []error {
	// TODO: WaitGroup + Mutex pattern
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error
	for _, operationFn := range operations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := operationFn()
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return errors
}

// Exercise 6: Channel-Based Error Collection
//
// PATTERN: Goroutines send errors to a buffered channel.
// Main goroutine drains the channel after all work completes.
//
// Requirements:
//   - Create a buffered error channel (capacity = len(operations))
//   - Launch one goroutine per operation
//   - Each goroutine sends its error to the channel (even if nil)
//   - Use a "closer goroutine": go func() { wg.Wait(); close(errCh) }()
//   - Range over errCh to collect non-nil errors
//   - Return all non-nil errors
//
// When to use: Prefer over mutex pattern — channels are idiomatic Go.
// Tradeoff: Allocates buffered channel; all results pass through it.
func CollectErrorsViaChan(operations []func() error) []error {
	// TODO: buffered channel + closer goroutine pattern
	var errCh = make(chan error, len(operations))
	var wg sync.WaitGroup
	for _, operation := range operations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- operation()
		}()
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()

	var output []error
	for v := range errCh {
		if v != nil {
			output = append(output, v)
		}
	}
	return output
}

// Exercise 7: Context-Aware Error Collection (errgroup-style)
//
// PATTERN: Like exercise 6, but stops remaining work when the
// first error occurs — using context cancellation.
//
// Requirements:
//   - Create a cancellable context from the provided parent ctx
//   - Launch one goroutine per operation, passing ctx to each
//   - Operations accept context: func(ctx context.Context) error
//   - On FIRST error: cancel context (signal others to stop)
//   - Collect all errors that arrived before/during cancellation
//   - Return the collected errors
//
// When to use: Fail-fast — no point continuing if one task failed.
// This is what errgroup.Group does internally.
func CollectErrorsWithCancel(ctx context.Context, operations []func(ctx context.Context) error) []error {
	// TODO: context cancellation + channel collection
	var wg sync.WaitGroup
	var errCh = make(chan error, len(operations))
	cancellableCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, operation := range operations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := operation(cancellableCtx)
			if err != nil {
				cancel()
			}
			errCh <- err
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var output []error
	for v := range errCh {
		if v != nil {
			output = append(output, v)
		}
	}

	return output
}
