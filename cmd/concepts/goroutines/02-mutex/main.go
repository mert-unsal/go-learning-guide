// Mutex — standalone demonstration of sync.Mutex, sync.RWMutex,
// race conditions, and safe concurrent data access.
//
// Run: go run ./cmd/concepts/goroutines/02-mutex
// Test for races: go run -race ./cmd/concepts/goroutines/02-mutex
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

// ============================================================
// RACE CONDITIONS AND MUTEX
// ============================================================
// Race condition: two goroutines access the same data concurrently
// and at least one is a write. This is UNDEFINED BEHAVIOR.
// Always run with -race flag during development: go test -race ./...
//
// sync.Mutex internals:
//   - Uses spinning (for short critical sections) + OS semaphore (for long waits)
//   - Starvation mode (Go 1.9+): after 1ms of waiting, switches to FIFO
//   - Never copy a Mutex after first use (contains internal state)

// UNSAFE counter — race condition!
// Multiple goroutines doing c.value++ concurrently is a read-modify-write
// that is NOT atomic. The race detector (go run -race) catches this.
type UnsafeCounter struct {
	value int
}

func (c *UnsafeCounter) Increment() {
	c.value++ // NOT SAFE: read-modify-write is not atomic
}

// SAFE counter with Mutex
// Lock() acquires exclusive access. defer Unlock() ensures release
// even if a panic occurs inside the critical section.
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// ============================================================
// RWMutex — allows multiple concurrent READERS, exclusive WRITER
// ============================================================
// Use RWMutex when reads vastly outnumber writes (e.g., config, cache).
// RLock() allows concurrent readers. Lock() blocks all readers AND writers.
// Under contention with many writers, readers can starve — understand
// the tradeoffs vs a regular Mutex.
type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock() // multiple goroutines can read simultaneously
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock() // exclusive write lock
	defer c.mu.Unlock()
	c.data[key] = value
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Mutex, RWMutex & Race Conditions        %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Unsafe counter (without mutex) ---
	fmt.Printf("%s▸ Without Mutex — race condition demo%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ Running 1000 goroutines incrementing a plain int (no lock)%s\n", yellow, reset)
	fmt.Printf("  %s⚠ This is a data race: read-modify-write on c.value is NOT atomic%s\n", yellow, reset)

	unsafeCounter := &UnsafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unsafeCounter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("  Unsafe counter result: %s%d%s  %s(expected 1000 — may be less due to lost updates!)%s\n", red+bold, unsafeCounter.value, reset, dim, reset)
	fmt.Printf("  %s⚠ 'go run -race' would flag this as a data race — always use -race in dev!%s\n\n", yellow, reset)

	// --- Safe counter with Mutex ---
	fmt.Printf("%s▸ With sync.Mutex — safe concurrent access%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Lock() acquires exclusive access; defer Unlock() ensures release even on panic%s\n", green, reset)

	counter := &SafeCounter{}
	var activeGoroutines atomic.Int64

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			activeGoroutines.Add(1)
			counter.Increment()
			activeGoroutines.Add(-1)
		}()
	}
	wg.Wait()
	fmt.Printf("  Safe counter result:   %s%d%s  %s(always exactly 1000 — mutex serializes access)%s\n", green+bold, counter.Value(), reset, dim, reset)
	fmt.Printf("  %s✔ Mutex uses spinning for short waits + OS semaphore for long waits%s\n", green, reset)
	fmt.Printf("  %s✔ Starvation mode (Go 1.9+): after 1ms of waiting, switches to FIFO fairness%s\n", green, reset)

	// --- Cache with RWMutex ---
	fmt.Printf("\n%s▸ sync.RWMutex — concurrent readers, exclusive writers%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ RLock() allows multiple simultaneous readers; Lock() is exclusive%s\n", green, reset)

	cache := &Cache{data: make(map[string]string)}

	// Concurrent writes
	fmt.Printf("  %sWriting 10 keys concurrently (each needs exclusive Lock)...%s\n", dim, reset)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i))
		}()
	}
	wg.Wait()
	fmt.Printf("  %s✔ All 10 writes completed — each held exclusive lock%s\n", green, reset)

	// Concurrent reads — all can proceed simultaneously with RLock
	fmt.Printf("  %sReading 10 keys concurrently (all share RLock)...%s\n", dim, reset)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			if v, ok := cache.Get(fmt.Sprintf("key%d", i)); ok {
				fmt.Printf("    cache[key%d] = %s%s%s\n", i, magenta, v, reset)
			}
		}()
	}
	wg.Wait()

	fmt.Printf("\n%s▸ When to use which%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Mutex:   simple mutual exclusion, balanced read/write workloads%s\n", green, reset)
	fmt.Printf("  %s✔ RWMutex: reads vastly outnumber writes (caches, config)%s\n", green, reset)
	fmt.Printf("  %s✔ atomic:  single counters/flags — no lock overhead at all%s\n", green, reset)
	fmt.Printf("  %s⚠ Under contention with many writers, RWMutex readers can starve%s\n", yellow, reset)
}
