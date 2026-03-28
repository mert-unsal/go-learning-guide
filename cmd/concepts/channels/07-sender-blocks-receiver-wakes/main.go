// Package main contains a standalone conceptual example for the Sender Blocks, Receiver Wakes pattern.
package main

import (
	"fmt"
	"time"
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
	fmt.Println("=== Sender Blocks, Receiver Wakes (sudog lifecycle) ===")
	fmt.Println()

	// Buffered channel with capacity 2.
	ch := make(chan string, 2)

	// Fill the buffer — these two sends return immediately.
	ch <- "val_A"
	ch <- "val_B"
	fmt.Println("  Buffer filled: [val_A | val_B]  (qcount=2, cap=2)")

	blocked := make(chan struct{})
	unblocked := make(chan struct{})
	var blockStart, blockEnd time.Time

	// Launch a sender goroutine that will block because the buffer is full.
	// At the runtime level:
	//   1. acquireSudog() → fill sg.g, sg.elem = &"val_C"
	//   2. enqueue sg on hchan.sendq
	//   3. gopark() → goroutine enters _Gwaiting
	go func() {
		close(blocked) // signal: we're about to block
		blockStart = time.Now()
		ch <- "val_C" // BLOCKS: buffer is full, sender parked on sendq
		blockEnd = time.Now()
		close(unblocked)
	}()

	<-blocked

	// Prove the sender is blocked by sleeping. During this time,
	// the sender's goroutine is in _Gwaiting state, consuming zero CPU.
	delay := 100 * time.Millisecond
	fmt.Printf("  Sender goroutine blocked (parked on sendq)...\n")
	fmt.Printf("  Sleeping %v to prove it...\n", delay)
	time.Sleep(delay)

	// Receive one value — this triggers the wake sequence:
	//   1. recv val_A from buf[recvx]
	//   2. dequeue sender's sudog from sendq
	//   3. copy sender's val_C into freed buffer slot
	//   4. goready(sender) → _Grunnable
	//   5. releaseSudog() → return to per-P cache
	v1 := <-ch
	fmt.Printf("  Received: %q → freed one buffer slot\n", v1)
	fmt.Println("  Runtime: dequeued sudog, copied val_C into buffer, woke sender")

	<-unblocked

	blockedDuration := blockEnd.Sub(blockStart)
	fmt.Printf("  Sender was blocked for: %v\n", blockedDuration.Round(time.Millisecond))

	// Drain remaining values to show FIFO order is preserved.
	v2 := <-ch
	v3 := <-ch
	fmt.Printf("  Remaining drain: %q, %q\n", v2, v3)
	fmt.Println()
	fmt.Println("  FIFO order: val_A → val_B → val_C")
	fmt.Println("  val_C was in the sender's sudog.elem, not in the buffer,")
	fmt.Println("  until the receiver freed a slot and the runtime copied it in.")
}
