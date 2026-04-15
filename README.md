# Go Under the Hood

A structured, hands-on Go curriculum for **senior engineers**. Runtime-level deep
dives, 100+ exercises with TDD tests, 200+ coding problems, and 67 runnable demos.

No boilerplate explanations. No trivial information. Every topic goes under the hood:
what the compiler and runtime actually do, why the Go team designed it that way,
and what breaks at scale.

> **Audience:** Engineers with 5+ years of experience transitioning to or deepening
> mastery in Go. You know what a pointer is. You want to know how Go's escape
> analysis decides whether it lives on the stack or the heap.

---

## 🗺️ Learning Path

This repo is designed to be followed in order. Each phase builds on the previous
one. Within each phase, **read the deep dive first**, then **run the demo**, then
**implement the exercises** (make the failing tests pass).

```
 Phase 1    Language Core         Read → Watch → Implement → Test
 Phase 2    Concurrency           The #1 mental model shift from other languages
 Phase 3    Standard Library      io, encoding, testing — the Go way
 Phase 4    Problem Solving       200+ problems, solved idiomatically in Go
```

### Phase 1 — Language Core

Learn Go's data structures, type system, and value semantics. Focus on how they
differ from Java/C#/Python under the hood.

| Step | Read (learnings/) | Run (cmd/concepts/) | Implement (exercises/) |
|------|-------------------|---------------------|----------------------|
| 1.1 | [04 — Variables & Pointers](learnings/04_variables_pointers_declarations.md) | `basics/01-variables` | `fundamentals/01_basics` |
| 1.2 | [22 — Control Flow Under the Hood](learnings/22_control_flow_under_the_hood.md) | `control-flow/*` | `fundamentals/02_control_flow` |
| 1.3 | [05 — Closures & funcval](learnings/05_closures_funcval_and_capture.md) | `functions/*` | `fundamentals/03_functions` |
| 1.4 | [04 — Pointers deep dive](learnings/04_variables_pointers_declarations.md) | `pointers/*` | `fundamentals/04_pointers` |
| 1.5 | [06 — Interfaces & iface/eface](learnings/06_interfaces_iface_eface_nil_trap.md) | `structs/*`, `interfaces/*` | `fundamentals/05_structs` then `06_interfaces` |
| 1.6 | [08 — Error Chains](learnings/08_error_chains_wrapping_strategy.md), [09 — errgroup](learnings/09_concurrent_errors_errgroup.md) | `error-handling/*` | `fundamentals/07_error_handling` |
| 1.7 | [01 — Slices: Three-Word Header](learnings/01_slices_three_word_header.md) | `arrays-slices/*` | `fundamentals/08_arrays_slices` |
| 1.8 | [02 — Maps: Buckets & Growth](learnings/02_maps_buckets_and_growth.md) | `maps/*` | `fundamentals/09_maps` |
| 1.9 | [03 — Strings: Immutability](learnings/03_strings_immutability_and_boxing.md) | — | — |

### Phase 2 — Concurrency

This is where Go diverges most from your background. Goroutines are not threads.
Channels are not queues. The GMP scheduler is the key to understanding everything.

| Step | Read (learnings/) | Run (cmd/concepts/) | Implement (exercises/) |
|------|-------------------|---------------------|----------------------|
| 2.1 | [10 — GMP Scheduler](learnings/10_goroutines_gmp_scheduler.md), [11 — Goroutine Stacks](learnings/11_goroutine_stacks_growth.md) | `goroutines/*` | `fundamentals/10_goroutines` |
| 2.2 | [12 — Channels & hchan](learnings/12_channels_hchan_select.md) | `channels/*` (16 demos!) | `fundamentals/11_channels` |
| 2.3 | [19 — Context Masterclass](learnings/19_context_interface_masterclass.md) | — | — |

### Phase 3 — Standard Library & Tooling

Master the patterns that make Go code idiomatic: `io.Reader`/`io.Writer`, struct
tags, table-driven tests, and the build toolchain.

