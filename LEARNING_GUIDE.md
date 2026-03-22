# 🗺️ Go Mastery Roadmap — Interview-Ready in 7 Milestones

> **Senior engineer track.** This is not a beginner guide. You already know
> programming — this roadmap takes you from "knows Go syntax" to
> "can explain the runtime, ace coding interviews, and ship production Go."

---

## ⚡ Where You Are Now

| Area | Status | Detail |
|------|--------|--------|
| Fundamentals (12 packages) | ✅ All exercises implemented, all tests pass | Core Go syntax mastered |
| Interfaces | ✅ Deep mastery | `iface`/`eface`, nil trap, method sets, three guards |
| HackerRank (9 problems) | ✅ Complete | All implemented and passing |
| Deep dive docs (11 documents) | ✅ Written | 7,400+ lines of runtime internals |
| Stdlib (7 packages) | 🔲 37 exercises TODO | Need implementation |
| LeetCode (14 categories) | 🔲 207 problems TODO | All stubs, 0% implemented |
| Performance tuning | 🔲 8 exercises TODO | Module created, not attempted |

---

## 🗺️ The Roadmap

```
 M1: Data Structures    ──►  M2: Language Core    ──►  M3: Concurrency
 (slices, maps)              (closures, ptrs, err)     (channels, goroutines)
       │                           │                          │
       ▼                           ▼                          ▼
 M4: Runtime & Perf      ◄── All theory done ──►  M5: Stdlib Mastery
 (GC, escape, profiling)                           (strings, io, json, sort)
       │                                                  │
       ▼                                                  ▼
 M6: Algorithm Patterns + LeetCode  ◄────────────────────┘
 (two-pointers, sliding-window, binary-search, DP, trees, graphs)
       │
       ▼
 M7: Production Readiness & Mock Interviews
 (pitfalls, system design, concurrency patterns)
```

**M1-M4** = deep reading (your 11 docs are the curriculum)
**M5** = stdlib fluency (37 coding exercises)
**M6** = the grind (207 LeetCode problems, pattern-by-pattern)
**M7** = synthesis (production patterns + mock interviews)

> M1-M4 can overlap with M5 — read a doc, then do related stdlib exercises.

---

## M1 — Data Structures Internals

**Goal:** Understand how Go stores data at the runtime level.

### 📖 Read

| Document | Key Topics |
|----------|-----------|
| [`learnings/01_slices_internals.md`](learnings/01_slices_internals.md) | Slice header (3-word struct), `runtime.growslice`, backing array sharing, memory leaks |
| [`learnings/02_maps_internals.md`](learnings/02_maps_internals.md) | `runtime.hmap`, buckets, tophash, load factor 6.5, incremental evacuation |
| [`learnings/10_sorting_pdqsort.md`](learnings/10_sorting_pdqsort.md) | How `slices.Sort` adaptively picks insertion sort, quicksort, or heapsort |

### 🔬 Practice

- Re-read `fundamentals/08_arrays_slices/concepts.go` with the internals lens
- Re-read `fundamentals/09_maps/concepts.go` — now you know what `hmap` does under the hood
- Run escape analysis on your code: `go build -gcflags='-m' ./fundamentals/08_arrays_slices/`

### ✅ Mastery Check — Can You Answer These?

1. What's the growth factor of `append()` and when does it change?
2. Why does `s2 := s[:3]` create a memory leak risk?
3. Why can't you write to a nil map but CAN append to a nil slice?
4. What triggers map evacuation and how does it work incrementally?
5. How does `tophash` optimize bucket lookup? What values are reserved?
6. What's the difference between `s[:3]` and `s[:3:3]`? (full slice expression)

---

## M2 — Language Core Internals

**Goal:** How Go's type system, scoping, and error model work under the hood.

### 📖 Read

| Document | Key Topics |
|----------|-----------|
| [`learnings/03_closures_and_scopes.md`](learnings/03_closures_and_scopes.md) | `funcval` struct, capture-by-reference, escape analysis, Go 1.22 loop fix |
| [`learnings/04_pointers_and_auto_deref.md`](learnings/04_pointers_and_auto_deref.md) | Auto-dereference, auto-address, addressability rules |
| [`learnings/06_error_handling_patterns.md`](learnings/06_error_handling_patterns.md) | Sentinel errors, `%w` wrapping, `errors.Is`/`As` chain walk, enterprise strategy |
| [`learnings/05_interfaces_internals.md`](learnings/05_interfaces_internals.md) | *(Quick refresher — already mastered)* |

