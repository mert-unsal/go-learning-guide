// Mutex — standalone demonstration of sync.Mutex, sync.RWMutex,
// race conditions, and safe concurrent data access.
//
// Run: go run ./cmd/concepts/goroutines/02-mutex
// Test for races: go run -race ./cmd/concepts/goroutines/02-mutex
package main

import (
	"fmt"
	"sync"
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
	// --- SafeCounter with Mutex ---
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Println("Safe counter:", counter.Value()) // always 1000

	// --- Cache with RWMutex ---
	cache := &Cache{data: make(map[string]string)}

	// Concurrent writes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i))
		}()
	}
	wg.Wait()

	// Concurrent reads — all can proceed simultaneously with RLock
	for i := 0; i < 10; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			if v, ok := cache.Get(fmt.Sprintf("key%d", i)); ok {
				fmt.Printf("  cache[key%d] = %s\n", i, v)
			}
		}()
	}
	wg.Wait()
}
