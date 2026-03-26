# Deep Dive: The `any` Type, Interface Boxing & Mixed-Type Collections

> How Go handles heterogeneous data, what it costs at the runtime level,
> and why understanding boxing is essential for production performance.

---

## Table of Contents

1. [What Is `any` and Where Did It Come From?](#1-what-is-any-and-where-did-it-come-from)
2. [Mixed-Type Collections: `[]any` Under the Hood](#2-mixed-type-collections-any-under-the-hood)
3. [Boxing: The Conversion from Concrete to Interface](#3-boxing-the-conversion-from-concrete-to-interface)
4. [The `convT` Family: How the Runtime Boxes Values](#4-the-convt-family-how-the-runtime-boxes-values)
5. [The `staticuint64s` Optimization: Free Boxing for Small Integers](#5-the-staticuint64s-optimization-free-boxing-for-small-integers)
6. [Direct Interface Types: Zero-Cost Boxing for Pointers](#6-direct-interface-types-zero-cost-boxing-for-pointers)
7. [Which Types Are Cheap vs Expensive to Box?](#7-which-types-are-cheap-vs-expensive-to-box)
8. [Unboxing: Type Assertions and Type Switches](#8-unboxing-type-assertions-and-type-switches)
9. [Equality and Comparison with `any`](#9-equality-and-comparison-with-any)
10. [Real-World Patterns Using `any`](#10-real-world-patterns-using-any)
11. [Generics vs `any`: The Modern Choice](#11-generics-vs-any-the-modern-choice)
12. [Performance Impact: Benchmarks and Escape Analysis](#12-performance-impact-benchmarks-and-escape-analysis)
13. [Quick Reference Card](#13-quick-reference-card)

---

## 1. What Is `any` and Where Did It Come From?

### The History

Before Go 1.18, the empty interface was written as `interface{}`. It has **zero methods**,
which means every type in Go satisfies it — because every type has at least zero methods.

```go
// Pre Go 1.18
func Print(v interface{}) { fmt.Println(v) }

// Go 1.18+: `any` is a built-in alias
// In the source code: type any = interface{}
func Print(v any) { fmt.Println(v) }
```

**`any` is not a new type.** It's a type alias defined in the `builtin` package:

```go
// From Go source: builtin/builtin.go
type any = interface{}
```

The `=` makes it a **true alias** — `any` and `interface{}` are identical at the compiler
level. No conversion needed, no runtime difference, no performance difference.

### Why Does Go Have This?

Go is **statically typed** — every variable has a known type at compile time. But
sometimes you genuinely don't know the type:
- JSON parsing of unknown structure
- `fmt.Println` accepting any arguments
- Plugin systems receiving dynamic data
- Database `Scan()` reading columns of unknown types

The empty interface is Go's **escape hatch from the type system**. It says: "I'll accept
anything, but I take responsibility for figuring out the type at runtime."

**Go Proverb #7: *"interface{} says nothing."*** — Use it only when you truly must.

---

## 2. Mixed-Type Collections: `[]any` Under the Hood

### Can You Have Mixed-Type Slices?

Yes. `[]any` holds values of different types:

```go
mixed := []any{42, "hello", true, 3.14, []int{1, 2, 3}}
```

But what does the runtime actually see? Not a nice heterogeneous list. It sees this:

```
mixed (slice header: 24 bytes)
┌──────────┬─────┬─────┐
│ ptr ─────┼─►   │ len │ cap │
└──────────┴─────┴─────┘
           │
           ▼
  Backing array of eface structs (each 16 bytes)
  ┌─────────────────────────┬─────────────────────────┐
  │ eface[0]                │ eface[1]                │
  │ _type → int descriptor  │ _type → string descriptor│
  │ data  → 42 (see §5)    │ data  → "hello" header  │
  ├─────────────────────────┼─────────────────────────┤
  │ eface[2]                │ eface[3]                │
  │ _type → bool descriptor │ _type → float64 desc    │
  │ data  → true            │ data  → 3.14            │
  ├─────────────────────────┼─────────────────────────┤
  │ eface[4]                │                         │
  │ _type → []int descriptor│                         │
  │ data  → heap ptr to     │                         │
  │         slice header    │                         │
  └─────────────────────────┴─────────────────────────┘
```

**Key observations:**
- The backing array contains `eface` structs (16 bytes each), NOT the actual values
- Each element is **boxed** — wrapped in a type descriptor + data pointer pair
- The actual data may live elsewhere (heap, static memory, or inlined — see sections below)
- You've lost all compile-time type information — the compiler can't help you

### Contrast with Other Languages

| Language | Mixed collection | How it works |
|----------|-----------------|--------------|
| **Python** | `[42, "hello", True]` | Everything is already a PyObject — always boxed |
| **Java** | `Object[]` or `List<Object>` | Autoboxing for primitives (`int` → `Integer`), reference types stored directly |
| **C#** | `object[]` | Boxing for value types, reference types stored directly |
| **Go** | `[]any` | Every element boxed into `eface{_type, data}` — honest about the cost |
| **Rust** | `Vec<Box<dyn Any>>` | Explicit boxing with trait objects — even more explicit than Go |

Go's approach: **no hidden cost, no autoboxing magic. You see what you pay for.**

---

## 3. Boxing: The Conversion from Concrete to Interface

"Boxing" is the process of wrapping a concrete typed value into an interface value.
Every time you assign a concrete value to `any` (or any interface), boxing happens.

```go
var x int = 42
var a any = x  // BOXING happens here
```

### What the Compiler Generates

When the compiler sees `var a any = x`, it emits code to:

1. **Determine the type descriptor** for `int` — this is a compile-time constant pointer
   to a static `_type` struct in the binary
2. **Copy the value** into a location that `eface.data` can point to
3. **Assemble the eface** — `{_type: &intType, data: pointerToValue}`

The critical question is: **where does the value copy live?** This is where the
runtime optimizations come in (sections 4-6).

### What Boxing Costs

```
Without boxing:  x is an int on the stack. 8 bytes. Zero GC involvement.

With boxing:     a is an eface (16 bytes) + potentially a heap copy of x.
                 The eface.data is a pointer → GC must scan it.
                 Type information is resolved at runtime for type assertions.
```

**The chain of costs:**
```
Boxing → heap allocation (maybe) → GC scanning → pointer indirection → no inlining
```

---

## 4. The `convT` Family: How the Runtime Boxes Values

When the compiler can't optimize boxing away, it calls one of the **`convT` functions**
in the runtime. These are specialized by type category for performance:

```
runtime/iface.go — the convT family
──────────────────────────────────────────────────────────────────

convT(t *_type, v unsafe.Pointer) unsafe.Pointer
  └─ General purpose: allocates heap memory, copies value. Slowest.

convTnoptr(t *_type, v unsafe.Pointer) unsafe.Pointer
  └─ For types with NO pointers (e.g., struct{x, y int}).
     Uses mallocgc with noscan=true → GC doesn't scan the allocated block.
     Faster because GC has less work.

convT16(v uint16) unsafe.Pointer
  └─ Specialized for 2-byte values. Checks staticuint64s first.

convT32(v uint32) unsafe.Pointer
  └─ Specialized for 4-byte values. Checks staticuint64s first.

convT64(v uint64) unsafe.Pointer
  └─ Specialized for 8-byte values. Checks staticuint64s first.

convTstring(v string) unsafe.Pointer
  └─ Specialized for strings (16 bytes: pointer + length).
     If empty string → returns pointer to a static zero value.
     Otherwise → allocates and copies the string header.

convTslice(v []byte) unsafe.Pointer
  └─ Specialized for slices (24 bytes: pointer + length + capacity).
     If nil/empty → returns pointer to a static zero value.
     Otherwise → allocates and copies the slice header.
```

### The Decision Tree

```
Compiler sees: var a any = someValue
                  │
                  ▼
        Is the value a pointer type?
        (*T, map, chan, func, unsafe.Pointer)
              │           │
             YES          NO
              │           │
              ▼           ▼
    Store pointer      Is it a small integer (0-255)?
    directly in        with size ≤ 8 bytes?
    eface.data              │           │
    (zero cost!)          YES          NO
                            │           │
                            ▼           ▼
                    Use staticuint64s   Call convT/convTnoptr
                    (zero allocation!)  (heap allocation)
```

---

## 5. The `staticuint64s` Optimization: Free Boxing for Small Integers

This is one of Go's cleverest runtime optimizations. In `runtime/iface.go`:

```go
// runtime/iface.go
var staticuint64s = [256]uint64{0, 1, 2, 3, ..., 255}
```

This is a **pre-allocated, static array** of 256 uint64 values (0 through 255).
It lives in the binary's data segment — never allocated, never garbage collected.

### How It Works

When you box a small value:

```go
var a any = 42
```

The compiler calls `convT64(42)`. Inside `convT64`:

```go
func convT64(v uint64) unsafe.Pointer {
    if v < 256 {
        // Point to the pre-allocated static value — NO ALLOCATION
        return unsafe.Pointer(&staticuint64s[v])
    }
    // Value too large for static cache — allocate on heap
    p := mallocgc(8, uint64Type, false)
    *(*uint64)(p) = v
    return p
}
```

### What This Means

```go
var a any = 42    // eface.data = &staticuint64s[42]   → NO allocation
var b any = 255   // eface.data = &staticuint64s[255]  → NO allocation
var c any = 256   // eface.data = heap-allocated copy   → ALLOCATES
var d any = -1    // int is signed, but bits: 0xFFFFFFFFFFFFFFFF > 255 → ALLOCATES
var e any = 0     // eface.data = &staticuint64s[0]    → NO allocation
```

### Which Types Benefit from `staticuint64s`?

Any type whose binary representation fits in a `uint64` AND has value 0-255:

| Type | Value range for zero-alloc | Example |
|------|---------------------------|---------|
| `int`, `int64`, `uint`, `uint64` | 0–255 | `var a any = 42` ✅ |
| `int32`, `uint32` | 0–255 | `var a any = int32(100)` ✅ |
| `int16`, `uint16` | 0–255 | `var a any = int16(200)` ✅ |
| `int8`, `uint8`/`byte` | 0–255 (full range for uint8) | `var a any = byte('A')` ✅ |
| `bool` | true (1), false (0) | `var a any = true` ✅ |
| `float32` | Only if bits 0-255 (NOT the float value!) | Rarely matches |
| `float64` | Only if bits 0-255 (NOT the float value!) | Rarely matches |

**Important:** For `float64`, the check is on the **binary representation**, not the
mathematical value. `float64(1.0)` has bits `0x3FF0000000000000` — way beyond 255.
So `var a any = 1.0` **DOES allocate**.

### Correction to a Common Misconception

Some sources claim that small values are "stored directly in the data pointer field"
(including my earlier explanation — I need to correct that). The truth is more nuanced:

**The `data` field is always a pointer.** For small values 0-255, it points to the
`staticuint64s` array (a static, pre-allocated memory region). The value is NOT stored
"inside" the pointer — the pointer points to a shared static location.

For pointer-shaped types (see next section), the actual pointer value IS stored in
`data` — but that's because the value itself is already a pointer.

---

## 6. Direct Interface Types: Zero-Cost Boxing for Pointers

Some types are inherently "pointer-shaped" — their value IS a pointer. For these types,
no copying or allocation is needed: the pointer value goes directly into `eface.data`.

### Pointer-Shaped Types (Direct Interface)

```go
// These types store their value directly in eface.data — ZERO allocation cost:

var a any = &User{Name: "sam"}   // *User → pointer already, just copy the pointer
var b any = myMap                 // map is *runtime.hmap under the hood
var c any = myChan                // chan is *runtime.hchan under the hood
var d any = myFunc                // func is *runtime.funcval under the hood
```

Why? Because these types are already pointers in memory. The compiler knows this
through a flag called `_type.kind & kindDirectIface`:

```go
// runtime/type.go
const kindDirectIface = 1 << 5  // bit flag: type is pointer-shaped

// When this bit is set, the value itself IS a pointer,
// so eface.data = the value itself (no wrapping needed)
```

### The Full Picture

```
                          Boxing a value into any
                          ─────────────────────────
                                    │
                      ┌─────────────┴──────────────┐
                      │                            │
              kindDirectIface?                Not pointer-shaped
              (pointer types)               (value types: int, string,
                      │                      struct, array, etc.)
                      │                            │
                      ▼                            ▼
              Store pointer value        Need to copy value somewhere
              directly in eface.data     eface.data points to copy
              ┌──────────────────┐              │
              │ ZERO allocation  │    ┌─────────┴─────────┐
              │ *T, map, chan,   │    │                    │
              │ func, unsafe.Ptr │    │              Value > size 0
              └──────────────────┘  Value == 0      and not zero?
                                   or zero-sized?          │
                                       │           ┌──────┴──────┐
                                       ▼           │             │
                                  Use static      Small          Large
                                  zero value    (0-255, ≤8B)   or non-trivial
                                       │           │             │
                                       ▼           ▼             ▼
                                  zeroVal[]   staticuint64s   mallocgc()
                                  (no alloc)  (no alloc)      (HEAP ALLOC)
```

---

## 7. Which Types Are Cheap vs Expensive to Box?

Here is the definitive cost table for boxing into `any`:

### 🟢 Free (Zero Allocation)

| Type | Why | Example |
|------|-----|---------|
| Any pointer `*T` | Pointer-shaped — direct interface | `var a any = &myStruct` |
| `map[K]V` | Internally `*hmap` — pointer-shaped | `var a any = myMap` |
| `chan T` | Internally `*hchan` — pointer-shaped | `var a any = myChan` |
| `func(...)` | Internally `*funcval` — pointer-shaped | `var a any = myFunc` |
| `unsafe.Pointer` | Pointer-shaped | `var a any = myPtr` |
| Small int 0-255 | `staticuint64s` optimization | `var a any = 42` |
| `bool` | `true`=1, `false`=0, both ≤ 255 | `var a any = true` |
| `byte` / `uint8` | Full range 0-255 covered | `var a any = byte('A')` |
| Zero value of any type | `zeroVal` static buffer | `var a any = 0` |

### 🟡 Cheap (Small Heap Allocation, No Pointer Scanning)

| Type | Size allocated | Why not free? |
|------|---------------|---------------|
| `int` (values > 255) | 8 bytes | Too large for `staticuint64s` |
| `float64` | 8 bytes | Binary representation ≠ small integer |
| `float32` | 4 bytes | Same reason |
| `int32` (values > 255) | 4 bytes | Too large for `staticuint64s` |
| `complex64` | 8 bytes | No static cache |
| `rune` (> 255) | 4 bytes | UTF-8 codepoints beyond ASCII |

These use `convTnoptr` — `mallocgc` with `noscan=true`, so GC allocates but skips scanning.

### 🔴 Expensive (Heap Allocation + GC Scanning)

| Type | Size allocated | Why expensive? |
|------|---------------|----------------|
| `string` | 16 bytes (header copy) | Contains a pointer → GC must scan |
| `[]T` (slice) | 24 bytes (header copy) | Contains a pointer → GC must scan |
| `struct` with pointer fields | struct size | Contains pointers → full GC scan |
| `[N]T` array (large) | N × sizeof(T) | Entire array copied to heap |

### Why `string` Boxing Is Particularly Interesting

A string in Go is:
```
string header (16 bytes)
┌──────────────────┬──────────────┐
│ ptr *byte        │ len int      │
│ (pointer to the  │ (byte count) │
│  underlying      │              │
│  byte array)     │              │
└──────────────────┴──────────────┘
```

When you box a string into `any`, the runtime copies this 16-byte header to the heap.
The actual character data is NOT copied — both the original string and the boxed copy
share the same underlying byte array. But the header copy itself allocates 16 bytes,
and because the header contains a pointer (`ptr *byte`), the GC must scan it.

**This is why `fmt.Sprintf` is slower than `strconv.Itoa` in hot paths:**
```go
fmt.Sprintf("%d", n)   // n gets boxed into any → allocation → GC pressure
strconv.Itoa(n)        // no interface, no boxing, can inline
```

---

## 8. Unboxing: Type Assertions and Type Switches

Getting values back OUT of `any` requires runtime type checking.

### Type Assertion

```go
var a any = 42

// Safe — comma-ok pattern
v, ok := a.(int)    // v = 42, ok = true
v, ok := a.(string) // v = "", ok = false

// Unsafe — panics on wrong type
v := a.(int)    // v = 42
v := a.(string) // PANIC: interface conversion: interface {} is int, not string
```

**Under the hood:** The runtime compares `eface._type` against the target type's
`_type` descriptor. This is a pointer comparison (fast) because type descriptors
are deduplicated — each type has exactly one `_type` in the binary.

### Type Switch

```go
func describe(a any) string {
    switch v := a.(type) {
    case int:
        return fmt.Sprintf("int: %d", v)
    case string:
        return fmt.Sprintf("string: %q", v)
    case bool:
        return fmt.Sprintf("bool: %t", v)
    case []int:
        return fmt.Sprintf("[]int with %d elements", len(v))
    case nil:
        return "nil"
    default:
        return fmt.Sprintf("unknown: %T", v)
    }
}
```

**Under the hood:** The compiler generates a series of `_type` pointer comparisons.
For switches with many cases, it may use the `_type.hash` field for a faster hash-based
lookup. Inside each case, `v` is already unboxed — no allocation cost.

### Performance: Type Switch vs Type Assertion vs Generics

```
Type assertion (single):     ~2ns — one pointer comparison
Type switch (few cases):     ~3-5ns — sequential pointer comparisons
Type switch (many cases):    ~5-10ns — hash-based lookup
Generic function:            ~0ns — type known at compile time, no runtime check
```

---

## 9. Equality and Comparison with `any`

This is a major **interview trap**:

```go
var a any = 42
var b any = 42
fmt.Println(a == b)  // true ✅ — ints are comparable

var c any = []int{1, 2}
var d any = []int{1, 2}
fmt.Println(c == d)  // PANIC! 💥 — slices are not comparable
```

### The Rule

`any` values can be compared with `==` only if the **underlying concrete type** is
comparable. The comparable types in Go:

```
✅ Comparable (== works):
   bool, int*, uint*, float*, complex*, string,
   pointer, channel, array (if element type is comparable),
   struct (if all fields are comparable), interface

❌ NOT comparable (== panics at runtime):
   slice, map, func
```

**This is a runtime check, not a compile-time check!** The compiler can't know what
concrete type is inside `any`, so it lets the code compile but panics at runtime.

### The `comparable` Constraint (Generics)

Go 1.18 introduced the `comparable` constraint to catch this at compile time:

```go
func Contains[T comparable](slice []T, target T) bool {
    for _, v := range slice {
        if v == target {
            return true
        }
    }
    return false
}

Contains([]int{1, 2, 3}, 2)          // ✅ compiles — int is comparable
Contains([][]int{{1}, {2}}, []int{1}) // ❌ compile error — []int is NOT comparable
```

---

## 10. Real-World Patterns Using `any`

### Pattern 1: JSON with Unknown Structure

```go
// When you don't know the JSON shape at compile time
var data any
json.Unmarshal([]byte(`{"name": "sam", "age": 30, "scores": [95, 87]}`), &data)

// data is now: map[string]any{
//   "name": "sam",        → string
//   "age":  float64(30),  → JSON numbers are always float64!
//   "scores": []any{float64(95), float64(87)},
// }

// Navigate with type assertions:
m := data.(map[string]any)
name := m["name"].(string)
age := m["age"].(float64)  // NOT int! JSON numbers → float64
```

**Interview trap:** JSON `Unmarshal` into `any` always produces `float64` for numbers,
`map[string]any` for objects, `[]any` for arrays. Never `int`, never `map[string]string`.

### Pattern 2: Variadic Functions

```go
// fmt.Println signature:
func Println(a ...any) (n int, err error)

// This is why you can pass anything:
fmt.Println(42, "hello", true, []int{1, 2, 3})
// Each argument gets boxed into an eface in the variadic slice
```

### Pattern 3: Database Row Scanning

```go
// database/sql — Scan takes ...any
rows.Scan(&id, &name, &age)
// Internally, Scan uses type switches to convert database values
// to the concrete types pointed to by your arguments
```

### Pattern 4: Plugin / Middleware Data

```go
// context.WithValue stores any
ctx = context.WithValue(ctx, "userID", 42)

// Retrieval requires type assertion
userID := ctx.Value("userID").(int)
```

### Anti-Patterns — When NOT to Use `any`

```go
// ❌ DON'T: Using any when you know the types
func Sum(values []any) float64 {  // Why not []float64 ?
    total := 0.0
    for _, v := range values {
        total += v.(float64)  // Unnecessary boxing + runtime type check
    }
    return total
}

// ❌ DON'T: Using any for "generic" code (use generics instead)
func Map(slice []any, f func(any) any) []any { ... }  // Pre-generics hack
// ✅ DO:
func Map[T, U any](slice []T, f func(T) U) []U { ... }  // Type-safe, zero boxing

// ❌ DON'T: Using any for dependency injection
type Service struct {
    repo any  // What methods does it have? Nobody knows!
}
// ✅ DO:
type Service struct {
    repo Repository  // Named interface — clear contract
}
```

---

## 11. Generics vs `any`: The Modern Choice

Since Go 1.18, most uses of `any` for "generic" code are **obsolete**. Here's the
decision framework:

```
Do I know the exact types at compile time?
    │
    YES → Use concrete types. Always.
    │
    NO → Is there a finite set of types?
         │
         YES → Use a named interface with methods, or a type constraint
         │
         NO → Is the set of types "anything comparable"?
              │
              YES → Use generics with `comparable` constraint
              │
              NO → Is it truly "anything at all"?
                   │
                   YES → Use `any` (interface{})
                   │
                   NO → Use a union constraint: `[T int | float64 | string]`
```

### Generics Under the Hood: GC Shape Stenciling

When you write generic code, Go doesn't use `any` internally. It uses **GC Shape
Stenciling** — the compiler generates one version per "GC shape":

```go
func Max[T int | float64 | string](a, b T) T {
    if a > b { return a }
    return b
}

// Compiler generates:
// Max_int(a, b int) int           — specialized for int
// Max_float64(a, b float64) float64 — specialized for float64
// Max_string(a, b string) string  — specialized for string
// All pointer types share ONE shape (they're all 8-byte pointers)
```

**Key insight:** Generics give you type safety AND performance (no boxing, inlineable).
`any` gives you flexibility but costs boxing, allocations, and runtime type checks.

---

## 12. Performance Impact: Benchmarks and Escape Analysis

### Boxing Cost Demonstration

```go
// Concrete — zero allocation
func sumConcrete(nums []int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// Boxed — allocation per element
func sumBoxed(nums []any) int {
    total := 0
    for _, n := range nums {
        total += n.(int)  // type assertion + unboxing
    }
    return total
}
```

Expected benchmark results:
```
BenchmarkSumConcrete-8    ~0 allocs/op    ~50ns for 1000 elements
BenchmarkSumBoxed-8       varies          ~500ns+ for 1000 elements
```

The concrete version is faster because:
1. **No type assertions** — CPU doesn't check types
2. **No pointer indirection** — values are contiguous in memory (cache-friendly)
3. **Inlineable** — compiler can optimize the loop body
4. **No GC scanning** — `[]int` has no pointers

### Escape Analysis Proof

```bash
$ go build -gcflags='-m' ./...
# var a any = x  →  "x escapes to heap"
# var a int = x  →  (no escape message — stays on stack)
```

---

## 13. Quick Reference Card

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        any / interface{} CHEAT SHEET                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  WHAT:  any = interface{} = empty interface = zero methods              │
│         Runtime struct: eface{_type *_type, data unsafe.Pointer}       │
│         Size: 16 bytes (two pointers)                                  │
│                                                                         │
│  BOXING COST:                                                          │
│    Pointer types (*T, map, chan, func)  → FREE (direct interface)      │
│    Small ints 0-255                     → FREE (staticuint64s)         │
│    bool                                 → FREE (0 or 1)               │
│    Larger scalars (int > 255, float64)  → 8B heap alloc (no GC scan)  │
│    string                               → 16B heap alloc (GC scans)   │
│    slice                                → 24B heap alloc (GC scans)   │
│    large struct                         → full struct heap alloc       │
│                                                                         │
│  TYPE RECOVERY:                                                        │
│    v, ok := a.(int)        // safe assertion                           │
│    v := a.(int)            // panics if wrong                          │
│    switch v := a.(type) {} // type switch (preferred)                  │
│                                                                         │
│  EQUALITY TRAP:                                                        │
│    any(42) == any(42)            → true  (ints comparable)             │
│    any([]int{}) == any([]int{})  → PANIC  (slices not comparable)     │
│                                                                         │
│  WHEN TO USE:                                                          │
│    ✅ JSON parsing unknown structure                                   │
│    ✅ fmt.Println variadic args                                        │
│    ✅ database/sql Scan                                                │
│    ✅ reflect-based frameworks                                         │
│    ❌ "Generic" code → use generics instead                            │
│    ❌ Dependency injection → use named interfaces                      │
│    ❌ When you know the types → use concrete types                     │
│                                                                         │
│  GENERICS vs any:                                                      │
│    Generics: type-safe, inlineable, zero boxing, compile-time errors   │
│    any: flexible, runtime type checks, boxing cost, runtime panics     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Key Takeaways

1. **`any` is `interface{}` — a type alias, nothing more.** Same runtime behavior.

2. **Boxing is the hidden cost.** Every value assigned to `any` gets wrapped in an
   `eface{_type, data}`. The cost depends on the value type.

3. **Three zero-cost boxing paths:** direct interface (pointers), staticuint64s
   (small ints), and zero values.

4. **The GC cost is the real killer.** Boxed values with pointers (strings, slices,
   structs with pointer fields) force GC scanning — this is where performance degrades
   at scale.

5. **Post-generics, the legitimate uses of `any` are narrow.** If you're reaching for
   `[]any`, ask: "Could I use `[T any]` instead?" If yes, always prefer generics.

6. **`any` is Go's honest version of Java's `Object`.** Same concept, but Go makes you
   see the boxing cost explicitly through escape analysis and benchmarks.

---

> *"interface{} says nothing."* — Go Proverbs
>
> But `any` says one thing clearly: **"I've traded compile-time safety for runtime
> flexibility."** Make sure the trade is worth it.