### 🔬 Practice

- Re-read `fundamentals/03_functions/` — focus on closure exercises
- Re-read `fundamentals/04_pointers/` — pointer manipulation with addressability understanding
- Re-read `fundamentals/07_error_handling/` — error wrapping exercises

### ✅ Mastery Check

1. How does a closure capture variables? Stack or heap? What does escape analysis say?
2. What changed in Go 1.22 with loop variable scoping and why?
3. When does Go auto-insert `*` or `&`? What is "addressability"?
4. Walk through `errors.Is()` chain walk — what happens with wrapped errors?
5. What's the difference between `errors.Is()` and `errors.As()`?

---

## M3 — Concurrency Internals

**Goal:** How Go runs things in parallel — the GMP model and channel machinery.

### 📖 Read

| Document | Key Topics |
|----------|-----------|
| [`learnings/07_channels_internals.md`](learnings/07_channels_internals.md) | `runtime.hchan`, ring buffer, `sudog` parking, select algorithm, nil channels |
| [`learnings/08_goroutines_and_scheduler.md`](learnings/08_goroutines_and_scheduler.md) | GMP model, work stealing, async preemption, syscall hand-off, network poller |

### 🔬 Practice

- Re-read `fundamentals/10_goroutines/` — goroutine exercises
- Re-read `fundamentals/11_channels/` — channel exercises
- Study `practical/06_concurrency_patterns/` — worker pool implementation
- Run ALL concurrency code with `-race`: `go test -race ./fundamentals/10_goroutines/ ./fundamentals/11_channels/`

### ✅ Mastery Check

1. Draw the GMP model. What happens when a goroutine enters a syscall?
2. What's inside an `hchan`? What are the 3 send paths?
3. When would you use a mutex over a channel? Give 3 concrete examples.
4. Design a graceful shutdown sequence for a service with 3 worker pools.
5. What does `GODEBUG=schedtrace=1000` show? What about `go tool trace`?
6. What happens when you send on a nil channel? Receive? Close?

---

## M4 — Runtime, GC & Performance Engineering

**Goal:** Memory management, GC tuning, escape analysis, profiling workflow.

### 📖 Read

| Document | Key Topics |
|----------|-----------|
| [`learnings/09_memory_gc_and_escape_analysis.md`](learnings/09_memory_gc_and_escape_analysis.md) | Stack vs heap, escape analysis rules, tri-color GC, write barrier, `GOGC`/`GOMEMLIMIT`, `sync.Pool` |

### 🔬 Practice — Complete All 8 Performance Exercises

All in `practical/07_performance_tuning/exercises.go`:

| # | Exercise | What You Learn |
|---|----------|---------------|
| 1 | Escape analysis | `fmt.Sprintf` boxing cost, pointer escape, `strconv` alternative |
| 2 | String building | `+=` concatenation O(n²) disaster → `strings.Builder` → pre-alloc |
| 3 | Slice pre-allocation | `make([]T, 0, n)` to avoid repeated `growslice` calls |
| 4 | `sync.Pool` | Buffer reuse pattern, 2-generation victim cache |
| 5 | Interface boxing | Concrete type vs interface dispatch in hot paths |
| 6 | Map pre-allocation | `make(map[K]V, n)` to avoid rehashing |
| 7 | GC pressure | Pointer-heavy structs vs value structs (89k vs 39k allocs) |
| 8 | Struct padding | Field ordering to minimize alignment padding |

### 🛠️ Tools to Master

```bash
go build -gcflags='-m'                    # escape analysis — what goes to heap
go build -gcflags='-m -m'                 # verbose escape reasons
go test -bench=. -benchmem                # benchmark + allocation counting
go test -cpuprofile=cpu.out -bench=.      # CPU profiling
go tool pprof cpu.out                     # analyze: top, list, web
go test -memprofile=mem.out -bench=.      # memory profiling
GODEBUG=gctrace=1 ./app                   # GC cycle stats in real time
go tool trace trace.out                   # visual goroutine timeline
```

