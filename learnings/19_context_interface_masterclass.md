# 12 — Context Interface Deep Dive

## Why `context.Context` Is the Best Example of Interfaces in Go

> This document explains what `context.Context` is, why it's an interface (not a struct),
> and how its design embodies Go's interface philosophy.

---

## 1. The Interface — Just 4 Methods

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

That's it. Four methods. No fields, no inheritance, no framework. Any type that has these four methods **is** a `Context` — the Go compiler guarantees it.

| Method       | Purpose                                                       |
|-------------|---------------------------------------------------------------|
| `Deadline()` | Returns the time when this context will be cancelled (if any) |
| `Done()`     | Returns a channel that closes when the context is cancelled   |
| `Err()`      | Returns **why** the context was cancelled (`Canceled` or `DeadlineExceeded`) |
| `Value(key)` | Retrieves a request-scoped value by key                       |

---

## 2. The Concrete Implementations (What the stdlib Actually Creates)

The `context` package has **four unexported structs**. Each one does exactly **one thing**:

### `emptyCtx` — The Root (Does Nothing)

```go
type emptyCtx struct{}

func (emptyCtx) Deadline() (time.Time, bool) { return time.Time{}, false }
func (emptyCtx) Done() <-chan struct{}        { return nil }
func (emptyCtx) Err() error                  { return nil }
func (emptyCtx) Value(key any) any           { return nil }
```

This is what `context.Background()` and `context.TODO()` return. It never cancels, has no deadline, stores no values. It's the **zero-value anchor** of the chain.

### `cancelCtx` — Adds Cancellation

```go
type cancelCtx struct {
    Context              // ← embeds the INTERFACE (parent)
    done    chan struct{} // closed on cancel
    err     error        // set on cancel
}
```

It **embeds the parent as an interface**, not as a concrete type. It only adds a `done` channel and an `err`. Everything else (`Deadline()`, `Value()`) is delegated to the parent via the embedded interface.

### `timerCtx` — Adds Deadline/Timeout

```go
type timerCtx struct {
    *cancelCtx           // embeds cancelCtx (inherits cancellation)
    timer    *time.Timer // fires at deadline
    deadline time.Time
}
```

Wraps a `cancelCtx`, adds a timer. Only overrides `Deadline()` — cancellation behavior comes from the embedded `cancelCtx`, and `Value()` still delegates up the chain.

### `valueCtx` — Adds One Key-Value Pair

```go
type valueCtx struct {
    Context        // ← embeds the INTERFACE (parent)
    key, val any   // ONE key-value pair, not a map
}

func (c *valueCtx) Value(key any) any {
    if c.key == key {
        return c.val  // found it
    }
    return c.Context.Value(key)  // not mine, ask parent
}
```

This is critical: each `WithValue()` call creates a **new node** holding exactly one key-value pair. It doesn't add to a map — it creates a linked list of single-value wrappers.

---

## 3. How the Chain Works — Step by Step

```go
ctx := context.Background()                        // emptyCtx{}
ctx  = context.WithValue(ctx, "userID", 42)         // valueCtx{parent: emptyCtx, key: "userID", val: 42}
ctx  = context.WithCancel(ctx)                      // cancelCtx{parent: valueCtx{...}}
ctx  = context.WithValue(ctx, "traceID", "abc-123") // valueCtx{parent: cancelCtx{...}, key: "traceID", val: "abc-123"}
```

The chain in memory looks like:

```
valueCtx("traceID"="abc-123")
  │
  └──▶ cancelCtx (has done channel)
         │
         └──▶ valueCtx("userID"=42)
                │
                └──▶ emptyCtx (the root)
```

### What Happens When You Call `ctx.Value("userID")`?

```
Step 1: valueCtx("traceID") → key is "traceID", not "userID" → ask parent
Step 2: cancelCtx           → doesn't override Value() → ask parent
Step 3: valueCtx("userID")  → key matches! → return 42
```

### What Happens When You Call `cancel()`?

The `cancelCtx` closes its `done` channel. Any goroutine doing this immediately unblocks:

```go
select {
case <-ctx.Done():
    return ctx.Err()  // context.Canceled
case result := <-workCh:
    return result
}
```

