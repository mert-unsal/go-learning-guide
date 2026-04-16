# 📦 Module 10 — Goroutines

> **Topics covered:** `go` keyword · `sync.WaitGroup` · `sync.Mutex` · `sync.Once` · Race conditions · Concurrent patterns

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
| Launching goroutines with `go` | Exercise 1 — `RunConcurrently` |
| `sync.WaitGroup` to wait for goroutines | Exercise 1, 3 |
| `sync.Mutex` to protect shared state | Exercise 2 — `ExCounter` |
| Race conditions and why they're dangerous | `concepts.go` |
| Splitting work across goroutines | Exercise 3 — `SumConcurrent` |
| `sync.Once` for one-time initialization | Exercise 4 — `RunOnce` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `RunConcurrently(n int, fn func(id int))` | Launch n goroutines, wait for all with WaitGroup |
| 2a | `(c *ExCounter) Inc()` | Thread-safe increment using Mutex |
| 2b | `(c *ExCounter) Value() int` | Thread-safe read using Mutex |
| 3 | `SumConcurrent(nums []int) int` | Sum two halves in separate goroutines |
| 4 | `RunOnce(setup func())` | Call setup exactly once using `sync.Once` |

---

## 🧪 Run Tests

> ⚠️ The `./exercises/fundamentals/...` paths work from the **project root** only.  
> If you are inside this folder, use `go test . -v` instead.

### From project root:
```bash
go test ./exercises/fundamentals/10_goroutines/... -v
```

### From inside this folder:
```bash
go test . -v
```

### Run with race detector (highly recommended for concurrency!):
```bash
go test . -v -race
```

### Run a single test (from inside this folder):
```bash
go test . -v -run TestRunConcurrently
go test . -v -run TestExCounter
go test . -v -run TestSumConcurrent
go test . -v -run TestRunOnce
```

---

## 💡 Key Hints

<details>
<summary>Exercise 1 — WaitGroup pattern hint</summary>

```go
var wg sync.WaitGroup
for i := 0; i < n; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fn(id)
    }(i)   // ← pass i as argument! Don't capture directly.
}
wg.Wait()
```
⚠️ Notice `go func(id int) {...}(i)` — always pass loop variables as arguments to avoid the classic goroutine closure bug.
</details>

<details>
<summary>Exercise 2 — Mutex pattern hint</summary>

```go
func (c *ExCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```
`defer Unlock()` ensures the lock is always released, even if the function panics.
</details>

<details>
<summary>Always run -race flag for concurrency code</summary>

The Go race detector catches data races at runtime:
```bash
go test ./exercises/fundamentals/10_goroutines/... -race
```
If it reports a race, your Mutex is missing somewhere.
</details>

---

## ✅ Done? Next Step

```bash
go test ./exercises/fundamentals/11_channels/... -v
```

---

## 📖 Companion Chapter

For the deep-dive theory behind these exercises, read:

- [10 — Goroutines & the GMP Scheduler](../../../learnings/10_goroutines_gmp_scheduler.md) — G/M/P model, goroutine states, work stealing, preemption, `GODEBUG=schedtrace`

