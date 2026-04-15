package concurrency

import (
	"context"
	"sort"
	"sync/atomic"
	"testing"
	"time"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: FanOut
// ────────────────────────────────────────────────────────────

func TestFanOut(t *testing.T) {
	tasks := []func() int{
		func() int { return 1 },
		func() int { return 2 },
		func() int { return 3 },
	}
	results := FanOut(tasks)
	if results == nil || len(results) != 3 {
		t.Fatal("❌ FanOut returned nil or wrong length\n\t\t" +
			"Hint: results := make([]int, len(tasks)); var wg sync.WaitGroup; " +
			"for i, task := range tasks { wg.Add(1); go func(i int, fn func() int) { " +
			"results[i] = fn(); wg.Done() }(i, task) }; wg.Wait(); return results")
	}
	if results[0] != 1 || results[1] != 2 || results[2] != 3 {
		t.Errorf("❌ FanOut = %v, want [1 2 3] (order preserved)", results)
	} else {
		t.Logf("✅ FanOut = %v", results)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: Pipeline
// ────────────────────────────────────────────────────────────

func TestPipeline(t *testing.T) {
	out := Pipeline([]int{1, 2, 3})
	if out == nil {
		t.Fatal("❌ Pipeline returned nil channel\n\t\t" +
			"Hint: stage1 sends numbers to ch1; stage2 reads ch1, doubles, sends to ch2; " +
			"stage3 reads ch2, adds 10, sends to out. Close channels when done")
	}
	var results []int
	for v := range out {
		results = append(results, v)
	}
	// 1→double→2→+10→12, 2→4→14, 3→6→16
	if len(results) != 3 {
		t.Fatalf("❌ Pipeline produced %d results, want 3", len(results))
	}
	sort.Ints(results)
	if results[0] != 12 || results[1] != 14 || results[2] != 16 {
		t.Errorf("❌ Pipeline = %v, want [12 14 16]", results)
	} else {
		t.Logf("✅ Pipeline([1,2,3]) → %v", results)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: WorkerPool
// ────────────────────────────────────────────────────────────

func TestWorkerPool(t *testing.T) {
	results := WorkerPool([]int{1, 2, 3, 4, 5}, 2, func(n int) int { return n * n })
	if results == nil || len(results) != 5 {
		t.Fatal("❌ WorkerPool returned nil or wrong length\n\t\t" +
			"Hint: jobs := make(chan int); resultCh := make(chan int); " +
			"start numWorkers goroutines reading from jobs; " +
			"send all jobs; close(jobs); collect results")
	}
	sort.Ints(results)
	if results[0] != 1 || results[4] != 25 {
		t.Errorf("❌ WorkerPool = %v, want [1 4 9 16 25]", results)
	} else {
		t.Logf("✅ WorkerPool (2 workers) = %v", results)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: Semaphore
// ────────────────────────────────────────────────────────────

func TestSemaphore(t *testing.T) {
	var running int32
	var maxSeen int32

	tasks := make([]func(), 10)
	for i := range tasks {
		tasks[i] = func() {
			cur := atomic.AddInt32(&running, 1)
			for {
				old := atomic.LoadInt32(&maxSeen)
				if cur <= old || atomic.CompareAndSwapInt32(&maxSeen, old, cur) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt32(&running, -1)
		}
	}

	Semaphore(tasks, 3)

	if atomic.LoadInt32(&maxSeen) > 3 {
		t.Errorf("❌ max concurrent = %d, want ≤ 3\n\t\t"+
			"Hint: sem := make(chan struct{}, maxConcurrent); "+
			"for each task: sem <- struct{}{}; go func() { task(); <-sem }(); "+
			"wait for all to finish",
			maxSeen)
	} else if atomic.LoadInt32(&maxSeen) == 0 {
		t.Error("❌ no tasks ran (maxSeen = 0)")
	} else {
		t.Logf("✅ Semaphore limited to max %d concurrent", maxSeen)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: WithTimeout
// ────────────────────────────────────────────────────────────

func TestWithTimeout(t *testing.T) {
	// Fast function
	result, err := WithTimeout(func() string { return "fast" }, 100*time.Millisecond)
	if err != nil || result != "fast" {
		t.Errorf("❌ WithTimeout(fast) = (%q, %v), want (\"fast\", nil)\n\t\t"+
			"Hint: use select with time.After(timeout) and a result channel",
			result, err)
	} else {
		t.Logf("✅ WithTimeout(fast) = %q", result)
	}

	// Slow function (should timeout)
	_, err = WithTimeout(func() string {
		time.Sleep(200 * time.Millisecond)
		return "slow"
	}, 50*time.Millisecond)
	if err != ErrTimeout {
		t.Errorf("❌ WithTimeout(slow) error = %v, want ErrTimeout", err)
	} else {
		t.Logf("✅ WithTimeout(slow) correctly timed out")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: Merge
// ────────────────────────────────────────────────────────────

func TestMerge(t *testing.T) {
	ch1 := make(chan int, 2)
	ch2 := make(chan int, 2)
	ch1 <- 1
	ch1 <- 2
	close(ch1)
	ch2 <- 3
	ch2 <- 4
	close(ch2)

	out := Merge(ch1, ch2)
	if out == nil {
		t.Fatal("❌ Merge returned nil\n\t\t" +
			"Hint: var wg sync.WaitGroup; for each ch: wg.Add(1); go func(c <-chan int) { " +
			"for v := range c { out <- v }; wg.Done() }(ch); go func() { wg.Wait(); close(out) }()")
	}

	var results []int
	for v := range out {
		results = append(results, v)
	}
	sort.Ints(results)
	if len(results) != 4 || results[0] != 1 || results[3] != 4 {
		t.Errorf("❌ Merge = %v, want [1 2 3 4]", results)
	} else {
		t.Logf("✅ Merge = %v", results)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 7: Generate
// ────────────────────────────────────────────────────────────

func TestGenerate(t *testing.T) {
	ch := Generate(5)
	if ch == nil {
		t.Fatal("❌ Generate returned nil\n\t\t" +
			"Hint: ch := make(chan int); go func() { for i := 0; i < n; i++ { ch <- i }; close(ch) }(); return ch")
	}
	var results []int
	for v := range ch {
		results = append(results, v)
	}
	if len(results) != 5 || results[0] != 0 || results[4] != 4 {
		t.Errorf("❌ Generate(5) = %v, want [0 1 2 3 4]", results)
	} else {
		t.Logf("✅ Generate(5) = %v", results)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: OrDone
// ────────────────────────────────────────────────────────────

func TestOrDone(t *testing.T) {
	ch := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		ch <- i
	}
	close(ch)

	ctx := context.Background()
	results := OrDone(ctx, ch)
	if len(results) != 5 {
		t.Errorf("❌ OrDone(full channel) = %v, want [1 2 3 4 5]\n\t\t"+
			"Hint: for { select { case <-ctx.Done(): return result; "+
			"case v, ok := <-ch: if !ok { return result }; result = append(result, v) } }",
			results)
	} else {
		t.Logf("✅ OrDone(full) = %v", results)
	}

	// Context cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	slowCh := make(chan int)
	results = OrDone(ctx, slowCh)
	if len(results) != 0 {
		t.Errorf("❌ OrDone(cancelled) = %v, want []", results)
	} else {
		t.Logf("✅ OrDone(cancelled context) returned immediately")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 9: RateLimited
// ────────────────────────────────────────────────────────────

func TestRateLimited(t *testing.T) {
	var count int32
	items := []int{1, 2, 3}
	start := time.Now()

	RateLimited(items, 100, func(n int) {
		atomic.AddInt32(&count, 1)
	})

	elapsed := time.Since(start)

	if atomic.LoadInt32(&count) != 3 {
		t.Errorf("❌ RateLimited processed %d items, want 3\n\t\t"+
			"Hint: ticker := time.NewTicker(time.Second / time.Duration(rate)); "+
			"for _, item := range items { <-ticker.C; fn(item) }",
			count)
	} else if elapsed < 20*time.Millisecond {
		t.Logf("⚠️ RateLimited finished in %v (rate limiting may not be working)", elapsed)
	} else {
		t.Logf("✅ RateLimited processed 3 items in %v", elapsed)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 10: SafeMap
// ────────────────────────────────────────────────────────────

func TestSafeMap(t *testing.T) {
	m := NewSafeMap()
	m.Set("a", 1)
	m.Set("b", 2)

	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Errorf("❌ Get(a) = (%d, %v), want (1, true)\n\t\t"+
			"Hint: s.mu.RLock(); defer s.mu.RUnlock(); v, ok := s.m[key]; return v, ok",
			v, ok)
	} else {
		t.Logf("✅ Get(a) = %d", v)
	}

	if m.Len() != 2 {
		t.Errorf("❌ Len = %d, want 2", m.Len())
	}

	// Concurrent safety test
	done := make(chan struct{})
	for i := 0; i < 100; i++ {
		go func(i int) {
			m.Set("key", i)
			m.Get("key")
			<-done
		}(i)
	}
	close(done)
}

// ────────────────────────────────────────────────────────────
// Exercise 11: Barrier
// ────────────────────────────────────────────────────────────

func TestBarrier(t *testing.T) {
	tasks := []func() int{
		func() int { return 10 },
		func() int { return 20 },
		func() int { return 30 },
	}
	result := Barrier(tasks, func(results []int) int {
		sum := 0
		for _, v := range results {
			sum += v
		}
		return sum
	})
	if result != 60 {
		t.Errorf("❌ Barrier = %d, want 60\n\t\t"+
			"Hint: run all tasks (use FanOut pattern), collect results, then call afterAll",
			result)
	} else {
		t.Logf("✅ Barrier(10+20+30) = %d", result)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: GracefulWorker
// ────────────────────────────────────────────────────────────

func TestGracefulWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan int, 5)
	jobs <- 1
	jobs <- 2
	jobs <- 3
	close(jobs)

	results := GracefulWorker(ctx, jobs, func(n int) int { return n * 10 })
	cancel()

	sort.Ints(results)
	if len(results) != 3 || results[0] != 10 || results[2] != 30 {
		t.Errorf("❌ GracefulWorker = %v, want [10 20 30]\n\t\t"+
			"Hint: for { select { case job, ok := <-jobs: if !ok { return }; "+
			"results = append(results, fn(job)); case <-ctx.Done(): "+
			"drain remaining jobs from channel, then return } }",
			results)
	} else {
		t.Logf("✅ GracefulWorker = %v", results)
	}
}
