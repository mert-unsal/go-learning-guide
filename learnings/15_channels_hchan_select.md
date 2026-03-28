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
6. [Close Semantics — Under the Hood](#6-close-semantics--under-the-hood)
7. [Select Statement Internals](#7-select-statement-internals)
8. [Nil Channel Behavior — Why It Blocks Forever](#8-nil-channel-behavior--why-it-blocks-forever)
9. [Channel Direction Types](#9-channel-direction-types)
10. [The `chan struct{}` Pattern — Why Not `chan bool`](#10-the-chan-struct-pattern--why-not-chan-bool)
11. [Production Patterns](#11-production-patterns)
12. [Performance Characteristics & When NOT to Use Channels](#12-performance-characteristics--when-not-to-use-channels)
13. [Common Bugs & Gotchas](#13-common-bugs--gotchas)
14. [Quick Reference Card](#14-quick-reference-card)
15. [Select in Practice — Patterns, Execution Model & Production Scenarios](#15-select-in-practice--patterns-execution-model--production-scenarios)

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

## 5b. Visual Guide: Channel Operations End-to-End

The following diagrams show the complete lifecycle of channel operations,
combining hchan, ring buffer, sudog, and goroutine state transitions.

### Scenario 1: Buffered Channel — Happy Path (No Blocking)

```
   ch := make(chan int, 3)      ← create channel with buffer of 3
   ch <- 10                     ← send (buffer has space)
   ch <- 20                     ← send (buffer has space)
   v := <-ch                    ← receive (buffer has data)

   ═══════════════════════════════════════════════════════════════
   AFTER make(chan int, 3):

   ch ──► hchan {
            qcount: 0          ← empty
            dataqsiz: 3        ← capacity
            buf: ──► [ _ | _ | _ ]
            sendx: 0   recvx: 0
            sendq: ∅   recvq: ∅
            closed: 0
          }

   ═══════════════════════════════════════════════════════════════
   AFTER ch <- 10:

   hchan {
     qcount: 1                 ← one element
     buf: ──► [ 10 | _ | _ ]
               ▲recvx  ▲sendx=1
     sendq: ∅   recvq: ∅      ← nobody waiting, no sudogs needed
   }

   G1 (sender) continues immediately — NOT blocked.

   ═══════════════════════════════════════════════════════════════
   AFTER ch <- 20:

   hchan {
     qcount: 2
     buf: ──► [ 10 | 20 | _ ]
               ▲recvx     ▲sendx=2
     sendq: ∅   recvq: ∅
   }

   ═══════════════════════════════════════════════════════════════
   AFTER v := <-ch:

   hchan {
     qcount: 1
     buf: ──► [ _ | 20 | _ ]
                    ▲recvx=1  ▲sendx=2
     sendq: ∅   recvq: ∅
   }

   v = 10 (the oldest value — FIFO order)
   G2 (receiver) continues immediately.
```

### Scenario 2: Unbuffered Channel — Direct Transfer

```
   ch := make(chan int)         ← unbuffered (no buffer at all)

   ═══════════════════════════════════════════════════════════════
   GOROUTINE G1: v := <-ch     (receiver arrives first, nobody sending)

   Step 1: G1 checks recvq → sendq is empty, no buffer → must BLOCK

   Step 2: Runtime gets sudog from P's cache:
     P2's cache: [s4, s5, s6] → pop s4
     Fill in: s4 = { g: G1, elem: &v, c: ch }

   Step 3: Enqueue s4 onto ch.recvq:

   hchan {
     dataqsiz: 0          ← unbuffered
     buf: nil
     recvq: ──► ┌────────┐
                │ sudog  │──► nil
                │ g: G1  │
                │ elem:&v│ ← points to v on G1's STACK
                └────────┘
     sendq: ∅
   }

   Step 4: gopark(G1) — G1 state: _Grunning → _Gwaiting
           G1 is now sleeping. M (OS thread) picks another goroutine.

   ═══════════════════════════════════════════════════════════════
   LATER — GOROUTINE G2: ch <- 42     (sender arrives)

   Step 1: G2 locks hchan.lock
   Step 2: G2 checks recvq → found s4 (G1 is waiting!)
   Step 3: DIRECT COPY — bypasses buffer entirely:

     G2's stack                G1's stack (frozen — G1 is parked)
     ┌──────────────┐          ┌──────────────┐
     │  val = 42    │ ─────►  │  v = 42      │
     └──────────────┘ typedmemmove()  └──────────────┘
                      (copies 8 bytes directly from stack to stack)

   Step 4: Dequeue s4 from recvq. Return s4 to P's cache.
   Step 5: goready(G1) — G1 state: _Gwaiting → _Grunnable
           G1 is placed on a run queue, will resume soon.
   Step 6: G2 unlocks hchan.lock. G2 continues immediately.

   ═══════════════════════════════════════════════════════════════
   G1 WAKES UP:

   G1 resumes execution at the point after <-ch.
   v is already 42 on its stack (written directly by G2).
   G1 continues with v = 42. No buffer was involved.
```

### Scenario 3: Buffered Channel Full — Sender Blocks, Then Receiver Arrives

```
   ch := make(chan int, 2)
   ch <- 10   ← buffer: [10, _]
   ch <- 20   ← buffer: [10, 20] — FULL

   ═══════════════════════════════════════════════════════════════
   GOROUTINE G1: ch <- 30     (sender, but buffer is full!)

   Step 1: G1 checks — no waiting receivers, buffer FULL → must BLOCK

   Step 2: Get sudog from P's cache:
     s1 = { g: G1, elem: &(30), c: ch }

   Step 3: Enqueue s1 onto sendq:

   hchan {
     qcount: 2 / dataqsiz: 2    ← buffer is full
     buf: ──► [ 10 | 20 ]
               ▲recvx  ▲sendx=0 (wrapped!)
     sendq: ──► ┌────────┐
                │ sudog  │──► nil
                │ g: G1  │
                │ elem:30│ ← the value G1 wants to send
                └────────┘
     recvq: ∅
   }

   Step 4: gopark(G1) — G1 sleeps.

   ═══════════════════════════════════════════════════════════════
   GOROUTINE G2: v := <-ch    (receiver arrives!)

   This case is special — buffer is full AND a sender is waiting.
   The runtime does a clever optimization:

   Step 1: G2 takes the OLDEST value from the buffer:
     v = buf[recvx] = buf[0] = 10

   Step 2: G1's value (30) fills the slot that was just freed:
     buf[recvx] = 30  (from s1.elem)

   Step 3: Advance recvx:

   hchan {
     qcount: 2 / dataqsiz: 2    ← still full!
     buf: ──► [ 30 | 20 ]       ← 30 replaced 10's slot
                    ▲recvx=1 ▲sendx=1
     sendq: ∅                    ← s1 removed
     recvq: ∅
   }

   Step 4: Wake G1, return sudog to cache.

   Result:
     v = 10 (G2 got the OLDEST value — FIFO maintained!)
     G1 unblocked (its value 30 is now in the buffer)
     Buffer is still full: [30, 20] — the receiver took one, sender added one
```

### Scenario 4: Select with Multiple Channels — Goroutine on Multiple Wait Queues

```
   select {
   case v := <-ch1:     // case 0
   case v := <-ch2:     // case 1
   case ch3 <- 42:      // case 2
   }

   None of the channels are ready. Goroutine G5 must block on ALL of them.

   ═══════════════════════════════════════════════════════════════
   Step 1: Runtime creates 3 sudogs (from P's cache):
     s1 = { g: G5, c: ch1 }    ← for receive on ch1
     s2 = { g: G5, c: ch2 }    ← for receive on ch2
     s3 = { g: G5, c: ch3 }    ← for send on ch3

   Step 2: Enqueue each on its channel:

   ch1.recvq: [..., s1{g:G5}]    ← G5 is waiting to receive from ch1
   ch2.recvq: [..., s2{g:G5}]    ← G5 is waiting to receive from ch2
   ch3.sendq: [..., s3{g:G5}]    ← G5 is waiting to send to ch3

   Step 3: gopark(G5) — G5 sleeps on ALL three channels simultaneously.

   ═══════════════════════════════════════════════════════════════
   LATER: Someone sends a value to ch2.

   ch2 processing:
     1. Dequeue s2 from ch2.recvq
     2. Copy value to G5's variable
     3. Wake G5

   G5 wakes up and CLEANS UP the other sudogs:
     4. Remove s1 from ch1.recvq  ← G5 is no longer waiting on ch1
     5. Remove s3 from ch3.sendq  ← G5 is no longer waiting on ch3
     6. Return s1, s2, s3 to P's cache

   selectgo() returns case index 1 (ch2 was ready).

   ═══════════════════════════════════════════════════════════════
   VISUAL: G5 parked on 3 channels, then woken by ch2:

   BEFORE (G5 sleeping):

   ch1.recvq: ─── ... ─── s1{G5} ───
   ch2.recvq: ─── ... ─── s2{G5} ───   ← all three queues have G5
   ch3.sendq: ─── ... ─── s3{G5} ───

           ┌────────────────┐
           │  G5 (_Gwaiting)│   sleeping, waiting for ANY channel
           └────────────────┘

   AFTER (ch2 delivers value, G5 wakes):

   ch1.recvq: ─── ... ───           ← s1 removed (cleanup)
   ch2.recvq: ─── ... ───           ← s2 removed (delivered)
   ch3.sendq: ─── ... ───           ← s3 removed (cleanup)

           ┌──────────────────┐
           │  G5 (_Grunnable) │   awake, processing case 1 (ch2)
           └──────────────────┘
```

### Scenario 5: Close Channel with Waiting Receivers

```
   ch := make(chan int, 3)
   ch <- 10
   ch <- 20
   // 2 goroutines are waiting to receive:

   hchan {
     qcount: 2
     buf: ──► [ 10 | 20 | _ ]
     recvq: ──► s1{G3} ──► s2{G7} ──► nil
     closed: 0
   }

   ═══════════════════════════════════════════════════════════════
   close(ch):

   Step 1: Set closed = 1

   Step 2: Wake ALL receivers in recvq:
     G3: receives 10 from buffer (ok=true, buffer had data)
     G7: receives 20 from buffer (ok=true, buffer had data)

   Step 3: If there were MORE receivers than buffered values:
     Remaining receivers get (zero value, ok=false)

   Step 4: If there were senders in sendq:
     All senders PANIC: "send on closed channel"

   ═══════════════════════════════════════════════════════════════
   AFTER close:

   hchan {
     qcount: 0
     buf: ──► [ _ | _ | _ ]       ← drained
     recvq: ∅                      ← all receivers woken
     sendq: ∅
     closed: 1                     ← permanently closed
   }

   Any future <-ch returns (0, false) immediately.
   Any future ch <- val PANICS.
```

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

### Pattern 1: Done / Cancellation (replaced by `context.Context`)

```go
// Legacy pattern (pre-context)
done := make(chan struct{})
go worker(done)
// ... later ...
close(done)  // signal worker to stop

// Modern pattern — use context
ctx, cancel := context.WithCancel(context.Background())
go worker(ctx)
// ... later ...
cancel()  // signals ctx.Done() channel
```

Under the hood, `context.WithCancel` creates a `chan struct{}` internally.
`ctx.Done()` returns it. `cancel()` calls `close()` on it. Same mechanism,
better API.

### Pattern 2: Worker Pool with Backpressure

```go
func workerPool(numWorkers int, jobs <-chan Job) <-chan Result {
    results := make(chan Result)
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {    // range blocks until jobs is closed
                results <- process(job)
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

**Backpressure**: The `jobs` channel buffer size controls how far ahead the
producer can get. When the buffer is full, the producer blocks — this is
natural backpressure without explicit rate limiting.

### Pattern 3: Semaphore (Bounded Concurrency)

```go
sem := make(chan struct{}, maxConcurrent)

for _, item := range items {
    sem <- struct{}{}         // acquire — blocks when maxConcurrent reached
    go func(it Item) {
        defer func() { <-sem }()  // release
        process(it)
    }(item)
}

// Wait for all to finish
for i := 0; i < maxConcurrent; i++ {
    sem <- struct{}{}
}
```

The buffered channel acts as a counting semaphore. Buffer capacity = max
concurrent goroutines. This is used in production to limit concurrent DB
connections, HTTP requests, or file handles.

### Pattern 4: Fan-Out / Fan-In

```
  Fan-Out: one producer, N consumers

       ┌──► worker 1 ──┐
  jobs ┼──► worker 2 ──┼──► merged results
       └──► worker 3 ──┘

  Fan-In: N producers, one consumer (merge)
```

Fan-out is just multiple goroutines reading from the same channel (the runtime
fairly distributes values). Fan-in merges multiple channels into one (see your
`Merge` exercise).

### Pattern 5: Pipeline

```go
// Each stage: goroutine reading from input channel, writing to output channel
func stage(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for v := range in {
            out <- transform(v)
        }
    }()
    return out
}

// Compose: input → stage1 → stage2 → stage3 → output
result := stage3(stage2(stage1(input)))
```

Each stage is a goroutine connected by channels. Cancellation flows backward
through the pipeline: when the consumer stops reading, buffered channels fill up,
and producers eventually block. For eager cancellation, thread a `context.Context`
or `done` channel through each stage.

### Pattern 6: Timeout with Context (Production-Grade)

```go
// DON'T use time.After in loops (it leaks timers)
// DO use context.WithTimeout

func fetchWithTimeout(ctx context.Context, url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()  // ALWAYS cancel to release resources

    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := http.DefaultClient.Do(req)
    // ...
}
```

**Why not `time.After`?** In a loop, each `time.After` creates a timer that
lives until it fires — even if the select case was not chosen. This leaks timers.
`context.WithTimeout` + `defer cancel()` cleans up immediately.

### Pattern 7: Or-Done Channel (Elegant Cancellation)

#### The Problem: Blocking Outside of Select

Imagine a pipeline stage that reads from `in` and sends to `out`:

```go
  // Attempt 1: single select — BROKEN
  for {
      select {
      case <-ctx.Done():
          return                   // ✅ cancel detected while READING
      case v, ok := <-in:
          if !ok { return }
          out <- v                 // ❌ BLOCKS HERE if nobody reads 'out'
      }                            //    ctx.Done() is NOT checked!
  }
```

The bug: after receiving `v` from `in`, we leave the select and do a plain
`out <- v`. If downstream stopped reading from `out`, this send blocks forever.
We're no longer inside a `select`, so `ctx.Done()` is invisible.

```
  Timeline of the goroutine leak:

  1. select: picks case 2 (value arrived from 'in')     ← leave select
  2. out <- v                                             ← BLOCKED here
  3. ctx gets cancelled                                   ← nobody checks!
  4. goroutine stuck forever on step 2                    ← LEAKED
```

#### The Solution: Double Select

Every blocking channel operation (both read AND write) must be inside a
`select` that also watches `ctx.Done()`:

```go
// orDone wraps a channel read to respect cancellation
func orDone(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case <-ctx.Done():             // OUTER select: cancel OR read
                return
            case v, ok := <-in:
                if !ok { return }
                select {
                case out <- v:             // INNER select: cancel OR write
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}
```

#### Every Scenario Traced

```
  Scenario 1: Normal flow (data available, downstream ready)
  ──────────────────────────────────────────────────────────
  outer select → case 2: v arrives from 'in'
  inner select → case 1: out <- v succeeds (someone is reading)
  → loop continues ✅

  Scenario 2: Context cancelled while waiting to READ
  ──────────────────────────────────────────────────
  outer select → case 1: ctx.Done() fires
  → return, goroutine exits cleanly ✅

  Scenario 3: Context cancelled while waiting to WRITE ← the bug we fixed
  ──────────────────────────────────────────────────────
  outer select → case 2: v arrives from 'in'
  inner select → case 2: ctx.Done() fires (out is blocked, nobody reading)
  → return, goroutine exits cleanly ✅

  Scenario 4: Input channel closes
  ─────────────────────────────────
  outer select → case 2: ok = false (in was closed)
  → return, goroutine exits cleanly ✅
```

```
  Visual — Single Select vs Double Select:

  Single select (BROKEN):              Double select (CORRECT):

  ┌──────────────────────┐             ┌──────────────────────┐
  │ select {             │             │ select {             │ ← outer
  │   case <-ctx.Done() │             │   case <-ctx.Done() │
  │   case v := <-in    │             │   case v := <-in ───────┐
  │ }                    │             │ }                    │  │
  │                      │             │                      │  ↓
  │ out <- v  ← DANGER! │             │   ┌────────────────┐ │
  │ (no cancel check)   │             │   │ select {       │ │ ← inner
  └──────────────────────┘             │   │  case out <- v │ │
                                       │   │  case <-ctx   │ │
                                       │   │ }              │ │
                                       │   └────────────────┘ │
                                       └──────────────────────┘
```

#### Why It's Called "Or-Done"

The name means: forward values from `in` **or** stop when done. It's a
**composable building block** — wrap any channel once, and downstream code
can use simple `range`:

```go
  // Without orDone: every pipeline stage needs double-select boilerplate
  // With orDone: wrap once, range simply

  out := orDone(ctx, someChannel)
  for v := range out {
      // if ctx cancels → out closes → range exits naturally
      // no double-select needed here
      process(v)
  }
```

`orDone` is NOT in the standard library. It's a community pattern from
*"Concurrency in Go"* (Katherine Cox-Buday, O'Reilly 2017). Go's stdlib
gives you the primitives (`chan`, `select`, `context`). You compose them.
This is deliberate — *"A little copying is better than a little dependency."*

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

## 15. Select in Practice — Patterns, Execution Model & Production Scenarios

### Select Executes ONCE — It Is NOT a Loop

This is the most common misconception. A `select` statement is like a `switch` — it
evaluates once, picks one case, executes it, and is done. **There is no implicit looping.**

```go
select {
case v := <-ch1:
    fmt.Println("got from ch1:", v)
case v := <-ch2:
    fmt.Println("got from ch2:", v)
default:
    fmt.Println("nothing ready")
}
fmt.Println("select is DONE — execution continues here")
```

```
  Execution flow:
  
  ┌──────────────────────────────────────┐
  │            select { ... }            │
  │                                      │
  │  Are any cases ready?                │
  │  ├── YES → pick one (random if many) │
  │  │         execute its body          │
  │  │         EXIT select               │
  │  └── NO  → is there a default?       │
  │       ├── YES → execute default body │
  │       │         EXIT select          │
  │       └── NO  → BLOCK (gopark)       │
  │                 wait until a case     │
  │                 becomes ready         │
  │                 execute its body      │
  │                 EXIT select          │
  └──────────────────────────────────────┘
  │
  ▼
  Next line of code runs here
```

### The `for-select` Loop — YOU Add the Loop

When you want continuous listening, you wrap `select` in a `for` loop yourself:

```go
for {
    select {
    case v := <-ch:
        process(v)
    case <-ctx.Done():
        return   // breaks out of BOTH select and for
    }
}
```

**Critical distinction:**

```
  select { ... }              → runs ONCE, like switch
  for { select { ... } }     → runs FOREVER, you must break/return to stop
  for v := range ch { ... }  → runs until channel is closed (simpler for single channel)
```

**When to use which:**

```
  ┌─────────────────────────────────┬─────────────────────────────────────────────┐
  │ Pattern                         │ When to Use                                 │
  ├─────────────────────────────────┼─────────────────────────────────────────────┤
  │ v := <-ch                       │ Single blocking receive, you know data      │
  │                                 │ will come, no cancellation needed            │
  ├─────────────────────────────────┼─────────────────────────────────────────────┤
  │ for v := range ch               │ Drain a single channel until it closes      │
  │                                 │ (simplest consumer pattern)                  │
  ├─────────────────────────────────┼─────────────────────────────────────────────┤
  │ select { cases... }             │ ONE-TIME check across multiple channels     │
  │ (no loop)                       │ or non-blocking tryReceive/trySend          │
  ├─────────────────────────────────┼─────────────────────────────────────────────┤
  │ for { select { cases... } }     │ Continuous listening on multiple channels   │
  │                                 │ with cancellation — the EVENT LOOP pattern  │
  └─────────────────────────────────┴─────────────────────────────────────────────┘
```

### Random Selection — Why and How

When multiple cases are ready simultaneously, `select` picks one **uniformly at random**:

```go
ch1 := make(chan string, 1)
ch2 := make(chan string, 1)
ch1 <- "from ch1"
ch2 <- "from ch2"

// Both cases are ready — which one runs?
select {
case v := <-ch1:
    fmt.Println(v)    // ~50% of the time
case v := <-ch2:
    fmt.Println(v)    // ~50% of the time
}
```

**Under the hood** (from `selectgo` in `runtime/select.go`):

```
  Phase 1: SHUFFLE — Fisher-Yates shuffle on case order
  
  Cases in source order:  [ch1, ch2]
  After shuffle:          [ch2, ch1]   ← random permutation
  
  Phase 3: POLL — walk cases in SHUFFLED order
  
  First shuffled case (ch2) → ready? YES → execute it, done.
  
  Next run, shuffle might produce [ch1, ch2]:
  First shuffled case (ch1) → ready? YES → execute it, done.
```

**Why random, not first-case-wins?**

```
  If select always chose the first ready case:
  
  for {
      select {
      case v := <-highPriority:     // always checked first
          process(v)
      case v := <-lowPriority:      // STARVED if highPriority is always busy
          process(v)
      }
  }
  
  With a busy highPriority channel, lowPriority would NEVER be selected.
  Random selection guarantees fairness — both channels get served.
```

**What if you WANT priority?** Use nested selects:

```go
for {
    // Try high priority first (non-blocking)
    select {
    case v := <-highPriority:
        process(v)
        continue       // go back to top, check high priority again
    default:
    }
    
    // Nothing high-priority — wait for anything
    select {
    case v := <-highPriority:
        process(v)
    case v := <-lowPriority:
        process(v)
    case <-ctx.Done():
        return
    }
}
```

### The `default` Case — Behavior and Dangers

`default` turns `select` from a blocking operation into a **non-blocking poll**:

```
  WITHOUT default:                     WITH default:
  ──────────────────                   ──────────────────
  select {                             select {
  case v := <-ch:                      case v := <-ch:
      process(v)                           process(v)
  }                                    default:
  // goroutine SLEEPS here                 // runs immediately
  // until ch has data                 }
  // (gopark → wait queue)             // goroutine NEVER sleeps
                                       // checks once, moves on
```

**The `for-select-default` danger — busy spin:**

```go
// ❌ CPU KILLER — this burns 100% CPU on one core
for {
    select {
    case v := <-ch:
        process(v)
    default:
        // nothing ready — loop immediately checks again
        // millions of iterations per second, all doing nothing
    }
}

// ✅ FIX 1: Remove default — let select block
for {
    select {
    case v := <-ch:
        process(v)
    case <-ctx.Done():
        return
    }
    // goroutine sleeps until ch or ctx.Done() is ready — zero CPU
}

// ✅ FIX 2: If you need default, add a sleep or yield
for {
    select {
    case v := <-ch:
        process(v)
    default:
        time.Sleep(10 * time.Millisecond)   // back off, don't spin
    }
}

// ✅ FIX 3: Use default only for ONE-TIME checks, not in loops
select {
case v := <-ch:
    process(v)
default:
    fallback()     // runs once, no loop
}
```

### Select Execution Model — Complete Mental Model

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
