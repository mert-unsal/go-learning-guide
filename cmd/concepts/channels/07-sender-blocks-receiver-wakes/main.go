// Package main contains a standalone conceptual example for the Sender Blocks, Receiver Wakes pattern.
package main

import (
	"fmt"
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
// Sender Blocks, Receiver Wakes — The sudog Lifecycle
// ============================================================
//
// The Problem:
//   When a buffered channel is full and a sender tries to write, what
//   exactly happens? The goroutine doesn't spin-wait or busy-loop.
//   The runtime parks it using a sudog — a "sleeping goroutine
//   descriptor" — and only wakes it when buffer space becomes available.
//   Understanding this mechanism explains channel backpressure, memory
//   cost of blocked goroutines, and why channels are efficient even
//   with millions of goroutines.
//
// The sudog lifecycle:
//   sudog is defined in runtime/runtime2.go. Key fields:
//     - g:       pointer to the blocked goroutine (runtime.g)
//     - elem:    unsafe.Pointer to the value being sent/received
//     - c:       pointer to the channel (runtime.hchan)
//     - next:    linked list pointer (sendq/recvq are FIFO lists)
//
//   Full sequence when buffer is full:
//
//   1. Sender calls chansend(). Acquires hchan.lock.
//   2. Checks: qcount == dataqsiz → buffer is full, cannot proceed.
//   3. runtime.acquireSudog(): get a sudog from per-P cache (or alloc).
//   4. Fill sudog: sg.g = current goroutine, sg.elem = &value, sg.c = ch.
//   5. Enqueue sg onto hchan.sendq (tail of FIFO linked list).
//   6. runtime.gopark(): change goroutine state from _Grunning to _Gwaiting.
//      The goroutine is removed from P's run queue. It consumes no CPU.
//   7. Release hchan.lock. The M (OS thread) picks up another G from
//      the run queue.
//
//   ... time passes, sender is sleeping ...
//
//   8. A receiver calls chanrecv(). Acquires hchan.lock.
//   9. Copies from buf[recvx] into receiver's variable. recvx++, qcount--.
//  10. Checks sendq: finds our sender's sudog at the head.
//  11. Dequeues the sudog. Copies the sender's elem value into the
//      freed buffer slot (buf[sendx]). sendx++, qcount++.
//  12. runtime.goready(sg.g): change goroutine state to _Grunnable,
//      put it back on a P's run queue.
//  13. runtime.releaseSudog(sg): return sudog to per-P cache for reuse.
//  14. Release hchan.lock. Sender goroutine will resume when scheduled.
//
//   ┌────────────────────────────────────────────────────────────┐
//   │  Buffered channel (cap=2), buffer full                    │
//   │                                                           │
//   │  buf: [ val_A | val_B ]    ← full (qcount == dataqsiz)   │
//   │                                                           │
//   │  Sender S3 tries to send val_C:                           │
//   │    → buffer full → create sudog → park on sendq           │
//   │                                                           │
//   │  sendq: [ sudog{g=S3, elem=&val_C} ]                     │
//   │  recvq: [ ]                                               │
//   │                                                           │
//   │  Receiver R1 arrives:                                     │
//   │    → recv val_A from buf[0] → dequeue S3's sudog          │
//   │    → copy val_C into freed slot → goready(S3)             │
//   │                                                           │
//   │  buf: [ val_C | val_B ]    ← S3's value is now in buffer  │
//   │  sendq: [ ]                ← S3 is awake and runnable     │
//   └────────────────────────────────────────────────────────────┘
//
// Memory cost of a blocked goroutine:
//   A parked goroutine holds its stack (~2-8 KB) and one sudog (~96 bytes).
//   This is why Go can handle millions of blocked goroutines — they're
//   just data structures in memory, not OS threads consuming kernel
//   resources.
//
// Production relevance:
//   This is exactly how backpressure works in producer-consumer systems.
//   A buffered channel of capacity N allows N items to be in-flight.
//   When the consumer falls behind, producers park via sudogs until
//   the consumer catches up. No data is lost, no explicit flow control
//   code is needed — the channel handles it.

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Sender Blocks, Receiver Wakes (sudog Lifecycle) %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	// Buffered channel with capacity 2.
	ch := make(chan string, 2)

	// Fill the buffer — these two sends return immediately.
	fmt.Printf("%s▸ Step 1: Fill the Buffer%s\n", cyan+bold, reset)
	ch <- "val_A"
	ch <- "val_B"
	fmt.Printf("  %s✔ Two sends returned immediately — buffer had space%s\n", green, reset)
	fmt.Printf("  Buffer state: %s[%sval_A%s|%sval_B%s]%s  (qcount=%s2%s, cap=%s2%s)\n",
		dim, magenta, dim, magenta, dim, reset, yellow, reset, yellow, reset)
	fmt.Printf("  %s⚠ Buffer is now FULL — next send will block the goroutine%s\n", yellow, reset)
	fmt.Println()

	blocked := make(chan struct{})
	unblocked := make(chan struct{})
	var blockStart, blockEnd time.Time

	// Launch a sender goroutine that will block because the buffer is full.
	// At the runtime level:
	//   1. acquireSudog() → fill sg.g, sg.elem = &"val_C"
	//   2. enqueue sg on hchan.sendq
	//   3. gopark() → goroutine enters _Gwaiting
	fmt.Printf("%s▸ Step 2: Sender Blocks on Full Buffer%s\n", cyan+bold, reset)
	go func() {
		close(blocked) // signal: we're about to block
		blockStart = time.Now()
		fmt.Printf("  [%sSENDER%s]  ch <- \"val_C\" — buffer full, creating sudog and parking on sendq\n", yellow+bold, reset)
		ch <- "val_C" // BLOCKS: buffer is full, sender parked on sendq
		blockEnd = time.Now()
		fmt.Printf("  [%sSENDER%s]  woke up! goready() moved goroutine back to %s_Grunnable%s\n", yellow+bold, reset, magenta, reset)
		close(unblocked)
	}()

	<-blocked

	// Prove the sender is blocked by sleeping. During this time,
	// the sender's goroutine is in _Gwaiting state, consuming zero CPU.
	delay := 100 * time.Millisecond
	fmt.Printf("  [%sSENDER%s]  goroutine parked in %s_Gwaiting%s on sendq — zero CPU consumed\n", yellow+bold, reset, magenta, reset)
	fmt.Printf("  %s✔ Sleeping %v to prove sender is truly blocked...%s\n", green, delay, reset)
	time.Sleep(delay)
	fmt.Println()

	// Receive one value — this triggers the wake sequence:
	//   1. recv val_A from buf[recvx]
	//   2. dequeue sender's sudog from sendq
	//   3. copy sender's val_C into freed buffer slot
	//   4. goready(sender) → _Grunnable
	//   5. releaseSudog() → return to per-P cache
	fmt.Printf("%s▸ Step 3: Receiver Triggers the Wake Sequence%s\n", cyan+bold, reset)
	v1 := <-ch
	fmt.Printf("  [%sRECEIVER%s] received: %s%q%s → freed one buffer slot\n", cyan+bold, reset, magenta, v1, reset)
	fmt.Printf("  %s✔ Runtime: dequeued sudog from sendq%s\n", green, reset)
	fmt.Printf("  %s✔ Runtime: copied val_C into freed buffer slot via typedmemmove%s\n", green, reset)
	fmt.Printf("  %s✔ Runtime: goready(sender) → sender is now _Grunnable%s\n", green, reset)
	fmt.Printf("  %s✔ Runtime: releaseSudog() → sudog returned to per-P cache%s\n", green, reset)
	fmt.Printf("  Buffer state: %s[%sval_C%s|%sval_B%s]%s  (val_C filled the freed slot)\n",
		dim, magenta, dim, magenta, dim, reset)

	<-unblocked

	blockedDuration := blockEnd.Sub(blockStart)
	fmt.Println()
	fmt.Printf("%s▸ Timing Analysis%s\n", cyan+bold, reset)
	fmt.Printf("  Sender was blocked for: %s%v%s\n", magenta, blockedDuration.Round(time.Millisecond), reset)
	fmt.Printf("  %s⚠ ≈100ms matches our sleep — proof that sender was truly parked, not spinning%s\n", yellow, reset)
	fmt.Println()

	// Drain remaining values to show FIFO order is preserved.
	fmt.Printf("%s▸ Step 4: Drain Remaining Values%s\n", cyan+bold, reset)
	v2 := <-ch
	v3 := <-ch
	fmt.Printf("  Received: %s%q%s, %s%q%s\n", magenta, v2, reset, magenta, v3, reset)
	fmt.Println()

	fmt.Printf("  %s✔ FIFO order preserved: %sval_A%s %s→%s %sval_B%s %s→%s %sval_C%s\n",
		green+bold, magenta, green+bold, dim, green+bold, magenta, green+bold, dim, green+bold, magenta, reset)
	fmt.Printf("  %s✔ val_C was in the sender's sudog.elem, not in the buffer,%s\n", green, reset)
	fmt.Printf("  %s  until the receiver freed a slot and the runtime copied it in.%s\n", green, reset)
}
