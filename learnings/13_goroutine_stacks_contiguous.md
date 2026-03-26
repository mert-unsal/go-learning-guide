# Deep Dive: Goroutine Stacks — Contiguous Stacks, Growth, Shrinking & Pointer Adjustment

> How Go manages per-goroutine stacks, why contiguous stacks replaced segmented
> stacks in Go 1.4, what happens when a stack grows, and why pointers into the
> stack must be adjusted.

---

## Table of Contents

1. [Every Goroutine Gets Its Own Stack](#1-every-goroutine-gets-its-own-stack)
2. [The History: Segmented Stacks (Go 1.0–1.3)](#2-the-history-segmented-stacks-go-10-13)
3. [Contiguous Stacks (Go 1.4+) — The Current Design](#3-contiguous-stacks-go-14--the-current-design)
4. [Stack Growth — Step by Step](#4-stack-growth--step-by-step)
5. [Pointer Adjustment — The Critical Step](#5-pointer-adjustment--the-critical-step)
6. [Stack Maps — How the Runtime Finds Pointers](#6-stack-maps--how-the-runtime-finds-pointers)
7. [Stack Shrinking — During GC](#7-stack-shrinking--during-gc)
8. [Stack Growth Detection — The Function Preamble](#8-stack-growth-detection--the-function-preamble)
9. [Stacks and Channels — The Direct Copy Connection](#9-stacks-and-channels--the-direct-copy-connection)
10. [Performance Implications](#10-performance-implications)
11. [Debugging Stack Behavior](#11-debugging-stack-behavior)
12. [Quick Reference Card](#12-quick-reference-card)

---

## 1. Every Goroutine Gets Its Own Stack

When you write `go func() { ... }()`, the runtime allocates a new stack for that
goroutine. This stack is **private** — no other goroutine can access it directly.

```
Process Memory
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │ G1 Stack │  │ G2 Stack │  │ G3 Stack │  │ G4 Stack │  │
│  │  2KB     │  │  8KB     │  │  2KB     │  │  32KB    │  │
│  │ (private)│  │ (grew!)  │  │ (private)│  │ (deep    │  │
│  │          │  │          │  │          │  │  recursion│  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
│       │              │            │              │         │
│       ▼              ▼            ▼              ▼         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              SHARED HEAP (one per process)          │   │
│  │  All goroutines can read/write heap objects          │   │
│  │  → channels, maps, slices, structs that escape      │   │
│  │  → GC scans and collects this memory                │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

**Key properties:**
- **Initial size: 2KB** (since Go 1.4) — tiny compared to OS thread stacks (~1-8MB)
- **Heap-allocated**: the stack memory itself is allocated from the heap by the runtime
- **Growable**: up to a maximum (default 1GB, configurable via `runtime.SetMaxStack`)
- **Shrinkable**: GC can halve stacks that are <25% utilized
- **Private**: goroutine stacks are never shared — no locking needed for stack variables

**Why 2KB?** It's a deliberate trade-off. Small initial stacks let you create millions of
goroutines without exhausting memory. An OS thread stack of 1MB × 1M goroutines = 1TB.
A Go stack of 2KB × 1M goroutines = 2GB — feasible on modern machines.

---

## 2. The History: Segmented Stacks (Go 1.0–1.3)

Before Go 1.4, Go used **segmented stacks** — a linked list of stack segments.

```
Segmented Stack (Go 1.0-1.3):

  Segment 1 (4KB)        Segment 2 (4KB)        Segment 3 (4KB)
  ┌──────────────┐       ┌──────────────┐       ┌──────────────┐
  │  main()      │       │  doWork()    │       │  deepCall()  │
  │  frame       │──────►│  frame       │──────►│  frame       │
  │              │ link   │              │ link   │              │
  └──────────────┘       └──────────────┘       └──────────────┘
```

### The "Hot Split" Problem

Segmented stacks had a critical performance bug called the **hot split problem**:

```go
func processItem(item Item) {  // this function needs more stack than available
    // stack is almost full...
    buffer := make([]byte, 1024)  // needs new segment!
    process(buffer)
    // returns → segment is freed
}

for _, item := range items {
    processItem(item)  // EACH call: allocate segment → use → free → repeat
    // This loop creates and destroys a stack segment on EVERY iteration!
}
```

```
Iteration 1: alloc segment → call processItem → return → free segment
Iteration 2: alloc segment → call processItem → return → free segment
Iteration 3: alloc segment → call processItem → return → free segment
... 1 million iterations → 1 million alloc/free cycles!
```

If a function call happens right at the stack boundary, every call causes a
segment allocation, and every return frees it. In a tight loop, this turned
into constant allocation/deallocation thrashing — sometimes **100x slowdown**.

**This was unpredictable.** Whether you hit the hot split depended on the exact
stack usage, which varied by platform, compiler version, and inlining decisions.
A seemingly innocent refactor could trigger it.

### The Go Team's Decision

From the Go 1.4 release notes:

> *"Segmented stacks cause performance problems. We're switching to contiguous
> stacks, which are simpler and solve the hot split problem."*

The idea: instead of adding a new segment, **copy the entire stack to a bigger
contiguous block**. More expensive per growth event, but growth events are rare
(only when the stack actually needs more space), and there's no hot split.

---

## 3. Contiguous Stacks (Go 1.4+) — The Current Design

A contiguous stack is a single, unbroken block of memory. When it's full, the
runtime allocates a new block **twice the size** and copies everything.

```
Before growth:                      After growth:
┌──────────────────┐               ┌──────────────────────────────────────┐
│  main()          │               │  main()          (same data,        │
│    x := 42       │               │    x := 42        new address)      │
│    p := &x ──►x  │    COPY ►    │    p := &x ──►x  (pointer UPDATED!) │
│  doWork()        │               │  doWork()                            │
│    y := 10       │               │    y := 10                           │
│    [FULL!]       │               │    needMore()     (room to grow)     │
└──────────────────┘               │                                      │
   2KB at addr 0x1000              │                                      │
                                   │    (empty space for future calls)    │
                                   └──────────────────────────────────────┘
                                      4KB at addr 0x5000

   Old 2KB block is freed.
```

**Growth factor: always 2×**

```
2KB → 4KB → 8KB → 16KB → 32KB → 64KB → ... → up to 1GB (default max)
```

---

## 4. Stack Growth — Step by Step

**Source:** `runtime/stack.go` — function `newstack()` and `copystack()`

```
Step 1: DETECT — "Do I need more stack?"
────────────────────────────────────────
  Every function (except leaf functions and nosplit) has a PREAMBLE
  inserted by the compiler:

    if SP < stackguard0 {    // SP = current stack pointer
        runtime.morestack()  // stackguard0 = "danger zone" boundary
    }

  This check is ~1 CPU instruction (compare + conditional jump).
  In the hot path (stack is big enough), the branch is NOT taken.
  CPU branch predictor learns this quickly → near-zero overhead.

Step 2: ALLOCATE — "Get me a bigger stack"
──────────────────────────────────────────
  newsize := oldsize * 2
  newstack := stackalloc(newsize)  // allocate from stack pool or heap

  Stack pools: the runtime maintains per-size free lists of previously
  used stacks. If a 4KB stack was freed earlier, it can be reused here.
  This avoids hitting the heap allocator for common sizes.

Step 3: COPY — "Move everything to the new stack"
──────────────────────────────────────────────────
  memmove(newstack, oldstack, oldsize)

  This copies ALL frames, ALL local variables, ALL saved registers.
  It's a bulk memory copy — very fast for small stacks (2-8KB),
  measurable for large stacks (100KB+).

Step 4: ADJUST POINTERS — "Fix all references into the stack"
─────────────────────────────────────────────────────────────
  The stack has moved! Any pointer that pointed into the old stack
  now points to freed memory. The runtime MUST find and update
  every such pointer.

  delta := newstack_base - oldstack_base  // address offset

  Walk every stack frame:
    For each frame, consult the STACK MAP (compiler-generated bitmap):
      For each pointer-sized slot that the bitmap says is a pointer:
        if it points into the old stack range:
          *slot += delta   // adjust by the offset

  (See Section 5 for the full pointer adjustment algorithm)

Step 5: FREE OLD STACK
──────────────────────
  Return old stack to the stack pool (or free to heap if too large).

Step 6: RESUME
──────────────
  The goroutine continues executing on the new stack.
  It doesn't know anything happened — all its variables are intact
  at new addresses, all pointers are updated.
```

### How Often Does Growth Happen?

In practice, most goroutines grow their stack **0-3 times** during their lifetime:

```
Typical goroutine lifecycle:
  Start:      2KB  (enough for simple functions)
  First grow: 4KB  (if it calls a few functions)
  Maybe:      8KB  (if moderately deep call chain)
  Rare:       16KB+ (deep recursion or large stack frames)
```

The 2× growth factor means growth events decrease exponentially.
After 2-3 growths, the stack is big enough for the goroutine's entire lifetime.

---

## 5. Pointer Adjustment — The Critical Step

This is the most subtle part of contiguous stacks. When the stack moves,
every pointer **into** the old stack becomes dangling. The runtime must fix them all.

### What Needs Adjustment?

```
NEEDS adjustment (points into the stack):
  - Local variable pointers:  p := &x  (where x is on the stack)
  - Slice headers pointing to stack arrays
  - Pointer arguments passed between stack frames
  - Saved frame pointers and return addresses
  - defer struct pointers
  - Pointers in the goroutine's g struct (stackbase, stackguard, etc.)

Does NOT need adjustment (points to heap or static data):
  - Pointers to heap objects:  p := &User{}  (escaped to heap)
  - Pointers to global variables
  - Pointers to string literals
  - Function pointers
  - Interface data pointers (heap-allocated)
```

### The Adjustment Algorithm

```
Given: old stack at [old_lo, old_hi)
       new stack at [new_lo, new_hi)
       delta = new_lo - old_lo

For each frame on the stack:
  1. Look up the function's stack map (compiled into the binary)
  2. The stack map is a bitmap: 1 bit per pointer-sized slot
     bit = 1 → this slot contains a pointer
     bit = 0 → this slot is a scalar (int, bool, etc.) — skip it

  3. For each slot marked as pointer:
     ptr := *slot
     if old_lo ≤ ptr < old_hi:    // does it point into the OLD stack?
         *slot = ptr + delta       // YES → adjust it
     // else: points to heap/static — leave it alone

Example:
  old_lo = 0x1000,  old_hi = 0x1800  (old 2KB stack)
  new_lo = 0x5000,  new_hi = 0x6000  (new 4KB stack)
  delta  = 0x4000

  Slot at frame offset 16 contains: 0x1010  (points into old stack!)
    0x1010 is in [0x1000, 0x1800) → YES
    *slot = 0x1010 + 0x4000 = 0x5010  ← now points to correct new location

  Slot at frame offset 24 contains: 0xc000a4000  (points to heap)
    0xc000a4000 is NOT in [0x1000, 0x1800) → skip, leave as-is
```

### Visual Example

```
OLD STACK (0x1000):                    NEW STACK (0x5000):

  main() frame:                         main() frame:
  ┌─────────────────────────┐           ┌─────────────────────────┐
  │ x   = 42                │ 0x1010    │ x   = 42                │ 0x5010
  │ y   = 3.14              │ 0x1018    │ y   = 3.14              │ 0x5018
  │ p   = 0x1010 ─►x  [PTR]│ 0x1020    │ p   = 0x5010 ─►x  [ADJ]│ 0x5020
  │ str = "hello"      [HEAP]│ 0x1028    │ str = "hello"      [SAME]│ 0x5028
  ├─────────────────────────┤           ├─────────────────────────┤
  │ doWork() frame:         │           │ doWork() frame:         │
  │ buf = [0x1010]     [PTR]│ 0x1040    │ buf = [0x5010]     [ADJ]│ 0x5040
  │ n   = 100          [INT]│ 0x1048    │ n   = 100          [INT]│ 0x5048
  │ hp  = 0xc000a0  [HEAP] │ 0x1050    │ hp  = 0xc000a0  [SAME] │ 0x5050
  └─────────────────────────┘           └─────────────────────────┘

  Stack map for main():    [0, 0, 1, 0]  → slot 2 (p) is a pointer
  Stack map for doWork():  [1, 0, 1]     → slot 0 (buf) and slot 2 (hp) are pointers

  Adjustment:
    p:   0x1010 → in old range → 0x1010 + 0x4000 = 0x5010 ✅
    str: points to heap literal → NOT in old range → unchanged ✅
    buf: 0x1010 → in old range → 0x1010 + 0x4000 = 0x5010 ✅
    hp:  0xc000a0 → NOT in old range → unchanged ✅
    n:   not a pointer (stack map says 0) → skipped ✅
```

---

## 6. Stack Maps — How the Runtime Finds Pointers

The compiler generates a **stack map** for every function at every "safe point"
(places where stack growth or GC can happen). Without stack maps, the runtime
couldn't distinguish pointers from integers on the stack.

**Source:** `cmd/compile/internal/liveness/planner.go`

```
func doWork(p *int, n int, hp *User) {
    // p is a pointer, n is an int, hp is a pointer
    // The compiler generates a bitmap for this frame:
    // [1, 0, 1] = [pointer, scalar, pointer]
}
```

### Where Stack Maps Are Stored

```
Binary layout:
┌──────────────────────────────────────┐
│  .text section (compiled code)       │
│    func doWork(...)                  │
│    func main(...)                    │
│                                      │
├──────────────────────────────────────┤
│  .rodata section (read-only data)    │
│    stack maps:                       │
│      doWork: bitmap [1, 0, 1]        │
│      main:   bitmap [0, 0, 1, 0]    │
│      ...one per function per PC      │
│                                      │
│    Each function has a funcdata that │
│    maps: PC offset → stack map index │
│                                      │
│    Same maps are used by GC to scan  │
│    goroutine stacks for live pointers│
└──────────────────────────────────────┘
```

### Why Stack Maps Are Per-PC (Program Counter)

A function's live pointer set changes as it executes:

```go
func example() {
    // PC=0: no pointers live on stack → bitmap: []
    x := 42
    // PC=1: still no pointers → bitmap: [0]
    p := &x
    // PC=2: p is a pointer → bitmap: [0, 1]
    doSomething()
    // PC=3: if p is no longer used, bitmap: [0, 0] — GC can ignore it
}
```

The compiler tracks pointer liveness at each potential safe point and generates
the appropriate bitmap. This is the same infrastructure used by the garbage
collector to scan goroutine stacks.

---

## 7. Stack Shrinking — During GC

Stacks can also **shrink**. During garbage collection, the runtime checks each
goroutine's stack utilization:

```
GC Stack Shrinking:
  1. During the mark phase, the GC scans each goroutine's stack
  2. It notes the stack's total size and how much is actually used
  3. If used < 25% of capacity:
       newsize = oldsize / 2
       Copy stack to smaller buffer (same pointer adjustment process)
       Free old stack

Example:
  G1 grew its stack to 32KB during a deep call chain.
  Later, G1 is handling simple requests using only 4KB.
  GC sees: 4KB used / 32KB capacity = 12.5% < 25%
  → Shrink to 16KB. Next GC might shrink to 8KB.
```

**Why shrink?** Without shrinking, a goroutine that once needed 1MB for a deep
recursive call would hold 1MB forever, even if it now only uses 4KB. In a server
with 10K goroutines, that's 10GB wasted.

**Trade-off:** Shrinking too aggressively causes repeated grow-shrink cycles.
The 25% threshold is a balance — shrink only when the waste is significant.

---

## 8. Stack Growth Detection — The Function Preamble

Every function (except `nosplit` functions) starts with a **stack check**:

```asm
; Assembly output (go build -gcflags='-S')
; Function preamble for doWork():

TEXT doWork(SB), $128     ; this function needs 128 bytes of stack
    MOVQ  (TLS), CX        ; CX = current goroutine (g)
    CMPQ  SP, 16(CX)       ; compare SP with g.stackguard0
    JLS   morestack         ; if SP < stackguard0: need more stack!
    ; ... normal function body ...

morestack:
    CALL  runtime.morestack_noctxt(SB)  ; grow the stack
    JMP   doWork(SB)                     ; restart the function
```

```
Stack layout with guard:

┌──────────────────────────┐  ← stack top (high address)
│                          │
│  (used frames)           │
│                          │
├──────────────────────────┤  ← SP (stack pointer, grows downward)
│                          │
│  (remaining space)       │
│                          │
├──────────────────────────┤  ← stackguard0 (danger zone!)
│  guard area (small)      │     If SP drops below here → grow
└──────────────────────────┘  ← stack bottom (low address)
```

**Cost of the check:** One COMPARE + one conditional JUMP. The branch predictor
learns this is almost always "not taken" (stack is usually big enough), so the
amortized cost is essentially **zero** in the hot path.

### `//go:nosplit` — Skipping the Check

Some runtime functions use `//go:nosplit` to skip the stack check:

```go
//go:nosplit
func fastPath() {
    // NO stack growth check here
    // This function MUST NOT call any function that might grow the stack
    // Used in scheduler code, signal handlers, and other critical paths
}
```

This is dangerous — if the function uses more stack than available, the program
crashes with a stack overflow. Only the runtime team and very low-level code uses this.

---

## 9. Stacks and Channels — The Direct Copy Connection

Now you can see why the channel "direct send" path is so elegant:

```
Goroutine G1 (sender):             Goroutine G2 (receiver):
┌──────────────────────┐           ┌──────────────────────┐
│  Stack (G1's own)    │           │  Stack (G2's own)    │
│                      │           │                      │
│  val := 42           │ 0xA010   │  v := <-ch           │ 0xB020
│  ch <- val           │           │  (v is here, waiting │
│                      │           │   for a value)       │
└──────────────────────┘           └──────────────────────┘
         │                                   ▲
         │     typedmemmove(&v, &val, 8)     │
         └───────────────────────────────────┘
              Direct copy: G1's stack → G2's stack
```

**This is safe because:**
1. G2 is parked (sleeping) — its stack won't move during the copy
2. The sudog stores `elem = &v` — a direct pointer to G2's stack variable
3. The runtime holds `hchan.lock` — no concurrent modification
4. After the copy, G2 is woken and finds `v = 42` on its stack

**What if G2's stack had grown between parking and receiving?**
It can't! A parked goroutine is not running — it can't call functions, so no
stack growth check executes. The stack is frozen at the exact state it was in
when the goroutine was parked. The sudog.elem pointer remains valid.

---

## 10. Performance Implications

### Stack Growth Cost

```
┌──────────────────────────┬──────────────────────────────────────┐
│ Operation                │ Approximate Cost                      │
├──────────────────────────┼──────────────────────────────────────┤
│ Stack check (no growth)  │ ~0.5-1ns (branch predictor handles)  │
│ Stack growth (2KB→4KB)   │ ~1-5μs (alloc + copy + adjust)       │
│ Stack growth (64KB→128KB)│ ~10-50μs (larger copy + more ptrs)   │
│ Stack shrink (during GC) │ ~1-10μs (same as growth)             │
│ Stack alloc from pool    │ ~100ns (reuse freed stack)            │
│ Stack alloc from heap    │ ~500ns-1μs (fresh allocation)         │
└──────────────────────────┴──────────────────────────────────────┘
```

### When Stack Growth Hurts

**Deep recursion** is the main pain point:

```go
func fibonacci(n int) int {
    if n <= 1 { return n }
    return fibonacci(n-1) + fibonacci(n-2)
}

fibonacci(40):
  - Call depth: up to 40 frames
  - Each frame: ~64-128 bytes
  - Total stack: ~5KB
  - Growth events: 2-3 (2KB → 4KB → 8KB)
  - Each growth copies the ENTIRE stack

fibonacci(10000):
  - Call depth: up to 10000 frames
  - Total stack: ~1MB
  - Growth events: ~10 (2KB → 4KB → ... → 2MB)
  - Later growths copy hundreds of KB!
  - Pointer adjustment walks thousands of frames
```

**Solutions:**
- Convert recursion to iteration (always preferred in Go)
- Use a manual stack (slice) for tree/graph traversal
- Pre-allocate goroutine stack size (not directly possible, but reducing
  initial call depth helps)

### Why Go Prefers Iteration Over Recursion

Go does not have **tail call optimization** (TCO). In languages with TCO (Scheme,
Erlang, some Scala), `return f(x)` reuses the current frame. Go deliberately
chose NOT to implement TCO because:

1. Stack traces would lose information (frames are reused, not stacked)
2. defer semantics would be ambiguous (when does the deferred function run?)
3. The Go team values debuggability over micro-optimization

This means every recursive call in Go adds a frame → grows the stack → potential
growth event. For hot-path code, always prefer iteration.

---

## 11. Debugging Stack Behavior

### See Stack Size at Runtime

```go
import "runtime/debug"

func printStackSize() {
    var buf [1]byte
    // The address of a local variable tells you roughly where SP is
    fmt.Printf("approximate SP: %p\n", &buf)
}

// Or check goroutine count (indirect measure of total stack memory):
fmt.Println("goroutines:", runtime.NumGoroutine())
```

### Force Stack Growth (for Testing)

```go
//go:noinline
func growStack(depth int) {
    var buf [1024]byte  // 1KB per frame
    buf[0] = byte(depth)
    if depth > 0 {
        growStack(depth - 1)  // each call adds ~1KB frame
    }
    _ = buf
}

// growStack(100) → forces stack to grow to ~100KB
```

### See Stack Map in Assembly

```bash
go build -gcflags='-S' ./mypackage/ 2>&1 | grep -A5 "gclocals"
# gclocals entries show the stack map bitmaps for each function
```

### Monitor Stack Memory with `runtime.MemStats`

```go
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Printf("Stack in use: %d bytes\n", m.StackInuse)
fmt.Printf("Stack from OS: %d bytes\n", m.StackSys)
```

---

## 12. Quick Reference Card

```
GOROUTINE STACKS
────────────────
  Initial size: 2KB (since Go 1.4)
  Growth: 2× each time (2KB → 4KB → 8KB → ...)
  Max size: 1GB default (configurable via runtime.SetMaxStack — runtime/debug)
  Shrink: during GC, halved if <25% used
  Private: each goroutine owns its stack, no sharing

CONTIGUOUS STACKS (Go 1.4+)
───────────────────────────
  Replaced segmented stacks to fix the "hot split" problem
  Growth = allocate 2× buffer → copy all frames → adjust pointers → free old
  Pointer adjustment uses compiler-generated stack maps (bitmaps)

STACK GROWTH DETECTION
─────────────────────
  Every function preamble: if SP < stackguard0 → call morestack
  Cost when no growth needed: ~0.5ns (branch predictor)
  //go:nosplit skips the check (dangerous, runtime-only)

POINTER ADJUSTMENT
─────────────────
  delta = new_base - old_base
  For each frame → consult stack map bitmap:
    If slot is pointer AND points into old stack range:
      *slot += delta
  Same stack maps used by GC to scan stacks

STACK MAPS
──────────
  Generated by compiler for every function at every safe point
  Bitmap: 1 = pointer, 0 = scalar
  Stored in binary's .rodata section
  Maps are per-PC (pointer liveness changes within a function)

PERFORMANCE
──────────
  Stack check: ~0.5ns (branch predicted, nearly free)
  Stack growth: ~1-50μs depending on stack size
  Deep recursion: each growth copies more data → prefer iteration
  No tail call optimization in Go (by design — debuggability)

CHANNELS + STACKS
────────────────
  Direct send: copies value from sender's stack to receiver's stack
  Safe because parked goroutine's stack is frozen (no growth possible)
  sudog.elem = pointer directly into receiver's stack variable

TOOLS
─────
  runtime.NumGoroutine()              # goroutine count
  runtime.ReadMemStats(&m)            # m.StackInuse, m.StackSys
  go build -gcflags='-S'              # see stack check preamble in assembly
  GODEBUG=schedtrace=1000             # scheduler trace (shows goroutine counts)
  runtime/debug.SetMaxStack(bytes)    # configure max stack size
```

---

## Key Takeaways

1. **Every goroutine has its own private stack (2KB initial, growable to 1GB).**
   Stacks are cheap — this is why you can have millions of goroutines.

2. **Contiguous stacks (Go 1.4+) fixed the hot split problem** of segmented stacks.
   Trade-off: growth is more expensive (must copy), but growth events are rare.

3. **Stack growth copies everything and adjusts all pointers.** The compiler generates
   stack maps (bitmaps) so the runtime knows which slots are pointers vs scalars.

4. **The heap is shared, stacks are private.** This is why heap access needs
   synchronization (mutexes/channels/atomics) but stack variables don't.

5. **Prefer iteration over recursion in Go.** No TCO means every recursive call
   adds a frame. Deep recursion causes repeated expensive stack growths.

6. **Channel direct send works because parked goroutines have frozen stacks.**
   The sudog stores a pointer directly into the receiver's stack — safe because
   the receiver can't grow its stack while parked.

---

> *"Goroutine stacks start small and grow as needed. This is central to Go's
> concurrency model — cheap goroutines enable the CSP style."*
> — Go runtime documentation
