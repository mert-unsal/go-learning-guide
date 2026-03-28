# Deep Dive: Production Go — Pitfalls, Library Choices & Enterprise Patterns

> Everything a senior engineer needs to know before shipping Go to production.
> The bugs that bite at 3 AM, the libraries worth adopting, and the patterns
> that survive at scale.

---

## Table of Contents

1. [The Top 15 Go Production Pitfalls](#1-the-top-15-go-production-pitfalls)
2. [Library & Framework Comparison Tables](#2-library--framework-comparison-tables)
3. [The Go Community's Strong Opinions](#3-the-go-communitys-strong-opinions)
4. [Project Structure for Enterprise](#4-project-structure-for-enterprise)
5. [Graceful Shutdown — The Production Pattern](#5-graceful-shutdown--the-production-pattern)
6. [Docker Multi-Stage Build — The Standard](#6-docker-multi-stage-build--the-standard)
7. [Middleware Pattern — The Go Way](#7-middleware-pattern--the-go-way)
8. [Configuration Management](#8-configuration-management)
9. [Quick Reference Card](#9-quick-reference-card)

---

## 1. The Top 15 Go Production Pitfalls

### Pitfall #1 — Goroutine Leak (blocked on channel, no cancellation)

**The Bug:**
```go
func fetch(url string) <-chan string {
    ch := make(chan string)
    go func() {
        resp, _ := http.Get(url)           // blocks forever if server hangs
        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)
        ch <- string(body)                 // blocks forever if nobody reads
    }()
    return ch
}
```

**Why it happens:** The goroutine has no way to be told "stop waiting." If the caller
abandons the channel or the HTTP request hangs, the goroutine lives forever. Each leaked
goroutine holds its stack (~2-8 KB), plus any heap objects it references — including the
`http.Response`, its TLS buffers, and the TCP connection. At scale, thousands of leaked
goroutines cause OOM.

**How to detect:**
```
┌───────────────────────────────────────────────────────────┐
│  runtime.NumGoroutine()    — monitor in /debug/vars       │
│  pprof goroutine profile   — GET /debug/pprof/goroutine   │
│  go tool pprof             — list goroutines by function  │
│  goleak (uber-go/goleak)   — fails tests on goroutine leak│
└───────────────────────────────────────────────────────────┘
```

**How to fix:** Always pass `context.Context` and `select` on `ctx.Done()`:
```go
func fetch(ctx context.Context, url string) (string, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    return string(body), err
}
```

---

### Pitfall #2 — Nil Interface Trap (typed nil vs true nil)

**The Bug:**
```go
func getError() error {
    var err *MyError = nil       // typed nil pointer
    return err                   // wraps in iface{*itab, nil} — NOT nil!
}

if getError() != nil {           // TRUE — even though the value is nil
    fmt.Println("unexpected!")   // this runs
}
```

**Why it happens:** An interface is a two-word struct `{type, data}`. When you assign
a typed nil pointer, the type word is non-nil (it points to `*MyError`'s type info).
`== nil` only returns true when BOTH words are zero.

```
  var err *MyError = nil       var err error = nil
  ┌──────────┬──────────┐     ┌──────────┬──────────┐
  │ *MyError │   nil    │     │   nil    │   nil    │
  └──────────┴──────────┘     └──────────┴──────────┘
  err == nil → FALSE ✗         err == nil → TRUE ✓
  (type word is non-nil)       (both words are nil)
```

**How to fix:** Never return a concrete typed nil from a function that returns an
interface. Return the bare `nil`:
```go
func getError() error {
    var err *MyError = validate()
    if err != nil {
        return err          // only return if non-nil
    }
    return nil              // return untyped nil
}
```

---

### Pitfall #3 — Slice Backing Array Memory Leak

**The Bug:**
```go
var bigSlice = make([]byte, 1<<20)  // 1 MB

func first10() []byte {
    return bigSlice[:10]            // returns 10 bytes, holds 1 MB in memory
}
```

**Why it happens:** A slice header is `{pointer, len, cap}`. Sub-slicing shares the
backing array. The GC cannot free the 1 MB array because the 10-byte slice still
references it.

```
  bigSlice[:10]
  ┌──────────┬──────┬──────────┐
  │ ptr ─────┼──►┌──┴──────────┴───────────────────────────┐
  │ len = 10 │   │ 10 bytes used │ ~~~ 1MB-10 wasted ~~~  │
  │ cap = 1M │   └─────────────────────────────────────────┘
  └──────────┘   backing array pinned by ptr — GC can't free it
```

**How to fix:** Copy the data to detach from the backing array:
```go
func first10() []byte {
    result := make([]byte, 10)
    copy(result, bigSlice[:10])
    return result                   // new 10-byte backing array
}
```

---

### Pitfall #4 — Map Never Shrinks

**The Bug:**
```go
m := make(map[int]string)
for i := 0; i < 1_000_000; i++ {
    m[i] = "data"
}
for i := 0; i < 1_000_000; i++ {
    delete(m, i)
}
// m now has 0 entries, but the internal bucket array is still sized for 1M entries
```

**Why it happens:** `runtime.hmap` allocates buckets in powers of 2. Deleting keys
marks slots as empty (`tophash = emptyOne`) but NEVER deallocates bucket memory.
The only way to reclaim is to create a new map and copy surviving entries.

**How to detect:** `runtime.ReadMemStats()` — watch `HeapInuse` after bulk deletes.

**How to fix:** Periodically rebuild the map, or use a map of pointers where the
pointed-to values can be GC'd independently:
```go
// Rebuild to reclaim memory
newMap := make(map[int]string, len(m))
for k, v := range m {
    newMap[k] = v
}
m = newMap
```

---

### Pitfall #5 — `time.After` in Loops (timer leak)

**The Bug:**
```go
for {
    select {
    case msg := <-ch:
        handle(msg)
    case <-time.After(5 * time.Second):  // NEW timer every iteration — never GC'd until fired
        log.Println("timeout")
    }
}
```

**Why it happens:** `time.After` creates a `*time.Timer` that cannot be garbage collected
until it fires. In a hot loop processing thousands of messages/sec, this creates thousands
of live timers, each holding a goroutine in the runtime timer heap.

**How to fix:** Reuse a single timer with `time.NewTimer` + `Reset`:
```go
timer := time.NewTimer(5 * time.Second)
defer timer.Stop()
for {
    timer.Reset(5 * time.Second)
    select {
    case msg := <-ch:
        handle(msg)
    case <-timer.C:
        log.Println("timeout")
    }
}
```

---

### Pitfall #6 — Forgetting to Check Error from `defer Close/Flush`

**The Bug:**
```go
func writeFile(path string, data []byte) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()                 // error from Close is silently swallowed!

    w := bufio.NewWriter(f)
    w.Write(data)
    return w.Flush()                // Flush succeeds, but Close may fail (disk full, NFS)
}
```

**Why it happens:** `defer` executes after the return statement. If `Close()` returns
an error (disk full, network filesystem timeout), it's silently discarded.
Buffered writers may not flush their last buffer until `Close`.

**How to fix:** Use a named return to capture the deferred error:
```go
func writeFile(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer func() {
        closeErr := f.Close()
        if err == nil {
            err = closeErr
        }
    }()

    w := bufio.NewWriter(f)
    if _, err = w.Write(data); err != nil {
        return err
    }
    return w.Flush()
}
```

---

### Pitfall #7 — Data Race on Map (fatal crash, not panic)

**The Bug:**
```go
m := make(map[string]int)
go func() { m["a"] = 1 }()
go func() { m["b"] = 2 }()
// fatal error: concurrent map writes — PROCESS KILLED, cannot recover()
```

**Why it happens:** Since Go 1.6, the runtime inserts a write flag check in map
operations. If two goroutines write simultaneously, the runtime calls `throw()`, not
`panic()`. This is an **unrecoverable fatal error** — `recover()` cannot catch it.

```
  Goroutine 1: map write → sets hashWriting flag
  Goroutine 2: map write → sees hashWriting flag → throw("concurrent map writes")
  ┌──────────────────────────────────────────┐
  │  FATAL: not a panic — cannot recover()   │
  │  Process exit code 2, stack trace dumped  │
  └──────────────────────────────────────────┘
```

**How to detect:** `go test -race ./...` — the race detector catches this before production.

**How to fix:** Use `sync.RWMutex` for general maps, or `sync.Map` for the two specific
patterns it's optimized for (write-once-read-many, disjoint key sets per goroutine):
```go
var mu sync.RWMutex
var m = make(map[string]int)

// Write
mu.Lock()
m["a"] = 1
mu.Unlock()

// Read
mu.RLock()
v := m["a"]
mu.RUnlock()
```

---

### Pitfall #8 — Closure Capturing Loop Variable (pre-Go 1.22)

**The Bug (Go < 1.22):**
```go
for _, v := range items {
    go func() {
        process(v)      // all goroutines see the LAST value of v
    }()
}
```

**Why it happens:** Before Go 1.22, the loop variable `v` was a single variable
reused across all iterations. Closures capture the variable by reference, not by
value. By the time goroutines execute, `v` holds the final iteration's value.

**Status:** Fixed in Go 1.22 — each iteration now gets a new variable. But legacy
code and older Go versions still have this bug. The fix for older code:
```go
for _, v := range items {
    v := v              // shadow with per-iteration copy
    go func() {
        process(v)      // correct — each goroutine has its own copy
    }()
}
```

---

### Pitfall #9 — String Concatenation in Loops (O(n²) allocations)

**The Bug:**
```go
var s string
for _, item := range items {
    s += item.Name + ","      // each += allocates a NEW string
}
```

**Why it happens:** Strings in Go are immutable `{pointer, length}` headers. Each
`+=` allocates a new backing array, copies the old content + new content. For n
items, this is O(n²) total bytes copied.

**How to fix:**
```go
var b strings.Builder
for _, item := range items {
    b.WriteString(item.Name)
    b.WriteByte(',')
}
s := b.String()             // single allocation at the end
```

`strings.Builder` uses an internal `[]byte` that grows with amortized O(1) appends.
Pre-allocate with `b.Grow(estimatedSize)` for hot paths.

---

### Pitfall #10 — JSON Unmarshaling to `interface{}` (numbers become float64)

**The Bug:**
```go
var data map[string]interface{}
json.Unmarshal([]byte(`{"id": 12345678901234}`), &data)
id := data["id"].(int)     // PANIC: interface conversion: float64, not int
// Even if it doesn't panic, float64 loses precision for large int64 values
```

**Why it happens:** The JSON spec has a single `number` type. Go's `encoding/json`
decodes all numbers into `float64` when the target is `interface{}`, because `float64`
can represent both integers and decimals. But `float64` only has 53 bits of mantissa —
integers > 2^53 lose precision silently.

**How to fix:** Use `json.Decoder` with `UseNumber()`, or unmarshal into a concrete struct:
```go
dec := json.NewDecoder(bytes.NewReader(raw))
dec.UseNumber()                                  // numbers become json.Number (string)
dec.Decode(&data)
id, _ := data["id"].(json.Number).Int64()        // safe conversion
```

---

### Pitfall #11 — `context.Value` Abuse (dependency injection via context)

**The Bug:**
```go
ctx = context.WithValue(ctx, "db", dbConn)
ctx = context.WithValue(ctx, "logger", logger)
ctx = context.WithValue(ctx, "userID", uid)

// deep in the call stack:
db := ctx.Value("db").(*sql.DB)     // type-unsafe, invisible dependency, O(n) lookup
```

**Why it happens:** `context.WithValue` creates a linked list. Each `Value()` call walks
the chain — O(n) in context depth. Worse, it hides dependencies: the function signature
doesn't reveal what the function needs, making testing and refactoring dangerous.

```
  ctx.Value("db") walks:
  ctx4{key:"userID"} → ctx3{key:"logger"} → ctx2{key:"db"} → FOUND
       ▲                    ▲                    ▲
       O(1)                 O(2)                 O(3) — linear walk
```

**Legitimate uses:** request-scoped metadata that crosses API boundaries (trace ID,
request ID, auth claims). NOT for services, connections, or config.

**How to fix:** Pass dependencies explicitly via constructor injection:
```go
func NewHandler(db *sql.DB, logger *slog.Logger) *Handler {
    return &Handler{db: db, logger: logger}
}
```

---

### Pitfall #12 — Unbounded Goroutine Spawning (no backpressure, OOM)

**The Bug:**
```go
for _, job := range jobs {          // 1 million jobs
    go process(job)                 // 1 million goroutines — OOM
}
```

**Why it happens:** Each goroutine costs ~2-8 KB stack + heap objects it allocates.
1 million goroutines = 2-8 GB just in stacks, plus all their allocations. No
backpressure means the system cannot push back when overwhelmed.

**How to fix:** Use a semaphore pattern or worker pool:
```go
sem := make(chan struct{}, maxWorkers)
var wg sync.WaitGroup
for _, job := range jobs {
    sem <- struct{}{}               // blocks when maxWorkers goroutines are running
    wg.Add(1)
    go func(j Job) {
        defer func() { <-sem; wg.Done() }()
        process(j)
    }(job)
}
wg.Wait()
```

---

### Pitfall #13 — Silent Panic in Goroutine (crashes entire process)

**The Bug:**
```go
go func() {
    result := riskyOperation()      // panics!
    ch <- result
}()
// panic in a goroutine is NOT caught by the caller's recover()
// the entire process crashes with a stack trace
```

**Why it happens:** `recover()` only works within the same goroutine's deferred
function. A panic in a child goroutine propagates up to `runtime.main` and kills
the entire process — even if the parent has `recover()`.

**How to fix:** Every goroutine you launch should have panic recovery:
```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("recovered panic: %v\n%s", r, debug.Stack())
        }
    }()
    result := riskyOperation()
    ch <- result
}()
```

For HTTP servers, `net/http` already recovers panics per-handler. For worker
goroutines, you MUST add recovery yourself.

---

### Pitfall #14 — Copying `sync` Types (Mutex, WaitGroup, Cond)

**The Bug:**
```go
type Service struct {
    mu sync.Mutex
    cache map[string]string
}

func (s Service) Get(key string) string {   // VALUE receiver — copies s, copies mu!
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.cache[key]                     // locking a COPY of the mutex — no protection
}
```

**Why it happens:** `sync.Mutex` contains internal state (wait queue, lock flag).
Copying it creates a second mutex with duplicated state — the copy and original are
independent. The `noCopy` sentinel embedded in sync types triggers `go vet` warnings.

**How to detect:** `go vet ./...` catches this: `copies lock value: sync.Mutex`.

**How to fix:** Use pointer receivers for any struct containing sync types:
```go
func (s *Service) Get(key string) string {  // POINTER receiver — no copy
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.cache[key]
}
```

---

### Pitfall #15 — HTTP Handler Not Draining Request Body

**The Bug:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    if r.ContentLength > maxSize {
        http.Error(w, "too large", 413)
        return                              // body not read — connection can't be reused!
    }
    // ...
}
```

**Why it happens:** HTTP/1.1 connection reuse (keep-alive) requires the full request
body to be consumed. If the handler returns without reading the body, the server must
close the connection — no reuse. Under load, this creates a flood of new TCP
connections (with TLS handshake cost).

**How to fix:** Drain and close the body:
```go
func handler(w http.ResponseWriter, r *http.Request) {
    defer func() {
        io.Copy(io.Discard, r.Body)
        r.Body.Close()
    }()
    if r.ContentLength > maxSize {
        http.Error(w, "too large", 413)
        return
    }
    // ...
}
```

---

## 2. Library & Framework Comparison Tables

### HTTP Routers

| | **net/http (stdlib)** | **chi** | **gin** | **echo** |
|---|---|---|---|---|
| **Description** | stdlib router; Go 1.22 added `GET /users/{id}` pattern matching | Lightweight, idiomatic; 100% compatible with `net/http` | Full framework; fastest benchmarks; Radix tree router | Full framework; similar to Express.js in feel |
| **Pros** | Zero deps; always maintained; core team support | Stdlib compatible; middleware-first; small API surface | Fast; huge ecosystem; battle-tested at scale | Good docs; WebSocket/HTTP2 built-in; nice API |
| **Cons** | No middleware chaining stdlib; new patterns are basic | Smaller community than gin; fewer tutorials | Uses own `*gin.Context` — not stdlib `http.Handler`; heavy | Uses own `echo.Context`; smaller community |
| **When to use** | Simple services; when zero deps matters | Best default choice for most services | High-perf APIs; team already knows gin | Teams from Node/Express background |
| **Who uses it** | Everyone (stdlib) | Cloudflare, Heroku | Uber, many startups | Various |
| **Maintenance** | Go team (forever) | Active | Active | Active |
| **Verdict** | ✅ Always safe | ✅ **Recommended default** | ⚠️ Good but not idiomatic | ⚠️ Good but not idiomatic |

### Logging

| | **log/slog (stdlib)** | **zap (Uber)** | **zerolog** | **logrus** |
|---|---|---|---|---|
| **Description** | Structured logging in stdlib since Go 1.21; `slog.Logger` with handlers | Uber's zero-allocation structured logger | Zero-allocation JSON logger | Oldest structured logger; widely used |
| **Pros** | Stdlib; standard API; handler interface extensible | Fastest; zero-alloc hot path; battle-tested at Uber scale | Very fast; small API; zero-alloc | Familiar API; huge plugin ecosystem |
| **Cons** | Newer; fewer handlers available | Large API surface; Uber-specific patterns | Fluent API feels non-Go; less flexible | **Maintenance mode**; slower than zap/zerolog |
| **When to use** | New projects; when stdlib preference is strong | Latency-critical services needing maximum perf | When binary size and speed matter most | Legacy projects already using it |
| **Who uses it** | Go community (growing) | Uber, many large companies | Various | Docker, many older projects |
| **Verdict** | ✅ **Recommended default** | ✅ Production-proven | ✅ Good alternative | ❌ Don't adopt new |

### Database

| | **database/sql** | **sqlx** | **pgx** | **GORM** | **sqlc** |
|---|---|---|---|---|---|
| **Description** | Stdlib SQL interface; driver-agnostic | Extension of `database/sql`; struct scanning | PostgreSQL-native driver; no `database/sql` wrapper | Full ORM; auto-migrations; associations | Code gen from SQL queries → type-safe Go |
| **Pros** | Stable; driver-agnostic; well understood | Drop-in upgrade; named params; struct scan | Fastest PG driver; LISTEN/NOTIFY; COPY; batch | Rapid prototyping; migrations; associations | Compile-time SQL validation; no runtime reflection |
| **Cons** | Verbose; manual scanning; no named params | Still limited to `database/sql` perf | PostgreSQL only; non-standard API | Slow; magic; hard to debug; fights Go philosophy | Only SQL → Go (not Go → SQL); learning curve |
| **When to use** | Multi-DB support; simple queries | Upgrade path from `database/sql` | PostgreSQL services needing max performance | Prototypes; CRUD-heavy apps (controversial) | Services with complex queries; type safety |
| **Who uses it** | Everyone | Many companies | Many PG-heavy services | Startups, rapid dev | Stripe, many Go teams |
| **Verdict** | ✅ Safe baseline | ✅ Easy upgrade | ✅ **Best for PostgreSQL** | ⚠️ Controversial | ✅ **Recommended for SQL-heavy** |

### Configuration

| | **os.Getenv + flags** | **viper** | **envconfig** | **koanf** |
|---|---|---|---|---|
| **Description** | Stdlib env vars + `flag` package | Full config solution: files, env, remote, flags | Struct tags → env vars mapping | Modular config library; provider-based |
| **Pros** | Zero deps; simple; explicit | Everything included; popular; multi-format | Dead simple; zero magic; one job done well | Clean API; composable providers; no globals |
| **Cons** | No file parsing; no struct mapping; manual | Heavy; global state; too much magic; slow | Only env vars; no file support | Less popular; newer |
| **When to use** | Simple services; 12-factor apps | Legacy; complex config sources | 12-factor microservices | Modern replacement for viper |
| **Verdict** | ✅ Start here | ⚠️ Overused | ✅ Good for microservices | ✅ Modern choice |

### Testing

| | **testing (stdlib)** | **testify** | **go-cmp (Google)** | **gomock / mockgen** |
|---|---|---|---|---|
| **Description** | Built-in test framework; `go test` | Assertions, require, suite, mock | Deep equality comparison with diff output | Interface mock generation |
| **Pros** | Always available; table-driven tests; subtests; fuzzing | Familiar assertions; reduces boilerplate; wide adoption | Semantic comparison; diff output; option composability | Auto-generates mocks; strict expectations |
| **Cons** | Verbose assertions (`if got != want`); no diff | Extra dep; `suite` promotes non-Go patterns; controversial | Only comparison, not full assertion library | Generated code is verbose; tight coupling to impl |
| **When to use** | Always — it's the foundation | When team prefers assertion style | Complex struct comparison; golden file tests | When you need strict call verification |
| **Verdict** | ✅ **Always use** | ⚠️ Divisive in community | ✅ **Recommended for comparisons** | ⚠️ Use sparingly |

### CLI Frameworks

| | **cobra + pflag** | **urfave/cli** |
|---|---|---|
| **Description** | Command-based CLI framework; used by kubectl, hugo, gh | Simpler CLI framework; flag-centric |
| **Pros** | Subcommands; completions; man pages; code gen; massive adoption | Simpler API; less boilerplate for small CLIs |
| **Cons** | Heavy for simple tools; code gen creates noise | Less powerful; fewer features |
| **When to use** | Any CLI with subcommands | Simple single-command tools |
| **Verdict** | ✅ **Industry standard** | ✅ Good for small tools |

### Observability

| | **OpenTelemetry** | **prometheus/client_golang** | **Datadog / vendor SDKs** |
|---|---|---|---|
| **Description** | Vendor-neutral traces, metrics, logs SDK | Prometheus metrics client; exposition format | Vendor-specific telemetry SDKs |
| **Pros** | Vendor-agnostic; CNCF standard; traces + metrics + logs | Simple; proven; huge ecosystem; pull-based | Turnkey; vendor support; advanced features |
| **Cons** | Large API surface; still maturing; complexity | Metrics only; no traces; push requires adapter | Vendor lock-in; cost |
| **When to use** | New services; when vendor flexibility matters | Prometheus/Grafana stack; Kubernetes native | When vendor is already chosen and paid for |
| **Verdict** | ✅ **Future standard** | ✅ Proven for metrics | ⚠️ Pragmatic but locked in |

### Dependency Injection

| | **Manual (main wiring)** | **wire (Google)** | **fx (Uber)** | **dig (Uber)** |
|---|---|---|---|---|
| **Description** | Explicit constructor calls in `main()` | Compile-time DI via code generation | Runtime DI container; lifecycle management | Runtime DI container (fx is built on dig) |
| **Pros** | Explicit; debuggable; no magic; Go-idiomatic | Compile-time errors; no runtime overhead | Lifecycle hooks; module system; good for large apps | Simpler than fx; just dependency resolution |
| **Cons** | Verbose `main()` for large apps; manual ordering | Code gen step; limited adoption; learning curve | Runtime errors; magic; opaque; hard to debug | Runtime errors; reflection-based |
| **When to use** | Most services (up to ~30 dependencies) | Large services where manual wiring is painful | Very large services with complex lifecycle | When you want dig without fx's lifecycle |
| **Verdict** | ✅ **Recommended default** | ✅ Good for large apps | ⚠️ Uber-style only | ⚠️ Uber-style only |

---

## 3. The Go Community's Strong Opinions

These aren't just conventions — they're battle-tested principles with real
production reasoning behind each one.

### The Opinions, Ranked by Consensus Strength

```
┌────┬──────────────────────────────────────┬──────────────────────────────────────────────┐
│ ## │ Opinion                              │ Why                                          │
├────┼──────────────────────────────────────┼──────────────────────────────────────────────┤
│  1 │ go test -race in CI is              │ Data races are UNRECOVERABLE fatals in Go.   │
│    │ non-negotiable                       │ Race detector has ~zero false positives.     │
│    │                                      │ Cost: 5-10x slowdown — worth it.             │
├────┼──────────────────────────────────────┼──────────────────────────────────────────────┤
│  2 │ Accept interfaces, return           │ Interfaces in params = flexible, testable.   │
│    │ concrete types                       │ Concrete returns = discoverable, no alloc.   │
│    │                                      │ Consumer defines the interface, not provider.│
├────┼──────────────────────────────────────┼──────────────────────────────────────────────┤
│  3 │ Context as first parameter           │ Convention enforced by go vet. Enables       │
│    │ (ctx context.Context)                │ cancellation, timeouts, tracing propagation. │
│    │                                      │ Never store context in a struct.             │
├────┼──────────────────────────────────────┼──────────────────────────────────────────────┤
│  4 │ Prefer stdlib when possible          │ "A little copying is better than a little    │
│    │                                      │ dependency." Fewer transitive deps = fewer   │
│    │                                      │ supply chain attacks, faster builds.         │
├────┼──────────────────────────────────────┼──────────────────────────────────────────────┤
│  5 │ Small interfaces (1-3 methods)       │ "The bigger the interface, the weaker the    │
│    │                                      │ abstraction." io.Reader (1 method) is the    │
│    │                                      │ most composed interface in the stdlib.        │
├────┼──────────────────────────────────────┼──────────────────────────────────────────────┤
│  6 │ Table-driven tests                   │ Reduces test boilerplate; easy to add cases; │
│    │                                      │ subtests enable t.Parallel() and -run filter.│
│    │                                      │ The Go standard, used throughout stdlib.     │
├────┼──────────────────────────────────────┼──────────────────────────────────────────────┤
│  7 │ Structured logging with slog         │ JSON logs are machine-parseable. slog is now │
│    │                                      │ stdlib (Go 1.21+). No reason for logrus      │
│    │                                      │ in new projects.                             │
├────┼──────────────────────────────────────┼──────────────────────────────────────────────┤
│  8 │ No DI frameworks needed              │ Wire it in main(). Explicit beats magic.     │
│    │                                      │ Constructor injection is simple and testable.│
│    │                                      │ If main() is too big, use wire (code gen).   │
├────┼──────────────────────────────────────┼──────────────────────────────────────────────┤
│  9 │ Avoid ORM in Go                      │ ORMs fight Go's philosophy. Use sqlc (code   │
│    │                                      │ gen) or sqlx (scanning helper). Write SQL.   │
│    │                                      │ Go devs WANT to see the queries.             │
├────┼──────────────────────────────────────┼──────────────────────────────────────────────┤
│ 10 │ Avoid global state                   │ Globals make testing hard, create hidden      │
│    │                                      │ coupling, prevent parallelism. Pass deps     │
│    │                                      │ through constructors. Even loggers.           │
└────┴──────────────────────────────────────┴──────────────────────────────────────────────┘
```

### The Interface Location Rule

This deserves special attention because it's counter-intuitive for Java/C# developers:

```
  ❌ Java/C# style: interface defined by the PROVIDER (the package that implements it)

  ┌──── provider/repo.go ────┐     ┌──── consumer/service.go ────┐
  │ type UserRepo interface { │     │ import "provider"            │
  │   Get(id int) (*User, err)│     │                              │
  │   Save(u *User) error     │◄────│ func NewService(             │
  │ }                         │     │   repo provider.UserRepo,    │
  │ type PGRepo struct{...}   │     │ ) *Service                   │
  │ func (p *PGRepo) Get(...) │     │                              │
  └───────────────────────────┘     └──────────────────────────────┘
  Consumer depends on provider package — tight coupling!

  ✅ Go style: interface defined by the CONSUMER (the package that uses it)

  ┌──── repo/pg.go ──────────┐     ┌──── service/service.go ──────┐
  │ type PGRepo struct{...}  │     │ type UserStore interface {    │
  │ func (p *PGRepo) Get(...)|     │   Get(id int) (*User, error)  │
  │ func (p *PGRepo) Save(..)|     │ }                             │
  └──────────────────────────┘     │ func NewService(              │
          │                         │   store UserStore,            │
          │ satisfies implicitly    │ ) *Service                    │
          └─────────────────────────►                              │
                                    └──────────────────────────────┘
  Consumer defines only the methods IT needs — decoupled!
```

---

## 4. Project Structure for Enterprise

### The Standard Layout

```
myservice/
├── cmd/                            # Entry points — one main package per binary
│   ├── api-server/
│   │   └── main.go                 # Wires dependencies, starts HTTP server
│   └── worker/
│       └── main.go                 # Wires dependencies, starts job consumer
│
├── internal/                       # Private packages — compiler-enforced boundary
│   ├── domain/                     # Business entities + repository interfaces
│   │   ├── user.go                 # type User struct; type UserStore interface
│   │   └── order.go                # type Order struct; type OrderStore interface
│   │
│   ├── service/                    # Business logic — orchestrates domain
│   │   ├── user_service.go         # func NewUserService(store domain.UserStore)
│   │   └── order_service.go
│   │
│   ├── repository/                 # Data access — implements domain interfaces
│   │   ├── postgres/
│   │   │   ├── user_repo.go        # type PGUserRepo struct; implements domain.UserStore
│   │   │   └── order_repo.go
│   │   └── redis/
│   │       └── cache.go
│   │
│   ├── handler/                    # HTTP/gRPC transport — depends on service
│   │   ├── user_handler.go         # Routes, request parsing, response writing
│   │   ├── middleware.go           # Logging, auth, recovery middleware
│   │   └── router.go              # Assembles routes + middleware
│   │
│   └── platform/                   # Cross-cutting concerns
│       ├── config/                 # Config loading + validation
│       ├── logging/                # slog setup + middleware
│       ├── metrics/                # Prometheus/OTel registration
│       └── database/              # Connection pool setup
│
├── pkg/                            # Public reusable packages (use sparingly!)
│   └── httputil/                   # Shared HTTP helpers (only if truly reusable)
│
├── api/                            # API definitions
│   ├── openapi.yaml                # REST API spec
│   └── proto/                      # gRPC protobuf definitions
│       └── user/v1/user.proto
│
├── migrations/                     # Database migrations (goose, golang-migrate)
│   ├── 001_create_users.up.sql
│   └── 001_create_users.down.sql
│
├── deployments/                    # Infrastructure
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── k8s/
│       ├── deployment.yaml
│       └── service.yaml
│
├── go.mod
├── go.sum
└── Makefile
```

### The Dependency Rule

Dependencies flow inward — outer layers depend on inner layers, never the reverse:

```
  ┌─────────────────────────────────────────────────────────────┐
  │                        cmd/main.go                          │
  │                    (composition root)                        │
  │  Wires everything: repo → service → handler → http.Server  │
  └─────────────┬──────────────┬──────────────┬─────────────────┘
                │              │              │
                ▼              ▼              ▼
  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
  │   handler/   │  │   service/   │  │  repository/ │
  │  depends on  │──►  depends on  │──►  implements  │
  │  service     │  │  domain      │  │  domain      │
  └──────────────┘  └──────┬───────┘  └──────┬───────┘
                           │                  │
                           ▼                  ▼
                    ┌──────────────────────────────┐
                    │          domain/              │
                    │  entities + interfaces        │
                    │  DEPENDS ON NOTHING           │
                    └──────────────────────────────┘

  Key rules:
  ├─ domain/ imports NOTHING from internal/ (no circular deps)
  ├─ service/ imports domain/ only (never handler/ or repository/)
  ├─ repository/ imports domain/ (implements its interfaces)
  ├─ handler/ imports service/ (calls business logic)
  └─ cmd/main.go imports everything (the composition root)
```

### Why `internal/` Matters — Compiler-Enforced Privacy

The `internal/` directory is **not a convention — it's enforced by the Go toolchain**.
The compiler refuses to build if you violate it. No linter needed, no config file — it's
baked into `cmd/go/internal/load/pkg.go`.

```
  myservice/internal/domain/user.go  → importable by myservice/cmd/api-server  ✅
  myservice/internal/domain/user.go  → NOT importable by ANY other module      ❌

  Compiler error:
  "use of internal package myservice/internal/domain not allowed"
```

#### The Parent-Scoping Rule

This is the part most engineers miss. The `internal/` boundary is scoped to its
**immediate parent directory**, not the module root:

```
  mymodule/
  ├── internal/                    ← parent = mymodule/ (module root)
  │   └── config/                  ← ANY package in mymodule can import
  │
  ├── api/
  │   ├── internal/                ← parent = api/
  │   │   └── validator/           ← ONLY api/ subtree can import
  │   └── handler/                 ← ✅ can import api/internal/validator
  │
  ├── worker/
  │   ├── internal/                ← parent = worker/
  │   │   └── retry/               ← ONLY worker/ subtree can import
  │   └── processor/               ← ✅ can import worker/internal/retry
  │
  └── cmd/
      └── server/                  ← ✅ can import internal/config
                                   ← ❌ CANNOT import api/internal/validator
                                   ← ❌ CANNOT import worker/internal/retry
```

The toolchain walks the import path segments. If `internal` appears in the path,
it checks whether the importing package's directory is a descendant of `internal/`'s
parent. If not — **compile error, full stop**.

#### Layered Privacy in Practice

This enables **team-level boundaries** in large codebases:

```
  ┌─────────────────────────────────────────────────────────┐
  │                  internal/config/                        │
  │              module-wide shared config                   │
  │         (any package in the module can import)           │
  └───────────────┬─────────────────────┬───────────────────┘
                  │                     │
    ┌─────────────▼──────────┐ ┌───────▼─────────────────┐
    │      api/ subtree      │ │     worker/ subtree      │
    │                        │ │                           │
    │  api/internal/         │ │  worker/internal/         │
    │    validator/          │ │    retry/                  │
    │    ratelimit/          │ │    deadletter/             │
    │                        │ │                           │
    │  api/handler/   ✅     │ │  worker/processor/  ✅    │
    │  api/middleware/ ✅    │ │  worker/consumer/   ✅    │
    └────────────────────────┘ └───────────────────────────┘
           ❌ cross-boundary imports blocked by compiler
```

The API team's validation logic stays invisible to the worker team.
The worker's retry/dead-letter logic stays invisible to the API team.
Shared infrastructure (config, logging) lives in module-root `internal/`.

### The `pkg/` Debate — Convention vs Noise

Unlike `internal/`, the `pkg/` directory has **zero compiler enforcement**.
It's purely a signal to humans: "these packages are designed for external consumption."

```go
  // Both are identical to the compiler:
  import "mymodule/pkg/client"
  import "mymodule/client"
```

#### Why The Community Is Moving Away From `pkg/`

The Go community has shifted against `pkg/`. Here's the timeline:

```
  2014-2018: pkg/ widely adopted (Kubernetes, Docker, Prometheus)
             Rationale: "clear signal of public API"

  2019+:     Pushback from Go team and community
             Russ Cox: "internal/ for private. Everything else is public.
                        pkg/ just adds a useless directory level."

  2020+:     Docker removed their pkg/ directory
             New projects avoid it
             Kubernetes regrets it but can't change (backward compat)
```

#### The Arguments

```
  Pro-pkg/:
  ├─ Clear boundary: "this is our public API" vs implementation details
  ├─ New contributors immediately know what they can depend on
  └─ Prevents accidental exposure of internals

  Anti-pkg/ (winning position):
  ├─ Redundant: if it's NOT in internal/, it's already public
  ├─ Adds noise: deeper import paths for no compiler benefit
  ├─ The Go standard library doesn't use pkg/
  ├─ Russ Cox explicitly recommends against it
  └─ You can't un-adopt it without breaking all importers
```

#### When `pkg/` Still Makes Sense

Despite the debate, `pkg/` can be valuable in **one specific case**:
a module that is primarily a **service** (cmd/) but also exposes a **client SDK**:

```
  myservice/
  ├── cmd/server/main.go            ← the service
  ├── internal/                     ← service implementation (private)
  │   ├── handler/
  │   └── repository/
  └── pkg/                          ← client library (public API)
      └── client/
          └── client.go             ← SDK for other services to call yours
```

Here `pkg/` clearly separates "this is what we export" from "this is our service."
But if your module IS a library (no cmd/), everything is public — `pkg/` is noise.

### When Do These Boundaries Matter?

```
  ┌──────────────────────────────┬────────────────────────────────┐
  │ Scenario                     │ Do you need internal/?          │
  ├──────────────────────────────┼────────────────────────────────┤
  │ Learning/reference repo      │ No — nobody imports your code  │
  │ CLI tool (single binary)     │ Yes — clean architecture       │
  │ Library published on pkg.dev │ Absolutely — API stability     │
  │ Microservice at company      │ Yes — prevents tight coupling  │
  │ Monorepo with multiple teams │ Yes + nested internal/         │
  └──────────────────────────────┴────────────────────────────────┘
```

The key insight: `internal/` is **API stability insurance**. Once you export a package,
changing it is a breaking change. `internal/` lets you iterate freely on implementation
without worrying about external consumers.

---

## 5. Graceful Shutdown — The Production Pattern

Order matters. Get this wrong and you lose in-flight requests or corrupt data.

```
  Signal received (SIGTERM/SIGINT)
  │
  ├─ 1. Stop accepting new connections
  │     └─ http.Server.Shutdown(ctx) — stops listener, waits for in-flight
  │
  ├─ 2. Drain in-flight requests
  │     └─ Shutdown waits until all active handlers return (or ctx deadline)
  │
  ├─ 3. Stop background workers
  │     └─ Cancel worker context → workers finish current job → exit
  │
  ├─ 4. Close external connections
  │     └─ Database pools, Redis, message queues — close in reverse init order
  │
  ├─ 5. Flush buffers
  │     └─ Flush logs, metrics, trace spans
  │
  └─ 6. Exit
        └─ os.Exit(0) — clean shutdown
```

### The Complete Implementation

```go
func main() {
    // ── Setup ──────────────────────────────────────────────────────────
    cfg := config.MustLoad()
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    db, err := sql.Open("pgx", cfg.DatabaseURL)
    if err != nil {
        logger.Error("failed to connect to database", "error", err)
        os.Exit(1)
    }

    svc := service.NewUserService(repository.NewPGUserRepo(db))
    handler := handler.NewRouter(svc, logger)

    srv := &http.Server{
        Addr:         cfg.Addr,
        Handler:      handler,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    // ── Start server in background ────────────────────────────────────
    go func() {
        logger.Info("server starting", "addr", cfg.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("server error", "error", err)
            os.Exit(1)
        }
    }()

    // ── Wait for shutdown signal ──────────────────────────────────────
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()
    <-ctx.Done()

    logger.Info("shutdown signal received, draining...")

    // ── Graceful shutdown sequence ────────────────────────────────────
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 1. Stop accepting new connections + drain in-flight requests
    if err := srv.Shutdown(shutdownCtx); err != nil {
        logger.Error("server shutdown error", "error", err)
    }

    // 2. Close database connections
    if err := db.Close(); err != nil {
        logger.Error("database close error", "error", err)
    }

    logger.Info("shutdown complete")
}
```

### Key Details

| Concern | Setting | Why |
|---|---|---|
| `ReadTimeout` | 5s | Prevent slowloris attacks |
| `WriteTimeout` | 10s | Bound total response time |
| `IdleTimeout` | 120s | Reuse connections but don't hold forever |
| Shutdown timeout | 30s | Kubernetes sends SIGKILL after `terminationGracePeriodSeconds` (default 30s) |
| `http.ErrServerClosed` | Check in ListenAndServe error | Normal during shutdown — not an error |

---

## 6. Docker Multi-Stage Build — The Standard

### The Dockerfile

```dockerfile
# ── Stage 1: Build ─────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Build with all optimizations
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/server ./cmd/api-server

# ── Stage 2: Runtime ───────────────────────────────────────────
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/server /server

# Container-aware runtime settings
ENV GOMEMLIMIT=450MiB
ENV GOMAXPROCS=2

EXPOSE 8080

ENTRYPOINT ["/server"]
```

### Build Flags Explained

```
CGO_ENABLED=0        → Pure Go binary, no C dependencies, works on scratch/distroless
GOOS=linux           → Cross-compile target OS
GOARCH=amd64         → Target architecture (arm64 for Graviton/M-series)
-ldflags="-s -w"     → Strip debug symbols and DWARF info (~30% smaller binary)
```

### Container Runtime Settings

```
┌──────────────────────────────────────────────────────────────────────────┐
│ GOMEMLIMIT (Go 1.19+)                                                   │
│                                                                          │
│ Sets a soft memory limit for the Go runtime. GC runs more aggressively   │
│ as heap approaches this limit. Set to ~90% of container memory limit.    │
│                                                                          │
│   Container limit: 512 MiB → GOMEMLIMIT=450MiB                          │
│   Container limit: 1 GiB   → GOMEMLIMIT=900MiB                          │
│                                                                          │
│ Without this, Go has no idea it's in a container and may OOM.            │
├──────────────────────────────────────────────────────────────────────────┤
│ GOMAXPROCS                                                               │
│                                                                          │
│ Number of OS threads for goroutine scheduling. Defaults to all host CPUs │
│ — in a container, this sees HOST CPUs, not the cgroup limit!             │
│                                                                          │
│ Fix: use uber-go/automaxprocs (reads cgroup CPU limit at startup):       │
│   import _ "go.uber.org/automaxprocs"                                    │
│                                                                          │
│ Or set manually: GOMAXPROCS=2 for a 2-CPU container.                     │
├──────────────────────────────────────────────────────────────────────────┤
│ Health Check Endpoints                                                   │
│                                                                          │
│   /healthz  (liveness)  → "am I alive?" → restart if unhealthy           │
│   /readyz   (readiness) → "can I serve?" → remove from load balancer     │
│                                                                          │
│   Readiness should check: DB connection, cache connection, config loaded │
│   Liveness should be cheap: return 200 OK (don't check dependencies)     │
└──────────────────────────────────────────────────────────────────────────┘
```

### Image Size Comparison

| Base Image | Size | Security | Use When |
|---|---|---|---|
| `golang:1.23` | ~800 MB | Full OS, many CVEs | Never for production |
| `alpine:3.19` | ~7 MB + binary | Minimal, musl libc | Need shell for debugging |
| `distroless/static` | ~2 MB + binary | No shell, no package manager | **Recommended default** |
| `scratch` | 0 MB + binary | Nothing at all | When you handle TLS certs yourself |

---

## 7. Middleware Pattern — The Go Way

The universal middleware signature creates a composable handler chain:

```go
type Middleware func(next http.Handler) http.Handler
```

```
  Request → [ Logging ] → [ Auth ] → [ RateLimit ] → [ Handler ] → Response
            (unwinding on the way back — each middleware runs "after" logic)
```

**Essential middleware for production:** Recovery (panic→500), RequestID
(trace correlation), Logging (structured request/response), Timeout
(`http.TimeoutHandler`), Auth (JWT/API key validation), RateLimit.

Compose with a `Chain` function that wraps handlers right-to-left:

```go
stack := Chain(Recovery(logger), RequestID, Logging(logger), Timeout(30*time.Second))
mux.Handle("GET /users/{id}", stack(userHandler))
```

> **Full deep dive:** See [Chapter 22 — The Middleware Pattern](./22_middleware_pattern.md)
> for how `http.HandlerFunc`, closures, and structural typing make this work,
> plus ResponseWriter wrapping and cross-language comparison.
> See also [Chapter 23 §1](./23_production_patterns.md) for architectural context.

---

## 8. Configuration Management

### The 12-Factor Approach

```
┌──────────────────────────────────────────────────────────────────────┐
│                     Configuration Hierarchy                          │
│                     (highest priority wins)                          │
│                                                                      │
│  Priority 1: Environment variables        ← secrets, per-env overrides│
│  Priority 2: Config file (YAML/JSON/TOML) ← non-sensitive defaults   │
│  Priority 3: Struct field defaults         ← sane fallbacks          │
│                                                                      │
│  Rule: NEVER commit secrets to code or config files.                 │
│  Rule: Config should be validated at startup, not at first use.      │
└──────────────────────────────────────────────────────────────────────┘
```

### The Config Struct Pattern

```go
type Config struct {
    Addr         string        `env:"ADDR"          default:":8080"`
    DatabaseURL  string        `env:"DATABASE_URL"  required:"true"`
    ReadTimeout  time.Duration `env:"READ_TIMEOUT"  default:"5s"`
    WriteTimeout time.Duration `env:"WRITE_TIMEOUT" default:"10s"`
    LogLevel     string        `env:"LOG_LEVEL"     default:"info"`
    MaxWorkers   int           `env:"MAX_WORKERS"   default:"10"`
}

func MustLoad() Config {
    var cfg Config
    if err := envconfig.Process("", &cfg); err != nil {
        log.Fatalf("config: %v", err)           // fail FAST at startup
    }
    if err := cfg.Validate(); err != nil {
        log.Fatalf("config validation: %v", err)
    }
    return cfg
}

func (c Config) Validate() error {
    if c.MaxWorkers < 1 || c.MaxWorkers > 1000 {
        return fmt.Errorf("max_workers must be 1-1000, got %d", c.MaxWorkers)
    }
    if c.ReadTimeout <= 0 {
        return fmt.Errorf("read_timeout must be positive")
    }
    return nil
}
```

### Feature Flags — Simple Pattern

```go
type FeatureFlags struct {
    EnableNewCheckout bool `env:"FF_NEW_CHECKOUT" default:"false"`
    EnableBetaAPI     bool `env:"FF_BETA_API"     default:"false"`
    MaxBatchSize      int  `env:"FF_MAX_BATCH"    default:"100"`
}

// In handler:
if cfg.Features.EnableNewCheckout {
    return newCheckoutHandler(w, r)
}
return legacyCheckoutHandler(w, r)
```

For production feature flag systems at scale, use a dedicated service (LaunchDarkly,
Unleash, Flipt) that supports gradual rollouts, user targeting, and kill switches
without redeployment.

### Hot Reload with `atomic.Value`

```go
var globalCfg atomic.Value  // stores *Config

func init() {
    globalCfg.Store(MustLoad())
}

func ReloadConfig() error {
    newCfg, err := Load()
    if err != nil {
        return fmt.Errorf("reload: %w", err)
    }
    globalCfg.Store(newCfg)         // atomic swap — no lock needed for readers
    return nil
}

func GetConfig() *Config {
    return globalCfg.Load().(*Config)   // lock-free read
}
```

**Caveat:** Hot reload works for feature flags and tuning parameters. Do NOT hot-reload
database URLs or listen addresses — those require a restart.

---

## 9. Quick Reference Card

### Before You Ship to Production — Checklist

```
┌──────────────────────────────────────────────────────────────────────┐
│                    PRODUCTION READINESS CHECKLIST                     │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  BUILD & TEST                                                        │
│  ☐ go test -race ./...              Race detector in CI (ALWAYS)     │
│  ☐ go vet ./...                     Static analysis                  │
│  ☐ golangci-lint run                Comprehensive linting            │
│  ☐ go test -count=1 ./...           No cached test results in CI     │
│  ☐ go test -cover ./...             Know your coverage (don't fake)  │
│                                                                      │
│  RUNTIME CONFIGURATION                                               │
│  ☐ GOMEMLIMIT set                   ~90% of container memory limit   │
│  ☐ GOMAXPROCS set or automaxprocs   Match container CPU limit        │
│  ☐ Config validated at startup      Fail fast, not on first request  │
│                                                                      │
│  RESILIENCE                                                          │
│  ☐ Graceful shutdown implemented    SIGTERM → drain → close → exit   │
│  ☐ Panic recovery middleware        Every goroutine, every handler   │
│  ☐ Context propagation              Timeouts and cancellation chain  │
│  ☐ HTTP client timeouts set         Never use http.DefaultClient     │
│  ☐ Database connection pool sized   SetMaxOpenConns, SetMaxIdleConns │
│  ☐ Retry with backoff               Exponential backoff on transient │
│                                                                      │
│  OBSERVABILITY                                                       │
│  ☐ Structured logging (slog/zap)    JSON format, request-scoped      │
│  ☐ Metrics exported                 RED: Rate, Errors, Duration      │
│  ☐ Health endpoints                 /healthz (live), /readyz (ready) │
│  ☐ pprof enabled (behind auth)      /debug/pprof/* for production    │
│  ☐ Request IDs propagated           Correlation across services      │
│                                                                      │
│  SECURITY                                                            │
│  ☐ TLS configured (MinVersion 1.2)  crypto/tls with sane defaults    │
│  ☐ Input validated at boundary      Never trust user input           │
│  ☐ SQL parameterized                $1, $2 — never string concat     │
│  ☐ Secrets in env vars / vault      Never in code or config files    │
│  ☐ Distroless/scratch Docker image  Minimal attack surface           │
│  ☐ Rate limiting enabled            Per-client, per-endpoint         │
│                                                                      │
│  GOROUTINE HYGIENE                                                   │
│  ☐ No unbounded goroutine spawning  Use semaphore / worker pool      │
│  ☐ Every goroutine has exit path    Context cancellation or signal   │
│  ☐ uber-go/goleak in tests          Detects goroutine leaks          │
│  ☐ No sync type copying             go vet catches this              │
│                                                                      │
│  PERFORMANCE BASELINE                                                │
│  ☐ Benchmarks for hot paths         go test -bench=. -benchmem       │
│  ☐ Escape analysis reviewed         go build -gcflags='-m' on hot    │
│  ☐ Slices/maps pre-allocated        Reduce GC pressure               │
│  ☐ sync.Pool for high-churn objects Buffers, temporary structs       │
│                                                                      │
│  DEPLOYMENT                                                          │
│  ☐ Multi-stage Docker build         Builder → distroless/scratch     │
│  ☐ Kubernetes probes configured     livenessProbe + readinessProbe   │
│  ☐ terminationGracePeriodSeconds    Match your shutdown timeout       │
│  ☐ Resource limits set              CPU + memory limits in k8s       │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

### Quick Debugging Commands

```bash
# Goroutine dump (find leaks)
curl http://localhost:6060/debug/pprof/goroutine?debug=2

# CPU profile (30 seconds)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Heap profile (current allocations)
go tool pprof http://localhost:6060/debug/pprof/heap

# GC trace (see GC pauses and frequency)
GODEBUG=gctrace=1 ./myservice

# Scheduler trace (see goroutine scheduling)
GODEBUG=schedtrace=1000 ./myservice

# Execution trace (visual timeline)
curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5
go tool trace trace.out

# Race detector (in tests)
go test -race -count=1 ./...

# Escape analysis (what allocates on heap)
go build -gcflags='-m -m' ./... 2>&1 | grep 'escapes to heap'
```

### The "Never Do This" List

```
┌──────────────────────────────────────────────────────────────────────┐
│  ❌ NEVER use http.DefaultClient in production                       │
│     → No timeouts. A hanging server blocks your goroutine forever.   │
│     → Create: &http.Client{Timeout: 10 * time.Second}               │
│                                                                      │
│  ❌ NEVER use http.ListenAndServe(addr, handler) directly            │
│     → No timeouts, no graceful shutdown.                             │
│     → Create: &http.Server{...} with all timeouts set.              │
│                                                                      │
│  ❌ NEVER log errors at every layer                                  │
│     → Log ONCE at the top (handler). Wrap with context at each layer.│
│     → return fmt.Errorf("repo.GetUser: %w", err) — don't log here.  │
│                                                                      │
│  ❌ NEVER use init() for complex setup                               │
│     → init() runs at import time, before main(). Hard to test.       │
│     → Put setup in main() or constructors.                           │
│                                                                      │
│  ❌ NEVER store context.Context in a struct field                     │
│     → Context is request-scoped. Struct fields outlive requests.     │
│     → Pass ctx as the first parameter to methods.                    │
│                                                                      │
│  ❌ NEVER use panic for expected error conditions                     │
│     → Panic is for programmer errors (bug in the code).              │
│     → Return errors for operational failures (network, disk, input). │
│                                                                      │
│  ❌ NEVER ignore the error from rows.Close() or resp.Body.Close()    │
│     → Wrap in a helper or use named returns to capture.              │
│                                                                      │
│  ❌ NEVER launch a goroutine without knowing how it stops            │
│     → Every go func() needs a cancellation path (context, channel,   │
│       or done signal).                                               │
└──────────────────────────────────────────────────────────────────────┘
```

---

> **Go Proverb:** *"Don't just check errors, handle them gracefully."*
>
> Production Go isn't about writing clever code — it's about writing code that
> fails predictably, recovers gracefully, and is debuggable at 3 AM with only
> `pprof` and structured logs. Every pattern in this document exists because
> someone learned it the hard way.

---

*Further reading:*
- `runtime/proc.go` — goroutine scheduler internals
- `runtime/map.go` — map implementation and why it never shrinks
- `net/http/server.go` — HTTP server shutdown implementation
- [Go Proverbs](https://go-proverbs.github.io/) — Rob Pike's design philosophy
- [Effective Go](https://go.dev/doc/effective_go) — official Go style guide
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) — enterprise patterns
