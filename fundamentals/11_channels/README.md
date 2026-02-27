# 📦 Module 11 — Channels

> **Topics covered:** Buffered/unbuffered channels · Directional channels · `select` · Pipeline pattern · Fan-in · Timeout with `time.After`

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
```

---

## 💡 Key Hints

<details>
<summary>Exercise 2 — Generator pattern hint</summary>

```go
func Generate(n int) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)   // always close when done!
        for i := 1; i <= n; i++ {
            ch <- i
        }
    }()
    return ch
}
```
The goroutine runs in the background; the caller ranges over the returned channel.
</details>

<details>
<summary>Exercise 5 — WithTimeout select hint</summary>

```go
select {
case v := <-ch:
    return v, true
case <-time.After(time.Duration(maxWaitMs) * time.Millisecond):
    return 0, false
}
```
`select` picks whichever case is ready first — either a value arrives, or the timeout fires.
</details>

<details>
<summary>Channels vs Mutexes — when to use which?</summary>

| Use **channels** | Use **mutexes** |
|-----------------|-----------------|
| Passing ownership of data | Protecting shared state |
| Pipeline / producer-consumer | Simple counters/caches |
| Fan-out / fan-in patterns | Struct field protection |

Go proverb: *"Don't communicate by sharing memory; share memory by communicating."*
</details>

---

## ✅ Done? Next Step

```bash
go test ./fundamentals/12_packages_modules/... -v
```

