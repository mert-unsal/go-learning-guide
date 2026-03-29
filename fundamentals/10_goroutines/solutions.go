package goroutines

import (
	"context"
	"sync"
)

// ============================================================
// SOLUTIONS — 10 Goroutines
// ============================================================
func RunConcurrentlySolution(n int, fn func(id int)) {
var wg sync.WaitGroup
for i := 0; i < n; i++ {
wg.Add(1)
i := i // capture loop variable — critical in Go < 1.22
go func() {
defer wg.Done()
fn(i)
}()
}
wg.Wait()
}
func (c *ExCounter) IncSolution() {
c.mu.Lock()
defer c.mu.Unlock()
c.value++
}
func (c *ExCounter) ValueSolution() int {
c.mu.Lock()
defer c.mu.Unlock()
return c.value
}
func SumConcurrentSolution(nums []int) int {
if len(nums) == 0 {
return 0
}
mid := len(nums) / 2
var mu sync.Mutex
total := 0
var wg sync.WaitGroup
sumHalf := func(slice []int) {
defer wg.Done()
s := 0
for _, v := range slice {
s += v
}
mu.Lock()
total += s
mu.Unlock()
}
wg.Add(2)
go sumHalf(nums[:mid])
go sumHalf(nums[mid:])
wg.Wait()
return total
}
func RunOnceSolution(setup func()) {
var once sync.Once
once.Do(setup) // runs setup
once.Do(setup) // does nothing — already ran
once.Do(setup) // does nothing — already ran
}

// ── Exercise 5 Solution: WaitGroup + Mutex ──

func CollectAllErrorsSolution(operations []func() error) []error {
	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	for _, op := range operations {
		wg.Add(1)
		go func(fn func() error) {
			defer wg.Done()
			if err := fn(); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(op)
	}

	wg.Wait()
	return errs
}

// ── Exercise 6 Solution: Channel-Based ──

func CollectErrorsViaChanSolution(operations []func() error) []error {
	errCh := make(chan error, len(operations))
	var wg sync.WaitGroup

	for _, op := range operations {
		wg.Add(1)
		go func(fn func() error) {
			defer wg.Done()
			errCh <- fn() // send even if nil — channel drains all
		}(op)
	}

	// Closer goroutine: waits for all workers, then closes channel
	go func() {
		wg.Wait()
		close(errCh)
	}()

	var errs []error
	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// ── Exercise 7 Solution: Context-Aware (errgroup-style) ──

func CollectErrorsWithCancelSolution(ctx context.Context, operations []func(ctx context.Context) error) []error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, len(operations))
	var wg sync.WaitGroup

	for _, op := range operations {
		wg.Add(1)
		go func(fn func(ctx context.Context) error) {
			defer wg.Done()
			if err := fn(ctx); err != nil {
				errCh <- err
				cancel() // signal others to stop
				return
			}
			errCh <- nil
		}(op)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var errs []error
	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}