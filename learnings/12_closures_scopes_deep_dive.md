# Deep Dive: Go Closures & Scopes — funcval, Capture by Reference, and the Compiler's Transformation

> Everything the compiler does when you create a closure: scope resolution,
> variable capture, heap escape, the funcval struct, and what breaks at scale.

---

## Table of Contents

1. [Block Scoping in Go — The Scope Chain](#1-block-scoping-in-go--the-scope-chain)
2. [Functions as First-Class Values — The `funcval` Struct](#2-functions-as-first-class-values--the-funcval-struct)
3. [Closures — The Compiler's Transformation](#3-closures--the-compilers-transformation)
4. [Capture by Reference, Not Value](#4-capture-by-reference-not-value)
5. [Escape Analysis and Closures](#5-escape-analysis-and-closures)
6. [Closures and Goroutines — Race Conditions](#6-closures-and-goroutines--race-conditions)
7. [The Classic Gotchas](#7-the-classic-gotchas)
8. [Method Values and Method Expressions](#8-method-values-and-method-expressions)
9. [Common Closure Patterns in Production](#9-common-closure-patterns-in-production)
10. [Performance Implications](#10-performance-implications)
11. [Quick Reference Card](#11-quick-reference-card)

---

## 1. Block Scoping in Go — The Scope Chain

Go resolves variable names at **compile time** by walking a scope chain. Every `{}` creates
a new lexical scope. The compiler never needs the runtime for name resolution — it's all
determined before the binary is built.

### The Four Scope Levels

```
┌─────────────────────────────────────────────────────────────────┐
│  UNIVERSE SCOPE  (predeclared: true, false, nil, len, cap,     │
│                   append, make, new, error, int, string, ...)  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  PACKAGE SCOPE  (package-level var, const, func, type)  │   │
│  │  ┌───────────────────────────────────────────────────┐   │   │
│  │  │  FILE SCOPE  (import declarations only)           │   │   │
│  │  │  ┌────────────────────────────────────────────┐   │   │   │
│  │  │  │  FUNCTION SCOPE  (params, named returns)  │   │   │   │
│  │  │  │  ┌─────────────────────────────────────┐   │   │   │   │
│  │  │  │  │  BLOCK SCOPE  (if, for, switch, {}) │   │   │   │   │
│  │  │  │  │  ┌──────────────────────────────┐    │   │   │   │   │
│  │  │  │  │  │  NESTED BLOCK SCOPE          │    │   │   │   │   │
│  │  │  │  │  └──────────────────────────────┘    │   │   │   │   │
│  │  │  │  └─────────────────────────────────┘   │   │   │   │   │
│  │  │  └────────────────────────────────────────┘   │   │   │   │
│  │  └───────────────────────────────────────────────┘   │   │   │
│  └──────────────────────────────────────────────────────┘   │   │
└─────────────────────────────────────────────────────────────────┘
```

**Resolution rule:** The compiler walks outward from the current scope. The first match wins.
An inner variable with the same name **shadows** the outer one — the outer variable still
exists but is unreachable by name within the inner scope.

### Multi-Level Shadowing Trace

```go
x := "package"                         // scope A

func example() {                       // scope B
    x := "function"                    // shadows A's x
    fmt.Println(x)                     // → "function"

    if true {                          // scope C
        x := "if-block"               // shadows B's x
        fmt.Println(x)                // → "if-block"

        for i := 0; i < 1; i++ {      // scope D
            x := "for-block"          // shadows C's x
            fmt.Println(x)            // → "for-block"
        }
        fmt.Println(x)                // → "if-block" (C's x, D is gone)
    }
    fmt.Println(x)                    // → "function" (B's x, C is gone)
}
```

```
Compiler's scope chain lookup for "x" inside the for-body:

   scope D: x = "for-block"  ← FOUND, stop looking
   scope C: x = "if-block"       (shadowed — unreachable from D)
   scope B: x = "function"       (shadowed)
   scope A: x = "package"        (shadowed)
```

### `:=` vs `var` — Both Create in Current Scope

Both `x := val` and `var x = val` declare a new variable in the **current** scope. The
`:=` operator is not "reassignment" — it is always a **new declaration**. This is why
`i := i` inside a loop body creates a separate variable: the right-hand `i` resolves to
the outer scope, the left-hand `i` declares in the current scope.

**Key insight:** Shadowing is resolved entirely at compile time. There is no runtime scope
chain walk. The compiler assigns each variable to a specific stack slot or heap location.

**Source:** `cmd/compile/internal/types2` — Go's type checker resolves all scopes at compile time.

---

## 2. Functions as First-Class Values — The `funcval` Struct

In Go, functions are values. They have a type (`func(params) returns`), can be assigned
to variables, passed as arguments, and returned from other functions. But under the hood,
a function value is not just a code pointer.

### The `runtime.funcval` Struct

Every function value in Go is a pointer to a `funcval` struct. The struct's first field
is always the function's code pointer.

```
runtime.funcval
┌──────────────────────────────────────┐
│  fn  uintptr   // pointer to the    │
│                // machine code       │
└──────────────────────────────────────┘
```

**Source:** `runtime/runtime2.go`
```go
type funcval struct {
    fn uintptr
    // variable-size, fn-specific data here
}
```

### Simple Function Value (No Closure)

```go
func add(a, b int) int { return a + b }

var f func(int, int) int = add
```

```
  Stack variable f:
  ┌──────────────┐
  │  f ──────────┼──► funcval{fn: 0x48a230}
  └──────────────┘             │
                               ▼
                         machine code for add()
                         at address 0x48a230
```

For a non-closure function value, the `funcval` contains only the code pointer. The
compiler often places these in read-only data — they are effectively singletons. Assigning
`add` to multiple variables all point to the same `funcval`.

### Why a Pointer to a Struct, Not Just a Code Pointer?

The calling convention requires it. When Go calls a function through a function value,
it passes the `funcval` pointer in a register (DX on amd64). For simple functions,
this pointer is unused inside the function. For closures, the function reads its captured
variables through this pointer. The uniform representation — always a `funcval*` —
means the caller doesn't need to know whether it's calling a closure or a plain function.

```
Call through function value:
  1. Load funcval pointer into DX register
  2. Load funcval.fn (code address) from [DX+0]
  3. CALL [DX+0]
  4. Inside the function: DX points to funcval — closure reads captured vars from [DX+8], [DX+16], ...
```

---

## 3. Closures — The Compiler's Transformation

A closure is a function literal that references variables from its enclosing scope. The
compiler transforms it into a `funcval` struct with captured variable pointers appended
after the code pointer.

### Source Code → Compiler Output

```go
func makeAdder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}

adder := makeAdder(5)
fmt.Println(adder(3))   // 8
```

### Step 1 — Compiler Identifies Free Variables

```
Compiler scans the function literal func(y int) int { return x + y }
  ├─ "y" → declared as parameter → local, not captured
  ├─ "x" → NOT declared in this function → FREE VARIABLE
  │         resolved to makeAdder's parameter "x"
  └─ This function literal is a CLOSURE (has free variables)
```

### Step 2 — Compiler Generates the Closure Struct

The compiler creates an anonymous struct type with the code pointer first, followed by
pointers to each captured variable:

```
// Compiler-generated (conceptual):
type closure_makeAdder_func1 struct {
    F    uintptr    // code pointer to the anonymous function
    x    *int       // pointer to captured variable x
}
```

### Step 3 — Memory Layout

```
  makeAdder(5) is called:

  x (heap-allocated because closure escapes):
  ┌─────────┐
  │  x = 5  │ ◄──────────────────────────────────────────┐
  └─────────┘                                            │
                                                         │
  Closure struct (heap-allocated — returned to caller):   │
  ┌────────────────────┬──────────────────┐              │
  │ fn = 0x48b100      │ x = 0xc00001a0b0 ┼─────────────┘
  │ (anonymous func    │ (pointer to the  │
  │  code address)     │  captured x)     │
  └────────────────────┴──────────────────┘

  Returned to caller:   adder ──► closure{fn: 0x48b100, x: ──► [5]}
```

### Step 4 — Calling the Closure

```
adder(3) execution:
  1. Load funcval pointer (closure struct) into DX
  2. Load DX.fn → 0x48b100 → CALL
  3. Inside the function body:
     ├─ y = 3             (normal parameter from stack/register)
     ├─ x = *([DX+8])     (load pointer from closure struct, dereference → 5)
     └─ return 5 + 3 = 8
```

**Source:** `cmd/compile/internal/walk` — The compiler's closure transformation happens
during the walk phase. See `cmd/compile/internal/walk/closure.go`.

---

## 4. Capture by Reference, Not Value

This is the single most important thing to understand about Go closures: **closures
capture VARIABLES, not values.** The closure struct holds a **pointer** to the variable,
not a copy of its value.

### Proof: Shared Variable

```go
x := 10
inc := func() { x++ }
get := func() int { return x }

inc()
fmt.Println(get())   // 11 — both closures share the SAME x
inc()
fmt.Println(get())   // 12
```

### Memory Layout — Two Closures, One Variable

```
  Stack/Heap:
  ┌───────────┐
  │   x = 12  │ ◄─────────────────────────┐
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

  Both closures hold a *int pointing to the SAME x.
  inc() modifies x through its pointer. get() reads x through its pointer.
```

### Why Pointers, Not Copies?

The Go spec says closures "may refer to variables defined in a surrounding function.
Those variables are then shared." The word "shared" is key — a copy would break this
contract. The compiler achieves sharing by moving captured variables to the heap (if they
escape) and replacing all references — in both the enclosing function and the closure —
with pointer dereferences through the same heap location.

### The `i := i` Fix — Creating a New Variable

```go
for i := 0; i < 3; i++ {
    i := i   // NEW variable — shadows the loop's i
    funcs[i] = func() { fmt.Println(i) }
}
```

```
  Without i := i (one shared variable):       With i := i (three separate variables):

  funcs[0].closure.i ──┐                      funcs[0].closure.i ──► [i₀ = 0]
  funcs[1].closure.i ──┼──► [ i = 3 ]         funcs[1].closure.i ──► [i₁ = 1]
  funcs[2].closure.i ──┘                       funcs[2].closure.i ──► [i₂ = 2]

  All see final value 3.                       Each sees its own snapshot.
```

### Go 1.22 Range Loop Change

Go 1.22 changed range loops so each iteration creates a **new loop variable**:

```go
// Go 1.22+: range loops create a new variable per iteration (FIXED)
for i, v := range slice {
    go func() { fmt.Println(i, v) }()   // ✅ safe — each iteration's i, v is unique
}

// C-style for loops are NOT changed — still one variable for the entire loop
for i := 0; i < len(slice); i++ {
    go func() { fmt.Println(i) }()      // ❌ STILL BROKEN — all share one i
}
```

**Source:** Go proposal [#60078](https://go.dev/blog/loopvar-preview). The compiler
inserts an implicit `i := i` at the top of each range loop iteration body starting in
Go 1.22.

---

## 5. Escape Analysis and Closures

When a closure captures a local variable, the compiler must decide: can that variable
stay on the stack, or must it move to the heap? This decision is **escape analysis**,
and closures are one of the most common triggers for heap escapes.

### The Rule

If a closure **escapes** the function where the captured variable is declared, the captured
variable MUST be heap-allocated. "Escapes" means: returned, passed to a goroutine, stored
in a struct/slice/map/package-var, or sent through a channel.

### Case 1 — Closure Escapes → Variable Escapes to Heap

```go
func makeCounter() func() int {
    count := 0                        // count ESCAPES to heap
    return func() int {               // closure ESCAPES (returned)
        count++
        return count
    }
}
```

```bash
$ go build -gcflags='-m' .
./main.go:3:2: moved to heap: count
./main.go:4:9: func literal escapes to heap
```

`count` must survive beyond `makeCounter`'s stack frame, so the compiler heap-allocates it.
The returned closure's funcval struct (also on heap) holds a pointer to that `count`.

### Case 2 — Closure Does NOT Escape → Variable Stays on Stack

```go
func process(items []int) int {
    sum := 0
    apply := func(x int) { sum += x }  // closure does NOT escape — sum stays on stack
    for _, item := range items { apply(item) }
    return sum
}
```

```bash
$ go build -gcflags='-m' .
./main.go:3:12: func literal does not escape
```

`apply` is never stored, returned, or goroutined. Everything stays on stack. **No GC pressure.**

### Case 3 — Goroutine Always Escapes

```go
func worker(tasks []Task) {
    for _, t := range tasks {
        t := t                  // new var per iteration (Go 1.22+ range: automatic)
        go func() { process(t) }()  // closure + t escape to heap
    }
}
```

**Rule of thumb:** `go func() { ... }()` always escapes. Every goroutine launch = heap allocation.

### Escape Analysis Decision Tree for Closures

```
Closure is created
├─ Returned from function?              → escapes → captured vars heap-allocated
├─ Passed to go func()?                 → escapes → captured vars heap-allocated
├─ Stored in struct/slice/map/pkg var?  → escapes → captured vars heap-allocated
├─ Sent through a channel?              → escapes → captured vars heap-allocated
└─ Only called within declaring func?   → does NOT escape → stack ✅
```

---

## 6. Closures and Goroutines — Race Conditions

Launching a goroutine with `go func() { ... }()` creates a closure that captures variables
by reference. If the enclosing function modifies those variables after the goroutine
launches, you have a **data race**.

### The Race

```go
func process() {
    data := "initial"
    go func() {
        fmt.Println(data)      // READ data — when? unpredictable
    }()
    data = "modified"          // WRITE data — race with the goroutine's read
}
```

### Race Timeline

```
  Main goroutine (G1)           Spawned goroutine (G2)
  ─────────────────────         ─────────────────────────
  data := "initial"
  go func() { ... }()
  │                             (G2 scheduled — maybe now, maybe later)
  data = "modified"  ←── WRITE
  │                             fmt.Println(data) ←── READ
  │                             │
  └─────── RACE ────────────────┘
           No synchronization between WRITE and READ.
           go vet / -race will flag this.
```

### The Fix: Pass as Argument (Copy)

```go
func process() {
    data := "initial"
    go func(d string) {        // d is a COPY of data at launch time
        fmt.Println(d)         // reads the copy — no race
    }(data)                    // data is copied into d HERE
    data = "modified"          // only modifies original — d is independent
}
```

### Detection

```bash
go test -race ./...            # ALWAYS run in CI — non-negotiable
go vet ./...                   # detects common closure-over-loop-var patterns
```

`go vet` specifically checks for `go func() { ... uses loop var ... }()` patterns and
reports: `loop variable i captured by func literal`.

---

## 7. The Classic Gotchas

### Gotcha 1 — Loop Variable Capture (C-Style For Loop)

**Still broken in Go 1.22+** — the range loop fix does NOT apply to C-style `for` loops.

```go
funcs := make([]func(), 3)
for i := 0; i < 3; i++ {
    funcs[i] = func() { fmt.Println(i) }   // all capture the SAME i
}
funcs[0]()   // 3 ❌
funcs[1]()   // 3 ❌
funcs[2]()   // 3 ❌
```

**Why:** One `i` variable exists for the entire loop. All closures hold `*int` pointing to
that single `i`. When the loop ends, `i == 3`.

**Fix:** Shadow it — `i := i` inside the body, or pass as argument.

### Gotcha 2 — Deferred Closure in a Loop

```go
for i := 0; i < 3; i++ {
    defer func() {
        fmt.Println(i)     // prints 3, 3, 3 (LIFO order, all see i=3)
    }()
}
```

Defers run when the **function** returns, not when the loop iteration ends. By that time,
`i` is already 3.

**Fix:**
```go
for i := 0; i < 3; i++ {
    defer func(n int) {
        fmt.Println(n)     // prints 2, 1, 0 (LIFO order, each sees its own n)
    }(i)
}
```

### Gotcha 3 — Goroutine Closure in a Loop

```go
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)     // race condition + wrong values
    }()
}
```

**Two bugs:** (a) all closures share one `i`, (b) race between `i++` and goroutine reads.

**Fix:**
```go
for i := 0; i < 3; i++ {
    go func(n int) {
        fmt.Println(n)     // each goroutine gets its own copy
    }(i)
}
```

### Gotcha 4 — Closure Capturing Mutated Outer Variable

```go
func makeHandlers() []http.HandlerFunc {
    handlers := make([]http.HandlerFunc, 0)
    prefix := "/api"
    for _, version := range []string{"v1", "v2", "v3"} {
        handlers = append(handlers, func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "%s/%s", prefix, version)
        })
    }
    prefix = "/internal"  // ALL handlers now print "/internal/vX" — shared by reference!
    return handlers
}
```

**Fix:** Snapshot `prefix` before creating closures: `p := prefix` and use `p` in closures.

### Go 1.22 Gotcha Summary

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

---

## 8. Method Values and Method Expressions

Method values are closures in disguise. When you write `p.String`, Go creates a closure
that binds the receiver.

### Method Value — Binds the Receiver

```go
type Person struct{ Name string }
func (p Person) String() string { return p.Name }

p := Person{Name: "Alice"}
f := p.String          // method value — f is a closure
fmt.Println(f())       // "Alice"
```

### Compiler Transformation

The compiler transforms `p.String` into a closure that captures a **copy** of the receiver:

```
// Compiler-generated (conceptual):
type methodValue_Person_String struct {
    F        uintptr   // code pointer to a wrapper function
    receiver Person    // COPY of p (value receiver → copied)
}

// The wrapper function:
func wrapper(mv *methodValue_Person_String) string {
    return mv.receiver.String()
}
```

### Memory Layout

```
  p := Person{Name: "Alice"}
  f := p.String

  f (funcval pointer):
  ┌──────────────┐
  │  f ──────────┼──► methodValue{
  └──────────────┘      fn: wrapper_code_ptr,
                        receiver: Person{Name: "Alice"}  ← COPY of p
                      }

  p.Name = "Bob"
  f()   → "Alice"   // f captured a COPY — doesn't see the change
```

**Key insight for value receivers:** The receiver is **copied** at the time the method
value is created. Later changes to `p` don't affect `f`.

### Pointer Receiver — Captures the Pointer

```go
func (p *Person) SetName(name string) { p.Name = name }

p := &Person{Name: "Alice"}
f := p.SetName                 // captures the POINTER
f("Bob")
fmt.Println(p.Name)            // "Bob" — f holds the same pointer
```

With pointer receivers, the method value captures the pointer itself — mutations through
`f` are visible through `p` (both point to the same `Person`).

### Method Expression — No Binding

```go
f := Person.String              // method expression — NOT a closure
fmt.Println(f(Person{"Carol"})) // must pass receiver explicitly — type: func(Person) string
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

## 9. Common Closure Patterns in Production

### Pattern 1 — Functional Options

```go
type Option func(*ServerConfig)

func WithAddr(addr string) Option {
    return func(c *ServerConfig) { c.Addr = addr }  // captures addr by reference
}

func NewServer(opts ...Option) *Server {
    cfg := ServerConfig{Addr: ":8080", ReadTimeout: 30 * time.Second}
    for _, opt := range opts { opt(&cfg) }  // each option closure mutates cfg
    return &Server{cfg: cfg}
}
```

Each `With*` function returns a closure capturing its argument. Short-lived, minimal GC impact.

### Pattern 2 — HTTP Middleware

```go
func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)       // closure captures 'next' from outer scope
        slog.Info("request", "method", r.Method, "duration", time.Since(start))
    })
}
// Usage: handler = logging(authCheck(myHandler))  — chain of closures
```

### Pattern 3 — Memoization

```go
func memoize(f func(int) int) func(int) int {
    cache := make(map[int]int)     // captured by the returned closure
    return func(n int) int {
        if v, ok := cache[n]; ok { return v }
        v := f(n)
        cache[n] = v
        return v
    }
}
```

**Warning:** NOT goroutine-safe. Multiple goroutines race on `cache`. Add `sync.Mutex`
for concurrent use.

### Pattern 4 — Deferred Cleanup

```go
func withTempFile(fn func(f *os.File) error) error {
    f, err := os.CreateTemp("", "work-*")
    if err != nil { return err }
    defer func() { f.Close(); os.Remove(f.Name()) }()  // closure captures f
    return fn(f)
}
```

The deferred closure captures `f`. Since it runs before `withTempFile` returns, `f`
stays on stack (no escape). Clean resource management.

### Pattern 5 — Iterator (Go 1.23+ Range-Over-Function)

```go
func Filter[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
    return func(yield func(T) bool) {
        for v := range seq {       // closure captures seq and pred
            if pred(v) && !yield(v) { return }
        }
    }
}
```

---

## 10. Performance Implications

### Closure Call vs Direct Call

```
Direct call:       CALL 0x48a230       → address known at compile time, CAN be inlined
Function value:    MOV DX, [fv_ptr]    → indirect call through funcval.fn, CANNOT be inlined
Closure call:      MOV DX, [cl_ptr]    → indirect + captured var loads from [DX+8], ...
```

### Allocation Costs

```
┌───────────────────────────────┬─────────────────┬────────────────────────────────┐
│ Scenario                      │ Heap Allocs      │ Why                            │
├───────────────────────────────┼─────────────────┼────────────────────────────────┤
│ Non-escaping closure, no      │ 0               │ Closure struct + captured vars │
│ captured vars escape          │                 │ all stay on stack              │
├───────────────────────────────┼─────────────────┼────────────────────────────────┤
│ Returned closure (escapes)    │ 1 + N           │ 1 for funcval struct,          │
│                               │                 │ N for each captured variable   │
│                               │                 │ that must move to heap         │
├───────────────────────────────┼─────────────────┼────────────────────────────────┤
│ go func() { ... }()           │ 1 + N           │ Goroutine closures always      │
│ with captured vars            │                 │ escape                         │
├───────────────────────────────┼─────────────────┼────────────────────────────────┤
│ Method value (p.Method)       │ 1               │ funcval struct with receiver   │
│                               │                 │ copy — always heap-allocated   │
├───────────────────────────────┼─────────────────┼────────────────────────────────┤
│ Method expression (T.Method)  │ 0               │ Static funcval — singleton     │
│                               │                 │ in read-only data              │
└───────────────────────────────┴─────────────────┴────────────────────────────────┘
```

### When to Avoid Closures in Hot Paths

In performance-critical inner loops (100k+ ops/sec), consider using a struct with methods:

```go
// Closure version — heap alloc, not inlineable:
func makeProcessor(threshold int) func(int) bool {
    return func(v int) bool { return v > threshold }
}

// Struct version — no closure, inlineable, no heap alloc:
type Processor struct{ Threshold int }
func (p Processor) Check(v int) bool { return v > p.Threshold }
```

### Verify with Escape Analysis

```bash
go build -gcflags='-m' ./...          # see what escapes
go build -gcflags='-m -m' ./...       # detailed reasons for each escape decision
go test -bench=. -benchmem ./...      # measure allocations per operation
```

---

## 11. Quick Reference Card

```
SCOPE CHAIN
  Universe → Package → File → Function → Block → Nested Block
  Compiler resolves innermost-first. := and var both declare in CURRENT scope.

FUNCVAL STRUCT
  Every function value = pointer to funcval{fn uintptr, ...captured vars}
  Simple function: singleton in rodata. Closure: fn + pointers to captured vars.
  Calling convention: funcval pointer in DX register (amd64).

CAPTURE RULES
  Closures capture VARIABLES (pointers), not values (copies).
  Multiple closures over same variable share it — mutations visible to all.
  Fix: i := i (shadow), or pass as function argument (copy).

GO 1.22 LOOP CHANGE
  Range loops: new variable per iteration (fixed).
  C-style for loops: ONE variable for entire loop (NOT fixed).

ESCAPE TRIGGERS
  Closure returned / goroutined / stored / sent on channel → heap escape.
  Closure used only locally → stack. Method values always heap-allocate.

GOROUTINE SAFETY
  go func() { uses x }()  → RACE if x modified after launch.
  go func(x T) { uses x }(x)  → safe (copy).

METHOD VALUES vs EXPRESSIONS
  p.Method → closure (binds receiver, heap)    type: func(...) ...
  T.Method → plain function value (no binding) type: func(T, ...) ...

TOOLS
  go build -gcflags='-m'        # escape analysis
  go build -gcflags='-S'        # assembly: funcval construction
  go test -race ./...           # race detector
  go vet ./...                  # loop variable capture warnings
```

---

## One-Line Summary

> A closure is a `funcval` struct `{code_ptr, &captured_var1, &captured_var2, ...}` —
> it captures variables by **pointer**, not by value. If the closure escapes, those
> variables move to the heap. This is why loop closures share one variable, why goroutine
> closures race, and why `i := i` fixes both.
