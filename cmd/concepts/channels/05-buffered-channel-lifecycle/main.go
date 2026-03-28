// Package main contains a standalone conceptual example for the Buffered Channel Lifecycle.
package main

import (
	"fmt"
	"strings"
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
// Buffered Channel Lifecycle — The hchan Ring Buffer Under the Hood
// ============================================================
//
// The Problem:
//   You know that make(chan int, 3) creates a buffered channel, and that
//   sends don't block until the buffer is full. But what is the runtime
//   actually doing on every send and receive? Understanding the hchan
//   struct internals lets you reason about channel performance, memory
//   layout, and when copies happen.
//
// What happens at the runtime level:
//   make(chan int, 3) allocates a runtime.hchan struct with:
//     - buf:      pointer to a circular ring buffer of 3 int-sized slots
//     - dataqsiz: 3 (buffer capacity)
//     - elemsize: 8 (sizeof(int) on 64-bit)
//     - qcount:   0 (number of elements currently in the buffer)
//     - sendx:    0 (next write index into buf)
//     - recvx:    0 (next read index from buf)
//     - sendq:    empty sudog list (blocked senders)
//     - recvq:    empty sudog list (blocked receivers)
//     - lock:     mutex protecting all fields
//
// The ring buffer:
//   buf is a contiguous array of dataqsiz elements. sendx and recvx
//   are indices that wrap around using modular arithmetic:
//     sendx = (sendx + 1) % dataqsiz
//     recvx = (recvx + 1) % dataqsiz
//
//   This gives FIFO ordering without shifting elements. The buffer
//   is "full" when qcount == dataqsiz, "empty" when qcount == 0.
//
//   ┌─────────────── hchan ───────────────┐
//   │                                     │
//   │  buf ──► [ slot0 | slot1 | slot2 ]  │  ← circular ring buffer
//   │            ▲                        │
//   │           sendx (next write)        │
//   │           recvx (next read)         │
//   │                                     │
//   │  qcount: items currently buffered   │
//   │  sendq:  blocked senders  (empty)   │
//   │  recvq:  blocked receivers (empty)  │
//   └─────────────────────────────────────┘
//
// Send path (when buffer has space, qcount < dataqsiz):
//   1. Lock hchan.lock
//   2. typedmemmove: copy sender's value into buf[sendx]
//   3. sendx = (sendx + 1) % dataqsiz
//   4. qcount++
//   5. Unlock — sender returns immediately (no blocking)
//
// Receive path (when buffer has data, qcount > 0):
//   1. Lock hchan.lock
//   2. typedmemmove: copy buf[recvx] into receiver's variable
//   3. recvx = (recvx + 1) % dataqsiz
//   4. qcount--
//   5. Unlock — receiver returns immediately (no blocking)
//
// Key insight:
//   Every send and receive through the buffer involves a memory copy
//   (typedmemmove). Go channels transfer ownership by copying values,
//   not by sharing pointers. This is what makes channels safe: once
//   you send a value, the sender's copy and the buffer's copy are
//   independent.

// bufferState tracks the conceptual hchan ring buffer state.
// We can't read the real hchan fields (they're unexported runtime
// internals), so we simulate the indices to show what the runtime does.
type bufferState struct {
	cap    int
	qcount int
	sendx  int
	recvx  int
	buf    []int // mirrors the ring buffer contents
}

func newBufferState(capacity int) *bufferState {
	return &bufferState{
		cap: capacity,
		buf: make([]int, capacity),
	}
}

func (s *bufferState) send(val int) {
	s.buf[s.sendx] = val
	s.sendx = (s.sendx + 1) % s.cap
	s.qcount++
}

func (s *bufferState) recv() int {
	val := s.buf[s.recvx]
	s.buf[s.recvx] = 0 // slot cleared after read
	s.recvx = (s.recvx + 1) % s.cap
	s.qcount--
	return val
}

func (s *bufferState) visualize() string {
	parts := make([]string, s.cap)
	for i := 0; i < s.cap; i++ {
		if s.buf[i] != 0 {
			parts[i] = fmt.Sprintf("%s%d%s", magenta, s.buf[i], reset)
		} else {
			parts[i] = fmt.Sprintf("%s__%s", dim, reset)
		}
	}
	return "[" + strings.Join(parts, "|") + "]"
}

func (s *bufferState) print(label string) {
	fill := fmt.Sprintf("Buffer: %s%d/%d%s slots used", yellow, s.qcount, s.cap, reset)
	fmt.Printf("  %s%-28s%s  qcount=%s%d%s  sendx=%s%d%s  recvx=%s%d%s  %s  %s\n",
		cyan, label, reset,
		yellow, s.qcount, reset,
		cyan, s.sendx, reset,
		cyan, s.recvx, reset,
		s.visualize(),
		fill)
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Buffered Channel Lifecycle (hchan Ring Buffer)  %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Step 1: make(chan int, 3) ---
	// Runtime allocates hchan + contiguous buf of 3 int slots.
	fmt.Printf("%s▸ Step 1: Channel Creation%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ make(chan int, 3) allocates hchan + contiguous ring buffer of 3 int-sized slots%s\n", green, reset)
	ch := make(chan int, 3)
	state := newBufferState(3)
	state.print("make(chan int, 3)")
	fmt.Println()

	// --- Step 2: send 3 values (fills the buffer) ---
	// Each send copies into buf[sendx], advances sendx, increments qcount.
	// No blocking because qcount < dataqsiz at each step.
	fmt.Printf("%s▸ Step 2: Fill the Buffer (3 sends)%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Each send: typedmemmove into buf[sendx], then sendx = (sendx+1) %% cap%s\n", green, reset)
	fmt.Printf("  %s✔ No blocking — qcount < dataqsiz at each step%s\n", green, reset)
	for _, v := range []int{10, 20, 30} {
		ch <- v
		state.send(v)
		state.print(fmt.Sprintf("send %d", v))
	}
	// sendx has wrapped back to 0: (0→1→2→0)
	fmt.Printf("  %s⚠ sendx wrapped back to 0 (0→1→2→0) — this is the circular ring in action%s\n", yellow, reset)
	fmt.Println()

	// --- Step 3: receive 2 values ---
	// Each receive copies from buf[recvx], advances recvx, decrements qcount.
	fmt.Printf("%s▸ Step 3: Receive 2 Values%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Each recv: typedmemmove from buf[recvx], then recvx = (recvx+1) %% cap%s\n", green, reset)
	for range 2 {
		v := <-ch
		state.recv()
		state.print(fmt.Sprintf("recv → %d", v))
	}
	// recvx is now at 2, qcount is 1, one value (30) still in buf[2]
	fmt.Printf("  %s✔ recvx is now at 2, qcount is 1 — value 30 still in buf[2]%s\n", green, reset)
	fmt.Println()

	// --- Step 4: send 1 more value (sendx wraps around) ---
	// sendx is currently 0 (wrapped after step 2). The new value goes
	// into buf[0], demonstrating the circular nature of the ring buffer.
	fmt.Printf("%s▸ Step 4: Send Again (sendx wraps around)%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ sendx is 0 (already wrapped) — new value goes into buf[0]%s\n", green, reset)
	ch <- 40
	state.send(40)
	state.print("send 40 (sendx wraps)")
	fmt.Printf("  %s⚠ Old slot 0 was cleared by recv — ring buffer reuses it without reallocation%s\n", yellow, reset)
	fmt.Println()

	// --- Step 5: drain remaining ---
	fmt.Printf("%s▸ Step 5: Drain Remaining%s\n", cyan+bold, reset)
	for range 2 {
		v := <-ch
		state.recv()
		state.print(fmt.Sprintf("recv → %d (drain)", v))
	}

	fmt.Println()
	fmt.Printf("  %s✔ Key takeaway: sendx and recvx chase each other around%s\n", green+bold, reset)
	fmt.Printf("  %s  the ring buffer. No element shifting, no reallocation.%s\n", green, reset)
	fmt.Printf("  %s  FIFO order is maintained purely by index arithmetic.%s\n", green, reset)
}
