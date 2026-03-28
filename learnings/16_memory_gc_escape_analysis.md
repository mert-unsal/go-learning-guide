# Deep Dive: Memory, Garbage Collection & Escape Analysis

> How Go decides where values live, how the allocator obtains memory,
> how the GC reclaims it concurrently, and how you control all of it in production.
>
> **Related:** [Chapter 06 §4–5](./06_closures_funcval_and_capture.md) covers escape analysis
> specifically for closures and goroutines (capture-by-reference, funcval escaping).

---

## Table of Contents

1. [Stack vs Heap — The Fundamental Decision](#1-stack-vs-heap--the-fundamental-decision)
2. [Escape Analysis — How the Compiler Decides](#2-escape-analysis--how-the-compiler-decides)
3. [The Memory Allocator — TCMalloc-Inspired](#3-the-memory-allocator--tcmalloc-inspired)
4. [The Garbage Collector — Tri-Color Mark and Sweep](#4-the-garbage-collector--tri-color-mark-and-sweep)
5. [Write Barrier — Protecting Concurrent Marking](#5-write-barrier--protecting-concurrent-marking)
6. [GC Pacing — GOGC and GOMEMLIMIT](#6-gc-pacing--gogc-and-gomemlimit)
7. [GC Phases and STW Pauses](#7-gc-phases-and-stw-pauses)
8. [sync.Pool — GC-Aware Object Reuse](#8-syncpool--gc-aware-object-reuse)
9. [Profiling Memory](#9-profiling-memory)
10. [Common Memory Optimization Patterns](#10-common-memory-optimization-patterns)
11. [Understanding Allocations with Benchmarks](#11-understanding-allocations-with-benchmarks)
12. [Quick Reference Card](#12-quick-reference-card)

---

## 1. Stack vs Heap — The Fundamental Decision

Every value lives in one of two places. The compiler chooses at compile time.

### Stack — Per-Goroutine, Ultra-Fast

Each goroutine owns a private stack. Allocation is just adjusting the stack pointer
(SP) — a single CPU instruction. Deallocation is free: function return restores SP.

```
  G1 Stack (2KB-1MB)            G2 Stack (independent)
  ┌──────────────────┐          ┌──────────────────┐
  │   main() frame   │          │  handler() frame │
  │  ┌────────────┐  │          │  ┌────────────┐  │
  │  │ local vars │  │          │  │ local vars │  │
  │  └────────────┘  │          │  └────────────┘  │
  ├──────────────────┤          ├──────────────────┤
  │  doWork() frame  │          │  query() frame   │
  │  ┌────────────┐  │          │  ┌────────────┐  │
  │  │ x := 42    │  │          │  │ buf [64]B  │  │
  │  │ y := 3.14  │  │          │  │ n := 0     │  │
  │  └────────────┘  │          │  └────────────┘  │
  ├──────────────────┤          ├──────────────────┤
  │  ▼ SP (grows     │          │  ▼ SP            │
  │    downward)     │          │                  │
  └──────────────────┘          └──────────────────┘
```

### Heap — Shared, GC-Managed

The heap is shared across all goroutines. The GC scans and reclaims it.

```
                  ┌───────────────────────── HEAP ──────────────────────┐
                  │  ┌───────────┐  ┌───────────┐  ┌───────────┐      │
    G1 ──────────►│  │ User{     │  │ []byte    │  │ *Config{} │◄──── G3
                  │  │  Name:"A" │  │ len=1024  │  │           │      │
    G2 ──────────►│  │ }         │  │           │  └───────────┘      │
                  │  └───────────┘  └───────────┘                      │
                  │  Managed by GC — concurrent mark and sweep         │
                  └────────────────────────────────────────────────────┘
```

### The Cost Gap

```
┌───────────────────┬───────────────────┬──────────────────────────────┐
│ Property          │ Stack             │ Heap                         │
├───────────────────┼───────────────────┼──────────────────────────────┤
│ Allocation cost   │ ~1ns (SP adjust)  │ ~25-50ns (runtime.mallocgc)  │
│ Deallocation      │ Free (SP restore) │ GC background cost           │
│ Synchronization   │ None (private)    │ GC coordination              │
│ Access speed      │ Hot in L1 cache   │ Pointer chase, cache miss    │
│ Who decides?      │ Compiler (escape  │ Compiler (escape analysis)   │
│                   │ analysis)         │                              │
└───────────────────┴───────────────────┴──────────────────────────────┘
```

**Source:** `runtime/stack.go` — goroutine stack management. Stacks start at 2-8KB,
grow by copying to a 2x buffer (detected at function preamble), shrink during GC
if <25% used.

---

## 2. Escape Analysis — How the Compiler Decides

Escape analysis answers: "Can this value live on the stack, or must it escape to heap?"

**Source:** `cmd/compile/internal/escape/`

```bash
go build -gcflags='-m' ./...        # what escapes
go build -gcflags='-m -m' ./...     # WHY it escapes
```

### Rule 1 — Pointer Outlives Function → Escapes

```go
func newUser(name string) *User {
    u := User{Name: name}       // ← escapes to heap
    return &u                   // pointer survives function return
}
// gcflags: "moved to heap: u"

func processUser(name string) string {
    u := User{Name: name}      // ← stays on stack
    return u.Name              // value copied, no pointer escape
}
```

### Rule 2 — Assigned to Interface → Usually Escapes

```go
func greet(u User) {
    fmt.Println(u)              // ← u escapes to heap
}
// Println(a ...any) — u boxed into eface{_type, data=unsafe.Pointer(&u)}
// Compiler can't prove Println won't store the pointer
```

### Rule 3 — Sent to Channel → Escapes

```go
func produce(ch chan *Data) {
    d := Data{Value: 42}       // ← escapes to heap
    ch <- &d                   // receiving goroutine may outlive sender
}
```

### Rule 4 — Closure Captures Escaping Variable → Escapes

```go
func makeCounter() func() int {
    count := 0                 // ← escapes to heap
    return func() int {        // closure outlives makeCounter
        count++
        return count
    }
}
```

### Rule 5 — Too Large for Stack → Escapes

```go
func bigAlloc() {
    buf := make([]byte, 1<<20)  // 1MB — escapes to heap
    _ = buf
}
```

### Common Surprise Escapes

```
┌─────────────────────────────────┬──────────────────────────────────────┐
│ Code                            │ Why It Escapes                       │
├─────────────────────────────────┼──────────────────────────────────────┤
│ fmt.Println(x)                  │ x boxed into interface{} arg         │
│ fmt.Sprintf("%d", n)            │ n boxed, result string allocated     │
│ return &localVar                │ pointer outlives stack frame          │
│ go func() { use(x) }()         │ closure captures x, goroutine escapes│
│ someSlice = append(s, v)        │ growth may allocate new backing array│
│ errors.New("msg")               │ returns *errorString on heap         │
└─────────────────────────────────┴──────────────────────────────────────┘
```

---

## 3. The Memory Allocator — TCMalloc-Inspired

**Source:** `runtime/malloc.go`, `runtime/mcache.go`, `runtime/mcentral.go`, `runtime/mheap.go`

### Three Size Classes

```
  Tiny (<16B, no ptrs)    Small (16B-32KB)        Large (>32KB)
  ────────────────────    ──────────────────      ──────────────
  Packed into 16-byte     ~67 size classes         Direct from mheap
  blocks. Multiple        (8,16,32,48,64...       Rounded to page
  small objects share     up to 32KB)              multiple (8KB pages)
  one allocation.
```

### The Three-Tier Allocation Path

```
  runtime.mallocgc(size, typ, needzero)
  │
  │  size > 32KB?
  │  ├─ YES ──────────────────────────────────► mheap (global, locked) → OS
  │  └─ NO
  │      │
  ▼      ▼
  ┌────────────────────────────────────────────────────────────────────┐
  │ TIER 1: mcache (per-P, NO LOCK)                                    │
  │  Each P has a private mcache with one mspan per size class.        │
  │  mcache.alloc[sizeclass] → mspan has free slot? → return. (~25ns)  │
  └─────────┬──────────────────────────────────────────────────────────┘
            │ span full
            ▼
  ┌────────────────────────────────────────────────────────────────────┐
  │ TIER 2: mcentral (per-size-class, LOCKED)                          │
  │  Global pool of partial/full spans per size class.                 │
  │  cacheSpan() → find partial span → move to mcache. (~100ns)        │
  └─────────┬──────────────────────────────────────────────────────────┘
            │ no spans available
            ▼
  ┌────────────────────────────────────────────────────────────────────┐
  │ TIER 3: mheap (global, LOCKED)                                     │
  │  Allocates new spans from free pages or requests from OS.          │
  │  mheap.alloc(npages) → carve span. (~500ns, or μs if OS call)     │
  └────────────────────────────────────────────────────────────────────┘
```

### Spans — The Unit of Memory Management

An `mspan` (`runtime/mheap.go`) is a contiguous run of 8KB pages divided into
fixed-size slots of a given size class.

```
  mspan (size class = 32 bytes, 256 objects per span)
  ┌──────┬──────┬──────┬──────┬──────┬──────┬─────────┐
  │ obj0 │ obj1 │ obj2 │ obj3 │ obj4 │ obj5 │ ...     │
  │ 32B  │ 32B  │ 32B  │ 32B  │ 32B  │ 32B  │         │
  │ used │ FREE │ used │ FREE │ used │ FREE │         │
  └──────┴──────┴──────┴──────┴──────┴──────┴─────────┘
  allocBits:  1  0  1  0  1  0  ...   (tracks which slots are in use)
  gcmarkBits: 1  0  1  0  0  0  ...   (tracks which are marked alive)
```

---

## 4. The Garbage Collector — Tri-Color Mark and Sweep

Concurrent, non-generational, non-compacting. Objects don't move once allocated.

**Source:** `runtime/mgc.go`, `runtime/mgcmark.go`, `runtime/mgcsweep.go`

### The Tri-Color Abstraction

```
  WHITE — not yet seen. If still white after marking → GARBAGE.
  GREY  — seen, but children (pointers) NOT yet scanned.
  BLACK — fully scanned. All children are grey or black. ALIVE.
```

### Step-by-Step: A GC Cycle

```
  STEP 1 — Mark Roots (STW pause #1, <200μs)
  ────────────────────────────────────────────
  Stop all goroutines. Enable write barrier.
  Mark root-reachable objects GREY.

  Stacks:        Heap objects:
  ┌─────┐       ┌───┐  ┌───┐  ┌───┐  ┌───┐  ┌───┐
  │  ●──┼──────►│ A │  │ B │  │ C │  │ D │  │ E │
  │  ●──┼──┐    │GRY│  │WHT│  │WHT│  │WHT│  │WHT│
  └─────┘  │    └─┬─┘  └───┘  └───┘  └───┘  └───┘
           │      │ points to B
           │      ▼
           │    ┌───┐
           └───►│ F │   A.ref→B,  B.ref→C, B.ref→D
                │GRY│   E has no references (garbage)
                └───┘
```

```
  STEP 2 — Concurrent Mark (goroutines resume)
  ─────────────────────────────────────────────
  GC goroutines scan grey objects alongside application.

  Scan A (grey→black): A→B → mark B grey
  ┌───┐  ┌───┐  ┌───┐  ┌───┐  ┌───┐  ┌───┐
  │ A │─►│ B │  │ C │  │ D │  │ E │  │ F │
  │BLK│  │GRY│  │WHT│  │WHT│  │WHT│  │GRY│
  └───┘  └─┬─┘  └───┘  └───┘  └───┘  └───┘

  Scan B (grey→black): B→C, B→D → mark C,D grey
  ┌───┐  ┌───┐  ┌───┐  ┌───┐  ┌───┐  ┌───┐
  │ A │─►│ B │─►│ C │  │ D │  │ E │  │ F │
  │BLK│  │BLK│  │GRY│  │GRY│  │WHT│  │GRY│
  └───┘  └───┘  └───┘  └───┘  └───┘  └───┘

  Scan C, D, F (grey→black): no new children
  ┌───┐  ┌───┐  ┌───┐  ┌───┐  ┌───┐  ┌───┐
  │ A │  │ B │  │ C │  │ D │  │ E │  │ F │
  │BLK│  │BLK│  │BLK│  │BLK│  │WHT│  │BLK│
  └───┘  └───┘  └───┘  └───┘  └───┘  └───┘
  No more grey → marking complete.
```

```
  STEP 3 — Mark Termination (STW pause #2, <200μs)
  ────────────────────────────────────────────────
  Stop goroutines. Disable write barrier. Finalize marking.

  STEP 4 — Concurrent Sweep (goroutines resume)
  ──────────────────────────────────────────────
  Reclaim WHITE objects. E is freed.
  ┌───┐  ┌───┐  ┌───┐  ┌───┐  ┌───┐  ┌───┐
  │ A │  │ B │  │ C │  │ D │  │ E │  │ F │
  │ ✓ │  │ ✓ │  │ ✓ │  │ ✓ │  │ ✗ │  │ ✓ │
  │live│  │live│  │live│  │live│  │FREE│  │live│
  └───┘  └───┘  └───┘  └───┘  └───┘  └───┘
```

### The Tri-Color Invariant

Black never points directly to white — enforced by the write barrier:

```
  ✅ ALLOWED:               ❌ VIOLATED:
  BLACK ──► GREY            BLACK ──► WHITE
  BLACK ──► BLACK           (write barrier prevents this)
  GREY  ──► WHITE/GREY/BLACK
```

---

## 5. Write Barrier — Protecting Concurrent Marking

**Source:** `runtime/mbarrier.go`

### The Problem Without It

```
  Time 0: GC scans A (→black). A.ref→B, so B marked grey.
  Time 1: App runs: A.ref = C; B.ref = nil  (moved pointer to C)
  Time 2: GC scans B (→black): no children.

  ┌───┐       ┌───┐       ┌───┐
  │ A │──────►│ C │       │ B │
  │BLK│       │WHT│ ☠️    │BLK│
  └───┘       └───┘       └───┘
  C is ALIVE (A→C) but WHITE → GC frees it → dangling pointer!
```

### Go's Hybrid Barrier (Dijkstra + Yuasa, since Go 1.8)

```
  writePointer(slot, new):
      shade(new)    // Dijkstra: if new is white → grey
      shade(*slot)  // Yuasa: if old is white → grey
      *slot = new   // actual pointer write
```

The barrier fires on every pointer store during concurrent mark phase.

```
┌───────────────────────────────────┬──────────────────────────────────┐
│ GC Phase                          │ Write Barrier Cost               │
├───────────────────────────────────┼──────────────────────────────────┤
│ Between GC cycles (sweep/idle)    │ OFF — zero overhead              │
│ Concurrent Mark                   │ ON — ~5% throughput overhead     │
└───────────────────────────────────┴──────────────────────────────────┘
```

---

## 6. GC Pacing — GOGC and GOMEMLIMIT

**Source:** `runtime/mgc.go` — `gcController`

### GOGC — The Growth Ratio (default: 100)

```
  GC trigger = live_heap × (1 + GOGC/100)

  Live heap = 100MB:
    GOGC=100:  trigger at 200MB  (2× live)
    GOGC=200:  trigger at 300MB  (3× live)
    GOGC=50:   trigger at 150MB  (1.5× live)
    GOGC=off:  GC disabled
```

```
┌───────────────┬──────────────────────┬───────────────────────────────┐
│ GOGC Value    │ GC Frequency         │ Tradeoff                      │
├───────────────┼──────────────────────┼───────────────────────────────┤
│ 50            │ Very frequent        │ Low memory, high CPU          │
│ 100 (default) │ Balanced             │ 2× live heap, moderate CPU    │
│ 200           │ Less frequent        │ 3× live heap, lower CPU       │
│ off           │ Never                │ Unbounded memory, zero GC CPU │
└───────────────┴──────────────────────┴───────────────────────────────┘
```

### GOMEMLIMIT — Soft Memory Limit (Go 1.19+)

As heap approaches the limit, GC runs more aggressively (dynamically lowers GOGC).

```
  Heap                                             ← container limit (512MB)
  Size
  (MB)
   450 ┤ ═══════════════════ GOMEMLIMIT ═══════════ ← soft target
       │        ╱╲        GC more aggressive here
   400 ┤──────╱──╲── normal GOGC trigger
       │    ╱      ╲
   300 ┤  ╱          ╲
       │╱              ╲
   200 ┤                 ╲─── live heap
       │
    0  ┼──────────────────────► time
```

```bash
# Container recipe (1GB limit):
GOMEMLIMIT=900MiB    # ~90% of container limit
GOGC=100             # or GOGC=off (GC runs only near limit)
```

**Replaces the ballast pattern** — no more allocating large `[]byte` to delay GC.

---

## 7. GC Phases and STW Pauses

```
  ◄────────────────── One GC Cycle ──────────────────────►
  ┌─────────┐ ┌────────────────────────┐ ┌─────────┐ ┌──────────────┐
  │ Mark    │ │  Concurrent Marking    │ │ Mark    │ │ Concurrent   │
  │ Setup   │ │  GC + app together     │ │ Termin. │ │ Sweep        │
  │ (STW)   │ │  write barrier ON      │ │ (STW)   │ │              │
  │ <200μs  │ │  ~5% overhead          │ │ <200μs  │ │ lazy reclaim │
  └─────────┘ └────────────────────────┘ └─────────┘ └──────────────┘
```

```
┌───────────────────┬─────────────┬────────────────────────────────────┐
│ Phase             │ Duration    │ What Happens                       │
├───────────────────┼─────────────┼────────────────────────────────────┤
│ Mark Setup        │ <200μs STW  │ Enable write barrier, enqueue      │
│                   │             │ root scan jobs, resume goroutines   │
├───────────────────┼─────────────┼────────────────────────────────────┤
│ Concurrent Mark   │ ms to 100s  │ GC goroutines (25% GOMAXPROCS) +   │
│                   │ of ms       │ mark assist from allocating Gs.    │
│                   │             │ Scan stacks, scan heap objects.    │
├───────────────────┼─────────────┼────────────────────────────────────┤
│ Mark Termination  │ <200μs STW  │ Drain remaining work, disable      │
│                   │             │ write barrier, compute next trigger│
├───────────────────┼─────────────┼────────────────────────────────────┤
│ Concurrent Sweep  │ Until next  │ Spans swept lazily on allocation.  │
│                   │ cycle       │ Returns pages to mheap / OS.       │
└───────────────────┴─────────────┴────────────────────────────────────┘
```

**Mark assist:** If a goroutine allocates faster than GC can mark, the runtime forces
it to help mark objects before its allocation proceeds — latency spike on that G.
Reduce allocations in hot paths to avoid the "mark assist tax."

---

## 8. sync.Pool — GC-Aware Object Reuse

Per-P object cache. Objects may be silently removed by GC — NOT a connection pool.

**Source:** `sync/pool.go` (registered with GC via `poolCleanup`)

### Internal Structure

```
  P0                   P1                   P2
  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
  │ private: *obj   │  │ private: *obj   │  │ private: *obj   │
  │ (fast, no lock) │  │ (fast, no lock) │  │ (fast, no lock) │
  ├─────────────────┤  ├─────────────────┤  ├─────────────────┤
  │ shared: []obj   │  │ shared: []obj   │  │ shared: []obj   │
  │ (lock-free,     │  │ (stealable)     │  │ (stealable)     │
  │  stealable)     │  │                 │  │                 │
  └─────────────────┘  └─────────────────┘  └─────────────────┘
         │                    │                    │
         ▼                    ▼                    ▼
  ┌──────────────────────────────────────────────────────────┐
  │     VICTIM CACHE (objects from previous GC cycle)         │
  │     Survive ONE more cycle, then dropped                  │
  └──────────────────────────────────────────────────────────┘
```

### Get/Put Flow

```
  pool.Get():
    1. private (current P) ─── found? → return (no lock)
    2. shared (current P)  ─── found? → return (lock-free pop)
    3. steal other P's shared ── found? → return
    4. victim cache ──────────── found? → return
    5. pool.New() ────────────── fresh allocation

  pool.Put(obj):
    1. private empty? → store there (no lock)
    2. else → append to shared (lock-free push)
```

### GC Interaction — Two-Generation Scheme

Each GC cycle: current pools → victim cache, previous victim → dropped.
Objects survive at most **two GC cycles** in the pool.

```
┌────────────────────────────────┬──────────────────────────────────────┐
│ ✅ Use For                      │ ❌ Don't Use For                      │
├────────────────────────────────┼──────────────────────────────────────┤
│ Temporary byte buffers         │ DB connection pools                  │
│ Serialization scratch space    │ Objects with lifecycle/state         │
│ bytes.Buffer / strings.Builder │ Anything requiring Close()           │
│ Encoder/decoder state          │ Objects that must persist across GC  │
└────────────────────────────────┴──────────────────────────────────────┘
```

---

## 9. Profiling Memory

### Heap Profile

```bash
go test -memprofile=mem.out -bench=BenchmarkXxx ./...
go tool pprof mem.out                   # interactive
go tool pprof -alloc_space mem.out      # total allocated (find hot spots)
go tool pprof -inuse_space mem.out      # currently alive (find leaks)
```

### runtime.ReadMemStats

```go
var m runtime.MemStats
runtime.ReadMemStats(&m)
// m.HeapAlloc    — bytes currently allocated (live objects)
// m.HeapSys      — bytes obtained from OS for heap
// m.NumGC        — completed GC cycles
// m.PauseTotalNs — total STW pause time
// m.Mallocs      — cumulative heap allocations
// m.Frees        — cumulative heap frees
```

### GODEBUG=gctrace=1

```bash
GODEBUG=gctrace=1 ./myservice
```

```
gc 1 @0.012s 2%: 0.044+1.2+0.030 ms clock, 0.35+0.82/1.8/0.15+0.24 ms cpu, 4->4->2 MB, 4 MB goal, 8 P
                  │          │                                          │
                  │          │                                          └─ heap: before→after→live
                  │          └─ CPU time: assist/background/idle
                  └─ wall clock: STW1 + concurrent + STW2
```

### Live Production Profiling

```go
import _ "net/http/pprof"  // registers /debug/pprof/ endpoints
// ALWAYS behind authentication!
```

```bash
go tool pprof http://localhost:6060/debug/pprof/heap
go tool pprof -diff_base=before.out after.out  # compare snapshots
```

---

## 10. Common Memory Optimization Patterns

### Pre-Allocate Slices and Maps

```go
// BAD — repeated grow: 0→1→2→4→8→16→32→64→128
results := []User{}
for _, row := range rows { results = append(results, parseUser(row)) }

// GOOD — single allocation
results := make([]User, 0, len(rows))
m := make(map[string]int, expectedSize)
```

### Avoid String Concatenation in Loops

```go
// BAD — O(n²)                        // GOOD — O(n)
s := ""                                var b strings.Builder
for _, item := range items {           b.Grow(len(items) * 20)
    s += item.Name + ","               for i, item := range items {
}                                          if i > 0 { b.WriteByte(',') }
                                           b.WriteString(item.Name)
                                       }
                                       s := b.String()
```

### Avoid fmt.Sprintf in Hot Paths

```go
// BAD — 2+ allocs                     // GOOD — 0-1 allocs
key := fmt.Sprintf("user:%d", id)      key := "user:" + strconv.Itoa(id)
```

### Reduce Pointer-Heavy Structures

```go
// MORE GC WORK — every pointer scanned    // LESS GC WORK — inline values
type Record struct {                        type Record struct {
    Name    *string                             Name     string
    Tags    []*string                           Tags     []string
    Parent  *Record                             ParentID int64
}                                           }
```

### Struct Field Ordering — Reduce Padding

```
  BadLayout (32B):                       GoodLayout (24B):
  ┌──┬───────┬────────┬──┬───────┬───────┐  ┌────────┬────────┬──┬──┬──────┐
  │a │ 7pad  │   b    │c │ 7pad  │   d   │  │   b    │   d    │a │c │ 6pad │
  │1B│       │  8B    │1B│       │  8B   │  │  8B    │  8B    │1B│1B│      │
  └──┴───────┴────────┴──┴───────┴───────┘  └────────┴────────┴──┴──┴──────┘
  bool,int64,bool,int64 = 32B               int64,int64,bool,bool = 24B (25% less)
```

---

## 11. Understanding Allocations with Benchmarks

### b.ReportAllocs

```go
func BenchmarkProcess(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ { process(testData) }
}
// Output: BenchmarkProcess-8  1000000  1050 ns/op  256 B/op  3 allocs/op
```

### testing.AllocsPerRun

```go
func TestZeroAlloc(t *testing.T) {
    allocs := testing.AllocsPerRun(100, func() { fastPath(input) })
    if allocs > 0 { t.Errorf("expected zero allocs, got %f", allocs) }
}
```

### Zero-Allocation Hot Path Pattern

```go
// Caller owns the buffer, reuses across calls
func formatKey(buf []byte, prefix string, id int) []byte {
    buf = buf[:0]                               // reset, keep capacity
    buf = append(buf, prefix...)
    buf = append(buf, ':')
    buf = strconv.AppendInt(buf, int64(id), 10)
    return buf
}
```

### Allocation Avoidance Quick Reference

```
┌──────────────────────────────────┬──────────────────────────────────────┐
│ Pattern That Allocates           │ Zero-Alloc Alternative               │
├──────────────────────────────────┼──────────────────────────────────────┤
│ fmt.Sprintf("%d", n)             │ strconv.Itoa(n) / AppendInt          │
│ fmt.Println(x)                   │ io.WriteString(w, s)                 │
│ return &localVar                 │ return value type (copy)             │
│ interface conversion in hot path │ use concrete type                    │
│ append beyond cap                │ pre-allocate: make([]T, 0, n)        │
│ errors.New("msg") per call       │ package-level sentinel error         │
│ closure capturing variable       │ pass as parameter instead            │
└──────────────────────────────────┴──────────────────────────────────────┘
```

---

## 12. Quick Reference Card

```
ESCAPE ANALYSIS
───────────────
  Triggers: return &x, interface boxing, channel send, escaping closure, too large
  Tools:    go build -gcflags='-m' ./...    (what escapes)
            go build -gcflags='-m -m' ./... (why it escapes)

MEMORY ALLOCATOR  (runtime/malloc.go)
────────────────
  Tiny (<16B)  → packed into 16-byte blocks
  Small (≤32KB) → mcache (per-P, no lock) → mcentral (locked) → mheap → OS
  Large (>32KB) → mheap directly → OS

GC PHASES  (runtime/mgc.go)
─────────
  1. Mark Setup      — STW <200μs — enable write barrier
  2. Concurrent Mark — app + GC together, ~5% overhead
  3. Mark Termination — STW <200μs — disable write barrier
  4. Concurrent Sweep — reclaim dead objects, fully concurrent
  Total STW: typically <1ms per cycle

GC TUNING
─────────
  GOGC=100 (default)   2× live heap.   Balanced.
  GOGC=50              1.5× live heap.  Less mem, more CPU.
  GOGC=200             3× live heap.    More mem, less CPU.
  GOMEMLIMIT=900MiB    Soft limit, GC aggressive near limit.
  Container: GOMEMLIMIT = container_memory × 0.9

PROFILING
─────────
  go test -memprofile=mem.out -bench=.     # heap profile
  go tool pprof -alloc_space mem.out       # where allocs happen
  go tool pprof -inuse_space mem.out       # what's holding memory
  GODEBUG=gctrace=1 ./myservice            # per-cycle stats
  go test -benchmem -bench=. ./...         # allocs per op

OPTIMIZATION CHECKLIST
──────────────────────
  □ Pre-allocate slices/maps with make(T, 0, n)
  □ sync.Pool for temporary buffers
  □ strings.Builder not += in loops
  □ strconv not fmt.Sprintf in hot paths
  □ Return values not pointers (avoid escape)
  □ Concrete types in hot paths (avoid interface boxing)
  □ Order struct fields by size (reduce padding)
  □ b.ReportAllocs() in every benchmark
```

---

## One-Line Summary

> The compiler's escape analysis decides stack (~1ns) vs heap (~25-50ns) at build
> time; the concurrent tri-color GC reclaims heap with <1ms STW pauses; tune via
> `GOGC`/`GOMEMLIMIT`, profile with `pprof`/`gctrace`, and reduce allocations —
> every heap alloc you avoid is a direct throughput win.
