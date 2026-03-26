# Deep Dive: Loop Variable Capture & The Fan-In Pattern

> Why closures in loops all see the same value, what the compiler does under
> the hood, how Go 1.22 partially fixed it, and how the fan-in pattern uses
> this knowledge to safely coordinate concurrent goroutines.

---

## Table of Contents

1. [The Problem — One Address, Many Closures](#1-the-problem--one-address-many-closures)
2. [Memory Layout — Step by Step Through the Loop](#2-memory-layout--step-by-step-through-the-loop)
3. [What the Compiler Actually Generates](#3-what-the-compiler-actually-generates)
4. [The Three Fixes and Their Mechanics](#4-the-three-fixes-and-their-mechanics)
5. [Go 1.22 Range Loop Change — What the Compiler Does Differently](#5-go-122-range-loop-change--what-the-compiler-does-differently)
6. [C-Style For Loops — Still Broken in Go 1.25+](#6-c-style-for-loops--still-broken-in-go-125)
7. [All the Ways This Bites You](#7-all-the-ways-this-bites-you)
8. [Production Case Study: The Fan-In Pattern — Line by Line](#8-production-case-study-the-fan-in-pattern--line-by-line)
9. [Fan-In Execution Timeline](#9-fan-in-execution-timeline)
10. [Fan-In Variations and Tradeoffs](#10-fan-in-variations-and-tradeoffs)
11. [Detection Tools](#11-detection-tools)
12. [How Other Languages Handle This](#12-how-other-languages-handle-this)
13. [Quick Reference Card](#13-quick-reference-card)

---

## 1. The Problem — One Address, Many Closures

Here's the classic broken code:

```go
funcs := make([]func(), 5)
for i := 0; i < 5; i++ {
    funcs[i] = func() { fmt.Println(i) }
}
for _, f := range funcs { f() }
// Output: 5, 5, 5, 5, 5   ← all print 5, not 0,1,2,3,4
```

### Why This Happens — The Single Variable Rule

In a C-style `for` loop, Go creates **exactly one variable `i`** for the entire loop.
Every closure created inside the loop captures **the address of that same `i`** — not
a copy of its current value.

```
The key insight: closures capture VARIABLES (addresses), not VALUES (copies).
```

Think of it this way: when you write `func() { fmt.Println(i) }`, the compiler
does NOT snapshot `i`'s current value. Instead, it stores a **pointer to `i`**.
When the closure is eventually called, it follows that pointer and reads whatever
value `i` has **at that moment** — which is 5 (the value that ended the loop).

This is identical to how it works in languages you already know:
- **JavaScript (pre-`let`)**: `var i` in a `for` loop had the same bug. `let` fixed it.
- **C#**: Before C# 5.0, `foreach` had this exact behavior. Fixed in C# 5.0.
- **Python**: Same problem with closures in list comprehensions over a loop variable.

Go's version is more subtle because it looks like `i` should be scoped per iteration,
but it isn't — until Go 1.22 for range loops, and **still isn't** for C-style loops.

---

## 2. Memory Layout — Step by Step Through the Loop

Let's trace exactly what happens in memory during the broken loop:

```go
funcs := make([]func(), 5)
for i := 0; i < 5; i++ {
    funcs[i] = func() { fmt.Println(i) }
}
```

### Before the Loop Starts

```
Stack frame of enclosing function:
┌──────────────────────────────────────────────────────┐
│  funcs: []func()  → backing array on heap            │
│  i:     int       → addr: 0xC000012080  value: 0     │  ← ONE variable
└──────────────────────────────────────────────────────┘

Note: i is allocated on the HEAP, not the stack. Why?
Because closures capture its address, and those closures
outlive the loop → escape analysis moves i to the heap.
```

### Iteration 0: i = 0

```
Heap:
  i (0xC000012080): 0

  funcs[0] ──► funcval {
                 fn:  closure_code_ptr
                 &i:  0xC000012080      ← points to the SAME i
               }
```

### Iteration 1: i = 1

```
Heap:
  i (0xC000012080): 1           ← value changed, address SAME

  funcs[0] ──► funcval { &i: 0xC000012080 }    ← still points here
  funcs[1] ──► funcval { &i: 0xC000012080 }    ← also points here
```

### Iteration 2: i = 2

```
Heap:
  i (0xC000012080): 2           ← value changed again, address SAME

  funcs[0] ──► funcval { &i: 0xC000012080 }    ← all three
  funcs[1] ──► funcval { &i: 0xC000012080 }    ← point to
  funcs[2] ──► funcval { &i: 0xC000012080 }    ← same address
```

### After Loop Ends (i = 5, loop condition false)

```
Heap:
  i (0xC000012080): 5           ← final value after i++ makes condition false

  funcs[0] ──► funcval { &i: 0xC000012080 }  ──┐
  funcs[1] ──► funcval { &i: 0xC000012080 }  ──┤
  funcs[2] ──► funcval { &i: 0xC000012080 }  ──┼── ALL point to same i
  funcs[3] ──► funcval { &i: 0xC000012080 }  ──┤     which now holds 5
  funcs[4] ──► funcval { &i: 0xC000012080 }  ──┘

  funcs[0]()  → follow pointer → read *0xC000012080 → 5
  funcs[1]()  → follow pointer → read *0xC000012080 → 5
  ...all print 5
```

### The Core Visual

```
  ┌─────────────────────────────────────────────────────────────────┐
  │                    THE BUG IN ONE PICTURE                       │
  │                                                                 │
  │   ONE variable i on the heap:    [  5  ]  ← addr 0xC000012080  │
  │                                   ▲ ▲ ▲ ▲ ▲                    │
  │                                   │ │ │ │ │                     │
  │   funcs[0].&i ─────────────────────┘ │ │ │ │                    │
  │   funcs[1].&i ───────────────────────┘ │ │ │                    │
  │   funcs[2].&i ─────────────────────────┘ │ │                    │
  │   funcs[3].&i ───────────────────────────┘ │                    │
  │   funcs[4].&i ─────────────────────────────┘                    │
  │                                                                 │
  │   5 closures × 1 variable = all see the same value              │
  └─────────────────────────────────────────────────────────────────┘
```

---

## 3. What the Compiler Actually Generates

You can see the compiler's closure transformation using:

```bash
go build -gcflags='-m -m' ./...     # escape analysis with reasons
go build -gcflags='-S' ./...        # assembly output
```

### The Broken Version — Compiler's View

```go
// What you write:
for i := 0; i < 5; i++ {
    funcs[i] = func() { fmt.Println(i) }
}

// What the compiler conceptually generates:
i := new(int)             // heap-allocated because closures capture &i
*i = 0
for ; *i < 5; (*i)++ {
    funcs[*i] = &funcval{
        fn:  closure_body_ptr,
        captured_i: i,        // pointer to the SAME heap int
    }
}

// The closure body:
func closure_body(fv *funcval) {
    i := fv.captured_i        // load pointer from funcval
    fmt.Println(*i)           // dereference → reads current value of i
}
```

### Escape Analysis Output

```bash
$ go build -gcflags='-m' main.go
./main.go:5:6: moved to heap: i            ← i escapes because closures capture it
./main.go:6:15: func literal escapes to heap ← closure stored in slice
```

**Key:** The "moved to heap: i" message is your signal. If `i` didn't escape
(e.g., closure only used locally), it would stay on the stack and the problem
wouldn't exist — but then you wouldn't be storing closures for later use.

---

## 4. The Three Fixes and Their Mechanics

### Fix 1: Variable Shadowing (`i := i`)

```go
for i := 0; i < 5; i++ {
    i := i    // NEW variable, shadows the loop's i
    funcs[i] = func() { fmt.Println(i) }
}
```

**What happens in memory:**

```
Iteration 0:
  loop_i (0xC000012080): 0
  shadow_i_0 (0xC000012088): 0     ← NEW heap allocation
  funcs[0].&i → 0xC000012088       ← points to shadow_i_0

Iteration 1:
  loop_i (0xC000012080): 1
  shadow_i_1 (0xC000012090): 1     ← ANOTHER new heap allocation
  funcs[1].&i → 0xC000012090       ← points to shadow_i_1

After loop:
  loop_i: 5
  shadow_i_0: 0   ← funcs[0] reads this → prints 0 ✅
  shadow_i_1: 1   ← funcs[1] reads this → prints 1 ✅
  shadow_i_2: 2   ← funcs[2] reads this → prints 2 ✅
  ...
```

```
  ┌─────────────────────────────────────────────────────────────────┐
  │                    THE FIX IN ONE PICTURE                       │
  │                                                                 │
  │   FIVE separate variables on the heap:                          │
  │                                                                 │
  │   shadow_i_0: [  0  ]  ← funcs[0].&i points here               │
  │   shadow_i_1: [  1  ]  ← funcs[1].&i points here               │
  │   shadow_i_2: [  2  ]  ← funcs[2].&i points here               │
  │   shadow_i_3: [  3  ]  ← funcs[3].&i points here               │
  │   shadow_i_4: [  4  ]  ← funcs[4].&i points here               │
  │                                                                 │
  │   5 closures × 5 variables = each sees its own value            │
  └─────────────────────────────────────────────────────────────────┘
```

**Cost:** One extra heap allocation per iteration. For most code, irrelevant.
In a hot path running millions of iterations, use Fix 2 instead.

### Fix 2: Pass as Function Argument (Copy)

```go
for i := 0; i < 5; i++ {
    funcs[i] = func(n int) func() {
        return func() { fmt.Println(n) }
    }(i)    // i is COPIED into n at call time
}
```

Or more commonly:

```go
for i := 0; i < 5; i++ {
    func(n int) {
        funcs[n] = func() { fmt.Println(n) }
    }(i)
}
```

**What happens:** `i` is copied by value into parameter `n`. Each invocation
creates a new `n` on the stack (or heap if the inner closure escapes). The
closure captures `&n`, which is unique per call.

**This is exactly what the `fanIn` pattern does:** `go func(s <-chan int){...}(src)`
copies `src` into `s` at launch time.

### Fix 3: Go 1.22+ Range Loops (Automatic)

```go
// Go 1.22+: range loop variables are per-iteration. No fix needed.
for i, v := range slice {
    funcs[i] = func() { fmt.Println(i, v) }   // ✅ each closure gets its own i, v
}
```

The compiler generates a new variable per iteration automatically. See Section 5.

---

## 5. Go 1.22 Range Loop Change — What the Compiler Does Differently

### The Proposal

Go proposal [#60078](https://github.com/golang/go/issues/60078) changed the semantics
of range loop variables starting in Go 1.22. This was one of the most significant
language changes since Go 1.0.

### Before Go 1.22 — Range Loop Compiled As:

```go
// Source:
for i, v := range slice { body(i, v) }

// Compiler generated (conceptual):
{
    i := 0              // ONE i for the whole loop
    v := zero_value_T   // ONE v for the whole loop
    for ; i < len(slice); i++ {
        v = slice[i]
        body(i, v)      // all iterations share the same i, v variables
    }
}
```

### Go 1.22+ — Range Loop Compiled As:

```go
// Source (same code!):
for i, v := range slice { body(i, v) }

// Compiler now generates (conceptual):
{
    for i_tmp := 0; i_tmp < len(slice); i_tmp++ {
        i := i_tmp          // NEW i per iteration
        v := slice[i_tmp]   // NEW v per iteration
        body(i, v)          // each iteration has its own variables
    }
}
```

### Memory Layout Comparison

```
BEFORE Go 1.22 (range loop):

  ONE i (0xA0): 0 → 1 → 2 → 3        ONE v (0xA8): "a" → "b" → "c" → "d"
                 ▲ ▲ ▲ ▲                              ▲   ▲   ▲   ▲
                 │ │ │ │                              │   │   │   │
  closures all point to same address    closures all point to same address
  → all see final value 3              → all see final value "d"


AFTER Go 1.22 (range loop):

  i_0 (0xA0): 0    v_0 (0xB0): "a"    ← iteration 0 variables
  i_1 (0xA8): 1    v_1 (0xB8): "b"    ← iteration 1 variables
  i_2 (0xC0): 2    v_2 (0xC8): "c"    ← iteration 2 variables
  i_3 (0xD0): 3    v_3 (0xD8): "d"    ← iteration 3 variables
       ▲                ▲
       │                │
  each closure          each closure
  has its own i         has its own v
```

### Why Only Range Loops?

The Go team's rationale (from Russ Cox's [blog post](https://go.dev/blog/loopvar-preview)):

1. **Range loops iterate over a known collection** — the variable represents "the
   current element." It's natural for each iteration to have its own element.

2. **C-style `for` loops are mutation-based** — `for i := 0; i < n; i++` explicitly
   says "I have a variable `i`, and I mutate it." Changing this semantic would break
   the mental model of what `:=` and `++` mean.

3. **Backward compatibility:** The Go team analyzed all public Go code on GitHub.
   Range loop variable reuse caused bugs. C-style loop variable reuse rarely did
   (because it's less common to capture `i` from a C-style loop in a closure).

### Verifying the Behavior

```bash
# Check your Go version's behavior:
go version

# If Go 1.22+, range loops are per-iteration by default.
# You can see what the compiler does:
GOEXPERIMENT=loopvar go build -gcflags='-m' ./...
```

---

## 6. C-Style For Loops — Still Broken in Go 1.25+

This is critical to remember: **even in Go 1.25.7 (your version), C-style loops
are NOT fixed:**

```go
// STILL BROKEN — even in Go 1.25.7:
funcs := make([]func(), 5)
for i := 0; i < 5; i++ {
    funcs[i] = func() { fmt.Println(i) }   // ❌ all see 5
}

// STILL NEED the fix:
for i := 0; i < 5; i++ {
    i := i                                   // ✅ shadow
    funcs[i] = func() { fmt.Println(i) }
}
```

### Why This Matters in Practice

C-style loops are common when:
- Iterating with custom step sizes: `for i := 0; i < n; i += 2`
- Countdown loops: `for i := n; i >= 0; i--`
- Complex loop conditions: `for i := start; i != end; i = nextNode(i)`
- Index-based iteration where you need to skip or repeat elements

In all these cases, if you capture `i` in a closure, you must still use `i := i`
or pass as a parameter.

---

## 7. All the Ways This Bites You

### 7.1 Slice of Closures (Most Common in Interviews)

```go
funcs := make([]func(), 3)
for i := 0; i < 3; i++ {
    funcs[i] = func() { fmt.Println(i) }
}
// All print 3
```

### 7.2 Goroutines in a Loop (Most Common in Production)

```go
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)    // race condition + wrong values
    }()
}
```

**Two bugs at once:**
1. All goroutines share one `i` → wrong values
2. `i++` in the loop races with goroutine reads → data race

```bash
$ go test -race ./...
WARNING: DATA RACE
  Write at 0xC000012080 by goroutine 1:    ← the for loop's i++
  Previous read at 0xC000012080 by goroutine 7:  ← the closure reading i
```

### 7.3 Deferred Closures in a Loop

```go
for i := 0; i < 3; i++ {
    defer func() { fmt.Println(i) }()   // all print 3 (LIFO: 3, 3, 3)
}
```

Defers run when the **enclosing function** returns, not when the iteration ends.
By that time, `i` is 3.

### 7.4 HTTP Handlers in a Loop

```go
for _, route := range routes {
    http.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Route: %s", route.Name)   // Go 1.22+: ✅ fixed for range
    })
}
```

**Pre-Go 1.22:** All handlers would serve the last route's name. This was one of the
most reported bugs in real Go web applications.

**Go 1.22+:** Fixed automatically for range loops.

### 7.5 Event Callbacks

```go
buttons := []Button{...}
for i := 0; i < len(buttons); i++ {
    buttons[i].OnClick = func() {
        fmt.Printf("Clicked button %d\n", i)   // ❌ all say last index
    }
}
```

### 7.6 The Subtle Struct Field Variant

```go
type Task struct {
    ID   int
    Name string
}

tasks := []Task{{1, "A"}, {2, "B"}, {3, "C"}}
var handlers []func()

// C-style loop — BROKEN even in Go 1.25+:
for i := 0; i < len(tasks); i++ {
    handlers = append(handlers, func() {
        fmt.Println(tasks[i].Name)    // ❌ i is captured, not the task
    })                                // panics if i == len(tasks) at call time
}

// Range loop — FIXED in Go 1.22+:
for _, t := range tasks {
    handlers = append(handlers, func() {
        fmt.Println(t.Name)           // ✅ t is per-iteration in Go 1.22+
    })
}
```

---

## 8. Production Case Study: The Fan-In Pattern — Line by Line

Now let's break down the `fanIn` function. This is a fundamental concurrency
pattern that demonstrates correct closure handling in production code.

### The Full Code

```go
func fanIn(sources ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(sources))

    for _, src := range sources {
        go func(s <-chan int) {
            defer wg.Done()
            for v := range s {
                out <- v
            }
        }(src)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### What This Function Does — The Big Picture

```
  ┌──────────┐
  │ source 1 │──► chan int ──┐
  └──────────┘               │
  ┌──────────┐               │    ┌─────────┐
  │ source 2 │──► chan int ──┼───►│   out   │──► chan int ──► consumer
  └──────────┘               │    │ (merged)│
  ┌──────────┐               │    └─────────┘
  │ source 3 │──► chan int ──┘
  └──────────┘

  "Fan-In" = many input channels merged into one output channel.
  
  Use cases:
  - Multiple producers, one consumer
  - Aggregating results from parallel workers
  - Merging log streams from multiple services
  - Combining search results from multiple backends
```

### Line-by-Line Breakdown

#### Line 1: `func fanIn(sources ...<-chan int) <-chan int {`

```
  sources ...<-chan int
  ├── ...           = variadic: accepts any number of arguments
  ├── <-chan int    = receive-only channel of ints
  └── The <- before chan means: callers can ONLY READ from these channels.
      fanIn does not need to send to them, so we restrict the type.

  <-chan int (return type)
  └── Returns a receive-only channel. The consumer can only read from it.
      fanIn owns the write side internally.

  Direction restriction is a COMPILE-TIME safety check:
  - If fanIn accidentally tried to send to a source: compile error
  - If the consumer accidentally tried to send to out: compile error
```

#### Line 2: `out := make(chan int)`

```
  out is an UNBUFFERED channel.
  
  Why unbuffered?
  - Backpressure: if the consumer stops reading, the goroutines block
    on "out <- v" instead of filling a buffer and wasting memory.
  - Simplicity: no need to guess buffer size.
  - Correctness: the consumer processes values at its own pace.

  Could use make(chan int, bufSize) for throughput optimization,
  but unbuffered is the correct default — optimize only with data.
  
  Memory layout (from our channels session):
  ┌───────────────────┐
  │ hchan (on heap)   │
  │  qcount:    0     │  ← no buffer
  │  dataqsiz:  0     │
  │  buf:       nil   │  ← no ring buffer allocated
  │  elemsize:  8     │  ← sizeof(int)
  │  closed:    0     │
  │  sendq:     nil   │  ← senders wait here
  │  recvq:     nil   │  ← receivers wait here
  │  lock:      0     │
  └───────────────────┘
```

#### Lines 3-4: `var wg sync.WaitGroup` / `wg.Add(len(sources))`

```
  WaitGroup is a concurrent counter:
  
  wg.Add(3)  → counter = 3
  wg.Done()  → counter-- (atomic decrement)
  wg.Wait()  → blocks until counter == 0

  We call Add(len(sources)) BEFORE launching goroutines — this is critical.
  
  ┌─────────────────────────────────────────────────────────────────────┐
  │  WHY Add() BEFORE the loop, not inside it?                         │
  │                                                                     │
  │  If you do wg.Add(1) inside the loop, there's a race:              │
  │                                                                     │
  │    go func() { wg.Wait(); close(out) }()  // coordinator starts    │
  │    for _, src := range sources {                                    │
  │        wg.Add(1)    // ← might run AFTER Wait() returns!           │
  │        go func(s <-chan int) { ... }(src)                           │
  │    }                                                                │
  │                                                                     │
  │  The Wait() goroutine could see counter==0 between iterations      │
  │  and close(out) prematurely.                                        │
  │                                                                     │
  │  FIX: wg.Add(len(sources)) before ANY goroutine launches.          │
  └─────────────────────────────────────────────────────────────────────┘
```

#### Lines 6-12: The Worker Loop

```go
    for _, src := range sources {
        go func(s <-chan int) {
            defer wg.Done()
            for v := range s {
                out <- v
            }
        }(src)
    }
```

Let's break this into layers:

**Layer 1: The Loop Variable Capture Fix**

```
  for _, src := range sources {
      go func(s <-chan int) { ... }(src)
  }
        │           │            │
        │           │            └── (src) = COPY src's value into s at launch time
        │           └── s <-chan int = parameter, receives the copy
        └── src = loop variable

  In Go 1.22+, range variables are per-iteration, so this would work:
      go func() { ... use src ... }()    // ✅ safe in Go 1.22+

  But the parameter-passing style (src → s) is STILL preferred because:
  1. It works in ALL Go versions
  2. It's self-documenting: you can SEE the value being passed
  3. It's the convention in the Go community
  4. It works for C-style loops too (future-proof)
```

**Layer 2: What Each Goroutine Does**

```
  go func(s <-chan int) {       // s is this goroutine's own copy of the channel
      defer wg.Done()           // when this goroutine exits, decrement the counter
      for v := range s {        // read from s until s is closed
          out <- v              // forward every value to the merged output
      }
  }(src)

  "for v := range s" desugars to:
      for {
          v, ok := <-s
          if !ok { break }     // channel closed, stop reading
          out <- v
      }

  When source channel closes → range exits → goroutine reaches end → 
  defer fires → wg.Done() → counter decrements
```

**Layer 3: Memory Layout for 3 Sources**

```
  After the loop, 3 goroutines are running independently:

  ┌─────────────────────────────────────────────────────────────┐
  │  G1 (goroutine)                                             │
  │  Stack: s = sources[0]  (own copy of channel pointer)       │
  │  Running: for v := range s { out <- v }                     │
  │  defer: wg.Done()                                           │
  ├─────────────────────────────────────────────────────────────┤
  │  G2 (goroutine)                                             │
  │  Stack: s = sources[1]  (own copy of channel pointer)       │
  │  Running: for v := range s { out <- v }                     │
  │  defer: wg.Done()                                           │
  ├─────────────────────────────────────────────────────────────┤
  │  G3 (goroutine)                                             │
  │  Stack: s = sources[2]  (own copy of channel pointer)       │
  │  Running: for v := range s { out <- v }                     │
  │  defer: wg.Done()                                           │
  └─────────────────────────────────────────────────────────────┘
  
  All 3 goroutines write to the SAME out channel.
  The channel's internal mutex serializes concurrent sends.
  Values arrive in non-deterministic order — whichever goroutine
  gets scheduled first writes first.
```

**Important:** A channel value (`chan int`) is already a pointer to the `hchan`
struct on the heap. So `s = src` copies the pointer, not the channel. Both `s` and
`src` refer to the same underlying `hchan`. This is fine — channels are designed
for concurrent access.

#### Lines 14-17: The Coordinator Goroutine

```go
    go func() {
        wg.Wait()    // wait for ALL senders to finish
        close(out)   // single close point
    }()
```

```
  This goroutine does TWO things:
  1. wg.Wait() — blocks until ALL worker goroutines call wg.Done()
  2. close(out) — signals the consumer that no more values are coming

  ┌─────────────────────────────────────────────────────────────────────┐
  │  WHY THIS PATTERN IS ESSENTIAL                                      │
  │                                                                     │
  │  Problem: Who closes `out`?                                         │
  │                                                                     │
  │  - NOT the workers: if G1 finishes first and closes out,            │
  │    G2 and G3 will panic ("send on closed channel") when they        │
  │    try to send their values.                                        │
  │                                                                     │
  │  - NOT the consumer: the consumer doesn't know when all             │
  │    producers are done.                                              │
  │                                                                     │
  │  - The COORDINATOR: a separate goroutine that waits for ALL         │
  │    workers to finish, THEN closes the channel. This is the          │
  │    only safe approach.                                              │
  │                                                                     │
  │  This is the same principle from our channels session:              │
  │  "Only the sender (or coordinator) should close the channel."       │
  └─────────────────────────────────────────────────────────────────────┘
```

**Why a separate goroutine for the coordinator?**

```
  If we did this synchronously:

    wg.Wait()      // ← BLOCKS here until all workers finish
    close(out)     // ← never reached until workers done
    return out     // ← never reached!

  The function would NEVER RETURN because wg.Wait() blocks.
  The consumer would never get the channel.
  
  By putting it in a goroutine:
  
    go func() { wg.Wait(); close(out) }()
    return out     // ← returns immediately! Consumer can start reading.

  The coordinator runs concurrently:
  - Consumer reads from out
  - Workers send to out
  - Coordinator waits in the background
  - When all workers finish → coordinator closes out → consumer's range exits
```

#### Line 19: `return out`

```
  Returns the merged channel immediately, before any values are produced.
  
  The consumer will use it like:
  
    merged := fanIn(ch1, ch2, ch3)
    for v := range merged {
        process(v)         // receives values from ALL sources, interleaved
    }
    // range exits when coordinator calls close(out)
    fmt.Println("all sources exhausted")
```

---

## 9. Fan-In Execution Timeline

Here's a complete timeline showing how all the pieces interact:

```
TIME  │ Main Goroutine    │ G1 (src[0])      │ G2 (src[1])     │ Coordinator
──────┼───────────────────┼──────────────────┼─────────────────┼──────────────
 t0   │ fanIn(s0, s1, s2) │                  │                 │
 t1   │ out = make(chan)   │                  │                 │
 t2   │ wg.Add(3)         │                  │                 │
 t3   │ spawn G1, G2, G3  │ start            │ start           │
 t4   │ spawn coordinator │                  │                 │ wg.Wait()
 t5   │ return out         │                  │                 │  (blocked,
      │                   │                  │                 │   counter=3)
──────┼───────────────────┼──────────────────┼─────────────────┼──────────────
 t6   │ consumer reads    │ v=10 from s0     │ (waiting for    │
      │ range merged      │ out <- 10        │  data on s1)    │
 t7   │ process(10)       │                  │ v=20 from s1    │
      │                   │                  │ out <- 20       │
 t8   │ process(20)       │ v=11 from s0     │                 │
      │                   │ out <- 11        │                 │
──────┼───────────────────┼──────────────────┼─────────────────┼──────────────
 t9   │ process(11)       │ s0 closed!       │ v=21 from s1    │
      │                   │ range exits      │ out <- 21       │
 t10  │ process(21)       │ wg.Done()        │                 │ counter: 3→2
      │                   │ G1 exits         │                 │
──────┼───────────────────┼──────────────────┼─────────────────┼──────────────
 t11  │                   │                  │ s1 closed!      │
      │                   │                  │ range exits     │
 t12  │                   │                  │ wg.Done()       │ counter: 2→1
      │                   │                  │ G2 exits        │
──────┼───────────────────┼──────────────────┼─────────────────┼──────────────
 t13  │                   │                  │                 │ (G3 also done)
      │                   │                  │                 │ counter: 1→0
      │                   │                  │                 │ Wait() returns
 t14  │                   │                  │                 │ close(out)
──────┼───────────────────┼──────────────────┼─────────────────┼──────────────
 t15  │ range exits       │                  │                 │ Coordinator
      │ (out is closed)   │                  │                 │ exits
      │ "all done!"       │                  │                 │
```

### The Critical Invariants

```
  1. wg.Add(N) happens BEFORE any goroutine starts     → no premature Wait() return
  2. Each goroutine calls wg.Done() exactly ONCE        → counter reaches exactly 0
  3. close(out) happens AFTER all goroutines finish      → no "send on closed channel"
  4. close(out) happens exactly ONCE                     → no "close of closed channel"
  5. return out happens BEFORE data flows                → consumer is ready to receive
```

---

## 10. Fan-In Variations and Tradeoffs

### Variation 1: With Context Cancellation

```go
func fanIn(ctx context.Context, sources ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(sources))

    for _, src := range sources {
        go func(s <-chan int) {
            defer wg.Done()
            for v := range s {
                select {
                case out <- v:
                case <-ctx.Done():    // stop if context cancelled
                    return            // drain stops, goroutine exits
                }
            }
        }(src)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

**Why:** Without context, if the consumer stops reading (crashes, timeout),
the worker goroutines block forever on `out <- v` → **goroutine leak**.
The `select` with `ctx.Done()` lets them exit cleanly.

### Variation 2: With Buffered Output Channel

```go
out := make(chan int, len(sources) * 10)    // buffer for burst absorption
```

**When:** High-throughput scenarios where producer bursts shouldn't block.
**Tradeoff:** Uses memory for the buffer. Values sit in buffer if consumer is slow.

### Variation 3: First-Value-Wins (Fan-In with Select)

```go
// Only works for exactly 2 sources (or use reflect.Select for N)
func fanInTwo(a, b <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for a != nil || b != nil {
            select {
            case v, ok := <-a:
                if !ok { a = nil; continue }
                out <- v
            case v, ok := <-b:
                if !ok { b = nil; continue }
                out <- v
            }
        }
    }()
    return out
}
```

**Tradeoff:** Single goroutine instead of N, but only practical for small N.
Uses `select` for multiplexing. Setting channel to `nil` disables that `select` case
(a `nil` channel blocks forever in `select`).

---

## 11. Detection Tools

### go vet — Catches Loop Capture in Goroutines

```bash
$ go vet ./...
./main.go:10:24: loop variable i captured by func literal
```

`go vet` specifically looks for the pattern:
`go func() { ... uses loop var ... }()`

### Race Detector — Catches Concurrent Access

```bash
$ go test -race ./...
WARNING: DATA RACE
  Write at 0x00c000012080 by goroutine 1:    ← loop's i++
  Read at 0x00c000012080 by goroutine 6:     ← closure reading i
```

### staticcheck / golangci-lint

```bash
$ staticcheck ./...
SA6005: loop variable captured by closure (staticcheck)

$ golangci-lint run
main.go:10:24: G601: loop variable i captured by func literal (gosec)
```

### Escape Analysis — See What Moves to Heap

```bash
$ go build -gcflags='-m' main.go
./main.go:5:6: moved to heap: i          ← i captured by escaping closure
./main.go:6:15: func literal escapes     ← closure stored/goroutined
```

---

## 12. How Other Languages Handle This

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

Key insight: Go 1.22's range fix mirrors C# 5.0's foreach fix. Both languages
decided the "iteration variable" semantic was more useful than the "mutating
counter" semantic for their foreach/range construct, while leaving the C-style
for loop unchanged.
```

---

## 13. Quick Reference Card

```
LOOP CAPTURE BUG
  Closures capture VARIABLES (pointers), not VALUES (copies).
  In a loop with ONE variable, all closures share that one address.

GO VERSION BEHAVIOR
  Go ≤1.21:  range AND C-style → ONE variable → BUG
  Go 1.22+:  range → per-iteration (FIXED) ✅
             C-style → ONE variable (STILL BROKEN) ❌

THREE FIXES
  1. Shadow:    i := i                          (new var per iteration)
  2. Param:     go func(n int) { use(n) }(i)   (copy into parameter)
  3. Go 1.22+:  range loops fixed automatically (C-style: still need fix)

FAN-IN PATTERN
  Pattern: N source channels → 1 output channel
  Key rules:
    - wg.Add(N) BEFORE goroutine loop
    - Each goroutine: range source, forward to out, defer wg.Done()
    - Coordinator goroutine: wg.Wait() then close(out)
    - Pass channel as parameter: go func(s <-chan int){...}(src)

  Why coordinator?
    - Workers can't close out (other workers still sending → panic)
    - Consumer can't close out (doesn't know when producers are done)
    - Coordinator waits for ALL workers, then safely closes

DETECTION
  go vet ./...                     # loop variable capture warnings
  go test -race ./...              # concurrent access detection
  go build -gcflags='-m' ./...     # escape analysis (heap moves)
  staticcheck ./...                # SA6005 loop capture warning

CHANNEL DIRECTION TYPES
  chan int       → bidirectional (send and receive)
  chan<- int     → send-only (can only write to it)
  <-chan int     → receive-only (can only read from it)
  Direction is enforced at COMPILE TIME — zero runtime cost.
```

---

## One-Line Summary

> Loop closures capture the **address** of the loop variable, not its value — all
> closures see the final value unless you shadow (`i := i`) or copy via parameter.
> The fan-in pattern applies this fix to safely merge N channels through N goroutines
> coordinated by a WaitGroup that ensures exactly one `close()`.
