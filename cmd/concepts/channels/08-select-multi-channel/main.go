// Package main contains a standalone conceptual example for the Select Multi-Channel pattern.
package main

import (
	"fmt"
	"time"
)

// ============================================================
// Select Multi-Channel — How selectgo() Watches N Channels
// ============================================================
//
// The Problem:
//   A goroutine needs to wait on multiple channels simultaneously.
//   select is the language primitive for this, but what does the runtime
//   actually do? Understanding selectgo() explains why select is
//   efficient, why case order doesn't matter for fairness, and how
//   a single goroutine can be on multiple wait queues at once.
//
// What happens at the runtime level (runtime/select.go → selectgo):
//
//   1. LOCK PHASE — Lock all channels in the select.
//      Channels are sorted by their memory address before locking.
//      This prevents deadlocks when two goroutines select on the
//      same set of channels in different order (lock ordering).
//
//   2. POLL PHASE — Check if any case can proceed immediately.
//      Walk each case: for send cases, check if buffer has space or
//      a receiver is waiting. For recv cases, check if buffer has data
//      or a sender is waiting. If a ready case is found, execute it
//      and return (fast path — no parking needed).
//      If multiple cases are ready, pick one at random (see below).
//
//   3. ENQUEUE PHASE — No case ready, goroutine must wait.
//      Create one sudog per case. Each sudog points to the same
//      goroutine (sg.g = current G) but different channels.
//      Enqueue each sudog on its channel's sendq or recvq.
//      The goroutine is now on N wait queues simultaneously.
//
//   4. PARK — gopark(). Goroutine enters _Gwaiting.
//      It sleeps until ANY channel wakes it.
//
//   5. WAKE — Some channel becomes ready.
//      That channel dequeues our sudog, performs the send/recv,
//      and calls goready() on our goroutine.
//
//   6. CLEANUP — Dequeue sudogs from ALL other channels.
//      The woken goroutine walks its list of sudogs and removes each
//      one from the other channels' wait queues. This is essential:
//      we must not leave stale sudogs on channels we didn't use.
//      All sudogs are returned to the per-P cache.
//
//   7. UNLOCK ALL — Release all channel locks and return the
//      winning case index.
//
//   ┌──────────────────────────────────────────────────────────────┐
//   │  select {                                                    │
//   │    case v := <-ch1:    // case 0                             │
//   │    case v := <-ch2:    // case 1                             │
//   │    case v := <-ch3:    // case 2                             │
//   │  }                                                           │
//   │                                                              │
//   │  Goroutine G parks on ALL three channels:                    │
//   │                                                              │
//   │  ch1.recvq: [ ... | sudog{g=G, case=0} ]                    │
//   │  ch2.recvq: [ ... | sudog{g=G, case=1} ]                    │
//   │  ch3.recvq: [ ... | sudog{g=G, case=2} ]                    │
//   │                                                              │
//   │  When ch2 gets a sender:                                     │
//   │    → dequeue sudog from ch2.recvq → goready(G)              │
//   │    → G wakes, removes sudogs from ch1.recvq and ch3.recvq   │
//   │    → selectgo returns case index 1                           │
//   └──────────────────────────────────────────────────────────────┘
//
// Random selection for fairness:
//   When selectgo finds multiple ready cases in the poll phase, it
//   picks one uniformly at random (using fastrandn). This prevents
//   starvation: if cases were evaluated in source order, the first
//   case would always win when multiple channels are ready. The
//   randomness ensures all channels get fair service over time.
//
// Lock ordering to prevent deadlock:
//   If goroutine A selects on {ch1, ch2} and goroutine B selects on
//   {ch2, ch1}, locking in source order would deadlock (A locks ch1,
//   B locks ch2, both wait for the other). Sorting by address gives
//   a globally consistent lock order, eliminating this class of deadlock.
//
// Performance:
//   A select with N cases creates N sudogs when it must park. Each
//   sudog is ~96 bytes. For a select with 3 cases that parks, that's
//   ~288 bytes from the per-P sudog cache. This is cheap, but a
//   select in a tight loop with many cases does generate GC work
//   from sudog allocation/deallocation. In hot paths, prefer fewer
//   select cases.

func main() {
	fmt.Println("=== Select Multi-Channel (selectgo internals) ===")
	fmt.Println()

	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)

	result := make(chan string, 1)

	// Launch a goroutine that enters select on all three channels.
	// At the runtime level:
	//   1. selectgo locks ch1, ch2, ch3 (sorted by address)
	//   2. Poll phase: none have data → must park
	//   3. Creates 3 sudogs, enqueues on ch1.recvq, ch2.recvq, ch3.recvq
	//   4. gopark() → goroutine sleeps on all three wait queues
	go func() {
		fmt.Println("  [select goroutine] entering select on 3 channels...")
		start := time.Now()

		select {
		case v := <-ch1:
			result <- fmt.Sprintf("ch1 won with %q after %v", v, time.Since(start).Round(time.Millisecond))
		case v := <-ch2:
			result <- fmt.Sprintf("ch2 won with %q after %v", v, time.Since(start).Round(time.Millisecond))
		case v := <-ch3:
			result <- fmt.Sprintf("ch3 won with %q after %v", v, time.Since(start).Round(time.Millisecond))
		}
		// After waking:
		//   - The winning channel's sudog was dequeued by the sender
		//   - This goroutine removes its sudogs from the other two channels
		//   - All 3 sudogs returned to per-P cache
	}()

	// Let the select goroutine park on all three channels.
	time.Sleep(50 * time.Millisecond)

	// Send to ch2 only. This triggers:
	//   1. ch2's chansend finds our sudog in ch2.recvq
	//   2. Direct copy (unbuffered) from sender to select goroutine
	//   3. goready() wakes the select goroutine
	//   4. Select goroutine cleans up sudogs from ch1 and ch3
	fmt.Println("  [main] sending to ch2...")
	ch2 <- "payload-from-ch2"

	winner := <-result
	fmt.Printf("  Result: %s\n", winner)
	fmt.Println()

	// Demonstrate random fairness: when multiple channels are ready
	// at the poll phase, selectgo picks at random.
	fmt.Println("  --- Random fairness demonstration ---")
	wins := map[string]int{"ch1": 0, "ch2": 0, "ch3": 0}
	const trials = 900

	for range trials {
		a := make(chan int, 1)
		b := make(chan int, 1)
		c := make(chan int, 1)

		// All three are ready before select enters.
		// selectgo's poll phase will find all three ready
		// and pick one with fastrandn().
		a <- 1
		b <- 2
		c <- 3

		select {
		case <-a:
			wins["ch1"]++
		case <-b:
			wins["ch2"]++
		case <-c:
			wins["ch3"]++
		}
	}

	fmt.Printf("  Over %d trials with all 3 ready:\n", trials)
	for name, count := range wins {
		pct := float64(count) / float64(trials) * 100
		fmt.Printf("    %s: %d wins (%.1f%%)\n", name, count, pct)
	}
	fmt.Println("  Each should be ~33% — select uses fastrandn for fairness.")
}
