# Deep Dive: Goroutine Scheduler Internals — The GMP Model

> How Go multiplexes millions of goroutines onto OS threads: the G-M-P model,
> work stealing, preemption, syscall handling, and the network poller.

---

## Table of Contents

1. [Goroutines vs OS Threads](#1-goroutines-vs-os-threads)
2. [The GMP Model](#2-the-gmp-model)
3. [The `runtime.g` Struct — Goroutine](#3-the-runtimeg-struct--goroutine)
4. [The `runtime.p` Struct — Processor](#4-the-runtimep-struct--processor)
5. [The `runtime.m` Struct — Machine](#5-the-runtimem-struct--machine)
6. [Scheduling Algorithm](#6-scheduling-algorithm)
7. [Work Stealing](#7-work-stealing)
8. [Preemption — Cooperative and Asynchronous](#8-preemption--cooperative-and-asynchronous)
9. [Syscall Handling — The Hand-Off](#9-syscall-handling--the-hand-off)
10. [Network Poller — Scalable I/O](#10-network-poller--scalable-io)
11. [Stack Management](#11-stack-management)
12. [Goroutine Lifecycle — Birth to Death](#12-goroutine-lifecycle--birth-to-death)
13. [GOMAXPROCS and Tuning](#13-gomaxprocs-and-tuning)
14. [Performance Characteristics](#14-performance-characteristics)
15. [Debugging Tools](#15-debugging-tools)
16. [Quick Reference Card](#16-quick-reference-card)

---

## 1. Goroutines vs OS Threads

Go doesn't map one goroutine to one OS thread. It uses **M:N scheduling** —
multiplexing M goroutines onto N OS threads.

```
┌─────────────────────────────────────────────────────────────────┐
│                    Goroutine vs OS Thread                        │
├───────────────────┬──────────────────┬──────────────────────────┤
│ Property          │ OS Thread        │ Goroutine                │
├───────────────────┼──────────────────┼──────────────────────────┤
│ Stack size        │ ~1-8 MB (fixed)  │ ~2-8 KB (growable)       │
│ Creation cost     │ ~1-10 μs         │ ~0.3 μs                  │
│ Context switch    │ ~1-5 μs (kernel) │ ~100-200 ns (user-space) │
│ Scheduling        │ Kernel (OS)      │ Go runtime (user-space)  │
│ Memory per 10k    │ ~10-80 GB        │ ~20-80 MB                │
│ Max practical     │ ~10k per process │ Millions                 │
│ Identity          │ OS thread ID     │ goid (not exposed)       │
└───────────────────┴──────────────────┴──────────────────────────┘
```

### Why Is a Context Switch Cheaper?

An OS thread context switch requires:
- Trap into kernel mode (ring 3 → ring 0)
- Save/restore all CPU registers, FPU state, TLB
- Kernel scheduler decision
- Return to user mode

A goroutine context switch only:
- Saves ~15 registers into `g.sched` (a `gobuf` struct)
- Updates a few pointers (current G on current M)
- Stays entirely in user space — no kernel involvement

```
  OS Thread Switch:                    Goroutine Switch:
  ┌──────────┐                        ┌──────────┐
  │ User     │ ←─ save regs           │ Save     │ ←─ ~15 registers
  │ → Kernel │ ←─ kernel scheduler    │ → Update │ ←─ swap G pointer
  │ → User   │ ←─ restore regs       │ → Run    │ ←─ restore regs
  └──────────┘                        └──────────┘
  ~1-5 μs                              ~100-200 ns
```

---

## 2. The GMP Model

The scheduler is built on three core entities:

```
  G = Goroutine      (the unit of work)
  M = Machine         (OS thread that executes code)
  P = Processor       (scheduling context with a run queue)
```

### The Relationship

```
                          ┌─────────────────────────────┐
                          │       Go Process             │
                          │                              │
  ┌─────────┐   ┌────────┴────────┐   ┌─────────┐      │
  │    P0   │   │       P1        │   │    P2   │      │
  │ ┌─────┐ │   │  ┌─────┐        │   │ ┌─────┐ │      │
  │ │ LRQ │ │   │  │ LRQ │        │   │ │ LRQ │ │      │
  │ │G G G│ │   │  │G G  │        │   │ │G    │ │      │
  │ └─────┘ │   │  └─────┘        │   │ └─────┘ │      │
  │    ↕    │   │     ↕           │   │    ↕    │      │
  │    M0   │   │     M1          │   │    M2   │      │
  │ (thread)│   │  (thread)       │   │ (thread)│      │
  └─────────┘   └─────────────────┘   └─────────┘      │
       │              │                    │             │
       │      ┌───────┴────────┐           │             │
       │      │  Global Run    │           │             │
       │      │  Queue (GRQ)   │           │             │
       │      │  G G G G G     │           │             │
       │      └────────────────┘           │             │
       └───────────────┴───────────────────┘             │
                          │                              │
                          └─────────────────────────────┘
```

**Key rules:**
- A G must be assigned to a P's run queue to execute
- A P must be attached to an M to run goroutines
- An M without a P is idle (parked or in syscall)
- Number of Ps = `GOMAXPROCS` (default: CPU core count)
- Number of Ms can exceed Ps (threads blocked in syscalls)

**Source:** `runtime/runtime2.go`, `runtime/proc.go`

---

## 3. The `runtime.g` Struct — Goroutine

Each goroutine is represented by a `runtime.g` struct (~400 bytes).

```go
// runtime/runtime2.go (simplified)
type g struct {
    stack       stack        // lo, hi — current stack bounds
    stackguard0 uintptr     // compared with SP for stack growth check
    m           *m           // the M currently running this G (nil if not running)
    sched       gobuf        // saved registers for context switch
    atomicstatus atomic.Uint32 // goroutine state
    goid         uint64      // goroutine ID (not exposed to user code)
    preempt      bool        // preemption signal
    // ... many more fields
}

type gobuf struct {
    sp   uintptr  // stack pointer
    pc   uintptr  // program counter (instruction pointer)
    g    guintptr // pointer back to this G
    ret  uintptr  // return value
    bp   uintptr  // base pointer (for frame pointer unwinding)
    // ... context, lr on ARM
}
```

### Goroutine States

```
                    ┌──────────┐
           ┌───────│ _Gidle   │ (just allocated, not yet used)
           │       └────┬─────┘
           │            │ newproc()
           │            ▼
           │       ┌──────────┐  ◄── schedule() picks this G
           │  ┌───►│_Grunnable│──────────────────────────────┐
           │  │    └──────────┘                              │
           │  │         │ execute()                           │
           │  │         ▼                                     │
           │  │    ┌──────────┐                              │
           │  │    │_Grunning │ ◄── G is executing on an M   │
           │  │    └──┬───┬───┘                              │
           │  │       │   │                                   │
           │  │       │   │ channel op / lock /               │
           │  │       │   │ I/O / sleep                       │
           │  │       │   ▼                                   │
           │  │       │ ┌──────────┐                          │
           │  │       │ │_Gwaiting │ (parked on something)    │
           │  │       │ └────┬─────┘                          │
           │  │       │      │ goready() — thing happened     │
           │  │       │      └────────────────────────────────┘
           │  │       │                     back to _Grunnable
           │  │       │ syscall
           │  │       ▼
           │  │  ┌──────────┐
           │  │  │_Gsyscall │ (blocked in OS syscall)
           │  │  └────┬─────┘
           │  │       │ syscall returns
           │  │       └───► _Grunnable (re-enqueue)
           │  │
           │  │  goexit() — function returned
           │  │       │
           │  │       ▼
           │  │  ┌──────────┐
           │  └──│ _Gdead   │ (finished, put on free list for reuse)
           │     └──────────┘
           │          │ reused by newproc()
           └──────────┘
```

---

## 4. The `runtime.p` Struct — Processor

P is the **scheduling context**. It holds the local run queue and per-P resources.

```go
// runtime/runtime2.go (simplified)
type p struct {
    id          int32
    status      uint32       // _Pidle, _Prunning, _Psyscall, _Pgcstop, _Pdead
    m           muintptr     // the M currently attached to this P
    mcache      *mcache      // per-P memory cache (no lock needed!)

    // Local run queue (LRQ) — lock-free circular buffer
    runqhead uint32
    runqtail uint32
    runq     [256]guintptr   // fixed-size array of 256 G pointers

    runnext  guintptr        // next G to run (fast path, skip queue)

    // Timers (time.After, time.Sleep, etc.)
    timers   timers

    // GC state
    gcBgMarkWorker guintptr
    // ... more fields
}
```

### The Local Run Queue (LRQ)

```
  P.runq — circular buffer of 256 slots:

  ┌───┬───┬───┬───┬───┬───┬───┬───┬─── ─── ───┬───┐
  │G12│G13│G14│G15│   │   │   │   │  ...       │   │
  └───┴───┴───┴───┴───┴───┴───┴───┴─── ─── ───┴───┘
   ^head              ^tail
   (next to dequeue)  (next to enqueue)

  Operations:
  - runqput():  enqueue G at tail (CAS, lock-free)
  - runqget():  dequeue G from head (CAS, lock-free)
  - runqsteal(): steal half from another P (CAS)

  When LRQ is full (256 Gs):
  → runtime moves HALF of the queue to the Global Run Queue (GRQ)
  → this is called "overflow" — prevents one P from hoarding
```

### The `runnext` Fast Path

`P.runnext` holds a single G that will be scheduled **next**, bypassing the queue.
When `go func()` is called, the new goroutine is placed in `runnext` — this gives
freshly created goroutines priority (good for producer-consumer patterns where the
producer immediately creates work for the consumer).

### Why P Exists (The Design Decision)

Before Go 1.1, there was no P — just G and M. Problems:
- Every goroutine operation required a global mutex on the run queue
- Memory caches were per-M, but Ms could be parked → cache wasted

P solves both: each P has a lock-free local queue and its own mcache.
The global queue is only touched when a P overflows or is empty.

---

## 5. The `runtime.m` Struct — Machine

M represents an OS thread. The runtime creates Ms as needed.

```go
// runtime/runtime2.go (simplified)
type m struct {
    g0        *g       // special goroutine for scheduling code
    curg      *g       // the user goroutine currently running
    p         puintptr // the P attached to this M (nil if no P)
    nextp     puintptr // P to attach on wakeup
    oldp      puintptr // P that was attached before syscall
    spinning  bool     // is this M looking for work?
    // ... thread ID, signal handling, etc.
}
```

### The `g0` — Scheduler Stack

Every M has a special `g0` goroutine that runs scheduler code. When a goroutine
needs to be scheduled (channel block, syscall, preemption), execution switches
from the user G's stack to `g0`'s stack to run `runtime.schedule()`.

```
  Normal execution:         Scheduling:
  ┌───────────┐             ┌───────────┐
  │ User G    │             │    g0     │
  │ (your     │ ──switch──► │ (scheduler│
  │  code)    │             │  code)    │
  │ stack:    │             │ stack:    │
  │ 2KB-1GB   │             │ ~8KB      │
  └───────────┘             └───────────┘
```

### M Lifecycle

```
  Need a new M?
  │
  ├─ Check sleeping M list (runtime.midle)
  │   └─ Found? → wake it up, attach P → done
  │
  └─ Not found? → create new OS thread (runtime.newm → pthread_create)
      └─ New M starts running runtime.mstart → schedule loop
```

The runtime doesn't eagerly create threads. It creates them on demand and
reuses parked ones. There's a hard cap of 10,000 Ms by default
(`runtime.SetMaxThreads()`).

---

## 6. Scheduling Algorithm

The core loop is `runtime.schedule()` — it runs on `g0`'s stack on every M.

```
runtime.schedule():
│
├─ 1. Is there a runnext G on my P?
│      └─ YES → run it immediately (fast path)
│
├─ 2. Check local run queue (P.runq)
│      └─ Non-empty → dequeue from head → run it
│
├─ 3. Every 61st schedule call: check Global Run Queue (GRQ)
│      └─ Non-empty → grab a batch → run one
│      (61 is prime — avoids synchronization patterns)
│
├─ 4. Check network poller (runtime.netpoll)
│      └─ Any goroutines ready from I/O? → grab them → run one
│
├─ 5. Work stealing — steal from another P
│      └─ Pick random P → steal half its LRQ → run one
│
└─ 6. Nothing anywhere?
       └─ Park this M (release P, go to sleep)
       └─ Will be woken when new work appears
```

### Visual Flow

```
  schedule() on M0/P0:

  ┌──────────────────────────────────────────────────────────┐
  │  P0.runnext ──► G17 → RUN IT                            │
  │       │                                                   │
  │       ▼ (empty)                                           │
  │  P0.runq ──► [G3, G5, G8] → dequeue G3 → RUN IT        │
  │       │                                                   │
  │       ▼ (empty)                                           │
  │  Global Queue ──► [G20, G21, G22] → grab batch → RUN G20│
  │       │                                                   │
  │       ▼ (empty)                                           │
  │  netpoll ──► [G30 (I/O ready)] → RUN G30                │
  │       │                                                   │
  │       ▼ (nothing)                                         │
  │  Steal from P2.runq ──► steal [G40, G41] → RUN G40      │
  │       │                                                   │
  │       ▼ (all Ps empty)                                    │
  │  Park M0 💤 (release P0, sleep until woken)              │
  └──────────────────────────────────────────────────────────┘
```

---

## 7. Work Stealing

When a P's local queue is empty and the global queue is empty, the M enters
**spinning** mode and tries to steal work from other Ps.

```
runtime.stealWork():
  1. Pick a random starting P (randomized to avoid thundering herd)
  2. For each P (round-robin from random start):
     a. Try to steal runnext → if success, done
     b. Try to steal half of P's runq → if success, done
  3. Check all Ps' timer heaps for expired timers
  4. If nothing found → give up, park M

  Stealing takes HALF the victim's queue:
  
  Before:
    P1.runq: [G1, G2, G3, G4, G5, G6]    P2.runq: []
  
  After steal:
    P1.runq: [G1, G2, G3]                 P2.runq: [G4, G5, G6]
```

### Why Half?

Stealing half (not one, not all) balances load without excessive stealing overhead.
If you steal one, you'll be back stealing again immediately. If you steal all,
the victim starves.

### Spinning Ms

The runtime limits the number of spinning Ms to avoid wasting CPU. At most
`GOMAXPROCS` Ms can be spinning. A spinning M actively looks for work rather
than parking — this reduces latency for newly created goroutines at the cost of
CPU usage.

---

## 8. Preemption — Cooperative and Asynchronous

### Cooperative Preemption (Pre-Go 1.14)

The compiler inserts a stack growth check at every function preamble:

```asm
; Function preamble (simplified)
MOVQ  (TLS), CX           ; load current G
CMPQ  SP, 16(CX)          ; compare SP with g.stackguard0
JLS   morestack           ; if SP < stackguard0 → grow stack / preempt
```

To preempt a goroutine, the runtime sets `g.stackguard0 = stackPreempt` (a special
value). At the next function call, the check triggers and the goroutine yields.

**Problem:** A tight loop with no function calls never hits the preamble check:

```go
// This goroutine can NEVER be preempted (pre-Go 1.14):
go func() {
    for {
        i++  // no function calls → no preamble check → never yields
    }
}()
// Other goroutines starve on this P!
```

### Asynchronous Preemption (Go 1.14+)

The runtime sends a **SIGURG** signal to the OS thread running the goroutine:

```
  1. Scheduler decides G has run too long (~10ms sysmon check)
  2. runtime.preemptM(mp) → sends SIGURG to thread M
  3. Signal handler (sighandler) runs on M:
     a. Checks if G is at a "safe point" (no unsafe pointer ops)
     b. If safe → saves registers → sets preempt flag
     c. G is suspended → scheduler picks next G
  4. If NOT at safe point → signal is noted, retry later

  This works even in tight loops with no function calls!
```

**Why SIGURG?** It's an obscure signal that nothing else uses (it's for TCP
out-of-band data, which Go's net package doesn't use). Safe to repurpose.

**Windows:** Uses `SuspendThread`/`GetThreadContext` instead of signals.

---

## 9. Syscall Handling — The Hand-Off

When a goroutine makes a blocking syscall (file I/O, CGo call), the M is blocked
in kernel space. The P can't just wait — it hands off to another M.

```
  Before syscall:
  ┌────┐     ┌────┐
  │ G5 │ ←── │ P0 │ ←── M0 (running G5)
  └────┘     └────┘

  G5 enters syscall (e.g., file read):
  ┌────┐     ┌────┐
  │ G5 │ ←── │    │      M0 (blocked in kernel with G5)
  └────┘     │ P0 │──►   M3 (woken or created, takes over P0)
             └────┘

  G5's syscall returns:
  M0 wakes up with G5
  ├─ Try to reacquire P0 → if free, take it → continue
  ├─ Try to get ANY idle P → if found, take it → continue
  └─ No P available → put G5 on Global Queue → M0 parks
```

### The `sysmon` Thread

`runtime.sysmon` is a special M that runs **without a P** — it's a background
monitor thread. It runs every 20μs-10ms and:

1. **Retakes Ps from syscalls**: if an M has been in a syscall for >20μs, sysmon
   steals its P and hands it to an idle M
2. **Preempts long-running Gs**: if a G has been running for >10ms, triggers
   preemption
3. **Polls the network**: calls `netpoll(0)` (non-blocking) to find ready I/O
4. **Triggers GC**: if needed based on heap growth

```
  sysmon (runs in a loop, no P needed):
  ┌──────────────────────────────────┐
  │ for {                            │
  │   sleep(20μs → 10ms, adaptive)   │
  │   retake Ps from slow syscalls   │
  │   preempt long-running Gs       │
  │   poll network (non-blocking)    │
  │   trigger GC if needed           │
  │ }                                │
  └──────────────────────────────────┘
```

---

## 10. Network Poller — Scalable I/O

Go wraps all network I/O with the runtime's integrated network poller, using
OS-specific mechanisms:

```
  Linux:   epoll    (epoll_create, epoll_ctl, epoll_wait)
  macOS:   kqueue   (kqueue, kevent)
  Windows: IOCP     (CreateIoCompletionPort, GetQueuedCompletionStatus)
```

### How It Works

```
  1. G calls net.Read(conn) → data not ready
  2. runtime.pollWait() → register fd with epoll
  3. G state → _Gwaiting (parked, NO thread consumed)
  4. M detaches G, picks another G from P's queue

  ... later, data arrives on the socket ...

  5. sysmon calls runtime.netpoll() → epoll_wait returns fd
  6. Find the G parked on that fd
  7. G state → _Grunnable → enqueued on a P
  8. G resumes net.Read() with data ready

  Key insight: parked goroutine consumes ZERO thread resources
  → 100K concurrent connections = 100K goroutines, NOT 100K threads
```

### Visual Flow

```
  G1: net.Read()          G2: net.Read()       G3: computing
       │                       │                    │
       ▼                       ▼                    │
  ┌─────────┐            ┌─────────┐               │
  │ _Gwaiting│           │_Gwaiting│               │ running
  │ fd=42    │           │ fd=87   │               │ on M/P
  └────┬─────┘           └────┬────┘               │
       │                      │                     │
       └──────┬───────────────┘                     │
              ▼                                     │
        ┌───────────┐                               │
        │  epoll    │  (kernel watches fd 42, 87)   │
        │  instance │                               │
        └─────┬─────┘                               │
              │ fd 42 ready!                        │
              ▼                                     │
        G1 → _Grunnable → enqueued → runs          │
```

This is why Go excels at network-heavy services — the goroutine-per-connection
model scales to hundreds of thousands of connections with minimal thread usage.

---

## 11. Stack Management

### Contiguous Stacks (Go 1.4+)

Each goroutine starts with a **2KB stack** (since Go 1.4). Stacks are contiguous
memory regions that grow and shrink dynamically.

```
  Initial goroutine stack:
  ┌──────────────────┐ high address (stack.hi)
  │                  │
  │  (unused space)  │  ← only 2KB initially
  │                  │
  ├──────────────────┤ ← SP (stack pointer)
  │  current frame   │
  └──────────────────┘ low address (stack.lo)
```

### Stack Growth

At every function preamble, the compiler inserts a check:

```
  1. Compare SP with g.stackguard0
  2. If SP < stackguard0 → stack is too small → call runtime.morestack()
  3. morestack():
     a. Allocate new stack = 2x current size
     b. Copy entire old stack to new stack
     c. Update ALL pointers into the stack (scan and adjust)
     d. Free old stack
     e. Resume function on new stack
```

```
  Before growth:                    After growth (2x):
  ┌──────────┐ 2KB                 ┌──────────────────────┐ 4KB
  │ frame 3  │                     │                      │
  │ frame 2  │      ──copy──►      │ frame 3              │
  │ frame 1  │                     │ frame 2              │
  │ frame 0  │                     │ frame 1              │
  └──────────┘                     │ frame 0              │
                                   └──────────────────────┘
```

### Stack Shrinking

During GC, if a goroutine's stack is less than 25% utilized, the runtime
halves it (copies to a smaller allocation). This reclaims memory from
goroutines that had deep call stacks temporarily.

### Implications

- **Pointers to stack variables can move** — this is why you can't pass Go
  pointers to C code (CGo pins them to prevent movement)
- **Deep recursion has a hidden cost** — each growth copies the entire stack
- **Function calls are not free** — the preamble check adds ~1-2ns overhead
  (but enables preemption and growth)

---

## 12. Goroutine Lifecycle — Birth to Death

```go
go myFunc(arg1, arg2)
```

### Birth: `runtime.newproc()`

```
  1. go myFunc(args) → compiler emits call to runtime.newproc()
  2. newproc():
     a. Get a free G from P's local free list (gfget)
        └─ If none → allocate new runtime.g + 2KB stack
     b. Set up G's stack frame:
        - Copy function args onto G's stack
        - Set g.sched.pc = myFunc entry point
        - Set return address = runtime.goexit (cleanup stub)
     c. Set g.status = _Grunnable
     d. Put G on current P:
        - runqput(P, G) → tries runnext first, then LRQ
        - If LRQ full → move half to Global Queue
     e. If there's an idle P, wake an M to run it (wakep)
```

### Running

```
  schedule() picks G → execute(G):
  1. Set g.status = _Grunning
  2. Set M.curg = G, G.m = M
  3. gogo(&g.sched) → restore registers from gobuf → jump to g.sched.pc
  4. G's function executes on G's stack
```

### Death: `runtime.goexit()`

```
  1. myFunc() returns → hits runtime.goexit() (planted as return address)
  2. goexit():
     a. Set g.status = _Gdead
     b. Clear G's fields (m, stack references, etc.)
     c. Put G on P's free list (gfput) for reuse
     d. Call schedule() to find next G to run
```

**Key insight:** Goroutines are **recycled**, not freed. The `runtime.g` struct
and its stack are reused for future `go` calls. This amortizes allocation cost.

---

## 13. GOMAXPROCS and Tuning

`GOMAXPROCS` sets the number of Ps — the maximum number of goroutines that can
execute **simultaneously** on CPU.

```
  GOMAXPROCS = 4:
  ┌────┐  ┌────┐  ┌────┐  ┌────┐
  │ P0 │  │ P1 │  │ P2 │  │ P3 │   ← 4 Ps, 4 goroutines run in parallel
  └────┘  └────┘  └────┘  └────┘
    ↕       ↕       ↕       ↕
   M0      M1      M2      M3      ← 4 OS threads (at minimum)
```

### Default Value

```go
runtime.GOMAXPROCS(0)  // returns current value without changing
// Default: runtime.NumCPU() — number of logical CPU cores
```

### In Containers

**Problem:** Inside a Docker container with `--cpus=2`, `runtime.NumCPU()` still
returns the **host's** CPU count (e.g., 64). GOMAXPROCS is set too high →
excessive context switching.

**Solution:** Use `go.uber.org/automaxprocs`:
```go
import _ "go.uber.org/automaxprocs"
// Automatically reads cgroup CPU limits and sets GOMAXPROCS correctly
```

Or set manually:
```bash
GOMAXPROCS=2 ./myservice  # match container CPU limit
```

### When to Tune

```
┌────────────────────────────┬────────────────────────────────────┐
│ Scenario                   │ GOMAXPROCS recommendation          │
├────────────────────────────┼────────────────────────────────────┤
│ CPU-bound work             │ = number of cores (default)        │
│ I/O-bound work             │ = number of cores (still default)  │
│ Container with CPU limit   │ = container CPU limit              │
│ Reduce GC STW latency      │ Higher can help (more GC workers)  │
│ Debugging / deterministic  │ = 1 (forces sequential execution)  │
└────────────────────────────┴────────────────────────────────────┘
```

---

## 14. Performance Characteristics

```
┌─────────────────────────┬──────────────┬────────────────────────────────┐
│ Operation               │ Approx Cost  │ What happens                   │
├─────────────────────────┼──────────────┼────────────────────────────────┤
│ go func() (create)      │ ~0.3 μs      │ Get G from free list, setup    │
│                         │              │ stack, enqueue on P            │
├─────────────────────────┼──────────────┼────────────────────────────────┤
│ Goroutine context       │ ~100-200 ns  │ Save/restore ~15 registers,    │
│ switch                  │              │ swap G pointer on M            │
├─────────────────────────┼──────────────┼────────────────────────────────┤
│ Channel send/recv       │ ~50-300 ns   │ Lock hchan + transfer +        │
│ (see channel doc)       │              │ possible context switch        │
├─────────────────────────┼──────────────┼────────────────────────────────┤
│ Syscall hand-off        │ ~1-5 μs      │ Detach P from M, find/create   │
│                         │              │ new M, attach P                │
├─────────────────────────┼──────────────┼────────────────────────────────┤
│ Stack growth (2x copy)  │ ~1-10 μs     │ Allocate new stack, copy all   │
│                         │              │ frames, update pointers        │
├─────────────────────────┼──────────────┼────────────────────────────────┤
│ Work steal              │ ~0.5-2 μs    │ Lock victim P, copy half queue │
├─────────────────────────┼──────────────┼────────────────────────────────┤
│ Memory per goroutine    │ ~2-8 KB      │ Stack (2KB initial) + g struct │
│ (initial)               │              │ (~400 bytes)                   │
├─────────────────────────┼──────────────┼────────────────────────────────┤
│ 1M goroutines memory    │ ~2-8 GB      │ Mostly stack memory            │
└─────────────────────────┴──────────────┴────────────────────────────────┘
```

### Practical Limits

- **Goroutine count:** Millions possible, but watch total memory. Use
  `runtime.NumGoroutine()` to monitor.
- **Thread count:** Defaults to max 10,000 (`runtime.SetMaxThreads()`). Hit when
  too many goroutines are blocked in syscalls simultaneously.
- **Stack size:** Each goroutine can grow to 1GB (default `runtime.SetMaxStack()`).
  In practice, deep recursion should be avoided.

---

## 15. Debugging Tools

### `GODEBUG=schedtrace=N`

Prints scheduler state every N milliseconds:

```bash
GODEBUG=schedtrace=1000 ./myservice
```

```
SCHED 1004ms: gomaxprocs=4 idleprocs=2 threads=6 spinningthreads=1
  idlethreads=1 runqueue=0 [3 0 1 0]
         │          │         │          │          │           └─ per-P LRQ sizes
         │          │         │          │          └─ global run queue size
         │          │         │          └─ spinning Ms looking for work
         │          │         └─ total OS threads created
         │          └─ Ps with no work (idle)
         └─ total Ps (GOMAXPROCS)
```

### `GODEBUG=scheddetail=1`

Add per-goroutine detail (very verbose):

```bash
GODEBUG=schedtrace=1000,scheddetail=1 ./myservice
```

### `go tool trace`

Visual timeline of all goroutine activity:

```bash
# Collect trace data
go test -trace=trace.out -bench=.

# Open in browser
go tool trace trace.out
```

Shows: goroutine creation/blocking/unblocking, GC phases, syscalls, network
I/O, per-P timeline, and scheduling latency.

### `runtime/pprof` — Goroutine Profile

```go
import "runtime/pprof"

// Dump all goroutine stacks (goroutine leak detection)
pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
```

Or via HTTP:
```go
import _ "net/http/pprof"
// GET /debug/pprof/goroutine?debug=1
```

### Goroutine Leak Detection in Tests

```go
func TestNoLeak(t *testing.T) {
    before := runtime.NumGoroutine()
    // ... run your code ...
    time.Sleep(100 * time.Millisecond) // let goroutines settle
    after := runtime.NumGoroutine()
    if after > before+1 {
        t.Errorf("goroutine leak: before=%d after=%d", before, after)
    }
}
```

Or use `go.uber.org/goleak` for robust leak detection:
```go
func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}
```

---

## 16. Quick Reference Card

```
GMP MODEL
─────────
  G (goroutine)  — unit of work, ~2KB stack, runtime.g struct
  M (machine)    — OS thread, runs Gs, runtime.m struct
  P (processor)  — scheduling context, LRQ of 256 Gs, count = GOMAXPROCS
  G needs P to run. P needs M to run. M without P is idle.

SCHEDULING ORDER
────────────────
  P.runnext → P.runq (local) → Global Queue (every 61st) → netpoll → steal

GOROUTINE STATES
────────────────
  _Grunnable → _Grunning → _Gwaiting (parked) → _Grunnable (woken)
  _Grunning → _Gsyscall → _Grunnable (syscall returned)
  _Grunning → _Gdead (function returned, recycled)

PREEMPTION
──────────
  Cooperative: stack check at function preamble (all versions)
  Asynchronous: SIGURG signal at safe points (Go 1.14+)
  sysmon triggers preemption after ~10ms of running

SYSCALL HAND-OFF
────────────────
  G enters syscall → P detached from M → P given to idle M
  Syscall returns → M tries to reacquire P → or G goes to global queue

NETWORK POLLER
──────────────
  net.Read() blocks → G parked, fd registered with epoll/kqueue/IOCP
  Data arrives → G made runnable → no thread consumed while waiting

STACK MANAGEMENT
────────────────
  Initial: 2KB. Growth: 2x copy. Shrink: halved during GC if < 25% used.
  Max: 1GB. Pointers into stack may move → can't pass to C.

TOOLS
─────
  GODEBUG=schedtrace=1000        scheduler state every second
  GODEBUG=scheddetail=1          per-goroutine detail
  go tool trace                  visual timeline
  runtime.NumGoroutine()         goroutine count
  go.uber.org/goleak             leak detection in tests
  go.uber.org/automaxprocs       container-aware GOMAXPROCS
```

---

## One-Line Summary

> Go's scheduler multiplexes goroutines (G) onto OS threads (M) through
> logical processors (P), using lock-free local run queues, work stealing
> for load balancing, async preemption via SIGURG, and an integrated network
> poller — enabling millions of goroutines on a handful of threads with
> ~200ns context switches.
