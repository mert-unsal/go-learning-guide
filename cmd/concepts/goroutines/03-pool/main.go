// Pool & Once — standalone demonstration of sync.Pool for object
// reuse and sync.Once for singleton initialization.
//
// Run: go run ./cmd/concepts/goroutines/03-pool
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

var allocCount atomic.Int64

var bufPool = sync.Pool{
	New: func() any {
		n := allocCount.Add(1)
		fmt.Printf("    %s⚡ Pool.New called — allocating buffer #%d (1024 bytes)%s\n", yellow, n, reset)
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
		fmt.Printf("    %s⚡ once.Do executing — initializing singleton%s\n", yellow, reset)
		instance = &Singleton{data: "initialized"}
	})
	return instance
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  sync.Pool & sync.Once                  %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- sync.Pool demo ---
	fmt.Printf("%s▸ sync.Pool — object reuse to reduce GC pressure%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Per-P (per logical processor) cache of reusable objects%s\n", green, reset)
	fmt.Printf("  %s⚠ NOT a connection pool — objects can vanish at any GC cycle%s\n\n", yellow, reset)

	// First Get — pool is empty, calls New
	fmt.Printf("  %s[Get #1] Pool is empty → calls New:%s\n", dim, reset)
	buf := bufPool.Get().([]byte)
	fmt.Printf("  Got buffer: len=%s%d%s  %s(active: 1, idle in pool: 0)%s\n", magenta, len(buf), reset, dim, reset)

	// Return to pool
	fmt.Printf("\n  %s[Put] Returning buffer to pool%s\n", dim, reset)
	bufPool.Put(buf)
	fmt.Printf("  %s✔ Buffer returned%s  %s(active: 0, idle in pool: 1)%s\n", green, reset, dim, reset)

	// Second Get — reuses the pooled buffer (no allocation)
	fmt.Printf("\n  %s[Get #2] Pool has an idle buffer → reuses it (no alloc):%s\n", dim, reset)
	buf2 := bufPool.Get().([]byte)
	fmt.Printf("  Got buffer: len=%s%d%s  %s(reused! active: 1, idle in pool: 0)%s\n", magenta, len(buf2), reset, dim, reset)
	fmt.Printf("  %s✔ Same backing memory — zero allocation on this Get()%s\n", green, reset)
	bufPool.Put(buf2)

	// Concurrent pool usage
	fmt.Printf("\n  %s[Concurrent] 5 goroutines sharing the pool:%s\n", dim, reset)
	var wg sync.WaitGroup
	var activeWorkers atomic.Int64
	for i := 0; i < 5; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			n := activeWorkers.Add(1)
			b := bufPool.Get().([]byte)
			fmt.Printf("    worker %d: got buffer len=%s%d%s  %s(active workers: %d)%s\n", i, magenta, len(b), reset, dim, n, reset)
			activeWorkers.Add(-1)
			bufPool.Put(b)
		}()
	}
	wg.Wait()
	fmt.Printf("  Total allocations by Pool.New: %s%d%s  %s(fewer than 5 means reuse worked!)%s\n", magenta, allocCount.Load(), reset, dim, reset)

	fmt.Printf("\n  %s✔ Production pattern: Get → Use → Reset/Clear → Put%s\n", green, reset)
	fmt.Printf("  %s✔ Great for: byte buffers, encoder state, temporary structs in hot paths%s\n", green, reset)
	fmt.Printf("  %s⚠ Objects survive one GC in victim cache, then freed — never rely on persistence%s\n", yellow, reset)

	// --- sync.Once demo ---
	fmt.Printf("\n%s▸ sync.Once — exactly-once initialization%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Atomic fast path + mutex slow path: guaranteed single execution%s\n", green, reset)
	fmt.Printf("  %s✔ Even under racing callers, only one executes — others block until done%s\n\n", green, reset)

	// Call GetSingleton multiple times — init runs only once
	fmt.Printf("  %s[Call 1] GetSingleton():%s\n", dim, reset)
	s1 := GetSingleton()
	fmt.Printf("  %s[Call 2] GetSingleton():%s  %s(no init — already done)%s\n", dim, reset, dim, reset)
	s2 := GetSingleton()
	fmt.Printf("  %s[Call 3] GetSingleton():%s  %s(no init — already done)%s\n", dim, reset, dim, reset)
	s3 := GetSingleton()

	fmt.Printf("\n  s1 == s2: %s%v%s,  s2 == s3: %s%v%s  ← same pointer, same instance\n", magenta, s1 == s2, reset, magenta, s2 == s3, reset)
	fmt.Printf("  Singleton data: %s%q%s\n", magenta, s1.data, reset)
	fmt.Printf("  %s✔ Common uses: singleton init, one-time config loading, driver registration%s\n", green, reset)
}
