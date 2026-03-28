package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============================================================
// Or-Done Channel — Cancellation-Aware Channel Forwarding
// ============================================================
//
// The Problem:
//   A pipeline stage reads from an input channel and sends to an output
//   channel. You use select to watch ctx.Done() alongside the read —
//   but the WRITE to the output channel is outside the select. If
//   downstream stops reading (crashed, slow, full buffer), the write
//   blocks forever, ignoring ctx.Done(). The goroutine leaks.
//
// Real-world example:
//   A log aggregation pipeline: reader goroutine pulls from a Kafka
//   partition and forwards to a processing channel. The processor
//   crashes. The reader's "out <- msg" blocks forever because nobody
//   reads from out. Even though ctx is cancelled (service shutting down),
//   the reader goroutine hangs — it never checks ctx.Done() on the
//   write side. After hours of restarts, hundreds of leaked goroutines
//   accumulate, exhausting memory.
//
// The Broken Pattern (single select):
//
//   func forward(ctx context.Context, in <-chan int, out chan<- int) {
//       for {
//           select {
//           case <-ctx.Done():
//               return              // ✅ cancels while waiting to READ
//           case v, ok := <-in:
//               if !ok { return }
//               out <- v            // ❌ BUG: blocks forever if out is full
//                                   //    ctx.Done() is NOT checked here
//           }
//       }
//   }
//
// The Fix (double select — the orDone pattern):
//
//   func forward(ctx context.Context, in <-chan int, out chan<- int) {
//       for {
//           select {                       // OUTER select: protects the READ
//           case <-ctx.Done():
//               return
//           case v, ok := <-in:
//               if !ok { return }
//               select {                   // INNER select: protects the WRITE
//               case <-ctx.Done():
//                   return
//               case out <- v:
//                   // sent successfully
//               }
//           }
//       }
//   }
//
// All 4 scenarios traced:
//
//   Scenario 1 — Normal flow (data arrives, downstream ready):
//     Outer select: in has data → receive v
//     Inner select: out is ready → send v
//     Both selects proceed immediately. No blocking.
//
//   Scenario 2 — Cancel while waiting to READ:
//     Outer select: in is empty, ctx.Done() fires → return
//     The outer select catches the cancellation. This works with
//     single-select too — it's not the bug we're fixing.
//
//   Scenario 3 — Cancel while waiting to WRITE (the bug we fixed):
//     Outer select: in has data → receive v
//     Inner select: out is full/nobody reading, ctx.Done() fires → return
//     WITHOUT the inner select, "out <- v" would block forever here.
//     This is the goroutine leak that orDone prevents.
//
//   Scenario 4 — Input channel closes:
//     Outer select: in is closed → ok=false → return
//     Normal pipeline termination. The for-range equivalent.
//
// Visual comparison:
//
//   Single select (BROKEN):
//   ┌────────────────────────────────┐
//   │ select {                       │
//   │   case <-ctx.Done(): ✅        │ ← protects read
//   │   case v := <-in:             │
//   │ }                             │
//   │ out <- v  ← UNPROTECTED ❌     │ ← blocks forever if out is full
//   └────────────────────────────────┘
//
//   Double select (orDone — CORRECT):
//   ┌────────────────────────────────┐
//   │ select {                       │
//   │   case <-ctx.Done(): ✅        │ ← protects read
//   │   case v := <-in:             │
//   │     select {                  │
//   │       case <-ctx.Done(): ✅    │ ← protects write
//   │       case out <- v:   ✅      │
//   │     }                         │
//   │ }                             │
//   └────────────────────────────────┘
//
// Composability:
//   orDone wraps any <-chan T into a cancellation-aware channel.
//   Downstream code can use simple for-range on the returned channel —
//   the cancellation logic is encapsulated inside orDone. This pattern
//   is from "Concurrency in Go" (Katherine Cox-Buday, O'Reilly 2017).
//   It's not in the stdlib because Go gives you primitives (select,
//   channels, context) and you compose them. orDone is a composition.

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

// orDone wraps an input channel with cancellation awareness.
// Every value from in is forwarded to the returned channel.
// If ctx is cancelled, forwarding stops and the output channel is closed.
// The caller can use simple for-range on the returned channel.
func orDone(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select { // OUTER select: protects the READ from in
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return // input closed — normal pipeline termination
				}
				select { // INNER select: protects the WRITE to out
				case <-ctx.Done():
					return
				case out <- v:
					// forwarded successfully
				}
			}
		}
	}()
	return out
}

