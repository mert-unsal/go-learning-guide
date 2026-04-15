# 03 Concurrency Patterns — Production Patterns

> **Companion chapters:**
> - [learnings/10_goroutines_gmp_scheduler.md](../../../learnings/10_goroutines_gmp_scheduler.md)
> - [learnings/12_channels_hchan_select.md](../../../learnings/12_channels_hchan_select.md)
> - [learnings/25_sync_primitives_under_the_hood.md](../../../learnings/25_sync_primitives_under_the_hood.md)

## Exercises

| # | Function | Pattern | Difficulty |
|---|---------|---------|------------|
| 1 | `FanOut` | Run N tasks, collect results | ⭐⭐ |
| 2 | `Pipeline` | Stage → Stage → Stage | ⭐⭐⭐ |
| 3 | `WorkerPool` | Bounded concurrency | ⭐⭐⭐ |
| 4 | `Semaphore` | Buffered channel as semaphore | ⭐⭐⭐ |
| 5 | `WithTimeout` | Cancel slow work | ⭐⭐ |
| 6 | `Merge` | Fan-in N channels | ⭐⭐⭐ |
| 7 | `Generate` | Lazy value production | ⭐⭐ |
| 8 | `OrDone` | Channel + context cancellation | ⭐⭐⭐ |
| 9 | `RateLimited` | Token bucket rate limiting | ⭐⭐⭐ |
| 10 | `SafeMap` | RWMutex-protected map | ⭐⭐ |
| 11 | `Barrier` | Wait for N, then proceed | ⭐⭐⭐ |
| 12 | `GracefulWorker` | Drain work, then exit | ⭐⭐⭐ |

## How to Practice

```bash
go test -race -v ./exercises/advanced/03_concurrency_patterns/
go test -race -run TestWorkerPool ./exercises/advanced/03_concurrency_patterns/
```

## Key Insights

- **Always use `-race`** — race conditions are silent killers
- **Fan-out**: use `sync.WaitGroup` + goroutines, index into result slice
- **Pipeline**: each stage is a goroutine reading from input channel, writing to output
- **Worker pool**: bounded by channel buffer or semaphore
- **Graceful shutdown**: drain remaining work before exiting
