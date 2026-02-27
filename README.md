# 🚀 Go Interview Prep — From Zero to Coding Exam Ready

> **Folder:** `go-interview-prep`  |  **Module:** `gointerviewprep`

Welcome! This project is a **structured, hands-on Go curriculum** designed to take you from beginner to confidently passing LeetCode Easy/Medium problems and writing idiomatic Go.

---

## 📁 Project Structure

```
go-interview-prep/
├── fundamentals/   ← Go language core — read these in order
├── stdlib/         ← Standard library deep-dives
├── patterns/       ← Algorithm pattern templates (reference)
├── leetcode/       ← Solved problems with explanations + tests
└── utils/          ← Shared helpers (ListNode, TreeNode, etc.)
```

---

## 📚 Learning Roadmap

### Track 1: Go Fundamentals (`fundamentals/`)
Work through each package **in order**. Each package contains:
- `concepts.go` — Heavily commented examples you can run
- `exercises.go` — Problems for you to solve *(01_basics only so far)*
- `solutions.go` — Reference solutions with explanations

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

### Track 2: Standard Library (`stdlib/`)
| # | Package | Topics |
|---|---------|--------|
| 01 | `strings_strconv` | strings, Builder, strconv, unicode |
| 02 | `sort` | sort.Slice, sort.Search, custom Less |
| 03 | `builtins` | **make, new, append, copy, len/cap, delete, conversions, panic/recover** |
| 04 | `io_files` | **os.ReadFile/WriteFile, bufio.Scanner, buffered writes, io.Reader/Writer, filepath** |
| 05 | `encoding_json` | **json.Marshal/Unmarshal, struct tags, Encoder/Decoder, config files** |
| 06 | `math` | math, rand, big numbers |
| 07 | `testing` | Table-driven tests, subtests, benchmarks |

### Track 3: Algorithm Patterns (`patterns/`)
One file with fully commented reference scaffolds:

| Pattern | Key Concept |
|---------|-------------|
| Binary Search | `left <= right`, left-bound variant |
| Sliding Window | expand right, shrink left |
| Two Pointers | converging from both ends |
| BFS | queue + visited + level tracking |
| DFS | recursive closure + iterative stack |
| Dynamic Programming | top-down memoization, bottom-up tabulation |
| Monotonic Stack | indices on stack, pop when invariant breaks |
| Union-Find | path compression + union by rank |
| Heap | `container/heap` interface implementation |

### Track 4: LeetCode Problems (`leetcode/`)

| # | Category | Problems Covered |
|---|----------|-----------------|
| 01 | Arrays | Two Sum, Max Profit, Product Except Self, Contains Duplicate, Max Subarray |
| 02 | Strings | Valid Anagram, Longest Substring, Valid Palindrome, LCP, Reverse Words |
| 03 | Linked List | Reverse, Merge Sorted, Has Cycle, Remove Nth, Middle Node |
| 04 | Stacks & Queues | Valid Parentheses, Min Stack, Daily Temperatures, Eval RPN |
| 05 | Binary Search | Search Rotated, Find Min Rotated, Classic BS, Search Matrix, Find Peak |
| 06 | Sliding Window | Max Average, Min Window Substring, Permutation in String, Fruit Baskets |
| 07 | Trees | Inorder, Max Depth, Level Order, LCA, Symmetric, Validate BST |
| 08 | Graphs | Number of Islands, Course Schedule, Clone Graph, Flood Fill |
| 09 | Dynamic Prog | Climb Stairs, Coin Change, House Robber, Unique Paths, LCS |
| 10 | Two Pointers | Container Water, 3Sum, Trapping Rain, Move Zeros, Remove Dups |

---

## 🧪 Running Tests

```bash
# Run all tests
go test ./...

# Run tests for one category
go test ./leetcode/01_arrays/...

# Run a specific test
go test -run TestTwoSum ./leetcode/01_arrays/

# Run with verbose output
go test -v ./leetcode/07_trees/...

# Run with coverage
go test -cover ./...
```

---

## 💡 How to Use This Project

1. **Start with `fundamentals/01_basics`** — read `concepts.go`, then try `exercises.go`
2. **Move to stdlib** — focus on `03_builtins`, `04_io_files`, `05_encoding_json` for real-world Go
3. **Study `patterns/templates.go`** — understand each pattern before solving problems
4. **Solve LeetCode problems** — read the problem comment, understand the approach, implement, then check tests pass
5. **Run `go test ./...`** often to validate your work

---

## 🔧 Setup

```bash
# Clone / open the project
cd go-interview-prep

# Verify everything builds
go build ./...

# Run all tests
go test ./...
```

> Go version: 1.21+  |  No external dependencies
