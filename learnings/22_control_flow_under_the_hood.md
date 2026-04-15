# Deep Dive: Control Flow Under the Hood

> How Go compiles `defer`, `range`, `switch`, and `for` — from the
> convenient syntax you write down to the runtime machinery and assembly
> the CPU actually executes.
>
> **Builds on:** [Chapter 5](./05_closures_funcval_and_capture.md) covers
> closure capture semantics, which interact with `defer`. [Chapter 11](./11_goroutine_stacks_growth.md)
> covers goroutine stacks, where defer records live.

---

## Table of Contents

1. [defer: More Than "Run at the End"](#1-defer-more-than-run-at-the-end)
2. [defer Internals: Three Implementations](#2-defer-internals-three-implementations)
3. [Named Returns + defer: The Modification Trick](#3-named-returns--defer-the-modification-trick)
4. [range: The Compiler Rewrite](#4-range-the-compiler-rewrite)
5. [range Gotchas: Value Copies and Loop Variables](#5-range-gotchas-value-copies-and-loop-variables)
6. [range Over int and func (Go 1.22/1.23)](#6-range-over-int-and-func-go-122123)
7. [switch: Jump Tables and Type Dispatch](#7-switch-jump-tables-and-type-dispatch)
8. [for: The Only Loop, Labeled Breaks, and Go 1.22](#8-for-the-only-loop-labeled-breaks-and-go-122)
9. [Quick Reference Card](#9-quick-reference-card)

---

## 1. defer: More Than "Run at the End"

`defer` schedules a function call to run when the surrounding function returns.
Most languages have try/finally or destructors for this. Go chose `defer` because
it keeps the cleanup code next to the acquisition code — a readability win.

But the mechanics have three rules that trip up even experienced engineers:

### Rule 1: Arguments Are Evaluated Immediately

```go
func example() {
    x := 10
    defer fmt.Println(x)  // x is evaluated NOW → captures 10
    x = 20
    // prints 10, not 20
}
```

The compiler evaluates `x` at the `defer` statement and stores the value `10`
in the defer record. When the function returns and the deferred call executes,
it uses the saved value. This is NOT a closure capturing a variable — it is a
function call with pre-evaluated arguments.

Compare with a closure-based defer:

```go
func example() {
    x := 10
    defer func() { fmt.Println(x) }()  // closure captures &x
    x = 20
    // prints 20 — the closure reads x at execution time
}
```

**Why this matters:** The argument-evaluation rule makes `defer f.Close()` safe
even if `f` is reassigned later. But it means `defer log.Printf("took %v", time.Since(start))`
computes the duration at defer time, not at return time. Use a closure if you need
the value at return time.

### Rule 2: Defers Execute LIFO (Stack Order)

```go
func cleanup() {
    defer fmt.Println("first")   // pushed first → runs last
    defer fmt.Println("second")  // pushed second → runs second
    defer fmt.Println("third")   // pushed last → runs first
}
// Output: third, second, first
```

Under the hood, each goroutine's `runtime.g` struct has a `_defer` linked list.
Each `defer` statement prepends a new `_defer` record to the head of this list.
At function return, the runtime walks the list from head to tail — LIFO.

```
  g._defer → [third] → [second] → [first] → nil

  On return: pop [third], execute
             pop [second], execute
             pop [first], execute
```

This stack ordering is deliberate: resources acquired later are released first.
Open file → open database connection → defer conn.Close() → defer file.Close()
means the connection closes before the file. This matches the natural dependency
order.

### Rule 3: Defers Run After Return Value Is Set, Before Caller Gets It

This is the subtlest rule and connects to named returns (section 3).

```
  Function execution timeline:
  ┌──────────────────────────────────────────────────────┐
  │  1. Function body executes                           │
  │  2. Return value is computed and assigned             │
  │  3. Deferred functions run (can read/modify returns)  │
  │  4. Function actually returns to caller               │
  └──────────────────────────────────────────────────────┘
```

---

## 2. defer Internals: Three Implementations

**Source:** `runtime/panic.go`, `cmd/compile/internal/ssagen/ssa.go`

Go has evolved through three defer implementations, each faster than the last:

### Era 1: Heap-Allocated Defer (Before Go 1.13)

Every `defer` statement allocated a `_defer` struct on the heap:

```go
type _defer struct {
    siz     int32    // size of args
    started bool     // whether the defer has started executing
    sp      uintptr  // stack pointer at time of defer
    pc      uintptr  // program counter (return address)
    fn      *funcval // the deferred function
    _panic  *_panic  // panic that is running the defer
    link    *_defer  // next defer in the chain (linked list)
    // ... followed by argument values copied inline
}
```

**Cost:** ~35ns per defer (heap allocation + GC pressure).

This was expensive enough that the Go community developed a culture of avoiding
`defer` in hot paths — a workaround for a language implementation issue, not a
language design issue.

### Era 2: Stack-Allocated Defer (Go 1.13)

The compiler learned to allocate `_defer` records on the stack when it could prove
they would not outlive the function. This covered ~90% of real-world defers.

**Cost:** ~6ns per defer (no heap allocation, no GC pressure).

### Era 3: Open-Coded Defer (Go 1.14+, Current)

**Source:** `cmd/compile/internal/ssagen/ssa.go`, search for `openDeferRecord`

For functions with ≤8 defers and no defers inside loops, the compiler eliminates
the `_defer` struct entirely. Instead, it:

1. Reserves a bitmask on the stack (1 bit per defer)
2. At each `defer` statement, sets the corresponding bit
3. At function return, checks bits and calls functions inline

```
  Traditional defer:                Open-coded defer:
  ┌──────────────────────┐         ┌──────────────────────────┐
  │ defer f() →           │         │ defer f() →               │
  │   alloc _defer struct │         │   set bit 0 in deferBits  │
  │   prepend to g._defer │         │   save args on stack      │
  │                       │         │                           │
  │ on return:            │         │ on return:                │
  │   walk _defer list    │         │   if deferBits & 1: f()   │
  │   call each fn        │         │   if deferBits & 2: g()   │
  └──────────────────────┘         └──────────────────────────┘
```

**Cost:** ~1ns per defer (just a bitwise OR).

This is why modern Go has virtually zero overhead for `defer`. The old advice to
"avoid defer in hot loops" is outdated for Go 1.14+, **unless** the defer is
inside the loop body (which prevents open-coded optimization).

### Why Defers Inside Loops Are Still Expensive

```go
// Open-coded: ONE bit set, nearly free
func openCoded() {
    f, _ := os.Open("file")
    defer f.Close()  // ← open-coded, ~1ns
    // ...
}

// NOT open-coded: heap-allocated each iteration
func loopDefer() {
    for _, name := range files {
        f, _ := os.Open(name)
        defer f.Close()  // ← heap-allocated, ~6ns EACH
        // BUG: files stay open until function returns, not loop iteration
    }
}

// CORRECT pattern: extract to a function
func loopDeferFixed() {
    for _, name := range files {
        processFile(name)  // defer inside processFile is open-coded
    }
}

func processFile(name string) {
    f, _ := os.Open(name)
    defer f.Close()  // ← scoped to this function, open-coded
    // ...
}
```

The compiler cannot open-code a defer inside a loop because it does not know how
many times the loop will execute at compile time, so it cannot pre-allocate bits.

---

## 3. Named Returns + defer: The Modification Trick

This is one of Go's most powerful (and most confusing) patterns. Because defers
run after the return value is set but before the caller receives it, a defer
with a named return can modify the returned value:

```go
func readFile(name string) (content string, err error) {
    f, err := os.Open(name)
    if err != nil {
        return "", err
    }
    defer func() {
        closeErr := f.Close()
        if err == nil {
            err = closeErr  // ← modifies the named return 'err'
        }
    }()

    data, err := io.ReadAll(f)
    if err != nil {
        return "", err
    }
    return string(data), nil
}
```

### What the Compiler Actually Does

A `return` statement with named returns compiles to:

```
  return string(data), nil
    ↓ compiler translates to:
  content = string(data)    // assign to named return variable
  err = nil                 // assign to named return variable
  // run deferred functions (they can read/write content, err)
  RET                       // actually return
```

The named return variables are just local variables on the stack frame. The `defer`
closure captures their addresses (capture-by-reference from Chapter 5). When the
deferred function writes `err = closeErr`, it is writing directly to the stack
slot that the caller will read.

### Production Use: Error Annotation

The most common production use of named returns + defer is adding context to errors:

```go
func (s *Service) CreateUser(ctx context.Context, req CreateReq) (user User, err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("CreateUser(%q): %w", req.Name, err)
        }
    }()

    // Every error from here on gets automatically wrapped with context
    user, err = s.repo.Insert(ctx, req)
    if err != nil {
        return  // bare return — err gets wrapped by defer
    }
    // ...
    return user, nil  // err is nil, defer does nothing
}
```

This eliminates repetitive `fmt.Errorf("CreateUser: %w", err)` at every return
point. One defer handles all of them.

---

## 4. range: The Compiler Rewrite

`range` is not a runtime operation. The compiler rewrites every `for range`
into a plain `for` loop before the code reaches the optimizer. You can see
this in `cmd/compile/internal/walk/range.go`.

### Range Over Slice

```go
// You write:
for i, v := range slice {
    process(i, v)
}

// Compiler rewrites to:
{
    _len := len(slice)
    _base := slice  // capture the slice header (pointer, len, cap)
    for _i := 0; _i < _len; _i++ {
        i, v := _i, _base[_i]  // v is a COPY of the element
        process(i, v)
    }
}
```

Key details:
- `len(slice)` is evaluated once at the start. Appending during iteration does
  not affect the loop count
- The slice header is captured, so reassigning `slice` inside the loop does not
  affect iteration
- `v` is a **copy** of the element (see section 5 for why this matters)

### Range Over Map

```go
// You write:
for k, v := range m {
    process(k, v)
}

// Compiler rewrites to (simplified):
{
    _hiter := runtime.mapiterinit(m)
    for ; _hiter.key != nil; runtime.mapiternext(_hiter) {
        k := *_hiter.key
        v := *_hiter.value
        process(k, v)
    }
}
```

The `runtime.mapiterinit` function picks a **random starting bucket** (using a
random seed stored in `hmap.hash0`). This is why map iteration order is
non-deterministic — it is deliberately randomized to prevent code from depending
on insertion order.

### Range Over String

```go
// You write:
for i, r := range "Go 🚀" {
    fmt.Printf("%d: %c\n", i, r)
}

// Compiler rewrites to (simplified):
{
    _s := "Go 🚀"
    for _i := 0; _i < len(_s); {
        r, size := utf8.DecodeRuneInString(_s[_i:])
        i := _i
        _i += size  // advance by RUNE width (1-4 bytes), not 1
        fmt.Printf("%d: %c\n", i, r)
    }
}
```

**Critical insight:** `i` is the **byte offset**, not the rune index.
For `"Go 🚀"`, the offsets are 0, 1, 2, 3 (space is 1 byte, rocket is 4 bytes).
The rocket starts at byte 3 even though it is the 4th character.

### Range Over Channel

```go
// You write:
for v := range ch {
    process(v)
}

// Compiler rewrites to:
{
    for {
        v, ok := <-ch
        if !ok {
            break  // channel closed and drained
        }
        process(v)
    }
}
```

`range` over a channel blocks on each receive and terminates when the channel
is closed AND drained. It is equivalent to receiving in a loop until `ok` is false.

---

## 5. range Gotchas: Value Copies and Loop Variables

### The Value Copy Trap

```go
type Player struct {
    Name  string
    Score int
}

players := []Player{{"Alice", 0}, {"Bob", 0}}

for _, p := range players {
    p.Score = 100  // modifies the COPY, not the slice element
}
// players[0].Score is still 0!
```

`p` is a copy of `players[i]`. Modifying `p` does not affect the original. This
catches every Java/C# developer who expects reference semantics from a for-each.

**Fix — use the index:**
```go
for i := range players {
    players[i].Score = 100  // modifies the original
}
```

**Fix — use a pointer slice:**
```go
players := []*Player{{"Alice", 0}, {"Bob", 0}}
for _, p := range players {
    p.Score = 100  // p is a copy of the pointer, but points to the same Player
}
```

### The Loop Variable Capture Trap (Pre Go 1.22)

Before Go 1.22, the loop variable was shared across iterations:

```go
// Go 1.21 and earlier — BUG
var funcs []func()
for _, v := range []int{1, 2, 3} {
    funcs = append(funcs, func() { fmt.Println(v) })
}
for _, f := range funcs {
    f()  // prints: 3, 3, 3 — all closures share the same 'v'
}
```

The compiler allocated ONE `v` variable for the entire loop. Every closure
captured the same address. By the time the closures ran, `v` held the last
value (3).

**Go 1.22 fixed this:** Each iteration now gets its own variable. The closures
capture different addresses and print 1, 2, 3. This was a breaking change enabled
by the `go` directive in `go.mod`:

```
// go.mod
module myapp
go 1.22  // enables per-iteration loop variable semantics
```

If your `go.mod` says `go 1.21` or earlier, you still get the old behavior.

---

## 6. range Over int and func (Go 1.22/1.23)

### Range Over Integer (Go 1.22)

```go
for i := range 5 {
    fmt.Println(i)  // 0, 1, 2, 3, 4
}
```

The compiler rewrites this to:

```go
for _i := 0; _i < 5; _i++ {
    i := _i
    fmt.Println(i)
}
```

Simple, clean, and eliminates the `for i := 0; i < n; i++` boilerplate.
The value is zero-indexed and excludes the upper bound (like Python's `range()`).

### Range Over Function (Go 1.23)

Go 1.23 added range over iterator functions — functions with a specific
signature that yield values via a callback:

```go
// Iterator function signature: func(yield func(V) bool)
func Fibonacci(n int) iter.Seq[int] {
    return func(yield func(int) bool) {
        a, b := 0, 1
        for range n {
            if !yield(a) {
                return  // caller broke out of the loop
            }
            a, b = b, a+b
        }
    }
}

// Usage — looks like a normal range loop
for v := range Fibonacci(10) {
    fmt.Println(v)
}
```

The `iter.Seq[V]` and `iter.Seq2[K, V]` types from the `iter` package define
the standard iterator signatures. The `slices` and `maps` packages in Go 1.23
provide iterator functions: `slices.All()`, `slices.Values()`, `maps.Keys()`, etc.

**Under the hood,** the compiler transforms range-over-func into a callback
invocation pattern. The `yield` function is the loop body wrapped in a function.
When the loop body executes `break`, the yield function returns `false`,
signaling the iterator to stop.

---

## 7. switch: Jump Tables and Type Dispatch

### Expression Switch

Go's `switch` does NOT fall through by default (unlike C/Java). Each case is
independent. Use `fallthrough` explicitly if you want C behavior (rare in Go).

```go
switch day {
case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
    return "weekday"
case "Saturday", "Sunday":
    return "weekend"
default:
    return "unknown"
}
```

The compiler chooses between two strategies:

```
  Few cases (< ~8):   Linear comparison chain
                      if day == "Monday" || day == "Tuesday" || ...
                      Simple, branch predictor handles it well

  Many cases:         Jump table (for integer switches)
                      or hash-based dispatch (for string switches)
                      O(1) lookup instead of O(n) comparisons
```

You can see the compiler's decision with:
```bash
go build -gcflags='-S' ./... 2>&1 | grep -i "jump\|switch"
```

### Type Switch

Type switches use the interface's type information (`iface.tab._type`) to
dispatch:

```go
func describe(i interface{}) string {
    switch v := i.(type) {
    case int:
        return fmt.Sprintf("int: %d", v)
    case string:
        return fmt.Sprintf("string: %s", v)
    case bool:
        return fmt.Sprintf("bool: %v", v)
    default:
        return "unknown"
    }
}
```

Under the hood, the compiler generates a series of type comparisons against
the `_type` pointer from the interface header:

```
  // Pseudo-assembly for type switch:
  load iface.tab._type into R1
  compare R1 with &type.int     → jump to case_int
  compare R1 with &type.string  → jump to case_string
  compare R1 with &type.bool    → jump to case_bool
  jump to default
```

Each type has a unique `_type` descriptor in the binary (assigned at compile time),
so the comparison is just a pointer equality check — very fast.

### Switch vs If-Else: When to Use Each

```
  Use switch when:
    - Comparing ONE value against multiple options
    - Type switching on an interface
    - Replacing long if-else chains (readability win)
    - You want the compiler to check for duplicate cases

  Use if-else when:
    - Conditions are independent expressions (not comparing same value)
    - You need complex boolean logic
    - Conditions depend on results of previous checks
```

---

## 8. for: The Only Loop, Labeled Breaks, and Go 1.22

Go has exactly one loop construct: `for`. There is no `while`, no `do-while`,
no `foreach`. This is a deliberate simplicity decision — one keyword, three forms:

```go
for i := 0; i < n; i++ { }   // C-style for
for condition { }              // while (other languages)
for { }                        // infinite loop
```

### Labeled Break and Continue

When you have nested loops, `break` only exits the innermost loop. Labels let
you break or continue an outer loop:

```go
outer:
    for i := 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            if matrix[i][j] == target {
                break outer  // exits BOTH loops
            }
        }
    }
```

Without the label, `break` would only exit the inner `j` loop. Labeled breaks
are common in Go for search patterns in nested data structures.

### Go 1.22: Per-Iteration Variable Semantics

As covered in section 5, Go 1.22 changed loop variable scoping. But there is
a subtle detail about **three-clause for loops** too:

```go
// Go 1.22+: each iteration gets its own 'i'
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // Go 1.22: prints 0,1,2,3,4 (in some order)
                         // Go 1.21: prints 5,5,5,5,5
    }()
}
```

The Go team retroactively enabled this for three-clause `for` loops as well,
not just `for range`. The rationale: "the bug from sharing loop variables is
never the behavior anyone wants" (Russ Cox, Go proposal #60078).

### The Compiler's Loop Optimizations

The Go compiler applies several optimizations to `for` loops:

```
  Bounds check elimination (BCE):
    If the compiler can prove i < len(slice), it eliminates the
    bounds check on slice[i]. Loop induction variable analysis
    is the main enabler.

  Loop unrolling:
    For very tight loops with known iteration counts, the compiler
    may unroll the loop body (repeat it N times to reduce branch
    overhead). Less aggressive than C compilers.

  Strength reduction:
    Multiplication by a constant in the loop index (e.g., i*8)
    is replaced by repeated addition (base += 8).
```

You can observe BCE with:
```bash
go build -gcflags='-d=ssa/check_bce/debug=1' ./...
```

---

## 9. Quick Reference Card

```
  ┌──────────────────────────────────────────────────────────────────┐
  │                  CONTROL FLOW QUICK REFERENCE                    │
  ├──────────────────────────────────────────────────────────────────┤
  │                                                                  │
  │  DEFER RULES                                                     │
  │    1. Arguments evaluated at defer statement (not execution)     │
  │    2. Defers execute LIFO (last deferred runs first)             │
  │    3. Defers run after return value set, before caller gets it   │
  │                                                                  │
  │  DEFER IMPLEMENTATIONS                                           │
  │    Pre-1.13:  Heap-allocated _defer struct       ~35ns           │
  │    Go 1.13:   Stack-allocated _defer             ~6ns            │
  │    Go 1.14+:  Open-coded (bitmask, ≤8 defers)   ~1ns            │
  │    In loops:  Falls back to stack-allocated       ~6ns           │
  │                                                                  │
  │  RANGE REWRITES                                                  │
  │    []T:     len evaluated once, v is a COPY of element           │
  │    map:     random start bucket, non-deterministic order         │
  │    string:  iterates RUNES (not bytes), i is byte offset         │
  │    chan:    blocks on receive, stops when closed+drained          │
  │    int:     range N = 0..N-1 (Go 1.22+)                         │
  │    func:    iter.Seq[V] callback pattern (Go 1.23+)              │
  │                                                                  │
  │  LOOP VARIABLE CHANGE (Go 1.22)                                  │
  │    Before:  One variable shared across all iterations             │
  │    After:   Each iteration gets its own copy                     │
  │    Applies: for-range AND three-clause for loops                 │
  │    Trigger: go directive in go.mod (go 1.22 or later)            │
  │                                                                  │
  │  SWITCH                                                          │
  │    No fallthrough by default (unlike C/Java)                     │
  │    Type switch: pointer comparison on _type descriptor           │
  │    Multiple values per case: case "a", "b", "c":                 │
  │                                                                  │
  │  LABELED BREAK/CONTINUE                                          │
  │    outer: for { for { break outer } }                            │
  │    Exits the labeled loop, not just the innermost                │
  │                                                                  │
  │  DEBUGGING                                                       │
  │    go build -gcflags='-S' ./...              # see defer impl    │
  │    go build -gcflags='-d=ssa/check_bce/debug=1'  # BCE checks   │
  │    go vet ./...                              # loop variable     │
  │                                                                  │
  └──────────────────────────────────────────────────────────────────┘
```

---

## Further Reading

- `runtime/panic.go` : the `_defer` struct definition and `deferreturn` function
  that processes the defer chain at function exit.
- `cmd/compile/internal/ssagen/ssa.go` : search for `openDeferRecord` to see
  the open-coded defer implementation decisions.
- `cmd/compile/internal/walk/range.go` : the compiler's range loop rewrites.
  Each range type (`TSLICE`, `TMAP`, `TSTRING`, `TCHAN`, `TINT`) has its own
  rewrite function.
- [Go 1.22 Release Notes — Range Over Integer](https://go.dev/doc/go1.22#language)
  explains the range-over-int addition and loop variable semantics change.
- [Go Proposal #60078](https://github.com/golang/go/issues/60078) — Russ Cox's
  proposal for per-iteration loop variable semantics, with data showing this class
  of bug in real Google codebases.
- [Go 1.23 Release Notes — Range Over Function](https://go.dev/doc/go1.23#language)
  explains the iterator function protocol.
- [Go Blog: "Backward Compatibility, Go 1.21, and Go 1.22"](https://go.dev/blog/compat)
  explains how the `go` directive in `go.mod` gates the loop variable change.
