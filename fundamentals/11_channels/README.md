# 📦 Module 11 — Channels

> **Topics covered:** Buffered/unbuffered channels · Directional channels · `select` · Pipeline pattern · Fan-in · Timeout with `time.After` · Nil-channel select · Semaphore · Close drain semantics · Non-blocking ops · Goroutine leak detection · OrDone pattern · Channels vs Atomics

---

## 🗺️ Learning Path

```
1. Read concepts.go        ← Theory with runnable examples
2. Open exercises.go       ← Implement the TODO functions yourself
3. Run the tests below     ← Instant feedback on your code
4. Stuck? Open solutions.go ← Only after you have tried!
```

---

## 📚 What You Will Learn

| Concept | Where |
|---------|-------|
| Creating and using channels | Exercise 1 — `SumAsync` |
| Send-only `chan<-` / receive-only `<-chan` | Exercise 2 — `Generate` |
| Pipeline pattern (chain of channels) | Exercise 3 — `Square` |
| Fan-in (merge multiple channels) | Exercise 4 — `Merge` |
| `select` statement | `concepts.go` + Exercise 5 |
| Timeout with `time.After` | Exercise 5 — `WithTimeout` |
| Closing channels and `range` over channel | Exercise 2, 3 |
| Nil-channel select pattern (dynamic disable) | Exercise 6 — `MergeN` |
| Semaphore (bounded concurrency) | Exercise 7 — `ProcessWithLimit` |
| Close drain semantics (buffer + closed) | Exercise 8 — `SendAndClose` |
| Non-blocking try-send / try-receive | Exercise 9 — `TrySend` / `TryReceive` |
| Goroutine leak detection with context | Exercise 10 — `SafeGenerator` |
| OrDone pattern (cancellation-aware wrapper) | Exercise 11 — `OrDone` |
| Channels vs Atomics (when NOT to use channels) | Bonus — `ChannelCounter` / `AtomicCounter` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `SumAsync(nums []int, ch chan<- int)` | Compute sum and send on channel |
| 2 | `Generate(n int) <-chan int` | Send 1..n on a channel, then close it |
| 3 | `Square(in <-chan int) <-chan int` | Square each value from input channel |
| 4 | `Merge(a, b <-chan int) <-chan int` | Fan-in: merge two channels into one |
| 5 | `WithTimeout(ch <-chan int, maxWaitMs int) (int, bool)` | Receive with a timeout |
| 6 | `MergeN(channels ...<-chan int) <-chan int` | Merge N channels using nil-channel disable |
| 7 | `ProcessWithLimit(items, max, fn)` | Semaphore: bounded concurrent processing |
| 8 | `SendAndClose(values, bufSize)` | Close semantics: buffer drain verification |
| 9 | `TrySend` / `TryReceive` | Non-blocking channel ops with select+default |
| 10 | `SafeGenerator(ctx) <-chan int` | Context-aware producer (no goroutine leak) |
| 11 | `OrDone(ctx, in) <-chan int` | Cancellation-aware channel wrapper |
| B | `ChannelCounter` / `AtomicCounter` | When channels are the wrong tool |

---

## 🧪 Run Tests

> ⚠️ The `./fundamentals/...` paths work from the **project root** only.  
> If you are inside this folder, use `go test . -v` instead.

### From project root:
```bash
go test ./fundamentals/11_channels/... -v
```

### From inside this folder:
```bash
go test . -v
```

### Run with race detector:
```bash
go test . -v -race
```

### Run a single test (from inside this folder):
```bash
go test . -v -run TestSumAsync
go test . -v -run TestGenerate
go test . -v -run TestSquare
go test . -v -run TestMerge
go test . -v -run TestWithTimeout
go test . -v -run TestMergeN
go test . -v -run TestProcessWithLimit
go test . -v -run TestSendAndClose
go test . -v -run TestTrySend
go test . -v -run TestTryReceive
go test . -v -run TestSafeGenerator
go test . -v -run TestOrDone
go test . -v -run TestChannelCounter
go test . -v -run TestAtomicCounter
```

---

## 💡 Key Hints

<details>
<summary>Exercise 6 — MergeN: nil channel select pattern</summary>

```go
// Setting a channel to nil in a select makes that case block forever.
// Use this to "disable" closed channels:
for i, ch := range channels {
    select {
    case v, ok := <-ch:
        if !ok {
            channels[i] = nil  // disable this case
            alive--
        } else {
            out <- v
        }
    default:
    }
}
```
This is the Go idiom for dynamically controlling select behavior without WaitGroup.
</details>

<details>
<summary>Exercise 7 — Semaphore: buffered channel as concurrency limiter</summary>

```go
sem := make(chan struct{}, maxConcurrent)
sem <- struct{}{}         // acquire — blocks when limit reached
defer func() { <-sem }()  // release — frees slot
```
Buffer capacity = max concurrent goroutines. Simple and elegant.
</details>

<details>
<summary>Exercise 9 — TrySend/TryReceive: select with default</summary>

```go
select {
case ch <- val:
    return true   // send succeeded
default:
    return false  // channel not ready — don't block
}
```
Under the hood: selectgo() polls once, sees default, returns immediately.
</details>

<details>
<summary>Exercise 10 — SafeGenerator: context cancellation</summary>

```go
select {
case <-ctx.Done():
    return           // context cancelled — exit goroutine
case ch <- value:
    // sent successfully
}
```
Always check `ctx.Done()` in select alongside channel ops. This prevents goroutine leaks.
</details>

<details>
<summary>Exercise 11 — OrDone: the double-select pattern</summary>

```go
// Outer select: receive from in or notice cancellation
case v, ok := <-in:
    // Inner select: forward to out or notice cancellation
    select {
    case out <- v:
    case <-ctx.Done():
        return
    }
```
Why inner select? Without it, blocking on `out <- v` ignores cancellation.
</details>

<details>
<summary>Channels vs Mutexes — when to use which?</summary>

| Use **channels** | Use **mutexes** |
|-----------------|-----------------|
| Passing ownership of data | Protecting shared state |
| Pipeline / producer-consumer | Simple counters/caches |
| Fan-out / fan-in patterns | Struct field protection |

Go proverb: *"Don't communicate by sharing memory; share memory by communicating."*

But also: *"Channels orchestrate; mutexes serialize."*  
For simple counters, `sync/atomic` is 5-10x faster than channels.
</details>

---

## 📖 Deep Dive

For runtime internals (hchan, sudog, selectgo algorithm, performance costs):  
[`learnings/07_channels_internals.md`](../../learnings/07_channels_internals.md)

---

## ✅ Done? Next Step

```bash
go test ./fundamentals/12_packages_modules/... -v
```

