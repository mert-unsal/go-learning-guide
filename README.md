# 🚀 Go Interview Prep — From Zero to Coding Exam Ready

> **Folder:** `go-interview-prep`  |  **Module:** `gointerviewprep`

Welcome! This project is a **structured, hands-on Go curriculum** covering language fundamentals, algorithm patterns, LeetCode, HackerRank, and real-world Go engineering practices.

---

## 📁 Project Structure

```
go-interview-prep/
├── fundamentals/   ← Go language core — read these in order
├── stdlib/         ← Standard library deep-dives
├── patterns/       ← Algorithm pattern templates (reference)
├── leetcode/       ← Solved LeetCode problems with explanations + tests
├── hackerrank/     ← HackerRank problems with full solutions + tests
├── practical/      ← Real-world Go: deps, build, Docker, debug, config
└── utils/          ← Shared helpers
```

---

## 📚 Learning Roadmap

### Track 1: Go Fundamentals (`fundamentals/`)

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
| 03 | `builtins` | make, new, append, copy, len/cap, delete, panic/recover |
| 04 | `io_files` | os.ReadFile/WriteFile, bufio.Scanner, io.Reader/Writer |
| 05 | `encoding_json` | json.Marshal/Unmarshal, struct tags, Encoder/Decoder |
| 06 | `math` | math, rand, big numbers |
| 07 | `testing` | Table-driven tests, subtests, benchmarks |

### Track 3: Algorithm Patterns (`patterns/`)

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

### Track 4: LeetCode Problems (`leetcode/`)

| # | Category | Difficulty | Problems |
|---|----------|-----------|---------|
| 01 | Arrays | Easy/Med | Two Sum, Max Profit, Product Except Self, Contains Dup, Max Subarray, Merge Sorted, Find Disappeared, Rotate, Find Min Rotated, Subarray Sum K |
| 02 | Strings | Easy/Med | Valid Anagram, Longest Substring, Valid Palindrome, LCP, Reverse Words, First Unique Char, Roman to Int, Count and Say, Group Anagrams, Encode & Decode |
| 03 | Linked List | Easy/Med | Reverse, Merge Sorted, Has Cycle, Remove Nth, Middle, Palindrome, Intersection, Add Two Numbers, Copy Random, Reorder |
| 04 | Stacks & Queues | Easy/Med/Hard | Valid Parens, Min Stack, Daily Temps, Eval RPN, Queue via Stacks, Next Greater, Decode String, Largest Rectangle |
| 05 | Binary Search | Easy/Med/Hard | Search Rotated, Find Min, Classic BS, Search Matrix, Find Peak, Guess Number, Count Negatives, Koko Bananas, First+Last Position, Median Two Arrays |
| 06 | Sliding Window | Easy/Med/Hard | Max Average, Min Window, Permutation in String, Fruit Baskets, Char Replacement, Max Points Cards, Min Subarray Sum |
| 07 | Trees | Easy/Med | Inorder, Max Depth, Level Order, LCA, Symmetric, Validate BST, Path Sum, Invert, Diameter, Build Tree, Right Side View, Kth Smallest |
| 08 | Graphs | Easy/Med | Islands, Course Schedule, Clone Graph, Flood Fill, Valid Path, Connected Components, Rotting Oranges, Pacific Atlantic |
| 09 | Dynamic Prog | Easy/Med | Climb Stairs, Coin Change, House Robber, Unique Paths, LCS, Min Cost Stairs, LIS, Word Break, Jump Game, Partition Equal Subset |
| 10 | Two Pointers | Easy/Med | Container Water, 3Sum, Trapping Rain, Move Zeros, Remove Dups, Triangle Number, Sorted Squares, 4Sum, Min Diff K Scores, Bag of Tokens |
| 11 | **Hard** | Hard | Merge K Lists, Trapping Rain, Word Ladder, Longest Valid Parens, Jump Game II, N-Queens, Serialize/Deserialize Tree, Min Window, Alien Dictionary, Regex Matching, Edit Distance, Job Scheduling |

### Track 5: HackerRank Problems (`hackerrank/`)

| Difficulty | Problems |
|------------|---------|
| Easy | Mini-Max Sum, FizzBuzz, Diagonal Difference, Counting Valleys, Sales by Match (Socks), Jumping on Clouds, Repeated String, Caesar Cipher, Pangrams, Mark and Toys |
| Medium | Encryption, Sherlock and the Valid String, Climbing the Leaderboard, Almost Sorted |
| Hard | Array Manipulation (Difference Array technique) |

### Track 6: Practical Go Engineering (`practical/`)

| # | Topic | What You Learn |
|---|-------|---------------|
| 01 | `dependency_management` | `go get`, `go mod tidy`, versioning, workspaces, private modules |
| 02 | `build_run_deploy` | `go build`, cross-compilation, ldflags, version injection, systemd |
| 03 | `docker` | Single-stage & multi-stage Dockerfiles, docker-compose, hot reload |
| 04 | `debugging` | Delve debugger, GoLand/VSCode setup, pprof, race detector, slog |
| 05 | `config_env_json_yaml` | `os.Getenv`, `.env` files, JSON config, YAML config, 12-factor pattern |

---

## 🛠️ Practical Quick Reference