### ✅ Mastery Check

1. What are the 5 common triggers for heap escape?
2. How does `GOGC=100` differ from `GOMEMLIMIT`? When use which?
3. Your service has 200ms p99 spikes every 2 min — diagnose with GC trace.
4. Explain `sync.Pool`'s victim cache. When do objects get collected?
5. Why is `fmt.Sprintf("%d", n)` slower than `strconv.Itoa(n)` in hot paths?

---

## M5 — Standard Library Mastery

**Goal:** Become fluent with Go's stdlib — the tools you'll use daily.

### 🔬 Practice — 37 Exercises Across 7 Packages

| Package | Exercises | Key APIs |
|---------|:---------:|----------|
| `stdlib/01_strings_strconv/` | 7 | `strings.Builder`, `Split/Join`, `strconv.Atoi/Itoa` |
| `stdlib/02_sort/` | 5 | `slices.Sort`, `sort.Slice`, custom comparators |
| `stdlib/03_builtins/` | 6 | `make`, `len`, `cap`, `copy`, `delete`, `clear` |
| `stdlib/04_io_files/` | 4 | `io.Reader/Writer`, `bufio.Scanner`, `os.ReadFile` |
| `stdlib/05_encoding_json/` | 4 | `json.Marshal/Unmarshal`, struct tags, streaming |
| `stdlib/06_math/` | 6 | `math.Max/Min/Abs`, `math/rand` |
| `stdlib/07_testing/` | 5 | Table-driven, `t.Run`, benchmarks, fuzzing |

### 📚 Key Stdlib Areas for Interviews (beyond exercises)

- **`context`** — `WithCancel`, `WithTimeout`, `WithValue`, propagation rules
- **`sync`** — `Mutex`, `RWMutex`, `WaitGroup`, `Once`, `Pool`, `Map`
- **`net/http`** — handlers, middleware pattern, `Server.Shutdown`
- **`encoding/json`** — struct tags, custom `MarshalJSON`, streaming decoder

```bash
# Run all stdlib tests
go test -race ./stdlib/...
```

---

## M6 — Algorithm Patterns & LeetCode

**Goal:** Build pattern recognition for coding interviews. 207 problems total.

### Phase A — Core Patterns (solve 3-5 per pattern first)

| Pattern | Folder | Problems | Key Problems |
|---------|--------|:--------:|-------------|
| Two Pointers | `leetcode/10_two_pointers/` | 7 | 3Sum, Container With Most Water |
| Sliding Window | `leetcode/06_sliding_window/` | 12 | Longest Substring Without Repeat |
| Binary Search | `leetcode/05_binary_search/` | 15 | Search in Rotated Sorted Array |
| Stacks & Queues | `leetcode/04_stacks_queues/` | 15 | Valid Parentheses, Min Stack |
| Arrays | `leetcode/01_arrays/` | 34 | Two Sum, Max Subarray, Product Except Self |
| Strings | `leetcode/02_strings/` | 22 | Valid Anagram, Group Anagrams |

### Phase B — Data Structure Patterns

| Pattern | Folder | Problems | Key Problems |
|---------|--------|:--------:|-------------|
| Linked Lists | `leetcode/03_linked_list/` | 25 | Reverse, Merge Two, Detect Cycle |
| Trees | `leetcode/07_trees/` | 10 | Max Depth, Level Order, Validate BST |
| Graphs | `leetcode/08_graphs/` | 6 | Number of Islands, Clone Graph |
| Heaps | `leetcode/14_heap_priority_queue/` | 7 | Top K Frequent, Merge K Lists |

### Phase C — Advanced Patterns

| Pattern | Folder | Problems | Key Problems |
|---------|--------|:--------:|-------------|
| Dynamic Programming | `leetcode/09_dynamic_prog/` | 7 | Climbing Stairs, Coin Change |
| Backtracking | `leetcode/12_backtracking/` | 9 | Permutations, N-Queens |
| Bit Manipulation | `leetcode/13_bit_manipulation/` | 9 | Single Number, Counting Bits |
| Hard | `leetcode/11_hard/` | 13 | Trapping Rain Water, Median of Two |

### Strategy

- Solve in **idiomatic Go** — not a translation from Python/Java
- Always run with `-race` on concurrent solutions
- After solving, compare with `solutions.go` — discuss tradeoffs with GoSensei
- Target: 3-5 problems/day

