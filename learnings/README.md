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
| 03 | [Strings: Immutability & Boxing](./03_strings_immutability_and_boxing.md) | 2-word header, runes vs bytes, substring memory leaks, concatenation costs, `strings.Builder`, interface boxing |

### Part II — Language Mechanics

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 04 | [Variables, Pointers & Declarations](./04_variables_pointers_declarations.md) | Three declaration forms, type inference, `:=` shadowing, auto-deref, auto-address, addressability rules |
| 05 | [Closures: funcval, Capture & Production Patterns](./05_closures_funcval_and_capture.md) | `funcval` struct, capture-by-reference, escape analysis, loop capture gotcha, defer internals, Go 1.22 changes, cross-language comparison |

### Part III — Type System & Interfaces

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 06 | [Interfaces: iface, eface & the Nil Trap](./06_interfaces_iface_eface_nil_trap.md) | `iface`/`eface`, `itab` dispatch, nil trap, three guards, struct embedding & method promotion |
| 07 | [The any Type: Boxing & Performance Cost](./07_any_type_boxing_and_cost.md) | Boxing/unboxing, `convT` family, `staticuint64s`, generics vs `any`, GC shape stenciling |

### Part IV — Error Handling

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 08 | [Error Chains: Wrapping, Is, As & Strategy](./08_error_chains_wrapping_strategy.md) | Sentinel errors, custom types, `%w` wrapping, `errors.Is/As` chain walk, enterprise error strategy |
| 09 | [Concurrent Errors: errgroup & Recovery](./09_concurrent_errors_errgroup.md) | WaitGroup+Mutex, channel-based collection, `errgroup.Group`, panic recovery in goroutines |

### Part V — Concurrency & Channels

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 10 | [Goroutines & the GMP Scheduler](./10_goroutines_gmp_scheduler.md) | GMP model, work stealing, async preemption, syscall hand-off, network poller, interview self-test |
| 11 | [Goroutine Stacks: Growth, Shrinking & Pointers](./11_goroutine_stacks_growth.md) | Contiguous stacks, 2KB→growable, pointer adjustment, stack maps, scope vs frame vs stack |
| 12 | [Channels: hchan, Select & Production Patterns](./12_channels_hchan_select.md) | `runtime.hchan`, ring buffer, `sudog` parking, select algorithm, nil channels, fan-in/out, pipelines |

### Part VI — Runtime & Performance

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 13 | [Memory, GC, Escape Analysis & Sorting](./13_memory_gc_escape_sorting.md) | Stack vs heap, escape analysis rules, tri-color GC, write barrier, GOGC/GOMEMLIMIT, sync.Pool, pdqsort |

### Part VII — Testing & Debugging

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 14 | [Testing Internals: go test Under the Hood](./14_testing_internals.md) | Code generation, test discovery, `_testmain.go`, white-box vs black-box, benchmarks, fuzz |
| 15 | [Debugging & Profiling: pprof, trace & dlv](./15_debugging_profiling.md) | pprof, go tool trace, dlv, GODEBUG, deadlock detection, goroutine leak monitoring |

### Part VIII — Design & Architecture

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 16 | [Go Design Philosophy: The Connected Architecture](./16_go_design_philosophy.md) | How immutability, interfaces, value semantics, and CSP concurrency form one coherent system |
| 17 | [The Middleware Pattern: Function Types & Closures](./17_middleware_pattern.md) | Methods on function types, HandlerFunc adapter, closure capture in middleware, Java comparison |
| 18 | [Production Patterns & Enterprise Go](./18_production_patterns_enterprise.md) | Top 15 pitfalls, library comparison, project structure, functional options, graceful shutdown, DI, Docker |

### Part IX — Cross-Cutting Concepts

| # | Document | What You'll Learn |
|---|----------|-------------------|
| 19 | [Context: The Interface Design Masterclass](./19_context_interface_masterclass.md) | 4 methods, cancelCtx/timerCtx/valueCtx chain, decorator pattern, production gotchas |
| 20 | [Practical Go Toolchain](./20_practical_go_toolchain.md) | Go modules, go build, cross-compilation, Docker multi-stage, config/env patterns, cheatsheet |
| 21 | [Zero Values, mallocgc, sync.Pool & duffzero](./21_zero_values_mallocgc_syncpool_duffzero.md) | Zeroing pipeline, `mallocgc` allocation paths, `mspan.needzero`, `duffzero` assembly, `sync.Pool` violation |
| 22 | [Control Flow Under the Hood](./22_control_flow_under_the_hood.md) | Defer internals (3 implementations), range compiler rewrites (6 types), switch dispatch, Go 1.22 loop variable change |
| 23 | [io.Reader/Writer Deep Dive](./23_io_reader_writer_deep_dive.md) | One-method contract, decorator pattern, LimitReader/MultiReader/TeeReader/Pipe, io.Copy 3 paths, sendfile(2), bufio internals, zero-copy |
| 24 | [encoding/json Under the Hood](./24_encoding_json_under_the_hood.md) | Reflect-driven encoding, encoder cache, struct tags, Decoder vs Unmarshal, custom marshalers, RawMessage, json.Number precision |
| 25 | [sync Primitives Under the Hood](./25_sync_primitives_under_the_hood.md) | Mutex spinning/starvation, RWMutex, WaitGroup, Once fast/slow path, Pool per-P cache, sync.Map dual-store, Cond, atomics |

---

## How to Use These Documents

1. **Deep dive**: Read Parts I–III. These cover the most common Go topics at senior level.
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
