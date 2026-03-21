# Go Senior Engineer — Training Sessions

> Each session teaches a concept, explains the **why**, shows the **how**, and includes **interview traps** you must know.

---

## Session 1: Interfaces & Type Systems ✅

### Why Interfaces? The Senior Engineer's Answer

In most languages, you implement an interface explicitly: `class Dog implements Animal`. Go is different — **interfaces are implicit**. A type satisfies an interface automatically if it has the right methods. This is called **structural typing** (also called "duck typing" — if it walks like a duck and quacks like a duck...).

**Why does this matter?** It decouples packages. Your code doesn't need to import the package that defines the interface — it just needs to implement the methods. This makes Go code extremely composable.

### Interfaces Under the Hood

An interface value has two parts internally:
```
interface value = (type pointer, data pointer)
```
When you assign a concrete type to an interface:
```go
var w io.Writer = os.Stdout
// w = (*os.File, pointer to stdout)
```
This is why a `nil` interface is NOT the same as an interface holding a `nil` pointer — a common gotcha:
```go
var p *MyType = nil
var i interface{} = p  // i is NOT nil! It has type info (*MyType) with a nil data pointer
fmt.Println(i == nil)  // false — TRAP!
```

### Key Rule: Accept Interfaces, Return Structs

```go
// ✅ GOOD — accepts interface, works with anything that has Write()
func SaveData(w io.Writer, data []byte) error { ... }

// ❌ BAD — forces caller to use *os.File, can't test with a mock
func SaveData(f *os.File, data []byte) error { ... }
```
This is the most important interface design principle. It makes functions testable and flexible.

### Type Switches — Runtime Type Inspection

When you receive `interface{}`, you need a type switch to safely extract the value:
```go
func Describe(i interface{}) string {
    switch v := i.(type) {  // v gets the concrete type inside each case
    case int:
        return fmt.Sprintf("int: %d", v)    // v is int here
    case string:
        return fmt.Sprintf("string: %s", v) // v is string here
    case bool:
        return fmt.Sprintf("bool: %t", v)
    default:
        return "unknown"
    }
}
```
**Interview trap**: What's the difference between `i.(int)` (type assertion) and `i.(type)` (type switch)?
- `i.(int)` — extracts and panics if wrong type. Use `v, ok := i.(int)` for safe version.
- `i.(type)` — only valid inside a `switch` statement.

---

## Session 2: Concurrency & The Go Memory Model ✅

### Goroutines — Not Threads

A goroutine starts with ~2KB of stack (vs ~1MB for an OS thread) and is managed by the Go runtime, not the OS. The Go scheduler multiplexes many goroutines onto a small number of OS threads (M:N scheduling).

**What this means in practice:**
- You can run 100,000 goroutines without issue.
- But goroutines are NOT free. Each has a stack, and the scheduler has overhead.
- **Senior insight**: don't create goroutines unboundedly. Use a Worker Pool (see below).

### The Go Memory Model — The Rule You Must Know

> **"If event A must happen before event B, you need a synchronization point."**

Without synchronization, the compiler and CPU can reorder your code in ways that break concurrent programs.

```go
// ❌ DATA RACE — no sync, one goroutine writes while another reads
x := 0
go func() { x = 1 }()
fmt.Println(x) // may print 0 or 1, or corrupt memory

// ✅ CORRECT — channel send/receive is a sync point
ch := make(chan int, 1)
go func() { x = 1; ch <- 1 }()
<-ch
fmt.Println(x) // guaranteed to print 1
```

**Always run `go test -race ./...` before shipping code.**

### WaitGroup Pattern

```go
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)              // BEFORE launching the goroutine
    go func(id int) {
        defer wg.Done()    // defer ensures it runs even on panic
        doWork(id)
    }(i)                   // pass i as argument — captures correct value
}
wg.Wait()
```
**Interview trap**: Why pass `i` as an argument instead of using it directly in the closure?
Because all goroutines would share the same `i` variable and by the time they run, the loop may have already finished — all goroutines would see the last value of `i`.

### Channels — Communication Over Shared Memory

Go's philosophy: **"Do not communicate by sharing memory; share memory by communicating."**

```
Unbuffered chan:  sender BLOCKS until receiver is ready. Pure synchronization.
Buffered chan:    sender only blocks when buffer is FULL. Decouples sender/receiver.
```

**Three channel patterns you must know:**

```go
// 1. Pipeline — chain stages
nums := Generate(5)    // produces 1,2,3,4,5
squares := Square(nums) // squares each, closes when input closes

// 2. Fan-In — merge multiple channels into one
merged := Merge(chanA, chanB) // one consumer reads all

// 3. Timeout — never block forever on a channel
select {
case v := <-ch:
    process(v)
case <-time.After(5 * time.Second):
    return errors.New("timed out")
}
```

