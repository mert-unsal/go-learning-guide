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

**Key insight:** sudogs are **pooled per P** (logical processor), not heap-allocated
each time. This keeps channel operations allocation-free in the common case.

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

```go
// orDone wraps a channel read to respect cancellation
func orDone(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case <-ctx.Done():
                return
            case v, ok := <-in:
                if !ok { return }
                select {
                case out <- v:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}
```

This pattern wraps any channel to make it cancellation-aware. Without it,
a goroutine blocked on `<-in` won't notice cancellation until the next value
arrives.

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
  └─ Every operation acquires hchan.lock (mutex)

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
