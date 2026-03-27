# 📚 Go Under the Hood — Deep Dive Series

Runtime-level explanations of Go's core systems — what the compiler and runtime
actually do under the hood. Every document includes ASCII diagrams, `runtime/` source
references, performance cost tables, and production implications.

---

## Reading Order

The documents are organized in a deliberate learning progression — from how Go
stores data in memory, through its type system and concurrency model, to production
debugging and enterprise patterns.

### Part I — Data Structures in Depth

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 01 | [Slices: The Three-Word Header](./01_slices_three_word_header.md) | Slice header `{ptr, len, cap}`, `runtime.growslice`, append algorithm, backing array sharing, memory leaks |
| 02 | [Maps: Buckets, Growth & the Never-Shrink Truth](./02_maps_buckets_and_growth.md) | `runtime.hmap`, buckets, tophash, load factor 6.5, incremental evacuation, concurrent access fatals |
| 03 | [Strings: Immutability, UTF-8 & the Substring Trap](./03_strings_immutability_and_utf8.md) | 2-word header, runes vs bytes, substring memory leaks, concatenation costs, `strings.Builder` |

### Part II — Functions, Closures & Variables

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 04 | [Variables: var vs := vs Explicit Type](./04_variables_var_vs_short_decl.md) | Three declaration forms, type inference rules, `:=` redeclaration, shadowing, idiomatic conventions |
| 05 | [Pointers: Auto-Deref & Auto-Address](./05_pointers_auto_deref.md) | Auto-dereference, auto-address, addressability rules, when Go silently inserts `*` and `&` |
| 06 | [Closures: funcval, Capture by Reference](./06_closures_funcval_and_capture.md) | Block scoping, `funcval` struct, capture-by-reference, escape analysis, defer internals, Go 1.22 changes |
| 07 | [Loop Variable Capture & the Fan-In Pattern](./07_loop_capture_and_fanin.md) | The classic loop bug, memory layout per iteration, Go 1.22 range fix, fan-in pattern line-by-line |

### Part III — Type System & Interfaces

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 08 | [Interfaces: iface, eface & the Nil Trap](./08_interfaces_iface_eface_nil_trap.md) | `iface`/`eface`, `itab` dispatch, nil trap, three guards, struct embedding & method promotion |
| 09 | [The any Type: Boxing & Performance Cost](./09_any_type_boxing_and_cost.md) | Boxing/unboxing, `convT` family, `staticuint64s`, generics vs `any`, GC shape stenciling |
| 10 | [Context: The Interface Design Masterclass](./10_context_interface_masterclass.md) | 4 methods, cancelCtx/timerCtx/valueCtx chain, decorator pattern, production gotchas |

### Part IV — Error Handling

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 11 | [Error Chains: Wrapping, Is, As & Strategy](./11_error_chains_wrapping_strategy.md) | Sentinel errors, custom types, `%w` wrapping, `errors.Is/As` chain walk, enterprise error strategy |
| 12 | [Concurrent Errors: errgroup & Recovery](./12_concurrent_errors_errgroup.md) | WaitGroup+Mutex, channel-based collection, `errgroup.Group`, panic recovery in goroutines |

### Part V — Concurrency & Channels

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 13 | [Goroutines & the GMP Scheduler](./13_goroutines_gmp_scheduler.md) | GMP model, work stealing, async preemption, syscall hand-off, network poller |
| 14 | [Goroutine Stacks: Growth, Shrinking & Pointers](./14_goroutine_stacks_growth.md) | Contiguous stacks, 2KB→growable, pointer adjustment, stack maps, scope vs frame vs stack |
| 15 | [Channels: hchan, Select & Production Patterns](./15_channels_hchan_select.md) | `runtime.hchan`, ring buffer, `sudog` parking, select algorithm, nil channels, fan-in/out, pipelines |

### Part VI — Runtime & Performance

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 16 | [Memory, GC & Escape Analysis](./16_memory_gc_escape_analysis.md) | Stack vs heap, escape analysis rules, tri-color GC, write barrier, GOGC/GOMEMLIMIT, sync.Pool |
| 17 | [Sorting: pdqsort](./17_sorting_pdqsort.md) | Pattern-defeating quicksort: insertion sort + quicksort + heapsort adaptive algorithm |

### Part VII — Testing & Debugging

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 18 | [Testing Internals: go test Under the Hood](./18_testing_internals.md) | Code generation, test discovery, `_testmain.go`, white-box vs black-box, benchmarks, fuzz |
| 19 | [Debugging & Profiling: pprof, trace & dlv](./19_debugging_profiling.md) | pprof, go tool trace, dlv, GODEBUG, deadlock detection, goroutine leak monitoring |

### Part VIII — Production Engineering

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 20 | [Production Pitfalls & Enterprise Patterns](./20_production_pitfalls_enterprise.md) | Top 15 production bugs, library comparison, project structure, graceful shutdown, Docker builds |

### Part IX — Design Philosophy (Living Document)

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 21 | [Go Design Philosophy: The Connected Architecture](./21_go_design_philosophy.md) | How immutability, interfaces, value semantics, and CSP concurrency form one coherent system |

---

## How to Use These Documents

1. **Interview prep**: Read Parts I–III. These cover the most common Go interview questions at senior level.
2. **System design**: Read Parts V–VI. Understanding the scheduler and GC is critical for high-throughput services.
3. **Production readiness**: Read Parts VII–VIII. Debugging tools + pitfalls will save you in production.
4. **Quick reference**: Every doc ends with a Quick Reference Card — use them as cheat sheets.

### Companion Tools

```bash
go build -gcflags='-m'           # escape analysis — what escapes to heap
go build -gcflags='-m -m'        # verbose escape analysis with reasons
go build -gcflags='-S'           # assembly output — see runtime calls
go test -race ./...              # data race detection — NON-NEGOTIABLE
go test -bench=. -benchmem       # benchmark with allocation counting
go tool pprof cpu.out            # analyze CPU profiles
go tool pprof -http=:8080 mem.out # web UI for memory profiles
go tool trace trace.out          # visual goroutine timeline
GODEBUG=gctrace=1 ./app          # GC cycle stats
GODEBUG=schedtrace=1000 ./app    # scheduler state every second
dlv debug ./cmd/server           # interactive debugger
```

---

> Each document follows the same structure: runtime struct layout → step-by-step
> trace → production implications → performance cost table → quick reference card.

