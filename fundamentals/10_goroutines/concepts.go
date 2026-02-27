// Package goroutines covers Go's concurrency model:
// goroutines, sync.WaitGroup, sync.Mutex, and race conditions.
package goroutines

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================
// 1. GOROUTINES
// ============================================================
// A goroutine is a lightweight thread managed by the Go runtime.
// Cost: ~2KB of stack (vs ~1MB for OS threads).
// Go can run millions of goroutines concurrently.
// Start with: go functionCall()

func DemonstrateGoroutines() {
	// Sequential vs Concurrent
	fmt.Println("Sequential:")
	doWork("A")
	doWork("B")
	doWork("C")

	fmt.Println("\nConcurrent:")
	var wg sync.WaitGroup

	for _, name := range []string{"A", "B", "C"} {
		wg.Add(1)    // increment counter BEFORE goroutine starts
		name := name // capture loop variable (important!)
		go func() {
			defer wg.Done() // decrement when done
			doWork(name)
		}()
	}

	wg.Wait() // block until all goroutines call Done()
	fmt.Println("All done!")
}

func doWork(name string) {
	fmt.Printf("  Worker %s starting\n", name)
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("  Worker %s done\n", name)
}

// ============================================================
// 2. RACE CONDITIONS AND MUTEX
// ============================================================
// Race condition: two goroutines access the same data concurrently
// and at least one is a write. This is UNDEFINED BEHAVIOR.

// UNSAFE counter — race condition!
type UnsafeCounter struct {
	value int
}

func (c *UnsafeCounter) Increment() {
	c.value++ // NOT SAFE: read-modify-write is not atomic
}

// SAFE counter with Mutex
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

// RWMutex: allows multiple concurrent READERS, exclusive WRITER
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

func DemonstrateMutex() {
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
}

// ============================================================
// 3. sync.Once — run code exactly once
// ============================================================

type Singleton struct {
	data string
}

var (
	instance *Singleton
	once     sync.Once
)

func GetSingleton() *Singleton {
	once.Do(func() {
		fmt.Println("Initializing singleton...")
		instance = &Singleton{data: "initialized"}
	})
	return instance
}

// ============================================================
// 4. sync.Pool — reuse objects to reduce GC pressure
// ============================================================

var bufPool = sync.Pool{
	New: func() any {
		return make([]byte, 1024)
	},
}

func DemonstratePool() {
	// Get a buffer from the pool
	buf := bufPool.Get().([]byte)
	// Use the buffer...
	_ = buf
	// Return to pool when done
	bufPool.Put(buf)
}

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Goroutines ===")
	DemonstrateGoroutines()
	fmt.Println("\n=== Mutex ===")
	DemonstrateMutex()
	fmt.Println("\n=== Singleton ===")
	_ = GetSingleton()
	_ = GetSingleton() // only initializes once
}
