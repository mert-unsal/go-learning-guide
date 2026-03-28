package concepts

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
//   ┌────────────────────────────┐
//   │ select {                   │
//   │   case <-ctx.Done(): ✅    │ ← protects read
//   │   case v := <-in:         │
//   │ }                         │
//   │ out <- v  ← UNPROTECTED ❌ │ ← blocks forever if out is full
//   └────────────────────────────┘
//
//   Double select (orDone — CORRECT):
//   ┌────────────────────────────┐
//   │ select {                   │
//   │   case <-ctx.Done(): ✅    │ ← protects read
//   │   case v := <-in:         │
//   │     select {              │
//   │       case <-ctx.Done(): ✅│ ← protects write
//   │       case out <- v:   ✅  │
//   │     }                     │
//   │ }                         │
//   └────────────────────────────┘
//
// Composability:
//   orDone wraps any <-chan T into a cancellation-aware channel.
//   Downstream code can use simple for-range on the returned channel —
//   the cancellation logic is encapsulated inside orDone. This pattern
//   is from "Concurrency in Go" (Katherine Cox-Buday, O'Reilly 2017).
//   It's not in the stdlib because Go gives you primitives (select,
//   channels, context) and you compose them. orDone is a composition.

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

// DemonstrateOrDone shows three scenarios:
//  1. orDone exits cleanly when context is cancelled (no goroutine leak)
//  2. orDone forwards values normally when data flows
//  3. Without orDone, a goroutine leaks when downstream stops reading
func DemonstrateOrDone() {
	// Scenario 1: Channel that never closes + context cancel → orDone exits cleanly
	fmt.Println("  Scenario 1: Cancel with orDone — clean exit")
	{
		ctx, cancel := context.WithCancel(context.Background())
		neverCloses := make(chan int) // nobody ever sends or closes this

		wrapped := orDone(ctx, neverCloses)

		// Cancel after a short delay
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		// for-range exits cleanly when orDone closes the output channel
		count := 0
		for range wrapped {
			count++
		}
		fmt.Println("    ✅ orDone exited cleanly — channel closed, no leak")
		fmt.Printf("    Values received before cancel: %d\n", count)
	}

	// Scenario 2: Normal forwarding — values pass through orDone
	fmt.Println("  Scenario 2: Normal forwarding through orDone")
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
			results = append(results, v)
		}
		fmt.Printf("    ✅ All values forwarded: %v\n", results)
	}

	// Scenario 3: Without orDone — demonstrate the leak problem
	fmt.Println("  Scenario 3: Without orDone — goroutine would leak")
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

		// Cancel context — but brokenForward is stuck on "out <- v",
		// so it can't check ctx.Done()
		time.Sleep(30 * time.Millisecond)
		cancel()

		// Give brokenForward time to hit the safety timeout
		time.Sleep(150 * time.Millisecond)

		mu.Lock()
		leakedCount := leaked
		mu.Unlock()

		if leakedCount > 0 {
			fmt.Println("    ⚠️  brokenForward blocked on write — ctx.Done() was ignored")
			fmt.Println("    In real code: goroutine leaks forever (no safety timeout)")
		}
		fmt.Println("    orDone's inner select would have caught ctx.Done() here")
	}
}
