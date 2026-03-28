# Deep Dive: Concurrent Error Patterns — From Goroutines to errgroup

> This document builds up from goroutine/channel basics to production-grade
> concurrent error handling patterns. Read this BEFORE doing the goroutine
> and channel exercises if you want to understand the "why" behind each pattern.
>
> **Related:** [Chapter 08 — Error Chains & Wrapping](./08_error_chains_wrapping_strategy.md) covers
> sequential error handling: sentinel errors, custom types, `errors.Is`/`As`, and panic/recover.

---

## Table of Contents

1. [Goroutines — The 30-Second Refresher](#1-goroutines--the-30-second-refresher)
2. [The Problem: Goroutines Can't Return Errors](#2-the-problem-goroutines-cant-return-errors)
3. [Channels — The Communication Pipe](#3-channels--the-communication-pipe)
4. [sync.WaitGroup — "Wait for Everyone to Finish"](#4-syncwaitgroup--wait-for-everyone-to-finish)
5. [Pattern 1: WaitGroup + Mutex (Collect All Errors)](#5-pattern-1-waitgroup--mutex-collect-all-errors)
6. [Pattern 2: Channel-Based Error Collection](#6-pattern-2-channel-based-error-collection)
7. [Pattern 3: errgroup.Group (The Standard Pattern)](#7-pattern-3-errgroupgroup-the-standard-pattern)
8. [When to Use Which Pattern](#8-when-to-use-which-pattern)
9. [Common Mistakes and Gotchas](#9-common-mistakes-and-gotchas)

---

## 1. Goroutines — The 30-Second Refresher

A goroutine is a lightweight thread managed by the Go runtime (~2-8KB stack vs ~1MB OS thread).

```go
// Normal function call — blocks until done:
result := doWork()       // current goroutine waits here

// Goroutine — fires and continues immediately:
go doWork()              // launches new goroutine, current goroutine moves to next line
fmt.Println("I don't wait")   // runs immediately, doWork() runs concurrently
```

**The critical mental model:** `go doWork()` is fire-and-forget. The launching goroutine
does NOT wait. It does NOT get a return value. It just keeps going.

```
Timeline:
  main goroutine:   [launch doWork] → [print "I don't wait"] → [exit program]
  doWork goroutine: [starts]────────→ [doing work]──→ [maybe finishes... or maybe main exits first]
```

**Problem #1:** If `main()` exits, ALL goroutines are killed immediately. No cleanup.

---

## 2. The Problem: Goroutines Can't Return Errors

This is the fundamental issue that all patterns below solve.

```go
// Normal function — caller gets the error:
err := fetchURL("https://example.com")
if err != nil { /* handle it */ }

// Goroutine — WHERE does the error go?
go fetchURL("https://example.com")   // returns error to... nobody
// The error is created, returned, and silently garbage collected.
// The main goroutine has no idea anything went wrong.
```

Why? Because `go func()` is syntactic sugar for "launch this function on a new goroutine."
There's no wire connecting the new goroutine back to the launching one. They're independent.

**To get results (or errors) back from a goroutine, you need a communication mechanism.**
Go gives you two:
1. **Channels** — send values between goroutines (Go's preferred way)
2. **Shared memory + mutex** — write to a shared variable with lock protection

---

## 3. Channels — The Communication Pipe

A channel is a typed pipe between goroutines. One goroutine sends, another receives.

```go
ch := make(chan string)    // create a pipe that carries strings

// Goroutine A: sends a value INTO the pipe
go func() {
    ch <- "hello"          // send — blocks until someone receives
}()

// Goroutine B (main): reads a value FROM the pipe
msg := <-ch                // receive — blocks until someone sends
fmt.Println(msg)           // "hello"
```

**Visual:**

```
Goroutine A                    Channel                    Goroutine B (main)
                           ┌───────────┐
  ch <- "hello"  ────────► │  "hello"  │ ────────►  msg := <-ch
                           └───────────┘
  (blocks until             (the pipe)              (blocks until
   receiver ready)                                   sender sends)
```

### Buffered vs Unbuffered

```go
// Unbuffered — sender blocks until receiver is ready (synchronous handoff):
ch := make(chan string)

// Buffered — sender can put N items without blocking (like a mailbox with N slots):
ch := make(chan string, 5)    // 5 slots — can send 5 times without a receiver
```

### Closing a Channel

```go
close(ch)    // signals "no more values will be sent"

// Receivers can detect closure:
val, ok := <-ch    // ok=false means channel is closed and empty
// OR use range (reads until channel is closed):
for val := range ch {
    fmt.Println(val)    // processes each value, exits loop when ch is closed
}
```

### Error Channel — The Key Idea

Since channels carry typed values, they can carry `error`:

```go
errCh := make(chan error, 1)    // a pipe that carries errors

go func() {
    err := riskyOperation()
    errCh <- err               // send error (or nil) back to main
}()

err := <-errCh                 // main waits here, receives the error
if err != nil { /* handle */ }
```

**This is how goroutines communicate errors back.** The channel is the wire.

---

## 4. sync.WaitGroup — "Wait for Everyone to Finish"

WaitGroup solves Problem #1: "main exits before goroutines finish."

It's a counter:
- `wg.Add(1)` — increment: "one more goroutine to wait for"
- `wg.Done()` — decrement: "one goroutine finished"
- `wg.Wait()` — block until counter reaches 0

```go
var wg sync.WaitGroup

for i := 0; i < 3; i++ {
    wg.Add(1)                 // counter: 1, 2, 3
    go func(id int) {
        defer wg.Done()       // counter: 2, 1, 0 (when each finishes)
        fmt.Println("worker", id)
    }(i)
}

wg.Wait()                     // blocks here until counter = 0
fmt.Println("all done")       // only runs after ALL goroutines finish
```

**Visual timeline:**

```
main:      [Add 1][Add 1][Add 1]──────────────[Wait...]──────────────[all done]
worker 0:        [start]──[work]──[Done]                    ▲
worker 1:              [start]──[work]──────[Done]          │ counter hits 0
worker 2:                    [start]──[work]──[Done]────────┘
```

**Important:** `defer wg.Done()` is the first line inside the goroutine — it ensures
Done() runs even if the goroutine panics. Without defer, a panic skips Done(),
and `wg.Wait()` blocks forever (deadlock).

---

## 5. Pattern 1: WaitGroup + Mutex (Collect All Errors)

The simplest pattern: shared `[]error` slice protected by a mutex.

```go
func fetchAll(urls []string) error {
    var (
        wg   sync.WaitGroup
        mu   sync.Mutex       // protects the errs slice
        errs []error
    )

    for _, url := range urls {
        wg.Add(1)
        go func(u string) {
            defer wg.Done()
            if err := fetch(u); err != nil {
                mu.Lock()                                    // lock before writing
                errs = append(errs, fmt.Errorf("%s: %w", u, err))
                mu.Unlock()                                  // unlock after writing
            }
        }(url)
    }

    wg.Wait()                  // wait for ALL goroutines to finish
    return errors.Join(errs...)  // combine into single error (nil if errs is empty)
}
```

**Why the mutex?** Multiple goroutines run concurrently. If two goroutines `append` to
`errs` at the same time without a lock, you get a **data race** — corrupted memory,
wrong values, or a crash. The mutex ensures only one goroutine writes at a time.

**Step-by-step trace (3 URLs, URL #2 fails):**

```
main:       [Add 1][Add 1][Add 1]─────────────────────[Wait...]───[Join errs]
goroutine 1:      [fetch url1]──── OK ────[Done]            ▲
goroutine 2:            [fetch url2]── FAIL ──[Lock][append err][Unlock]──[Done]  │
goroutine 3:                  [fetch url3]──── OK ────[Done]────────────────┘
```

**Pros:** Simple, collects ALL errors.
**Cons:** All goroutines run to completion even if one fails — no cancellation.

---

## 6. Pattern 2: Channel-Based Error Collection

Instead of a shared slice + mutex, send errors through a channel.

```go
func fetchAll(ctx context.Context, urls []string) error {
    errCh := make(chan error, len(urls))    // buffered: won't block senders
    var wg sync.WaitGroup

    // Launch all goroutines
    for _, url := range urls {
        wg.Add(1)
        go func(u string) {
            defer wg.Done()
            if err := fetch(ctx, u); err != nil {
                errCh <- fmt.Errorf("%s: %w", u, err)    // send error into channel
            }
        }(url)
    }

    // Close channel AFTER all goroutines finish (in a separate goroutine)
    go func() {
        wg.Wait()        // wait for all workers
        close(errCh)     // signal: no more errors will be sent
    }()

    // Collect all errors from the channel
    var errs []error
    for err := range errCh {    // reads until channel is closed
        errs = append(errs, err)
    }
    return errors.Join(errs...)
}
```

**Let's break down the tricky parts:**

### Why `make(chan error, len(urls))`?

The buffer size `len(urls)` guarantees no goroutine will block on send. In the worst case,
ALL URLs fail — that's `len(urls)` errors. With a buffer that large, every goroutine
can send its error and call `Done()` without waiting for a receiver.

Without the buffer (or with a too-small buffer), a goroutine trying to send on a full
channel blocks — and if the receiver isn't reading yet, you get a deadlock.

### Why the separate `go func() { wg.Wait(); close(errCh) }()`?

This is the coordination trick. We need two things to happen:

1. **All goroutines must finish** before we close the channel
2. **The main goroutine must read from the channel** until it's closed

These can't happen on the same goroutine — the `for range errCh` loop blocks waiting
for the channel to close, and `wg.Wait()` blocks waiting for goroutines to finish.
If both are on main, main is stuck in `for range` and can never call `wg.Wait()`.

Solution: put `wg.Wait() + close()` on a separate goroutine:

```
main goroutine:           [launch workers]──[for range errCh... reading errors...]──[done]
closer goroutine:         [wg.Wait()... waiting for workers...]──[close(errCh)]──►  ▲
                                                                    signals main ───┘
worker goroutines:        [fetch]──[send err if failed]──[Done]
```

### Why `for err := range errCh`?

`range` over a channel reads values one at a time until the channel is **closed**.
It's syntactic sugar for:

```go
for {
    err, ok := <-errCh
    if !ok { break }           // channel closed, no more values
    errs = append(errs, err)
}
```

**Pros:** No mutex needed — channels handle synchronization. More "Go-idiomatic."
**Cons:** More moving parts (the closer goroutine). Still no cancellation on first error.

---

## 7. Pattern 3: errgroup.Group (The Standard Pattern)

`errgroup` (from `golang.org/x/sync/errgroup`) combines WaitGroup + error collection +
context cancellation into one clean API. **This is what you use in production.**

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) error {
    g, ctx := errgroup.WithContext(ctx)    // creates group + derived context

    for _, url := range urls {
        g.Go(func() error {                // launches a goroutine
            return fetch(ctx, url)         // return error like a normal function!
        })
    }

    return g.Wait()    // waits for ALL goroutines, returns FIRST error
}
```

**That's it.** Three lines replace all the WaitGroup/mutex/channel boilerplate.

### What `errgroup` Does Under the Hood

`errgroup.Group` is roughly this (simplified):

```go
type Group struct {
    wg      sync.WaitGroup
    err     error              // stores the FIRST error
    once    sync.Once          // ensures only first error is captured
    cancel  context.CancelFunc // cancels context on first error
}

func (g *Group) Go(f func() error) {
    g.wg.Add(1)
    go func() {
        defer g.wg.Done()
        if err := f(); err != nil {
            g.once.Do(func() {     // only runs ONCE — captures first error
                g.err = err
                if g.cancel != nil {
                    g.cancel()     // cancels the context — signals other goroutines
                }
            })
        }
    }()
}

func (g *Group) Wait() error {
    g.wg.Wait()       // waits for ALL goroutines to finish
    return g.err      // returns the first error (or nil)
}
```

### Step-by-Step Trace

```go
g, ctx := errgroup.WithContext(ctx)
// g = new Group with cancel function
// ctx = child context that will be cancelled on first error

g.Go(func() error { return fetch(ctx, "url1") })    // launches goroutine 1
g.Go(func() error { return fetch(ctx, "url2") })    // launches goroutine 2
g.Go(func() error { return fetch(ctx, "url3") })    // launches goroutine 3

return g.Wait()    // blocks until all three finish
```

**Scenario: url2 fails first:**

```
goroutine 1: [fetch url1]────────────────[ctx cancelled! → returns early]──[Done]
goroutine 2: [fetch url2]──FAIL──[once.Do: save err, cancel ctx]──[Done]
goroutine 3: [fetch url3]──────────[ctx cancelled! → returns early]──[Done]
                                        │
g.Wait():   [...waiting...]─────────────┴──────────[all Done]──returns first error
```

**Key behavior:**
1. `g.Go()` launches a goroutine — you return `error` like a normal function
2. **First error** is captured via `sync.Once` — subsequent errors are discarded
3. **Context is cancelled** on first error — other goroutines can detect this via `ctx.Done()`
4. `g.Wait()` blocks until ALL goroutines finish (not just until first error)
5. Returns the first error, or nil if all succeeded

### Why Context Cancellation Matters

When url2 fails and cancels the context, goroutines 1 and 3 can detect this IF they
check `ctx`:

```go
func fetch(ctx context.Context, url string) error {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    // When ctx is cancelled, this request is immediately aborted
    // The goroutine doesn't waste time waiting for a response
    resp, err := http.DefaultClient.Do(req)
    if err != nil { return err }  // returns context.Canceled
    // ...
}
```

Without context-aware operations, the goroutines would run to completion even though
we already know the overall operation failed. Context cancellation enables **fast failure**.

---

## 8. When to Use Which Pattern

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Question                              │ Pattern                            │
├───────────────────────────────────────┼────────────────────────────────────┤
│ Need ALL errors?                      │ WaitGroup + Mutex + errors.Join   │
│                                       │ OR Channel-based collection       │
├───────────────────────────────────────┼────────────────────────────────────┤
│ Need FIRST error + cancel remaining?  │ errgroup.WithContext              │
├───────────────────────────────────────┼────────────────────────────────────┤
│ Need ALL errors + cancellation?       │ Custom: channel + context         │
├───────────────────────────────────────┼────────────────────────────────────┤
│ Simple fire-and-forget?               │ WaitGroup alone (no errors)       │
├───────────────────────────────────────┼────────────────────────────────────┤
│ Production service?                   │ errgroup (95% of the time)        │
└───────────────────────────────────────┴────────────────────────────────────┘
```

**In practice, `errgroup` covers 95% of concurrent error handling needs.**
Use the channel pattern when you genuinely need every error collected.

---

## 9. Common Mistakes and Gotchas

### Mistake 1: Forgetting to close the channel

```go
// BUG: for range errCh blocks FOREVER — nobody closes the channel
for err := range errCh {
    errs = append(errs, err)
}
```

### Mistake 2: Closing channel from multiple goroutines

```go
go func() {
    close(errCh)    // PANIC if two goroutines close the same channel
}()
```

Only ONE goroutine should close a channel — typically the "closer" goroutine
that waits for all senders to finish.

### Mistake 3: Data race on shared slice without mutex

```go
// BUG: multiple goroutines appending concurrently — data race
go func() { errs = append(errs, err1) }()
go func() { errs = append(errs, err2) }()
// Run with: go test -race → will flag this
```

### Mistake 4: Not using context in goroutines

```go
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error {
    return fetch(url)        // ❌ ignores ctx — won't cancel on first error
    return fetch(ctx, url)   // ✅ respects ctx — cancels when group fails
})
```

### Mistake 5: Panic in a goroutine kills the ENTIRE program

```go
go func() {
    panic("oops")    // kills main() and ALL other goroutines — no recovery
}()
```

A panic in a goroutine is NOT caught by `recover()` in the launching goroutine.
Each goroutine needs its own `defer/recover` if panics are possible:

```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("goroutine panic: %v", r)
        }
    }()
    riskyCode()
}()
```

---

## One-Line Summary

> Goroutines can't return errors — use channels or shared memory to communicate them back,
> `sync.WaitGroup` to wait for completion, and `errgroup.Group` to combine both with
> automatic context cancellation in production code.
