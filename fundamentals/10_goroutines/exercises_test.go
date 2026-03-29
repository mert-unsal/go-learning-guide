package goroutines

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunConcurrently(t *testing.T) {
	var counter int64
	RunConcurrently(100, func(id int) {
		atomic.AddInt64(&counter, 1)
	})
	if counter != 100 {
		t.Errorf("❌ RunConcurrently: counter = %d, want 100  ← Hint: use WaitGroup", counter)
	} else {
		t.Logf("✅ RunConcurrently ran 100 goroutines, counter = %d", counter)
	}
}

func TestExCounter(t *testing.T) {
	c := &ExCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	got := c.Value()
	if got != 1000 {
		t.Errorf("❌ ExCounter = %d, want 1000  ← Hint: possible race — use sync.Mutex", got)
	} else {
		t.Logf("✅ ExCounter after 1000 concurrent Inc() = %d", got)
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
		t.Run(fmt.Sprintf("SumConcurrent(%v)", tt.nums), func(t *testing.T) {
			got := SumConcurrent(tt.nums)
			if got != tt.want {
				t.Errorf("❌ SumConcurrent(%v) = %d, want %d", tt.nums, got, tt.want)
			} else {
				t.Logf("✅ SumConcurrent(%v) = %d", tt.nums, got)
			}
		})
	}
}

func TestRunOnce(t *testing.T) {
	calls := 0
	RunOnce(func() { calls++ })
	RunOnce(func() { calls++ })
	RunOnce(func() { calls++ })
	if calls != 1 {
		t.Errorf("❌ setup ran %d times, want exactly 1  ← Hint: use sync.Once", calls)
	} else {
		t.Logf("✅ RunOnce called 3 times but setup ran only %d time", calls)
	}
}

// ── Exercise 5: WaitGroup + Mutex Error Collection ──

func TestCollectAllErrors(t *testing.T) {
	t.Run("all succeed", func(t *testing.T) {
		ops := []func() error{
			func() error { return nil },
			func() error { return nil },
			func() error { return nil },
		}
		errs := CollectAllErrors(ops)
		if len(errs) != 0 {
			t.Errorf("❌ got %d errors, want 0", len(errs))
		} else {
			t.Log("✅ all succeed → 0 errors collected")
		}
	})

	t.Run("some fail", func(t *testing.T) {
		ops := []func() error{
			func() error { return nil },
			func() error { return fmt.Errorf("db timeout") },
			func() error { return nil },
			func() error { return fmt.Errorf("api 503") },
		}
		errs := CollectAllErrors(ops)
		if len(errs) != 2 {
			t.Errorf("❌ got %d errors, want 2  ← Hint: append non-nil errors under mutex", len(errs))
		} else {
			t.Logf("✅ 2 of 4 ops failed → collected %d errors", len(errs))
		}
	})

	t.Run("all fail", func(t *testing.T) {
		ops := []func() error{
			func() error { return fmt.Errorf("err1") },
			func() error { return fmt.Errorf("err2") },
			func() error { return fmt.Errorf("err3") },
		}
		errs := CollectAllErrors(ops)
		if len(errs) != 3 {
			t.Errorf("❌ got %d errors, want 3", len(errs))
		} else {
			t.Logf("✅ all 3 ops failed → collected %d errors", len(errs))
		}
	})

	t.Run("empty operations", func(t *testing.T) {
		errs := CollectAllErrors(nil)
		if len(errs) != 0 {
			t.Errorf("❌ got %d errors for nil input, want 0", len(errs))
		} else {
			t.Log("✅ nil operations → 0 errors")
		}
	})

	t.Run("concurrent safety", func(t *testing.T) {
		ops := make([]func() error, 1000)
		for i := range ops {
			ops[i] = func() error { return fmt.Errorf("err") }
		}
		errs := CollectAllErrors(ops)
		if len(errs) != 1000 {
			t.Errorf("❌ got %d errors, want 1000  ← Hint: race condition — protect with mutex", len(errs))
		} else {
			t.Logf("✅ 1000 concurrent errors collected safely")
		}
	})
}

// ── Exercise 6: Channel-Based Error Collection ──

func TestCollectErrorsViaChan(t *testing.T) {
	t.Run("all succeed", func(t *testing.T) {
		ops := []func() error{
			func() error { return nil },
			func() error { return nil },
		}
		errs := CollectErrorsViaChan(ops)
		if len(errs) != 0 {
			t.Errorf("❌ got %d errors, want 0", len(errs))
		} else {
			t.Log("✅ all succeed → 0 errors via channel")
		}
	})

	t.Run("mixed results", func(t *testing.T) {
		ops := []func() error{
			func() error { return nil },
			func() error { return fmt.Errorf("timeout") },
			func() error { return nil },
			func() error { return fmt.Errorf("connection refused") },
			func() error { return fmt.Errorf("not found") },
		}
		errs := CollectErrorsViaChan(ops)
		if len(errs) != 3 {
			t.Errorf("❌ got %d errors, want 3  ← Hint: buffer channel, use closer goroutine", len(errs))
		} else {
			t.Logf("✅ 3 of 5 ops failed → collected %d errors via channel", len(errs))
		}
	})

	t.Run("no deadlock on empty", func(t *testing.T) {
		errs := CollectErrorsViaChan(nil)
		if len(errs) != 0 {
			t.Errorf("❌ got %d errors for nil input, want 0", len(errs))
		} else {
			t.Log("✅ nil operations → no deadlock, 0 errors")
		}
	})
}

// ── Exercise 7: Context-Aware Error Collection ──

func TestCollectErrorsWithCancel(t *testing.T) {
	t.Run("all succeed", func(t *testing.T) {
		ctx := context.Background()
		ops := []func(ctx context.Context) error{
			func(ctx context.Context) error { return nil },
			func(ctx context.Context) error { return nil },
		}
		errs := CollectErrorsWithCancel(ctx, ops)
		if len(errs) != 0 {
			t.Errorf("❌ got %d errors, want 0", len(errs))
		} else {
			t.Log("✅ all succeed → 0 errors, context not cancelled")
		}
	})

	t.Run("first error cancels context", func(t *testing.T) {
		ctx := context.Background()
		var cancelled atomic.Bool

		ops := []func(ctx context.Context) error{
			func(ctx context.Context) error {
				return fmt.Errorf("immediate failure")
			},
			func(ctx context.Context) error {
				// Simulate slow work that checks context
				select {
				case <-ctx.Done():
					cancelled.Store(true)
					return ctx.Err()
				case <-time.After(2 * time.Second):
					return nil
				}
			},
		}

		errs := CollectErrorsWithCancel(ctx, ops)
		if len(errs) == 0 {
			t.Error("❌ got 0 errors, want at least 1  ← Hint: cancel context on first error")
		} else {
			t.Logf("✅ collected %d error(s), context cancellation propagated: %v", len(errs), cancelled.Load())
		}
	})

	t.Run("respects parent cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // already cancelled

		ops := []func(ctx context.Context) error{
			func(ctx context.Context) error {
				return ctx.Err()
			},
		}

		errs := CollectErrorsWithCancel(ctx, ops)
		if len(errs) == 0 {
			t.Error("❌ got 0 errors, want 1 (parent was cancelled)")
		} else {
			t.Logf("✅ parent cancellation propagated → %d error(s)", len(errs))
		}
	})
}