Cancellation propagates **downward** — children of the cancelled context also get cancelled.
It does **not** propagate upward — the parent remains alive.

---

## 4. What If Context Were a Struct?

This is where the interface design becomes obvious.

### The Struct Version (Hypothetical)

```go
type Context struct {
    deadline  time.Time
    hasDL     bool
    done      chan struct{}
    err       error
    values    map[any]any
    parent    *Context
    cancelFn  func()
    timer     *time.Timer
}
```

### Problems with the Struct Approach

#### Problem 1: Every Context Pays for Every Feature

```go
ctx := context.Background()
// With struct: allocates deadline, timer, values map, done channel...
//              even though Background() needs NONE of them

// With interface: emptyCtx is literally struct{} — zero bytes of useful data
```

In a high-throughput server handling 50k requests/sec, each request creates a context chain.
The struct approach wastes memory on unused fields across millions of allocations.

#### Problem 2: You Can't Extend It

What if a library needs a context with custom behavior — say, a context that logs every
cancellation? With the interface:

```go
type loggingCtx struct {
    context.Context              // embed the interface
    logger *slog.Logger
}

func (c *loggingCtx) Done() <-chan struct{} {
    ch := c.Context.Done()
    c.logger.Info("someone checked Done()")
    return ch
}
```

Done. It satisfies `context.Context`. It wraps any existing context. It works with
every function in the stdlib that accepts `context.Context`.

With a struct? **Impossible.** You can't override methods on a struct you don't own.
You'd need the `context` package to anticipate every possible extension — which violates
Go's design philosophy entirely.

#### Problem 3: Single Responsibility Collapses

The interface design gives each wrapper **one job**:

| Type         | Responsibility           | Fields Added        |
|-------------|--------------------------|---------------------|
| `emptyCtx`  | Be the root              | none                |
| `cancelCtx` | Cancellation             | `done`, `err`       |
| `timerCtx`  | Deadline-based cancel    | `timer`, `deadline` |
| `valueCtx`  | One key-value pair       | `key`, `val`        |

With a struct, all responsibilities collapse into one type. Adding a new feature
(like `context.WithoutCancel` added in Go 1.21) means modifying the struct,
adding fields, adding branching logic. With the interface, it's just a new wrapper type.

#### Problem 4: Testing Becomes Rigid

With the interface, you can create test-specific contexts:

```go
type alwaysCancelledCtx struct{}

func (alwaysCancelledCtx) Deadline() (time.Time, bool)  { return time.Time{}, false }
func (alwaysCancelledCtx) Done() <-chan struct{}         { ch := make(chan struct{}); close(ch); return ch }
func (alwaysCancelledCtx) Err() error                   { return context.Canceled }
func (alwaysCancelledCtx) Value(key any) any             { return nil }
```

Pass this to any function that takes `context.Context` — it immediately triggers the
cancellation path. With a struct, you'd need to create a real context chain and actually
cancel it, or the struct would need test-only fields.

---

## 5. The Design Pattern: Decorator / Wrapper

What `context.Context` demonstrates is the **decorator pattern** in Go:

```
┌─────────────────────────────────────────────────────┐
│  Each wrapper:                                       │
│    1. Holds the parent as an INTERFACE               │
│    2. Overrides only the methods it cares about      │
│    3. Delegates everything else to the parent        │
│    4. Is itself an implementation of the interface    │
│       → so it can be wrapped by the NEXT layer       │
└─────────────────────────────────────────────────────┘
```

This pattern appears everywhere in Go's stdlib:

```go
// io.Reader chain — same decorator pattern
file, _ := os.Open("data.gz")         // *os.File       → io.Reader
buf := bufio.NewReader(file)           // bufio.Reader   → wraps io.Reader
gz, _ := gzip.NewReader(buf)           // gzip.Reader    → wraps io.Reader
limited := io.LimitReader(gz, 1024)    // LimitedReader  → wraps io.Reader

// http.Handler chain — same pattern
handler := myHandler{}                                    // implements http.Handler
handler = loggingMiddleware(handler)                       // wraps http.Handler
handler = authMiddleware(handler)                          // wraps http.Handler
handler = rateLimitMiddleware(handler)                     // wraps http.Handler
```

