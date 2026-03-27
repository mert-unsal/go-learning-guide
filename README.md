# Go Under the Hood

> **By Mert Unsal** — Senior Software Engineer

A structured, hands-on Go curriculum with **runtime-level deep dives**. No boilerplate explanations, no trivial information — only details worth mentioning.

This repository covers Go language fundamentals, runtime internals, algorithm patterns, 200+ LeetCode/HackerRank problems, and production engineering practices.

---

## 📁 Project Structure

```
go-interview-prep/
├── learnings/      ← 📖 Go Under the Hood — 20 deep-dive chapters (start here)
├── fundamentals/   ← Go language core — concepts, exercises, tests
├── stdlib/         ← Standard library deep-dives (7 packages, 37 exercises)
├── patterns/       ← Algorithm pattern templates (9 patterns)
├── leetcode/       ← 200+ problems across 14 categories
├── hackerrank/     ← 15 problems (Easy → Hard)
├── practical/      ← Real-world Go: Docker, debugging, config, performance
├── tools/          ← md2pdf — generates the companion PDF book
└── utils/          ← Shared helpers (ListNode, TreeNode, etc.)
```

---

## 📖 Go Under the Hood — Deep Dive Series

The [`learnings/`](learnings/) folder contains **20 chapters** organized into 8 parts — runtime source code references, ASCII memory diagrams, escape analysis traces, and production gotchas.

**[→ Full Table of Contents](learnings/README.md)**

| Part | Chapters | Focus |
|------|----------|-------|
| I — Data Structures | 01-03 | Slices, Maps, Strings — memory layout & runtime internals |
| II — Functions & Variables | 04-07 | Closures, pointers, `funcval`, capture-by-reference, Go 1.22 changes |
| III — Type System | 08-10 | `iface`/`eface`, `itab` dispatch, `any` boxing, context internals |
| IV — Error Handling | 11-12 | Error chains, `%w` wrapping, `errgroup`, panic recovery |
| V — Concurrency | 13-15 | GMP scheduler, contiguous stacks, `hchan`, select algorithm |
| VI — Runtime & Performance | 16-17 | GC, escape analysis, `GOGC`/`GOMEMLIMIT`, pdqsort |
| VII — Testing & Debugging | 18-19 | Testing internals, pprof, `go tool trace`, Delve |
| VIII — Production | 20 | Enterprise pitfalls, library comparisons, graceful shutdown |

A **PDF version** (Go_Under_the_Hood.pdf) can be generated with:
```bash
pip install markdown pymupdf
python tools/md2pdf.py --config tools/book_config.json
```

---

## 🏗️ Hands-On Tracks

### Fundamentals (`fundamentals/`)

Each package contains `concepts.go` (annotated examples), `exercises.go` (implement these), and tests.

| # | Package | Topics |
|---|---------|--------|
| 01 | `basics` | Variables, constants, types, zero values, fmt |
| 02 | `control_flow` | if/else, switch, for loops, defer |
| 03 | `functions` | Signatures, variadic, closures, multiple returns |
| 04 | `pointers` | Address-of, dereferencing, nil, pointer receivers |
| 05 | `structs` | Definition, embedding, tags, methods |
| 06 | `interfaces` | Implicit impl, empty interface, type assertion |
| 07 | `error_handling` | error type, custom errors, panic/recover |
| 08 | `arrays_slices` | Arrays vs slices, append, copy, 2D slices |
| 09 | `maps` | CRUD, existence check, maps of slices |
| 10 | `goroutines` | go keyword, WaitGroup, race conditions, Mutex |
| 11 | `channels` | Buffered/unbuffered, select, patterns |
| 12 | `packages_modules` | Visibility, init order, go.mod, build tags |

### Standard Library (`stdlib/`)

| # | Package | Topics |
|---|---------|--------|
| 01 | `strings_strconv` | strings.Builder, strconv, unicode |
| 02 | `sort` | sort.Slice, sort.Search, custom Less |
| 03 | `builtins` | make, new, append, copy, len/cap, delete, panic/recover |
| 04 | `io_files` | os.ReadFile/WriteFile, bufio.Scanner, io.Reader/Writer |
| 05 | `encoding_json` | json.Marshal/Unmarshal, struct tags, streaming |
| 06 | `math` | math, rand, big numbers |
| 07 | `testing` | Table-driven tests, subtests, benchmarks |

