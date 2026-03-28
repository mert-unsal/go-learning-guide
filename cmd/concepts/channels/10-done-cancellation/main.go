package main

import (
	"context"
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

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Done / Cancellation — Legacy vs Modern Context  %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Legacy Pattern: done = make(chan struct{})%s\n", cyan+bold, reset)
	fmt.Printf("  %sPre-context era — close(done) broadcasts to all goroutines%s\n\n", dim, reset)
	legacyDoneCancellation()

	fmt.Println()
	fmt.Printf("%s▸ Modern Pattern: context.WithCancel%s\n", cyan+bold, reset)
	fmt.Printf("  %scancel() closes ctx.Done() internally — same broadcast mechanism%s\n", dim, reset)
	fmt.Printf("  %s✔ Adds hierarchy, deadlines, and idempotent cancel()%s\n\n", green, reset)
	contextCancellation()

	fmt.Println()
	fmt.Printf("%s▸ Comparison%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Both use the same underlying mechanism: close(chan struct{})%s\n", green, reset)
	fmt.Printf("  %s⚠ Legacy: manual, no hierarchy, no deadlines, easy to leak%s\n", yellow, reset)
	fmt.Printf("  %s✔ Modern: hierarchical cancel, WithTimeout/WithDeadline, idempotent%s\n", green, reset)
	fmt.Printf("  %s✔ Always prefer context in new code — it's the standard API boundary%s\n", green, reset)
}

// legacyDoneCancellation demonstrates the pre-context done channel pattern.
func legacyDoneCancellation() {
	workerColors := []string{cyan, yellow, magenta}
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Launch 3 workers that watch the done channel
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			select {
			case <-done:
				color := workerColors[id]
				fmt.Printf("    %s%s[worker-%d]%s received done signal %s← close(done) woke this goroutine%s\n",
					color, bold, id, reset, dim, reset)
			case <-time.After(5 * time.Second):
				fmt.Printf("    %s⚠ worker %d: timed out (should not happen)%s\n", red, id, reset)
			}
		}(i)
	}

	// Broadcast shutdown to all workers
	fmt.Printf("  %s%s✖ close(done)%s — broadcasting to all workers\n", bold, red, reset)
	close(done)
	wg.Wait()
	fmt.Printf("  %s✔ All 3 workers received the done signal simultaneously%s\n", green, reset)
}

// contextCancellation demonstrates the modern context.WithCancel pattern.
func contextCancellation() {
	workerColors := []string{cyan, yellow, magenta}
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Launch 3 workers that watch ctx.Done()
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				color := workerColors[id]
				fmt.Printf("    %s%s[worker-%d]%s ctx.Err() = %s%v%s %s← cancel() closed ctx.Done()%s\n",
					color, bold, id, reset, magenta, ctx.Err(), reset, dim, reset)
			case <-time.After(5 * time.Second):
				fmt.Printf("    %s⚠ worker %d: timed out (should not happen)%s\n", red, id, reset)
			}
		}(i)
	}

	// cancel() closes ctx.Done() — same broadcast mechanism as close(done)
	fmt.Printf("  %s%s✖ cancel()%s — closes ctx.Done() channel internally\n", bold, red, reset)
	cancel()
	wg.Wait()
	fmt.Printf("  %s✔ All 3 workers detected cancellation via ctx.Done()%s\n", green, reset)
}
