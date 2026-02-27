# 🗺️ Your Personal Learning Roadmap
## Go — From Beginner to Interview Ready

> **How to use this guide:** Follow the phases in order. Each phase
> has a "Why", a "What to do", and "How to verify you're done."
> Don't skip phases — each one builds on the last.

---

## ⚡ The Big Picture (read this first)

```
PHASE 1 ──► PHASE 2 ──► PHASE 3 ──► PHASE 4 ──► PHASE 5 ──► PHASE 6
 Language     Standard    Algorithm   LeetCode    HackerRank   Real World
 Basics       Library     Patterns    Problems    Problems     Engineering
 (2 weeks)    (1 week)    (1 week)    (4 weeks)   (1 week)     (1 week)
```

> Total: ~10 weeks at 1–2 hours/day. Adjust to your pace.

---

## ✅ PHASE 1 — Go Language Fundamentals
### Folder: `fundamentals/`
### Duration: ~2 weeks

This is your foundation. Everything else depends on it.
Read each `concepts.go` file slowly. **Don't just skim — type the code yourself.**

---

### 📅 Week 1 — Core Language

#### Day 1 → `fundamentals/01_basics/concepts.go`
**What you learn:** variables, constants, types, zero values, fmt printing

**Read in this order:**
1. Section 2: Variables — understand `:=` vs `var`
2. Section 3: Constants — understand `iota`
3. Section 4: Basic Types — understand `int`, `string`, `bool`, `float64`
4. Section 5: Zero Values — **this is unique to Go, very important**
5. Section 6: fmt verbs — `%v`, `%T`, `%d`, `%s`, `%f`

**Then do the exercises:**
```bash
# Open the exercise file
# File: fundamentals/01_basics/exercises.go
# Implement the TODO functions WITHOUT looking at solutions.go

# Check your work:
go test ./fundamentals/01_basics/...
```

**Key things to remember from Day 1:**
- `:=` is shorthand for `var x type = value` (only inside functions)
- Go NEVER implicitly converts types — you must do `float64(myInt)` explicitly
- Zero values mean variables are ALWAYS safe to use without initializing
- `len("hello")` returns BYTES, use `range` to iterate characters

---

#### Day 2 → `fundamentals/02_control_flow/concepts.go`
**What you learn:** if/else, switch, for loops, defer

**Key things that are different from other languages:**
```go
// Go has ONLY one loop keyword: for
// (no while, no do-while)

for i := 0; i < 10; i++ { }     // classic for
for i < 10 { }                  // while-style
for { }                          // infinite loop (use break to exit)
for i, v := range slice { }     // iterate over slice/array/string/map

// if with init statement (very Go-idiomatic):
if err := doSomething(); err != nil {
    // handle error
}

// switch — no fallthrough by default (opposite of C/Java)
switch day {
case "Mon", "Tue": fmt.Println("weekday")   // multiple values per case
case "Sat", "Sun": fmt.Println("weekend")
default:           fmt.Println("unknown")
}

// defer — runs when the surrounding function returns
defer fmt.Println("I run last")  // used for cleanup (close files, unlock)
```

**Exercise:** Write a FizzBuzz in Go using only `for` and `switch`. 
No if/else allowed.

---

#### Day 3 → `fundamentals/03_functions/concepts.go`
**What you learn:** function signatures, multiple returns, closures, variadic

**Key things that are different:**
```go
// Multiple return values — used everywhere in Go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}

result, err := divide(10, 2)
if err != nil { /* handle */ }

// Named return values (can be confusing, use sparingly)
func minMax(nums []int) (min, max int) {
    min, max = nums[0], nums[0]
    for _, n := range nums[1:] {
        if n < min { min = n }
        if n > max { max = n }
    }
    return  // "naked return" — returns named values
}

// Closures capture variables from outer scope
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}
counter := makeCounter()
counter() // 1
counter() // 2
```

**Exercise:** Write a function `filter(nums []int, fn func(int) bool) []int`
that returns only elements where fn returns true.

---

#### Day 4 → `fundamentals/04_pointers/concepts.go`
**What you learn:** `&` (address-of), `*` (dereference), nil pointers, pointer receivers

