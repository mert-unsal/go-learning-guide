// Goroutines — standalone demonstration of goroutine creation,
// sync.WaitGroup, and loop variable capture.
//
// Run: go run ./cmd/concepts/goroutines/01-goroutines
package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================
// GOROUTINES
// ============================================================
// A goroutine is a lightweight thread managed by the Go runtime.
// Cost: ~2KB of stack (vs ~1MB for OS threads).
// Go can run millions of goroutines concurrently.
// Start with: go functionCall()
//
// Under the hood (GMP model):
//   G (goroutine) — user-space thread, ~2-8KB initial stack (growable)
//   M (machine)   — OS thread, maps to kernel thread
//   P (processor) — logical processor, holds local run queue
//
// The scheduler is cooperatively preemptive (async preemption since Go 1.14).
// Goroutines blocked on I/O are parked by the network poller (epoll/kqueue/IOCP)
// and do NOT consume an OS thread while waiting.

func doWork(name string) {
	fmt.Printf("  Worker %s starting\n", name)
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("  Worker %s done\n", name)
}

func main() {
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
		// Without this re-declaration, all goroutines would share the
		// same loop variable and likely all print the last value ("C").
		// Go 1.22+ changes loop variable semantics, but this pattern
		// remains idiomatic for backward compatibility.
		go func() {
			defer wg.Done() // decrement when done
			doWork(name)
		}()
	}

	wg.Wait() // block until all goroutines call Done()
	fmt.Println("All done!")
}
