package main

import (
	"context"
	"fmt"
	"time"
)

// ============================================================
// Timeout with Context — The Correct Timeout Pattern
// ============================================================
//
// The Problem:
//   You need to abandon a slow operation after a deadline. Two approaches
//   exist: time.After and context.WithTimeout. They look similar but have
//   very different resource implications, especially inside loops.
//
// Real-world example:
//   An HTTP handler calls a downstream service. If the service doesn't
//   respond in 2 seconds, the handler must return HTTP 504 Gateway Timeout.
//   The handler processes hundreds of requests per second — each needs
//   its own timeout, and leaked timers would accumulate in memory.
//
// The Anti-Pattern (time.After in a loop):
//
//   for {
//       select {
//       case v := <-workCh:
//           process(v)
//       case <-time.After(5 * time.Second):  // ← BUG: leaks timer every iteration
//           return
//       }
//   }
//
//   Every iteration creates a new *time.Timer via runtime.startTimer().
//   The timer is registered in the runtime's timer heap (runtime.timers,
//   a per-P min-heap). Even if the select picks workCh, the timer from
//   time.After lives in the heap until it fires — the GC CANNOT collect
//   it because the runtime holds a reference. At 10k iterations/sec,
//   that's 10k leaked timers per second accumulating in the heap.
//
// The Correct Pattern (context.WithTimeout):
//
//   ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
//   defer cancel()  // ← this calls runtime.stopTimer(), removing the
//                   //   timer from the heap IMMEDIATELY
//
//   select {
//   case v := <-workCh:
//       process(v)
//   case <-ctx.Done():
//       return ctx.Err()
//   }
//
// Why context.WithTimeout is correct:
//   1. Single timer: context.WithTimeout creates ONE timer internally.
//      The cancel() function calls runtime.stopTimer() which removes
//      it from the timer heap immediately — no accumulation.
//
//   2. Propagation: the context carries the deadline to ALL downstream
//      calls. If your handler calls a DB query and an HTTP client, both
//      respect the same deadline automatically via ctx.
//
//   3. Composability: context.WithTimeout returns a child context. The
//      parent's deadline is inherited — if the parent has a 10s deadline
//      and you set a 5s timeout, the child uses 5s. If the parent has
//      a 3s deadline, the child uses 3s (tighter deadline wins).
//
// When time.After IS fine:
//   - Single-use timeout outside a loop (e.g., one select in main())
//   - The timer fires in all paths (no leak because it always completes)
//
// When to use time.NewTimer + Reset:
//   - Repeated timeouts in a for/select loop where you need per-iteration
//     timeout reset (the DualTimeoutWorker pattern from exercises)
//   - Remember: timer.Stop(), drain if needed, then timer.Reset()
//
//   ┌─────────── time.After in loop ──────────┐
//   │ iter 1: timer1 created ──► leaks        │
//   │ iter 2: timer2 created ──► leaks        │
//   │ iter 3: timer3 created ──► leaks        │
//   │ ...1000 iterations = 1000 live timers   │
//   └─────────────────────────────────────────┘
//
//   ┌─────── context.WithTimeout ─────────────┐
//   │ ctx, cancel := WithTimeout(parent, 5s)  │
//   │ defer cancel()                          │
//   │ ONE timer. cancel() removes it.         │
//   │ Zero accumulation. Zero leaks.          │
//   └─────────────────────────────────────────┘

// simulateWork simulates an operation that takes the given duration.
// It respects context cancellation — if ctx is cancelled before the
// work completes, it returns ctx.Err() immediately.
func simulateWork(ctx context.Context, duration time.Duration) (string, error) {
	select {
	case <-time.After(duration):
		return "result-ok", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// timeoutOperation runs an operation with a timeout using context.WithTimeout.
// Returns the result and whether the operation completed or timed out.
func timeoutOperation(parentCtx context.Context, workDuration, timeout time.Duration) (string, error) {
	// Create a child context with timeout — single timer, cancel cleans it up
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel() // ALWAYS defer cancel: stops the timer, releases resources

	result, err := simulateWork(ctx, workDuration)
	if err != nil {
		return "", fmt.Errorf("operation failed: %w", err)
	}
	return result, nil
}

func main() {
	parentCtx := context.Background()

	// Case 1: Fast operation completes within timeout
	fmt.Println("  Case 1: Fast operation (50ms work, 200ms timeout)")
	result, err := timeoutOperation(parentCtx, 50*time.Millisecond, 200*time.Millisecond)
	if err != nil {
		fmt.Printf("    ❌ %v\n", err)
	} else {
		fmt.Printf("    ✅ Completed: %s\n", result)
	}

	// Case 2: Slow operation hits timeout
	fmt.Println("  Case 2: Slow operation (500ms work, 100ms timeout)")
	result, err = timeoutOperation(parentCtx, 500*time.Millisecond, 100*time.Millisecond)
	if err != nil {
		fmt.Printf("    ⏱  Timed out: %v\n", err)
	} else {
		fmt.Printf("    ✅ Completed: %s\n", result)
	}

	// Case 3: Demonstrate that cancel() cleans up immediately
	fmt.Println("  Case 3: Cancel cleans up resources")
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)

	// Cancel immediately — the 5-second timer is removed from the runtime
	// timer heap right now, not after 5 seconds
	cancel()

	// ctx.Done() is already closed — any select on it returns immediately
	select {
	case <-ctx.Done():
		fmt.Printf("    ✅ Context cancelled immediately: %v\n", ctx.Err())
		fmt.Println("    Timer removed from runtime heap — zero resource leak")
	default:
		fmt.Println("    ❌ Should not reach here")
	}
}
