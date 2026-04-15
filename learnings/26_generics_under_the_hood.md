# 26 — Generics Under the Hood: GC Shape Stenciling, Dictionaries & When NOT to Use Them

> **Companion exercises:** [exercises/advanced/01_generics](../exercises/advanced/01_generics/)

---

## Table of Contents

1. [Why Go Waited Until 1.18](#1-why-go-waited-until-118)
2. [Type Parameters & Constraints](#2-type-parameters--constraints)
3. [Implementation: GC Shape Stenciling](#3-implementation-gc-shape-stenciling)
4. [The Dictionary Mechanism](#4-the-dictionary-mechanism)
5. [Constraint Interfaces in Depth](#5-constraint-interfaces-in-depth)
6. [Type Inference](#6-type-inference)
7. [Performance: Generics vs Interfaces vs Concrete](#7-performance-generics-vs-interfaces-vs-concrete)
8. [Generic Data Structures](#8-generic-data-structures)
9. [When NOT to Use Generics](#9-when-not-to-use-generics)
10. [Patterns That Shine](#10-patterns-that-shine)
11. [Comparison with Java, C#, Rust](#11-comparison-with-java-c-rust)
12. [Cost Table](#12-cost-table)
13. [Quick Reference Card](#13-quick-reference-card)
14. [Further Reading](#14-further-reading)

---

## 1. Why Go Waited Until 1.18

Go was released in 2009 without generics. This was deliberate.

Rob Pike and the Go team believed that **premature abstraction is worse than
duplication**. The Java/C# experience showed that generics, once available,
get overused: every factory, every builder, every service layer becomes
parameterized, often needlessly.

The Go team waited until they had a design that:
- Preserved **fast compilation** (no C++ template explosion)
- Kept the **language simple** (no Rust-style trait bounds complexity)
- Maintained **runtime performance** without boxing (unlike Java's type erasure)
- Was **backward compatible** with all existing Go code

The result: the **type parameters proposal** by Ian Lance Taylor and Robert
Griesemer (2021), accepted after 10+ years of experimentation with at least
six different designs.

**Key design document:** [Type Parameters Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)

---

## 2. Type Parameters & Constraints

### The Syntax

```go
func Min[T cmp.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

`[T cmp.Ordered]` means: "T is a type parameter constrained to types that
support ordering operators (`<`, `>`, `<=`, `>=`)."

### What Constraints Really Are

A constraint is just an **interface**. But since Go 1.18, interfaces can
contain more than method sets:

```go
// From the cmp package (Go 1.21+)
type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64 |
    ~string
}
```

The `~int` syntax means "any type whose **underlying type** is int" including
defined types like `type Score int`.

### The Three Built-in Constraints

| Constraint | What It Allows | Defined In |
|-----------|---------------|-----------|
| `any` | Everything (alias for `interface{}`) | builtin |
| `comparable` | Types supporting `==` and `!=` | builtin |
| `cmp.Ordered` | Types supporting `<` `>` `<=` `>=` | `cmp` package |

**Why `comparable` is special:** It's the only constraint that allows the `==`
operator. It includes all basic types, structs of comparable fields, arrays,
pointers, channels, and interfaces. It does NOT include slices, maps, or
functions.

---

## 3. Implementation: GC Shape Stenciling

This is where Go differs fundamentally from Java, C#, and Rust.

### The Problem

When the compiler sees:
```go
func Min[T cmp.Ordered](a, b T) T { ... }

x := Min(3, 5)          // T = int
y := Min("go", "rust")  // T = string
```

Should it generate **two separate functions** (like C++ templates and Rust
monomorphization)? Or **one function with runtime dispatch** (like Java type
erasure)?

### Go's Answer: GC Shape Stenciling (Hybrid)

Go generates **one version per GC shape**, where a GC shape is defined by
how the garbage collector sees the type:

| GC Shape | Types | One Compiled Function For |
|----------|-------|--------------------------|
| Pointer shape | `*int`, `*string`, `*MyStruct`, `any` | All pointer-sized types share one version |
| int shape | `int` | Its own version |
| int64 shape | `int64`, `time.Duration` | Shared (same underlying) |
| string shape | `string` | Its own version |
| `[3]int` shape | `[3]int` | Its own version |

**The key insight:** All pointer types (including interface types) have the
same GC shape (one machine word pointing to heap). So `Min[*int]` and
`Min[*string]` share the same compiled function body.

Value types each get their own version because their sizes and layouts differ.

### What Happens at Compile Time

```
Source: Min[int](3, 5)    Min[string]("go", "rust")
         │                        │
         ▼                        ▼
    Min·shape·int            Min·shape·string
    (dedicated code)         (dedicated code)

Source: Min[*int](p1, p2)  Min[*string](p3, p4)
         │                        │
         ▼                        ▼
         └──── Min·shape·pointer ─┘
               (shared code + dictionary)
```

### Viewing the Instantiations

```bash
# See which functions get generated
go build -gcflags='-m' ./... 2>&1 | grep 'instantiat'

# See assembly for a specific generic function
go build -gcflags='-S' ./... 2>&1 | grep 'Min\['
```

### Trade-off Compared to Alternatives

| Approach | Binary Size | Speed | GC Integration |
|----------|-------------|-------|----------------|
| Full monomorphization (Rust/C++) | Large | Fastest | N/A (no GC) |
| Type erasure (Java) | Small | Slower (boxing) | N/A (reference types only) |
| GC Shape Stenciling (Go) | Medium | Fast for values, slightly slower for pointers | Integrated |

Go's approach avoids Java's boxing overhead for value types while avoiding
Rust's binary bloat. The tradeoff: pointer-type generic calls go through a
dictionary (slight indirection).

---

## 4. The Dictionary Mechanism

When multiple types share a GC shape, the compiled function needs a way to
perform type-specific operations. This is done via a **runtime dictionary**.

### What the Dictionary Contains

The dictionary is a struct passed as a hidden first argument to shape-shared
generic functions:

```
┌─────────────────────────────────────┐
│          Generic Dictionary          │
├─────────────────────────────────────┤
│ Type descriptor for T               │ ← needed for reflect, type assertions
│ Itab pointer (if T satisfies iface) │ ← for interface method calls
│ Sub-dictionaries (nested generics)  │
│ Derived type pointers               │ ← for make([]T, n), new(T)
└─────────────────────────────────────┘
```

### When the Dictionary Is Used

```go
func Print[T fmt.Stringer](v T) {
    fmt.Println(v.String())  // needs dictionary to find .String() method
}
```

For pointer-shaped instantiations, the dictionary lookup adds one extra
indirection compared to a monomorphized version. For value-typed
instantiations, the function is fully specialized and the dictionary is
either inlined or not needed.

### Performance Implication

In microbenchmarks, the dictionary overhead is ~1-3ns per call for
pointer-type generics. In practice, this is negligible compared to actual
work. But in extremely hot inner loops processing millions of pointer-typed
elements, you might measure it.

---

## 5. Constraint Interfaces in Depth

### Union Elements

```go
type Number interface {
    int | float64 | complex128
}
```

This constrains T to EXACTLY these types. Named types based on them
(like `type Celsius float64`) are excluded.

### Approximation Elements (~)

```go
type Number interface {
    ~int | ~float64
}
```

The `~` tilde means "any type whose underlying type is X". Now `type Celsius
float64` satisfies this constraint.

### Method + Union Constraints

```go
type StringerOrdered interface {
    fmt.Stringer      // must have String() string method
    ~int | ~string    // must also be one of these underlying types
}
```

Both conditions must be true. This is an intersection, not a union.

### Interface Types That Cannot Be Used as Constraints

```go
type NotAConstraint interface {
    int | string
    Read([]byte) (int, error)
}
```

This is a valid interface definition, but it can only be used as a
**constraint** (in `[T NotAConstraint]`), not as a regular interface type
(you cannot declare `var x NotAConstraint`). Interfaces with type elements
are constraint-only.

### The `comparable` Subtlety

`comparable` includes types that support `==` at compile time. But there's a
nuance: `interface{}` is comparable (you can `==` two interface values), so
`any` satisfies `comparable`. But comparing interface values can panic at
runtime if the dynamic types are not comparable (e.g., slices in interfaces).

Go 1.20 resolved this: `comparable` now includes all types that are
**spec-comparable** (can be compared without panicking).

---

## 6. Type Inference

Go's type inference is intentionally limited compared to Rust or Haskell.

### What Gets Inferred

```go
// Argument types → T is inferred
Min(3, 5)             // T = int (from argument types)
Contains(s, "go")     // T = string (from slice element type)

// Return types are NOT inferred
var x int = Min(3, 5) // T inferred from args, not from x
```

### When Inference Fails

```go
// Cannot infer from return type alone
result := Map(ints, strconv.Itoa)  // ✅ T=int, U=string (from func signature)
result := Err[int](err)            // ❌ Must specify: T has no argument to infer from
```

### Constraint Type Inference

```go
type Slice[E any] interface{ ~[]E }

func Contains[S Slice[E], E comparable](s S, v E) bool { ... }

Contains([]int{1,2,3}, 2)
// S = []int (from first arg)
// E = int (inferred from S matching ~[]E)
```

The compiler can infer `E` from the constraint relationship. This is
**constraint type inference** and it's what powers the `slices` and `maps`
packages.

---

## 7. Performance: Generics vs Interfaces vs Concrete

### The Three Approaches

```go
// 1. Concrete (no generics)
func SumInts(s []int) int { ... }

// 2. Interface dispatch
func Sum(s []interface{ Add(int) int }) int { ... }

// 3. Generic
func Sum[T constraints.Integer](s []T) T { ... }
```

### Benchmark Reality

```
BenchmarkSumConcrete-8    500000000    2.34 ns/op    0 allocs/op
BenchmarkSumGenericInt-8  500000000    2.34 ns/op    0 allocs/op  ← same!
BenchmarkSumInterface-8   100000000   12.80 ns/op    1 allocs/op  ← boxing
```

For **value types**, generic code compiles to the same machine code as
hand-written concrete code. The GC shape stenciling gives each value type
its own version.

For **pointer types**, there's a small dictionary overhead, but it's
typically lost in the noise of actual work (memory access, computation).

### Escape Analysis Impact

```bash
go build -gcflags='-m' ./...
```

Generic functions with value-type instantiations behave identically to
concrete functions for escape analysis. Interface-based approaches cause
values to escape to heap (boxing).

---

## 8. Generic Data Structures

### Why Generics Shine Here

Before Go 1.18, a generic container required either:
- Code generation (`go generate`)
- `interface{}` with type assertions (losing type safety)
- Copy-paste for each type

Now:

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T)        { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() (T, bool)  { ... }
```

### The Receiver Syntax

Generic type methods repeat the type parameter on the receiver:

```go
func (s *Stack[T]) Push(v T) { ... }
//          ^^^
// This is the type parameter from the struct definition
// You CANNOT add new type parameters to methods
```

**Limitation:** Methods on generic types cannot introduce new type parameters.
This is intentional: it prevents the complexity explosion seen in languages
like Scala.

```go
// ❌ NOT ALLOWED — no extra type params on methods
func (s *Stack[T]) Map[U any](fn func(T) U) *Stack[U] { ... }

// ✅ Use a free function instead
func MapStack[T, U any](s *Stack[T], fn func(T) U) *Stack[U] { ... }
```

### Standard Library Generic Types (Go 1.21+)

The `slices`, `maps`, and `cmp` packages provide generic utilities:

```go
slices.Sort(s)                    // sort any ordered slice
slices.Contains(s, v)             // search
slices.Index(s, v)                // find index
maps.Keys(m)                      // extract keys
maps.Values(m)                    // extract values
cmp.Or(a, b, c)                   // first non-zero value
```

These use constraint type inference to accept named slice/map types
(not just `[]T`).

---

## 9. When NOT to Use Generics

The Go team's guidance (from the official blog and talks):

### Don't Use Generics When:

**1. Interface types work fine**
```go
// ❌ Unnecessary generic
func Print[T fmt.Stringer](v T) { fmt.Println(v.String()) }

// ✅ Just use the interface
func Print(v fmt.Stringer) { fmt.Println(v.String()) }
```

If the function body only calls methods from an interface, use the interface
directly. Generics add nothing here except complexity.

**2. The type parameter is only used once**
```go
// ❌ T is only used in one place
func Wrap[T any](v T) interface{} { return v }

// ✅ Just use interface{}
func Wrap(v interface{}) interface{} { return v }
```

**3. Implementation doesn't change based on type**
```go
// ❌ Over-abstraction
type Repository[T any] interface {
    FindByID(id string) (T, error)
    Save(entity T) error
}
// In Go, the "T" rarely helps here. Each repo has type-specific queries.
// Generic repos are a pattern from Java/C# that doesn't translate well.
```

**4. You're building a framework**
Go is not Java. Frameworks with deep generic hierarchies fight against
Go's philosophy. Prefer simple, explicit code over clever abstractions.

### Do Use Generics When:

1. **Collection operations** (Map, Filter, Reduce, Contains, Sort)
2. **Data structures** (Stack, Queue, Set, Tree, Cache)
3. **Type-safe wrappers** (Result[T], Optional[T], Future[T])
4. **Algorithm libraries** (binary search, merge, partition)

The rule of thumb from Ian Lance Taylor:

> *"If you find yourself writing the exact same code multiple times, where
> the only difference is the type, that's a sign you should use generics."*

---

## 10. Patterns That Shine

### Pattern 1: Functional Slice Operations

```go
func Map[T, U any](s []T, fn func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = fn(v)
    }
    return result
}
```

Pre-allocate with `make([]U, len(s))` not `append` to avoid growslice overhead.

### Pattern 2: Type-Safe Set

```go
type Set[T comparable] struct {
    m map[T]struct{}
}
```

`struct{}` takes zero bytes. Before generics, you'd need `map[string]struct{}`
or `map[int]struct{}` with no way to abstract over the key type.

### Pattern 3: Result Type

```go
type Result[T any] struct {
    value T
    err   error
    ok    bool
}
```

Go's error handling doesn't need monadic Results (the `if err != nil` pattern
is idiomatic). But for library APIs that transform chains of operations,
Result[T] can clean up error propagation.

### Pattern 4: Constraint for Numeric Operations

```go
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64
}

func Sum[T Number](s []T) T {
    var total T
    for _, v := range s {
        total += v
    }
    return total
}
```

### Pattern 5: Generic Middleware / Decorator

```go
func WithRetry[T any](fn func() (T, error), maxRetries int) (T, error) {
    var result T
    var err error
    for i := 0; i <= maxRetries; i++ {
        result, err = fn()
        if err == nil {
            return result, nil
        }
    }
    return result, err
}
```

---

## 11. Comparison with Java, C#, Rust

### Java: Type Erasure

```java
List<String> list = new ArrayList<>();
// At runtime: just a List of Object. Type info erased.
// Cannot do: new T(), T.class, instanceof T
```

Java erases generic type information at compile time. The JVM sees only
`Object`. This means:
- **No specialization** for primitives (must box `int` → `Integer`)
- **No runtime type info** (`T.class` is illegal)
- **Cannot create generic arrays** (`new T[10]` is illegal)
- Backward compatible with pre-generics code

### C#: Reification

```csharp
List<int> ints = new List<int>();  // truly specialized for int
List<string> strs = new List<string>();  // separate type at runtime
```

C# preserves full type information at runtime. The CLR JIT-compiles
specialized versions:
- **Value type specialization** (no boxing for `List<int>`)
- **Runtime type info** (`typeof(T)` works)
- Larger memory footprint (each closed type is distinct)

### Rust: Full Monomorphization

```rust
fn min<T: Ord>(a: T, b: T) -> T { if a < b { a } else { b } }
// Compiler generates min_i32, min_String, min_MyStruct, ...
```

Rust generates completely separate functions for every type:
- **Maximum performance** (identical to hand-written code)
- **Binary bloat** (many copies of the same logic)
- **Long compile times** (each instantiation must be compiled)
- **No runtime overhead** (no dictionary, no dispatch)

### Go: GC Shape Stenciling (the Middle Path)

| Feature | Java | C# | Rust | Go |
|---------|------|-----|------|-----|
| Value type specialization | ❌ (boxing) | ✅ | ✅ | ✅ |
| Pointer type sharing | N/A | ❌ | ❌ | ✅ |
| Runtime type info | ❌ (erased) | ✅ | ❌ (no runtime) | Partial (dictionary) |
| Binary size impact | None | Medium | Large | Small-Medium |
| Compile time impact | None | Small | Large | Small |
| Performance | Slow (boxing) | Fast | Fastest | Fast |

Go's approach is uniquely suited to a garbage-collected language: by sharing
code across GC-equivalent types (all pointers look the same to the GC), it
minimizes binary size while maintaining good performance.

---

## 12. Cost Table

| Operation | Cost | Notes |
|-----------|------|-------|
| Generic call (value type) | Same as concrete | Fully specialized |
| Generic call (pointer type) | +1 dictionary lookup | ~1-3ns per call |
| Constraint check | Compile time only | Zero runtime cost |
| Type inference | Compile time only | Zero runtime cost |
| Generic struct instantiation | Per GC shape | Shared where possible |
| Interface → generic migration | Better for value types | Eliminates boxing/allocation |

---

## 13. Quick Reference Card

```text
┌─────────────────────────────────────────────────────────────────┐
│                      GENERICS CHEAT SHEET                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  func Name[T constraint](args) result                           │
│  type Name[T constraint] struct { ... }                         │
│                                                                 │
│  Built-in constraints:                                          │
│    any          = interface{} (anything)                         │
│    comparable   = supports == and !=                             │
│    cmp.Ordered  = supports < > <= >=                             │
│                                                                 │
│  Custom constraint:                                             │
│    type Numeric interface { ~int | ~float64 }                   │
│    ~T means "underlying type is T"                              │
│                                                                 │
│  Method + type constraint:                                      │
│    type Constraint interface { ~int; String() string }          │
│                                                                 │
│  Methods cannot add type params:                                │
│    ❌ func (s *Stack[T]) Map[U any](...) ...                    │
│    ✅ func MapStack[T, U any](s *Stack[T], ...) ...             │
│                                                                 │
│  Type inference:                                                │
│    Min(3, 5)        ← T=int inferred from args                  │
│    Err[int](err)    ← must specify, no arg to infer from        │
│                                                                 │
│  Inspect:                                                       │
│    go build -gcflags='-m' (escape + instantiation)              │
│    go build -gcflags='-S' (assembly output)                     │
│                                                                 │
│  Use generics for: collections, data structures, algorithms     │
│  Don't use for: DI, repos, single-use wrappers                 │
│                                                                 │
│  stdlib generic packages: slices, maps, cmp (Go 1.21+)         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 14. Further Reading

- [Type Parameters Proposal (design doc)](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md) — the authoritative design document
- [Tutorial: Getting Started with Generics](https://go.dev/doc/tutorial/generics) — official tutorial
- [When to Use Generics](https://go.dev/blog/when-generics) — Ian Lance Taylor's blog post
- [GC Shape Stenciling Design](https://github.com/golang/go/issues/47514) — implementation discussion
- [`slices` package source](https://cs.opensource.google/go/go/+/master:src/slices/) — excellent examples
- [`maps` package source](https://cs.opensource.google/go/go/+/master:src/maps/) — how constraints compose
- [`cmp` package source](https://cs.opensource.google/go/go/+/master:src/cmp/) — `Ordered` and `Or`
- [Generics implementation: dictionaries and stenciling](https://github.com/golang/go/issues/48822) — compiler internals

---

> **Next:** [27 — reflect Under the Hood](27_reflect_under_the_hood.md)
>
> **Companion exercises:** [exercises/advanced/01_generics](../exercises/advanced/01_generics/)
