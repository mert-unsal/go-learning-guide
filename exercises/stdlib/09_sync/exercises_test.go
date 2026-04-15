package sync_exercises

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: SafeCounter
// ────────────────────────────────────────────────────────────

func TestSafeCounter(t *testing.T) {
	c := NewSafeCounter()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment("key")
		}()
	}
	wg.Wait()
	got := c.Value("key")
	if got != 1000 {
		t.Errorf("❌ SafeCounter = %d after 1000 increments, want 1000\n\t\t"+
			"Hint: c.mu.Lock(); c.m[key]++; c.mu.Unlock(). "+
			"Without the lock, concurrent map writes cause fatal crash (not just wrong count)",
			got)
	} else {
		t.Logf("✅ SafeCounter = %d (1000 concurrent increments)", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: RWCache
// ────────────────────────────────────────────────────────────

func TestRWCache(t *testing.T) {
	c := NewRWCache()
	c.Set("greeting", "hello")

	// Concurrent reads should not block each other
	var wg sync.WaitGroup
	var readCount int64
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, ok := c.Get("greeting")
			if ok && val == "hello" {
				atomic.AddInt64(&readCount, 1)
			}
		}()
	}
	wg.Wait()

	if readCount != 100 {
		t.Errorf("❌ RWCache reads = %d, want 100\n\t\t"+
			"Hint: c.mu.RLock() for reads (shared), c.mu.Lock() for writes (exclusive). "+
			"RWMutex allows concurrent readers — Mutex does not",
			readCount)
	} else {
		t.Logf("✅ RWCache: 100 concurrent reads OK")
	}

	// Write should be exclusive
	c.Set("greeting", "world")
	val, _ := c.Get("greeting")
	if val != "world" {
		t.Errorf("❌ after Set(\"world\"), Get = %q", val)
	} else {
		t.Logf("✅ RWCache: Set/Get = %q", val)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: WaitForAll
// ────────────────────────────────────────────────────────────

func TestWaitForAll(t *testing.T) {
	var count int64
	got := WaitForAll(50, func(i int) {
		atomic.AddInt64(&count, 1)
	})
	if got != 50 {
		t.Errorf("❌ WaitForAll returned %d, want 50\n\t\t"+
			"Hint: var wg sync.WaitGroup; wg.Add(n); for i := 0; i < n; i++ { "+
			"go func() { defer wg.Done(); fn(i) }() }; wg.Wait(); return n",
			got)
	} else {
		t.Logf("✅ WaitForAll completed %d goroutines", got)
	}
	if count != 50 {
		t.Errorf("❌ fn called %d times, want 50", count)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: InitOnce
// ────────────────────────────────────────────────────────────

func TestInitOnce(t *testing.T) {
	s := &Service{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(dsn string) {
			defer wg.Done()
			s.Init(dsn)
		}(fmt.Sprintf("dsn-%d", i))
	}
	wg.Wait()

	conn := s.ConnString()
	if conn == "" {
		t.Error("❌ ConnString() = \"\", want non-empty\n\t\t" +
			"Hint: sync.Once.Do(fn) guarantees fn runs exactly once. " +
			"All other callers block until the first Do completes")
	} else {
		t.Logf("✅ ConnString() = %q (initialized exactly once)", conn)
	}

	// Verify subsequent Init calls don't change it
	original := conn
	s.Init("different-dsn")
	if s.ConnString() != original {
		t.Error("❌ ConnString changed after second Init — sync.Once should prevent this")
	} else {
		t.Logf("✅ Second Init correctly ignored")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: BufferPool
// ────────────────────────────────────────────────────────────

func TestBufferPool(t *testing.T) {
	w := GetBuffer()
	if w == nil {
		t.Fatal("❌ GetBuffer() = nil\n\t\t" +
			"Hint: bufferPool.Get().(*bytes.Buffer) — then Reset() to clear. " +
			"sync.Pool.New creates one if pool is empty")
	}
	buf, ok := w.(*bytes.Buffer)
	if !ok {
		t.Fatal("❌ GetBuffer() did not return *bytes.Buffer\n\t\t" +
			"Hint: Return type is io.Writer but underlying is *bytes.Buffer")
	}

	buf.WriteString("test data")
	PutBuffer(buf)

	// Get again — should be reset
	w2 := GetBuffer()
	buf2 := w2.(*bytes.Buffer)
	if buf2.Len() != 0 {
		t.Errorf("❌ buffer from pool has Len=%d, want 0 (should be reset)\n\t\t"+
			"Hint: Reset() the buffer before Put, or after Get — pick one consistently. "+
			"sync.Pool objects can disappear between GC cycles",
			buf2.Len())
	} else {
		t.Logf("✅ BufferPool: Get→Write→Put→Get cycle OK")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: ConcurrentSum
// ────────────────────────────────────────────────────────────

func TestConcurrentSum(t *testing.T) {
	nums := make([]int, 1000)
	expected := 0
	for i := range nums {
		nums[i] = i + 1
		expected += i + 1
	}

	got := ConcurrentSum(nums, 4)
	if got != expected {
		t.Errorf("❌ ConcurrentSum = %d, want %d\n\t\t"+
			"Hint: Split slice into chunks. Each goroutine sums its chunk. "+
			"Use mu.Lock() to add partial sum to total. WaitGroup for sync. "+
			"Handle uneven chunk sizes (len(nums) not divisible by workers)",
			got, expected)
	} else {
		t.Logf("✅ ConcurrentSum(1..1000, 4 workers) = %d", got)
	}

	// Edge cases
	if ConcurrentSum(nil, 4) != 0 {
		t.Error("❌ ConcurrentSum(nil) should be 0")
	}
	if ConcurrentSum([]int{42}, 4) != 42 {
		t.Error("❌ ConcurrentSum([42]) should be 42")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 7: SyncMapStore
// ────────────────────────────────────────────────────────────

func TestSyncMapStore(t *testing.T) {
	pairs := map[string]int{
		"alpha": 1, "beta": 2, "gamma": 3, "delta": 4, "epsilon": 5,
	}
	got := SyncMapStore(pairs)
	if got != 5 {
		t.Errorf("❌ SyncMapStore count = %d, want 5\n\t\t"+
			"Hint: var m sync.Map; goroutine per pair: m.Store(k, v); "+
			"after wg.Wait(), m.Range(func(k,v any) bool { count++; return true }). "+
			"sync.Map is optimized for disjoint-key writes — perfect here",
			got)
	} else {
		t.Logf("✅ SyncMapStore: stored and counted %d pairs", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: OnceValue (Go 1.21+)
// ────────────────────────────────────────────────────────────

func TestMakeLazy(t *testing.T) {
	var callCount int64
	factory := func() string {
		atomic.AddInt64(&callCount, 1)
		return "expensive-result"
	}

	getter := MakeLazy(factory)

	// Call from multiple goroutines
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val := getter()
			if val != "expensive-result" {
				t.Errorf("❌ getter() = %q, want \"expensive-result\"", val)
			}
		}()
	}
	wg.Wait()

	if callCount != 1 {
		t.Errorf("❌ factory called %d times, want exactly 1\n\t\t"+
			"Hint: sync.OnceValue(factory) returns a func() T that calls factory "+
			"exactly once, caches the result, and returns it on all subsequent calls. "+
			"Added in Go 1.21 — cleaner than sync.Once + separate variable",
			callCount)
	} else {
		t.Logf("✅ MakeLazy: factory called exactly once, 100 goroutines got result")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 9: TimedMutex
// ────────────────────────────────────────────────────────────

func TestTimedMutex(t *testing.T) {
	m := NewTimedMutex()

	// Basic lock/unlock
	m.Lock()
	m.Unlock()

	// TryLock should succeed when unlocked
	timeout := make(chan struct{})
	if !m.TryLock(timeout) {
		t.Error("❌ TryLock failed on unlocked mutex\n\t\t" +
			"Hint: select { case <-m.ch: return true; case <-timeout: return false }")
	} else {
		t.Logf("✅ TryLock succeeded on unlocked mutex")
	}

	// TryLock should fail when locked (with immediate timeout)
	close(timeout) // already closed = immediate timeout
	if m.TryLock(timeout) {
		t.Error("❌ TryLock succeeded on locked mutex with expired timeout")
		m.Unlock()
	} else {
		t.Logf("✅ TryLock correctly failed with expired timeout")
	}

	m.Unlock()

	// TryLock with real timeout
	m.Lock()
	timeoutCh := make(chan struct{})
	go func() {
		time.Sleep(50 * time.Millisecond)
		close(timeoutCh)
	}()
	start := time.Now()
	if m.TryLock(timeoutCh) {
		t.Error("❌ TryLock succeeded but mutex was locked")
		m.Unlock()
	} else {
		elapsed := time.Since(start)
		if elapsed < 40*time.Millisecond {
			t.Errorf("❌ TryLock returned too fast: %v", elapsed)
		} else {
			t.Logf("✅ TryLock timed out after %v", elapsed)
		}
	}
	m.Unlock()
}

// ────────────────────────────────────────────────────────────
// Exercise 10: Gate (sync.Cond broadcast)
// ────────────────────────────────────────────────────────────

func TestGate(t *testing.T) {
	g := NewGate()
	var passed int64
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.Wait()
			atomic.AddInt64(&passed, 1)
		}()
	}

	// Give goroutines time to block on Wait
	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt64(&passed) != 0 {
		t.Error("❌ goroutines passed through closed gate")
	}

	g.Open()
	wg.Wait()

	if passed != 10 {
		t.Errorf("❌ %d goroutines passed, want 10\n\t\t"+
			"Hint: sync.Cond wraps a Locker. Wait() atomically unlocks and sleeps. "+
			"Broadcast() wakes ALL waiters. Signal() wakes ONE. "+
			"Always check condition in a for loop (spurious wakeups)",
			passed)
	} else {
		t.Logf("✅ Gate: all 10 goroutines released on Open()")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 11: AtomicConfig
// ────────────────────────────────────────────────────────────

func TestAtomicConfig(t *testing.T) {
	ac := &AtomicConfig{}

	// Default should be zero Config
	cfg := ac.Current()
	if cfg.Host != "" || cfg.Port != 0 {
		t.Errorf("❌ default Config = %+v, want zero value", cfg)
	}

	ac.Update(Config{Host: "localhost", Port: 8080})

	// Concurrent reads
	var wg sync.WaitGroup
	var readCount int64
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := ac.Current()
			if c.Host == "localhost" && c.Port == 8080 {
				atomic.AddInt64(&readCount, 1)
			}
		}()
	}
	wg.Wait()

	if readCount != 100 {
		t.Errorf("❌ %d reads got correct config, want 100\n\t\t"+
			"Hint: atomic.Value.Store(cfg) and atomic.Value.Load().(Config). "+
			"No mutex needed — atomic operations are lock-free. "+
			"Using sync.Map works too: Store(\"cfg\", cfg) / Load(\"cfg\")",
			readCount)
	} else {
		t.Logf("✅ AtomicConfig: 100 concurrent reads OK")
	}

	// Update should be visible to subsequent reads
	ac.Update(Config{Host: "prod", Port: 443})
	cfg = ac.Current()
	if cfg.Host != "prod" {
		t.Errorf("❌ after update, Host = %q, want \"prod\"", cfg.Host)
	} else {
		t.Logf("✅ AtomicConfig: update visible immediately")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: SingleFlight
// ────────────────────────────────────────────────────────────

func TestSingleFlight(t *testing.T) {
	sf := NewSingleFlight()
	var callCount int64

	fn := func() (string, error) {
		atomic.AddInt64(&callCount, 1)
		time.Sleep(50 * time.Millisecond) // simulate expensive work
		return "result", nil
	}

	// Launch 10 concurrent calls with same key
	var wg sync.WaitGroup
	results := make([]string, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			val, err := sf.Do("same-key", fn)
			if err != nil {
				t.Errorf("❌ Do error: %v", err)
			}
			results[idx] = val
		}(i)
	}
	wg.Wait()

	if callCount != 1 {
		t.Errorf("❌ fn called %d times, want 1\n\t\t"+
			"Hint: This is the singleflight pattern (golang.org/x/sync/singleflight). "+
			"Lock → check map → if in-flight: wait on wg → else: add to map, execute, signal. "+
			"Prevents thundering herd on cache miss",
			callCount)
	} else {
		t.Logf("✅ SingleFlight: fn called once for 10 concurrent requests")
	}

	for i, r := range results {
		if r != "result" {
			t.Errorf("❌ results[%d] = %q, want \"result\"", i, r)
		}
	}
	t.Logf("✅ All 10 callers received \"result\"")

	// Different key should trigger new call
	callCount = 0
	sf.Do("different-key", fn)
	if callCount != 1 {
		t.Errorf("❌ different key didn't trigger fn")
	} else {
		t.Logf("✅ Different key correctly triggers new call")
	}
}