```bash
# Run tests for a specific category
go test -v ./leetcode/01_arrays/...

# Run a single problem
go test -run TestTwoSum ./leetcode/01_arrays/
```

---

## M7 — Production Readiness & Mock Interviews

**Goal:** Tie everything together — system design, pitfalls, real-world Go.

### 📖 Read

| Document | Key Topics |
|----------|-----------|
| [`learnings/11_production_go_pitfalls.md`](learnings/11_production_go_pitfalls.md) | Top 15 production bugs, library comparison tables (8 categories), project structure, graceful shutdown, Docker multi-stage builds |

### 🔬 Practice

- Review `practical/05_config_env_json_yaml/` — config management patterns
- Review `practical/01-04/` concepts — dependency mgmt, build/deploy, Docker, debugging

### 🎤 Mock Interview Topics

| Topic | What to Prepare |
|-------|----------------|
| System Design | Design a rate-limited HTTP API with middleware, context, graceful shutdown |
| Concurrency | Design a concurrent file processor with backpressure and error propagation |
| Runtime | Explain Go's GC to an interviewer (tri-color, write barrier, tuning knobs) |
| Debugging | Walk through diagnosing a memory leak using `pprof` |
| Code Review | Spot concurrency bugs, escape issues, error handling gaps in sample code |
| Architecture | Explain `internal/` boundary, clean architecture in Go, DI without frameworks |

---

## 📊 Progress Tracker

```
M1 — Data Structures Internals
[ ] Read slices internals doc
[ ] Read maps internals doc
[ ] Read sorting/pdqsort doc
[ ] Review slice/map exercises with internals lens

M2 — Language Core Internals
[x] Interfaces deep dive (COMPLETE)
[ ] Read closures & scopes doc
[ ] Read pointers & auto-deref doc
[ ] Read error handling patterns doc
[ ] Review functions/pointers/errors exercises

M3 — Concurrency Internals
[ ] Read channels internals doc
[ ] Read goroutines & scheduler doc
[ ] Review concurrency exercises + worker pool

M4 — Runtime & Performance
[ ] Read memory/GC/escape analysis doc
[ ] Complete 8 performance tuning exercises
[ ] Master profiling tools (pprof, trace, gcflags)

M5 — Stdlib Mastery
[ ] Complete 37 stdlib exercises (7 packages)

M6 — Algorithm Patterns & LeetCode
[ ] Phase A: Core patterns (~100 problems)
[ ] Phase B: Data structure patterns (~48 problems)
[ ] Phase C: Advanced patterns (~38 problems)

M7 — Production & Mock Interviews
[ ] Read production pitfalls doc
[ ] Mock interview preparation
```

---

## 🛠️ Essential Commands

```bash
# Build & test
go build ./...                            # compile everything
go test ./...                             # run all tests
go test -race ./...                       # race detector (NON-NEGOTIABLE)
go test -v -run TestName ./path/          # single test, verbose
go test -cover ./...                      # coverage report

# Performance
go test -bench=. -benchmem ./path/        # benchmarks + allocations
go build -gcflags='-m' ./path/            # escape analysis
go build -gcflags='-S' ./path/            # assembly output
go test -cpuprofile=cpu.out -bench=.      # CPU profile
go tool pprof cpu.out                     # analyze profile

# Debugging
GODEBUG=gctrace=1 ./app                  # GC trace
GODEBUG=schedtrace=1000 ./app            # scheduler trace
go tool trace trace.out                   # visual timeline

# Documentation
go doc fmt.Sprintf                        # function docs
go doc -all strings.Builder               # full type docs
```

---

## 💡 Tips for This Stage

1. **Read the deep dive docs actively** — don't just skim. Draw the diagrams yourself.
2. **Run escape analysis** on code you read — predict what escapes before checking.
3. **Always `-race`** — make it muscle memory. No exceptions.
4. **Solve LeetCode in Go idiomatically** — use slices, maps, goroutines naturally.
5. **When stuck on a problem**, identify the pattern first (Two Pointers? Sliding Window? DP?).
6. **Read Go standard library source code** — it's the best Go code you'll ever read.

> *"Clear is better than clever."* — Go Proverbs
