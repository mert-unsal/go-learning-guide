// Package main contains a standalone conceptual example for the Unbuffered Direct Transfer pattern.
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
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Unbuffered Direct Transfer (Goroutine Rendezvous)%s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	ch := make(chan int) // dataqsiz == 0, buf == nil

	senderReady := make(chan struct{})
	senderDone := make(chan struct{})

	var sendStart, sendEnd time.Time

	fmt.Printf("%s▸ Channel Created%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ make(chan int) → hchan with dataqsiz=0, buf=nil — no ring buffer at all%s\n", green, reset)
	fmt.Printf("  %s✔ Transfers use sendDirect(): one memmove from sender's stack to receiver's stack%s\n", green, reset)
	fmt.Println()

	// Sender goroutine — will block until main goroutine receives.
	fmt.Printf("%s▸ Launching Sender Goroutine%s\n", cyan+bold, reset)
	go func() {
		close(senderReady) // signal that we're about to send
		sendStart = time.Now()
		fmt.Printf("  [%sSENDER%s]  ch <- 42  — no receiver yet, goroutine parks in %s_Gwaiting%s\n", yellow+bold, reset, magenta, reset)
		ch <- 42 // BLOCKS here: no buffer, no receiver yet
		sendEnd = time.Now()
		fmt.Printf("  [%sSENDER%s]  unblocked! sendDirect() completed — value delivered\n", yellow+bold, reset)
		close(senderDone)
	}()

	<-senderReady // wait for sender goroutine to be running

	// Deliberately delay the receive to prove the sender is blocked.
	// During this 100ms, the sender's goroutine is in _Gwaiting state,
	// parked on hchan.sendq. It consumes zero CPU.
	delay := 100 * time.Millisecond
	fmt.Printf("  [%sRECEIVER%s] sleeping %s%v%s to prove sender blocks...\n", cyan+bold, reset, magenta, delay, reset)
	fmt.Printf("  %s⚠ During this sleep, sender is parked on hchan.sendq — zero CPU consumed%s\n", yellow, reset)
	time.Sleep(delay)

	fmt.Println()
	fmt.Printf("%s▸ Receive Operation%s\n", cyan+bold, reset)
	recvStart := time.Now()
	val := <-ch // runtime.sendDirect: memmove from sender's stack to ours
	recvEnd := time.Now()

	<-senderDone

	senderBlocked := sendEnd.Sub(sendStart)
	recvLatency := recvEnd.Sub(recvStart)

	fmt.Printf("  [%sRECEIVER%s] received value: %s%d%s\n", cyan+bold, reset, magenta, val, reset)
	fmt.Println()

	fmt.Printf("%s▸ Timing Analysis%s\n", cyan+bold, reset)
	fmt.Printf("  Sender was blocked for:    %s%v%s (≈ receiver's sleep)\n", magenta, senderBlocked.Round(time.Millisecond), reset)
	fmt.Printf("  Receive operation latency: %s%v%s (sub-microsecond = direct copy)\n", magenta, recvLatency, reset)
	fmt.Println()

	fmt.Printf("  %s✔ The sender blocked ~100ms — it was parked in _Gwaiting%s\n", green+bold, reset)
	fmt.Printf("  %s  until our receive call triggered sendDirect().%s\n", green, reset)
	fmt.Printf("  %s✔ The receive itself was nearly instant: one memmove,%s\n", green+bold, reset)
	fmt.Printf("  %s  no buffer involved. This is the fastest goroutine-to-goroutine transfer.%s\n", green, reset)
}