| Step | Read (learnings/) | Run (cmd/concepts/) | Implement (exercises/) |
|------|-------------------|---------------------|----------------------|
| 3.1 | [03 — Strings](learnings/03_strings_immutability_and_boxing.md) | `stdlib/01-strings-strconv` | `stdlib/01_strings_strconv` |
| 3.2 | [13 — Memory & Sorting](learnings/13_memory_gc_escape_sorting.md) | — | `stdlib/02_sort` |
| 3.3 | — | — | `stdlib/03_builtins` |
| 3.4 | — | — | `stdlib/04_io_files` |
| 3.5 | — | — | `stdlib/05_encoding_json` |
| 3.6 | — | — | `stdlib/06_math` |
| 3.7 | [14 — Testing Internals](learnings/14_testing_internals.md) | — | `stdlib/07_testing` |
| 3.8 | [12 — Packages & Modules](learnings/20_practical_go_toolchain.md) | `packages-modules/*` | `fundamentals/12_packages_modules` |

### Phase 4 — Problem Solving in Go

200+ problems across 14 categories. Each problem is a stub with hints — implement
the function, run the test, make it pass. Focus on **idiomatic Go**, not just
correctness.

| # | Category | Count | Key Problems |
|---|----------|:-----:|-------------|
| 01 | Arrays | 25 | Two Sum, Max Subarray, Array Manipulation |
| 02 | Strings | 18 | Valid Anagram, Group Anagrams, Sherlock Valid String |
| 03 | Linked List | 14 | Reverse, Merge Two, Detect Cycle |
| 04 | Stacks & Queues | 10 | Valid Parens, Min Stack, Largest Rectangle |
| 05 | Binary Search | 11 | Search Rotated, Climbing Leaderboard |
| 06 | Sliding Window | 7 | Min Window Substring, Fruit Baskets |
| 07 | Trees | 4 | Level Order, Validate BST, LCA |
| 08 | Graphs | 6 | Number of Islands, Course Schedule |
| 09 | Dynamic Prog | 8 | Coin Change, LIS, Jumping on Clouds |
| 10 | Two Pointers | 7 | 3Sum, Container With Most Water |
| 11 | Hard | 13 | N-Queens, Edit Distance, Alien Dictionary |
| 12 | Backtracking | 9 | Permutations, Subsets |
| 13 | Bit Manipulation | 9 | Single Number, Counting Bits |
| 14 | Heap | 7 | Top K Frequent, Merge K Lists |

### Bonus — Deep Dives for Production Mastery

Read these after completing Phases 1-3. They cover runtime internals, debugging,
and enterprise architecture — the knowledge that separates "writes Go" from
"engineers production Go systems."

| Chapter | Topic | When to Read |
|---------|-------|-------------|
| [07 — any Type & Boxing](learnings/07_any_type_boxing_and_cost.md) | `convT` family, `staticuint64s`, generics vs `any` | After interfaces (1.5) |
| [13 — Memory, GC & Escape Analysis](learnings/13_memory_gc_escape_sorting.md) | Stack vs heap, tri-color GC, `GOGC`/`GOMEMLIMIT` | After concurrency (Phase 2) |
| [15 — Debugging & Profiling](learnings/15_debugging_profiling.md) | pprof, `go tool trace`, Delve, `GODEBUG` | When debugging real code |
| [16 — Go Design Philosophy](learnings/16_go_design_philosophy.md) | How immutability, interfaces, CSP form one system | When you want the "why" |
| [17 — Middleware Pattern](learnings/17_middleware_pattern.md) | Function types, `HandlerFunc`, closure capture | Before building HTTP services |
| [18 — Production Patterns](learnings/18_production_patterns_enterprise.md) | Top 15 pitfalls, DI, graceful shutdown, Docker | Before deploying to production |
| [21 — Zero Values & mallocgc](learnings/21_zero_values_mallocgc_syncpool_duffzero.md) | Zeroing pipeline, `sync.Pool`, `duffzero` assembly | For runtime internals mastery |

