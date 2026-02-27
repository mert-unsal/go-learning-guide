package goroutines

import (
	"sync"
	"sync/atomic"
	"testing"
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
		got := SumConcurrent(tt.nums)
		if got != tt.want {
			t.Errorf("❌ SumConcurrent(%v) = %d, want %d", tt.nums, got, tt.want)
		} else {
			t.Logf("✅ SumConcurrent(%v) = %d", tt.nums, got)
		}
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
