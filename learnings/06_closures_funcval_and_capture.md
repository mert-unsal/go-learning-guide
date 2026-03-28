# Deep Dive: Go Closures — funcval, Capture, Escape, and Production Patterns

> The funcval struct, capture-by-reference mechanics, escape analysis triggers,
> loop capture gotchas, and production closure patterns.
>
> **Related:** [Chapter 16 §2](./16_memory_gc_escape_analysis.md) covers escape analysis
> in the broader memory/GC context (5 escape rules, allocator, profiling).

---

## Table of Contents

1. [Functions as Values — The `funcval` Struct](#1-functions-as-values--the-funcval-struct)
2. [Closures — The Compiler's Transformation](#2-closures--the-compilers-transformation)
3. [Capture by Reference](#3-capture-by-reference)
4. [Escape Analysis and Closures](#4-escape-analysis-and-closures)
5. [The Loop Capture Gotcha](#5-the-loop-capture-gotcha)
6. [Method Values and Method Expressions](#6-method-values-and-method-expressions)
7. [Defer Internals](#7-defer-internals)
8. [Production Patterns](#8-production-patterns)
9. [Performance](#9-performance)
10. [Detection and Cross-Language](#10-detection-and-cross-language)
11. [Quick Reference Card](#11-quick-reference-card)

---

## 1. Functions as Values — The `funcval` Struct

Every function value in Go is a pointer to a `funcval` struct. The first field
is always the code pointer; closures append captured variable pointers after it.

**Source:** `runtime/runtime2.go`
```go
type funcval struct {
    fn uintptr
    // variable-size, fn-specific data here
}
```

```
  Stack variable f:
  ┌──────────────┐
  │  f ──────────┼──► funcval{fn: 0x48a230}
  └──────────────┘             │
                               ▼
                         machine code for add()
```

For a plain function value, the `funcval` is a singleton in read-only data.
Assigning `add` to multiple variables all point to the same `funcval`.

### Calling Convention

When Go calls through a function value, the `funcval` pointer goes in a register
(DX on amd64). Plain functions ignore it; closures read captured vars through it.
The caller never needs to know which kind it's calling.

```
Call through function value:
  1. Load funcval pointer into DX register
  2. Load funcval.fn (code address) from [DX+0]
  3. CALL [DX+0]
  4. Inside the function: closure reads captured vars from [DX+8], [DX+16], ...
```

---

## 2. Closures — The Compiler's Transformation

A closure is a function literal that references variables from its enclosing scope.
The compiler transforms it into a `funcval` struct with captured variable pointers
appended after the code pointer.

```go
func makeAdder(x int) func(int) int {
    return func(y int) int { return x + y }
}
adder := makeAdder(5)
fmt.Println(adder(3))   // 8
```

### Step 1 — Identify Free Variables

```
Compiler scans func(y int) int { return x + y }
  ├─ "y" → parameter → local
  ├─ "x" → not declared here → FREE VARIABLE → resolved to makeAdder's x
  └─ This is a CLOSURE (has free variables)
```

### Step 2 — Generate Closure Struct

```
type closure_makeAdder_func1 struct {
    F    uintptr    // code pointer
    x    *int       // pointer to captured x
}
```

### Step 3 — Memory Layout

```
  x (heap-allocated — closure escapes):
  ┌─────────┐
  │  x = 5  │ ◄──────────────────────────────────────────┐
  └─────────┘                                            │
                                                         │
  Closure struct (heap-allocated — returned to caller):   │
  ┌────────────────────┬──────────────────┐              │
  │ fn = 0x48b100      │ x = 0xc00001a0b0 ┼─────────────┘
  │ (code address)     │ (ptr to x)       │
  └────────────────────┴──────────────────┘

  Returned:   adder ──► closure{fn: 0x48b100, x: ──► [5]}
```

### Step 4 — Calling the Closure

```
adder(3):
  1. Load funcval pointer (closure struct) into DX
  2. Load DX.fn → CALL
  3. Inside:  y = 3 (parameter),  x = *([DX+8]) → 5,  return 8
```

**Source:** `cmd/compile/internal/walk/closure.go`

---

## 3. Capture by Reference

Closures capture **variables** (pointers to them), not values (copies).
The closure struct holds a `*T` to the captured variable.

```go
x := 10
inc := func() { x++ }
get := func() int { return x }

inc()
fmt.Println(get())   // 11 — both closures share the SAME x
```

### Memory Layout — Two Closures, One Variable

```
  ┌───────────┐
  │   x = 11  │ ◄─────────────────────────┐
  └───────────┘                            │
       ▲                                   │
       │                                   │
  ┌────┴──────────────────┐    ┌───────────┴──────────────┐
  │ closure "inc"         │    │ closure "get"            │
  │ ┌────────┬──────────┐ │    │ ┌────────┬────────────┐  │
  │ │ fn     │ x ───────┼─┼────│─│ fn     │ x ─────────┼──┘
  │ │ inc's  │ (ptr to  │ │    │ │ get's  │ (ptr to    │
  │ │ code   │  same x) │ │    │ │ code   │  same x)   │
  │ └────────┴──────────┘ │    │ └────────┴────────────┘  │
  └───────────────────────┘    └──────────────────────────┘
```

The Go spec says closures "may refer to variables defined in a surrounding function.
Those variables are then shared." The compiler achieves sharing by heap-allocating
captured variables (when they escape) and replacing all references with pointer
dereferences through the same location.

### The `i := i` Fix — Separate Variables Per Iteration

```go
for i := 0; i < 3; i++ {
    i := i   // NEW variable — shadows the loop's i
    funcs[i] = func() { fmt.Println(i) }
}
```

```
  Without i := i (shared):                   With i := i (separate):

  funcs[0].closure.i ──┐                     funcs[0].closure.i ──► [i₀ = 0]
  funcs[1].closure.i ──┼──► [ i = 3 ]        funcs[1].closure.i ──► [i₁ = 1]
  funcs[2].closure.i ──┘                      funcs[2].closure.i ──► [i₂ = 2]

  All see final value 3.                      Each sees its own snapshot.
```

### Go 1.22 Range Loop Change

```
┌───────────────────────────────────┬──────────────┬────────────────────────┐
│ Loop Type                         │ Go ≤ 1.21    │ Go 1.22+               │
├───────────────────────────────────┼──────────────┼────────────────────────┤
│ for i, v := range ...            │ ONE var      │ NEW var per iteration  │
│                                   │ (gotcha!)    │ (fixed ✅)             │
├───────────────────────────────────┼──────────────┼────────────────────────┤
│ for i := 0; i < n; i++           │ ONE var      │ ONE var                │
│                                   │ (gotcha!)    │ (STILL a gotcha ❌)    │
├───────────────────────────────────┼──────────────┼────────────────────────┤
│ Variables declared OUTSIDE loop   │ shared       │ shared                 │
│ but captured by closures in loop  │ (by design)  │ (by design — be aware) │
└───────────────────────────────────┴──────────────┴────────────────────────┘
```

The compiler inserts an implicit `i := i` at the top of each range loop iteration
body. C-style `for` loops remain unchanged — mutation-based semantics where the
developer explicitly owns the variable.

**Source:** Go proposal [#60078](https://go.dev/blog/loopvar-preview)

---

## 4. Escape Analysis and Closures

When a closure captures a variable, the compiler decides: stack or heap? If the
closure **escapes** its declaring function, captured variables must be heap-allocated.

### Escape Analysis Decision Tree

```
Closure is created
├─ Returned from function?              → escapes → captured vars heap-allocated
├─ Passed to go func()?                 → escapes → captured vars heap-allocated
├─ Stored in struct/slice/map/pkg var?  → escapes → captured vars heap-allocated
├─ Sent through a channel?              → escapes → captured vars heap-allocated
└─ Only called within declaring func?   → does NOT escape → stack ✅
```

### Cases

```go
// Case 1 — Escapes (returned): count heap-allocated
func makeCounter() func() int {
    count := 0
    return func() int { count++; return count }
}
// go build -gcflags='-m' → "moved to heap: count"

// Case 2 — Stays on stack (local use only): no GC pressure
func process(items []int) int {
    sum := 0
    apply := func(x int) { sum += x }
    for _, item := range items { apply(item) }
    return sum
}

// Case 3 — Goroutine always escapes
go func() { process(t) }()   // every goroutine launch = heap allocation
```

---

## 5. The Loop Capture Gotcha

### Gotcha 1 — Slice of Closures (C-style loop)

C-style `for` loops create **one variable** for the entire loop. All closures
share that single address.

```go
funcs := make([]func(), 3)
for i := 0; i < 3; i++ {
    funcs[i] = func() { fmt.Println(i) }   // all capture the SAME i
}
funcs[0]()   // 3
funcs[1]()   // 3
funcs[2]()   // 3
```

```
  ┌─────────────────────────────────────────────────────────────────┐
  │                    THE BUG IN ONE PICTURE                       │
  │                                                                 │
  │   ONE variable i on the heap:    [  3  ]  ← addr 0xC000012080  │
  │                                   ▲ ▲ ▲                        │
  │                                   │ │ │                         │
  │   funcs[0].&i ────────────────────┘ │ │                         │
  │   funcs[1].&i ──────────────────────┘ │                         │
  │   funcs[2].&i ────────────────────────┘                         │
  │                                                                 │
  │   3 closures × 1 variable = all see the final value             │
  └─────────────────────────────────────────────────────────────────┘
```

**Fix:** `i := i` inside the body, or pass as parameter.

### Gotcha 2 — Goroutine in a Loop

```go
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }()   // race + wrong values
}
```

Two problems: (a) shared variable → wrong values, (b) `i++` races with reads.

**Fix:**
```go
for i := 0; i < 3; i++ {
    go func(n int) { fmt.Println(n) }(i)   // copy into parameter
}
```

### Gotcha 3 — Defer in a Loop

```go
for i := 0; i < 3; i++ {
    defer func() { fmt.Println(i) }()   // prints 3, 3, 3 (LIFO)
}
```

Defers run when the **function** returns, not when the iteration ends.

**Fix:**
```go
for i := 0; i < 3; i++ {
    defer func(n int) { fmt.Println(n) }(i)   // 2, 1, 0 (LIFO, own values)
}
```

---

## 6. Method Values and Method Expressions

Method values are closures — `p.String` creates a closure binding the receiver.

```go
type Person struct{ Name string }
func (p Person) String() string { return p.Name }

p := Person{Name: "Alice"}
f := p.String          // closure capturing a COPY of p (value receiver)
fmt.Println(f())       // "Alice"
p.Name = "Bob"
fmt.Println(f())       // "Alice" — f has its own copy

// Pointer receivers capture the pointer — mutations visible through both:
func (p *Person) SetName(name string) { p.Name = name }
g := p.SetName; g("Carol"); fmt.Println(p.Name) // "Carol"
```

### Method Expression — No Binding

```go
f := Person.String              // type: func(Person) string — pass receiver explicitly
```

```
┌─────────────────────┬───────────────┬────────────────────────────────┐
│                     │ Method Value  │ Method Expression              │
├─────────────────────┼───────────────┼────────────────────────────────┤
│ Syntax              │ p.String      │ Person.String                  │
│ Resulting type      │ func() string │ func(Person) string            │
│ Is a closure?       │ YES           │ NO                             │
│ Receiver bound?     │ YES (copied   │ NO (must pass explicitly)      │
│                     │  or pointed)  │                                │
│ Heap allocation?    │ YES (funcval) │ NO (static funcval, singleton) │
└─────────────────────┴───────────────┴────────────────────────────────┘
```

---

## 7. Defer Internals

`defer` captures a function and its arguments for later execution. The
implementation has evolved for performance:

```
Go Version    Implementation           Cost per defer
──────────    ──────────────────────   ───────────────
Go 1.0–1.12   Heap-allocated record    ~50ns (malloc + free)
Go 1.13       Stack-allocated record   ~35ns (no heap)
Go 1.14+      Open-coded defer         ~0ns  (compiler inlines it)
```

### Open-Coded Defer (Go 1.14+)

The compiler inlines the deferred call at every return point — zero overhead.
Falls back to record-based when: (1) defer inside a loop, (2) >8 defers,
(3) closure capturing variables.

### Defer in Loops — Extract to Helper

```go
// ❌ SLOW — N defer records, all files open simultaneously:
for _, path := range paths {
    f, _ := os.Open(path)
    defer f.Close()              // runs at FUNCTION return, not iteration end
}

// ✅ FAST — one open-coded defer per call, file closed each iteration:
for _, path := range paths {
    if err := processOne(path); err != nil { return err }
}
func processOne(path string) error {
    f, err := os.Open(path)
    if err != nil { return err }
    defer f.Close()              // open-coded, runs at processOne return
    return process(f)
}
```

### Argument Evaluation: Defer-Time vs Call-Time

```go
x := 10
defer fmt.Println(x)                // evaluates x NOW → prints 10
defer func() { fmt.Println(x) }()   // closure captures &x → prints 20
x = 20
```

---

## 8. Production Patterns

### Pattern 1 — Functional Options

```go
type Option func(*ServerConfig)

func WithAddr(addr string) Option {
    return func(c *ServerConfig) { c.Addr = addr }
}

func NewServer(opts ...Option) *Server {
    cfg := ServerConfig{Addr: ":8080", ReadTimeout: 30 * time.Second}
    for _, opt := range opts { opt(&cfg) }
    return &Server{cfg: cfg}
}
```

### Pattern 2 — HTTP Middleware

```go
func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)       // closure captures 'next'
        slog.Info("request", "method", r.Method, "duration", time.Since(start))
    })
}
```

### Pattern 3 — Memoization

```go
func memoize(f func(int) int) func(int) int {
    cache := make(map[int]int)
    return func(n int) int {
        if v, ok := cache[n]; ok { return v }
        v := f(n)
        cache[n] = v
        return v
    }
}
// Add sync.Mutex for concurrent use — cache is shared mutable state.
```

### Pattern 4 — Deferred Cleanup

```go
func withTempFile(fn func(f *os.File) error) error {
    f, err := os.CreateTemp("", "work-*")
    if err != nil { return err }
    defer func() { f.Close(); os.Remove(f.Name()) }()
    return fn(f)
}
```

### Pattern 5 — Iterator (Go 1.23+ Range-Over-Function)

```go
func Filter[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
    return func(yield func(T) bool) {
        for v := range seq {
            if pred(v) && !yield(v) { return }
        }
    }
}
```

### Pattern 6 — Fan-In (Closure + Goroutine Application)

The fan-in pattern merges N channels into one. Each goroutine is launched with
a closure that captures `out` and `wg`, while the loop variable `src` is passed
as a parameter to avoid the capture gotcha:

```go
for _, src := range sources {
    go func(s <-chan int) {         // s = copy of src (capture fix)
        defer wg.Done()             // closure captures wg
        for v := range s { out <- v } // closure captures out
    }(src)
}
```

> **Full pattern with timeline & variations:** See [Chapter 15 §6](./15_channels_hchan_select.md)
> (channels production patterns) for the complete fan-in implementation,
> coordinator goroutine, and context cancellation variant.

---

## 9. Performance

### Call Overhead

```
Direct call:       CALL 0x48a230       → compile-time address, CAN be inlined
Function value:    MOV DX, [fv_ptr]    → indirect call, CANNOT be inlined
Closure call:      MOV DX, [cl_ptr]    → indirect + captured var loads from [DX+8]
```

### Allocation Costs

```
┌───────────────────────────────┬─────────────────┬────────────────────────────────┐
│ Scenario                      │ Heap Allocs      │ Why                            │
├───────────────────────────────┼─────────────────┼────────────────────────────────┤
│ Non-escaping closure          │ 0               │ Everything stays on stack      │
├───────────────────────────────┼─────────────────┼────────────────────────────────┤
│ Returned closure (escapes)    │ 1 + N           │ 1 funcval + N captured vars    │
├───────────────────────────────┼─────────────────┼────────────────────────────────┤
│ go func() { ... }()           │ 1 + N           │ Goroutine closures always      │
│ with captured vars            │                 │ escape                         │
├───────────────────────────────┼─────────────────┼────────────────────────────────┤
│ Method value (p.Method)       │ 1               │ funcval + receiver copy        │
├───────────────────────────────┼─────────────────┼────────────────────────────────┤
│ Method expression (T.Method)  │ 0               │ Static singleton in rodata     │
└───────────────────────────────┴─────────────────┴────────────────────────────────┘
```

### Hot Path Alternative — Struct with Methods

```go
// Closure: heap alloc, not inlineable
func makeProcessor(threshold int) func(int) bool {
    return func(v int) bool { return v > threshold }
}

// Struct: no closure, inlineable, no heap alloc
type Processor struct{ Threshold int }
func (p Processor) Check(v int) bool { return v > p.Threshold }
```

```bash
go build -gcflags='-m' ./...          # what escapes
go test -bench=. -benchmem ./...      # allocations per operation
```

---

## 10. Detection and Cross-Language

### Detection Tools

```bash
go vet ./...                          # loop variable capture warnings
go test -race ./...                   # concurrent access detection
go build -gcflags='-m' ./...          # escape analysis (heap moves)
staticcheck ./...                     # SA6005 loop capture warning
golangci-lint run                     # G601 loop variable capture (gosec)
```

### How Other Languages Handle Loop Capture

```
┌──────────────┬───────────────────────────────────────────────────────────┐
│ Language     │ Behavior                                                  │
├──────────────┼───────────────────────────────────────────────────────────┤
│ JavaScript   │ var: shared (broken). let/const: per-iteration (fixed).  │
│ (pre-ES6)    │ Same bug! for(var i=0;...) { setTimeout(()=>log(i)) }    │
│ (ES6+)       │ for(let i=0;...) creates new i per iteration. Fixed.     │
├──────────────┼───────────────────────────────────────────────────────────┤
│ Python       │ Closures capture by name, evaluated at call time.         │
│              │ Same bug: [lambda: i for i in range(5)] → all return 4.  │
│              │ Fix: [lambda i=i: i for i in range(5)] (default arg).    │
├──────────────┼───────────────────────────────────────────────────────────┤
│ C#           │ Pre-5.0: foreach shared the variable (broken).            │
│              │ C# 5.0+: foreach creates new variable per iteration.      │
│              │ for(int i=...) still shares. Exact same split as Go!      │
├──────────────┼───────────────────────────────────────────────────────────┤
│ Rust         │ Move semantics by default. Closures move or borrow.       │
│              │ The borrow checker prevents this bug at compile time.     │
│              │ You literally cannot have this bug in safe Rust.          │
├──────────────┼───────────────────────────────────────────────────────────┤
│ Java         │ Lambda can only capture "effectively final" variables.    │
│              │ Loop variable is mutated → compiler rejects capture.      │
│              │ Bug is impossible, but also less flexible.                 │
├──────────────┼───────────────────────────────────────────────────────────┤
│ Go ≤1.21     │ Both range and C-style: ONE variable. Both broken.        │
│ Go 1.22+     │ Range: per-iteration (fixed). C-style: still ONE (broken)│
└──────────────┴───────────────────────────────────────────────────────────┘

Key insight: Go 1.22's range fix mirrors C# 5.0's foreach fix — both changed
the "iteration variable" semantic while leaving C-style loops unchanged.
```

---

## 11. Quick Reference Card

```
FUNCVAL STRUCT
  Every function value = pointer to funcval{fn uintptr, ...captured vars}
  Simple function: singleton in rodata. Closure: fn + pointers to captured vars.
  Calling convention: funcval pointer in DX register (amd64).

CAPTURE RULES
  Closures capture VARIABLES (pointers to them), not values (copies).
  Multiple closures over same variable share it — mutations visible to all.
  Fix: i := i (shadow), or pass as function argument (copy).

GO 1.22 LOOP CHANGE
  Range loops: new variable per iteration (fixed ✅).
  C-style for loops: ONE variable for entire loop (still a gotcha ❌).

ESCAPE TRIGGERS
  Closure returned / goroutined / stored / sent on channel → heap escape.
  Closure used only locally → stack. Method values always heap-allocate.

GOROUTINE SAFETY
  go func() { uses x }()        → RACE if x modified after launch.
  go func(x T) { uses x }(x)   → safe (copy).

METHOD VALUES vs EXPRESSIONS
  p.Method → closure (binds receiver, heap)    type: func(...) ...
  T.Method → plain function value (no binding) type: func(T, ...) ...

FAN-IN PATTERN
  N source channels → 1 output channel.
  wg.Add(N) BEFORE goroutine loop.
  Each worker: range source → forward to out → defer wg.Done().
  Coordinator goroutine: wg.Wait() → close(out).
  Add select { case out<-v: case <-ctx.Done(): } to prevent goroutine leaks.

DETECTION TOOLS
  go vet ./...                    # loop variable capture warnings
  go test -race ./...             # concurrent access detection
  go build -gcflags='-m' ./...    # escape analysis
  staticcheck ./...               # SA6005 loop capture
```
