# Deep Dive: Goroutine Scheduler Internals — The GMP Model

> How Go multiplexes millions of goroutines onto OS threads: the G-M-P model,
> work stealing, preemption, syscall handling, and the network poller.

---

## Table of Contents

1. [Goroutines vs OS Threads](#1-goroutines-vs-os-threads)
2. [The GMP Model](#2-the-gmp-model)
3. [The Three Structs — What G, M, and P Actually Hold](#3-the-three-structs--what-g-m-and-p-actually-hold)
4. [Scheduling Algorithm](#4-scheduling-algorithm)
5. [Work Stealing](#5-work-stealing)
6. [Preemption — Cooperative and Asynchronous](#6-preemption--cooperative-and-asynchronous)
7. [Syscall Handling — The Hand-Off](#7-syscall-handling--the-hand-off)
8. [Network Poller — Scalable I/O](#8-network-poller--scalable-io)
9. [Stack Management](#9-stack-management)
10. [Goroutine Lifecycle](#10-goroutine-lifecycle)
11. [GOMAXPROCS and Tuning](#11-gomaxprocs-and-tuning)
12. [Performance Characteristics](#12-performance-characteristics)
13. [Debugging Tools](#13-debugging-tools)
14. [Interview Self-Test](#14-interview-self-test)
15. [Quick Reference Card](#15-quick-reference-card)

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

**Why goroutine switches are cheaper:** They stay in user space (no kernel trap),
save only ~15 registers into `g.sched` (a `gobuf` struct), and swap a pointer.
OS thread switches require ring 3→0 kernel transition, full register/FPU/TLB save, and kernel scheduling.

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

---

## 3. The Three Structs — What G, M, and P Actually Hold

### G (Goroutine) — `runtime.g`

Each goroutine is a `runtime.g` struct (~400 bytes). Key fields:

- **stack** `{lo, hi}` — current stack bounds (starts at 2KB, grows dynamically)
- **sched** (`gobuf`) — saved registers for context switch: stack pointer (SP), program counter (PC), base pointer (BP). This is how Go suspends and resumes a goroutine
- **atomicstatus** — current goroutine state (see diagram below)
- **goid** — goroutine ID (internal, not exposed to user code)
- **m** — pointer to the M currently running this G (nil if not running)
- **preempt** — flag set by runtime to signal "yield at next safe point"

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

### P (Processor) — `runtime.p`

P is the **scheduling context** — it owns a run queue and per-P resources. Key fields:

- **runq** — local run queue (LRQ), a lock-free circular buffer of 256 G pointers
- **runnext** — single G slot for "next to run" (fast path, bypasses queue)
- **mcache** — per-P memory allocator cache (no locking needed!)
- **timers** — per-P timer heap (`time.After`, `time.Sleep`, etc.)
- **m** — the M currently attached to this P

### The Local Run Queue (LRQ)

```
  P.runq — lock-free circular buffer of 256 slots:

  ┌───┬───┬───┬───┬───┬───┬───┬───┬─── ─── ───┬───┐
  │G12│G13│G14│G15│   │   │   │   │  ...       │   │
  └───┴───┴───┴───┴───┴───┴───┴───┴─── ─── ───┴───┘
   ^head              ^tail

  runqput (enqueue, CAS) · runqget (dequeue, CAS) · runqsteal (steal half, CAS)
  When full (256 Gs) → move HALF to Global Run Queue (prevents hoarding)
```

**`P.runnext`:** holds a single G scheduled **next**, bypassing the queue. When
`go func()` is called, the new goroutine goes into `runnext` — giving freshly
created goroutines priority (great for producer-consumer patterns).

**Why P exists:** Before Go 1.1, there was no P — just G and M. Every goroutine
operation required a global mutex on the run queue, and memory caches were per-M
(wasted when Ms parked). P solves both: lock-free local queues + per-P mcache.
The global queue is only touched on overflow or starvation.

### M (Machine) — `runtime.m`

M represents an OS thread. The runtime creates Ms as needed. Key fields:

- **g0** — special goroutine with ~8KB stack for running scheduler code
- **curg** — the user goroutine currently executing on this M
- **p** — the P attached to this M (nil if idle or in syscall)
- **spinning** — true when M is actively looking for work (work stealing)

### The `g0` — Scheduler Stack

When a goroutine needs to be scheduled (channel block, syscall, preemption),
execution switches from the user G's stack to `g0`'s stack to run `schedule()`.

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

Need a new M? → Check idle list (`runtime.midle`) → Found? Wake it up. Not found?
Create new OS thread. Max: 10,000 Ms (`runtime.SetMaxThreads()`). Created on demand,
reused when parked.

## 4. Scheduling Algorithm

The core loop is `runtime.schedule()` — it runs on `g0`'s stack on every M.

```
runtime.schedule():
│
├─ 1. P.runnext? → run it (fast path)
├─ 2. P.runq (local)? → dequeue from head → run it
├─ 3. Every 61st call: check Global Queue → grab batch → run one
│     (61 is prime — avoids synchronization patterns)
├─ 4. netpoll → any goroutines ready from I/O? → run one
├─ 5. Work steal → pick random P → steal half its LRQ → run one
└─ 6. Nothing anywhere → park M (release P, go to sleep)
```

## 5. Work Stealing

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

**Why half?** One is too few (you'll steal again immediately). All is too many
(victim starves). Half balances load in one operation.

**Spinning Ms:** At most `GOMAXPROCS` Ms can spin. Spinning trades CPU for lower
latency — a spinning M picks up new goroutines faster than a parked one.

## 6. Preemption — Cooperative and Asynchronous

### Cooperative Preemption (Pre-Go 1.14)

The compiler inserts a **stack growth check** at every function preamble: it
compares SP with `g.stackguard0`. To preempt, the runtime sets `stackguard0`
to a special `stackPreempt` value — the next function call triggers a yield.

**Problem:** A tight loop like `for { i++ }` has no function calls, so the
preamble check never fires and the goroutine never yields (starving others).

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

## 7. Syscall Handling — The Hand-Off

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

## 8. Network Poller — Scalable I/O

Go wraps all network I/O with the runtime's integrated network poller
(epoll on Linux, kqueue on macOS, IOCP on Windows).

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

This is why Go excels at network-heavy services — the goroutine-per-connection
model scales to hundreds of thousands of connections with minimal thread usage.

## 9. Stack Management

Each goroutine starts with a **2KB contiguous stack** that grows dynamically (2x
copies) and shrinks during GC. The compiler inserts a stack-check preamble at
every function entry, which also serves as the cooperative preemption point.

> **Full deep dive:** See [Chapter 11 — Goroutine Stacks & Growth](./11_goroutine_stacks_growth.md)
> for contiguous vs segmented stack history, pointer adjustment algorithm,
> stack maps, shrinking mechanics, and performance implications.

## 10. Goroutine Lifecycle

```
  ┌─── BIRTH: runtime.newproc() ───────────────────────────────────────┐
  │  go myFunc(args)                                                    │
  │    → Get free G from P's free list (or allocate new G + 2KB stack)  │
  │    → Copy args onto G's stack, set PC = myFunc entry point          │
  │    → Set return address = runtime.goexit (cleanup stub)             │
  │    → G status = _Grunnable                                         │
  │    → Place on P: runnext first, then LRQ, overflow → GRQ            │
  │    → If idle P exists, wake an M to run it                          │
  ├─── RUNNING: schedule() → execute() ────────────────────────────────┤
  │  schedule() picks G → set _Grunning → M.curg = G, G.m = M          │
  │    → gogo(&g.sched) → restore registers → jump to PC               │
  │    → G's function executes on G's stack                             │
  ├─── DEATH: runtime.goexit() ────────────────────────────────────────┤
  │  myFunc() returns → hits goexit() (planted as return address)       │
  │    → G status = _Gdead → clear fields                              │
  │    → Put G on P's free list for reuse → call schedule()             │
  └─────────────────────────────────────────────────────────────────────┘
```

**Key insight:** Goroutines are **recycled**, not freed. The `runtime.g` struct
and its stack are reused for future `go` calls. This amortizes allocation cost.

## 11. GOMAXPROCS and Tuning

`GOMAXPROCS` sets the number of Ps — the maximum number of goroutines that can
execute **simultaneously** on CPU.

```
  GOMAXPROCS = 4:   P0↔M0   P1↔M1   P2↔M2   P3↔M3
  (4 Ps = 4 goroutines run in parallel, each on an OS thread)
```

**Default:** `runtime.NumCPU()` (logical core count). Query with `runtime.GOMAXPROCS(0)`.

**In containers:** `runtime.NumCPU()` returns the **host's** CPU count, not the
container limit → excessive context switching. Fix with `go.uber.org/automaxprocs`
(auto-reads cgroup limits) or set `GOMAXPROCS` manually to match container CPU limit.

**When to tune:** Usually don't — default (= core count) works for both CPU-bound
and I/O-bound workloads. In containers, match the CPU limit. For debugging, set
to 1 to force sequential execution.

## 12. Performance Characteristics

```
┌───────────────────────┬─────────────┬─────────────────────────────────────┐
│ Operation             │ Approx Cost │ What happens                        │
├───────────────────────┼─────────────┼─────────────────────────────────────┤
│ go func() (create)    │ ~0.3 μs     │ Get G from free list, setup, enqueue│
│ Goroutine ctx switch  │ ~100-200 ns │ Save/restore ~15 regs, swap G on M  │
│ Channel send/recv     │ ~50-300 ns  │ Lock hchan + transfer + maybe switch│
│ Syscall hand-off      │ ~1-5 μs     │ Detach P, find/create M, attach P   │
│ Stack growth (2x)     │ ~1-10 μs    │ Alloc new stack, copy, fix pointers │
│ Work steal            │ ~0.5-2 μs   │ Lock victim P, copy half its queue  │
│ Memory per goroutine  │ ~2-8 KB     │ Stack (2KB) + g struct (~400 bytes) │
│ 1M goroutines memory  │ ~2-8 GB     │ Mostly stack memory                 │
└───────────────────────┴─────────────┴─────────────────────────────────────┘
```

**Practical limits:** Millions of goroutines are possible (watch memory via
`runtime.NumGoroutine()`). Max OS threads: 10,000 (`runtime.SetMaxThreads()`) —
hit when too many goroutines block in syscalls simultaneously.

## 13. Debugging Tools

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

### Other Essential Tools

- **`go tool trace`**: Visual timeline of goroutine scheduling, GC, syscalls.
  Collect with `go test -trace=trace.out -bench=.`, view with `go tool trace trace.out`.
- **`runtime.NumGoroutine()`**: Check goroutine count at runtime — useful for leak detection.
- **`go.uber.org/goleak`**: Robust goroutine leak detection in tests via `goleak.VerifyTestMain(m)`.

> **See [Chapter 15](./15_debugging_profiling.md)** for full coverage of `pprof`, goroutine profiling, `scheddetail`,
> and production debugging workflows.

## 14. Interview Self-Test

Can you answer these without scrolling up? If yes, you own this topic.

**Q1: "What are G, M, and P?"**

> G = Goroutine — the unit of work (~2KB stack, growable).
> M = Machine — an OS thread that executes code.
> P = Processor — logical scheduling context that holds the local run queue.
> GOMAXPROCS controls how many Ps exist. G needs P to run. P needs M to run.

**Q2: "How does Go handle a blocking syscall?"**

> M detaches from P so the P isn't wasted. P finds another M (wakes an idle
> one or creates a new thread). The goroutine stays on the blocked M in
> kernel space. When the syscall returns, G re-enters a P's run queue.

**Q3: "What is work stealing?"**

> When a P's local run queue is empty, it steals half of another P's queue.
> This keeps all cores busy without a central scheduler lock. Half is the
> sweet spot — stealing one means you're back stealing immediately, stealing
> all starves the victim.

**Q4: "How does preemption work?"**

> Pre-1.14: cooperative only — the compiler inserts a stack check at every
> function preamble. Tight loops with no function calls could never be
> preempted. Go 1.14+: async preemption via SIGURG signal. The runtime's
> sysmon thread detects goroutines running >10ms and sends SIGURG to the OS
> thread, which suspends the goroutine at a safe point.

**Q5: "Goroutines vs OS threads?"**

> Goroutines are user-space, start at 2KB (growable to 1GB), context switch
> in ~200ns (just save/restore ~15 registers, no kernel involved). OS threads
> cost ~1MB stack, ~1-5μs kernel context switch. Go uses M:N scheduling —
> millions of goroutines on a handful of OS threads.

**Q6: "What happens when you `go func()`?"**

> The compiler emits `runtime.newproc()`. It grabs a G from the P's free list
> (or allocates a new one + 2KB stack), copies the function args onto G's
> stack, sets the program counter, and puts G on the current P's run queue
> (runnext slot for priority, or LRQ). Goroutines are recycled after they
> finish — the struct and stack are reused.

---

## 15. Quick Reference Card

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
  go tool trace                  visual timeline
  runtime.NumGoroutine()         goroutine count
  go.uber.org/goleak             leak detection in tests
  go.uber.org/automaxprocs       container-aware GOMAXPROCS
  See Chapter 15                 full debugging & profiling coverage
```

---

## One-Line Summary

> Go's scheduler multiplexes goroutines (G) onto OS threads (M) through
> logical processors (P), using lock-free local run queues, work stealing
> for load balancing, async preemption via SIGURG, and an integrated network
> poller — enabling millions of goroutines on a handful of threads with
> ~200ns context switches.

---

## Further Reading

- [runtime/proc.go](https://cs.opensource.google/go/go/+/master:src/runtime/proc.go) — Core scheduler source: `schedule()`, `findrunnable()`, work stealing, and the `sysmon` thread
- [Scalable Go Scheduler Design Doc](https://docs.google.com/document/d/1TTj4T2JO42uD5ID9e89oa0sLKhJYD0Y_kqxDv3I3XMw) — Dmitry Vyukov's original design document for the G-M-P scheduler model
- [runtime/runtime2.go](https://cs.opensource.google/go/go/+/master:src/runtime/runtime2.go) — Runtime definitions for the `g`, `m`, and `p` structs with all key fields
- [Issue #24543 — Non-cooperative preemption](https://github.com/golang/go/issues/24543) — Proposal for asynchronous preemption via SIGURG signals (implemented in Go 1.14)
- [runtime/netpoll.go](https://cs.opensource.google/go/go/+/master:src/runtime/netpoll.go) — Integrated network poller abstraction over epoll/kqueue/IOCP
- [Go's work-stealing scheduler](https://go.dev/src/runtime/proc.go) — The `stealWork` and `runqsteal` functions implementing the work-stealing algorithm
