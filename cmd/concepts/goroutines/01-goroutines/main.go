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

func doWork(name string, start time.Time) {
	launchDelay := time.Since(start)
	fmt.Printf("  %s⚡ Worker %s launched%s  %s(+%v since main)%s\n", green, name, reset, dim, launchDelay, reset)
	time.Sleep(10 * time.Millisecond)
	elapsed := time.Since(start)
	fmt.Printf("  %s✔ Worker %s done%s     %s(+%v since main)%s\n", green, name, reset, dim, elapsed, reset)
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Goroutines & WaitGroup                 %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Sequential ---
	fmt.Printf("%s▸ Sequential execution%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Each worker blocks main until it finishes — total time ≈ 3 × 10ms%s\n", green, reset)
	seqStart := time.Now()
	doWork("A", seqStart)
	doWork("B", seqStart)
	doWork("C", seqStart)
	seqElapsed := time.Since(seqStart)
	fmt.Printf("  Total sequential time: %s%v%s\n", magenta, seqElapsed, reset)

	// --- Concurrent ---
	fmt.Printf("\n%s▸ Concurrent execution with goroutines%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ All workers launch nearly instantly — total time ≈ 10ms (parallel)%s\n", green, reset)
	fmt.Printf("  %s✔ 'go func()' spawns a goroutine: ~2KB stack, managed by GMP scheduler%s\n", green, reset)

	concStart := time.Now()
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
			doWork(name, concStart)
		}()
	}

	wg.Wait() // block until all goroutines call Done()
	concElapsed := time.Since(concStart)
	fmt.Printf("  Total concurrent time: %s%v%s\n", magenta, concElapsed, reset)

	// --- Comparison ---
	fmt.Printf("\n%s▸ Key observations%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Sequential: ~%v  vs  Concurrent: ~%v%s\n", green, seqElapsed.Round(time.Millisecond), concElapsed.Round(time.Millisecond), reset)
	fmt.Printf("  %s✔ WaitGroup: Add(1) before launch, defer Done() inside goroutine, Wait() to block%s\n", green, reset)
	fmt.Printf("  %s⚠ Always capture loop variables — without 'name := name' all goroutines see last value%s\n", yellow, reset)
	fmt.Printf("  %s⚠ Goroutine order is non-deterministic — the scheduler decides who runs when%s\n", yellow, reset)
}
