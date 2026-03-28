# Deep Dive: Go Channel Internals — hchan, Scheduling, Select & Production Patterns

> Everything the runtime does when you create a channel, send a value,
> receive a value, close it, or multiplex with select.

---

## Table of Contents

1. [The Runtime Struct: `hchan`](#1-the-runtime-struct-hchan)
2. [Buffered vs Unbuffered — What Actually Differs](#2-buffered-vs-unbuffered--what-actually-differs)
3. [Step-by-Step: Send Operation](#3-step-by-step-send-operation)
4. [Step-by-Step: Receive Operation](#4-step-by-step-receive-operation)
5. [The `sudog` — Goroutine Parking](#5-the-sudog--goroutine-parking)
5b. [Channel State Machine](#5b-channel-state-machine)
6. [Close Semantics — Under the Hood](#6-close-semantics--under-the-hood)
7. [Select Statement Internals](#7-select-statement-internals)
8. [Nil Channel Behavior — Why It Blocks Forever](#8-nil-channel-behavior--why-it-blocks-forever)
9. [Channel Direction Types](#9-channel-direction-types)
10. [The `chan struct{}` Pattern — Why Not `chan bool`](#10-the-chan-struct-pattern--why-not-chan-bool)
11. [Production Patterns](#11-production-patterns) ← moved to `concepts/` folder
12. [Performance Characteristics](#12-performance-characteristics--when-not-to-use-channels)
13. [Common Bugs & Gotchas](#13-common-bugs--gotchas)
14. [Quick Reference Card](#14-quick-reference-card)
15. [Select in Practice — Quick Reference](#15-select-in-practice--quick-reference)

---

## 1. The Runtime Struct: `hchan`

When you write `make(chan int)` or `make(chan int, 5)`, the runtime allocates a `runtime.hchan`
struct on the heap. A channel variable is always a pointer to this struct.

**Source:** `runtime/chan.go`
```go
type hchan struct {
    qcount   uint           // current number of elements in the buffer
    dataqsiz uint           // capacity of the circular buffer (0 for unbuffered)
    buf      unsafe.Pointer // pointer to the circular ring buffer
    elemsize uint16         // size of each element in bytes
    closed   uint32         // 0 = open, 1 = closed
    elemtype *_type          // type descriptor for elements (for GC)
    sendx    uint           // send index into the ring buffer
    recvx    uint           // receive index into the ring buffer
    recvq    waitq          // linked list of blocked receiver goroutines
    sendq    waitq          // linked list of blocked sender goroutines
    lock     mutex          // protects all fields of hchan
}
```

### Visual Layout

```
  make(chan int, 5) allocates:

  hchan (on heap)
  ┌─────────────────────────────────────────────────────────┐
  │  qcount   = 0          (elements currently in buffer)   │
  │  dataqsiz = 5          (buffer capacity)                │
  │  buf ──────────────────► ┌───┬───┬───┬───┬───┐          │
  │                          │   │   │   │   │   │ (5 slots)│
  │                          └───┴───┴───┴───┴───┘          │
  │  elemsize = 8          (sizeof(int) on 64-bit)          │
  │  closed   = 0          (channel is open)                │
  │  sendx    = 0          (next write position)            │
  │  recvx    = 0          (next read position)             │
  │  recvq    = {nil, nil} (no waiting receivers)           │
  │  sendq    = {nil, nil} (no waiting senders)             │
  │  lock     = mutex{}    (unlocked)                       │
  └─────────────────────────────────────────────────────────┘
```

### Key Insight: Channel Variables Are Pointers

```go
ch := make(chan int, 5)
```

`ch` itself is just a pointer (8 bytes on 64-bit) to the `hchan` on the heap.
Passing a channel to a function copies the pointer, not the struct — both the
caller and callee operate on the **same** channel. This is why channels work
for cross-goroutine communication without additional indirection.

```
  goroutine 1                    goroutine 2
  ┌────────┐                    ┌────────┐
  │ ch ────┼───────┐    ┌──────┼── ch   │
  └────────┘       │    │      └────────┘
                   ▼    ▼
               ┌──────────┐
               │  hchan   │  ← same struct
               └──────────┘
```

---

## 2. Buffered vs Unbuffered — What Actually Differs

The **only** structural difference is whether `dataqsiz` is zero.

### Unbuffered: `make(chan int)`

```
hchan:
  dataqsiz = 0    ← no buffer
  buf      = nil  ← no ring buffer allocated
  sendq/recvq     ← all synchronization happens here
```

An unbuffered channel has **no storage**. Every send must rendezvous directly with
a receiver. This makes unbuffered channels a **synchronization primitive**, not a
data structure. The sender blocks until a receiver is ready, and vice versa.

Under the hood, when a send finds a waiting receiver:
- The runtime copies the value **directly from the sender's stack to the receiver's stack**
- No buffer is touched — `runtime.send()` calls `typedmemmove()` directly
- This is the fastest path — one memory copy, no ring buffer overhead

### Buffered: `make(chan int, 5)`

```
hchan:
  dataqsiz = 5
  buf ──► [slot0][slot1][slot2][slot3][slot4]   ← circular ring buffer
  sendx, recvx wrap around using modulo
```

A buffered channel decouples sender and receiver in time. The sender only blocks
when the buffer is full; the receiver only blocks when the buffer is empty.

The ring buffer works like this:

```
  After: ch <- 10; ch <- 20; ch <- 30

  buf: [  10  |  20  |  30  |     |     ]
        ^recvx              ^sendx

  After: <-ch  (receives 10)

  buf: [      |  20  |  30  |     |     ]
               ^recvx       ^sendx

  Indices wrap around with modulo:
    sendx = (sendx + 1) % dataqsiz
    recvx = (recvx + 1) % dataqsiz
```

### `make(chan int, 0)` vs `make(chan int)` — Same Thing

Both create an unbuffered channel. `dataqsiz == 0`, no ring buffer.

### `make(chan int, 1)` — The Special Case

A buffer of 1 is **not** the same as unbuffered. It allows one send to proceed
without a waiting receiver. This is used as a lightweight semaphore:

```go
sem := make(chan struct{}, 1)
sem <- struct{}{}   // acquire (doesn't block — buffer has space)
// ... critical section ...
<-sem               // release
```

### Many Senders, Many Receivers — How Channels Handle It

A single channel can have **any number of goroutines** sending and receiving
simultaneously. Each value goes to exactly **one** receiver — no duplication,
no loss.

```go
ch := make(chan int, 5)

// 3 receivers competing on the same channel
go func() { fmt.Println("R1 got:", <-ch) }()
go func() { fmt.Println("R2 got:", <-ch) }()
go func() { fmt.Println("R3 got:", <-ch) }()

ch <- 42  // exactly ONE of R1, R2, R3 gets this value
```

Under the hood, when all 3 receivers arrive before any value is sent,
they each get parked on `recvq`:

```
hchan.recvq (waiting receivers — FIFO linked list):
┌────────┐    ┌────────┐    ┌────────┐
│ sudog  │───►│ sudog  │───►│ sudog  │───► nil
│ g: R1  │    │ g: R2  │    │ g: R3  │
└────────┘    └────────┘    └────────┘
   first                      last

ch <- 42 arrives:
  1. Lock hchan.lock
  2. Dequeue FIRST sudog (R1) from recvq    ← FIFO order
  3. Copy 42 directly to R1's stack variable
  4. Wake R1 (goready)
  5. Unlock hchan.lock
  6. R2 and R3 remain parked — they wait for the next values
```

The same applies to senders. Multiple senders park on `sendq` in FIFO order:

```
hchan.sendq (waiting senders — FIFO linked list):
┌────────┐    ┌────────┐
│ sudog  │───►│ sudog  │───► nil
│ g: S1  │    │ g: S2  │
│ elem:42│    │ elem:99│
└────────┘    └────────┘

<-ch arrives:
  1. Lock hchan.lock
  2. Dequeue FIRST sudog (S1) from sendq
  3. Copy 42 from S1's stack to receiver's variable
  4. Wake S1 (goready)
  5. Unlock hchan.lock
  6. S2 remains parked — waits until another receiver arrives
```

**This is the Worker Pool pattern built into the runtime:**

```
                            ┌── Worker 1 (receiver) ── process job
 jobs channel ──────────────┼── Worker 2 (receiver) ── process job
                            └── Worker 3 (receiver) ── process job

 // All 3 workers call <-jobs — runtime distributes one job to one worker
 for job := range jobs {
     process(job)
 }
```

Whichever worker finishes first and calls `<-jobs` gets the next job. The
`hchan.lock` mutex guarantees no two receivers grab the same value. Zero
extra code needed for work distribution.

**Summary:**

```
┌─────────────────────────────────────────────────────────┐
│           Many-to-Many Channel Guarantees                │
│                                                          │
│  Senders blocked?   → parked on sendq (FIFO queue)      │
│  Receivers blocked? → parked on recvq (FIFO queue)      │
│                                                          │
│  Value sent     → delivered to exactly ONE receiver      │
│  Value received → comes from exactly ONE sender          │
│                                                          │
│  No duplication. No loss. FIFO ordering within queues.   │
│  All protected by a single hchan.lock mutex.             │
└─────────────────────────────────────────────────────────┘
```

### The `hchan.lock` Mutex — How It Works Under the Hood

Every channel operation (send, receive, close) acquires `hchan.lock` first.
This is not a `sync.Mutex` — it's a lower-level **runtime mutex** (`runtime.mutex`)
that uses a different strategy than what you get from the `sync` package.

**Source:** `runtime/lock_futex.go` (Linux) / `runtime/lock_sema.go` (Windows/macOS)

```go
// runtime/runtime2.go
type mutex struct {
    lockRankStruct          // for debugging lock ordering
    key uintptr             // 0 = unlocked, other values = locked state
}
```

The runtime mutex uses a **two-phase locking strategy:**

```
Phase 1: SPIN (optimistic — bet that the lock is released quickly)
───────────────────────────────────────────────────────────────────
  The goroutine tries an atomic compare-and-swap (CAS) on mutex.key:
    CAS(&key, 0, locked)  →  if key was 0, set to locked. Done!

  If CAS fails (someone else holds the lock):
    Spin for a few iterations — execute PAUSE instructions (busy-wait).
    Why? If the lock holder is on another CPU core and will release
    in nanoseconds, spinning is FASTER than sleeping.

    Spin count is limited (~4 iterations on most architectures).
    Spinning only helps for very short critical sections.

Phase 2: SLEEP (pessimistic — lock is contended, stop wasting CPU)
───────────────────────────────────────────────────────────────────
  If spinning didn't acquire the lock:

  On Linux:
    futex(FUTEX_WAIT) — kernel puts the thread to sleep on the mutex address.
    When the lock holder unlocks: futex(FUTEX_WAKE) wakes ONE sleeping thread.
    Cost: ~1-2 microseconds (kernel context switch).

  On Windows:
    Uses OS semaphore — similar mechanism, different syscall.

  On macOS:
    Uses pthread_mutex or os_unfair_lock depending on Go version.
```

**Why does this matter for channels?**

```
┌─────────────────────────────────────────────────────────────────────┐
│ Every channel send/receive locks hchan.lock:                        │
│                                                                     │
│   Uncontended (one goroutine using the channel):                    │
│     → CAS succeeds immediately → ~5-10ns overhead                  │
│     → This is the fast path. Channels are cheap when not contested. │
│                                                                     │
│   Low contention (a few goroutines):                                │
│     → Short spin → CAS succeeds → ~20-50ns overhead                │
│     → Still fast. The spin avoids the expensive kernel call.        │
│                                                                     │
│   High contention (100 goroutines on one channel):                  │
│     → Spin fails → futex sleep → kernel wakeup → ~1-2μs overhead   │
│     → The channel becomes a BOTTLENECK.                             │
│     → All goroutines serialize through one mutex.                   │
│     → Solution: shard into multiple channels, or use sync primitives│
└─────────────────────────────────────────────────────────────────────┘
```

**Comparison: `hchan.lock` (runtime mutex) vs `sync.Mutex`:**

```
┌──────────────────────┬─────────────────────────┬──────────────────────────┐
│                      │ runtime.mutex            │ sync.Mutex               │
│                      │ (used by channels)       │ (used by your code)      │
├──────────────────────┼─────────────────────────┼──────────────────────────┤
│ Where used           │ Internal runtime only    │ User-space Go code       │
│ Goroutine-aware?     │ No — blocks the OS       │ Yes — parks the          │
│                      │ thread (M), not just     │ goroutine (G), releases  │
│                      │ the goroutine            │ the thread (M) for       │
│                      │                          │ other goroutines         │
│ Spin phase           │ Yes (limited)            │ Yes (adaptive since      │
│                      │                          │ Go 1.9, starvation mode) │
│ Blocking mechanism   │ futex / OS semaphore     │ runtime.gopark (parks G, │
│                      │ (puts OS thread to       │ M picks another G from   │
│                      │ sleep)                   │ run queue)               │
│ Why different?       │ Runtime can't use its    │ Built on top of the      │
│                      │ own scheduler to park    │ runtime scheduler — can   │
│                      │ goroutines while inside  │ cooperatively yield      │
│                      │ the scheduler code       │                          │
└──────────────────────┴─────────────────────────┴──────────────────────────┘
```

**The critical insight:** `hchan.lock` blocks the **OS thread**, not just the goroutine.
This is necessary because channel operations are part of the scheduler itself — you can't
use the scheduler to sleep while you're inside scheduler code. This is why high-contention
channels are expensive: they waste OS threads, not just goroutines.

However, the lock is held for an extremely short time (copy a value, update an index,
maybe wake a goroutine). So in practice, contention only matters when you have many
goroutines hammering the same channel at extreme throughput (>1M ops/sec).

---

## 3. Step-by-Step: Send Operation

When you write `ch <- value`, the runtime calls `runtime.chansend()`.
There are **three possible paths**, checked in order:

### Path 1: Receiver Already Waiting (fastest)

```
ch <- 42   (and there's a goroutine blocked on <-ch)

  1. Lock hchan.lock
  2. Dequeue a waiting receiver from recvq
  3. Copy value DIRECTLY from sender's stack into receiver's variable
     → runtime.send() → typedmemmove(elem, src, dst)
     → NO ring buffer involvement
  4. Wake the receiver goroutine (change state _Gwaiting → _Grunnable)
  5. Unlock hchan.lock

  This is the DIRECT SEND path — fastest because it bypasses the buffer entirely.
```

```
  Sender goroutine G1              Receiver goroutine G2
  ┌───────────────┐                ┌───────────────┐
  │ stack: val=42  │ ──────────►  │ stack: v=42    │
  └───────────────┘   direct      └───────────────┘
                      copy           (was parked on
                                      recvq, now
                                      woken up)
```

### Path 2: Buffer Has Space

```
ch <- 42   (no waiting receiver, but buffer is not full)

  1. Lock hchan.lock
  2. Copy value into buf[sendx]
     → typedmemmove(elemtype, bufSlot, &value)
  3. sendx = (sendx + 1) % dataqsiz
  4. qcount++
  5. Unlock hchan.lock
  6. Sender continues immediately — no blocking
```

### Path 3: Buffer Full (or Unbuffered With No Receiver) — Block

```
ch <- 42   (buffer full or unbuffered, no waiting receiver)

  1. Lock hchan.lock
  2. Create a sudog for this goroutine (see Section 5)
  3. Enqueue sudog onto sendq
  4. Call runtime.gopark() — suspends this goroutine
     → goroutine state: _Grunning → _Gwaiting
     → scheduler picks another goroutine to run
  5. (later, when a receiver arrives, this goroutine is woken)
  6. Unlock hchan.lock
```

### Visual Summary

```
  ch <- value
  │
  ├─ recvq non-empty?
  │   └─ YES → DIRECT SEND: copy to receiver's stack, wake receiver
  │
  ├─ buffer has space? (qcount < dataqsiz)
  │   └─ YES → BUFFER SEND: copy to buf[sendx], advance sendx
  │
  └─ neither?
      └─ PARK: create sudog, enqueue on sendq, goroutine sleeps
```

---

## 4. Step-by-Step: Receive Operation

When you write `v := <-ch`, the runtime calls `runtime.chanrecv()`.
Again, three paths:

### Path 1: Sender Already Waiting (unbuffered or full buffer)

**Case A — Unbuffered channel with waiting sender:**
```
  1. Lock hchan.lock
  2. Dequeue sender from sendq
  3. Copy value DIRECTLY from sender's stack into receiver's variable
  4. Wake sender goroutine
  5. Unlock hchan.lock
```

**Case B — Buffered channel that is full, with waiting sender:**
```
  1. Lock hchan.lock
  2. Copy buf[recvx] into receiver's variable (the oldest element)
  3. Copy waiting sender's value into buf[recvx] (fill the slot just freed)
  4. Advance recvx
  5. Wake sender goroutine
  6. Unlock hchan.lock

  This maintains FIFO order — the receiver gets the oldest buffered value,
  and the sender's value goes to the end of the buffer.
```

### Path 2: Buffer Has Data (no waiting senders)

```
v := <-ch   (buffer has elements, no waiting senders)

  1. Lock hchan.lock
  2. Copy buf[recvx] into receiver's variable
  3. recvx = (recvx + 1) % dataqsiz
  4. qcount--
  5. Unlock hchan.lock
  6. Receiver continues with the value
```

### Path 3: Buffer Empty (or Unbuffered With No Sender) — Block

```
v := <-ch   (nothing to receive)

  1. Lock hchan.lock
  2. Create a sudog for this goroutine
  3. Enqueue sudog onto recvq
  4. Call runtime.gopark() — goroutine sleeps
  5. (later, when a sender arrives, this goroutine is woken)
  6. Value is written directly by the sender (Path 1 of send)
  7. Unlock hchan.lock
```

### The Two-Value Receive: `v, ok := <-ch`

```go
v, ok := <-ch
```

- `ok == true` → value was sent by a sender (normal operation)
- `ok == false` → channel was **closed** and buffer is **empty**; `v` is the zero value

Under the hood, `chanrecv()` returns a `received` boolean. When `closed == 1` and
`qcount == 0`, it sets `received = false` and zeroes the receiver's memory.

---

## 5. The `sudog` — Goroutine Parking

When a goroutine blocks on a channel (send to full, receive from empty), the runtime
creates a `sudog` ("sudo-goroutine") — a wrapper that links the goroutine to the
channel's wait queue.

**Source:** `runtime/runtime2.go`
```go
type sudog struct {
    g       *g              // the blocked goroutine
    elem    unsafe.Pointer  // pointer to the data being sent/received
    c       *hchan          // the channel this sudog is parked on
    next    *sudog          // linked list → next waiting goroutine
    prev    *sudog          // linked list → previous waiting goroutine
    // ... other fields for select, timers
}
```

### Wait Queue Structure

```
  hchan.sendq (goroutines blocked on send):
  ┌────────┐    ┌────────┐    ┌────────┐
  │ sudog  │───►│ sudog  │───►│ sudog  │───► nil
  │ g: G4  │    │ g: G7  │    │ g: G12 │
  │ elem:42│    │ elem:99│    │ elem:7 │
  └────────┘    └────────┘    └────────┘
     first         ...          last

  hchan.recvq (goroutines blocked on receive):
  ┌────────┐    ┌────────┐
  │ sudog  │───►│ sudog  │───► nil
  │ g: G2  │    │ g: G9  │
  │ elem:&v│    │ elem:&w│      ← elem points to the receiver's variable
  └────────┘    └────────┘
```

### sudog Lifecycle

```
  1. Goroutine G wants to send to full channel
  2. Runtime acquires sudog from per-P cache (or allocates new one)
  3. Fills in: g = G, elem = &value, c = ch
  4. Enqueues onto ch.sendq
  5. Calls gopark(G) — G becomes _Gwaiting
  6. ... time passes ...
  7. Receiver arrives, dequeues this sudog from sendq
  8. Copies value from sudog.elem directly
  9. Calls goready(G) — G becomes _Grunnable, added to run queue
  10. sudog returned to per-P cache for reuse
```

### Per-P sudog Cache — Why Channel Ops Are Allocation-Free

To understand the "per-P cache", you first need to know what a **P** is.

**The GMP Model (quick recap):**

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Go Scheduler: GMP Model                      │
│                                                                     │
│  G = Goroutine    (your code — millions of these)                  │
│  M = Machine      (OS thread — typically tens of these)             │
│  P = Processor    (logical CPU — count = GOMAXPROCS, e.g., 8)      │
│                                                                     │
│  Rule: A goroutine (G) can only run when an M has a P.             │
│  Each P has its own LOCAL resources to avoid global locks.          │
│                                                                     │
│   P0                P1                P2                P3          │
│   ┌─────────┐      ┌─────────┐      ┌─────────┐      ┌─────────┐ │
│   │ run queue│      │ run queue│      │ run queue│      │ run queue│ │
│   │ mcache   │      │ mcache   │      │ mcache   │      │ mcache   │ │
│   │ sudog    │      │ sudog    │      │ sudog    │      │ sudog    │ │
│   │  cache   │      │  cache   │      │  cache   │      │  cache   │ │
│   └────┬────┘      └────┬────┘      └────┬────┘      └────┬────┘ │
│        │                │                │                │        │
│        ▼                ▼                ▼                ▼        │
│   M0 (thread)      M1 (thread)     M2 (thread)      M3 (thread)  │
│   running G5       running G12     running G3       running G8    │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

Each P has its own **sudog free list** — a stack of pre-allocated, reusable sudog structs.

**Why "per-P" instead of one global pool?**

```
Global pool approach (BAD):
  ┌──────────────────────────┐
  │  Global sudog free list  │  ← single lock protects this
  │  ┌───┬───┬───┬───┬───┐  │
  │  │ s │ s │ s │ s │ s │  │
  │  └───┴───┴───┴───┴───┘  │
  └──────────────────────────┘
       ▲    ▲    ▲    ▲
       │    │    │    │
      P0   P1   P2   P3      ALL P's compete for ONE lock
                              → contention at high throughput!

Per-P pool approach (GOOD — what Go actually does):
  P0 sudog cache    P1 sudog cache    P2 sudog cache    P3 sudog cache
  ┌───┬───┬───┐    ┌───┬───┬───┐    ┌───┬───┬───┐    ┌───┬───┬───┐
  │ s │ s │ s │    │ s │ s │ s │    │ s │ s │ s │    │ s │ s │ s │
  └───┴───┴───┘    └───┴───┴───┘    └───┴───┴───┘    └───┴───┴───┘
       ▲                ▲                ▲                ▲
       │                │                │                │
      P0 only          P1 only          P2 only          P3 only
      NO lock needed!  NO lock needed!  NO lock needed!  NO lock needed!
```

**The key insight:** When a goroutine running on P2 needs a sudog (because it's
about to block on a channel), it grabs one from P2's local cache — **no lock,
no contention, no heap allocation**. When the sudog is returned (after the
goroutine is woken), it goes back to P2's cache for reuse.

This is the same design pattern used throughout the Go runtime:
- **Per-P sudog cache** — for channel blocking (this section)
- **Per-P mcache** — for small memory allocations (avoids global allocator lock)
- **Per-P run queue** — for goroutine scheduling (avoids global queue lock)
- **Per-P defer pool** — for defer struct allocation

The design philosophy: **"If each CPU core has its own resources, cores don't
fight over shared locks."**

### sudog Cache Lifecycle — Full Picture

```
  ┌─────────────────────────────────────────────────────────────────────────┐
  │                    sudog Cache Lifecycle                                │
  │                                                                         │
  │  STEP 1: Goroutine G (running on P2) needs to block on channel ch      │
  │                                                                         │
  │    P2's sudog cache: [s1, s2, s3]   ← 3 sudogs available              │
  │                                                                         │
  │    Runtime pops s1 from P2's cache   ← NO lock, NO heap allocation     │
  │    Fills in: s1.g = G, s1.elem = &value, s1.c = ch                    │
  │    Enqueues s1 onto ch.sendq                                           │
  │    Parks G (gopark)                                                    │
  │                                                                         │
  │    P2's sudog cache: [s2, s3]        ← 2 remaining                    │
  │                                                                         │
  │  STEP 2: A receiver on P0 dequeues s1 from ch.sendq                   │
  │                                                                         │
  │    Copies value from s1.elem                                           │
  │    Wakes G (goready) → G goes to a run queue                           │
  │    Returns s1 to... which P's cache?                                    │
  │                                                                         │
  │    → The P that is running the wakeup code (P0 in this case)           │
  │    → s1 goes to P0's cache, NOT P2's!                                  │
  │    → This is fine — sudogs are fungible (interchangeable)              │
  │                                                                         │
  │  STEP 3: If a P's cache is empty?                                      │
  │                                                                         │
  │    → Allocate a new sudog from the heap (mallocgc)                     │
  │    → This is the SLOW path — only happens when the cache is exhausted  │
  │    → In practice, sudogs are quickly returned and reused               │
  │                                                                         │
  │  STEP 4: If a P's cache is too full?                                   │
  │                                                                         │
  │    → Excess sudogs are returned to a global pool                       │
  │    → Other P's with empty caches can take from the global pool         │
  │    → Global pool access DOES require a lock, but rarely happens        │
  │                                                                         │
  └─────────────────────────────────────────────────────────────────────────┘
```

### Why This Matters for Performance

```
Channel send to full buffer (common case):
  1. Need a sudog → grab from P's cache → ~5ns (pointer pop from list)
  2. Park goroutine → ~50ns (scheduler state change)
  3. ... blocked ...
  4. Woken by receiver → ~50ns
  5. Return sudog to P's cache → ~5ns (pointer push to list)

  Total sudog overhead: ~10ns. NO heap allocation. NO GC pressure.

Without per-P cache (hypothetical):
  1. Need a sudog → mallocgc(sizeof(sudog)) → ~200-500ns (heap alloc)
  2. ... same blocking ...
  3. sudog becomes garbage → GC must scan and collect → adds GC pressure

  Total sudog overhead: ~500ns+ and GC work.
  At 1M channel ops/sec, that's 1M allocations/sec → significant GC pressure.
```

---

## 5b. Channel State Machine

The individual scenarios (buffered lifecycle, unbuffered direct transfer, sender
blocking, select multi-channel, close broadcasting) are demonstrated as runnable
Go code in `fundamentals/11_channels/concepts/` (files 05–09).

Below is the master diagram summarizing every path a send or receive can take:

### Master Diagram: Complete Channel State Machine

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CHANNEL OPERATION STATE MACHINE                       │
│                                                                         │
│  SEND (ch <- val):                                                      │
│  ═════════════════                                                      │
│                                                                         │
│    ┌──────────────┐     YES    ┌──────────────────────────┐             │
│    │ ch == nil?   ├──────────► │ BLOCK FOREVER (gopark)   │             │
│    └──────┬───────┘            └──────────────────────────┘             │
│           │ NO                                                          │
│           ▼                                                             │
│    ┌──────────────┐     YES    ┌──────────────────────────┐             │
│    │ ch.closed?   ├──────────► │ PANIC "send on closed"   │             │
│    └──────┬───────┘            └──────────────────────────┘             │
│           │ NO                                                          │
│           ▼                                                             │
│    ┌──────────────────┐  YES   ┌──────────────────────────┐             │
│    │ recvq non-empty? ├──────► │ DIRECT SEND to receiver  │  (fastest) │
│    └──────┬───────────┘        │ Copy val → receiver stack│             │
│           │ NO                  │ Wake receiver goroutine  │             │
│           ▼                    └──────────────────────────┘             │
│    ┌──────────────────┐  YES   ┌──────────────────────────┐             │
│    │ buffer has space?├──────► │ BUFFER SEND              │  (fast)    │
│    │ qcount<dataqsiz  │        │ Copy val → buf[sendx]    │             │
│    └──────┬───────────┘        │ sendx++, qcount++        │             │
│           │ NO                  └──────────────────────────┘             │
│           ▼                                                             │
│    ┌──────────────────────────┐                                         │
│    │ BLOCK (sudog + gopark)  │  (slow — goroutine sleeps)              │
│    │ Get sudog from P's cache│                                         │
│    │ Enqueue on sendq        │                                         │
│    │ G state → _Gwaiting     │                                         │
│    └──────────────────────────┘                                         │
│                                                                         │
│  RECEIVE (v := <-ch):                                                   │
│  ════════════════════                                                   │
│                                                                         │
│    ┌──────────────┐     YES    ┌──────────────────────────┐             │
│    │ ch == nil?   ├──────────► │ BLOCK FOREVER (gopark)   │             │
│    └──────┬───────┘            └──────────────────────────┘             │
│           │ NO                                                          │
│           ▼                                                             │
│    ┌───────────────────────┐ YES ┌────────────────────────┐             │
│    │ closed && qcount==0?  ├───► │ Return (zero, false)   │  (instant) │
│    └──────┬────────────────┘     └────────────────────────┘             │
│           │ NO                                                          │
│           ▼                                                             │
│    ┌──────────────────┐  YES   ┌──────────────────────────┐             │
│    │ sendq non-empty? ├──────► │ DIRECT RECV from sender  │  (fastest) │
│    │ (unbuf or full)  │        │ If unbuf: copy from sender│            │
│    └──────┬───────────┘        │ If full: take buf, put    │            │
│           │ NO                  │ sender's val in freed slot│            │
│           ▼                    └──────────────────────────┘             │
│    ┌──────────────────┐  YES   ┌──────────────────────────┐             │
│    │ buffer has data? ├──────► │ BUFFER RECV              │  (fast)    │
│    │ qcount > 0       │        │ Copy buf[recvx] → v      │             │
│    └──────┬───────────┘        │ recvx++, qcount--        │             │
│           │ NO                  └──────────────────────────┘             │
│           ▼                                                             │
│    ┌──────────────────────────┐                                         │
│    │ BLOCK (sudog + gopark)  │  (slow — goroutine sleeps)              │
│    │ Get sudog from P's cache│                                         │
│    │ Enqueue on recvq        │                                         │
│    │ G state → _Gwaiting     │                                         │
│    └──────────────────────────┘                                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Close Semantics — Under the Hood

When you write `close(ch)`, the runtime calls `runtime.closechan()`:

```
close(ch):
  1. Lock hchan.lock
  2. Set ch.closed = 1
  3. Iterate ALL waiting receivers (recvq):
     → For each: set their elem to zero value, wake them
     → They receive (zero value, ok=false)
  4. Iterate ALL waiting senders (sendq):
     → For each: PANIC — "send on closed channel"
     → Actually, these are woken and panic in their goroutine
  5. Unlock hchan.lock
```

### The Rules (with runtime explanation)

```
┌──────────────────────────────────┬──────────────────────────────────────────┐
│ Operation on closed channel      │ What happens (and WHY)                   │
├──────────────────────────────────┼──────────────────────────────────────────┤
│ ch <- value                      │ PANIC: "send on closed channel"          │
│                                  │ Runtime checks ch.closed before send.    │
│                                  │ Rationale: sending to closed channel     │
│                                  │ would lose data silently — Go prefers    │
│                                  │ to crash loud over silent data loss.     │
├──────────────────────────────────┼──────────────────────────────────────────┤
│ v := <-ch  (buffer non-empty)    │ Returns buffered value, ok=true          │
│                                  │ Closed channels drain remaining buffer   │
│                                  │ before returning zero values.            │
├──────────────────────────────────┼──────────────────────────────────────────┤
│ v := <-ch  (buffer empty)        │ Returns zero value immediately, ok=false │
│                                  │ Never blocks — closed + empty = done.    │
├──────────────────────────────────┼──────────────────────────────────────────┤
│ close(ch) again                  │ PANIC: "close of closed channel"         │
│                                  │ Runtime checks ch.closed at entry.       │
├──────────────────────────────────┼──────────────────────────────────────────┤
│ close(nil)                       │ PANIC: "close of nil channel"            │
│                                  │ No hchan to set closed on.               │
└──────────────────────────────────┴──────────────────────────────────────────┘
```

### The "Only Sender Closes" Rule

This isn't enforced by the compiler — it's a **convention** backed by the panic semantics:

- If receiver closes, sender panics on next send → program crashes
- If sender closes, receiver just gets zero values → graceful degradation

Multiple senders sharing one channel? Use a `sync.Once` or `sync.WaitGroup` to
coordinate a single close after all senders finish. Never have each sender try to close.

```go
// Pattern: multiple senders, one coordinated close
func fanIn(sources ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(sources))

    for _, src := range sources {
        go func(s <-chan int) {
            defer wg.Done()
            for v := range s {
                out <- v
            }
        }(src)
    }

    go func() {
        wg.Wait()    // wait for ALL senders to finish
        close(out)   // single close point
    }()

    return out
}
```

---

## 7. Select Statement Internals

`select` is the most complex channel operation. It multiplexes across multiple channels.

### Select is NOT Switch — The Mental Model

`switch` evaluates cases **top to bottom** and runs the first match.
`select` evaluates ALL cases **simultaneously** and runs whichever channel is ready.

```
  switch:                          select:
  ┌────────────────┐               ┌────────────────────────────────┐
  │ case a:  ← 1st │               │ case <-ch1:  ← checked        │
  │ case b:  ← 2nd │               │ case <-ch2:  ← simultaneously │
  │ case c:  ← 3rd │               │ case ch3<-v: ← with all others│
  │ default:       │               │ default:                       │
  └────────────────┘               └────────────────────────────────┘
  First true case wins.            First READY channel wins.
  Deterministic order.             Random if multiple ready.
```

### What Each Case Type Means

Every `case` in a `select` is a channel operation — either a send or a receive.
The `select` asks: "which of these operations can proceed right now?"

```go
  select {
  case v, ok := <-in:        // "TRY to receive from 'in'"
                              //   ready if: in has data, or in is closed
                              //   blocks if: in is open and empty

  case out <- v:              // "TRY to send v to 'out'"
                              //   ready if: someone is receiving, or buffer has space
                              //   blocks if: buffer full and no receiver

  case <-ctx.Done():          // "TRY to receive from ctx's internal channel"
                              //   ready if: ctx was cancelled (channel closed)
                              //   blocks if: ctx still active (channel open)

  default:                    // "if NOTHING is ready, run this immediately"
                              //   makes select non-blocking
  }
```

### ctx.Done() — It's Just a Channel

There is no magic signal behind context cancellation. `ctx.Done()` returns a
plain `<-chan struct{}` that gets **closed** when someone calls `cancel()`.

```
  ctx, cancel := context.WithCancel(background)

  ctx.Done() returns: <-chan struct{}

  Before cancel():
    <-ctx.Done()    → BLOCKS (channel is open, nothing sent, nobody will send)

  After cancel():
    close(done)     → happens internally when cancel() is called
    <-ctx.Done()    → RETURNS immediately (zero value from closed channel)
    <-ctx.Done()    → RETURNS immediately AGAIN (closed channels never block)
```

This is why `ctx.Done()` works seamlessly in `select` — the runtime treats it
as any other channel receive. When cancelled, the closed channel is always
"ready," so `select` picks that case.

```
  select {
  case v := <-dataCh:       // channel receive: waiting for data
  case <-ctx.Done():        // channel receive: waiting for cancellation
  }

  The runtime doesn't know one is "cancellation" — it just sees two
  channels and picks whichever becomes ready first.
```

### What the Compiler Generates

```go
select {
case v := <-ch1:
    fmt.Println(v)
case ch2 <- 42:
    fmt.Println("sent")
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("nothing ready")
}
```

The compiler translates this into a call to `runtime.selectgo()`, passing an array
of `scase` structs (one per case):

```go
// runtime/select.go
type scase struct {
    c    *hchan         // the channel for this case
    elem unsafe.Pointer // pointer to send/receive data
}
```

### The `selectgo()` Algorithm

```
runtime.selectgo(cases []scase, order []uint16):

  Phase 1: RANDOM SHUFFLE
  ────────────────────────
  Shuffle the polling order randomly.
  WHY? To prevent starvation. If two cases are always ready,
  both get a fair chance. Without shuffling, the first case
  in source order would always win.

  Phase 2: LOCK ALL CHANNELS
  ──────────────────────────
  Lock every unique channel's mutex (in address order to prevent deadlock).
  This is the expensive part — acquiring N mutexes.

  Phase 3: POLL (non-blocking check)
  ──────────────────────────────────
  Walk cases in shuffled order:
    For each receive case:
      - sendq non-empty? → direct receive from sender → DONE
      - buffer non-empty? → buffer receive → DONE
    For each send case:
      - recvq non-empty? → direct send to receiver → DONE
      - buffer has space? → buffer send → DONE

  If any case is ready → unlock all, return that case.
  If default exists and nothing ready → unlock all, return default.

  Phase 4: ENQUEUE ON ALL CHANNELS (nothing ready, no default)
  ─────────────────────────────────────────────────────────────
  Create a sudog for THIS goroutine for EACH case.
  Enqueue each sudog onto the corresponding channel's sendq or recvq.
  The goroutine is now parked on MULTIPLE channels simultaneously.
  Call gopark() — goroutine sleeps.

  Phase 5: WAKEUP
  ──────────────
  When ANY of the channels becomes ready, the goroutine is woken.
  Dequeue this goroutine's sudogs from ALL OTHER channels.
  Return the case that triggered the wakeup.
```

### Visual: Goroutine Parked on Multiple Channels

```
  select {
  case v := <-ch1:     // case 0
  case v := <-ch2:     // case 1
  case ch3 <- 42:      // case 2
  }

  Goroutine G5 is now waiting on 3 channels:

  ch1.recvq: [..., sudog{g:G5, case:0}]
  ch2.recvq: [..., sudog{g:G5, case:1}]
  ch3.sendq: [..., sudog{g:G5, case:2}]

  When ch2 gets a value:
    1. ch2 dequeues sudog{g:G5} from recvq
    2. Copies value to G5's variable
    3. Wakes G5
    4. G5 removes its sudogs from ch1.recvq and ch3.sendq
    5. selectgo returns case index 1
```

### Performance: Select is Expensive

```
  Cost of select with N cases:
  ─ N mutex lock/unlock operations (Phase 2 & unlock)
  ─ N sudog allocations (if blocking — Phase 4)
  ─ Random shuffle (O(N) Fisher-Yates)
  ─ Lock ordering sort (O(N log N) on channel addresses)

  For 2-3 cases: negligible
  For 10+ cases: measurable — consider redesigning
```

### Select with Default — Non-Blocking Operation

```go
select {
case v := <-ch:
    process(v)
default:
    // ch not ready — do something else
}
```

With `default`, `selectgo` never reaches Phase 4. It polls once and either
finds a ready case or falls through to default. This makes it a **non-blocking**
channel operation — the Go equivalent of `tryReceive()`.

---

## 8. Nil Channel Behavior — Why It Blocks Forever

```go
var ch chan int   // ch is nil — no hchan exists
<-ch             // blocks forever
ch <- 42         // blocks forever
```

When you operate on a nil channel, `runtime.chansend/chanrecv` check for `c == nil`
at the very top:

```go
// runtime/chan.go (simplified)
func chansend(c *hchan, ...) bool {
    if c == nil {
        gopark(nil, nil, waitReasonChanSendNilChan, ...)
        // goroutine parked forever — nothing will wake it
        // no hchan exists → no sendq/recvq to enqueue on
        // it's just... gone
    }
    ...
}
```

### Why Is This Useful?

Nil channels are a **select control mechanism**. You can dynamically disable
select cases by setting their channel to nil:

```go
func merge(ch1, ch2 <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for ch1 != nil || ch2 != nil {
            select {
            case v, ok := <-ch1:
                if !ok {
                    ch1 = nil  // disable this case — nil channel blocks forever
                    continue   // in select, won't be selected again
                }
                out <- v
            case v, ok := <-ch2:
                if !ok {
                    ch2 = nil  // disable this case
                    continue
                }
                out <- v
            }
        }
    }()
    return out
}
```

Without nil channel semantics, you'd need complex flag logic. Setting `ch1 = nil`
makes `select` naturally skip that case forever — the runtime says "nil channel
will never be ready" and excludes it from polling.

### Why `nil` Instead of `close`? — The Critical Difference

This is a common question: "Why not just close the channel instead of setting it to nil?"
Because they have **opposite behavior** inside a `for-select` loop:

```
  Closed channel in select:
  ┌────────────────────────────────────────────────┐
  │  <-closedCh  →  returns (zero, false) INSTANTLY │
  │                  select sees it as "ready"      │
  │                  picks it over and over and over │
  │                  CPU: 100% 🔥 (infinite loop)   │
  └────────────────────────────────────────────────┘

  Nil channel in select:
  ┌────────────────────────────────────────────────┐
  │  <-nilCh     →  blocks FOREVER                  │
  │                  select SKIPS this case entirely │
  │                  as if the case doesn't exist    │
  │                  CPU: 0% ✅                      │
  └────────────────────────────────────────────────┘
```

```go
// ❌ WRONG — closing inside for-select creates infinite busy loop:
for {
    select {
    case v, ok := <-ch:
        if !ok {
            // ch is closed, but select will pick this case
            // EVERY iteration — returns (0, false) instantly
            // This becomes an infinite loop burning 100% CPU!
        }
    case v := <-other:
        process(v)    // never gets a chance — closed ch always "wins"
    }
}

// ✅ CORRECT — setting to nil disables the case:
for ch != nil || other != nil {
    select {
    case v, ok := <-ch:
        if !ok {
            ch = nil       // select ignores this case from now on
            continue
        }
        process(v)
    case v, ok := <-other:
        if !ok {
            other = nil
            continue
        }
        process(v)
    }
}
```

**One-line rule:**

```
  close(ch)  = "I'm done sending"      → receivers get zero values FOREVER
  ch = nil   = "pretend this doesn't exist" → select skips it entirely
```

### Nil Channel Summary

```
┌───────────────────┬────────────────────────────────────────────┐
│ Operation         │ Nil channel behavior                       │
├───────────────────┼────────────────────────────────────────────┤
│ <-ch              │ Blocks forever (goroutine parked)           │
│ ch <- value       │ Blocks forever (goroutine parked)           │
│ close(ch)         │ PANIC: close of nil channel                │
│ len(ch)           │ 0                                          │
│ cap(ch)           │ 0                                          │
│ select case <-ch  │ Case is NEVER selected (effectively skip)  │
└───────────────────┴────────────────────────────────────────────┘
```

---

## 9. Channel Direction Types

Go channels have three direction types. The compiler enforces these at compile time.

```go
chan T      // bidirectional — can send and receive
chan<- T    // send-only — can only send
<-chan T    // receive-only — can only receive
```

### Conversions Are One-Way

```
  Bidirectional → send-only:  ✅ implicit conversion
  Bidirectional → recv-only:  ✅ implicit conversion
  Send-only → bidirectional:  ❌ compile error
  Recv-only → bidirectional:  ❌ compile error
  Send-only → recv-only:     ❌ compile error
```

Under the hood, all three are the same `*hchan` pointer. The direction restriction
is purely a **compile-time constraint** — zero runtime cost. The runtime doesn't
know or care about direction.

### Why Direction Matters

Direction communicates **ownership semantics**:

```go
// The function PRODUCES values — caller can only READ
func generate(n int) <-chan int {
    ch := make(chan int)  // bidirectional internally
    go func() {
        for i := 0; i < n; i++ { ch <- i }
        close(ch)
    }()
    return ch  // returned as receive-only → caller can't send or close
}

// The function CONSUMES values — caller sends, function reads
func process(jobs <-chan Job) {
    for job := range jobs {
        // ...
    }
}
```

This is Go's approach to enforcing **channel ownership at the API boundary**.
The sender controls the channel lifecycle (including close). The receiver just
reads until closed.

---

## 10. The `chan struct{}` Pattern — Why Not `chan bool`

You'll see `chan struct{}` throughout Go codebases for signaling:

```go
done := make(chan struct{})
close(done)  // signal all listeners
```

### Why `struct{}` instead of `bool`?

```
  sizeof(struct{}) = 0 bytes   ← zero memory for the value
  sizeof(bool)     = 1 byte    ← wastes space for information we don't use

  The signal is in the OPERATION (send/close), not the VALUE.
  Using bool implies the value matters — "was it true or false?"
  Using struct{} makes intent explicit — "this is a signal, not data."
```

Under the hood, `runtime.hchan` for `chan struct{}`:
- `elemsize = 0` — no space allocated per buffer slot
- The ring buffer exists but elements take zero bytes
- Send/receive still synchronize goroutines, but copy zero bytes of data

### When To Use What

```
chan struct{}    → pure signaling (done, quit, semaphore)
chan bool        → rare, when the true/false value actually matters
chan error       → when you need to communicate success/failure with detail
chan T           → when you're passing actual data
```

---

## 11. Production Patterns

All production patterns are implemented as runnable Go code with full explanations
in `fundamentals/11_channels/concepts/`:

| # | File | Pattern |
|---|------|---------|
| 01 | `01_worker_event_loop.go` | for/select event loop (Pub/Sub worker) |
| 02 | `02_heartbeat_monitor.go` | Timer reset inactivity detection |
| 03 | `03_first_response_wins.go` | Fan-out with buffered channel + cancel |
| 04 | `04_token_bucket_rate_limiter.go` | Buffered channel as counting semaphore |
| 05 | `05_buffered_channel_lifecycle.go` | Ring buffer walkthrough (sendx/recvx) |
| 06 | `06_unbuffered_direct_transfer.go` | Direct goroutine-to-goroutine copy |
| 07 | `07_sender_blocks_receiver_wakes.go` | sudog lifecycle (park → wake) |
| 08 | `08_select_multi_channel.go` | Goroutine on multiple wait queues |
| 09 | `09_close_wakes_receivers.go` | close() as broadcast signal |
| 10 | `10_done_cancellation.go` | Legacy done channel vs context.WithCancel |
| 11 | `11_worker_pool_backpressure.go` | N workers sharing a jobs channel |
| 12 | `12_semaphore_bounded_concurrency.go` | Buffered channel as max concurrency gate |
| 13 | `13_fan_out_fan_in.go` | One producer → N consumers → merged output |
| 14 | `14_pipeline.go` | Composable stages connected by channels |
| 15 | `15_timeout_with_context.go` | context.WithTimeout vs time.After leak |
| 16 | `16_or_done.go` | Double-select for cancellation-aware forwarding |

Each file explains: the problem → why channels fit → the pattern → runtime internals → runnable demo.

---

## 12. Performance Characteristics & When NOT to Use Channels

### Cost of Channel Operations

```
┌─────────────────────────────────┬───────────────┬───────────────────────────────┐
│ Operation                       │ Approx Cost   │ Why                           │
├─────────────────────────────────┼───────────────┼───────────────────────────────┤
│ Unbuffered send+recv (paired)   │ ~200-300ns    │ Mutex lock + gopark/goready + │
│                                 │               │ context switch                │
├─────────────────────────────────┼───────────────┼───────────────────────────────┤
│ Buffered send (not full)        │ ~50-100ns     │ Mutex lock + memcopy          │
├─────────────────────────────────┼───────────────┼───────────────────────────────┤
│ Buffered recv (not empty)       │ ~50-100ns     │ Mutex lock + memcopy          │
├─────────────────────────────────┼───────────────┼───────────────────────────────┤
│ Select with 2 cases             │ ~300-500ns    │ 2 mutex locks + poll + sudog  │
├─────────────────────────────────┼───────────────┼───────────────────────────────┤
│ Select with N cases             │ ~N × 150ns    │ N mutex locks + shuffle +     │
│                                 │               │ sort + poll                   │
├─────────────────────────────────┼───────────────┼───────────────────────────────┤
│ atomic.AddInt64                 │ ~5-15ns       │ Single CPU instruction (LOCK  │
│                                 │               │ XADD)                         │
├─────────────────────────────────┼───────────────┼───────────────────────────────┤
│ sync.Mutex Lock+Unlock          │ ~20-30ns      │ Atomic CAS + possible OS wait │
│ (uncontended)                   │               │                               │
└─────────────────────────────────┴───────────────┴───────────────────────────────┘
```

### When NOT to Use Channels

**"Channels orchestrate; mutexes serialize."** Use the right tool:

```
DON'T use channels for:
─────────────────────────
  ✗ Simple counter/accumulator       → use sync/atomic
  ✗ Protecting shared data structure → use sync.Mutex or sync.RWMutex
  ✗ One-time initialization          → use sync.Once
  ✗ Object pooling                   → use sync.Pool
  ✗ Read-heavy shared config         → use atomic.Value or sync.RWMutex

DO use channels for:
────────────────────
  ✓ Goroutine-to-goroutine communication (passing data/ownership)
  ✓ Signaling events (done, quit, ready)
  ✓ Coordinating pipeline stages
  ✓ Fan-out / fan-in patterns
  ✓ Rate limiting / semaphore
  ✓ Streaming results from producer to consumer
```

### The Contention Problem

A single channel used by 100 goroutines means 100 goroutines competing for
**one mutex** (`hchan.lock`). This serializes what should be concurrent work.

Solutions:
- **Shard channels**: instead of one channel, use N channels and hash-assign goroutines
- **Batch sends**: accumulate values locally, send a slice periodically
- **Use sync primitives**: for simple shared state, mutexes are faster

---

## 13. Common Bugs & Gotchas

### Bug 1: Goroutine Leak — Forgetting to Close

```go
// BUG: if the caller stops reading, the goroutine blocks forever on send
func generate() <-chan int {
    ch := make(chan int)
    go func() {
        for i := 0; ; i++ {
            ch <- i  // blocks forever if nobody reads
        }
    }()
    return ch
}

// FIX: accept a done/context for cancellation
func generate(ctx context.Context) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; ; i++ {
            select {
            case ch <- i:
            case <-ctx.Done():
                return
            }
        }
    }()
    return ch
}
```

### Bug 2: Send on Closed Channel — Panic

```go
// BUG: race between close and send
close(ch)
ch <- 42  // PANIC

// This happens in fan-out when multiple goroutines share a channel
// and one closes it while others are still sending.
// FIX: coordinate with WaitGroup (see Section 6)
```

### Bug 3: Range Over Unclosed Channel — Deadlock

```go
ch := make(chan int)
go func() {
    ch <- 1
    ch <- 2
    // forgot close(ch)!
}()

for v := range ch {  // blocks forever after receiving 2 — waiting for close
    fmt.Println(v)
}
// DEADLOCK: "all goroutines are asleep"
```

### Bug 4: time.After in a Loop — Timer Leak

```go
// BUG: each iteration creates a new timer that won't be GC'd until it fires
for {
    select {
    case v := <-ch:
        process(v)
    case <-time.After(5 * time.Second):  // NEW timer every iteration!
        return
    }
}

// FIX: reuse a timer
timer := time.NewTimer(5 * time.Second)
defer timer.Stop()
for {
    select {
    case v := <-ch:
        process(v)
        if !timer.Stop() {
            <-timer.C
        }
        timer.Reset(5 * time.Second)
    case <-timer.C:
        return
    }
}
```

### Bug 5: Sending Large Structs Through Channels — Hidden Copies

```go
type BigStruct struct {
    Data [1024]byte
    // ... many fields
}

ch := make(chan BigStruct)  // each send copies 1024+ bytes into the buffer!
ch <- bigValue              // memcopy into hchan.buf

// FIX: send a pointer (8 bytes copied instead of 1024)
ch := make(chan *BigStruct)
ch <- &bigValue
// BUT: now you share mutable state — the sender must not modify after sending
```

---

## 14. Quick Reference Card

```
HCHAN STRUCTURE
───────────────
runtime.hchan { buf, qcount, dataqsiz, sendx, recvx, sendq, recvq, lock, closed }
  └─ Channel variable is a pointer to hchan (8 bytes)
  └─ buf is a circular ring buffer (nil for unbuffered)
  └─ sendq/recvq are doubly-linked lists of sudog (parked goroutines)
  └─ Every operation acquires hchan.lock (runtime mutex)

MANY-TO-MANY
─────────────
  Multiple senders OK → parked on sendq in FIFO order
  Multiple receivers OK → parked on recvq in FIFO order
  Each value → exactly ONE receiver. No duplication, no loss.
  Worker pool = multiple goroutines ranging over same channel

HCHAN.LOCK (runtime mutex internals)
─────────────────────────────────────
  Phase 1: Spin — atomic CAS on mutex.key (~5-10ns if uncontended)
  Phase 2: Sleep — futex/semaphore (OS puts thread to sleep, ~1-2μs)
  Key: blocks the OS THREAD, not just the goroutine
  Why: channel code IS scheduler code, can't use scheduler to park
  Impact: high contention (100 goroutines, 1 channel) wastes OS threads

SEND PATHS (ch <- value)
────────────────────────
  1. Receiver waiting?   → direct copy to receiver's stack (fastest)
  2. Buffer has space?   → copy to buf[sendx]
  3. Neither?            → park goroutine on sendq

RECEIVE PATHS (v := <-ch)
─────────────────────────
  1. Sender waiting?     → direct copy from sender's stack
  2. Buffer has data?    → copy from buf[recvx]
  3. Neither?            → park goroutine on recvq

CLOSE SEMANTICS
───────────────
  close(ch)             → sets closed=1, wakes all waiters
  send to closed        → PANIC
  recv from closed      → returns buffered values, then zero values (ok=false)
  close nil             → PANIC
  close closed          → PANIC

NIL CHANNEL
───────────
  send to nil           → blocks forever (goroutine parked permanently)
  recv from nil         → blocks forever
  close nil             → PANIC
  select case on nil    → case is never selected (use for dynamic disable)

SELECT INTERNALS
────────────────
  1. Shuffle case order randomly (fairness)
  2. Lock all channels (address order to prevent deadlock)
  3. Poll each case (non-blocking check)
  4. If nothing ready & no default → park on all channels
  5. First channel ready → wake, dequeue from others, return case

DIRECTION TYPES (compile-time only, zero runtime cost)
──────────────────────────────────────────────────────
  chan T       bidirectional
  chan<- T     send-only (producer side)
  <-chan T     receive-only (consumer side)

PERFORMANCE RULES
─────────────────
  Channels orchestrate, mutexes serialize
  Buffered send/recv: ~50-100ns (mutex + memcopy)
  Unbuffered pair:    ~200-300ns (mutex + context switch)
  Select with N:      ~N × 150ns (N mutexes + shuffle)
  Use atomic for counters, mutex for shared state, channels for communication

TOOLS
─────
  GODEBUG=schedtrace=1000    # goroutine scheduler trace
  go tool trace              # visual timeline of goroutines, channels, GC
  go test -race ./...        # data race detection — NON-NEGOTIABLE
  runtime.NumGoroutine()     # detect goroutine leaks in tests
```

---

## One-Line Summary

> A channel is a pointer to `runtime.hchan` — a mutex-protected struct with a circular
> ring buffer and two wait queues (sendq/recvq). Sends and receives follow three paths
> (direct transfer → buffer → park), select multiplexes by parking on all channels
> simultaneously, and nil channels block forever — a feature, not a bug, for dynamic
> select control.

---

## 15. Select in Practice — Quick Reference

`select` is Go's multiplexer — it lets a goroutine wait on multiple channel
operations simultaneously. Understanding three things covers every select use case:

1. **One-shot vs loop**: `select` runs ONCE. You wrap it in `for` to make an event loop.
2. **Blocking vs non-blocking**: adding `default` turns a blocking wait into a one-time poll.
3. **Fairness**: when multiple cases are ready, the runtime picks one at random (Fisher-Yates
   shuffle inside `selectgo`), preventing starvation. If you need priority, use nested selects
   (non-blocking high-priority check first, then blocking wait for any).

```
  ┌──────────────────────────────────────────────────────────────────────┐
  │                      SELECT EXECUTION SUMMARY                       │
  │                                                                     │
  │  1. select runs ONCE per encounter (like switch, not like for)      │
  │  2. If multiple cases ready → pick one at RANDOM (fairness)         │
  │  3. If no cases ready + default → execute default, exit             │
  │  4. If no cases ready + no default → BLOCK (gopark on all channels) │
  │  5. When blocked, first channel that becomes ready → wake + execute │
  │                                                                     │
  │  default makes it non-blocking:                                     │
  │    select + default = tryReceive/trySend (check once, don't wait)   │
  │    select - default = blocking wait (sleep until something happens) │
  │                                                                     │
  │  for-select = event loop (continuous listening):                    │
  │    for { select { ... } }  = keep checking, handle events forever  │
  │    for + select + default  = BUSY SPIN (danger! burns CPU)          │
  │    for + select - default  = EFFICIENT event loop (sleeps when idle)│
  └──────────────────────────────────────────────────────────────────────┘
```