**Key insight — why pointers matter:**
```go
// Without pointer: function gets a COPY, original unchanged
func doubleWrong(x int) {
    x = x * 2  // changes the COPY, not the original
}

// With pointer: function gets the ADDRESS, can change original
func doubleRight(x *int) {
    *x = *x * 2  // * dereferences: "go to the address and change the value"
}

n := 5
doubleRight(&n)  // & takes the address of n
fmt.Println(n)   // 10 ✓

// When to use pointers:
// 1. When you want to MODIFY the original value
// 2. When the struct is LARGE (avoid copying)
// 3. When you need to represent "no value" (nil)
// 4. Method receivers (see Day 5)
```

---

#### Day 5 → `fundamentals/05_structs/concepts.go`
**What you learn:** struct definition, methods, embedding (Go's version of inheritance)

```go
type Person struct {
    Name string
    Age  int
}

// Value receiver — works on a COPY (use when not modifying)
func (p Person) Greet() string {
    return "Hi, I'm " + p.Name
}

// Pointer receiver — works on the ORIGINAL (use when modifying)
func (p *Person) Birthday() {
    p.Age++  // modifies the original Person
}

// Embedding — Go's composition over inheritance
type Employee struct {
    Person          // embedded — Employee gets all Person methods!
    Company string
}

e := Employee{Person: Person{Name: "Alice", Age: 30}, Company: "Go Corp"}
e.Greet()     // works! Promoted from Person
e.Birthday()  // works! Promoted from Person
```

**Rule of thumb:** If ANY method needs a pointer receiver, make ALL methods pointer receivers.

---

### 📅 Week 2 — Go's Unique Features

#### Day 6 → `fundamentals/06_interfaces/concepts.go`
**What you learn:** implicit interface implementation, polymorphism, type assertions

**This is one of Go's most powerful features:**
```go
// Interface = a contract (list of method signatures)
type Animal interface {
    Sound() string
    Name() string
}

// Dog implements Animal IMPLICITLY
// No "implements" keyword needed — if it has the methods, it satisfies the interface
type Dog struct{ name string }
func (d Dog) Sound() string { return "Woof" }
func (d Dog) Name() string  { return d.name }

type Cat struct{ name string }
func (c Cat) Sound() string { return "Meow" }
func (c Cat) Name() string  { return c.name }

// Now you can write ONE function that works for ANY Animal
func MakeSound(a Animal) {
    fmt.Printf("%s says %s\n", a.Name(), a.Sound())
}

MakeSound(Dog{name: "Rex"})  // Rex says Woof
MakeSound(Cat{name: "Whiskers"})  // Whiskers says Meow

// The MOST important interface in Go:
type error interface {
    Error() string  // any type with this method IS an error
}
```

---

#### Day 7 → `fundamentals/07_error_handling/concepts.go`
**What you learn:** Go's error pattern, custom errors, panic vs error

**Go error handling is different from exceptions — understand this deeply:**
```go
// In Go: errors are VALUES returned from functions
// You MUST check them (the compiler helps enforce this)

file, err := os.Open("myfile.txt")
if err != nil {
    // handle it — don't ignore!
    return fmt.Errorf("opening file: %w", err)  // %w wraps the error
}
defer file.Close()

// Check error type:
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Path error:", pathErr.Path)
}

// Sentinel errors (check with errors.Is):
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("File doesn't exist")
}

// Custom error type:
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s — %s", e.Field, e.Message)
}
```

**Rule:** Use `panic` ONLY for programming errors (like nil pointer you never expect).
Use `error` returns for ALL expected failure cases.

---

#### Day 8 → `fundamentals/08_arrays_slices/concepts.go`
**What you learn:** arrays vs slices, append, copy, 2D slices

**Slices are used EVERYWHERE. Understand them deeply:**
```go
// Array: fixed size, rarely used directly
var arr [5]int  // [0,0,0,0,0]

// Slice: dynamic, backed by an array — this is what you use daily
s := []int{1, 2, 3}
s = append(s, 4, 5)    // grows automatically

// CRITICAL: slice header = {pointer, length, capacity}
a := []int{1, 2, 3, 4, 5}
b := a[1:3]    // b = [2, 3] — shares the SAME underlying array!
b[0] = 99      // MODIFIES a too! a = [1, 99, 3, 4, 5]

// Safe copy:
c := make([]int, len(a))
copy(c, a)    // c is independent of a

// Slice tricks used in interviews:
s = append(s[:i], s[i+1:]...)  // delete element at index i
s = append(s, 0); copy(s[i+1:], s[i:]); s[i] = v  // insert at i
```

---

#### Day 9 → `fundamentals/09_maps/concepts.go`
**What you learn:** map CRUD, safe existence check, maps of slices

```go
// Always check existence — never assume a key exists
m := map[string]int{"a": 1, "b": 2}

val, ok := m["c"]    // ok=false, val=0 (zero value)
if !ok {
    fmt.Println("key not found")
}

// Maps are REFERENCE types — assigning copies the reference
m2 := m       // both m and m2 point to the SAME map!
m2["a"] = 99  // changes m too!

// Safe copy:
m3 := make(map[string]int)
for k, v := range m { m3[k] = v }

// IMPORTANT: map iteration order is RANDOM in Go (by design)
// Don't rely on order — sort keys if you need deterministic output
```

---

#### Day 10 → `fundamentals/10_goroutines/concepts.go`
#### Day 11 → `fundamentals/11_channels/concepts.go`
**What you learn:** concurrency — Go's killer feature

```go
// Goroutine: lightweight thread, started with 'go' keyword
go func() {
    fmt.Println("I run concurrently")
}()

// WaitGroup: wait for goroutines to finish
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Printf("worker %d done\n", id)
    }(i)
}
wg.Wait()

// Channel: communicate between goroutines (don't share memory!)
ch := make(chan int)
go func() { ch <- 42 }()  // send
val := <-ch                 // receive (blocks until value arrives)

// Buffered channel: doesn't block until buffer is full
ch := make(chan int, 3)
ch <- 1; ch <- 2; ch <- 3  // doesn't block

// Select: like switch but for channels
select {
case msg := <-ch1: fmt.Println("ch1:", msg)
case msg := <-ch2: fmt.Println("ch2:", msg)
case <-time.After(1 * time.Second): fmt.Println("timeout")
}
```

**Rule:** "Don't communicate by sharing memory; share memory by communicating."

---

#### Day 12 → `fundamentals/12_packages_modules/concepts.go`

**Quick build commands before moving on:**
```bash
go build ./...         # compile everything — 0 errors = you're ready
go test ./...          # run all tests — all pass = you're ready
go vet ./...           # static analysis — 0 issues = you're ready
```

---

## ✅ PHASE 2 — Standard Library
### Folder: `stdlib/`
### Duration: ~1 week (1 file per day)

Don't try to memorize everything. **Learn to look things up quickly.**
The goal is to know WHAT EXISTS so you know where to look.

| Day | File | Must-Know Things |
|-----|------|-----------------|
| 1 | `01_strings_strconv` | `strings.Contains/HasPrefix/Split/Join/TrimSpace`, `strconv.Atoi/Itoa` |
| 2 | `02_sort` | `sort.Ints/Strings`, `sort.Slice`, `sort.Search` |
| 3 | `03_builtins` | `len/cap`, `make/new`, `append/copy`, `delete`, `panic/recover` |
| 4 | `04_io_files` | `os.ReadFile/WriteFile`, `bufio.Scanner`, `fmt.Fprintf` |
| 5 | `05_encoding_json` | `json.Marshal/Unmarshal`, struct tags, `json.NewEncoder` |
| 6 | `06_math` | `math.Max/Min/Abs`, `math.Sqrt`, `math/rand` |
| 7 | `07_testing` | table-driven tests, `t.Run`, benchmarks |

**After each file, write a small program that uses what you learned.**

---

## ✅ PHASE 3 — Algorithm Patterns
### Folder: `patterns/templates.go`
### Duration: ~1 week

**Read the templates file ONCE thoroughly, then memorize these 8 patterns.**
Every LeetCode/HackerRank problem maps to one of these:

| # | Pattern | When to use | Time |
|---|---------|-------------|------|
| 1 | **Two Pointers** | Sorted array, find pair/triplet | O(n) |
| 2 | **Sliding Window** | Subarray/substring with constraint | O(n) |
| 3 | **Binary Search** | Sorted array, monotonic function | O(log n) |
| 4 | **BFS** | Shortest path, level-order traversal | O(V+E) |
| 5 | **DFS** | Explore all paths, tree traversal | O(V+E) |
| 6 | **Dynamic Programming** | Optimal substructure, overlapping subproblems | O(n²) |
| 7 | **Monotonic Stack** | Next greater/smaller element | O(n) |
| 8 | **Union-Find** | Connected components, cycle detection | O(α(n)) |

**For each pattern, answer these 3 questions:**
1. What does the template look like?
2. What problems does it solve?
3. What are the edge cases?

---

## ✅ PHASE 4 — LeetCode Problems
### Folder: `leetcode/`
### Duration: ~4 weeks

**The right way to use this repo for LeetCode:**

```
Step 1: READ the problem comment at the top of the function
Step 2: UNDERSTAND the approach section (don't skip this!)
Step 3: CLOSE the file and try to implement it YOURSELF
Step 4: RUN the tests — if they fail, debug before looking at the solution
Step 5: Compare your solution with the one in the file
Step 6: UNDERSTAND every line of the solution
```

**Weekly Schedule:**

### Week 1 — Arrays, Strings, Linked Lists (Easy focus)
```bash
# Start here — these are ALWAYS asked in interviews
go test ./leetcode/01_arrays/...     # arrays
go test ./leetcode/02_strings/...    # strings  
go test ./leetcode/03_linked_list/... # linked list
```
**Must solve without looking:** Two Sum, Valid Anagram, Reverse Linked List

### Week 2 — Stacks, Binary Search, Sliding Window
```bash
go test ./leetcode/04_stacks_queues/...
go test ./leetcode/05_binary_search/...
go test ./leetcode/06_sliding_window/...
```
**Must solve without looking:** Valid Parentheses, Binary Search, Longest Substring Without Repeat

### Week 3 — Trees, Graphs
```bash
go test ./leetcode/07_trees/...
go test ./leetcode/08_graphs/...
```
**Must solve without looking:** Max Depth, Level Order Traversal, Number of Islands

### Week 4 — Dynamic Programming, Two Pointers, Hard
```bash
go test ./leetcode/09_dynamic_prog/...
go test ./leetcode/10_two_pointers/...
go test ./leetcode/11_hard/...
```
**Must solve without looking:** Climbing Stairs, Coin Change, 3Sum

---

## ✅ PHASE 5 — HackerRank
### Folder: `hackerrank/`
### Duration: ~1 week

HackerRank is great for **interview warm-up** and **contest-style** problems.
The style is slightly different from LeetCode — input/output focused.

**Recommended order:**
1. All Easy problems first (they're fast, confidence builders)
2. Medium problems 
3. Array Manipulation (Hard — great for difference array technique)

```bash
go test ./hackerrank/... -v   # see each test name passing
```

---

## ✅ PHASE 6 — Practical Engineering
### Folder: `practical/`
### Duration: ~1 week

This is what separates interview candidates from working engineers.

| Day | File | Key Command to Try |
|-----|------|--------------------|
| 1 | `01_dependency_management` | `go get github.com/joho/godotenv` |
| 2 | `02_build_run_deploy` | `go build -ldflags "-s -w" -o myapp .` |
| 3 | `03_docker` | Build and run the Dockerfile in this guide |
| 4 | `04_debugging` | Install Delve, set a breakpoint in GoLand |
| 5 | `05_config_env_json_yaml` | `go test ./practical/05_config_env_json_yaml/...` |

---

## 🧪 How to Test Your Progress

### Run a single test:
```bash
go test -run TestTwoSum ./leetcode/01_arrays/
```

### Run a whole category:
```bash
go test ./leetcode/01_arrays/...
```

### Run everything and see pass rate:
```bash
go test ./... 2>&1
```

### Run with verbose output (see each test name):
```bash
go test -v ./leetcode/01_arrays/...
```

### Check coverage:
```bash
go test -cover ./leetcode/01_arrays/...
```

### Build to verify no compile errors:
```bash
go build ./...
```

---

## 🚦 Self-Assessment Checkpoints

After Phase 1, you should be able to answer YES to all:
- [ ] I know the difference between `:=` and `var`
- [ ] I know what zero values are and why they matter
- [ ] I can write a function with multiple return values
- [ ] I understand pointer receivers vs value receivers
- [ ] I know when to use an interface
- [ ] I understand Go error handling (return error, check nil)
- [ ] I know the difference between array and slice
- [ ] I can create and iterate a map safely
- [ ] I understand what a goroutine is and how to use WaitGroup
- [ ] I know what a channel is and can do basic send/receive

After Phase 4, you should be able to answer YES to all:
- [ ] I can solve Easy LeetCode problems in < 10 minutes
- [ ] I can solve Medium LeetCode problems in < 25 minutes
- [ ] I know which pattern to apply within 2 minutes of reading a problem
- [ ] I can explain my time and space complexity for every solution

---

## 💡 Tips for Success

### 1. Type the code, don't copy-paste
Reading code and understanding code are different skills.
Writing code builds muscle memory and forces you to think.

### 2. Understand before moving on
If you don't understand something, don't skip it.
Go is consistent — if you understand the basics, everything else makes sense.

### 3. Test-driven approach
Every function in this repo has a test.
Red → Green → Refactor. Run tests constantly.

### 4. Read the error messages
Go's error messages are very clear. Learn to read them:
```
./main.go:15:9: undefined: foo        → you used something that doesn't exist
./main.go:8:2: imported and not used  → remove unused import
./main.go:12:5: cannot use x (type int) as type string → type mismatch
```

### 5. Use `go doc`
```bash
go doc fmt.Sprintf          # see docs for a function
go doc strings              # see all functions in a package
go doc -all strings.Builder # full docs for a type
```

### 6. Interview Pattern
When given a problem in an interview:
```
1. Read out loud and confirm understanding (2 min)
2. Identify the pattern (Binary Search? DP? Two Pointers?)
3. State your approach before coding (2 min)
4. Code from the outside in (function signature first)
5. Walk through an example as you code
6. State the time/space complexity at the end
```

---

## 📌 Your First 30 Minutes — RIGHT NOW

Open your terminal and do this:

```bash
cd C:\Users\samte\GolandProjects\go-interview-prep

# 1. Confirm the project builds
go build ./...

# 2. Confirm all tests pass  
go test ./...

# 3. Open the first file and read it
# File: fundamentals/01_basics/concepts.go  ← START HERE

# 4. Try the exercises
# File: fundamentals/01_basics/exercises.go  ← IMPLEMENT THESE

# 5. Run your solutions
go test ./fundamentals/01_basics/...
```

**That's your first step. Open `fundamentals/01_basics/concepts.go` right now.**

---

## 📊 Progress Tracker

Copy this and keep it somewhere visible:

```
PHASE 1 — Fundamentals
[ ] 01_basics          (Day 1)
[ ] 02_control_flow    (Day 2)
[ ] 03_functions       (Day 3)
[ ] 04_pointers        (Day 4)
[ ] 05_structs         (Day 5)
[ ] 06_interfaces      (Day 6)
[ ] 07_error_handling  (Day 7)
[ ] 08_arrays_slices   (Day 8)
[ ] 09_maps            (Day 9)
[ ] 10_goroutines      (Day 10)
[ ] 11_channels        (Day 11)
[ ] 12_packages_modules(Day 12)

PHASE 2 — Standard Library
[ ] 01_strings_strconv
[ ] 02_sort
[ ] 03_builtins
[ ] 04_io_files
[ ] 05_encoding_json
[ ] 06_math
[ ] 07_testing

PHASE 3 — Patterns
[ ] Read patterns/templates.go
[ ] Memorize 8 patterns

PHASE 4 — LeetCode
[ ] 01_arrays (10 problems)
[ ] 02_strings (10 problems)
[ ] 03_linked_list (10 problems)
[ ] 04_stacks_queues (8 problems)
[ ] 05_binary_search (10 problems)
[ ] 06_sliding_window (7 problems)
[ ] 07_trees (12 problems)
[ ] 08_graphs (8 problems)
[ ] 09_dynamic_prog (10 problems)
[ ] 10_two_pointers (10 problems)
[ ] 11_hard (12 problems)

PHASE 5 — HackerRank
[ ] All Easy (10 problems)
[ ] All Medium (4 problems)
[ ] Hard (1 problem)

PHASE 6 — Practical
[ ] 01_dependency_management
[ ] 02_build_run_deploy
[ ] 03_docker
[ ] 04_debugging
[ ] 05_config_env_json_yaml
```

