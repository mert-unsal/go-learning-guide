package concepts

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============================================================
// Done / Cancellation Pattern — Legacy and Modern Approaches
// ============================================================
//
// The Problem:
//   You launch N background goroutines (workers, watchers, pollers).
//   When the service shuts down or a request is cancelled, ALL of them
//   must stop promptly. You can't send N values — you may not know N.
//   You need a broadcast signal that wakes every listener at once.
//
// Real-world example:
//   An HTTP server with per-request goroutines: a database query, a
//   cache lookup, and an external API call run concurrently. When the
//   client disconnects (or the 5s deadline hits), all three must abort
//   immediately to free resources.
//
// The Pattern:
//
//   Legacy (pre-context):
//     done := make(chan struct{})
//     // pass done to all goroutines
//     // goroutines select on <-done
//     close(done)  // broadcasts to ALL listeners
//
//   Modern (context.WithCancel):
//     ctx, cancel := context.WithCancel(parent)
//     // pass ctx to all goroutines
//     // goroutines select on <-ctx.Done()
//     cancel()  // closes ctx.Done() internally — same broadcast
//
// Why channels work here:
//   close() is a broadcast. Both patterns use the exact same mechanism:
//   closing a chan struct{} to wake all waiters simultaneously.
//
// Under the hood — context.WithCancel:
//   When you call context.WithCancel(parent), the runtime:
//   1. Creates a cancelCtx struct containing a chan struct{} (the done channel)
//   2. ctx.Done() returns this channel (receive-only)
//   3. cancel() closes the channel — same as close(done)
//   4. Closing propagates: parent cancel → closes child done channels too
//
//   The context package wraps the done-channel pattern with:
//   - Hierarchical cancellation (parent → children)
//   - Deadline/timeout support (WithDeadline, WithTimeout)
//   - Value propagation (WithValue — use sparingly)
//   - Thread-safe cancel() that is idempotent (safe to call multiple times)
//
//   ┌───────────────────────────────────────────────┐
//   │             context.WithCancel                │
//   │                                               │
//   │  parent ctx ──► cancelCtx {                   │
//   │                   done: make(chan struct{})    │
//   │                   children: map[...]          │
//   │                 }                             │
//   │                                               │
//   │  ctx.Done() ──► returns the chan struct{}      │
//   │  cancel()   ──► close(done) + cancel children │
//   │                                               │
//   │  Same mechanism as legacy done channel,       │
//   │  with hierarchy and deadline support.         │
//   └───────────────────────────────────────────────┘

// DemonstrateDoneCancellation shows the legacy done-channel pattern
// and the modern context.WithCancel pattern side by side, proving
// both wake multiple goroutines on cancellation.
func DemonstrateDoneCancellation() {
	// --- Legacy pattern: done channel ---
	fmt.Println("  Legacy done channel pattern:")
	legacyDoneCancellation()

	// --- Modern pattern: context.WithCancel ---
	fmt.Println("  Modern context.WithCancel pattern:")
	contextCancellation()
}

// legacyDoneCancellation demonstrates the pre-context done channel pattern.
func legacyDoneCancellation() {
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Launch 3 workers that watch the done channel
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			select {
			case <-done:
				fmt.Printf("    worker %d: received done signal\n", id)
			case <-time.After(5 * time.Second):
				fmt.Printf("    worker %d: timed out (should not happen)\n", id)
			}
		}(i)
	}

	// Broadcast shutdown to all workers
	close(done)
	wg.Wait()
}

// contextCancellation demonstrates the modern context.WithCancel pattern.
func contextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Launch 3 workers that watch ctx.Done()
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				fmt.Printf("    worker %d: ctx cancelled (%v)\n", id, ctx.Err())
			case <-time.After(5 * time.Second):
				fmt.Printf("    worker %d: timed out (should not happen)\n", id)
			}
		}(i)
	}

	// cancel() closes ctx.Done() — same broadcast mechanism as close(done)
	cancel()
	wg.Wait()
}
