# 📚 Go Internals Deep Dive Series

Runtime-level explanations of Go's core systems — what the compiler and runtime
actually do under the hood. Every document includes ASCII diagrams, `runtime/` source
references, performance cost tables, and production implications.

---

## Reading Order

The documents are organized in a deliberate learning progression:

### Phase 1: Data Structures (how Go stores things in memory)

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 01 | [Slices Internals](./01_slices_internals.md) | Slice header (3-word struct), `runtime.growslice`, append algorithm, backing array sharing, memory leaks |
| 02 | [Maps Internals](./02_maps_internals.md) | `runtime.hmap`, buckets, tophash, load factor, incremental evacuation, concurrent access fatals |

### Phase 2: Language Core (how Go's type system and semantics work)

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 03 | [Closures & Scopes](./03_closures_and_scopes.md) | Block scoping, `funcval` struct, capture-by-reference, escape analysis, loop gotchas, Go 1.22 changes |
| 04 | [Pointers & Auto-Deref](./04_pointers_and_auto_deref.md) | Auto-dereference, auto-address, addressability rules, when Go silently inserts `*` and `&` |
| 05 | [Interfaces Internals](./05_interfaces_internals.md) | `iface`/`eface`, `itab` dispatch table, nil trap, three guards, method sets, performance costs |
| 06 | [Error Handling Patterns](./06_error_handling_patterns.md) | Sentinel errors, custom types, `%w` wrapping, `errors.Is/As` chain walk, enterprise error strategy |

### Phase 3: Concurrency (how Go runs things in parallel)

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 07 | [Channels Internals](./07_channels_internals.md) | `runtime.hchan`, ring buffer, `sudog` parking, select algorithm, nil channel semantics |
| 08 | [Goroutines & Scheduler](./08_goroutines_and_scheduler.md) | GMP model, work stealing, async preemption, syscall hand-off, network poller, stack management |

### Phase 4: Runtime (how Go manages memory and performance)

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 09 | [Memory, GC & Escape Analysis](./09_memory_gc_and_escape_analysis.md) | Stack vs heap, escape analysis rules, tri-color GC, write barrier, GOGC/GOMEMLIMIT, sync.Pool |
| 10 | [Sorting — pdqsort](./10_sorting_pdqsort.md) | Pattern-defeating quicksort: how `slices.Sort` adaptively picks insertion sort, quicksort, or heapsort |

### Phase 5: Production (how to ship Go to production)

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 11 | [Production Go: Pitfalls & Choices](./11_production_go_pitfalls.md) | Top 15 production bugs, library comparison tables, project structure, graceful shutdown, Docker builds |

---

## How to Use These Documents

1. **Interview prep**: Read Phase 1-3. These cover the most common Go interview questions at senior level.
2. **System design**: Read Phase 3-4. Understanding the scheduler and GC is critical for designing high-throughput services.
3. **Production readiness**: Read Phase 5. The pitfalls doc alone will save you from the most common production bugs.
4. **Deep debugging**: Docs 08 and 09 give you the tools (`pprof`, `trace`, `GODEBUG`) to diagnose production issues.

### Companion Tools

```bash
go build -gcflags='-m'           # escape analysis — what escapes to heap
go build -gcflags='-m -m'        # verbose escape analysis with reasons
go build -gcflags='-S'           # assembly output — see runtime calls
go test -race ./...              # data race detection — NON-NEGOTIABLE
go test -bench=. -memprofile=mem.out  # memory profiling
go tool pprof mem.out            # analyze profiles
go tool trace trace.out          # visual goroutine timeline
GODEBUG=gctrace=1 ./app          # GC cycle stats
GODEBUG=schedtrace=1000 ./app    # scheduler state every second
```

---

> Each document follows the same structure: runtime struct layout → step-by-step
> trace → production implications → performance cost table → quick reference card.