### System Design: Worker Pool

**Problem**: You have 10,000 incoming HTTP requests, each needing a DB query. You can't create 10,000 goroutines — you'll exhaust DB connections.

**Solution**: Worker Pool — `N` workers, one job queue.

```
[Request 1] ─┐
[Request 2] ─┤──► [Job Queue (buffered chan)] ──► Worker 1
[Request 3] ─┤                               ──► Worker 2
[...10000 ] ─┘                               ──► Worker 3
```

```go
type WorkerPool struct {
    jobQueue chan Job        // buffered: absorbs bursts
    wg       sync.WaitGroup
    quit     chan struct{}   // for graceful shutdown
}

func (wp *WorkerPool) worker() {
    defer wp.wg.Done()
    for {
        select {
        case job, ok := <-wp.jobQueue:
            if !ok { return }  // channel closed = drain complete
            process(job)
        case <-wp.quit:
            return             // hard stop signal
        }
    }
}
```

**Two shutdown strategies:**
| Strategy | How | When to use |
|---|---|---|
| **Drain** | `close(jobQueue)` | Let workers finish all pending jobs |
| **Hard Stop** | `close(quit)` | Emergency stop, drop pending jobs |

---

## Session 3: Error Handling & Building Resilient Systems

### Go's Error Philosophy — Errors Are Values

In Java/Python, exceptions are thrown and caught, interrupting the normal flow. In Go, **errors are just values returned from functions**. You check them immediately. This is intentional — it forces you to think about every failure point.

```go
// Every function that can fail returns (value, error)
result, err := riskyOperation()
if err != nil {
    // handle it HERE, right now
    return fmt.Errorf("riskyOperation failed: %w", err)
}
// continue with result
```

**Why is this better than exceptions?**
- You know exactly which operations can fail by reading the signature.
- No hidden control flow.
- Forces callers to handle errors explicitly.

### Three Levels of Error Design

**Level 1: Simple string error** — for leaf-level, internal errors.
```go
return errors.New("cannot divide by zero")
return fmt.Errorf("key %q not found in map", key)
```

**Level 2: Sentinel errors** — predefined package-level errors for callers to check against.
```go
var ErrNotFound = errors.New("not found")
var ErrPermission = errors.New("permission denied")

// Caller checks:
if errors.Is(err, ErrNotFound) { ... }
```
**When to use**: When the caller needs to make a decision based on the error type.

**Level 3: Custom error types** — when you need to attach structured data.
```go
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s — %s", e.Field, e.Message)
}

// Caller extracts the data:
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Println("Invalid field:", ve.Field)
}
```
**When to use**: When the caller needs structured info (field name, HTTP status code, retry delay, etc.).

### Error Wrapping — Adding Context Without Losing Identity

Every time you return an error from a lower layer, **add context**:
```go
func GetUser(id int) (*User, error) {
    row, err := db.Query(...)
    if err != nil {
        // ✅ Wrap with %w — caller can still use errors.Is/errors.As on original
        return nil, fmt.Errorf("GetUser(id=%d): %w", id, err)
        
        // ❌ Don't do this — loses original error identity
        return nil, fmt.Errorf("GetUser failed: %v", err)
    }
}
```

**Error chain unwrapping:**
```
fmt.Errorf("A: %w", fmt.Errorf("B: %w", ErrNotFound))
        ↓
errors.Is(err, ErrNotFound) = true  ✅ — unwraps the whole chain
```

### `errors.Is` vs `errors.As` — The Critical Distinction

```go
// errors.Is — checks IDENTITY (is this exact error in the chain?)
errors.Is(err, ErrNotFound)        // true if ErrNotFound is anywhere in the chain

// errors.As — checks TYPE (is there an error of this type in the chain?)
var ve *ValidationError
errors.As(err, &ve)                // true if *ValidationError is anywhere in the chain
                                   // AND populates ve with that error's value
```

**Interview trap**: Why use `errors.Is` instead of `err == ErrNotFound`?
Because `err` might be a wrapped error. `errors.Is` unwraps the chain recursively. `==` only compares the outermost error.

### Exercises in `fundamentals/07_error_handling/exercises.go`

Let's now implement all 5 exercises. Here's what each one teaches:

| Exercise | Teaches |
|---|---|
| `Divide` | Basic error return |
| `Validate` | Custom error type (`*ValidationError`) |
| `SafeGet` | Error with context (`fmt.Errorf`) |
| `FindUser` | Sentinel errors (`errors.Is`) |
| `WrapError` | Error wrapping with `%w` |