### Algorithm Patterns (`patterns/`)

| Pattern | Key Concept |
|---------|-------------|
| Binary Search | `left <= right`, left-bound variant |
| Sliding Window | expand right, shrink left |
| Two Pointers | converging from both ends |
| BFS / DFS | queue + visited / recursive closure + stack |
| Dynamic Programming | top-down memoization, bottom-up tabulation |
| Monotonic Stack | indices on stack, pop when invariant breaks |
| Union-Find | path compression + union by rank |
| Heap | `container/heap` interface |

### LeetCode (`leetcode/`) — 200+ Problems

| # | Category | Count | Key Problems |
|---|----------|:-----:|-------------|
| 01 | Arrays | 34 | Two Sum, Max Subarray, Product Except Self |
| 02 | Strings | 22 | Valid Anagram, Group Anagrams, Longest Substring |
| 03 | Linked List | 25 | Reverse, Merge Two, Detect Cycle |
| 04 | Stacks & Queues | 15 | Valid Parens, Min Stack, Largest Rectangle |
| 05 | Binary Search | 15 | Search Rotated, Median of Two Arrays |
| 06 | Sliding Window | 12 | Min Window Substring, Fruit Baskets |
| 07 | Trees | 10 | Level Order, Validate BST, LCA |
| 08 | Graphs | 6 | Number of Islands, Course Schedule |
| 09 | Dynamic Prog | 7 | Coin Change, LIS, Word Break |
| 10 | Two Pointers | 7 | 3Sum, Container With Most Water |
| 11 | Hard | 13 | N-Queens, Edit Distance, Alien Dictionary |
| 12 | Backtracking | 9 | Permutations, Subsets |
| 13 | Bit Manipulation | 9 | Single Number, Counting Bits |
| 14 | Heap | 7 | Top K Frequent, Merge K Lists |

### HackerRank (`hackerrank/`)

| Difficulty | Problems |
|------------|---------|
| Easy | Mini-Max Sum, FizzBuzz, Diagonal Difference, Counting Valleys, Caesar Cipher, Pangrams, and more |
| Medium | Encryption, Sherlock Valid String, Climbing Leaderboard |
| Hard | Array Manipulation (Difference Array) |

### Practical Engineering (`practical/`)

| # | Topic | What You Learn |
|---|-------|---------------|
| 01 | Dependency Management | go get, go mod tidy, versioning, workspaces |
| 02 | Build & Deploy | Cross-compilation, ldflags, version injection |
| 03 | Docker | Multi-stage builds, docker-compose, hot reload |
| 04 | Debugging | Delve, pprof, race detector, slog |
| 05 | Config | os.Getenv, JSON/YAML config, 12-factor |
| 06 | Concurrency Patterns | Worker pool, graceful shutdown |
| 07 | Performance Tuning | Escape analysis, sync.Pool, struct padding, GC pressure |
| 08 | Error Recovery | defer/recover, retry with backoff, context-aware retry |

---

## 🧪 Running Tests

```bash
go test ./...                              # run everything
go test -race ./...                        # with race detector (mandatory)
go test -cover ./...                       # with coverage
go test -v -run TestName ./path/           # single test, verbose
go test -bench=. -benchmem ./path/         # benchmarks + allocations
```

## 🔬 Performance Analysis

```bash
go build -gcflags='-m' ./path/             # escape analysis
go build -gcflags='-S' ./path/             # assembly output
go test -cpuprofile=cpu.out -bench=.       # CPU profile
go tool pprof cpu.out                      # analyze: top, list, web
GODEBUG=gctrace=1 ./app                    # GC trace
```

---

## 🔧 Setup

```bash
cd go-interview-prep
go build ./...        # verify everything compiles
go test ./...         # run all tests
```

> **Go version:** 1.25.7+ | **Module:** `gointerviewprep` | **Dependencies:** stdlib only
