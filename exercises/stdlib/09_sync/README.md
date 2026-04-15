# 📦 Module 09 — sync: Mutexes, Pools, Once & Concurrency Primitives

> **Topics covered:** sync.Mutex · sync.RWMutex · sync.WaitGroup · sync.Once · sync.Pool · sync.Map · sync.Cond · atomic.Value · singleflight pattern
>
> **Deep dive:** [Chapter 25 — sync Primitives Under the Hood](../../../learnings/25_sync_primitives_under_the_hood.md)

---

## 🗺️ Learning Path

```
1. Read: Chapter 25 — sync internals (Mutex starvation, Pool victim cache)
2. Open exercises.go                                    ← Implement 12 exercises
3. Run go test -race -v ./...                           ← Make them all pass
```

---

## 📚 What You Will Learn

| Concept | Exercise | Key Insight |
|---------|----------|-------------|
| `sync.Mutex` | Ex 1 | **Exclusive lock** — protect shared map |
| `sync.RWMutex` | Ex 2 | **Shared reads, exclusive writes** |
| `sync.WaitGroup` | Ex 3 | **Fan-out and wait** for goroutines |
| `sync.Once` | Ex 4 | **Exactly-once initialization** |
| `sync.Pool` | Ex 5 | **Buffer reuse** to reduce GC pressure |
| Mutex for aggregation | Ex 6 | **Concurrent sum** with chunked work |
| `sync.Map` | Ex 7 | **Concurrent map** for disjoint keys |
| `sync.OnceValue` (1.21+) | Ex 8 | **Lazy singleton** with cached result |
| Channel-based semaphore | Ex 9 | **TryLock with timeout** pattern |
| `sync.Cond` broadcast | Ex 10 | **Gate pattern** — wake all waiters |
| `atomic.Value` / `sync.Map` | Ex 11 | **Lock-free config reload** |
| Singleflight pattern | Ex 12 | **Deduplicate concurrent calls** |

---

## ✏️ Exercises

| # | Function / Type | What to implement |
|---|----------------|------------------|
| 1 | `SafeCounter` | Mutex-protected map counter |
| 2 | `RWCache` | RWMutex read-heavy cache |
| 3 | `WaitForAll(n, fn)` | WaitGroup fan-out |
| 4 | `Service.Init(dsn)` | sync.Once initialization |
| 5 | `GetBuffer / PutBuffer` | sync.Pool buffer reuse |
| 6 | `ConcurrentSum(nums, workers)` | Mutex aggregation |
| 7 | `SyncMapStore(pairs)` | sync.Map concurrent store |
| 8 | `MakeLazy(factory)` | sync.OnceValue lazy singleton |
| 9 | `TimedMutex.TryLock(timeout)` | Channel semaphore |
| 10 | `Gate.Wait / Open` | sync.Cond broadcast |
| 11 | `AtomicConfig.Update / Current` | Lock-free config |
| 12 | `SingleFlight.Do(key, fn)` | Deduplicate calls |

---

## 🧪 Run Tests

```bash
go test -race -v -timeout=30s ./exercises/stdlib/09_sync/
```

---

## ✅ Done? Next Step

Explore advanced concurrency patterns:
```bash
go test -race -v ./exercises/stdlib/08_context/
```