### Adding a dependency
```bash
go get github.com/some/package@latest      # add latest
go get github.com/some/package@v1.2.3      # pin version
go mod tidy                                # clean up unused deps
```

### Build & Run
```bash
go run .                                   # dev mode (no binary)
go build -o myapp .                        # build binary
go build -ldflags "-s -w" -o myapp .       # stripped production binary

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o myapp-linux .
```

### Docker
```bash
docker build -t myapp:latest .             # build image
docker run -p 8080:8080 myapp:latest       # run container
docker run --env-file .env myapp:latest    # with .env file
docker compose up --build                  # full stack local dev
```

### Debug with Delve
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug .                                # start debug session
# Inside dlv: break main.go:42 → continue → print x → next → step
```

### Read Environment Variables
```go
port := os.Getenv("PORT")                         // "" if not set
port, ok := os.LookupEnv("PORT")                  // ok=false if missing
port = GetEnvOrDefault("PORT", "8080")            // with fallback
```

### Read JSON Config
```go
data, _ := os.ReadFile("config.json")
var cfg AppConfig
json.Unmarshal(data, &cfg)
fmt.Println(cfg.App.Port)  // → 8080
```

### Read YAML Config
```go
// go get gopkg.in/yaml.v3
data, _ := os.ReadFile("config.yaml")
var cfg AppConfigYAML
yaml.Unmarshal(data, &cfg)
```

---

## 🧪 Running Tests

```bash
go test ./...                              # run everything
go test ./leetcode/01_arrays/...           # one category
go test -run TestTwoSum ./leetcode/01_arrays/
go test -v ./hackerrank/...               # verbose
go test -cover ./...                      # with coverage
go test -race ./...                       # race detector
go test -count=1 ./...                    # bypass cache
```

---

## 💡 How to Use This Project

1. **Fundamentals first** — read `concepts.go` in each `fundamentals/` package in order
2. **Stdlib** — focus on `03_builtins`, `04_io_files`, `05_encoding_json`
3. **Patterns** — read `patterns/templates.go` before solving problems
4. **LeetCode** — read the problem comment → understand approach → implement → `go test`
5. **HackerRank** — same approach, great for interview warm-up
6. **Practical** — read `practical/` packages to understand real Go engineering

---

## 🔧 Setup

```bash
cd go-interview-prep
go build ./...        # verify everything compiles
go test ./...         # verify all tests pass
```

> Go version: 1.25.7+  |  No external dependencies (standard library only)

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

| # | Category | Difficulty | Problems Covered |
|---|----------|-----------|-----------------|
| 01 | Arrays | Easy/Medium | Two Sum, Max Profit, Product Except Self, Contains Duplicate, Max Subarray, **Merge Sorted Array, Find Disappeared Numbers, Rotate Array, Find Min Rotated, Subarray Sum Equals K** |
| 02 | Strings | Easy/Medium | Valid Anagram, Longest Substring, Valid Palindrome, LCP, Reverse Words, **First Unique Char, Roman to Integer, Count and Say, Group Anagrams, Encode & Decode Strings** |
| 03 | Linked List | Easy/Medium | Reverse, Merge Sorted, Has Cycle, Remove Nth, Middle Node, **Palindrome List, Intersection of Two Lists, Add Two Numbers, Copy List with Random Pointer, Reorder List** |
| 04 | Stacks & Queues | Easy/Medium/Hard | Valid Parentheses, Min Stack, Daily Temperatures, Eval RPN, **Queue using Stacks, Next Greater Element, Decode String, Largest Rectangle in Histogram** |
| 05 | Binary Search | Easy/Medium/Hard | Search Rotated, Find Min Rotated, Classic BS, Search Matrix, Find Peak, **Guess Number, Count Negatives, Koko Eating Bananas, Find First & Last Position, Median of Two Sorted Arrays** |
| 06 | Sliding Window | Easy/Medium/Hard | Max Average, Min Window Substring, Permutation in String, Fruit Baskets, **Longest Char Replacement, Max Points from Cards, Min Size Subarray Sum** |
| 07 | Trees | Easy/Medium | Inorder, Max Depth, Level Order, LCA, Symmetric, Validate BST, **Path Sum, Invert Tree, Diameter, Build from Pre+Inorder, Right Side View, Kth Smallest** |
| 08 | Graphs | Easy/Medium | Number of Islands, Course Schedule, Clone Graph, Flood Fill, **Valid Path, Connected Components, Rotting Oranges, Pacific Atlantic Water Flow** |
| 09 | Dynamic Prog | Easy/Medium | Climb Stairs, Coin Change, House Robber, Unique Paths, LCS, **Min Cost Stairs, LIS, Word Break, Jump Game, Partition Equal Subset Sum** |
| 10 | Two Pointers | Easy/Medium | Container Water, 3Sum, Trapping Rain, Move Zeros, Remove Dups, **Valid Triangle Number, Sorted Squares, 4Sum, Min Difference K Scores, Bag of Tokens** |
| 11 | **Hard** | Hard | Merge K Lists, Trapping Rain Water, Word Ladder, Longest Valid Parentheses, Jump Game II, N-Queens, Serialize/Deserialize Tree, Min Window Substring, Alien Dictionary, Regex Matching, Edit Distance, Job Scheduling |

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

> Go version: 1.25.7+  |  No external dependencies
