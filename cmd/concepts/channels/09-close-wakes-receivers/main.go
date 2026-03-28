package main

import (
	"fmt"
	"sync"
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
// Close Wakes All Receivers — Broadcast via close()
// ============================================================
//
// The Problem:
//   You have N goroutines all waiting for a signal. Sending N values
//   into a channel is fragile — you need to know N. You need a way
//   to wake all receivers at once, regardless of how many there are.
//
// Real-world example:
//   A service starts with multiple subsystems (HTTP server, gRPC server,
//   metrics exporter, background workers). On shutdown, ONE close(quit)
//   wakes them all — no need to track how many subsystems exist.
//
// The Pattern:
//   close(ch) wakes every goroutine blocked on <-ch simultaneously.
//   Each receiver gets: (zero-value, false). This is how Go implements
//   broadcast — there is no explicit broadcast primitive because
//   close() already does it.
//
// Why channels work here:
//   close() is the ONLY broadcast mechanism for channels. A normal
//   send wakes exactly one receiver. close() wakes ALL of them.
//   This is fundamental to shutdown patterns, done channels, and
//   context cancellation.
//
// Under the hood — runtime.closechan():
//
//   When you call close(ch), the runtime (runtime/chan.go):
//
//   1. Sets ch.closed = 1
//   2. Walks the ENTIRE recvq linked list (all waiting receivers):
//      - Each receiver's sudog is dequeued
//      - Its elem pointer is set to the zero value for the type
//      - The goroutine is marked _Grunnable and put on the run queue
//      - All receivers get: (zero-value-of-T, ok=false)
//   3. Walks the ENTIRE sendq linked list (all waiting senders):
//      - Each sender PANICS with "send on closed channel"
//      - This is why you must close from the sender side, never the receiver
//   4. Releases the channel lock
//
//   ┌──────────────────────────────────────────────────┐
//   │                   close(ch)                      │
//   │                                                  │
//   │  recvq: G1 → G2 → G3   (all waiting on <-ch)    │
//   │    │      │      │                               │
//   │    ▼      ▼      ▼                               │
//   │  wake!  wake!  wake!    (all get zero, false)     │
//   │                                                  │
//   │  sendq: G4 → G5        (any waiting senders)     │
//   │    │      │                                      │
//   │    ▼      ▼                                      │
//   │  PANIC! PANIC!          ("send on closed chan")   │
//   └──────────────────────────────────────────────────┘
//
// Buffered channel nuance:
//   If the channel has remaining data in its buffer when closed:
//   - Receivers get buffered data FIRST (with ok=true)
//   - AFTER the buffer is drained, receivers get (zero, false)
//   This means close() does NOT discard buffered data — it only
//   prevents future sends and signals EOF after the buffer empties.

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Close Wakes All Receivers — Broadcast via close %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	receiverColors := []string{cyan, yellow, magenta}
	receiverNames := []string{"receiver-0", "receiver-1", "receiver-2"}

	fmt.Printf("%s▸ Unbuffered Channel — close() Broadcasts to All Receivers%s\n", cyan+bold, reset)
	fmt.Printf("  %s3 goroutines will park on ch.recvq, then close() wakes them all%s\n\n", dim, reset)

	ch := make(chan int) // unbuffered — receivers will block immediately
	var mu sync.Mutex
	results := make([]string, 3)
	var wg sync.WaitGroup

	// Launch 3 goroutines, all blocking on <-ch
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			val, ok := <-ch // blocks until close(ch)
			mu.Lock()
			results[id] = fmt.Sprintf("val=%s%d%s, ok=%s%v%s",
				magenta, val, reset, red+bold, ok, reset)
			mu.Unlock()
		}(i)
	}

	// All 3 goroutines are now parked in recvq.
	// close(ch) walks recvq and wakes all of them.
	fmt.Printf("  %s%s✖ close(ch) — runtime walks recvq and wakes ALL receivers%s\n", bold, red, reset)
	fmt.Printf("  %s↳ each receiver gets (zero-value, ok=false)%s\n\n", dim, reset)
	close(ch)
	wg.Wait()

	for i, r := range results {
		color := receiverColors[i]
		name := receiverNames[i]
		fmt.Printf("    %s%s%s: %s\n", color+bold, name, reset, r)
	}
	fmt.Printf("\n  %s⚠ All receivers got val=0 (zero value of int) and ok=false%s\n", yellow, reset)
	fmt.Printf("  %s✔ close() is Go's ONLY channel broadcast — no explicit broadcast primitive needed%s\n\n", green, reset)

	// --- Buffered channel: data first, then zero/false ---
	fmt.Printf("%s▸ Buffered Channel — Data First, Then Zero/False%s\n", cyan+bold, reset)
	fmt.Printf("  %sBuffered data is NOT discarded by close() — receivers drain it first%s\n\n", dim, reset)

	buffered := make(chan int, 3)
	buffered <- 10
	buffered <- 20
	close(buffered) // buffer has 2 items, then closed

	for i := 1; i <= 4; i++ {
		val, ok := <-buffered
		if ok {
			fmt.Printf("    %s✔ recv %d:%s val=%s%d%s, ok=%s%v%s  %s← buffered data still available%s\n",
				green, i, reset, magenta, val, reset, green+bold, ok, reset, dim, reset)
		} else {
			fmt.Printf("    %s⚠ recv %d:%s val=%s%d%s, ok=%s%v%s  %s← buffer drained, channel closed%s\n",
				yellow, i, reset, magenta, val, reset, red+bold, ok, reset, dim, reset)
		}
	}
	fmt.Printf("\n  %s✔ close() signals EOF after buffer empties — never discards in-flight data%s\n", green, reset)
}
