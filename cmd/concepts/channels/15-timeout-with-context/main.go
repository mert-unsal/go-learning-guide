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
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Timeout with Context — Correct Timeout Pattern  %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	parentCtx := context.Background()

	// Case 1: Fast operation completes within timeout
	fmt.Printf("%s▸ Case 1: Fast Operation (work < timeout)%s\n", cyan+bold, reset)
	fmt.Printf("  Work: %s50ms%s  Timeout: %s200ms%s — work finishes before deadline\n", magenta, reset, magenta, reset)
	start := time.Now()
	result, err := timeoutOperation(parentCtx, 50*time.Millisecond, 200*time.Millisecond)
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("  %s✘ %v%s\n", red, err, reset)
	} else {
		fmt.Printf("  %s✔ Completed: %s%s (took %s%v%s)\n", green, result, reset, magenta, elapsed.Round(time.Millisecond), reset)
		fmt.Printf("  %s✔ time.After(50ms) fired first → select picked work result%s\n", green, reset)
	}
	fmt.Println()

	// Case 2: Slow operation hits timeout
	fmt.Printf("%s▸ Case 2: Slow Operation (work > timeout)%s\n", cyan+bold, reset)
	fmt.Printf("  Work: %s500ms%s  Timeout: %s100ms%s — deadline expires before work completes\n", magenta, reset, magenta, reset)
	start = time.Now()
	result, err = timeoutOperation(parentCtx, 500*time.Millisecond, 100*time.Millisecond)
	elapsed = time.Since(start)
	if err != nil {
		fmt.Printf("  %s✘ Timed out: %v%s (after %s%v%s)\n", red, err, reset, magenta, elapsed.Round(time.Millisecond), reset)
		fmt.Printf("  %s✔ ctx.Done() fired first → select picked cancellation branch%s\n", green, reset)
		fmt.Printf("  %s⚠ The 500ms timer from time.After still lives in the heap until it fires%s\n", yellow, reset)
	} else {
		fmt.Printf("  %s✔ Completed: %s%s\n", green, result, reset)
	}
	fmt.Println()

	// Case 3: Demonstrate that cancel() cleans up immediately
	fmt.Printf("%s▸ Case 3: Immediate Cancel — Resource Cleanup%s\n", cyan+bold, reset)
	fmt.Printf("  Creating context with %s5s%s timeout, then cancelling immediately\n", magenta, reset)
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)

	fmt.Printf("  %sLifecycle:%s create → %scancel()%s → timer removed from runtime heap\n", dim, reset, yellow+bold, reset)

	// Cancel immediately — the 5-second timer is removed from the runtime
	// timer heap right now, not after 5 seconds
	cancel()

	// ctx.Done() is already closed — any select on it returns immediately
	select {
	case <-ctx.Done():
		fmt.Printf("  %s✔ Context cancelled immediately: %v%s\n", green, ctx.Err(), reset)
		fmt.Printf("  %s✔ Timer removed from runtime heap — zero resource leak%s\n", green, reset)
		fmt.Printf("  %s✔ defer cancel() guarantees cleanup even on early return%s\n", green, reset)
	default:
		fmt.Printf("  %s✘ Should not reach here%s\n", red, reset)
	}

	fmt.Printf("\n%s▸ Key Observations%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ context.WithTimeout creates ONE timer; cancel() removes it immediately%s\n", green, reset)
	fmt.Printf("  %s✔ The deadline propagates to ALL downstream calls via ctx%s\n", green, reset)
	fmt.Printf("  %s✔ Tighter parent deadline wins: WithTimeout(3s-parent, 5s) → uses 3s%s\n", green, reset)
	fmt.Printf("  %s⚠ time.After in a loop leaks one timer per iteration — avoid in hot paths%s\n", yellow, reset)
	fmt.Printf("  %s⚠ Always defer cancel() — even if the timeout fires, cancel() is idempotent%s\n", yellow, reset)
}
