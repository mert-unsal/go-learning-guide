package main

import (
	"fmt"
	"sync"
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
			results[id] = fmt.Sprintf("receiver %d: val=%d, ok=%v", id, val, ok)
			mu.Unlock()
		}(i)
	}

	// All 3 goroutines are now parked in recvq.
	// close(ch) walks recvq and wakes all of them.
	close(ch)
	wg.Wait()

	fmt.Println("  Unbuffered channel — close() wakes all receivers:")
	for _, r := range results {
		fmt.Printf("    %s\n", r)
	}

	// --- Buffered channel: data first, then zero/false ---
	buffered := make(chan int, 3)
	buffered <- 10
	buffered <- 20
	close(buffered) // buffer has 2 items, then closed

	fmt.Println("  Buffered channel — data first, then zero/false:")
	for i := 1; i <= 4; i++ {
		val, ok := <-buffered
		fmt.Printf("    recv %d: val=%d, ok=%v\n", i, val, ok)
	}
}
