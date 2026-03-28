# Deep Dive: Debugging & Profiling — pprof, go tool trace, dlv & GODEBUG

> How to find and fix performance problems, deadlocks, goroutine leaks,
> and memory issues in Go — from development to production.

---

## Table of Contents

1. [Go's Built-In Deadlock Detector — And Why It's Not Enough](#1-gos-built-in-deadlock-detector--and-why-its-not-enough)
2. [pprof — The Swiss Army Knife](#2-pprof--the-swiss-army-knife)
3. [CPU Profiling — Find What's Slow](#3-cpu-profiling--find-whats-slow)
4. [Memory Profiling — Find What Allocates](#4-memory-profiling--find-what-allocates)
5. [Goroutine Profiling — Find What's Stuck](#5-goroutine-profiling--find-whats-stuck)
6. [Block & Mutex Profiling — Find Contention](#6-block--mutex-profiling--find-contention)
7. [go tool trace — The Visual Timeline](#7-go-tool-trace--the-visual-timeline)
8. [GODEBUG — Runtime Diagnostic Flags](#8-godebug--runtime-diagnostic-flags)
9. [dlv (Delve) — The Interactive Debugger](#9-dlv-delve--the-interactive-debugger)
10. [Goroutine Leak Detection](#10-goroutine-leak-detection)
11. [Production Profiling — Safe Practices](#11-production-profiling--safe-practices)
12. [Debugging Workflow — Decision Tree](#12-debugging-workflow--decision-tree)
13. [Quick Reference Card](#13-quick-reference-card)

---

## 1. Go's Built-In Deadlock Detector — And Why It's Not Enough

Go's runtime includes a narrow deadlock detector. It fires when **all** goroutines
are asleep simultaneously:

```go
func main() {
    ch := make(chan int)
    <-ch   // only goroutine blocks → runtime detects deadlock
}
// fatal error: all goroutines are asleep - deadlock!
```

### Why It Almost Never Fires in Production

```go
func main() {
    ch := make(chan int)
    go func() { time.Sleep(time.Hour) }()  // one goroutine "alive"
    <-ch   // blocks forever — NO detection
}
// Runtime: "G2 is sleeping on timer, not stuck — not a deadlock"
```

In production, HTTP servers, background workers, and timers are always running,
so the all-goroutines-asleep condition never triggers.

```
┌─────────────────────────────────────────────────────────────────────┐
│  Go's deadlock detector:                                            │
│  ✅ Catches: ALL goroutines blocked (toy programs, tests)           │
│  ❌ Misses:  ONE goroutine blocked while others run (production)    │
│  ❌ Misses:  Goroutine leaks (blocked goroutines that pile up)      │
│  ❌ Cannot:  Predict future deadlocks at compile time               │
│                                                                     │
│  This is equivalent to the Halting Problem — unsolvable in general. │
│  Go's answer: give you TOOLS to find problems, not prevention.      │
└─────────────────────────────────────────────────────────────────────┘
```

### What Go Provides Instead

```
Tool                          What It Finds
────────────────────────────  ──────────────────────────────────────────
context.WithTimeout()         YOUR code sets a deadline (prevention)
go test -race                 Data races (not deadlocks)
runtime.NumGoroutine()        Goroutine count growth (leak detection)
pprof goroutine profile       WHERE goroutines are stuck (diagnosis)
go tool trace                 VISUAL timeline of blocking (diagnosis)
GODEBUG=schedtrace=1000       Scheduler state snapshots (diagnosis)
goleak (uber library)         Test-time goroutine leak detection
```

---

## 2. pprof — The Swiss Army Knife

`pprof` is Go's built-in profiling system. It's in the standard library — no
external tools needed.

### Two Ways to Use pprof

```
1. FROM TESTS (development):
   go test -cpuprofile=cpu.out -bench=.
   go tool pprof cpu.out

2. FROM A RUNNING SERVER (production):
   import _ "net/http/pprof"    // register handlers
   go http.ListenAndServe(":6060", nil)
   
   go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

### Available Profiles

```
┌──────────────────────┬─────────────────────────────────────────────────┐
│ Profile              │ What It Measures                                │
├──────────────────────┼─────────────────────────────────────────────────┤
│ /debug/pprof/profile │ CPU: where is time spent? (sampling-based)     │
│ /debug/pprof/heap    │ Memory: what allocated? what's still alive?     │
│ /debug/pprof/goroutine │ Goroutines: where are they? what are they doing? │
│ /debug/pprof/block   │ Blocking: where do goroutines wait?            │
│ /debug/pprof/mutex   │ Mutex: where is lock contention?               │
│ /debug/pprof/threadcreate │ OS threads: how many created?             │
│ /debug/pprof/allocs  │ Allocations: all allocations (past + present)  │
│ /debug/pprof/trace   │ Execution trace: full goroutine timeline       │
└──────────────────────┴─────────────────────────────────────────────────┘
```

### pprof Navigation Commands

```bash
$ go tool pprof cpu.out

(pprof) top 10          # top 10 functions by CPU time
(pprof) top -cum        # top by cumulative time (includes callees)
(pprof) list funcName   # source-annotated view of a specific function
(pprof) web             # open interactive graph in browser
(pprof) peek funcName   # show callers and callees of a function
(pprof) disasm funcName # assembly-level view
```

### Web UI (Recommended)

```bash
go tool pprof -http=:8080 cpu.out
# Opens browser with interactive flame graph, source view, graph view
```

---

## 3. CPU Profiling — Find What's Slow

### From Tests

```bash
# Profile a specific benchmark
go test -cpuprofile=cpu.out -bench=BenchmarkProcess ./...
go tool pprof -http=:8080 cpu.out

# Profile all tests (useful for slow test suites)
go test -cpuprofile=cpu.out -run=TestSlow ./...
go tool pprof cpu.out
```

### From Production

```go
import _ "net/http/pprof"

// In main() or init():
go func() {
    // IMPORTANT: bind to localhost or internal network only!
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

```bash
# Capture 30 seconds of CPU profile:
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Or save to file:
curl -o cpu.out http://localhost:6060/debug/pprof/profile?seconds=30
go tool pprof -http=:8080 cpu.out
```

### How CPU Profiling Works Internally

```
The runtime sets a timer (SIGPROF on Unix) that fires 100 times/second.
Each tick:
  1. Interrupt the current goroutine
  2. Capture its full stack trace
  3. Store in profile buffer

After N seconds, you have N×100 stack samples.
Functions that appear in many samples = spending more CPU time.

This is SAMPLING — not tracing. Overhead is ~1-3%. Safe for production.
```

### Reading CPU Profile Output

```
(pprof) top 5
Showing nodes accounting for 4.2s, 84% of 5s total
      flat  flat%   sum%        cum   cum%
     1.8s 36.00% 36.00%      1.8s 36.00%  runtime.memmove
     0.9s 18.00% 54.00%      2.1s 42.00%  encoding/json.(*Decoder).Decode
     0.6s 12.00% 66.00%      0.6s 12.00%  runtime.mallocgc
     0.5s 10.00% 76.00%      0.5s 10.00%  syscall.Syscall
     0.4s  8.00% 84.00%      3.5s 70.00%  main.handleRequest

flat  = time spent IN this function (excluding callees)
cum   = time spent IN this function + all functions it calls
```

```
When flat ≈ cum  → this function IS the bottleneck (doing the work itself)
When flat << cum → this function CALLS the bottleneck (look at its callees)
```

---

## 4. Memory Profiling — Find What Allocates

### From Tests

```bash
go test -memprofile=mem.out -bench=BenchmarkProcess ./...
go tool pprof -http=:8080 mem.out
```

### Four Memory Views

```
┌──────────────────┬──────────────────────────────────────────────────────┐
│ View             │ What It Shows                                        │
├──────────────────┼──────────────────────────────────────────────────────┤
│ inuse_space      │ Bytes currently allocated (what's alive right now)   │
│                  │ → Find memory leaks, large allocations               │
├──────────────────┼──────────────────────────────────────────────────────┤
│ inuse_objects    │ Number of objects currently alive                     │
│                  │ → Find object count issues (too many small objects)   │
├──────────────────┼──────────────────────────────────────────────────────┤
│ alloc_space      │ Total bytes allocated over time (including freed)    │
│                  │ → Find GC pressure sources (hot allocation paths)    │
├──────────────────┼──────────────────────────────────────────────────────┤
│ alloc_objects    │ Total number of allocations over time                 │
│                  │ → Find allocation-heavy code (even if objects are     │
│                  │   small, many allocs = GC pressure)                  │
└──────────────────┴──────────────────────────────────────────────────────┘
```

```bash
# Switch between views in pprof:
(pprof) sample_index = inuse_space     # what's alive now (default)
(pprof) sample_index = alloc_space     # total allocated over time
(pprof) top 10
```

### Combining with Escape Analysis

```bash
# Step 1: Find what allocates (pprof)
go test -memprofile=mem.out -bench=. && go tool pprof mem.out

# Step 2: Find WHY it allocates (escape analysis)
go build -gcflags='-m' ./...
# ./handler.go:42:6: moved to heap: result   ← compiler tells you WHY

# Step 3: Verify fix with benchmarks
go test -bench=BenchmarkProcess -benchmem ./...
# BenchmarkProcess-8   500000   3200 ns/op   256 B/op   4 allocs/op
#                                              ↑          ↑
#                                     bytes allocated   number of heap allocs
```

---

## 5. Goroutine Profiling — Find What's Stuck

This is your primary tool for finding **blocked goroutines** and **goroutine leaks**.

### Quick Count Check

```bash
# From a running server:
curl http://localhost:6060/debug/pprof/goroutine?debug=1 | head -1
# goroutine profile: total 847

# If this number grows monotonically → you have a goroutine leak
```

### Full Goroutine Dump

```bash
# Human-readable dump of ALL goroutine stacks:
curl http://localhost:6060/debug/pprof/goroutine?debug=2
```

Output:

```
goroutine 4821 [chan receive, 47 minutes]:
main.processOrder.func1()
    /app/handler.go:142 +0x8c

goroutine 4822 [chan receive, 47 minutes]:
main.processOrder.func1()
    /app/handler.go:142 +0x8c

goroutine 4823 [chan receive, 46 minutes]:
main.processOrder.func1()
    /app/handler.go:142 +0x8c
```

```
Red flags in goroutine dumps:

  [chan receive, 47 minutes]     ← blocked on channel for 47 min
  [chan send, 2 hours]           ← blocked trying to send for 2 hours
  [select, 30 minutes]          ← stuck in select for 30 min
  [semacquire, 15 minutes]      ← waiting for mutex for 15 min
  [IO wait, 1 hour]             ← stuck on I/O for 1 hour

  Healthy goroutines show:
  [running]                     ← actively executing
  [sleep]                       ← time.Sleep (intentional)
  [IO wait, 50ms]               ← normal network I/O
  [chan receive]                 ← no duration = just started waiting
```

### Grouped View (pprof)

```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
(pprof) top 10
# Shows which functions have the most goroutines stuck in them

(pprof) list processOrder
# Source-level view: which LINE has goroutines blocked
```

---

## 6. Block & Mutex Profiling — Find Contention

### Block Profiling — Where Goroutines Wait

```go
// Enable block profiling (not enabled by default — has overhead):
runtime.SetBlockProfileRate(1)  // 1 = record every blocking event
                                // N = sample 1 in N events (for production)
```

```bash
go tool pprof http://localhost:6060/debug/pprof/block
(pprof) top 5
# Shows: which operations cause goroutines to WAIT the longest
# Typical findings: channel receives, mutex locks, select statements
```

### Mutex Profiling — Where Locks Compete

```go
// Enable mutex profiling:
runtime.SetMutexProfileFraction(1)  // 1 = record every contention event
```

```bash
go tool pprof http://localhost:6060/debug/pprof/mutex
(pprof) top 5
# Shows: which mutexes have the most contention (multiple goroutines
# trying to acquire the same lock simultaneously)
```

#### Understanding Block Profiling in Depth

`runtime.SetBlockProfileRate(n)` records how long goroutines spend **blocked** — waiting
on channel ops, mutex acquisition, select statements, and condition variables. Use it when
you suspect goroutines are spending time **waiting** rather than working. The profile shows
the **call stack** where the block occurred and the **total accumulated blocked time**.

- Rate `n` = nanoseconds threshold. `n=1` records **all** blocking events. Higher values
  (e.g., `n=1000`) only record blocks longer than that threshold, reducing overhead.
- Scope: channels, mutexes, select, `sync.Cond` — anything that calls `gopark()`.

#### Understanding Mutex Profiling in Depth

`runtime.SetMutexProfileFraction(n)` records contention on `sync.Mutex` and `sync.RWMutex`
**only**. Use it when you suspect lock contention is degrading throughput. The profile shows
**which locks** are contended and **how long** goroutines waited to acquire them.

- Fraction `n`: `1/n` of mutex contention events are recorded. `n=1` records all events,
  `n=5` records ~20% of events (lower overhead, statistically representative).
- Scope: **only** `sync.Mutex` and `sync.RWMutex` — does **not** cover channels or select.

#### Block vs Mutex — When to Use Which

- **Block profiling** = broad net (channels + mutexes + select + cond vars). Start here.
- **Mutex profiling** = focused lens (only mutexes, lower noise). Use when you've narrowed
  the problem to lock contention specifically.
- **Production overhead**: Block profiling with a high rate (e.g., `n=100000`) has minimal
  overhead. Mutex profiling at `fraction=5` is safe for production use.

### When to Use Each

```
"My service is slow"
├── CPU profile shows low CPU usage
│   ├── Block profile → goroutines waiting on channels/IO
│   └── Mutex profile → goroutines fighting over locks
└── CPU profile shows high CPU usage
    └── CPU profile → which function is hot
```

---

## 7. go tool trace — The Visual Timeline

`go tool trace` gives you a **visual, interactive timeline** of everything the
runtime does: goroutine scheduling, GC pauses, syscalls, network I/O, blocking.

### Capture a Trace

```bash
# From tests:
go test -trace=trace.out -bench=. ./...
go tool trace trace.out

# From production (short capture):
curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5
go tool trace trace.out

# Programmatically:
import "runtime/trace"

f, _ := os.Create("trace.out")
trace.Start(f)
defer trace.Stop()
```

### What the Trace Shows

```
┌─────────────────────────────────────────────────────────────────────┐
│  Timeline (horizontal = time, vertical = goroutines/procs)         │
│                                                                     │
│  Proc 0: ██G1████░░░░██G3████░░GC░░██G1████                       │
│  Proc 1: ░░░░██G2████████░░░░░░GC░░░░██G4████                     │
│  Proc 2: ██G5████░░syscall░░░░░░░░░░██G5████                      │
│                                                                     │
│  ██ = running                                                       │
│  ░░ = idle (no goroutine scheduled on this P)                      │
│  GC = garbage collection (STW phases shown in red)                 │
│  syscall = goroutine in syscall (detached from P)                  │
│                                                                     │
│  Click any block to see: goroutine ID, function, stack trace,      │
│  why it was blocked, how long it ran, what unblocked it.           │
└─────────────────────────────────────────────────────────────────────┘
```

### What to Look For

```
Problem                     What Trace Shows
─────────────────────────   ──────────────────────────────────────────
Latency spikes              Long GC pauses (red STW bars)
Underutilization            Many idle procs (not enough goroutines)
Lock contention             Goroutines blocked on sync events
Channel bottleneck          Many goroutines blocked on same channel
Scheduling delays           Long gaps between runnable → running
Syscall storms              Many goroutines in syscall simultaneously
```

### pprof vs trace — When to Use Which

```
pprof:  "WHERE is the problem?" (statistical, low overhead, production-safe)
        → Use for: CPU hotspots, memory leaks, goroutine counts

trace:  "WHAT HAPPENED and WHEN?" (deterministic, higher overhead, short captures)
        → Use for: latency investigation, scheduler behavior, GC impact, concurrency bugs
```

---

## 8. GODEBUG — Runtime Diagnostic Flags

Environment variables that make the Go runtime print diagnostic information.
Zero code changes needed — just set the variable before running.

### GC Tracing

```bash
GODEBUG=gctrace=1 ./myapp
```

Output:

```
gc 1 @0.012s 2%: 0.021+0.54+0.003 ms clock, 0.17+0.32/0.41/0+0.025 ms cpu, 4->4->1 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 8 P
│   │      │    │                     │                                    │              │             │
│   │      │    │                     │                                    │              │             └── GOMAXPROCS
│   │      │    │                     │                                    │              └── goal heap size
│   │      │    │                     │                                    └── heap: before→after→live
│   │      │    │                     └── CPU time: assist/background/idle
│   │      │    └── wall clock: STW1+concurrent+STW2
│   │      └── % of CPU used for GC
│   └── time since program start
└── GC cycle number
```

**What to look for:**
- `2%` → GC using 2% of CPU — healthy
- `25%` → GC using 25% of CPU — too much allocation, reduce heap pressure
- `STW` phases > 1ms → investigate large pointer-heavy data structures
- `before→after→live`: if `live` keeps growing → memory leak

### Scheduler Tracing

```bash
GODEBUG=schedtrace=1000 ./myapp      # print every 1000ms
GODEBUG=schedtrace=1000,scheddetail=1 ./myapp  # detailed per-P info
```

Output:

```
SCHED 1004ms: gomaxprocs=8 idleprocs=6 threads=10 spinningthreads=0
  idlethreads=5 runqueue=0 [0 0 0 0 0 0 0 0]
                                              └── per-P local run queue lengths
                           └── global run queue length
              └── how many Ps are idle (6 of 8 = underutilized)
```

**What to look for:**
- `idleprocs=7` of `gomaxprocs=8` → only 1 P is busy — not enough parallelism
- `runqueue=500 [100 50 200 ...]` → goroutines piling up, can't keep up
- `threads=200` → too many OS threads — goroutines stuck in syscalls

### Other Useful GODEBUG Flags

```
GODEBUG=madvdontneed=1     # return memory to OS more aggressively (Linux)
GODEBUG=asyncpreemptoff=1  # disable async preemption (debugging only)
GODEBUG=invalidptr=1       # crash on invalid pointer (finding corruption)
GODEBUG=cgocheck=2         # strict cgo pointer checks
```

---

## 9. dlv (Delve) — The Interactive Debugger

Delve is Go's purpose-built debugger. It understands goroutines, channels,
interfaces, and defer stacks — unlike GDB which doesn't.

### Installation

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

### Basic Usage

```bash
# Debug a program:
dlv debug ./cmd/server

# Debug a test:
dlv test ./pkg/handler -- -test.run TestProcess

# Attach to running process:
dlv attach <PID>

# Core dump analysis:
dlv core ./myapp core.12345
```

### Essential Commands

```
(dlv) break main.handleRequest    # set breakpoint
(dlv) break handler.go:42         # breakpoint at line
(dlv) condition 1 id == 42        # conditional breakpoint
(dlv) continue                    # run until breakpoint
(dlv) next                        # step over (next line)
(dlv) step                        # step into function call
(dlv) stepout                     # step out of current function
(dlv) print myVar                 # inspect variable
(dlv) print *myPointer            # dereference pointer
(dlv) locals                      # show all local variables
(dlv) args                        # show function arguments
(dlv) stack                       # show stack trace
(dlv) goroutines                  # list ALL goroutines
(dlv) goroutine 15                # switch to goroutine 15
(dlv) goroutine 15 stack          # stack of specific goroutine
```

### Debugging Goroutines

```
(dlv) goroutines
  Goroutine 1 - User: main.main (main.go:15) (thread 12345)
  Goroutine 5 - User: main.worker (worker.go:42) [chan receive]
  Goroutine 6 - User: main.worker (worker.go:42) [chan receive]
  Goroutine 7 - User: main.processor (proc.go:18) [running]
  
  [chan receive] = blocked on channel — dlv shows you the state!
  
(dlv) goroutine 5          # switch to goroutine 5
(dlv) stack                # see its full call stack
(dlv) frame 2              # navigate to frame 2 in the stack
(dlv) locals               # see local variables in that frame
(dlv) print ch             # inspect the channel it's blocked on
```

### Debugging Channel State

```
(dlv) print ch
chan int {
    qcount: 3,          # 3 values in buffer
    dataqsiz: 5,        # buffer capacity = 5
    buf: *[5]int{10, 20, 30, 0, 0},  # buffer contents
    closed: 0,          # not closed
    sendx: 3,           # next send index
    recvx: 0,           # next receive index
    recvq: waitq{...},  # waiting receivers
    sendq: waitq{...},  # waiting senders
}
```

### When to Use dlv vs pprof

```
dlv:   "I need to inspect SPECIFIC STATE at a SPECIFIC MOMENT"
       → Breakpoints, variable inspection, step-through
       → Use for: logic bugs, unexpected values, control flow issues

pprof: "I need to find PATTERNS across the WHOLE PROGRAM"
       → Statistical profiling, aggregate views
       → Use for: performance issues, memory leaks, goroutine leaks
```

---

## 10. Goroutine Leak Detection

A goroutine leak is when goroutines are created but never exit — they pile up
over time, consuming memory and eventually crashing the process.

### Common Leak Patterns

```go
// Leak 1: Sending to a channel nobody reads
func leak1() {
    ch := make(chan int)
    go func() {
        result := expensiveWork()
        ch <- result      // blocks forever if nobody reads ch
    }()
    // ch goes out of scope but goroutine is stuck on send
}

// Leak 2: Receiving from a channel nobody writes to or closes
func leak2() {
    ch := make(chan int)
    go func() {
        for v := range ch {  // waits forever — ch is never closed
            process(v)
        }
    }()
}

// Leak 3: Missing context cancellation
func leak3() {
    go func() {
        ticker := time.NewTicker(time.Second)
        for range ticker.C {
            doWork()     // runs forever — no stop condition
        }
    }()
}
```

### Detection in Tests — goleak

```go
import "go.uber.org/goleak"

func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)    // fails if goroutines leak between tests
}

// Or per-test:
func TestProcess(t *testing.T) {
    defer goleak.VerifyNone(t)
    // ... test code ...
}
```

### Detection in Tests — Manual

```go
func TestNoLeak(t *testing.T) {
    before := runtime.NumGoroutine()
    
    // ... exercise the code ...
    
    time.Sleep(100 * time.Millisecond)  // let goroutines settle
    after := runtime.NumGoroutine()
    
    if after > before {
        t.Errorf("goroutine leak: before=%d after=%d", before, after)
        // Dump goroutine stacks for debugging:
        buf := make([]byte, 1<<20)
        n := runtime.Stack(buf, true)  // true = all goroutines
        t.Logf("goroutine dump:\n%s", buf[:n])
    }
}
```

### Detection in Production — Metrics

```go
// Expose as Prometheus metric:
prometheus.NewGaugeFunc(prometheus.GaugeOpts{
    Name: "go_goroutines",
    Help: "Number of goroutines that currently exist.",
}, func() float64 {
    return float64(runtime.NumGoroutine())
})

// Alert rule (Prometheus):
// ALERT GoroutineLeak
//   IF go_goroutines > 10000
//   OR rate(go_goroutines[5m]) > 10
//   FOR 5m
```

### Prevention Patterns

```go
// ALWAYS pair goroutine creation with a cancellation mechanism:

// Pattern 1: Context cancellation
func safeWorker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return          // exits when context cancelled
        case job := <-jobs:
            process(job)
        }
    }
}

// Pattern 2: Done channel
func safeProducer(done <-chan struct{}) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 0; ; i++ {
            select {
            case <-done:
                return      // exits when done closed
            case out <- i:
            }
        }
    }()
    return out
}

// Pattern 3: Buffered channel (prevents send-side leak)
func safeOneShot() <-chan Result {
    ch := make(chan Result, 1)     // buffer of 1!
    go func() {
        ch <- expensiveWork()     // even if nobody reads, goroutine exits
    }()                           // because buffer absorbs the value
    return ch
}
```

---

## 11. Production Profiling — Safe Practices

### Security: Don't Expose pprof to the Internet

```go
// ❌ DANGEROUS — pprof accessible from anywhere:
go http.ListenAndServe(":6060", nil)

// ✅ SAFE — bind to localhost only:
go http.ListenAndServe("localhost:6060", nil)

// ✅ SAFE — separate mux with auth middleware:
debugMux := http.NewServeMux()
debugMux.HandleFunc("/debug/pprof/", pprof.Index)
debugMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
debugMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
// ... add authentication middleware ...
go http.ListenAndServe("localhost:6060", authMiddleware(debugMux))
```

### Performance Impact

```
┌─────────────────────────┬──────────────┬─────────────────────────────┐
│ Profile Type            │ Overhead     │ Production-Safe?            │
├─────────────────────────┼──────────────┼─────────────────────────────┤
│ CPU profile             │ ~1-3%        │ ✅ Yes (sampling-based)     │
│ Heap profile            │ ~0%          │ ✅ Yes (always collected)   │
│ Goroutine profile       │ ~0%          │ ✅ Yes (snapshot)           │
│ Block profile           │ ~5-10%       │ ⚠️  Rate-limit in prod     │
│ Mutex profile           │ ~5-10%       │ ⚠️  Rate-limit in prod     │
│ Execution trace         │ ~10-20%      │ ⚠️  Short captures only    │
└─────────────────────────┴──────────────┴─────────────────────────────┘
```

### Continuous Profiling in Production

```go
// Pattern: periodic profile capture for monitoring platforms
func startContinuousProfiling(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // Capture 10-second CPU profile
            f, _ := os.CreateTemp("", "cpuprof-*.out")
            pprof.StartCPUProfile(f)
            time.Sleep(10 * time.Second)
            pprof.StopCPUProfile()
            f.Close()
            uploadToMonitoring(f.Name())
            os.Remove(f.Name())
        }
    }
}
```

Services like **Google Cloud Profiler**, **Datadog Continuous Profiler**, and
**Pyroscope** automate this — they continuously sample and aggregate profiles.

---

## 12. Debugging Workflow — Decision Tree

```
SYMPTOM                                TOOL                        LOOK FOR
──────────────────────────────────────────────────────────────────────────────

"Service is slow"
├── High CPU?
│   └── CPU pprof                      → top functions by flat time
│       go tool pprof profile
│
├── Low CPU but high latency?
│   ├── Block pprof                    → channel/mutex wait times
│   │   runtime.SetBlockProfileRate(1)
│   ├── Mutex pprof                    → lock contention
│   │   runtime.SetMutexProfileFraction(1)
│   └── go tool trace                  → scheduling delays, GC pauses
│
└── Latency spikes every N minutes?
    └── GODEBUG=gctrace=1             → GC pause correlation
        go tool trace                  → STW phase durations

"Memory keeps growing"
├── Heap pprof (inuse_space)           → what's alive and not freed
│   go tool pprof heap
├── Goroutine pprof                    → goroutine leak (each = 2-8KB stack)
│   curl .../debug/pprof/goroutine
└── GODEBUG=gctrace=1                 → is live heap growing each cycle?

"Goroutines keep increasing"
├── Goroutine pprof (debug=2)          → full stack of each goroutine
│   curl .../debug/pprof/goroutine?debug=2
├── goleak in tests                    → catch leaks before production
└── runtime.NumGoroutine() metric      → monitor trend over time

"Deadlock / hang"
├── Goroutine pprof                    → find which goroutines are stuck
├── dlv attach <PID>                   → inspect channel/mutex state
│   (dlv) goroutines
│   (dlv) goroutine N stack
└── GODEBUG=schedtrace=1000           → is runqueue growing?

"Wrong behavior / logic bug"
├── dlv debug                          → breakpoints, variable inspection
│   (dlv) break handler.go:42
│   (dlv) print myStruct
└── go test -race                      → data race detection
```

---

## 13. Quick Reference Card

```
CPU PROFILING
─────────────
  go test -cpuprofile=cpu.out -bench=.
  go tool pprof -http=:8080 cpu.out
  Production: import _ "net/http/pprof"

MEMORY PROFILING
────────────────
  go test -memprofile=mem.out -bench=.
  go tool pprof mem.out → sample_index=inuse_space (leaks)
                        → sample_index=alloc_space (GC pressure)

GOROUTINE PROFILING
───────────────────
  curl localhost:6060/debug/pprof/goroutine?debug=1 | head -1   # count
  curl localhost:6060/debug/pprof/goroutine?debug=2              # all stacks
  Look for: [chan receive, 47 minutes] ← blocked too long

BLOCK & MUTEX PROFILING
───────────────────────
  runtime.SetBlockProfileRate(1)        # enable block profiling
  runtime.SetMutexProfileFraction(1)    # enable mutex profiling
  go tool pprof localhost:6060/debug/pprof/block
  go tool pprof localhost:6060/debug/pprof/mutex

EXECUTION TRACE
───────────────
  go test -trace=trace.out -bench=.
  go tool trace trace.out               # visual timeline in browser
  curl -o trace.out localhost:6060/debug/pprof/trace?seconds=5

GODEBUG FLAGS
─────────────
  GODEBUG=gctrace=1 ./app              # GC cycle stats
  GODEBUG=schedtrace=1000 ./app        # scheduler state every 1s
  GODEBUG=schedtrace=1000,scheddetail=1 ./app  # per-P detail

DELVE DEBUGGER
──────────────
  dlv debug ./cmd/server               # debug program
  dlv test ./pkg/handler               # debug tests
  dlv attach <PID>                     # attach to running process
  (dlv) goroutines                     # list all goroutines
  (dlv) goroutine N stack              # stack of goroutine N
  (dlv) print ch                       # inspect channel state

GOROUTINE LEAK DETECTION
─────────────────────────
  goleak.VerifyTestMain(m)             # uber/goleak in TestMain
  runtime.NumGoroutine()               # manual count in tests
  runtime.Stack(buf, true)             # dump all goroutine stacks
  Prometheus: go_goroutines metric     # production monitoring

ESCAPE ANALYSIS (what goes to heap)
───────────────────────────────────
  go build -gcflags='-m' ./...         # basic escape decisions
  go build -gcflags='-m -m' ./...      # detailed reasons

RACE DETECTOR
─────────────
  go test -race ./...                  # NON-NEGOTIABLE in CI
  go run -race ./cmd/server            # during development
```

---

## One-Line Summary

> Go has no compile-time deadlock prevention — instead it gives you `pprof` for
> profiling (CPU, memory, goroutines, contention), `go tool trace` for visual
> timelines, `GODEBUG` for runtime diagnostics, `dlv` for interactive debugging,
> and `goleak` for test-time leak detection. Your job is to use them.