// brokenForward demonstrates the single-select bug: it protects the read
// but NOT the write. If out is full and nobody reads, this goroutine
// leaks even when ctx is cancelled.
func brokenForward(ctx context.Context, in <-chan int, out chan<- int, leaked *int32, mu *sync.Mutex) {
	for {
		select {
		case <-ctx.Done():
			return // ✅ protects the read
		case v, ok := <-in:
			if !ok {
				return
			}
			// ❌ BUG: this send is outside select — blocks forever
			// if nobody reads from out. ctx.Done() is NOT checked.
			//
			// We add a safety timeout here so the demo doesn't hang,
			// but in real code this would be an infinite block.
			select {
			case out <- v:
			case <-time.After(100 * time.Millisecond):
				mu.Lock()
				*leaked++
				mu.Unlock()
				return // timed out — would be an infinite block in real code
			}
		}
	}
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Or-Done Channel — Cancellation-Aware Forwarding %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ The Double-Select Protection%s\n", cyan+bold, reset)
	fmt.Printf("  %sOUTER select%s: guards the %sREAD%s from input channel\n", yellow+bold, reset, cyan, reset)
	fmt.Printf("  %sINNER select%s: guards the %sWRITE%s to output channel\n", yellow+bold, reset, cyan, reset)
	fmt.Printf("  Both check ctx.Done() — cancellation works at every blocking point\n\n")

	// Scenario 1: Channel that never closes + context cancel → orDone exits cleanly
	fmt.Printf("%s▸ Scenario 1: Cancel with orDone — Clean Exit%s\n", cyan+bold, reset)
	fmt.Printf("  Input channel never closes; cancel() is the only way out\n")
	{
		ctx, cancel := context.WithCancel(context.Background())
		neverCloses := make(chan int) // nobody ever sends or closes this

		wrapped := orDone(ctx, neverCloses)

		// Cancel after a short delay
		go func() {
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("  %s→ cancel() called%s — ctx.Done() channel closes\n", yellow+bold, reset)
			cancel()
		}()

		// for-range exits cleanly when orDone closes the output channel
		count := 0
		for range wrapped {
			count++
		}
		fmt.Printf("  %s✔ orDone exited cleanly — output channel closed, no goroutine leak%s\n", green, reset)
		fmt.Printf("  %s✔ Values received before cancel: %s%d%s (expected 0 — nothing was sent)%s\n", green, magenta, count, green, reset)
		fmt.Printf("  %s✔ OUTER select caught ctx.Done() while waiting to read from neverCloses%s\n", green, reset)
	}

	fmt.Println()

	// Scenario 2: Normal forwarding — values pass through orDone
	fmt.Printf("%s▸ Scenario 2: Normal Forwarding Through orDone%s\n", cyan+bold, reset)
	fmt.Printf("  Buffered input with 5 values; orDone forwards all transparently\n")
	{
		ctx := context.Background()
		in := make(chan int, 5)
		for i := 1; i <= 5; i++ {
			in <- i
		}
		close(in)

		wrapped := orDone(ctx, in)

		var results []int
		for v := range wrapped {
			fmt.Printf("    %s→%s forwarded %s%d%s through orDone\n", green, reset, magenta, v, reset)
			results = append(results, v)
		}
		fmt.Printf("  %s✔ All values forwarded: %v%s\n", green, results, reset)
		fmt.Printf("  %s✔ Input closed → ok=false in OUTER select → orDone returns%s\n", green, reset)
	}

	fmt.Println()

	// Scenario 3: Without orDone — demonstrate the leak problem
	fmt.Printf("%s▸ Scenario 3: Without orDone — Goroutine Leak%s\n", cyan+bold, reset)
	fmt.Printf("  brokenForward protects the READ but %sNOT%s the WRITE\n", red+bold, reset)
	{
		ctx, cancel := context.WithCancel(context.Background())
		in := make(chan int, 3)
		out := make(chan int) // unbuffered — nobody reads from this

		in <- 42
		in <- 43
		in <- 44

		var leaked int32
		var mu sync.Mutex

		// brokenForward will try to write to out, but nobody reads.
		// In real code, this blocks forever. Our demo version times out.
		go brokenForward(ctx, in, out, &leaked, &mu)

		fmt.Printf("  %s→ brokenForward reads %s42%s from in...%s\n", dim, magenta, dim, reset)
		fmt.Printf("  %s→ tries to write to out... but nobody reads from out!%s\n", dim, reset)

		// Cancel context — but brokenForward is stuck on "out <- v",
		// so it can't check ctx.Done()
		time.Sleep(30 * time.Millisecond)
		fmt.Printf("  %s→ cancel() called%s — but brokenForward is stuck on write\n", yellow+bold, reset)
		cancel()

		// Give brokenForward time to hit the safety timeout
		time.Sleep(150 * time.Millisecond)

		mu.Lock()
		leakedCount := leaked
		mu.Unlock()

		if leakedCount > 0 {
			fmt.Printf("  %s⚠ brokenForward blocked on write — ctx.Done() was ignored%s\n", yellow, reset)
			fmt.Printf("  %s⚠ In real code: goroutine leaks forever (no safety timeout)%s\n", yellow, reset)
		}
		fmt.Printf("  %s✔ orDone's INNER select would have caught ctx.Done() here%s\n", green, reset)
	}

	fmt.Printf("\n%s▸ Key Observations%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ orDone encapsulates the double-select so callers use simple for-range%s\n", green, reset)
	fmt.Printf("  %s✔ Without INNER select, a blocked write ignores ctx.Done() forever%s\n", green, reset)
	fmt.Printf("  %s✔ Pattern origin: \"Concurrency in Go\" (Cox-Buday, O'Reilly 2017)%s\n", green, reset)
	fmt.Printf("  %s⚠ Every channel forward in a pipeline needs both selects — audit your code%s\n", yellow, reset)
	fmt.Printf("  %s⚠ This is NOT in stdlib — Go gives primitives; you compose patterns%s\n", yellow, reset)
}