---

## 📁 Repository Structure

```
go-learning-guide/
├── learnings/          ← 📖 22 deep-dive chapters (Go Under the Hood series)
├── cmd/concepts/       ← 🎯 67 runnable demos (go run each one)
├── exercises/
│   ├── fundamentals/   ← ✏️ 12 packages — language core exercises
│   └── stdlib/         ← ✏️ 7 packages — standard library exercises
├── problems/           ← 🧩 200+ coding problems across 14 categories
├── tools/              ← 🔧 md2pdf — generates the companion PDF book
└── utils/              ← 📦 Shared helpers (ListNode, TreeNode, etc.)
```

## 📐 How Exercises Work

Every exercise follows the TDD pattern: **the tests are the spec**.

1. Open `exercises.go` — you'll see stub functions that return zero values
2. Read the function signature, comments, and the corresponding test file
3. Implement the function to make the tests pass
4. Run with the race detector — always

```bash
# Run one package
go test -race -v ./exercises/fundamentals/06_interfaces/

# Run all exercises
go test -race ./exercises/...

# Run all problems
go test -race ./problems/...
```

All tests will **FAIL** on a fresh clone. That's by design. Your job is to make
them pass.

---

## 📖 Go Under the Hood — Full Chapter List

The [`learnings/`](learnings/) directory contains **22 chapters** organized into
9 parts. Each chapter includes runtime source references, ASCII memory diagrams,
performance cost tables, and a quick reference card.

**[→ Full Table of Contents](learnings/README.md)**

| Part | Chapters | Focus |
|------|----------|-------|
| I — Data Structures | 01-03 | Slices, Maps, Strings — memory layout & runtime internals |
| II — Language Mechanics | 04-05 | Variables, pointers, closures, `funcval`, capture semantics |
| III — Type System | 06-07 | `iface`/`eface`, `itab` dispatch, `any` boxing cost |
| IV — Error Handling | 08-09 | Error chains, `%w` wrapping, `errgroup`, panic recovery |
| V — Concurrency | 10-12 | GMP scheduler, goroutine stacks, `hchan`, select algorithm |
| VI — Runtime & Performance | 13 | GC, escape analysis, `GOGC`/`GOMEMLIMIT`, pdqsort |
| VII — Testing & Debugging | 14-15 | Testing internals, pprof, `go tool trace`, Delve |
| VIII — Design & Architecture | 16-18 | Design philosophy, middleware pattern, enterprise patterns |
| IX — Cross-Cutting | 19-22 | Context, Go toolchain, zero values & `mallocgc`, control flow internals |

A **PDF version** can be generated with:
```bash
pip install markdown pymupdf
python tools/md2pdf.py --config tools/book_config.json
```

---

## 🧪 Running Tests

```bash
go test ./...                              # run everything
go test -race ./...                        # with race detector (MANDATORY)
go test -cover ./...                       # with coverage
go test -v -run TestName ./path/           # single test, verbose
go test -bench=. -benchmem ./path/         # benchmarks + allocations
```

## 🔬 Performance Analysis

```bash
go build -gcflags='-m' ./path/             # escape analysis
go build -gcflags='-m -m' ./path/          # verbose escape analysis with reasons
go build -gcflags='-S' ./path/             # assembly output
go test -cpuprofile=cpu.out -bench=.       # CPU profile
go tool pprof cpu.out                      # analyze: top, list, web
GODEBUG=gctrace=1 ./app                    # GC trace
GODEBUG=schedtrace=1000 ./app              # scheduler state every second
```

---

## 🔧 Setup

```bash
git clone https://github.com/mert-unsal/go-learning-guide.git
cd go-learning-guide
go build ./...        # verify everything compiles
go test ./...         # all tests FAIL — that's correct, you implement them
```

> **Go version:** 1.24+ | **Module:** `go-learning-guide` | **Dependencies:** stdlib only
