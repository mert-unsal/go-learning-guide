# Deep Dive: Go Slice Internals — Header, growslice, Backing Arrays & Memory Leaks

> Everything the runtime does when you create a slice, append to it,
> share a backing array, or accidentally leak memory through one.

---

## Table of Contents

1. [The Slice Header (`runtime.slice`)](#1-the-slice-header-runtimeslice)
2. [Arrays vs Slices — What Actually Differs](#2-arrays-vs-slices--what-actually-differs)
3. [`make()` Under the Hood](#3-make-under-the-hood)
4. [Append Internals — `runtime.growslice()`](#4-append-internals--runtimegrowslice)
5. [The Full Slice Expression: `s[low:high:max]`](#5-the-full-slice-expression-slowhighmax)
6. [Shared Backing Array — The Mutation Trap](#6-shared-backing-array--the-mutation-trap)
7. [Copy Semantics](#7-copy-semantics)
8. [Nil Slice vs Empty Slice](#8-nil-slice-vs-empty-slice)
9. [Memory Leaks via Slices](#9-memory-leaks-via-slices)
10. [Slice Tricks — Performance Analysis](#10-slice-tricks--performance-analysis)
11. [Performance Characteristics](#11-performance-characteristics)
12. [Quick Reference Card](#12-quick-reference-card)

---

## 1. The Slice Header (`runtime.slice`)

A slice in Go is a **3-word value type** — a struct containing a pointer, a length, and
a capacity. It is NOT a pointer. It is NOT a reference type. It is a **value that contains
a pointer**.

```
runtime.slice (24 bytes on 64-bit)
┌──────────────────────┬──────────────────────┬──────────────────────┐
│  array unsafe.Pointer│  len   int           │  cap   int           │
│  (pointer to first   │  (number of elements │  (total elements in  │
│   element of the     │   accessible via      │   backing array from │
│   backing array)     │   this slice header)  │   array ptr onward)  │
└──────────────────────┴──────────────────────┴──────────────────────┘
       8 bytes                 8 bytes                8 bytes
```

**Source:** `runtime/slice.go`
```go
// The compiler treats slice values as this 3-word struct:
type slice struct {
    array unsafe.Pointer
    len   int
    cap   int
}
```

### Why 3 Words?

```
  s := make([]int, 3, 5)

  Stack (slice header):                  Heap (backing array):
  ┌────────────────────┐                 ┌─────┬─────┬─────┬─────┬─────┐
  │ array ─────────────┼────────────────►│  0  │  0  │  0  │  ?  │  ?  │
  ├────────────────────┤                 └─────┴─────┴─────┴─────┴─────┘
  │ len = 3            │                 index:  0     1     2     3     4
  ├────────────────────┤                         ◄──── len ────►
  │ cap = 5            │                         ◄──────── cap ────────►
  └────────────────────┘
```

- **`array`** — points to index 0 of the backing array. This is where the data lives.
- **`len`** — how many elements you can read/write via `s[i]`. Bounds-checked at runtime.
- **`cap`** — how many elements exist in the backing array from `array` onward. Determines
  whether `append` needs to allocate a new array.

### Value Type With Pointer Semantics

A slice header is passed by value (copied on assignment and function calls), but the
copy shares the same backing array — cheap "passing" (24-byte copy) with shared data.

```go
a := []int{1, 2, 3}
b := a             // COPY of the header, SAME backing array
b[0] = 99          // mutates the shared backing array
fmt.Println(a[0])  // 99 — a sees the mutation!
```

```
  After b := a:
  a: ┌──────────┐                          ┌─────┬─────┬─────┐
     │ array ───┼────────────────────┐     │  1  │  2  │  3  │
     │ len = 3  │                    ├────►│     │     │     │
     │ cap = 3  │                    │     └─────┴─────┴─────┘
     └──────────┘                    │
  b: ┌──────────┐                    │
     │ array ───┼────────────────────┘
     │ len = 3  │
     │ cap = 3  │
     └──────────┘
```

---

## 2. Arrays vs Slices — What Actually Differs

### Arrays: Value Types Where Size Is Part of the Type

```go
var a [5]int        // type is [5]int — the 5 is part of the type
var b [3]int        // type is [3]int — DIFFERENT type from [5]int
// a = b            // compile error: cannot use [3]int as [5]int
```

An array is a **contiguous block of memory** — no header, no indirection. The
entire array IS the value. Size = `N × sizeof(T)`.

### Assignment copies ALL data

```go
a := [5]int{10, 20, 30, 40, 50}
b := a            // copies ALL 40 bytes — deep copy, completely independent
b[0] = 99
fmt.Println(a[0]) // 10 — a is unchanged
```

### How `s := arr[1:3]` Creates a Slice Header

```go
arr := [5]int{10, 20, 30, 40, 50}
s := arr[1:3]
```

```
  arr (40 bytes on stack):
  ┌─────┬─────┬─────┬─────┬─────┐
  │ 10  │ 20  │ 30  │ 40  │ 50  │
  └─────┴──▲──┴─────┴─────┴─────┘
  index: 0 │  1     2     3     4
           │
  s (slice header — 24 bytes):
  ┌────────────────┐
  │ array = &arr[1]│──► points to index 1
  │ len = 2        │    (high - low = 3 - 1)
  │ cap = 4        │    (len(arr) - low = 5 - 1)
  └────────────────┘
  s[0] == arr[1] == 20,  s[1] == arr[2] == 30
```

### Side-by-Side Comparison

```
┌──────────────────┬────────────────────┬───────────────────────┐
│ Feature          │ Array              │ Slice                 │
├──────────────────┼────────────────────┼───────────────────────┤
│ Type identity    │ [N]T (N in type)   │ []T (size not in type)│
│ Memory size      │ N × sizeof(T)      │ 24B header + backing  │
│ Assignment       │ Deep copy all data │ Shallow (header only) │
│ Func arg pass    │ Copies all data    │ Copies 24B header     │
│ Can resize?      │ No                 │ Yes (append)          │
│ Comparable (==)? │ Yes (element-wise) │ No (use slices.Equal) │
│ Can be map key?  │ Yes                │ No                    │
└──────────────────┴────────────────────┴───────────────────────┘
```

---

## 3. `make()` Under the Hood

```go
s := make([]int, 3, 5)
```

### What the Compiler Generates

The compiler translates `make([]T, len, cap)` into `runtime.makeslice()`:

**Source:** `runtime/slice.go`
```go
func makeslice(et *_type, len, cap int) unsafe.Pointer {
    mem, overflow := math.MulUintptr(et.Size_, uintptr(cap))
    if overflow || mem > maxAlloc || len < 0 || len > cap {
        panicmakeslicecap()
    }
    return mallocgc(mem, et, true)  // allocate cap × sizeof(T) bytes, zeroed
}
```

### The Allocation Split: Header on Stack, Backing Array on Heap

```
  s := make([]int, 3, 5)

  Stack:                                 Heap (mallocgc(40 bytes, zeroed)):
  ┌────────────────────┐                 ┌─────┬─────┬─────┬─────┬─────┐
  │ array = 0xc000... ─┼────────────────►│  0  │  0  │  0  │  0  │  0  │
  │ len   = 3          │                 └─────┴─────┴─────┴─────┴─────┘
  │ cap   = 5          │
  └────────────────────┘
```

`makeslice` returns `unsafe.Pointer` to the backing array. The compiler constructs the
24-byte header on the caller's stack.

### Escape Analysis: When Does the Header Escape?

The slice header stays on the stack UNLESS it escapes:

```go
func local() {
    s := make([]int, 0, 100)  // header: stack, backing array: heap
    s = append(s, 1, 2, 3)
    fmt.Println(len(s))
}

func leaks() []int {
    s := make([]int, 0, 100)  // header AND backing array: both on heap
    return s                   // slice header escapes
}
```

Verify with:
```bash
go build -gcflags='-m' ./...
# make([]int, 0, 100) escapes to heap     ← backing array
# s escapes to heap                        ← header (only when returned)
```

If the compiler can prove the backing array doesn't escape AND the size is small
(known at compile time), it MAY stack-allocate the backing array too. Check with
`-gcflags='-m'`.

---

## 4. Append Internals — `runtime.growslice()`

`append` is the most misunderstood slice operation. Let's trace every path.

```go
s = append(s, elem)
```

### Critical Insight: `append` Returns a NEW Slice Header

`append` **always** returns a new slice header. You **must** reassign:

```go
s = append(s, 42)   // ✅ reassign
append(s, 42)        // ❌ BUG — return value discarded
```

### Case A: `len < cap` — No Growth Needed

```go
s := make([]int, 2, 5)
s[0], s[1] = 10, 20
s = append(s, 30)   // len(2) < cap(5) → just write s[2]=30, len=3. Same backing array.
```

### Case B: `len == cap` — Growth Required

```go
s := []int{10, 20, 30}   // len=3, cap=3
s = append(s, 40)         // len == cap → must grow!
```

The compiler calls `runtime.growslice()`:

```
  1. append detects: len(s) == cap(s) → no room
  2. Calls: runtime.growslice(typeof(int), old_slice, newLen=4)
  3. growslice computes new capacity (see growth algorithm below)
  4. Allocates NEW backing array via mallocgc(newCap × sizeof(T))
  5. Copies old data → new array via memmove
  6. Writes new element at position old_len
  7. Returns NEW slice header: {new_array_ptr, newLen, newCap}
```

```
  BEFORE (len=3, cap=3):          AFTER growslice (len=4, cap=6):
  ┌─────┬─────┬─────┐            ┌─────┬─────┬─────┬─────┬─────┬─────┐
  │ 10  │ 20  │ 30  │  (GC'd)   │ 10  │ 20  │ 30  │ 40  │  0  │  0  │
  └─────┴─────┴─────┘            └─────┴─────┴─────┴─────┴─────┴─────┘
                                                          ▲new   ▲cap=6
```

### The Growth Algorithm (Go 1.18+)

**Source:** `runtime/slice.go` → `nextslicecap()` (changed in CL 347823)

```
  if newLen > 2×oldCap     → return newLen        (caller needs more than 2x)
  if oldCap < 256          → return 2×oldCap      (small: double)
  else                     → oldCap += (oldCap + 768) / 4  (≈1.25x + 192)
                              repeat until >= newLen
```

```
  Growth examples:  0→1  1→2  4→8  128→256  256→512  512→832  1024→~1472
  After nextslicecap, runtime rounds up to nearest memory allocator size class
  (runtime/malloc.go → roundupsize(), ~70 classes from 8B to 32KB).
```

```go
s := make([]int, 0)
for i := 0; i < 5; i++ {
    s = append(s, i)
    fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
}
// Actual caps: 1, 2, 4, 4, 8 — rounded up to size class boundaries
```

### Append Cost Analysis: Two Paths, Vastly Different Costs

A common misconception is that `append` is always expensive. It's not.
`append` has **two completely different execution paths**:

```
                    append(s, elem)
                         │
                   ┌─────┴─────┐
                   │ len < cap? │
                   └─────┬─────┘
                  yes/         \no
                  │             │
          ┌───────▼──────┐  ┌──▼──────────────────────────────────┐
          │ FAST PATH     │  │ SLOW PATH (growslice)               │
          │ ~2-5ns        │  │ ~100-500ns+                         │
          │               │  │                                     │
          │ 1. Write elem │  │ 1. Compute new capacity             │
          │    at s[len]  │  │ 2. mallocgc: allocate NEW array     │
          │ 2. len++      │  │ 3. memmove: copy all old elements   │
          │ 3. Done       │  │ 4. Write new element                │
          │               │  │ 5. Return new slice header          │
          │ No allocation │  │ 6. Old array → garbage collector    │
          │ No copy       │  │                                     │
          │ No GC work    │  │ Allocation + copy + GC pressure     │
          └───────────────┘  └─────────────────────────────────────┘
```

**Fast path (len < cap):** Nearly free. Just writes to the next slot in the
already-allocated backing array and increments `len`. No allocation, no copy,
no GC involvement. This is comparable to `s[i] = elem` in cost.

**Slow path (len == cap):** Expensive. Allocates an entirely new backing array,
copies every existing element via `memmove`, writes the new element, and leaves
the old array for GC to collect. The cost grows with slice size — copying 1M
elements is far more expensive than copying 10.

### The Practical Lesson: Pre-Allocate to Stay on the Fast Path

```go
// ❌ Starts at cap=0 — triggers growslice repeatedly
// For 1000 elements: ~10 growslice calls, ~10 dead arrays for GC
results := []int{}
for i := 0; i < 1000; i++ {
    results = append(results, i)
}

// ✅ Starts at cap=1000 — every append hits the fast path
// For 1000 elements: 0 growslice calls, 0 dead arrays, 1 allocation total
results := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    results = append(results, i)
}
```

**When you know the size** (or even a rough estimate), `make([]T, 0, n)` eliminates
the slow path entirely. Every `append` becomes a ~2ns slot write instead of a
potential allocation+copy+GC event.

**When you don't know the size**, `append` is still safe — the amortized cost is O(1)
per element thanks to the exponential growth strategy. But in hot paths processing
thousands of requests per second, those growslice allocations add up as GC pressure.

---

## 5. The Full Slice Expression: `s[low:high:max]`

The three-index slice controls the **capacity** of the resulting slice.

```go
s := []int{10, 20, 30, 40, 50}
t := s[1:3:4]    // len = high-low = 2,  cap = max-low = 3
```

```
  s:  ┌─────┬─────┬─────┬─────┬─────┐
      │ 10  │ 20  │ 30  │ 40  │ 50  │
      └─────┴──▲──┴─────┴─────┴─────┘
               t.array   len=2, cap=3
```

Without the third index, a sub-slice inherits the full remaining capacity:

```go
sub := original[1:3]             // cap=4 — append can overwrite original[3]!
sub2 := original[1:3:3]          // cap=2 — append triggers growslice → detach
```

**Rule of thumb:** Use `s[low:high:high]` to prevent accidental mutation of shared data.

---

## 6. Shared Backing Array — The Mutation Trap

This is the single most common source of slice bugs in Go. Two slices sharing a backing
array can silently corrupt each other.

### Step-by-Step Lifecycle

```go
a := []int{10, 20, 30, 40, 50}   // Step 1: create
b := a[1:3]                        // Step 2: shared backing array
b[0] = 99                          // Step 3: mutates a[1]!
b = append(b, 77)                  // Step 4: overwrites a[3] (within cap)
b = append(b, 88)                  // Step 5: overwrites a[4] (within cap)
b = append(b, 66)                  // Step 6: exceeds cap → b detaches!
```

```
  Steps 1-2: b = a[1:3], b shares a's backing array
  a: ──► ┌─────┬─────┬─────┬─────┬─────┐
         │ 10  │ 20  │ 30  │ 40  │ 50  │   b.len=2, b.cap=4
         └─────┴──▲──┴─────┴─────┴─────┘
                  b

  Step 3: b[0] = 99 → a[1] is now 99!

  Steps 4-5: append within cap → overwrites a[3]=77, a[4]=88

  Step 6: len(b)==cap(b) → growslice → NEW backing array
  a: ──► ┌─────┬─────┬─────┬─────┬─────┐
         │ 10  │ 99  │ 30  │ 77  │ 88  │  ← a keeps old array
         └─────┴─────┴─────┴─────┴─────┘
  b: ──► ┌─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┐
         │ 99  │ 30  │ 77  │ 88  │ 66  │  0  │  0  │  0  │  ← NEW array
         └─────┴─────┴─────┴─────┴─────┴─────┴─────┴─────┘
         a and b are now COMPLETELY INDEPENDENT
```

### How to Safely Detach

```go
b := make([]int, len(a[1:3]))     // Method 1: explicit copy
copy(b, a[1:3])

b := a[1:3:3]                      // Method 2: full slice expression (detach on first append)

b := append([]int(nil), a[1:3]...) // Method 3: append to nil
```

---

## 7. Copy Semantics

`copy` creates a **completely independent** copy of elements.

```go
src := []int{10, 20, 30, 40, 50}
dst := make([]int, 3)
n := copy(dst, src)   // n = 3 — copies min(len(dst), len(src)) elements
```

### Under the Hood: `runtime.memmove`

**Source:** `runtime/memmove_*.s` (assembly, platform-specific)

`copy(dst, src)` compiles to `runtime.memmove(n × sizeof(T) bytes)`. Uses SIMD
instructions (AVX2/SSE) for large copies on amd64, REP MOVSB for medium copies,
direct MOV for tiny copies. Handles overlapping regions safely.

### Copy Semantics Table

```
┌────────────────────────────┬───────────────────────────────────────┐
│ Expression                 │ What Happens                          │
├────────────────────────────┼───────────────────────────────────────┤
│ n := copy(dst, src)        │ Copies min(len(dst), len(src)) elems │
│ copy(s[2:], s[1:])         │ Overlapping copy — safe (memmove)    │
│ copy(dst, "hello")         │ String → []byte copy, no alloc       │
│ b := a (assignment)        │ NOT data copy! 24B header copy only  │
└────────────────────────────┴───────────────────────────────────────┘
```

---

## 8. Nil Slice vs Empty Slice

These look the same from the outside but differ at the memory level.

```
  Nil slice (var s []int):        {array: nil,       len: 0, cap: 0} ← 24 zero bytes
  Empty slice (s := []int{}):     {array: &zerobase, len: 0, cap: 0} ← non-nil pointer
```

`runtime.zerobase` (`runtime/malloc.go`) is a single global variable. ALL zero-length
allocations point here. It's never dereferenced — just provides a non-nil address.

### Behavioral Comparison

```
┌──────────────────────┬──────────────┬──────────────┐
│ Operation            │ nil slice    │ empty slice  │
├──────────────────────┼──────────────┼──────────────┤
│ len(s)               │ 0            │ 0            │
│ cap(s)               │ 0            │ 0            │
│ s == nil             │ true         │ false ⚠️      │
│ append(s, 1)         │ works ✅      │ works ✅      │
│ for range s {}       │ 0 iters ✅    │ 0 iters ✅    │
│ json.Marshal(s)      │ "null" ⚠️    │ "[]"         │
│ reflect.DeepEqual    │ nil ≠ []     │ [] ≠ nil     │
│ fmt.Sprintf("%#v",s) │ []int(nil)   │ []int{}      │
└──────────────────────┴──────────────┴──────────────┘
```

### The JSON Trap — Production Impact

```go
type Response struct {
    Items []string `json:"items"`
}

r1 := Response{}                 // nil:   {"items":null}  ← may break clients
r2 := Response{Items: []string{}} // empty: {"items":[]}   ← safe
```

**Rule:** For API responses, always initialize slices to empty, never leave them nil.

### Best Practice

```go
var s []int                      // nil — semantic: "no data"
s := []int{}                     // empty — semantic: "data exists, it's empty"
var results []int
results = append(results, 1)    // works — append handles nil slices
```

---

## 9. Memory Leaks via Slices

Slices are the #1 source of subtle memory leaks in Go. The GC cannot collect a backing
array if ANY slice header still points into it.

### Leak 1: Large Backing Array Retained by Small Slice

```go
func processData() []byte {
    data := make([]byte, 1<<20)  // 1 MB
    // ... fill data ...
    return data[:3]               // 3 bytes — but 1MB backing array stays alive!
}
```

**Fix:** Copy the needed data to a new, small slice:

```go
func processData() []byte {
    data := make([]byte, 1<<20)
    result := make([]byte, 3)
    copy(result, data[:3])   // new 3-byte backing array, 1MB becomes GC-eligible
    return result
}
```

### Leak 2: Slice of Pointers — Deleted Elements Still Referenced

When you shrink a slice of pointers, elements beyond the new length but within
capacity still hold references:

```go
users := []*User{{"A"}, {"B"}, {"C"}, {"D"}}
users = users[:2]   // backing array still holds pointers to C and D!
```

```
  After users = users[:2]:
  Backing array:
  ┌─────────┬─────────┬─────────┬─────────┐
  │ *User A │ *User B │ *User C │ *User D │  ← C, D outside len but GC sees them
  └─────────┴─────────┴────▲────┴────▲────┘
  len=2, cap=4              │         │
                  GC keeps C and D alive!
```

**Fix:** Nil out elements before shrinking:

```go
for i := 2; i < len(users); i++ {
    users[i] = nil   // break the reference
}
users = users[:2]
```

### Leak 3: Substrings Retaining the Parent String

Strings are immutable byte slices (`{Data uintptr, Len int}`). Substrings share
the parent's backing bytes:

```go
huge := loadHugeString()        // 10 MB string
prefix := huge[:10]              // 10 bytes — but 10 MB stays alive!

// Fix:
prefix := strings.Clone(huge[:10])  // Go 1.20+: allocates new backing bytes
```

---

## 10. Slice Tricks — Performance Analysis

All tricks below modify slices in-place using `append`, `copy`, and re-slicing.
Understanding their allocation behavior is critical for hot paths.

### Delete Element at Index i (Order-Preserving)

```go
s = append(s[:i], s[i+1:]...)   // O(n), 0 allocs, modifies original
```

This one-liner does a lot. Let's break it apart step by step.

#### What Each Part Means

```go
s[:i]       // slice from start up to (not including) index i
s[i+1:]     // slice from index i+1 to end
...         // spread operator — unpacks s[i+1:] into individual arguments for append
append(a, b...)  // append all elements of b onto a
```

#### Step-by-Step Example: Delete Index 2 from `[10, 20, 30, 40, 50]`

```go
s := []int{10, 20, 30, 40, 50}   // len=5, cap=5
i := 2                            // delete element 30

s = append(s[:2], s[3:]...)
//         ^^^^   ^^^^^
//       [10,20]  [40,50]
```

**Step 1 — `s[:2]` creates a sub-slice:**
```
s[:2] = [10, 20]    (len=2, cap=5, shares backing array with s)
```

**Step 2 — `s[3:]` creates another sub-slice:**
```
s[3:] = [40, 50]    (len=2, shares backing array with s)
```

**Step 3 — `append(s[:2], 40, 50)` writes into the backing array:**

Since `s[:2]` has `len=2, cap=5`, there's room. `append` writes `40` at index 2
and `50` at index 3 — directly in the **original backing array**:

```
BEFORE:
┌─────┬─────┬─────┬─────┬─────┐
│ 10  │ 20  │ 30  │ 40  │ 50  │
└─────┴─────┴─────┴─────┴─────┘
  [0]   [1]   [2]   [3]   [4]

AFTER append(s[:2], s[3:]...):
┌─────┬─────┬─────┬─────┬─────┐
│ 10  │ 20  │ 40  │ 50  │ 50  │
└─────┴─────┴─────┴─────┴─────┘
  [0]   [1]   [2]   [3]   [4]
                          ▲
               s has len=4, this slot is
               inaccessible but still exists
```

The returned slice has `len=4`: `[10, 20, 40, 50]`. Element `30` is gone.
Element `50` is duplicated at index 4 but inaccessible (beyond `len`).

#### Why It's O(n) and 0 Allocations

- **O(n):** `append` must shift every element after index `i` one position left.
  Internally this is a `memmove` of `(len - i - 1)` elements.
- **0 allocations:** `s[:i]` has the same backing array with enough capacity,
  so `append` uses the fast path — no `growslice`, no new array.

#### The Gotcha: Original Slice Is Modified

Because this operates on the **same backing array**, any other slice sharing that
backing array sees the mutation:

```go
original := []int{10, 20, 30, 40, 50}
alias := original[:]                    // shares backing array

original = append(original[:2], original[3:]...)
// original = [10, 20, 40, 50]
// alias    = [10, 20, 40, 50, 50]  ← corrupted! alias still has len=5
```

If other slices reference the same backing array, **copy first**:

```go
s2 := make([]int, len(s))
copy(s2, s)
s2 = append(s2[:i], s2[i+1:]...)  // safe — independent backing array
```

#### Go 1.21+: `slices.Delete`

```go
s = slices.Delete(s, i, i+1)   // same operation, clearer intent
```

Under the hood, `slices.Delete` does the same `append` trick but also **zeroes the
vacated slot** at the end to prevent memory leaks when the slice holds pointers or
structs with pointer fields.

### Delete Element at Index i (Unordered — Fast)

```go
s[i] = s[len(s)-1]   // swap with last
s = s[:len(s)-1]      // shrink — O(1), 0 allocs
```

### Insert Element at Index i

```go
s = append(s, zero)      // grow by one (0-1 allocs)
copy(s[i+1:], s[i:])     // shift right — O(n)
s[i] = value
```

### Filter In-Place (Zero Allocation)

```go
n := 0
for _, v := range s {
    if keep(v) { s[n] = v; n++ }
}
s = s[:n]   // O(n), 0 allocs
```

### Reverse

```go
for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 { s[i], s[j] = s[j], s[i] }
// Or: slices.Reverse(s) — Go 1.21+
```

### Tricks Summary Table

```
┌──────────────────────┬───────┬────────┬────────────────────────────┐
│ Operation            │ Time  │ Allocs │ Notes                      │
├──────────────────────┼───────┼────────┼────────────────────────────┤
│ Delete (ordered)     │ O(n)  │ 0      │ Preserves order, in-place  │
│ Delete (unordered)   │ O(1)  │ 0      │ Changes order, in-place    │
│ Insert at index      │ O(n)  │ 0-1    │ Preserves order            │
│ Filter in-place      │ O(n)  │ 0      │ Single pass, in-place      │
│ Reverse              │ O(n)  │ 0      │ In-place                   │
│ Pop last             │ O(1)  │ 0      │ In-place                   │
│ Push (append)        │ O(1)* │ 0-1    │ Amortized O(1)             │
└──────────────────────┴───────┴────────┴────────────────────────────┘
```

---

## 11. Performance Characteristics

### Operation Cost Table

```
┌──────────────────────┬──────────────┬────────────────────────────────┐
│ Operation            │ Cost         │ Why                            │
├──────────────────────┼──────────────┼────────────────────────────────┤
│ make([]T, n)         │ ~50-200ns    │ mallocgc + zeroing memory      │
│ s[i] (index)         │ ~1ns         │ Bounds check + load (cached)   │
│ append (no grow)     │ ~2-5ns       │ Write element + increment len  │
│ append (grow)        │ ~100-500ns+  │ mallocgc + memmove (amort O(1))│
│ copy(dst, src)       │ ~0.5ns/byte  │ memmove (SIMD on amd64)        │
│ range iteration      │ ~1ns/elem    │ Sequential = cache-friendly    │
│ s[low:high] reslice  │ ~0.5ns       │ Pointer arithmetic, no copy    │
│ len(s), cap(s)       │ ~0ns         │ Inlined: direct field read     │
└──────────────────────┴──────────────┴────────────────────────────────┘
```

### Pre-Allocation: The #1 Performance Optimization

```go
// BAD — repeated growslice calls
var results []int
for _, v := range input {
    results = append(results, transform(v))
}

// GOOD — single allocation
results := make([]int, 0, len(input))
for _, v := range input {
    results = append(results, transform(v))
}
// 10K elements: BAD = ~14 allocs + ~20K copies. GOOD = 1 alloc + 0 copies.
```

### `sync.Pool` for Hot-Path Slice Reuse

In high-throughput code, reuse slices to avoid repeated allocation:

```go
var bufPool = sync.Pool{
    New: func() any {
        b := make([]byte, 0, 4096)
        return &b
    },
}

func handleRequest() {
    bp := bufPool.Get().(*[]byte)
    buf := (*bp)[:0]          // reset length to 0, keep capacity
    buf = append(buf, data...)
    *bp = buf                  // put back (may have grown)
    bufPool.Put(bp)
}
```

**Warning:** `sync.Pool` objects can be GC'd at any time (survive at most two GC cycles).
This is NOT a connection pool — it's a best-effort recycling bin.

### Cache Locality: Why Slices Beat Linked Lists

```
  Slice: contiguous memory — CPU prefetcher works perfectly
  ┌───┬───┬───┬───┬───┬───┬───┬───┐
  │ 0 │ 1 │ 2 │ 3 │ 4 │ 5 │ 6 │ 7 │  ← one cache line (64 bytes)
  └───┴───┴───┴───┴───┴───┴───┴───┘    ~1ns per element (L1 cache hits)

  Linked list: scattered pointers — cache misses everywhere
  ┌───┬─┐     ┌───┬─┐     ┌───┬─┐
  │ 0 │─┼────►│ 1 │─┼────►│ 2 │─┼────► ...
  └───┴─┘     └───┴─┘     └───┴─┘      ~5-100ns per element
  0x1000       0x5000       0x9000
```

Iteration over a 1M-element slice is typically **10-100x faster** than a linked list
of the same size, purely due to cache behavior.

### Escape Analysis Tips

```bash
go build -gcflags='-m' ./...       # what escapes to heap
go build -gcflags='-m -m' ./...    # verbose — see WHY something escapes
```

---

## 12. Quick Reference Card

```
SLICE HEADER (24 bytes on 64-bit)   Source: runtime/slice.go
─────────────────────────────────
{array unsafe.Pointer, len int, cap int}
  └─ value type containing a pointer — copies share the backing array

MAKE:  make([]T, len, cap) → runtime.makeslice() → mallocgc(cap × sizeof(T))
APPEND: s = append(s, ...) — ALWAYS reassign! len<cap → same array; len==cap → growslice
GROWTH: 2x for cap < 256, then ~1.25x + 192 (Go 1.18+)
FULL SLICE: s[a:b:c] → len=b-a, cap=c-a — use s[a:b:b] to prevent shared-array mutation
NIL vs EMPTY: nil={nil,0,0} json→"null" | empty={&zerobase,0,0} json→"[]"
COPY:  n := copy(dst, src) → min(len(dst), len(src)) — fully independent (memmove)

MEMORY LEAKS
  1. Small reslice of large array → copy to new slice
  2. Slice of pointers after delete → nil out removed elements
  3. Substring of huge string → strings.Clone()

TOOLS
  go build -gcflags='-m'      # escape analysis
  go build -gcflags='-m -m'   # verbose escape analysis
  go build -gcflags='-S'      # assembly output
  go test -bench=. -benchmem  # benchmark + allocations
  go test -race ./...         # race detector
```

---

## One-Line Summary

> A slice is a 24-byte value type `{array, len, cap}` — `array` points to the backing
> data, `append` returns a new header (must reassign), growth allocates a new backing
> array via `runtime.growslice`, and any sub-slice shares the backing array until
> capacity is exceeded. Pre-allocate to avoid growslice, nil-out deleted pointer
> elements to avoid leaks, and use the full slice expression `s[a:b:b]` to prevent
> accidental mutation of shared data.
