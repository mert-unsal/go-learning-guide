# Deep Dive: Zero Values, mallocgc, sync.Pool & duffzero

> How Go's zero value guarantee connects to the runtime's memory zeroing
> machinery, why sync.Pool deliberately breaks it, and how duffzero clears
> memory at near-hardware speed.
>
> **Builds on:** [Chapter 13](./13_memory_gc_escape_sorting.md) covers the
> allocator hierarchy, GC, and escape analysis. This chapter goes deeper into
> the zeroing pipeline itself.

---

## Table of Contents

1. [Zero Values Are a Runtime Contract, Not Just Syntax](#1-zero-values-are-a-runtime-contract-not-just-syntax)
2. [The 2KB Stack Connection](#2-the-2kb-stack-connection)
3. [mallocgc: Where Zeroing Actually Happens](#3-mallocgc-where-zeroing-actually-happens)
4. [needzero on mspan: Lazy Zeroing for Performance](#4-needzero-on-mspan-lazy-zeroing-for-performance)
5. [duffzero: Zeroing Memory at Assembly Level](#5-duffzero-zeroing-memory-at-assembly-level)
6. [sync.Pool: The Deliberate Zero Value Violation](#6-syncpool-the-deliberate-zero-value-violation)
7. [Production Patterns: Working With and Around Zeroing](#7-production-patterns-working-with-and-around-zeroing)
8. [Quick Reference Card](#8-quick-reference-card)

---

## 1. Zero Values Are a Runtime Contract, Not Just Syntax

When you write `var x int`, the Go specification guarantees `x` is `0`. This is
not just a language convenience. It is a contract enforced at every layer of the
runtime: the compiler, the stack manager, and the heap allocator all cooperate
to guarantee that freshly allocated memory is always zeroed.

### Why This Matters Architecturally

Most languages handle uninitialized memory in one of three ways:

```
  C/C++:    Uninitialized memory contains GARBAGE
            → Fast, but undefined behavior bugs, security vulnerabilities
            → Heartbleed (2014) was reading uninitialized heap memory

  Java/C#:  Default initialization via constructors + VM zeroing
            → Safe, but constructor cost per object, complex initialization chains
            → NullPointerException at runtime for reference types

  Go:       ALL memory is zeroed before use, ALWAYS
            → Safe: no garbage reads, no undefined behavior
            → Cheap: zeroing is one memclr call, no constructor logic
            → Deterministic: same zero value every time, every platform
```

The Go team chose zeroing because it sits at the sweet spot: nearly as fast as
C's "do nothing" approach, but with the safety guarantees of managed languages.

### The Cascading Design Decision

Zero values are not an isolated feature. They cascade through the entire
language design:

```
  Cheap goroutine stacks (2KB, millions of them)
    └─► Stack setup must be nearly free
        └─► Just zero the memory, no constructors needed
            └─► Every type MUST have a meaningful zero value
                └─► Standard library designed around useful zero values
                    └─► sync.Mutex{}, bytes.Buffer{}, http.Server{} all work at zero
                        └─► No "forgot to call constructor" class of bugs
```

This chain is why Rob Pike's proverb says "Make the zero value useful." It is
not just style advice. It is the logical consequence of Go's memory model.

### What Zero Means for Each Type Category

```
  ┌─────────────────────┬──────────────┬──────────────────────────────┐
  │ Type Category       │ Zero Value   │ Under the Hood               │
  ├─────────────────────┼──────────────┼──────────────────────────────┤
  │ bool                │ false        │ 1 byte: 0x00                 │
  │ int, float, complex │ 0            │ all bytes 0x00               │
  │ string              │ ""           │ StringHeader{ptr:nil, len:0} │
  │ pointer             │ nil          │ 8 bytes: 0x00 (64-bit)       │
  │ slice               │ nil          │ {ptr:nil, len:0, cap:0}      │
  │ map                 │ nil          │ pointer: 0x00                │
  │ channel             │ nil          │ pointer: 0x00                │
  │ interface           │ nil          │ {type:nil, data:nil}         │
  │ function            │ nil          │ pointer: 0x00                │
  │ struct              │ recursive    │ each field zeroed recursively │
  │ array               │ recursive    │ each element zeroed           │
  └─────────────────────┴──────────────┴──────────────────────────────┘
```

The critical insight: **all zero values are literally all-zero bytes**. This is
not a coincidence. Go's type system is designed so that the bit pattern
`0x000...0` is always a valid, meaningful value for every type. This means the
runtime can use a single `memclr` operation to initialize any type, regardless
of its structure.

---

## 2. The 2KB Stack Connection

### Why the Stack Size Matters for Zeroing

Go goroutines start with a 2KB stack (since Go 1.4). Compare:

```
  OS thread (Linux default):  8MB stack     → 10,000 threads  = 80 GB
  OS thread (custom):         1MB stack     → 10,000 threads  = 10 GB
  Go goroutine:               2KB stack     → 10,000 routines = 20 MB
  Go goroutine:               2KB stack     → 1,000,000       = 2 GB
```

The 2KB size only works because goroutine creation is nearly free. The cost
breakdown of creating a goroutine:

```
  1. Allocate 2KB from the heap             ~25-50ns (runtime.mallocgc)
  2. Zero the stack memory                  ~2-5ns  (memclr for 2KB)
  3. Set up the goroutine descriptor (g)    ~10ns   (runtime.newproc1)
  4. Put g on the run queue                 ~5ns    (runqput)
  ─────────────────────────────────────────────────
  Total:                                    ~50-70ns per goroutine
```

Step 2 is the zeroing. If Go used constructor-based initialization (like Java),
this step would involve walking the type hierarchy, calling initialization
methods, setting up vtables. Instead, it is one `memclr` call.

### Stack Frame Zeroing

When a function is called, the compiler knows the exact size of its stack frame
(all local variables, return values, spill slots). The function preamble zeros
this frame:

```go
func process(data []byte) error {
    var buf [256]byte       // 256 bytes, zeroed
    var count int           // 8 bytes, zeroed
    var result strings.Builder  // 24 bytes, zeroed (all fields zero)
    // ...
}
```

The compiler generates a preamble that looks like this in pseudo-assembly:

```asm
  process:
    SUBQ  $304, SP          // grow stack by 304 bytes (256+8+24+padding+...)
    // Zero the frame:
    LEAQ  (SP), DI           // destination = stack pointer
    MOVQ  $38, CX            // 304/8 = 38 quadwords
    XORQ  AX, AX             // AX = 0
    REP   STOSQ              // zero 38 quadwords
    // ... function body ...
```

For small frames (under ~128 bytes), the compiler uses individual `MOVQ $0`
instructions instead. For larger frames, it calls `runtime.duffzero` or uses
`REP STOSQ`.

### The Dead Store Optimization

The compiler is smart enough to skip zeroing when it is unnecessary:

```go
func example() int {
    var x int       // normally zeroed
    x = 42          // immediately overwritten
    return x
}
```

The compiler sees that the zero value of `x` is never read, so it eliminates
the zeroing entirely (dead store elimination). You can verify with:

```bash
go build -gcflags='-m -m' ./...
```

This means the zero value guarantee has near-zero cost in practice for
variables that are immediately assigned.

---

## 3. mallocgc: Where Zeroing Actually Happens

**Source:** `runtime/malloc.go`

`mallocgc` is the single entry point for ALL heap allocations in Go. Every
`new()`, `make()`, composite literal, and escaped variable ends up here.

### Function Signature

```go
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer
```

Three parameters:
- `size`: how many bytes to allocate
- `typ`: type metadata (pointer layout for GC scanning, or nil for untyped)
- `needzero`: whether the caller needs the memory to be zeroed

### The Three Allocation Paths

The allocator takes different paths based on object size. Each path handles
zeroing differently:

```
  mallocgc(size, typ, needzero)
  │
  ├─ size == 0?
  │   └─ return &zerobase (a global zero-sized object, no allocation)
  │
  ├─ Tiny path: size < 16 bytes AND no pointers AND needzero == false
  │   └─ Pack multiple tiny objects into a single 16-byte block
  │      (uses mcache.tiny, mcache.tinyoffset)
  │
  ├─ Small path: size <= 32KB
  │   └─ Look up size class → get mspan from mcache
  │      → allocate from mspan's free list
  │
  └─ Large path: size > 32KB
      └─ Allocate directly from mheap (page-aligned)
```

### Tiny Allocation Path (< 16 bytes, no pointers)

This is Go's most aggressive optimization. Small non-pointer values (like small
integers, bools, single bytes) get packed together:

```
  mcache.tiny ──► ┌────────────────┐  16-byte block
                  │ obj1 │ obj2 │··│
                  │  3B  │  1B  │  │  ← tinyoffset tracks next free byte
                  └────────────────┘

  If a new 4-byte object arrives and fits:
                  ┌────────────────┐
                  │ obj1 │ obj2 │o3│
                  │  3B  │  1B  │4B│  ← packed, 8 bytes remaining
                  └────────────────┘
```

Zeroing for tiny objects: the 16-byte block is zeroed when first allocated from
mspan. Subsequent objects packed into it rely on the fact that the block was
already zeroed. This is why tiny allocation requires `needzero == false`: the
caller promises they will write the full object immediately.

### Small Allocation Path (16 bytes to 32KB)

```
  1. Determine size class (e.g., 24-byte request → size class 3, 32 bytes)
  2. Get mspan from mcache.alloc[sizeclass]
  3. Find free slot using allocBits bitmap
  4. Check mspan.needzero
  5. If needzero && caller needs zero → zero the slot
  6. Return pointer to slot
```

The key decision point is step 4. If the mspan's `needzero` flag is set, the
allocator must zero the memory before returning it. More on this in section 4.

### Large Allocation Path (> 32KB)

```
  1. Round size up to page multiple (8KB pages)
  2. Lock mheap
  3. Find contiguous free pages (treap-based free list)
  4. If no pages: grow the heap (mmap from OS)
  5. Create mspan for this allocation
  6. Zero if needed (large objects always check mspan.needzero)
  7. Unlock mheap, return pointer
```

For very large allocations, the OS may provide already-zeroed pages through
`mmap(MAP_ANONYMOUS)`, which lets the allocator skip explicit zeroing.

### The Zeroing Decision Tree Inside mallocgc

```
  Is memory fresh from OS (mmap)?
  ├─ YES → Already zeroed by kernel (security guarantee) → skip zeroing
  └─ NO  → Memory recycled from freed objects
           │
           Is mspan.needzero set?
           ├─ YES → Memory contains stale data → MUST zero
           └─ NO  → Memory was zeroed when span was last swept → safe
```

This is "lazy zeroing": the allocator avoids zeroing until absolutely necessary,
and tracks whether it has already been done using the `needzero` flag on each
mspan.

---

## 4. needzero on mspan: Lazy Zeroing for Performance

**Source:** `runtime/mheap.go`

### The mspan Structure (Relevant Fields)

```go
type mspan struct {
    startAddr uintptr     // start of the span's memory
    npages    uintptr     // number of 8KB pages
    nelems    uintptr     // number of objects in this span

    allocBits  *gcBits    // bitmap: which slots are allocated
    gcmarkBits *gcBits    // bitmap: which slots are marked by GC

    sweepgen   uint32     // sweep generation (for concurrent sweep)
    spanclass  spanClass  // size class + noscan flag

    needzero   uint8      // THIS IS THE KEY FIELD
    // if needzero == 1: free slots contain stale data, must zero on alloc
    // if needzero == 0: free slots are already zeroed, safe to return as-is
}
```

### When needzero Gets Set

```
  mspan lifecycle:
  ┌──────────────────────────────────────────────────────────────────┐
  │                                                                  │
  │  1. Fresh from OS (mmap)                                        │
  │     needzero = 0  (OS zeroed the pages)                         │
  │                                                                  │
  │  2. Objects allocated, used, then freed by GC sweep              │
  │     needzero = 1  (freed slots contain old object data)         │
  │                                                                  │
  │  3. Span returned to mcentral with partial free slots            │
  │     needzero stays 1  (some slots have stale data)              │
  │                                                                  │
  │  4. Span fully freed, returned to mheap page cache               │
  │     needzero = 1  (all slots have stale data)                   │
  │                                                                  │
  │  5. Span reused for a DIFFERENT size class                       │
  │     needzero = 1  (layout changed, all memory is stale)         │
  │                                                                  │
  └──────────────────────────────────────────────────────────────────┘
```

### The Performance Win

Without lazy zeroing, every allocation would zero memory even if it was already
clean. With `needzero`, the allocator can skip zeroing for freshly mmap'd pages:

```
  Fresh page from OS:     allocate → return (0 zeroing cost)
  Recycled page:          allocate → zero slot → return (~2-10ns per slot)

  In a steady-state service:
  - Most allocations come from recycled spans (needzero = 1)
  - But the zeroing cost is per-slot (32 bytes, 64 bytes, etc.), not per-page
  - This is MUCH cheaper than zeroing an entire 8KB page upfront
```

### How the Sweeper Interacts With needzero

During GC sweep, the sweeper walks through spans and identifies dead objects:

```
  Before sweep:
    allocBits:   1 0 1 1 0 1 0 0
    gcmarkBits:  1 0 0 1 0 0 0 0
                      ↑        ↑
                      dead     dead (not marked)

  After sweep:
    allocBits:   1 0 0 1 0 0 0 0  (dead objects freed)
    needzero:    1                  (freed slots have stale data)
```

The sweeper does NOT zero the freed slots immediately. It just updates the
bitmap and sets `needzero = 1`. The actual zeroing happens later, when
`mallocgc` reuses those slots. This lazy approach spreads the zeroing cost
across future allocations rather than concentrating it during GC.

---

## 5. duffzero: Zeroing Memory at Assembly Level

**Source:** `runtime/duff_amd64.s`

### What Is Duff's Device?

Tom Duff invented this technique at Lucasfilm in 1983. The idea: create an
unrolled loop with multiple entry points, then jump into the middle based on how
many iterations you need. This eliminates the loop overhead for the remainder.

Traditional loop:
```
  for i := 0; i < n; i++ {
      mem[i] = 0
  }
  // Each iteration: compare, branch, increment, store = 4 instructions
```

Duff's device:
```
  Jump to entry point based on n % unroll_factor
  Then execute full unrolled iterations for the rest
  // Per 8 elements: just 8 stores, 1 compare, 1 branch = 10 instructions
  // vs traditional: 32 instructions (4 * 8)
```

### How Go Uses duffzero

Go's `runtime/duff_amd64.s` contains a zeroing function structured as a long
sequence of `MOVUPS` (128-bit store of zeros) instructions with `ADDQ` to
advance the pointer. The compiler computes an offset to jump into this sequence
at exactly the right point.

Simplified structure (the real code uses 128-bit SSE stores):

```asm
TEXT runtime·duffzero(SB), NOSPLIT|NOFRAME, $0-0
    // DI = pointer to memory to zero
    // Compiler jumps to a computed offset within this function

    MOVUPS  X15, 0(DI)       // zero bytes 0-15
    MOVUPS  X15, 16(DI)      // zero bytes 16-31
    ADDQ    $32, DI           // advance pointer

    MOVUPS  X15, 0(DI)       // zero bytes 32-47
    MOVUPS  X15, 16(DI)      // zero bytes 48-63
    ADDQ    $32, DI           // advance pointer

    // ... repeats for a total of 16 blocks (512 bytes total) ...

    MOVUPS  X15, 0(DI)       // zero bytes 480-495
    MOVUPS  X15, 16(DI)      // zero bytes 496-511
    ADDQ    $32, DI           // advance pointer

    RET
```

Each block zeroes 32 bytes (two 16-byte SSE stores + pointer advance). With 16
blocks, `duffzero` can zero up to 512 bytes in a single call with no loop
overhead. The X15 register is pre-zeroed (the SSE register holds 128 bits of
zeros).

### The Compiler's Decision: How to Zero

The compiler in `cmd/compile/internal/ssagen` picks the zeroing strategy based
on the size of memory to clear:

```
  Size             Strategy             Why
  ──────────────   ────────────────     ──────────────────────────────
  0 bytes          Nothing              No-op

  1-8 bytes        MOV $0, (addr)       Single instruction, inline
                                        Cheapest possible: 1 cycle

  9-32 bytes       2-4 MOV instructions Still inline, no function call
                                        Branch predictor loves this

  33-512 bytes     DUFFZERO             Jump into unrolled sequence
                                        No loop overhead, ~1 cycle per 32B
                                        Perfect for typical stack frames

  513 bytes-       REP STOSQ            Hardware-accelerated bulk zero
  large                                 Uses ERMSB (Enhanced REP MOV/STOS)
                                        on modern CPUs: ~40 GB/s throughput
                                        Best for large heap allocations
```

### Why Not Always Use REP STOSQ?

REP STOSQ has a startup cost. The CPU microcode needs to set up the fast-path
internally (prefetch, cache line alignment). For small sizes, this startup cost
dominates:

```
  Zeroing 64 bytes:
    4 MOV instructions:  ~4 cycles  (winner)
    DUFFZERO:            ~6 cycles  (small overhead for jump + return)
    REP STOSQ:           ~30 cycles (startup cost dominates)

  Zeroing 4KB:
    512 MOV instructions: ~512 cycles (too many instructions, icache pressure)
    DUFFZERO:             ~130 cycles (good, but function call overhead repeats)
    REP STOSQ:            ~50 cycles  (winner: startup amortized, hardware fast)
```

### Viewing duffzero in Your Own Code

You can see the compiler choosing duffzero by looking at the assembly output:

```bash
go build -gcflags='-S' ./cmd/concepts/basics/04-zero-values/ 2>&1 | grep -i duff
```

Look for lines like:
```
DUFFZERO $224
```

The `$224` is the offset into the duffzero function. The compiler computed it
from the size of memory to zero:

```
  offset = 512 - sizeToZero    (simplified)
  Jump to duffzero+offset → zeroes exactly sizeToZero bytes
```

---

## 6. sync.Pool: The Deliberate Zero Value Violation

**Source:** `sync/pool.go`

### The Core Problem

Everything we discussed above establishes Go's guarantee: freshly allocated
memory is always zeroed. But `sync.Pool` returns **previously used** objects.
These objects contain data from their last use. Pool does NOT zero them.

This is a deliberate violation of the zero value principle for performance.

### Why Pool Does Not Zero

```
  sync.Pool's purpose: avoid repeated heap allocation + GC pressure
  The win: reuse already-allocated memory without going through mallocgc

  If Pool zeroed objects on Get():
    ┌─────────────────────────────────────────────────────────────────┐
    │ 1. Get() retrieves object from per-P cache                     │
    │ 2. Zero the entire object (memclr)     ← ADDED COST           │
    │ 3. Return to caller                                            │
    │ 4. Caller initializes the object       ← ALSO HAPPENS         │
    │                                                                 │
    │ Result: zeroed twice (once by Pool, once by caller)             │
    │ Cost: wasted CPU cycles on the first zeroing                    │
    └─────────────────────────────────────────────────────────────────┘

  By NOT zeroing:
    ┌─────────────────────────────────────────────────────────────────┐
    │ 1. Get() retrieves object from per-P cache                     │
    │ 2. Return to caller immediately                                │
    │ 3. Caller resets/initializes as needed   ← ONLY ONCE           │
    │                                                                 │
    │ Result: caller controls what to reset, can be selective         │
    │ Win: only clear what you actually need to clear                 │
    └─────────────────────────────────────────────────────────────────┘
```

### The Security Implication

This non-zeroing behavior has real security consequences:

```go
// DANGEROUS: data leaks between requests
var bufPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    buf := bufPool.Get().(*bytes.Buffer)
    defer bufPool.Put(buf)

    // BUG: buf may contain data from a PREVIOUS request!
    // If the previous request wrote sensitive data (tokens, passwords),
    // and this handler reads from buf before writing, it gets old data.

    buf.WriteString("new data")
    // The buffer now contains: "old sensitive data" + "new data"
    w.Write(buf.Bytes())  // LEAKS old request's data!
}
```

The fix is always to reset before use:

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    buf := bufPool.Get().(*bytes.Buffer)
    defer bufPool.Put(buf)

    buf.Reset()  // CRITICAL: clear previous data
    buf.WriteString("new data")
    w.Write(buf.Bytes())
}
```

### The Standard Library's Pattern

The standard library uses sync.Pool extensively and always follows the
"reset on Get" pattern. Here is how `encoding/json` does it:

```go
// From encoding/json/encode.go (simplified)
var encodeStatePool sync.Pool

func newEncodeState() *encodeState {
    if v := encodeStatePool.Get(); v != nil {
        e := v.(*encodeState)
        e.Reset()              // ← always reset before use
        return e
    }
    return new(encodeState)
}
```

And `fmt` does the same:

```go
// From fmt/print.go (simplified)
var ppFree = sync.Pool{
    New: func() any { return new(pp) },
}

func newPrinter() *pp {
    p := ppFree.Get().(*pp)
    p.panicking = false        // ← explicit field reset
    p.erroring = false
    p.wrapErrs = false
    p.fmt.init(&p.buf)         // ← reinitialize formatter
    return p
}
```

Notice: `fmt` does not zero the entire `pp` struct. It resets only the fields
that matter. This is more efficient than a full `memclr` when the struct is
large but only a few fields need resetting.

### Pool Internals: Why Zeroing Would Be Technically Hard

The Pool stores objects as `any` (interface{}). It does not know the concrete
type or size of the stored object. To zero it, Pool would need to:

```
  1. Extract the type information from the interface header (iface.tab._type)
  2. Get the size from _type.Size_
  3. memclr(iface.data, _type.Size_)
```

This is technically possible but adds overhead to every Get() call. The Go
team decided the caller is in a better position to know what needs resetting.

### The Design Philosophy Tension

```
  Go's promise:     "Fresh memory is always zeroed"
  Pool's reality:   "Recycled memory is NOT zeroed"

  This creates a mental model split:
  ┌───────────────────────────────────────────────────────────────┐
  │                                                               │
  │  new(T)        →  zeroed (runtime guarantee)                 │
  │  make([]T, n)  →  zeroed (runtime guarantee)                 │
  │  var x T       →  zeroed (compiler guarantee)                │
  │  pool.Get()    →  STALE DATA (caller's responsibility)       │
  │                                                               │
  │  The rule: if you got it from a Pool, treat it as dirty.     │
  │                                                               │
  └───────────────────────────────────────────────────────────────┘
```

This is an intentional tradeoff. Pool exists for hot paths where allocation
pressure would cause GC latency spikes. In those paths, the caller is already
a performance-conscious developer who understands the lifecycle. The Go team
trusts this developer to reset correctly, and the payoff is significant:

```
  Without Pool (allocate each time):
    50,000 req/s × 4KB buffer = 200 MB/s allocation rate
    GC runs every ~100ms, 200ms p99 spikes

  With Pool (reuse buffers):
    50,000 req/s × 0 allocations (recycled) = near-zero GC pressure
    Smooth p99 latency
```

### Victim Cache: Objects Survive Two GC Cycles

When GC runs, Pool does not discard everything immediately. It uses a
two-generation scheme:

```
  GC Cycle 1:
    pool.local (current objects) → pool.victim (preserved for one more cycle)
    pool.victim (from last cycle) → DROPPED (garbage collected)

  GC Cycle 2:
    pool.local (new objects) → pool.victim
    pool.victim (objects from cycle 1) → DROPPED

  Timeline:
    ┌──────────┬────────────┬────────────┬────────────┐
    │          │  Cycle N   │  Cycle N+1 │  Cycle N+2 │
    ├──────────┼────────────┼────────────┼────────────┤
    │ Object A │  local     │  victim    │  DROPPED   │
    │ Object B │  -         │  local     │  victim    │
    │ Object C │  -         │  -         │  local     │
    └──────────┴────────────┴────────────┴────────────┘
```

Get() checks the victim cache after exhausting the current generation. This
prevents thrashing in services with bursty traffic: a burst creates many pooled
objects, and the victim cache keeps them alive through one quiet period before
discarding.

---

## 7. Production Patterns: Working With and Around Zeroing

### Pattern 1: Reset on Get, Not on Put

```go
// RECOMMENDED: reset when you Get
func acquire() *bytes.Buffer {
    buf := pool.Get().(*bytes.Buffer)
    buf.Reset()
    return buf
}

// NOT recommended: reset when you Put
// Why? The object might be GC'd before anyone Gets it,
// wasting the reset work. Also, the next Get might come
// from a different pool (victim cache) or New(), so you
// cannot rely on Put having cleaned it.
```

### Pattern 2: Type-Safe Pool Wrapper (Go 1.18+ Generics)

```go
type Pool[T any] struct {
    pool sync.Pool
    reset func(*T)
}

func NewPool[T any](newFn func() *T, resetFn func(*T)) *Pool[T] {
    return &Pool[T]{
        pool: sync.Pool{
            New: func() any { return newFn() },
        },
        reset: resetFn,
    }
}

func (p *Pool[T]) Get() *T {
    obj := p.pool.Get().(*T)
    p.reset(obj)  // automatic reset on every Get
    return obj
}

func (p *Pool[T]) Put(obj *T) {
    p.pool.Put(obj)
}

// Usage:
var bufPool = NewPool(
    func() *bytes.Buffer { return new(bytes.Buffer) },
    func(b *bytes.Buffer) { b.Reset() },
)
```

This wrapper guarantees the reset happens. The caller cannot forget.

### Pattern 3: Pre-allocate to Avoid Zeroing Cost

When you know the size upfront, pre-allocate to avoid repeated
allocation-and-zeroing cycles:

```go
// BAD: each append may trigger growslice → new allocation → zeroing
func collect(items []Event) []string {
    var result []string
    for _, item := range items {
        result = append(result, item.Name)
    }
    return result
}

// GOOD: one allocation, one zeroing pass, no growslice
func collect(items []Event) []string {
    result := make([]string, 0, len(items))
    for _, item := range items {
        result = append(result, item.Name)
    }
    return result
}
```

### Pattern 4: Avoid Zeroing Large Arrays on the Stack

```go
// This zeros 1MB on the stack every time the function is called
func risky() {
    var buf [1 << 20]byte  // 1MB, zeroed by duffzero/REP STOSQ
    // ...
}

// Better: allocate once, reuse via Pool or package-level var
var largeBuf = make([]byte, 1<<20)
```

### Pattern 5: Detecting Unnecessary Allocations

Use escape analysis to find where zeroing is happening unnecessarily:

```bash
go build -gcflags='-m -m' ./... 2>&1 | grep "escapes to heap"
```

Each "escapes to heap" line represents a heap allocation that goes through
`mallocgc`, including zeroing. Reducing escapes reduces zeroing overhead.

---

## 8. Quick Reference Card

```
  ┌──────────────────────────────────────────────────────────────────────┐
  │                     ZERO VALUE QUICK REFERENCE                      │
  ├──────────────────────────────────────────────────────────────────────┤
  │                                                                      │
  │  GUARANTEE                                                           │
  │    var x T        → zeroed by compiler (stack frame zeroing)         │
  │    new(T)         → zeroed by mallocgc (heap)                        │
  │    make([]T, n)   → zeroed by mallocgc                               │
  │    pool.Get()     → NOT ZEROED (caller must reset)                   │
  │                                                                      │
  │  MALLOCGC PATHS                                                      │
  │    Tiny  (<16B, no ptrs): packed into 16B blocks, pre-zeroed         │
  │    Small (16B-32KB):      mspan slot, check needzero flag            │
  │    Large (>32KB):         mheap pages, may skip if fresh from OS     │
  │                                                                      │
  │  ZEROING STRATEGIES (by size)                                        │
  │    1-8B:     single MOV $0 instruction (~1 cycle)                    │
  │    9-32B:    2-4 MOV instructions (~2-4 cycles)                      │
  │    33-512B:  DUFFZERO (jump into unrolled sequence, ~1 cycle/32B)    │
  │    513B+:    REP STOSQ (hardware accelerated, ~40 GB/s)              │
  │                                                                      │
  │  NEEDZERO FLAG                                                       │
  │    mspan.needzero == 0: memory is clean (fresh from OS or zeroed)    │
  │    mspan.needzero == 1: memory has stale data (freed objects)        │
  │    Set by: sweeper after freeing objects                              │
  │    Checked by: mallocgc before returning slot                        │
  │                                                                      │
  │  SYNC.POOL RULES                                                     │
  │    ✓ Always Reset() on Get, not on Put                               │
  │    ✓ Never store objects that need Close() or cleanup                 │
  │    ✓ Objects survive at most 2 GC cycles (victim cache)              │
  │    ✓ Use for temporary buffers, scratch space, encoders              │
  │    ✗ Never for DB connections, file handles, stateful objects        │
  │                                                                      │
  │  DEBUGGING COMMANDS                                                  │
  │    go build -gcflags='-S' ./...         # see DUFFZERO in assembly   │
  │    go build -gcflags='-m -m' ./...      # see escape + alloc reasons │
  │    go test -bench=. -benchmem ./...     # measure allocs per op      │
  │    GODEBUG=gctrace=1 ./binary           # watch GC + heap growth     │
  │                                                                      │
  └──────────────────────────────────────────────────────────────────────┘
```

---

## Further Reading

- `runtime/malloc.go` : the `mallocgc` function, trace the three allocation
  paths and where `memclr` gets called.
- `runtime/mheap.go` : search for `needzero` to see exactly when the flag
  flips.
- `runtime/duff_amd64.s` : the actual duffzero implementation. Count the
  `MOVUPS` instructions and verify the 32-bytes-per-block structure.
- `sync/pool.go` : read `poolCleanup()` to see the victim cache rotation
  during GC.
- `cmd/compile/internal/ssagen/ssa.go` : search for `OpZeroRange` and
  `duffDevice` to see the compiler's size threshold decisions.
- [Go Blog: "Go GC: Prioritizing low latency and simplicity"](https://blog.golang.org/go15gc)
  covers why zeroing costs are acceptable given the GC design.
- [Tom Duff's original Usenet post (1983)](http://www.lysator.liu.se/c/duffs-device.html)
  explaining the device that inspired Go's implementation.
