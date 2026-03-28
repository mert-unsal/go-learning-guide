// Pool & Once — standalone demonstration of sync.Pool for object
// reuse and sync.Once for singleton initialization.
//
// Run: go run ./cmd/concepts/goroutines/03-pool
package main

import (
	"fmt"
	"sync"
)

// ============================================================
// sync.Pool — reuse objects to reduce GC pressure
// ============================================================
// sync.Pool is a per-P (per logical processor) cache of reusable objects.
// Objects may be evicted at any GC cycle (they survive one cycle in
// the victim cache, then are freed).
//
// KEY: sync.Pool is NOT a connection pool. Objects can disappear.
// Use it for short-lived allocations in hot paths: byte buffers,
// temporary structs, encoder/decoder state.
//
// Production pattern:
//   1. Get() retrieves an object (or calls New if pool is empty)
//   2. Use the object
//   3. Reset/clear the object
//   4. Put() returns it to the pool

var bufPool = sync.Pool{
	New: func() any {
		fmt.Println("  Pool: allocating new buffer")
		return make([]byte, 1024)
	},
}

// ============================================================
// sync.Once — run code exactly once (thread-safe)
// ============================================================
// Uses atomic fast path + mutex slow path internally.
// Guarantees exactly one execution even under racing callers.
// Common use: singleton initialization, one-time config loading.

type Singleton struct {
	data string
}

var (
	instance *Singleton
	once     sync.Once
)

func GetSingleton() *Singleton {
	once.Do(func() {
		fmt.Println("  Initializing singleton...")
		instance = &Singleton{data: "initialized"}
	})
	return instance
}

func main() {
	// --- sync.Pool demo ---
	fmt.Println("=== sync.Pool ===")

	// First Get — pool is empty, calls New
	buf := bufPool.Get().([]byte)
	fmt.Printf("  Got buffer of len %d\n", len(buf))

	// Return to pool
	bufPool.Put(buf)

	// Second Get — reuses the pooled buffer (no allocation)
	buf2 := bufPool.Get().([]byte)
	fmt.Printf("  Got buffer of len %d (reused)\n", len(buf2))
	bufPool.Put(buf2)

	// --- sync.Once demo ---
	fmt.Println("\n=== sync.Once (Singleton) ===")

	// Call GetSingleton multiple times — init runs only once
	s1 := GetSingleton()
	s2 := GetSingleton()
	s3 := GetSingleton()
	fmt.Printf("  s1 == s2: %v, s2 == s3: %v\n", s1 == s2, s2 == s3)
	fmt.Printf("  Singleton data: %q\n", s1.data)
}
