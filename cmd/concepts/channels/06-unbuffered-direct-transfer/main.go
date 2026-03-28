// Package main contains a standalone conceptual example for the Unbuffered Direct Transfer pattern.
package main

import (
	"fmt"
	"time"
)

// ============================================================
// Unbuffered Direct Transfer — Zero-Copy Goroutine Rendezvous
// ============================================================
//
// The Problem:
//   Unbuffered channels are not "buffered channels with capacity 0."
//   They follow a fundamentally different code path in the runtime.
//   Understanding this path explains why unbuffered channels are the
//   fastest way to hand a value between two goroutines — and why
//   sender and receiver must meet at the same instant.
//
// What happens at the runtime level:
//   make(chan int) allocates an hchan with dataqsiz == 0 and buf == nil.
//   There is no ring buffer. The send and receive paths are:
//
//   Case 1 — Receiver arrives first:
//     1. Receiver calls chanrecv(), finds no sender → creates a sudog
//        with a pointer to the receiver's destination variable (elem).
//     2. sudog enqueued on hchan.recvq. Goroutine parks (_Gwaiting).
//     3. Sender calls chansend(), finds waiting receiver in recvq.
//     4. runtime.sendDirect(): memmove from sender's stack variable
//        directly into the receiver's stack variable (the elem pointer
//        in the sudog). ONE copy, no intermediate buffer.
//     5. Sender wakes the receiver goroutine → both continue.
//
//   Case 2 — Sender arrives first:
//     1. Sender calls chansend(), finds no receiver → creates sudog
//        with elem pointing to the value to send.
//     2. sudog enqueued on hchan.sendq. Goroutine parks.
//     3. Receiver calls chanrecv(), finds waiting sender in sendq.
//     4. Copies from sender's elem pointer into receiver's variable.
//     5. Receiver wakes the sender goroutine → both continue.
//
//   ┌──────────────────────────────────────────────────────────┐
//   │          Unbuffered Channel: Direct Transfer             │
//   │                                                         │
//   │  Sender goroutine          Receiver goroutine            │
//   │       │                         │                        │
//   │       ▼                         ▼                        │
//   │   chansend()               chanrecv()                    │
//   │       │                         │                        │
//   │       │    ┌─── hchan ───┐      │                        │
//   │       │    │ buf = nil   │      │                        │
//   │       │    │ sendq / recvq│     │                        │
//   │       │    └─────────────┘      │                        │
//   │       │                         │                        │
//   │       ├──── memmove ───────────►│  (one copy, no buffer) │
//   │       │     sendDirect()        │                        │
//   │       ▼                         ▼                        │
//   │   (continues)              (continues)                   │
//   └──────────────────────────────────────────────────────────┘
//
// Why this matters for performance:
//   - ONE memory copy total (sender → receiver), not two (sender → buf,
//     buf → receiver) like buffered channels.
//   - No ring buffer allocation, so make(chan T) is cheaper than
//     make(chan T, 1) in both memory and setup cost.
//   - The tradeoff: both goroutines must synchronize. The sender is
//     guaranteed to block until a receiver is ready (or vice versa).
//     This makes unbuffered channels a synchronization primitive,
//     not just a data pipe.
//
// Comparison with other languages:
//   - Java's SynchronousQueue is the closest equivalent: zero capacity,
//     handoff blocks until the other side arrives.
//   - Rust's std::sync::mpsc is buffered by default (unbounded).
//   - Go makes unbuffered the default (make(chan T)) — this is a
//     deliberate design choice to encourage synchronization-first thinking.

func main() {
	fmt.Println("=== Unbuffered Direct Transfer (no buffer, one copy) ===")
	fmt.Println()

	ch := make(chan int) // dataqsiz == 0, buf == nil

	senderReady := make(chan struct{})
	senderDone := make(chan struct{})

	var sendStart, sendEnd time.Time

	// Sender goroutine — will block until main goroutine receives.
	go func() {
		close(senderReady) // signal that we're about to send
		sendStart = time.Now()
		ch <- 42 // BLOCKS here: no buffer, no receiver yet
		sendEnd = time.Now()
		close(senderDone)
	}()

	<-senderReady // wait for sender goroutine to be running

	// Deliberately delay the receive to prove the sender is blocked.
	// During this 100ms, the sender's goroutine is in _Gwaiting state,
	// parked on hchan.sendq. It consumes zero CPU.
	delay := 100 * time.Millisecond
	fmt.Printf("  Receiver sleeping %v to prove sender blocks...\n", delay)
	time.Sleep(delay)

	recvStart := time.Now()
	val := <-ch // runtime.sendDirect: memmove from sender's stack to ours
	recvEnd := time.Now()

	<-senderDone

	senderBlocked := sendEnd.Sub(sendStart)
	recvLatency := recvEnd.Sub(recvStart)

	fmt.Printf("  Received value: %d\n", val)
	fmt.Printf("  Sender was blocked for:    %v (≈ receiver's sleep)\n", senderBlocked.Round(time.Millisecond))
	fmt.Printf("  Receive operation latency: %v (sub-microsecond = direct copy)\n", recvLatency)
	fmt.Println()
	fmt.Println("  The sender blocked ~100ms — it was parked in _Gwaiting")
	fmt.Println("  until our receive call triggered sendDirect().")
	fmt.Println("  The receive itself was nearly instant: one memmove,")
	fmt.Println("  no buffer involved.")
}
