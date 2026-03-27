# Deep Dive: Go Error Handling — Patterns, Internals & Enterprise Strategy

> How Go represents errors at runtime, how wrapping builds linked chains,
> how errors.Is/As walk those chains, and how to architect error handling
> across service layers in production.

---

## Table of Contents

1. [The `error` Interface — Under the Hood](#1-the-error-interface--under-the-hood)
2. [Creating Errors — The Four Ways](#2-creating-errors--the-four-ways)
3. [Sentinel Errors](#3-sentinel-errors)
4. [Custom Error Types](#4-custom-error-types)
5. [Error Wrapping — The `%w` Verb](#5-error-wrapping--the-w-verb)
6. [`errors.Is()` — Sentinel Matching Through the Chain](#6-errorsis--sentinel-matching-through-the-chain)
7. [`errors.As()` — Type Extraction Through the Chain](#7-errorsas--type-extraction-through-the-chain)
    - [7b. Type Assertions — The Runtime Machinery](#7b-type-assertions--the-runtime-machinery-behind-error-inspection)
    - [7c. Why Errors Are Interfaces, Not Strings](#7c-why-errors-are-interfaces-not-strings)
    - [7d. Java-to-Go Error Model Mental Bridge](#7d-java-to-go-error-model-mental-bridge)
8. [`errors.Join()` (Go 1.20+)](#8-errorsjoin-go-120)
9. [Enterprise Error Handling Strategy](#9-enterprise-error-handling-strategy)
10. [Panic and Recover](#10-panic-and-recover)
11. [Error Handling Anti-Patterns](#11-error-handling-anti-patterns)
12. [Error Handling in Concurrent Code](#12-error-handling-in-concurrent-code)
13. [Performance Characteristics](#13-performance-characteristics)
14. [Quick Reference Card](#14-quick-reference-card)

---

## 1. The `error` Interface — Under the Hood

The `error` type is not a keyword — it's a **plain interface** in the `builtin` package:

```go
type error interface {
    Error() string
}
```

One method. Any type with `Error() string` satisfies it.

### Runtime Representation

Since `error` has one method, it is a **non-empty interface**. At runtime, every `error`
value is a `runtime.iface` — the same 16-byte two-word struct from the interface deep dive:

```
error value at runtime = runtime.iface (16 bytes)
┌──────────────────────┬──────────────────────┐
│  tab  *itab          │  data unsafe.Pointer │
│  {                   │  (pointer to the     │
│   inter: error       │   concrete value:    │
│   _type: *errorString│   *errorString,      │
│   fun[0]: Error()    │   *PathError, etc.)  │
│  }                   │                      │
└──────────────────────┴──────────────────────┘
```

**Key connection:** The `itab` is cached globally per (error, concrete type) pair.
Method dispatch through `error` is an indirect call via `itab.fun[0]` — cannot be inlined.

### The Zero Value: nil = Success

```
var err error   // iface{nil, nil} — both words zero
err == nil → TRUE → this IS the success case in Go
```

### Why Not Exceptions?

Rob Pike and the Go team chose explicit error returns because:

1. **Exceptions hide control flow** — Go's `(result, error)` makes the error path visible
2. **Exception handling is lazy** — try/catch swallows errors; Go forces per-call handling
3. **Stack unwinding is expensive** — Go errors are just values, no runtime cost until wrapped
4. **"Errors are values"** — you can store them, pass them, aggregate them in slices

**Source:** Go FAQ ("Why does Go not have exceptions?"), Rob Pike's "Errors are Values" (2015).

---

## 2. Creating Errors — The Four Ways

### (a) `errors.New("msg")` — The Simplest Error

Returns `*errorString` (unexported struct):

```go
// errors/errors.go
func New(text string) error { return &errorString{text} }
type errorString struct{ s string }
func (e *errorString) Error() string { return e.s }
```

```
err (iface):
┌──────────────────┬──────────────────┐
│ tab: itab{       │ data ────────────┼──► errorString{s: "connection refused"}
│   error,         │                  │    (heap allocated — 1 alloc)
│   *errorString}  │                  │
└──────────────────┴──────────────────┘
```

### (b) `fmt.Errorf("%w", err)` — Wrapping with Context

With `%w`, returns `*fmt.wrapError`:

```go
// fmt/errors.go (simplified)
type wrapError struct {
    msg string
    err error   // the wrapped error
}
func (e *wrapError) Error() string { return e.msg }
func (e *wrapError) Unwrap() error { return e.err }
```

### (c) Custom Error Types — Struct with `Error()`

```go
type ValidationError struct{ Field, Message string }
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s — %s", e.Field, e.Message)
}
```

### (d) Sentinel Errors — Package-Level Variables

```go
var ErrNotFound = errors.New("not found")
var ErrTimeout  = errors.New("operation timed out")
```

### When to Use Each

```
┌─────────────────────┬──────────────────────────────────────────────────┐
│ Method              │ When to Use                                      │
├─────────────────────┼──────────────────────────────────────────────────┤
│ errors.New()        │ Simple, one-off errors. Not part of API contract.│
├─────────────────────┼──────────────────────────────────────────────────┤
│ fmt.Errorf("%w")    │ Adding context while preserving the cause.       │
├─────────────────────┼──────────────────────────────────────────────────┤
│ Custom error type   │ Domain errors carrying structured data.          │
├─────────────────────┼──────────────────────────────────────────────────┤
│ Sentinel error      │ Well-known conditions callers must check for.    │
└─────────────────────┴──────────────────────────────────────────────────┘
```

---

## 3. Sentinel Errors

A sentinel error is a **package-level variable** representing a specific, known condition.

```go
var ErrNotFound   = errors.New("not found")       // Convention: prefix with "Err"
var ErrPermission = errors.New("permission denied")
```

### Standard Library Examples

```go
io.EOF                    // end of stream — a signal, not really an "error"
sql.ErrNoRows             // query returned zero rows
os.ErrNotExist            // file/dir does not exist
context.Canceled          // context was canceled
context.DeadlineExceeded  // context deadline passed
```

### How to Check: `errors.Is()`, Not `==`

```go
if errors.Is(err, io.EOF) { }     // GOOD — works through wrapping chains
if err == io.EOF { }               // BAD — breaks if err is wrapped
if err.Error() == "EOF" { }        // TERRIBLE — fragile, breaks on typo
```

### API Contract Implications

Sentinel errors are part of your **public API.** Once exported, callers depend on them.
You can't rename, change meaning, or remove without breaking downstream code.

```
 Package "repo"                     Caller code
 ─────────────────                  ──────────────────────────
 var ErrNotFound = errors.New(...)  if errors.Is(err, repo.ErrNotFound) {
 This is an API contract ──────────►    // handle 404
                                    }
```

---

## 4. Custom Error Types

When an error needs to carry **structured data** beyond a message string:

```go
type ValidationError struct{ Field, Message string }
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed: field %q — %s", e.Field, e.Message)
}

type NotFoundError struct{ Resource, ID string }
func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s %q not found", e.Resource, e.ID)
}
```

### Interface Satisfaction & Extraction

```
Compiler sees: var err error = &ValidationError{...}
  Does *ValidationError have Error() string?  ✅
  └─ Satisfies error — iface{itab, data} is constructed
```

```go
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Printf("field: %s, message: %s\n", ve.Field, ve.Message)  // structured access
}
```

**Rule:** Use sentinels when callers only need the category. Use custom types when
callers need to extract data (field name, status code, resource ID).

---

## 5. Error Wrapping — The `%w` Verb

### The Core Mechanism

```go
original := errors.New("connection refused")
wrapped  := fmt.Errorf("database: query users: %w", original)
// wrapped.Unwrap() returns original — chain is preserved
```

### The Wrapping Chain as a Linked List

```go
err1 := errors.New("connection refused")
err2 := fmt.Errorf("postgres: query failed: %w", err1)
err3 := fmt.Errorf("user service: get user 42: %w", err2)
err4 := fmt.Errorf("GET /api/users/42: %w", err3)
```

```
err4 (surface)
┌───────────────────────────────────────────┐
│ wrapError{                                │
│   msg: "GET /api/users/42: ..."           │
│   err ────────────────────────────────┐   │
│ }                                     │   │
└───────────────────────────────────────┼───┘
                                        ▼
err3  ┌─────────────────────────────────────┐
      │ wrapError{                          │
      │   msg: "user service: get user 42:…"│
      │   err ──────────────────────────┐   │
      │ }                               │   │
      └─────────────────────────────────┼───┘
                                        ▼
err2  ┌─────────────────────────────────────┐
      │ wrapError{                          │
      │   msg: "postgres: query failed: …"  │
      │   err ──────────────────────────┐   │
      │ }                               │   │
      └─────────────────────────────────┼───┘
                                        ▼
err1  ┌─────────────────────────────────────┐
      │ errorString{s: "connection refused"}│
      │ (no Unwrap — chain ends here)       │
      └─────────────────────────────────────┘
```

### `%w` vs `%v` — The Critical Difference

```go
fmt.Errorf("query: %w", err)   // WRAPS: chain preserved, errors.Is/As work
fmt.Errorf("query: %v", err)   // FORMATS: chain BROKEN, original err lost
```

Use `%w` when callers need to inspect the cause. Use `%v` to intentionally **hide**
internal error types at an API boundary.

---

## 6. `errors.Is()` — Sentinel Matching Through the Chain

Walks the entire `Unwrap()` chain looking for a target match.

### The Algorithm

```go
// errors/wrap.go (simplified)
func Is(err, target error) bool {
    for {
        if err == target { return true }                     // pointer equality
        if x, ok := err.(interface{ Is(error) bool }); ok {
            if x.Is(target) { return true }                  // custom Is()
        }
        switch x := err.(type) {
        case interface{ Unwrap() error }:
            err = x.Unwrap()                                 // walk deeper
        case interface{ Unwrap() []error }:
            for _, e := range x.Unwrap() {
                if Is(e, target) { return true }             // recurse branches
            }
            return false
        default:
            return false                                     // chain exhausted
        }
    }
}
```

### Step-by-Step Trace

```go
var ErrConnRefused = errors.New("connection refused")
err := fmt.Errorf("service: %w", fmt.Errorf("repo: %w", ErrConnRefused))
errors.Is(err, ErrConnRefused)   // → true
```

```
Step 1: err = wrapError{"service: repo: connection refused"}
        err == ErrConnRefused? NO → Unwrap()
Step 2: err = wrapError{"repo: connection refused"}
        err == ErrConnRefused? NO → Unwrap()
Step 3: err = errorString{"connection refused"}  ← IS ErrConnRefused
        err == ErrConnRefused? YES ✅
```

### Custom `Is()` Method

```go
type HTTPError struct{ Code int; Message string }
func (e *HTTPError) Error() string { return fmt.Sprintf("%d: %s", e.Code, e.Message) }
func (e *HTTPError) Is(target error) bool {
    t, ok := target.(*HTTPError)
    return ok && e.Code == t.Code   // match by status code only
}

err := &HTTPError{Code: 404, Message: "user not found"}
errors.Is(err, &HTTPError{Code: 404, Message: ""})  // true — codes match
```

---

## 7. `errors.As()` — Type Extraction Through the Chain

Walks the chain looking for an error assignable to the target type. Uses reflection.

### Step-by-Step Trace

```go
inner := &os.PathError{Op: "open", Path: "/etc/secret", Err: os.ErrPermission}
wrapped := fmt.Errorf("config: load: %w", inner)

var pe *os.PathError
errors.As(wrapped, &pe)   // → true
```

```
Step 1: err = wrapError{"config: load: open /etc/secret: permission denied"}
        Is *wrapError assignable to *os.PathError? NO → Unwrap()
Step 2: err = &os.PathError{Op:"open", Path:"/etc/secret"}
        Is *os.PathError assignable to *os.PathError? YES ✅
        pe.Op = "open", pe.Path = "/etc/secret", pe.Err = os.ErrPermission
```

### `errors.As()` Runtime Internals — What Actually Happens

`errors.As` uses `reflectlite` (a slimmed-down reflection package) to inspect types
at runtime. Here's the simplified source from `src/errors/wrap.go`:

```go
func As(err error, target any) bool {
    val := reflectlite.ValueOf(target)       // reflection on target pointer
    targetType := val.Type().Elem()          // the type target points to (*DNSError → DNSError)

    for {
        // KEY: type assignability check via reflection
        if reflectlite.TypeOf(err).AssignableTo(targetType) {
            val.Elem().Set(reflectlite.ValueOf(err))   // set *target = err
            return true
        }
        // Check if err implements custom As(target) method
        if x, ok := err.(interface{ As(any) bool }); ok {
            if x.As(target) { return true }
        }
        // Unwrap and continue down the chain
        err = Unwrap(err)
        if err == nil { return false }
    }
}
```

**What's happening at the iface/itab level:**

```
errors.As(err, &dnsErr)
  │
  ├─ reflectlite reads err's iface → itab → _type (concrete type descriptor)
  ├─ compares _type against *DNSError's type descriptor
  │   ├─ match? → set target via reflection, return true
  │   └─ no match? → type-assert err.(interface{ Unwrap() error })
  │                    └─ itab lookup: does concrete type have Unwrap()?
  │                        ├─ yes → call Unwrap(), loop with inner error
  │                        └─ no → return false
  └─ walks the entire wrap chain until match or nil
```

**Performance note:** `errors.As` uses reflection, making it ~10× slower than a direct
type assertion. In hot paths where you know the error isn't wrapped, prefer direct:

```go
// Fast — direct itab lookup, pointer comparison, no reflection (~5ns):
if dnsErr, ok := err.(*DNSError); ok { ... }

// Slower but walks wrapped chains — reflection at each level (~50-100ns per level):
var dnsErr *DNSError
if errors.As(err, &dnsErr) { ... }
```

Use `errors.As` when errors might be wrapped with `fmt.Errorf("%w")`.
Use direct type assertions when you know the error is unwrapped.

### `errors.Is()` vs `errors.As()`

```
errors.Is(err, target)   → "Does the chain contain this VALUE?"
                            Use with sentinels: errors.Is(err, io.EOF)
errors.As(err, &target)  → "Does the chain contain this TYPE?"
                            Use with custom types: errors.As(err, &ve)
```

---

## 7b. Type Assertions — The Runtime Machinery Behind Error Inspection

Type assertions are a **Go 1.0 language feature** — they've existed since Go's first release
(2012). They work **only on interface values** because that's the only case where the
concrete type is unknown at compile time.

### Why Only Interfaces?

With a concrete type (`var x int`), the compiler already knows the type — there's nothing
to "assert." With an interface, the concrete type is hidden behind the `iface`/`eface`
wrapper, so runtime machinery is needed to recover it.

### How It Works at Runtime

```go
var err error = &DNSError{Host: "example.com"}
dnsErr, ok := err.(*DNSError)   // type assertion
```

```
err is an iface { tab *itab, data unsafe.Pointer }
                    │
                    └─ itab { inter *interfacetype,  // error interface
                              _type *_type,           // concrete: *DNSError
                              fun   [1]uintptr }      // method: Error()

Type assertion does:
1. Read itab._type from the iface
2. Compare itab._type against the _type descriptor for *DNSError
3. This is a POINTER COMPARISON — each type has exactly one _type in the binary
4. Match? → return ((*DNSError)(iface.data), true)
5. No match? → return (nil, false)
```

**Cost:** Essentially a pointer comparison + pointer cast. Near zero (~1-2ns).
Much cheaper than Java's `instanceof` which may walk the class hierarchy.

### The Two Forms

```go
// Safe form — program continues:
val, ok := iface.(ConcreteType)   // ok=false if wrong type

// Panic form — crashes if wrong:
val := iface.(ConcreteType)       // panics if assertion fails — avoid in production
```

### Connection to errors.As

`errors.As` is essentially a loop of type-assertion-like checks using reflection,
walking the `Unwrap()` chain. A direct type assertion is the fast path; `errors.As`
is the general-purpose path for wrapped error chains.

```
Direct type assertion:  iface._type == target._type         → ~1-2ns
errors.As:              reflectlite.TypeOf(err) per level    → ~50-100ns × chain depth
```

---

## 7c. Why Errors Are Interfaces, Not Strings

If `error` were just a `string`, you could only read the message. The interface enables
**programmatic inspection**:

```go
// With strings — fragile string matching:
if err == "connection refused" { retry() }   // breaks if wording changes

// With error interface — carry structured data + identity:
type DNSError struct {
    Host    string
    Timeout bool
}
func (e *DNSError) Error() string { return "dns lookup failed: " + e.Host }

var dnsErr *DNSError
if errors.As(err, &dnsErr) {
    if dnsErr.Timeout { retry() }        // structured inspection
    else { alertOncall() }               // branch on data, not strings
}
```

**What the interface buys over strings:**

| Strings                       | `error` Interface                              |
|-------------------------------|------------------------------------------------|
| Compare by text (fragile)     | Compare by identity with `errors.Is()`         |
| No metadata                   | Carry fields: status codes, retry info, context|
| Can't type-switch             | `errors.As()` unwraps to concrete types        |
| Can't wrap                    | `fmt.Errorf("ctx: %w", err)` builds chains     |
| Dead data                     | Programmable: store in maps, send on channels  |

**Design insight:** `Error() string` is for **humans** (logging, display).
The **type itself** is for **code** (branching, inspection, wrapping).
That separation is the entire design.

### The Cross-Cutting Design: Why Error() Returns `string`, Not `[]byte`

This is where Go's design decisions connect across the entire language. The `error`
interface is an intersection of three independent design choices, each reinforcing
the others:

```
  ┌─────────────────────────────────────────────────────────────────────┐
  │              Three Pillars of Go's Error Design                     │
  │                                                                     │
  │  ┌──────────────┐   ┌──────────────────┐   ┌────────────────────┐  │
  │  │   INTERFACE   │   │ RETURNS string   │   │  VALUES, NOT       │  │
  │  │  (not string) │   │ (not []byte)     │   │  EXCEPTIONS        │  │
  │  ├──────────────┤   ├──────────────────┤   ├────────────────────┤  │
  │  │• Structured   │   │• Immutable       │   │• Returned, not     │  │
  │  │  context      │   │• Goroutine-safe  │   │  thrown             │  │
  │  │• Wrapping     │   │• No defensive    │   │• Explicit control  │  │
  │  │  chains       │   │  copies needed   │   │  flow              │  │
  │  │• nil = no err │   │• Safe as map key │   │• Cheap to create   │  │
  │  │• errors.Is/As │   │• Safe to log     │   │• No stack unwinding│  │
  │  │  type-safe    │   │  concurrently    │   │  overhead          │  │
  │  └──────────────┘   └──────────────────┘   └────────────────────┘  │
  └─────────────────────────────────────────────────────────────────────┘
```

**Why `string` and not `[]byte`?** Imagine if `Error()` returned `[]byte`:

```go
  // HYPOTHETICAL: Error() returns []byte
  err := fetchUser(42)

  // Goroutine 1: logging
  go func() {
      msg := err.Error()        // gets []byte — mutable!
      log.Println(string(msg))  // reading bytes...
  }()

  // Goroutine 2: also logging
  go func() {
      msg := err.Error()        // gets SAME []byte? or a copy?
      msg[0] = 'X'              // if shared → DATA RACE
  }()
```

With `[]byte`, every caller of `Error()` would need to either:
- Return a **new copy** every time (allocation per call — expensive)
- Share the backing array (data race risk)

With `string`, the return is **immutable by design**:
- Multiple goroutines can read `err.Error()` simultaneously — zero risk
- The string header is copied (16 bytes, free), backing bytes are shared safely
- No locks, no channels, no atomics needed

**This connects to how errors flow through concurrent systems:**

```go
  // Real production pattern — errors from multiple goroutines
  func processAll(items []Item) error {
      var mu sync.Mutex
      var errs []error

      var wg sync.WaitGroup
      for _, item := range items {
          wg.Add(1)
          go func() {
              defer wg.Done()
              if err := process(item); err != nil {
                  mu.Lock()
                  errs = append(errs, err)   // safe: error is interface value
                  mu.Unlock()                 // (we lock for the slice, not the error)
                  // Meanwhile, err.Error() can be called from ANY goroutine
                  // without coordination — the string it returns is immutable
              }
          }()
      }
      wg.Wait()
      return errors.Join(errs...)
  }
```

The mutex protects the **slice** (mutable), not the **errors** (immutable strings inside).
If `Error()` returned `[]byte`, you'd need to protect every error message too.

**The bigger picture — Go's immutability strategy:**

```
  ┌────────────────────────────────────────────────────────────┐
  │ Go uses immutability selectively, not universally:         │
  │                                                            │
  │  strings    → immutable  → safe to share across goroutines │
  │  []byte     → mutable    → must protect with mutex/channel │
  │  map keys   → must be comparable (strings work, slices don't) │
  │  error msgs → string     → safe to log/store/compare anywhere│
  │  channels   → thread-safe by design (internal mutex)       │
  │                                                            │
  │ Go doesn't make EVERYTHING immutable (unlike Rust/Haskell) │
  │ It makes the things that TRAVEL BETWEEN GOROUTINES safe:   │
  │ strings, error messages, channel values (copied on send)   │
  └────────────────────────────────────────────────────────────┘
```

This is why understanding slices, strings, and interfaces as **connected systems**
— not isolated topics — reveals Go's design philosophy. Each choice constrains and
enables the others. The language is small because the pieces fit together tightly.

---

## 7d. Java-to-Go Error Model Mental Bridge

For engineers coming from Java/C#, the error model is the biggest mental shift.

### Mapping Table

| Java/C#                           | Go                              | Key Difference                          |
|-----------------------------------|---------------------------------|-----------------------------------------|
| `throw new Exception("msg")`      | `return fmt.Errorf("msg")`      | Java unwinds stack; Go just returns     |
| `try { ... } catch (E e) { ... }` | `if err != nil { ... }`         | Go has no implicit control flow         |
| `throws IOException` in signature | `error` as return type          | Both visible; Go is simpler             |
| Unchecked `RuntimeException`      | `panic()` (rare, not for errors)| Go reserves panic for truly fatal bugs  |
| `finally { ... }`                 | `defer func() { ... }()`        | Go's defer runs on ALL exit paths       |
| `instanceof` (class hierarchy)    | Type assertion (pointer compare)| Go is faster: no hierarchy walk         |
| Stack trace captured at creation  | No stack trace by default       | Go errors are cheap; add context via %w |

### The Fundamental Shift

```
Java:   CREATING an exception and THROWING it are separate but almost always paired.
        The "throw" triggers stack unwinding — an expensive, invisible control flow jump.

Go:     There is NO "throw" step. You create the error and RETURN it.
        The caller checks it explicitly. No invisible jumps, no unwinding.
        The only thing that gives errors power is YOUR CODE checking them.
```

### What Happens When You Ignore an Error

```go
// Java — compiler forces you to handle checked exceptions:
file.read();   // won't compile without try/catch or throws declaration

// Go — compiles and runs fine, error silently dropped:
divide(10, 0)              // both return values discarded — no warning
fmt.Println("life goes on") // runs as if nothing happened

// The danger in production:
db.Exec("DELETE FROM users WHERE id = $1", userID)   // what if this failed?
fmt.Println("user deleted!")                           // lies.
```

Go's compiler **trusts you**. The language enforces that declared variables are used,
but if you never assign the returns, it's legal. Tooling catches this:
- `go vet` — warns about some ignored errors
- `errcheck` linter — catches ALL unhandled error returns
- `golangci-lint` — bundles errcheck and more

**In enterprise Go, running linters in CI is non-negotiable.**

---

## 8. `errors.Join()` (Go 1.20+)

Combines multiple errors into one. Returns `Unwrap() []error` (multi-error tree).

```go
// errors/join.go
type joinError struct{ errs []error }
func (e *joinError) Unwrap() []error { return e.errs }  // ← SLICE, not single
```

```
joinError
┌────────────────────────────────────────┐
│ errs []error                           │
│  ├─ [0] "name required"               │
│  ├─ [1] "email invalid"               │  errors.Is/As check ALL branches
│  └─ [2] "age positive"                │
└────────────────────────────────────────┘
```

### Use Cases

```go
// Multi-field validation
func validate(u User) error {
    var errs []error
    if u.Name == "" { errs = append(errs, fmt.Errorf("name: %w", ErrRequired)) }
    if !validEmail(u.Email) { errs = append(errs, fmt.Errorf("email: %w", ErrInvalid)) }
    return errors.Join(errs...)   // nil if errs is empty
}

// Parallel operations — collect all failures
func fetchAll(ctx context.Context, urls []string) error {
    var mu sync.Mutex
    var errs []error
    var wg sync.WaitGroup
    for _, url := range urls {
        wg.Add(1)
        go func(u string) {
            defer wg.Done()
            if err := fetch(ctx, u); err != nil {
                mu.Lock()
                errs = append(errs, fmt.Errorf("%s: %w", u, err))
                mu.Unlock()
            }
        }(url)
    }
    wg.Wait()
    return errors.Join(errs...)
}
```

---

## 9. Enterprise Error Handling Strategy

Errors flow **up** through layers, gaining context. Logged **once** at the top.

```
┌───────────────────────────────────────────────────────────────────┐
│  HANDLER — Maps domain errors → HTTP codes, logs ONCE            │
│                                                                   │
│  user, err := h.service.GetUser(ctx, id)                         │
│  if err != nil { h.handleError(w, err); return }                 │
└───────────────────────────────┬───────────────────────────────────┘
                                │
┌───────────────────────────────▼───────────────────────────────────┐
│  SERVICE — Wraps with business context, does NOT log             │
│                                                                   │
│  user, err := s.repo.FindByID(ctx, id)                           │
│  if err != nil { return nil, fmt.Errorf("get user %d: %w", id, err) }
└───────────────────────────────┬───────────────────────────────────┘
                                │
┌───────────────────────────────▼───────────────────────────────────┐
│  REPOSITORY — Translates DB errors → domain errors               │
│                                                                   │
│  if errors.Is(err, sql.ErrNoRows) {                              │
│      return nil, fmt.Errorf("user %d: %w", id, ErrNotFound)      │
│  }                                                                │
└───────────────────────────────────────────────────────────────────┘
```

### Error Mapping at the Handler

```go
func (h *Handler) handleError(w http.ResponseWriter, err error) {
    h.logger.Error("request failed", "error", err)   // log FULL chain ONCE

    switch {
    case errors.Is(err, ErrNotFound):
        http.Error(w, "resource not found", http.StatusNotFound)
    case errors.Is(err, ErrPermission):
        http.Error(w, "forbidden", http.StatusForbidden)
    default:
        http.Error(w, "internal server error", http.StatusInternalServerError)
    }
}
```

### The Error Chain Tells the Full Story

```
err.Error(): "get user 42: user 42: not found"

errors.Is(err, ErrNotFound) walks:
  wrapError{"get user 42: user 42: not found"}
    → Unwrap() → wrapError{"user 42: not found"}
      → Unwrap() → errorString{"not found"} ← IS ErrNotFound ✅

The chain answers: WHAT (not found) + WHERE (user 42, repo, service)
```

### The Golden Rules

```
1. WRAP at every layer     — fmt.Errorf("context: %w", err)
2. TRANSLATE at boundaries — sql.ErrNoRows → ErrNotFound
3. LOG once at the top     — handler/middleware only
4. MAP to status codes     — handler only, using errors.Is/As
5. SANITIZE for clients    — never expose internal error messages
```

---

## 10. Panic and Recover

`panic()` is for **truly unrecoverable** situations — programmer errors, not runtime conditions.

```go
// YES — valid panic (bug in code, detected at init)
func MustCompile(pattern string) *regexp.Regexp {
    re, err := regexp.Compile(pattern)
    if err != nil { panic(fmt.Sprintf("invalid regex %q: %v", pattern, err)) }
    return re
}

// NO — return an error (expected runtime condition)
func ReadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil { return nil, fmt.Errorf("read config %s: %w", path, err) }
    // ...
}
```

### Stack Unwinding

```
  CALL STACK (growing down)        PANIC UNWIND (going back up)
  ─────────────────────────        ────────────────────────────
  main()                           ⑥ main: deferred (defer runs)
    │                              ⑤ panic continues up
    ├─ levelOne()                  ④ one: deferred  (defer runs)
    │    │                         ③ panic continues up
    │    ├─ levelTwo()             ② two: deferred  (defer runs)
    │    │    │                    ① panic("broke") ← starts here
    │    │    └─ PANIC ════════════╝
```

Defers run in LIFO order during unwind. If no `recover()` is found, program crashes.

### `recover()` — Catching a Panic

Only works **inside a deferred function**. This is Go's equivalent of try/catch.

```go
func safeOperation() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()
    riskyCode()
    return nil
}
```

### Why `recover()` MUST Be Inside `defer`

When `panic()` fires, execution **jumps** — lines after the panic never run.
The **only** code that executes during a panic unwind is **deferred functions**.

```go
func broken() {
    recover()        // ❌ useless — normal execution, no panic is happening
    riskyCode()      // panic here — stack starts unwinding
    recover()        // ❌ never reached — we already panicked past this line
}

func works() {
    defer func() {
        recover()    // ✅ runs during unwind — catches the panic
    }()
    riskyCode()      // panic here — defer runs, recover() catches it
}
```

### Why the Anonymous Function?

`defer` needs a **function call**. You can't write meaningful recovery without one:

```go
defer recover()           // ❌ recovers but throws away the value — can't log, can't set err

defer fmt.Println("bye")  // ✅ valid but only one statement — can't do recovery logic

defer func() {            // ✅ anonymous function lets you:
    r := recover()        //    1. capture the panic value
    err = fmt.Errorf(...) //    2. assign to the named return
    log.Error(...)        //    3. log it
    metrics.Inc(...)      //    4. record metrics
}()                       //    ← note the () — you're CALLING the function
```

### How Named Returns Make This Work

The critical question: **how does the defer set the return value without a `return` statement?**

Named return values are the answer. With `(err error)` in the signature, `err` is not
a local variable — it IS the **return slot** on the stack frame:

```
Stack frame of safeOperation:
┌─────────────────────────┐
│  return slot: err error  │  ← allocated at function entry, initialized to nil
├─────────────────────────┤
│  local variables...      │
└─────────────────────────┘
```

The deferred closure **captures** `err` — the same memory location as the return slot.
Writing to `err` inside the defer directly modifies what the caller will receive:

```go
func safeOperation() (err error) {       // return slot: err = nil
    defer func() {                        // closure captures &err (the return slot)
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)  // writes directly to the return slot
            // NO return needed — we're modifying the slot, not a local copy
        }
    }()
    riskyCode()    // PANIC! → defer runs → recover → err = error → function exits
    return nil     // never reached on panic path
}
// function exits → caller reads the return slot → finds the error we wrote
```

**Step-by-step trace — panic path:**

```
① func safeOperation() (err error)    → return slot allocated, err = nil
② defer func() { ... }()              → anonymous function registered (not run yet)
③ riskyCode()                         → PANIC! stack starts unwinding
④ defer fires                         → anonymous function runs
⑤ r := recover()                      → catches panic, r = the panic value
⑥ err = fmt.Errorf("panic: %v", r)   → writes error to return slot
⑦ function exits                      → caller reads return slot → gets the error
```

**Step-by-step trace — happy path:**

```
① func safeOperation() (err error)    → return slot allocated, err = nil
② defer func() { ... }()              → registered
③ riskyCode()                         → completes normally
④ return nil                          → writes nil to return slot
⑤ defer fires                         → recover() returns nil → if block skipped
⑥ function exits                      → caller reads return slot → gets nil
```

**This would NOT work without named returns:**

```go
func broken() error {           // unnamed return — no variable to reference
    defer func() {
        if r := recover(); r != nil {
            // HOW do I set the return value?
            // There's no variable name to write to.
            // I can't "return" from inside a defer — return exits the DEFER, not the outer func.
        }
    }()
    riskyCode()
    return nil
}
```

**The Java parallel:**

```
Java:   try { riskyCode(); } catch (Exception e) { return new Error("recovered: " + e); }
Go:     defer func() { if r := recover(); r != nil { err = fmt.Errorf("recovered: %v", r) } }()
```

Both achieve the same result. Go is more explicit — no hidden control flow, the recovery
mechanism is visible as a deferred function call.

**Production rule:** You almost never call `recover()` in business logic. It belongs in:
1. **HTTP middleware** — recover panics so one bad request doesn't crash the server
2. **Goroutine launchers** — a panic in a goroutine kills the **entire program**, not just that goroutine

### HTTP Recovery Middleware

```go
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                slog.Error("panic recovered", "panic", rec,
                    "stack", string(debug.Stack()), "path", r.URL.Path)
                http.Error(w, "internal server error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

---

## 11. Error Handling Anti-Patterns

### (a) Swallowing Errors

```go
result, _ := doSomething()            // BAD — silent failure
result, err := doSomething()           // FIX — handle it
if err != nil { return fmt.Errorf("do something: %w", err) }
```

### (b) Logging AND Returning

```go
// BAD — logged at every layer = duplicate logs
log.Printf("failed: %v", err)
return nil, fmt.Errorf("get user: %w", err)

// FIX — wrap and return only; handler logs once
return nil, fmt.Errorf("get user %d: %w", id, err)
```

### (c) Wrapping Without Context

```go
return fmt.Errorf("%w", err)                   // BAD — adds nothing
return fmt.Errorf("repo: find user %d: %w", id, err)  // FIX — WHERE + WHAT
```

### (d) Panic for Control Flow

```go
panic("user not found")              // BAD — expected condition
return nil, ErrNotFound              // FIX — return error value
```

### (e) Comparing Error Strings

```go
if err.Error() == "connection refused" { }    // BAD — fragile
if errors.Is(err, ErrConnRefused) { }         // FIX — use sentinels
```

### (f) Ignoring Close/Write Errors

```go
defer f.Close()                      // BAD — buffered writes may be lost
defer func() {                       // FIX — capture close error
    if cerr := f.Close(); cerr != nil && err == nil { err = cerr }
}()
```

### Summary

```
┌───────────────────────────┬──────────────────────────────────────────┐
│ Anti-Pattern              │ Fix                                      │
├───────────────────────────┼──────────────────────────────────────────┤
│ _ = f()                   │ Always handle or explicitly document why │
│ log + return              │ Wrap and return; log once at top         │
│ fmt.Errorf("%w", err)     │ Add context: "layer: operation: %w"     │
│ panic for expected errors │ Return error values                      │
│ err.Error() == "..."      │ errors.Is() or errors.As()              │
│ defer f.Close()           │ Capture and check close error            │
└───────────────────────────┴──────────────────────────────────────────┘
```

---

## 12. Error Handling in Concurrent Code

### Panics Don't Cross Goroutine Boundaries

```
  main goroutine              child goroutine
  ──────────────              ────────────────
  main()                      func() {
    │  defer func() {           panic("boom")
    │    recover() ← NEVER      │
    │    TRIGGERED!              ↓
    │  }()                     no recover → PROGRAM CRASH
    └─ go func()...
```

**Always** `recover()` inside every goroutine that might panic.

### `errgroup.Group` — The Standard Pattern

```go
g, ctx := errgroup.WithContext(ctx)
for _, url := range urls {
    g.Go(func() error { return fetch(ctx, url) })
}
return g.Wait()   // first error cancels ctx, waits for all to finish
```

### Channel-Based Error Collection (All Errors)

```go
errCh := make(chan error, len(urls))
var wg sync.WaitGroup
for _, url := range urls {
    wg.Add(1)
    go func(u string) {
        defer wg.Done()
        if err := fetch(ctx, u); err != nil {
            errCh <- fmt.Errorf("%s: %w", u, err)
        }
    }(url)
}
go func() { wg.Wait(); close(errCh) }()

var errs []error
for err := range errCh { errs = append(errs, err) }
return errors.Join(errs...)
```

---

## 13. Performance Characteristics

### Cost Table

```
┌──────────────────────────────┬─────────────┬─────────────────────────────────────┐
│ Operation                    │ Approx Cost │ Why                                 │
├──────────────────────────────┼─────────────┼─────────────────────────────────────┤
│ errors.New("msg")            │ ~50-80ns    │ 1 heap alloc (errorString struct)   │
│                              │ 1 alloc     │ + iface construction                │
├──────────────────────────────┼─────────────┼─────────────────────────────────────┤
│ fmt.Errorf("ctx: %w", err)   │ ~200-400ns  │ String formatting + 1 heap alloc   │
│                              │ 1-2 allocs  │ (wrapError struct with msg + err)   │
├──────────────────────────────┼─────────────┼─────────────────────────────────────┤
│ errors.Is() — depth 1        │ ~5-10ns     │ Pointer comparison + type assertion │
│                              │ 0 allocs    │                                     │
├──────────────────────────────┼─────────────┼─────────────────────────────────────┤
│ errors.Is() — depth N        │ ~5-10ns × N │ Linear walk: Unwrap + compare each │
│                              │ 0 allocs    │                                     │
├──────────────────────────────┼─────────────┼─────────────────────────────────────┤
│ errors.As() — depth 1        │ ~50-100ns   │ Uses reflect.TypeOf internally     │
│                              │ 0-1 allocs  │                                     │
├──────────────────────────────┼─────────────┼─────────────────────────────────────┤
│ errors.As() — depth N        │ ~50-100ns×N │ Reflection at each level + Unwrap  │
│                              │ 0-1 allocs  │                                     │
├──────────────────────────────┼─────────────┼─────────────────────────────────────┤
│ panic + recover              │ ~1-5µs      │ Stack unwinding + deferred func     │
│                              │ varies      │ execution. 100-1000× costlier      │
├──────────────────────────────┼─────────────┼─────────────────────────────────────┤
│ Stack trace (debug.Stack())  │ ~5-20µs     │ Walks call stack, formats PCs to   │
│                              │ many allocs │ function names                      │
└──────────────────────────────┴─────────────┴─────────────────────────────────────┘
```

### Go Errors vs Exceptions

```
Java/C# Exception:                    Go Error:
─────────────────                     ─────────
• Created: captures full stack trace  • Created: allocates a small struct
  (~5-20µs, many allocations)           (~50-80ns, 1 allocation)
• Thrown: unwinds the call stack      • Returned: normal function return
  (walks frames, runs finally)          (register/stack copy, ~1ns)
• Caught: exception table lookup      • Checked: if err != nil (~0.5ns)

Go errors are 25-100× cheaper than Java exceptions.
```

### Happy Path Is Free

The CPU branch predictor learns `err != nil` is almost always `false`. The nil
pointer comparison is essentially free. Error creation cost only matters on failure.

---

## 14. Quick Reference Card

### Decision Tree

```
You need to signal an error condition
│
├─ Known condition callers must check?
│   ├─ Simple (not found, timeout) → Sentinel: var ErrX = errors.New("...")
│   └─ Needs data (field, code)    → Custom type: type XError struct{...}
│
├─ Adding context to existing error? → fmt.Errorf("context: %w", err)
├─ Combining multiple errors?        → errors.Join(err1, err2, ...)
└─ Simple internal error?            → errors.New("description")
```

### Checklist

```
✅ ALWAYS check returned errors        ❌ DON'T log AND return
✅ WRAP with context at each layer     ❌ DON'T compare error strings
✅ USE errors.Is/As not == or strings  ❌ DON'T panic for expected errors
✅ TRANSLATE at boundaries             ❌ DON'T wrap without adding context
✅ LOG once at the top                 ❌ DON'T export types unless needed
✅ RECOVER panics in goroutines
✅ CHECK Close() errors for writers
```

### Common Patterns

```go
// Sentinel
var ErrNotFound = errors.New("not found")

// Custom type
type ValidationError struct{ Field, Message string }
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s — %s", e.Field, e.Message)
}

// Wrap, check, extract
return fmt.Errorf("service: get user %d: %w", id, err)
if errors.Is(err, ErrNotFound) { /* sentinel match */ }
var ve *ValidationError
if errors.As(err, &ve) { /* use ve.Field */ }

// Goroutine recovery
defer func() { if r := recover(); r != nil { log.Printf("recovered: %v", r) } }()

// Safe close
defer func() { if cerr := f.Close(); cerr != nil && err == nil { err = cerr } }()
```

### Tools

```bash
go vet ./...                    # common error handling mistakes
go build -gcflags='-m' ./...    # escape analysis: error alloc locations
errcheck ./...                  # finds unchecked error returns
go test -race ./...             # race detector: always run
```

---

## One-Line Summary

> An `error` is a one-method interface (`Error() string`) stored as `iface{tab, data}` —
> wrapping with `%w` builds a linked chain walked by `errors.Is`/`errors.As`, sentinels
> define API contracts checked by value, custom types carry structured data extracted by
> type, and production code wraps at every layer, translates at boundaries, and logs once
> at the top.