In every case, the interface is what makes the stacking possible. Each wrapper
doesn't know or care what's inside it — it only knows the interface contract.

---

## 6. Summary — Why Context Must Be an Interface

| Concern                    | Struct                                      | Interface ✅                                     |
|---------------------------|---------------------------------------------|--------------------------------------------------|
| **Memory efficiency**      | All fields always allocated                 | Each wrapper carries only what it needs          |
| **Extensibility**          | Modify the one struct (breaking change)     | Add a new wrapper type (non-breaking)            |
| **Single responsibility**  | One type does everything                    | Each type does one thing                         |
| **Third-party extension**  | Impossible without forking                  | Implement 4 methods, done                        |
| **Testability**            | Must use real context chain                 | Create a fake in 10 lines                        |
| **Composition**            | Flat, rigid                                 | Stackable, infinitely composable                 |

### The Go Philosophy This Embodies

1. **"The bigger the interface, the weaker the abstraction"** — 4 methods is all you need
2. **"Accept interfaces, return concrete types"** — `WithCancel` accepts `Context` (interface), returns `Context` (interface wrapping a concrete `cancelCtx`)
3. **Consumer defines the contract** — any package can create a custom `Context` implementation
4. **Composition over inheritance** — behavior is built by stacking wrappers, not by subclassing

---

## 7. Production Gotchas

### Value Lookup is O(n)

Each `Value()` call walks the linked list. Deep chains = slow lookups:

```go
// 100 WithValue calls = 100-node linked list
// Value() on the outermost walks up to 100 nodes
```

**Rule:** Keep context values shallow. Use them for request-scoped metadata
(trace IDs, request IDs), **not** for dependency injection.

### Cancellation is One-Way (Parent → Children)

```
parent ──cancel──▶ child1 ✓
                   child2 ✓

child1 ──cancel──▶ parent ✗  (parent is unaffected)
```

### Don't Store Context in Structs

```go
// ❌ Anti-pattern
type Service struct {
    ctx context.Context  // which request does this belong to?
}

// ✅ Pass as first parameter
func (s *Service) Process(ctx context.Context, data []byte) error { ... }
```

Context is **request-scoped** and short-lived. Struct fields are long-lived. Mixing them
creates confusing lifecycles and potential goroutine leaks.

### Always Check `ctx.Err()` Before Expensive Operations

```go
func process(ctx context.Context, items []Item) error {
    for _, item := range items {
        if err := ctx.Err(); err != nil {
            return err  // bail early if cancelled
        }
        // ... expensive work ...
    }
    return nil
}
```

---

## 8. Mental Model — The Final Picture

Think of `context.Context` as a **Russian nesting doll** (matryoshka):

```
┌──────────────────────────────────────────────┐
│  valueCtx("traceID" = "abc")                 │
│  ┌────────────────────────────────────────┐   │
│  │  cancelCtx (done channel)              │   │
│  │  ┌──────────────────────────────────┐  │   │
│  │  │  valueCtx("userID" = 42)         │  │   │
│  │  │  ┌────────────────────────────┐   │  │   │
│  │  │  │  emptyCtx (Background)     │   │  │   │
│  │  │  └────────────────────────────┘   │  │   │
│  │  └──────────────────────────────────┘  │   │
│  └────────────────────────────────────────┘   │
└──────────────────────────────────────────────┘
```

Each doll:
- Adds exactly one capability
- Delegates everything else inward
- Only knows the **interface** of the doll inside it
- **Is itself** a valid `Context` that can be put inside another doll

This is only possible because the "shape" each doll must match is an **interface**, not a
specific struct. That's the entire lesson.

---

### Go 1.21+ Additions

**`context.WithoutCancel(parent)`** — creates a child that is NOT cancelled when parent is cancelled. Use case: cleanup operations that must complete even after request cancellation (e.g., audit logging, releasing resources).

**`context.AfterFunc(ctx, func())`** — registers a function to run after ctx is done. Safer than spawning a goroutine with `go func() { <-ctx.Done(); cleanup() }()` because AfterFunc handles the goroutine lifecycle and returns a `stop` function to deregister the callback if it's no longer needed.
