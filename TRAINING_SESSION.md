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

### Exercises in `fundamentals/07_error_handling/exercises.go` 🔲

All 5 exercises have been **reset to stubs** — implement them yourself! Run tests with:
```bash
go test -v ./fundamentals/07_error_handling/
```

**Exercises overview:**
1. `Divide` — Basic `(value, error)` return pattern
2. `Validate` — Custom error type with `*ValidationError`
3. `SafeGet` — Descriptive errors with `fmt.Errorf`
4. `FindUser` — Sentinel errors (`ErrUserNotFound`, `ErrAccessDenied`)
5. `WrapError` — Error wrapping with `%w` to preserve identity

Reference solutions are in `solutions.go` if you get stuck.

### 🆕 Deep Dive: Error Recovery & Retry — `practical/08_error_recovery_retry/`

Coming from Java, you're used to `try-catch-finally`. Go doesn't have that — and it's intentional.
This module bridges the mental model gap with 6 exercises covering:

| Exercise | Java Equivalent | Go Pattern |
|----------|----------------|-----------|
| `SafeExecute` | `try { fn() } catch (Exception e)` | `defer func() { recover() }()` |
| `SafeGoroutine` | Thread.setUncaughtExceptionHandler | Recovery in spawned goroutines |
| `Retry` | `@Retryable(maxAttempts=3)` | Manual retry loop |
| `RetryWithBackoff` | Resilience4j / Spring Retry | Exponential backoff with jitter |
| `RetryClassified` | Catching specific exception types | `PermanentError` marker pattern |
| `RetryWithContext` | `Future.cancel()` / timeout | `context.Context` + `select` |

```bash
# Read concepts first, then implement exercises
cat practical/08_error_recovery_retry/concepts.go
go test -v ./practical/08_error_recovery_retry/
```

---

## Session 4: Slices & Maps — The Building Blocks of Every Algorithm ✅

### Slices — Not Arrays

A slice is a **view into an underlying array**. It has three fields:
```
slice header = { pointer to array, length, capacity }
```

This is the most important thing to understand:
```go
a := []int{1, 2, 3, 4, 5}
b := a[1:4]  // b = [2, 3, 4]  ← shares the SAME underlying array as a
b[0] = 99
fmt.Println(a) // [1, 99, 3, 4, 5] — a was modified! b is a VIEW.
```

**Interview trap**: When does `append` create a new array?
```go
s := make([]int, 3, 5)  // len=3, cap=5 — room for 2 more without realloc
s = append(s, 4, 5)     // uses existing array (cap was enough)
s = append(s, 6)        // CAP EXCEEDED → new array allocated, old array abandoned
```
Rule: **After append, always assume the underlying array may have changed.**

### Key Slice Patterns

| Pattern | Use Case | Code |
|---|---|---|
| Two Pointers | Reverse, palindrome check | `left, right := 0, len(s)-1` |
| Write Pointer | Remove dups in-place | `write := 0; s[write] = s[read]; write++` |
| Sliding Window | Subarray sum/max | Expand right, shrink left |
| Three Reversal | Rotate in-place | reverse all → reverse parts |

### Maps — Hashmaps with Zero-Value Magic

Go maps have a useful property: accessing a missing key returns the **zero value** for the value type — no crash, no error.
```go
var counts map[string]int  // nil map — reading is safe, writing PANICS
counts["x"]++              // PANIC: assignment to nil map

counts = make(map[string]int) // initialized map
counts["x"]++              // safe: counts["x"] starts as 0, becomes 1
```

**The map toolkit** — patterns every senior engineer knows:
```go
// 1. Frequency counter — most common map pattern
freq := make(map[rune]int)
for _, ch := range s { freq[ch]++ }

// 2. Existence check — the "comma ok" idiom
if v, ok := m[key]; ok { /* key exists, v has value */ }

// 3. Map of slices — group items
groups := make(map[byte][]string)
groups[word[0]] = append(groups[word[0]], word)  // append to nil is safe

// 4. Seen set — use struct{} for memory efficiency, bool for readability
seen := make(map[int]bool)
if seen[n] { return n }
seen[n] = true
```

**Interview trap**: Are maps thread-safe?
**No.** Concurrent reads are safe, but concurrent read+write or write+write causes a panic (`concurrent map read and map write`). Use `sync.RWMutex` or `sync.Map` for concurrent maps.

---

## Session 5: Packages, Modules & Visibility ✅

### Visibility — Capital Letter IS the Access Modifier

Go has no `public` / `private` / `protected` keywords. The rule is simple:
```
Starts with Uppercase → exported (visible from any package)
Starts with lowercase → unexported (visible only within its package)
```

This applies to **everything**: functions, types, struct fields, methods, constants, variables.
```go
type User struct {
    Name  string  // exported — JSON/ORM can see this
    email string  // unexported — internal only, JSON ignores it
}
```

### init() — Automatic Setup

Every Go file can have an `init()` function. It runs **automatically** when the package is imported, before `main()`.

```go
var db *sql.DB

func init() {
    db = connectToDatabase()  // runs automatically on import
}
```

**Order of execution:**
```
1. Package-level var declarations (in declaration order)
2. init() functions (in import order, then file order)
3. main()
```

**When to use init()**:
- Registering database drivers: `import _ "github.com/lib/pq"`
- Registering HTTP handlers, codecs, plugins
- Setting up package-level defaults

**When NOT to use init()**: Do not do heavy work or anything that can fail (no error return from init).

### Blank Import — Side-Effect Only Import

```go
import _ "image/png"  // Registers PNG decoder via init() — no direct use
import _ "github.com/lib/pq"  // Registers PostgreSQL driver
```
The `_` tells the compiler: "I know I'm not using this package's exports, but run its init()."

### Module Path = Your Package Address

```
go.mod:  module github.com/mycompany/myservice

Package at:  myservice/internal/config/loader.go
Import path: github.com/mycompany/myservice/internal/config
```

---

## Coming Up — Session 6: System Design Deep Dive

### Topics:
1.  **Context propagation** — cancellation, timeouts, and request-scoped values.
2.  **Standard Project Layout** — how senior engineers structure a real Go service.
3.  **Dependency Injection** — making code testable without a framework.
4.  **Rate Limiting** — the Token Bucket pattern implemented in Go.

---
---

# ═══════════════════════════════════════════════
# SYSTEM DESIGN INTERVIEW TRACK (Language Agnostic)
# ═══════════════════════════════════════════════

> These sessions prepare you for the **45-60 minute system design interview**.
> Every major tech company (FAANG, unicorns, mid-size) asks these.
> The goal is not to memorize answers — it's to **think out loud like a senior engineer**.

---

## SD Session 1: The Interview Framework — How to Approach Any System Design Question ✅

### The #1 Mistake Candidates Make

They jump straight into drawing boxes: "OK so we'll have a load balancer, then a fleet of servers, then a Postgres database..."

This is **wrong**. You haven't understood what you're building yet.

Senior engineers always clarify requirements first. An interviewer who asks "design Twitter" might actually want you to design only the "post tweet + news feed" part, or only the search, or maybe a read-heavy analytics system. You won't know unless you ask.

### The 5-Step Framework (Use This Every Time)

```
Step 1 — Clarify Requirements        (~5 min)
Step 2 — Estimate Scale              (~5 min)
Step 3 — Define the API / Data Model (~5 min)
Step 4 — High-Level Design           (~15 min)
Step 5 — Deep Dive on Components     (~20 min)
```

---

### Step 1: Clarify Requirements

Split requirements into two buckets:

**Functional requirements** — what the system *does*:
```
"Design a URL shortener"

Functional:
- User submits a long URL → gets a short URL (e.g., bit.ly/abc123)
- User visits short URL → redirected to original
- (Optional) Custom aliases, expiry dates, analytics
```

**Non-functional requirements** — how the system *performs*:
```
Non-functional:
- How many URLs shortened per day? (100M? 1M?)
- How many redirects per day? (read:write ratio — often 100:1)
- Latency requirement for redirect? (<10ms? <100ms?)
- Availability requirement? (99.9%? 99.99%?)
- Data durability? (can we lose some URLs?)
- Global or single-region?
```

**Questions to ask every time:**
1. What is the scale? (DAU, QPS, data size)
2. Read-heavy or write-heavy?
3. Strong consistency required, or eventual consistency OK?
4. Single region or global?
5. What is the SLA for availability?

---

### Step 2: Estimate Scale (Back-of-the-Envelope)

Interviewers want to see that you can reason about numbers. You don't need exact answers — **orders of magnitude** are fine.

**Memorize these reference numbers:**

| Resource | Latency |
|---|---|
| L1 cache hit | 0.5 ns |
| L2 cache hit | 7 ns |
| RAM read | 100 ns |
| SSD read | 150 µs |
| HDD read | 10 ms |
| Network within datacenter | 0.5 ms |
| Cross-region network | 150 ms |

| Unit | Value |
|---|---|
| 1 million req/day | ~12 req/sec |
| 1 billion req/day | ~12,000 req/sec |
| 1 KB × 1M records | 1 GB |
| 1 KB × 1B records | 1 TB |

**Example estimation for URL shortener:**
```
Write QPS:   100M URLs/day ÷ 86,400 sec ≈ 1,200 writes/sec
Read QPS:    100:1 read:write ratio → 120,000 reads/sec
Storage:     1 record ≈ 500 bytes (URL + metadata)
             100M/day × 365 days × 10 years = 365B records
             365B × 500 bytes ≈ 180 TB total storage
```

---

### Step 3: Define the API and Data Model

**API first** — forces clarity about what the system does:
```
POST /shorten
  Body: { "url": "https://example.com/very/long/path" }
  Response: { "short_url": "https://bit.ly/abc123" }

GET /{shortCode}
  Response: HTTP 301 Redirect to original URL
```

**Data model** — keep it simple, evolve later:
```
URLs table:
  short_code  VARCHAR(8)   PRIMARY KEY
  long_url    TEXT         NOT NULL
  user_id     BIGINT
  created_at  TIMESTAMP
  expires_at  TIMESTAMP    NULLABLE
```

---

### Step 4: High-Level Design

Draw the minimal architecture that satisfies requirements:
```
Client
  │
  ▼
[Load Balancer]
  │
  ├──► [Write Service] ──► [SQL Database (URLs)]
  │
  └──► [Read Service]  ──► [Cache (Redis)] ──► [SQL Database]
```

At this stage: **don't over-engineer**. Get something working end-to-end, then optimize.

---

### Step 5: Deep Dive

The interviewer will pick one or two components to go deep on. Typical deep-dives:
- "How do you generate the short code?" → Hash vs Counter, collision handling
- "What happens when the cache is full?" → Eviction policies
- "How do you handle 10x traffic spike?" → Auto-scaling, rate limiting
- "What if the database goes down?" → Replication, failover

**The golden rule of deep dives**: Always discuss **trade-offs**. There is no single right answer — show you understand the pros and cons.

---

## SD Session 2: Scalability — Handling Growth ✅

### Vertical vs Horizontal Scaling

**Vertical scaling** (scale up): Give your one server more RAM, CPU, SSD.
```
Pros:  Simple, no code changes, works immediately
Cons:  Hard upper limit (biggest machine on AWS), single point of failure,
       expensive at the top, downtime during upgrade
```

**Horizontal scaling** (scale out): Add more servers.
```
Pros:  Near-unlimited capacity, no single point of failure,
       cheaper at scale, zero-downtime deploys
Cons:  Requires stateless design, coordination overhead,
       distributed system complexity
```

**The rule**: Start vertical. When you hit limits, go horizontal. Most senior design questions assume horizontal.

**For horizontal scaling to work, your services must be stateless:**
```
❌ Stateful — session stored in server memory:
   User logs in → Server A handles it → stores session in memory
   Next request hits Server B → "Who are you?"

✅ Stateless — session stored externally:
   User logs in → Server A stores session in Redis
   Next request hits Server B → fetches session from Redis → works
```

---

### Load Balancing

A Load Balancer (LB) distributes incoming requests across your server pool.

**Layer 4 vs Layer 7 LB:**
```
L4 (Transport Layer): Routes based on IP + TCP port only.
    Fast, can't inspect content. Used for raw TCP (databases, game servers).

L7 (Application Layer): Inspects HTTP headers, URL, cookies.
    Can route /api/* to API servers, /static/* to CDN, sticky sessions.
    This is what Nginx, HAProxy, AWS ALB do.
```

**Load balancing algorithms:**

| Algorithm | How | When to use |
|---|---|---|
| **Round Robin** | Requests rotate through servers 1→2→3→1 | Servers are identical |
| **Least Connections** | Send to server with fewest active connections | Requests have varying processing time |
| **IP Hash** | Hash(client IP) → same server always | Sticky sessions without shared state |
| **Weighted** | Servers get proportional traffic by weight | Mixed server sizes |
| **Random** | Pick a server randomly | Simple, works surprisingly well |

**What the LB does beyond routing:**
- **Health checks**: Removes unhealthy servers from the pool automatically
- **SSL termination**: Decrypts HTTPS once at the LB, plain HTTP to servers (saves CPU)
- **Rate limiting**: Reject clients sending too many requests
- **Connection pooling**: Reuses connections to backends

---

### Stateless Architecture — The Foundation of Scale

```
The Three-Tier Pattern:

[Clients]
    │
[Load Balancer]   ← Routes, terminates SSL, health checks
    │
[App Servers]     ← Stateless, horizontally scaled, interchangeable
    │
[Data Tier]       ← State lives here: DB, Cache, Object Storage, Message Queue
```

Any app server can handle any request because all state is in the data tier. You can:
- Add/remove app servers without coordination
- Deploy new code with rolling restarts (zero downtime)
- Auto-scale based on CPU/memory

---

### When One Load Balancer Isn't Enough

For truly high-scale systems, the LB itself becomes a bottleneck or single point of failure:

```
DNS Load Balancing → Multiple LB clusters → App Servers
    │
    ├── Anycast (BGP): One IP, traffic routes to nearest datacenter
    └── DNS Round Robin: Multiple A records for same domain
```

AWS/GCP handle this for you. In interviews, mention "we'd use a managed LB like AWS ALB in an active-active configuration."

---

## SD Session 3: Caching — Making Everything Faster ✅

### Why Cache?

RAM is ~1000x faster than SSD, and ~10x faster than network.
A cache stores a **hot subset** of your data in memory so most reads never touch the database.

**The 80/20 rule of caching**: 20% of data gets 80% of traffic. Cache that 20% and you eliminate most DB load.

---

### Where to Cache

```
1. Client-Side Cache
   Browser caches HTML/CSS/JS.
   HTTP headers control it: Cache-Control: max-age=3600

2. CDN (Content Delivery Network)
   Edge servers worldwide cache static assets close to users.
   User in Tokyo gets JS from Tokyo, not your US datacenter.

3. Application Cache (e.g., Redis, Memcached)
   In-memory store shared by all app servers.
   Used for: session data, DB query results, computed values.

4. Database Cache
   DB's own query result cache, buffer pool (InnoDB buffer pool).
   Automatic, you don't control it.
```

---

### Cache Write Strategies

This is where most candidates struggle. Caching reads is easy. What about writes?

**Cache-Aside (Lazy Loading)** — Most common pattern:
```
Read:   1. Check cache. Hit → return. Miss → query DB → write to cache → return.
Write:  1. Write to DB. 2. Invalidate (delete) cache entry.

Pros:   Simple, cache only holds data that's actually read.
Cons:   Cache miss = 3 operations. Data briefly stale after write.
        "Cache stampede": many simultaneous misses hit DB at once.
```

**Write-Through** — Write cache and DB synchronously:
```
Write:  1. Write to cache. 2. Write to DB (in the same call).
Read:   Cache is always warm.

Pros:   Cache always consistent with DB.
Cons:   Write latency doubles. Cache fills with rarely-read data.
```

**Write-Behind (Write-Back)** — Write to cache only, DB async:
```
Write:  1. Write to cache. Return immediately.
        2. Background job flushes cache → DB asynchronously.

Pros:   Lowest write latency. Batches DB writes.
Cons:   Risk of data loss if cache crashes before flush.
        Complex to implement correctly.
```

**Which to use?**
- General purpose: **Cache-Aside** (safe default)
- Read-heavy, writes are fine being slow: **Write-Through**
- High-throughput, can tolerate some data loss: **Write-Behind**

---

### Cache Eviction Policies

When the cache is full, something must go:

| Policy | How | Use When |
|---|---|---|
| **LRU** (Least Recently Used) | Evict the item not accessed for the longest time | Access pattern has temporal locality (recent = hot) |
| **LFU** (Least Frequently Used) | Evict the item accessed fewest times | Some items are always hot (celebrity data) |
| **FIFO** | Evict the oldest inserted item | Simple, predictable expiry |
| **TTL** | Evict items after a time-to-live expires | Data that naturally goes stale (rate limits, sessions) |

Redis default: **LRU**. Usually the right choice.

---

### Cache Pitfalls — Interview Gold

**1. Cache Stampede (Thundering Herd)**
```
Popular cache entry expires.
1000 requests hit DB simultaneously.
DB crashes.

Fix: Probabilistic early expiration (re-warm before expiry), or
     Mutex lock: first miss takes the lock, others wait for it.
```

**2. Cache Penetration**
```
Requests for keys that DON'T EXIST in DB bypass cache every time.
(Common in abuse scenarios: user queries random user IDs)

Fix: Cache negative results ("key → null") with a short TTL.
     Or use a Bloom Filter to pre-check if key can possibly exist.
```

**3. Cache Avalanche**
```
Many cache entries expire at the same time.
Massive DB traffic spike.

Fix: Randomize TTLs (instead of all 1h, use 1h ± random minutes).
     Or use "eternal cache + background refresh" pattern.
```

---

### Redis vs Memcached

| Feature | Redis | Memcached |
|---|---|---|
| Data structures | Strings, Lists, Sets, Sorted Sets, Hashes | Strings only |
| Persistence | Yes (RDB + AOF) | No |
| Replication | Yes | No |
| Cluster mode | Yes | Client-side sharding only |
| Pub/Sub | Yes | No |
| Use for | Sessions, leaderboards, rate limits, queues | Pure cache |

**Default choice: Redis.** Memcached only makes sense if you need simpler ops and maximum throughput for pure string caching.

---

## SD Session 4: Databases & Storage — Choosing the Right Tool ✅

### The Most Important Decision: SQL vs NoSQL

**Don't let anyone tell you "NoSQL scales, SQL doesn't."** This was a myth from 2010. **PostgreSQL can handle millions of QPS at petabyte scale with proper sharding/replication.** Choose based on data model and access patterns, not a trend.

---

### SQL (Relational Databases)

**When to use SQL:**
- Your data has clear relationships (users have orders, orders have items)
- You need ACID transactions (financial data, inventory)
- Your queries are complex (JOINs, aggregations, analytics)
- Your schema is stable

**ACID Properties:**
```
Atomicity:   A transaction is all-or-nothing.
             (Transfer $100: debit AND credit, or neither)

Consistency: Every transaction moves the DB from one valid state to another.
             (No negative bank balances)

Isolation:   Concurrent transactions don't interfere with each other.
             (Two people booking the last seat don't both succeed)

Durability:  Committed data survives crashes.
             (Disk write, WAL — Write-Ahead Log)
```

**Indexing — The Most Impactful Optimization:**
```sql
-- Without index: full table scan = O(n)
SELECT * FROM users WHERE email = 'alice@example.com';

-- With index: B-tree lookup = O(log n)
CREATE INDEX idx_users_email ON users(email);
```

**Composite index rule**: The leftmost prefix must be used.
```sql
-- Index on (last_name, first_name)
-- ✅ Uses index: WHERE last_name = 'Smith'
-- ✅ Uses index: WHERE last_name = 'Smith' AND first_name = 'Alice'
-- ❌ Skips index: WHERE first_name = 'Alice'  (leftmost column missing)
```

---

### NoSQL Databases — Four Types

**1. Key-Value Stores** (Redis, DynamoDB, Riak)
```
Data: blob associated with a key. No schema.
Access: Only by key (no queries by value).
Use: Sessions, user profiles, shopping carts, caches.

DynamoDB: key-value + document, fully managed, scales to any load.
```

**2. Document Stores** (MongoDB, CouchDB, Firestore)
```
Data: JSON/BSON documents. Flexible schema.
Access: Query by any field, nested documents.
Use: Catalogs, content management, user preferences, event logs.

Gotcha: No multi-document transactions (historically). MongoDB 4.0+ added them.
```

**3. Wide-Column Stores** (Cassandra, HBase, BigTable)
```
Data: Rows with dynamic columns. Data sorted by row key.
Access: Optimized for sequential reads of huge datasets.
Use: Time-series data, IoT sensor readings, activity logs, inbox.

Cassandra write path: Write to commit log + memtable → flushed to SSTable.
                      No locking. Near-linear write throughput at scale.
```

**4. Graph Databases** (Neo4j, Amazon Neptune)
```
Data: Nodes and edges (relationships are first-class citizens).
Access: Traverse relationships efficiently (what's hard in SQL).
Use: Social networks, recommendation engines, fraud detection.
```

---

### Database Replication

**Why replicate?** High availability + Read scaling.

```
Primary-Replica (Master-Slave):

[Primary]  ──writes──► [Replica 1]
           ──writes──► [Replica 2]  ← reads served here
           ──writes──► [Replica 3]

Writes: Primary only.
Reads: Replicas (scale reads horizontally).
Failover: Replica promoted to primary if primary dies.

Lag: Replication is usually asynchronous → slight staleness on replicas.
```

**Synchronous vs Asynchronous replication:**
```
Sync:  Write confirmed only after replica acknowledges.
       Pro: No data loss. Con: Latency, replica slowdown hurts primary.

Async: Write confirmed when primary writes to disk. Replica catches up later.
       Pro: Fast writes. Con: If primary dies, recent writes lost.

Semi-sync: At least ONE replica must acknowledge. Common default.
```

---

### Database Sharding — Horizontal Partitioning

Sharding splits data across multiple DB nodes. Each shard owns a subset.

**Sharding strategies:**

**Range sharding**: Shard by value range.
```
Shard 1: user_id 1–1M
Shard 2: user_id 1M–2M
Shard 3: user_id 2M–3M

Pro: Range queries fast, simple.
Con: Hot spots (new users always hit last shard = "hot shard" problem).
```

**Hash sharding**: `shard = hash(key) % num_shards`
```
Pro: Even distribution.
Con: Range queries need all shards. Resharding is painful (all keys move).
```

**Consistent Hashing** — The senior engineer's answer to resharding:
```
Place shards on a ring of 360°.
Hash each key → position on ring → go clockwise to find responsible shard.

Add a new shard: Only the next shard's keys move. Other shards unaffected.
Remove a shard: Only that shard's keys move to the next one.

Used by: Cassandra, DynamoDB, Memcached, Riak.
```

---

### ACID vs BASE

| | ACID (SQL) | BASE (NoSQL) |
|---|---|---|
| Stands for | Atomicity, Consistency, Isolation, Durability | Basically Available, Soft State, Eventually Consistent |
| Consistency model | Strong consistency | Eventual consistency |
| Performance | Lower write throughput | Higher write throughput |
| Example | PostgreSQL, MySQL | Cassandra, DynamoDB |

---

## SD Session 5: Distributed Systems — The Hard Problems ✅

### CAP Theorem — What Every Engineer Must Know

**"In a distributed system, you can only guarantee 2 of 3 properties simultaneously:"**

```
C — Consistency:   Every read returns the most recent write (or an error).
A — Availability:  Every request gets a response (no error, possibly stale).
P — Partition Tolerance: System works even when network messages are lost/delayed.
```

**The key insight**: In real distributed systems, **network partitions WILL happen** (a network switch fails, a link is saturated). So P is not optional — you must tolerate partitions. The real choice is:

```
CP: When partitioned, sacrifice availability (return error rather than stale data).
    → Financial systems, inventory, any "correctness over uptime" use case.
    Examples: ZooKeeper, HBase, MongoDB (with writeConcern:majority)

AP: When partitioned, sacrifice consistency (return stale data rather than error).
    → DNS, shopping carts, social media counts, any "availability over accuracy" case.
    Examples: DynamoDB, Cassandra, CouchDB, DNS
```

**Interview answer to "is X CP or AP?"**: Almost nothing is purely one or the other. Most systems let you tune the trade-off per operation (DynamoDB's ConsistencyLevel, Cassandra's quorum reads).

---

### Consistency Models — The Spectrum

Stronger consistency = lower availability and higher latency. These are in order from strongest to weakest:

**Linearizability (Strongest)**
```
Every operation appears instantaneous and takes effect at one point in time.
The gold standard. Used in: Single-master databases, Spanner, etcd.
Cost: High latency, low availability.
```

**Sequential Consistency**
```
All operations appear in some total order. The order respects each process's
individual order, but NOT wall-clock time across processes.
Easier to implement than linearizability.
```

**Causal Consistency**
```
Operations that are causally related appear in causal order.
"If I post a comment and you reply to it, everyone sees my post before your reply."
Used by some distributed databases and version control.
```

**Eventual Consistency (Weakest)**
```
If no new writes happen, all replicas will converge to the same value eventually.
When? Could be milliseconds, could be seconds.
Used by: DNS propagation, CDN, DynamoDB (default), Cassandra (default).
```

---

### Consistent Hashing — In Detail

This solves the resharding problem. Here's the full mental model:

```
Imagine a clock face (ring) numbered 0 to 2^32.

1. Hash each server to a position on the ring:
   Server A → position 90°
   Server B → position 180°
   Server C → position 270°

2. Hash each key to a position on the ring:
   Key "user:123" → position 45°  → goes to Server A (next clockwise)
   Key "user:456" → position 200° → goes to Server C
   Key "user:789" → position 300° → goes to Server A (wraps around)

3. Add Server D at position 135°:
   Only keys between 90° and 135° (previously going to B) now go to D.
   All other keys unaffected. ← This is the big win.

4. Virtual nodes: Each physical server maps to multiple ring positions.
   This prevents hotspots when servers have different capacities.
   Cassandra uses 256 virtual nodes per server by default.
```

---

### The Two Generals Problem — Why Distributed Consensus is Hard

```
General A and General B want to attack at the same time (must coordinate).
They communicate via messengers who might be captured.

General A sends: "Attack at dawn."
Did General B receive it? A doesn't know.
B sends back: "Confirmed."
Did A receive the confirmation? B doesn't know.
A sends: "Confirmed your confirmation."
...

There is no protocol that guarantees both generals act simultaneously
when the communication channel is unreliable.
```

This is why distributed transactions and consensus algorithms are fundamentally hard. It's not a software problem — it's a mathematical impossibility in an unreliable network.

Real systems work around this with:
- **Raft/Paxos consensus** (leader election, log replication)
- **2-Phase Commit** (distributed transactions, with caveats)
- **Saga pattern** (eventual consistency through compensation)

---

## SD Session 6: Messaging & Event-Driven Architecture ✅

### Why Message Queues?

Without queues, services call each other directly (synchronous):
```
User Service → Order Service → Inventory Service → Payment Service
              ↑ if any service is slow or down, the whole chain fails
```

With a queue, services are **decoupled**:
```
User Service → [Queue] → Order Service
                       → Inventory Service (independently)
                       → Notification Service (independently)

Benefits:
- Resilience: If Order Service is down, messages wait in the queue
- Decoupling: Services don't know about each other
- Load leveling: Queue absorbs bursts (Black Friday spike)
- Replay: Reprocess messages for new consumers
```

---

### Point-to-Point vs Pub/Sub

**Point-to-Point (Queue)**:
```
Producer → [Queue] → Consumer
One message is processed by exactly ONE consumer.
Use: Task queues, job processing, email sending.
Examples: SQS, RabbitMQ
```

**Publish/Subscribe (Topic)**:
```
Producer → [Topic] → Consumer A
                   → Consumer B
                   → Consumer C
One message is delivered to ALL subscribers.
Use: Event broadcasting, audit logs, fan-out notifications.
Examples: SNS, Kafka topics, Google Pub/Sub
```

---

### Kafka — The Senior Engineer's Deep Dive

Kafka is not "just a queue" — it's a **distributed, persistent, ordered log**.

**Core concepts:**
```
Topic:      A named stream of records (like a database table name).
Partition:  A topic is split into partitions for parallelism.
            Partition 0: [msg1, msg4, msg7]
            Partition 1: [msg2, msg5, msg8]
            Partition 2: [msg3, msg6, msg9]
Offset:     The position of a message within a partition (immutable, sequential).
Consumer Group: Multiple consumers sharing a topic — each partition assigned to one.
```

**Why Kafka is different from SQS:**
| Feature | Kafka | SQS |
|---|---|---|
| Message retention | Days/weeks (configurable) | 14 days max, deleted after ack |
| Replay | Yes (re-read from any offset) | No |
| Ordering | Per-partition | No guarantee (FIFO queue exists) |
| Consumer groups | Many can read same topic | Queue = one consumer group |
| Throughput | Millions/sec | Thousands/sec |
| Complexity | High (ops burden) | Low (managed) |

---

### Idempotency — The Key to Reliable Messaging

In distributed systems, messages can be delivered **more than once** (at-least-once delivery). Your consumers must handle duplicates:

```
❌ Non-idempotent:
   "Transfer $100 from A to B" processed twice → $200 transferred. Bug.

✅ Idempotent:
   "Set balance of A to $900, B to $1100" (idempotent: same result if repeated)
   Or: Include a unique transaction_id. Check if already processed. Skip if so.
```

**The idempotency key pattern:**
```
1. Client generates a unique request_id (UUID) for each operation.
2. Server checks: "Have I seen request_id abc-123 before?"
3. Yes → return cached result.
4. No → process, store (request_id → result), return result.

Used by: Stripe, Twilio, all payment APIs.
```

---

### The Outbox Pattern — Reliable Event Publishing

**Problem**: You save to the DB and send a message to Kafka. What if Kafka is down?
```
❌ Write to DB, then publish to Kafka:
   DB write succeeds.
   Kafka publish fails.
   → Event lost. Other services never notified.

❌ Publish to Kafka, then write to DB:
   Kafka publish succeeds.
   DB write fails.
   → Event published but no data to back it up.
```

**Outbox Pattern (correct solution):**
```
1. In the SAME database transaction:
   - Write your business data (e.g., create order)
   - Write event to an "outbox" table in same DB

2. Outbox processor (separate process):
   - Reads unprocessed events from outbox table
   - Publishes to Kafka
   - Marks as processed

→ If Kafka is down: event stays in outbox until Kafka recovers.
→ Outbox + business data are in same transaction: atomicity guaranteed.
```

---

## SD Session 7: API Design ✅

### REST vs GraphQL vs gRPC

**REST** — The default choice:
```
Resources map to URLs. HTTP verbs have semantic meaning.

GET    /users/123       → Fetch user 123
POST   /users           → Create user
PUT    /users/123       → Replace user 123
PATCH  /users/123       → Update fields of user 123
DELETE /users/123       → Delete user 123

Pros: Universal, cacheable, stateless, human-readable, browser support.
Cons: Over-fetching (getting all user fields when you need 2),
      Under-fetching (need 5 API calls to build one page).
```

**GraphQL** — For complex, client-driven data fetching:
```
POST /graphql
Body:
  query {
    user(id: "123") {
      name
      email
      posts(last: 3) { title }  ← Exactly what the client needs
    }
  }

Pros: Client specifies exactly what data it needs (solves over/under-fetching).
      Single endpoint. Self-documenting (introspection).
Cons: Complex caching (queries are dynamic), N+1 query problem,
      Over-permission risk (clients can query anything schema exposes).
```

**gRPC** — For internal microservice communication:
```
Uses Protocol Buffers (binary encoding, ~5-10x smaller than JSON).
HTTP/2 (multiplexing, streaming, binary framing).
Strong typing via .proto contract files.

Pros: Extremely fast, streaming support (server-streaming, bi-directional),
      Generated client/server code in any language.
Cons: Not human-readable, poor browser support, requires proto tooling.

Use: Internal service-to-service calls where performance matters.
     Don't use for public-facing APIs (use REST/GraphQL for that).
```

---

### Pagination — You Must Know All Three

**Offset Pagination** (simple, use for admin UIs):
```
GET /posts?page=3&limit=20  →  OFFSET 40 LIMIT 20

Pros: Jump to any page, easy to implement.
Cons: Slow for large offsets (DB must scan and discard 40 rows).
      Inconsistent (new item inserted shifts all pages).
```

**Cursor Pagination** (production standard):
```
GET /posts?cursor=eyJpZCI6MTIzfQ==&limit=20
  → WHERE id > 123 ORDER BY id LIMIT 20

Pros: O(log n) regardless of page depth (uses index).
      Stable (inserts/deletes don't shift pages).
Cons: Can't jump to page 5. Cursor must be opaque to clients.
```

**Keyset Pagination** (same idea as cursor, explicit):
```
GET /posts?after_id=123&limit=20
```

**When to use each:**
- User-facing feeds/timelines → **Cursor**
- Admin tables where jumping to page N matters → **Offset**
- High-performance at scale → **Cursor / Keyset**

---

### Rate Limiting — Four Algorithms

**1. Fixed Window Counter**:
```
Bucket per time window (e.g., per minute). Request increments counter.
If counter > limit, reject.

Problem: "Burst at boundary" — 100 requests at 0:59 + 100 at 1:01 = 200 req in 2 sec,
         but both windows report 100.
```

**2. Sliding Window Log**:
```
Store timestamp of every request. On each request, count timestamps in last N seconds.
If count > limit, reject.

Pro: Accurate.
Con: Memory-intensive (store all timestamps).
```

**3. Sliding Window Counter** (best balance):
```
Combine two fixed windows with weighted interpolation:
current_count = prev_window_count × (overlap_%) + current_window_count

Pro: Memory-efficient. More accurate than fixed window.
Con: Approximation (not perfect).
```

**4. Token Bucket** (most flexible):
```
Bucket holds up to `capacity` tokens.
Tokens added at rate `refill_rate` per second.
Each request consumes 1 token.
Reject if bucket empty.

Pro: Allows short bursts (up to bucket capacity). Simple to implement.
Con: Slightly bursty behavior.

Used by: AWS API Gateway, Stripe, most cloud rate limiters.
```

---

## SD Session 8: Resilience Patterns ✅

### Why Resilience Patterns Exist

In a distributed system with 10 services each at 99.9% uptime:
```
Combined uptime = 0.999^10 = 99.0%
→ System is down ~1% of the time = ~87 hours/year
```
Failures are inevitable. Your system must **fail gracefully**.

---

### Circuit Breaker

Named after electrical circuit breakers that prevent overload damage.

```
Three states:

CLOSED (normal):
  Requests pass through. Failure counter tracked.
  If failures > threshold in window → trip to OPEN.

OPEN (tripped):
  All requests fail immediately (no downstream call).
  After timeout, move to HALF-OPEN.

HALF-OPEN (testing):
  Allow one test request through.
  If success → back to CLOSED.
  If failure → back to OPEN.

Why?
  Without circuit breaker: slow downstream service blocks your threads,
  causing your service to also become slow/unresponsive (cascading failure).
  Circuit breaker stops the cascade.
```

---

### Retry with Exponential Backoff + Jitter

```
Naive retry (bad):
  Fail → wait 1s → retry → fail → wait 1s → retry
  If 1000 clients all retry at the same time → thundering herd on recovering server.

Exponential backoff (better):
  attempt 1: wait 1s
  attempt 2: wait 2s
  attempt 3: wait 4s
  attempt 4: wait 8s

Exponential backoff + jitter (best):
  wait = min(base × 2^attempt + random(0, 1000ms), max_wait)
  Random jitter spreads retries across time → no thundering herd.
```

**Only retry idempotent operations** — retrying a non-idempotent operation (payment charge) can cause duplicates.

---

### Bulkhead Pattern

Named after ship bulkheads that contain flooding to one compartment.

```
Without bulkhead:
  One slow dependency uses all 200 threads → entire application hangs.

With bulkhead:
  Thread pool A (50 threads) → Dependency X
  Thread pool B (50 threads) → Dependency Y
  Thread pool C (100 threads) → Core functionality

  If Dependency X hangs → Thread pool A exhausted.
  Thread pools B and C unaffected. Core functionality survives.
```

---

### Timeout — Always Set Them

```
❌ No timeout:
   Service A calls Service B.
   Service B hangs (slow DB query).
   Service A's thread blocks forever.
   New requests pile up → Service A's thread pool exhausted → Service A down.

✅ With timeout (e.g., 500ms):
   After 500ms, Service A gives up and returns an error.
   Thread is freed immediately.
   Service A stays healthy.

Rule: Always set a timeout on every outbound network call.
      Timeout value = P99 latency of the dependency + buffer.
```

---

### Saga Pattern — Distributed Transactions Without 2PC

**Problem**: You have a business operation spanning multiple services (place order → reserve inventory → charge payment). How do you keep consistency without a distributed transaction?

**Two-Phase Commit (2PC)** is the classic answer, but it has problems:
- Requires all participants to hold locks during the transaction
- If coordinator crashes, participants are blocked forever
- Poor availability

**Saga pattern** breaks the operation into a sequence of local transactions, each with a **compensating transaction** if something fails:

```
Choreography Saga (event-driven):

1. Order Service: Create order (status=PENDING) → publish "OrderCreated"
2. Inventory Service: Reserve items → publish "InventoryReserved"
3. Payment Service: Charge card → publish "PaymentProcessed"
4. Order Service: Update order (status=CONFIRMED)

If Payment fails:
3. Payment Service: publish "PaymentFailed"
2. Inventory Service: Release reservation (compensating tx)
1. Order Service: Cancel order (compensating tx)

Pro: Loose coupling, no coordinator needed.
Con: Complex to reason about, eventual consistency, hard to debug.
```

---

## SD Session 9: Case Studies — Putting It All Together ✅

### Case Study 1: Design a URL Shortener

**Requirements:**
- 100M URLs shortened/day (1,200 writes/sec)
- 10B redirects/day (120,000 reads/sec)
- URLs expire after 1 year
- Custom aliases optional

**The interesting problem: Short code generation**

```
Option A: Hash (MD5/SHA1, take first 7 chars)
  MD5("https://example.com/...") → "a1b2c3d"
  Pro: Deterministic (same URL → same code), stateless.
  Con: Collisions possible. Need collision resolution. Hard to guarantee uniqueness.

Option B: Base62 Counter (recommended)
  Global counter (auto-increment): 1, 2, 3, ...
  Convert to base62: 1 → "1", 62 → "10", 3844 → "100"
  7 chars of base62 = 62^7 = 3.5 trillion unique URLs.
  Pro: No collisions, simple.
  Con: Sequential IDs leak business volume. Single counter = bottleneck.

Option C: Distributed Counter (Twitter Snowflake)
  Each app server has a unique machine_id.
  ID = timestamp(41 bits) + machine_id(10 bits) + sequence(12 bits)
  Pro: No central bottleneck, globally unique, k-sortable.
```

**Architecture:**
```
Write Path:
  POST /shorten → API Server → Generate code → Write to DB + Cache
  
Read Path (120K req/sec):
  GET /{code} → API Server → Redis cache hit → HTTP 301 Redirect
                                ↓ miss
                             DB lookup → cache → HTTP 301 Redirect

Cache: Redis. TTL = 24h. 120K reads/sec is easily handled by Redis.
DB: PostgreSQL. With index on short_code, 1,200 writes/sec is trivial.
```

---

### Case Study 2: Design a News Feed (Twitter/Instagram)

**The core problem: When user A posts, their 10M followers need to see it.**

**Two approaches:**

**Pull (Fan-Out on Read):**
```
When user loads feed:
1. Fetch IDs of all accounts they follow
2. Query posts from each account
3. Merge and sort by time

Pro: No storage overhead, always fresh.
Con: Extremely slow for users following 10,000 accounts (10K DB queries per load).
```

**Push (Fan-Out on Write):**
```
When user posts a tweet:
1. Find all followers
2. Write tweet ID into each follower's feed (pre-computed timeline)
3. User loads feed: single read from their pre-computed timeline

Pro: Feed reads are O(1) — instant.
Con: "Celebrity problem" — Kylie Jenner posts → write to 200M feed lists.
     Writing to 200M records takes minutes.
```

**Hybrid (used by Twitter):**
```
Regular users (< 1M followers): Fan-out on write (fast reads).
Celebrities (> 1M followers): Fan-out on read (skip writing to 200M feeds).

When loading feed:
1. Read pre-computed timeline (fast path)
2. Query posts from any celebrities you follow (slow path, few of them)
3. Merge

This is how Twitter's actual architecture worked.
```

---

### Case Study 3: Design a Distributed Rate Limiter

**Requirements:**
- Rate limit per user: 1000 requests/hour
- Works across multiple API servers (centralized enforcement)
- Low latency overhead (<5ms added to request)

**Solution: Redis + Token Bucket**

```lua
-- Redis Lua script (atomic execution):
local key = "ratelimit:" .. user_id
local limit = 1000
local window = 3600  -- 1 hour in seconds

local current = redis.call("INCR", key)
if current == 1 then
  redis.call("EXPIRE", key, window)
end
if current > limit then
  return 0  -- rejected
end
return 1  -- allowed
```

**Why a Lua script?** Redis executes Lua scripts atomically — no race condition between INCR and EXPIRE.

**For sliding window rate limiting in Redis:**
```
Use a Sorted Set (ZSET):
- Score = request timestamp
- Member = unique request ID
- On each request:
  1. Remove old entries (ZREMRANGEBYSCORE key 0 (now - window))
  2. Count remaining (ZCARD key)
  3. If count < limit: add new entry (ZADD key timestamp request_id)
  4. Else: reject
```

---

## SD Session 10: Advanced Deep Dives ✅

### Consensus Algorithms: Raft

**Problem**: How do distributed nodes agree on a value when some nodes can fail?

Raft is the easier-to-understand alternative to Paxos. It underpins etcd, Consul, CockroachDB, and TiDB.

**Key roles:**
```
Leader:    The one server that handles all writes. There is exactly ONE.
Follower:  Passive, replicate the leader's log.
Candidate: Running for election (temporary state).
```

**Leader election:**
```
1. All nodes start as followers with a random election timeout (150–300ms).
2. If no heartbeat from leader before timeout → become candidate.
3. Candidate increments its term, votes for itself, sends RequestVote to peers.
4. Majority votes → becomes leader.
5. Leader sends heartbeats to reset all followers' timeouts.

Safety: Two nodes can't both win an election in the same term (majority quorum).
```

**Log replication:**
```
1. Client sends write to leader.
2. Leader appends to its log (not committed yet).
3. Leader sends AppendEntries to followers.
4. Once majority (quorum) acknowledge → leader commits.
5. Leader notifies followers to commit. Returns success to client.
```

---

### Distributed Caching at Scale: Consistent Hashing with Virtual Nodes

In production Redis clusters or distributed caches, you need to distribute keys across nodes without centralized coordination:

```
Naive sharding: shard = hash(key) % 3
  Add a 4th node → hash(key) % 4 changes for 75% of keys.
  Cache miss storm. All traffic hits DB simultaneously.

Consistent hashing: Add 4th node → only ~25% of keys move.
  Virtual nodes ensure even distribution even with heterogeneous hardware.
```

---

### The Interview Cheat Sheet

Keep this mental model for EVERY system design interview:

```
1. Clarify (5 min)
   - Users? Scale? Read/write ratio?
   - Strong consistency or eventual OK?
   - Global or single-region? SLA?

2. Estimate (5 min)
   - QPS, storage, bandwidth
   - Derive whether you need caching, sharding, etc.

3. API + Data Model (5 min)
   - What are the endpoints?
   - What does the schema look like?

4. High-level (15 min)
   - Client → LB → Services → Cache → DB
   - Write path and read path separately

5. Deep-dive trade-offs (20 min)
   Every decision: Pro? Con? Why this over alternatives?

Key trade-offs to discuss:
  - Consistency vs Availability (CAP)
  - Latency vs Throughput
  - Fan-out on write vs Fan-out on read
  - Normalized vs Denormalized data
  - Strong consistency vs Eventual consistency
  - Push vs Pull
  - Sync vs Async processing
```

```

---
---

# ═══════════════════════════════════════════════════════════
# ORDER MANAGEMENT DOMAIN — Complete Knowledge Base
# ═══════════════════════════════════════════════════════════

> You are interviewing for teams working in **Order Management**, **Order Promises**,
> and **Order Allocation**. This section teaches you the domain from zero —
> what these systems are, why they exist, how they connect, and what problems they solve.
> By the end, you will speak fluently about this domain in your interview.

---

## Domain Overview: The Big Picture

When a customer clicks **"Buy Now"** on Amazon, what actually happens?
Most people think: "The order is saved, then it ships." The reality is a deeply
coordinated chain of 10+ systems across multiple companies and hundreds of decisions.

Here is the **complete order journey** — this is the domain you are entering:

```
Customer clicks "Buy"
       │
       ▼
① ORDER CAPTURE ──────────────────────── "I want this"
   Order Management System (OMS)
   - Creates order record
   - Validates cart, pricing, discounts
   - Fraud check
       │
       ▼
② ORDER PROMISE ──────────────────────── "We commit to deliver by Tuesday"
   Available-to-Promise (ATP) System
   - Checks real inventory across all warehouses
   - Calculates earliest delivery date
   - Reserves the stock so nobody else gets it
       │
       ▼
③ PAYMENT ─────────────────────────────── "Money changes hands"
   Payment Service (NOT your domain, but you integrate with it)
   - Authorizes and captures the payment
   - On failure: releases inventory reservation
       │
       ▼
④ ORDER ALLOCATION ────────────────────── "Which warehouse sends this?"
   Allocation Engine
   - Decides which Fulfillment Center (FC) will ship the order
   - Balances: delivery speed, cost, inventory levels
   - Can split across multiple FCs
       │
       ▼
⑤ WAREHOUSE FULFILLMENT ───────────────── "Workers pick items from shelves"
   Warehouse Management System (WMS) — NOT your domain
   - Pick list generated for warehouse workers
   - Items picked, verified, packed
   - Shipping label printed
       │
       ▼
⑥ SHIPPING / LAST MILE ────────────────── "Package travels to customer"
   Transportation Management System (TMS)
   - Carrier selected (FedEx, UPS, DHL, own fleet)
   - Tracking number generated
   - Real-time tracking events
       │
       ▼
⑦ DELIVERY ──────────────────────────────
   - Delivered to customer
   - Confirmation event triggers post-delivery flows
       │
       ▼
⑧ POST-DELIVERY ────────────────────────── "What if customer is unhappy?"
   Returns Management / Reverse Logistics
   - Customer initiates return
   - Return shipped back to warehouse
   - Refund processed
   - Inventory restocked
```

This entire chain is called the **Order-to-Cash (O2C)** process.
Your target teams own stages ①, ②, and ④.

---

## Domain 1: Order Management System (OMS)

### What Problem Does an OMS Solve?

Imagine a large retailer (like Zalando, ASOS, or a major bank's e-commerce platform)
that sells across:
- Their own website
- Mobile app
- Physical stores
- Amazon marketplace
- Instagram Shopping

Each channel has its own technical format for an "order." Without an OMS, you'd need to build order processing logic 5 times, and there would be no single source of truth for "what did this customer order?"

**The OMS is the single system of record for all orders, regardless of channel.**

### What Does an OMS Actually Do?

```
1. ORDER INGESTION
   Receives orders from all channels.
   Normalizes them into a standard internal format.
   "A Shopify order and a POS order look different externally,
    but the OMS stores them the same way internally."

2. ORDER VALIDATION
   - Is the product still available?
   - Is the price still valid? (Price changed during checkout?)
   - Is the address valid? (Address validation API)
   - Does the customer exist? (Fraud check)
   - Are there any promotions to apply?

3. ORDER LIFECYCLE MANAGEMENT
   Tracks every order through its states:
   PENDING → CONFIRMED → ALLOCATED → PICKED → SHIPPED → DELIVERED
   (with side paths for CANCELLED, FAILED, RETURNED)

4. ORDER VISIBILITY
   - Customer portal: "Where is my order?"
   - CS portal: "This customer called, what happened to order #12345?"
   - Operations dashboard: "How many orders are stuck in ALLOCATED state?"

5. INTEGRATIONS HUB
   OMS is the central hub that integrates with:
   - Inventory system (check availability)
   - Payment service (capture money)
   - ATP system (get delivery promise)
   - Allocation engine (decide which warehouse)
   - WMS (tell warehouse to pick-and-pack)
   - Carrier/TMS (get tracking number)
   - Notification service (update customer)
```

### The Order Lifecycle — A Finite State Machine

This is the most important concept in OMS. An order is always in **exactly one state**,
and can only move to specific next states. This prevents bugs like "shipped before paid."

```
                         ┌──────────────────────────────────────────────────┐
                         │                                                  │
PENDING ──payment ok──► CONFIRMED ──allocation ok──► ALLOCATED ──WMS ok──► IN_FULFILLMENT
   │          │               │              │              │
   │        payment         cancel         cancel         cancel
   │        failed            │              │              │
   ▼           │              ▼              ▼              ▼
PAYMENT_FAILED │         CANCELLED       CANCELLED     CANCELLED (complex —
               ▼                                        need to recall from WMS)
           PAYMENT_FAILED
                                                         │
                                          IN_FULFILLMENT ──► SHIPPED ──► DELIVERED
                                                                              │
                                                                    customer return
                                                                              │
                                                                              ▼
                                                                    RETURN_REQUESTED
                                                                              │
                                                                    return received
                                                                              │
                                                                              ▼
                                                                          RETURNED
```

**Why states matter in engineering:**
- You can never make a database `UPDATE orders SET status = 'SHIPPED'` directly.
  The system must validate the transition: `current_status == 'IN_FULFILLMENT'`.
- Every transition is an **event** that gets logged: who changed it, when, why.
- This makes debugging trivial: "Show me every event for order #12345."

### Key OMS Concepts You Must Know

**Order vs Order Line Item:**
```
Order: The top-level container.
  - One order per checkout session
  - Has one billing address, one payment method
  - Has one customer
  - May have multiple shipments

Order Line Item (or Order Line): One product in the order.
  - "3 × Blue Nike Trainers Size 42"
  - Each line item has: SKU, quantity, unit price, line total
  - Each line item is allocated, shipped, and tracked independently
  - A single order might have line items shipped from different warehouses
```

**Idempotency in order creation:**
```
Problem: Customer clicks "Place Order" twice (network lag).
         Two identical orders created. Customer charged twice. Disaster.

Solution: Idempotency Key
  - Client generates a UUID when checkout page loads: checkout_session_id
  - Sent with every "create order" request
  - OMS checks: "Have I seen this checkout_session_id before?"
    - YES → return the existing order (don't create a new one)
    - NO  → create order, store checkout_session_id → order mapping

This is how all payment APIs (Stripe, Adyen) work too.
```

**Optimistic Locking — Preventing Race Conditions:**
```sql
-- When updating an order, always include the version number
UPDATE orders
SET status = 'CONFIRMED', version = version + 1
WHERE id = 'abc-123'
  AND version = 5;   -- Only succeeds if nobody else updated first

-- If 0 rows updated: someone else changed the order → retry or reject
```

---

## Domain 2: Order Promises — Available to Promise (ATP)

### What Problem Does ATP Solve?

When you shop on Amazon and it says **"Delivery by Tuesday, Jan 7"** — that is a promise.
Amazon is committing to deliver that item to you by that specific date.

This promise is not random or made up. Amazon has a sophisticated system that calculates:
- Do I have this item in stock? Where?
- Which warehouse can fulfill it and get it to you in time?
- How long does shipping take to your zip code?
- Is there enough buffer so that the promise is reliable?

**Without ATP:**
- You'd always show generic promises ("Ships in 3-5 business days")
- You'd lose sales to competitors showing precise dates
- You'd oversell and then have to apologize to customers ("Sorry, we don't actually have it")

**With ATP:**
- Precise, reliable promises → higher conversion rate
- Real inventory commitment → no overselling
- Competitive advantage: showing "Arrive by Tomorrow, 11 PM" beats "3-5 days"

### Core ATP Concepts

**On-Hand Inventory:**
```
The physical quantity of a SKU currently sitting in a warehouse.
"Right now, in FC-London, there are 500 units of SKU ABC123."
This is the starting point for all ATP calculations.
```

**Soft Reservation vs Hard Commitment:**
```
Soft Reservation (during checkout):
  Customer adds to cart. We "soft-reserve" 1 unit.
  The unit is no longer shown as available to other customers.
  But: if customer abandons checkout, the reservation expires (e.g., after 15 min).
  Stored in Redis with a TTL.

Hard Commitment (after payment):
  Customer completes payment. The reservation becomes a commitment.
  The unit is now definitively allocated to this order.
  Stored permanently in the database.

Why the distinction?
  If we only reserved after payment, two customers could both see
  "1 left in stock" and both try to buy. One would fail after payment.
  Bad experience. The soft reservation prevents this.
```

**Lead Time:**
```
The time from "order placed" to "delivered to customer."
Components:
  Fulfillment Lead Time: How long does the warehouse need to pick, pack, and hand to carrier?
                         (Usually same-day or next-day for standard items)
  Transit Time: How long does shipping take from the warehouse to the customer's zip code?
                (Depends on carrier, service level, distance)

Total Lead Time = Fulfillment Time + Transit Time

Example:
  FC-London → Customer in Manchester
  Fulfillment: 4 hours (same-day FC)
  Transit: 1 day (standard next-day delivery)
  Promise: "Order by 2 PM today → Delivered tomorrow"
```

**Available-to-Promise (ATP) Formula:**
```
ATP for a time period T =
    On-Hand Inventory
  + Scheduled Inbound (POs arriving by T)
  - Committed Demand (orders already placed for delivery by T)
  - Safety Stock (buffer — don't promise the very last unit)

If ATP(T) > 0 → Can promise delivery at or before T
If ATP(T) = 0 → Earliest promise date is when next replenishment arrives
If no replenishment → Out of stock, cannot promise
```

**Capable to Promise (CTP):**
```
ATP answers: "Do I have inventory?"
CTP answers: "Can I MAKE it or PROCURE it in time?"

CTP is used when:
- Made-to-order products (custom furniture, bespoke clothing)
- Procurement-based businesses (B2B with supplier lead times)
- ATP = 0 but a purchase order can be placed and arrive in time

CTP is more complex — it considers supplier lead times, production capacity, etc.
For your interview context (e-commerce/retail OMS), ATP is primary.
```

**Why Accuracy Matters:**
```
If you promise too aggressively (promise dates you can't hit):
  → Broken promises → Customer dissatisfaction → Refunds → Returns
  → Trust loss: customer doesn't buy from you again

If you promise too conservatively (add 3 extra days to be safe):
  → Customer sees "5-7 days" while competitor shows "2 days"
  → You lose the sale

The target: Promise accuracy > 99.5% (less than 0.5% broken promises)
This requires:
  - Accurate real-time inventory data
  - Realistic transit time models (carrier SLA data)
  - Safety stock buffer calibration
  - Monitoring and alerting on broken promises
```

### Promise Scenarios You Must Know

**Scenario 1: Simple In-Stock Item**
```
Customer in Berlin orders a TV → FC-Frankfurt has 50 in stock
ATP = 50 - 0 (no prior commitments) - 5 (safety stock) = 45 → Can promise
Transit time Berlin ← Frankfurt: 1 day
Promise: "Delivery Tomorrow"
```

**Scenario 2: Multi-FC Selection**
```
Customer in Edinburgh orders a laptop
  FC-London:    2 in stock, 1.5 day transit → Promise: Day After Tomorrow
  FC-Glasgow:   0 in stock, 0.5 day transit → Cannot promise from here
  FC-Birmingham: 5 in stock, 2 day transit → Promise: 2 days

→ Pick FC-London (fastest that has stock)
```

**Scenario 3: Split Inventory (bundles)**
```
Customer orders: 1× keyboard + 1× mouse + 1× monitor (as separate items)
  Keyboard: FC-Amsterdam has stock, transit 1 day
  Mouse: FC-Amsterdam has stock, transit 1 day
  Monitor: Only FC-Frankfurt has stock, transit 2 days

Options:
  A. Promise entire order for Day 2 (slowest item) — single shipment
  B. Split: keyboard+mouse arrive Day 1, monitor Day 2 — two shipments
  
Customer preference + business rules decide. Usually configurable.
```

**Scenario 4: Backorder / Out of Stock**
```
Customer wants SKU ABC, currently out of stock everywhere.
Next PO (Purchase Order) from supplier arrives in 7 days.
Transit time: 1 day.

ATP answer: "Available in 8 days"
OMS can:
  A. Show "Out of Stock" — no promise
  B. Accept pre-order with "Estimated delivery in 8 days"
  C. Offer "Notify me when available" (email when back in stock)

Business decision, but the system must support all three.
```

---

## Domain 3: Order Allocation

### What Problem Does Allocation Solve?

You have 1,000 orders to fulfill today. You have 8 warehouses across Europe.
Each warehouse has different inventory levels. Each warehouse is at a different
distance from each customer. Shipping from different warehouses has different costs.

**Allocation decides: which warehouse ships each order.**

This sounds simple but is actually a **constrained optimization problem**.

### Key Constraints in Allocation

```
Hard Constraints (cannot be violated):
  ✓ FC must have sufficient stock of all items in the order
  ✓ FC must be capable of shipping the item (hazmat restrictions, temperature control)
  ✓ FC must be open and operational (no weather closures)
  ✓ Carrier cutoff time not passed (can't ship same-day after 3 PM)

Soft Constraints (optimize for):
  ✓ Minimize transit time to customer (faster = better promise)
  ✓ Minimize shipping cost (business efficiency)
  ✓ Avoid overloading one FC (distribute work evenly)
  ✓ Minimize number of shipments (split shipments cost more)
  ✓ Match promised delivery date
```

### Allocation Strategies

**Strategy 1: Proximity-Based (simple)**
```
Always ship from the closest warehouse to the customer.
Pro: Fastest delivery, lowest shipping cost.
Con: One warehouse gets all orders for a region and runs out of stock
     while others are full. Imbalanced.
```

**Strategy 2: Inventory-Based**
```
Ship from the warehouse with the most stock (to avoid stockouts).
Pro: Evenly depletes inventory across all FCs.
Con: May ship from far away, breaking the promise date.
```

**Strategy 3: Scored Allocation (what production systems use)**
```
Score every candidate FC for each order:
  score = (delivery_speed_weight × speed_score)
        + (cost_weight × cost_score)
        + (inventory_weight × inventory_score)
        + (capacity_weight × capacity_score)

Pick the highest scoring FC.
The weights are business configuration — tuneable by ops team.
```

**Strategy 4: Zone-Based Allocation**
```
Pre-assign geographic zones to FCs.
  Zone A (UK): → FC-Manchester
  Zone B (Germany, Austria): → FC-Frankfurt
  Zone C (France, Benelux): → FC-Paris

Simple, predictable, low latency (no algorithm, just a lookup).
Used when simplicity is more important than optimization.
Downside: No fallback when the assigned FC is out of stock.
```

### Split Fulfillment

**What is it?**
```
A single order fulfilled from multiple warehouses.
Customer orders: blue jacket + white shirt + black trousers
  Blue jacket: only in FC-North (stock = 2)
  White shirt: only in FC-South (stock = 5)
  Black trousers: available in both, better from FC-South

Allocation decides:
  Shipment 1 from FC-North: blue jacket
  Shipment 2 from FC-South: white shirt + black trousers
```

**Pros and Cons:**
```
Pros:
  - Can fulfill the order at all (no single FC has everything)
  - Can deliver partial order faster (don't wait for slowest item)
  - Better inventory utilization across FCs

Cons:
  - Higher shipping cost (two packages, two labels)
  - More complex tracking (customer gets two tracking numbers)
  - Higher chance of partial delivery issues
  - Customer confusion ("where is the rest of my order?")
```

**Business Rule: Max Splits**
```
Most businesses cap splits at 2 or 3.
If an order requires 4+ FCs to fulfill: flag it for manual review
or trigger a backorder/pre-order flow.
```

### Reallocation

**What happens when allocation needs to change after it's been made?**

```
Scenarios:
  1. Warehouse discovers the item is missing or damaged during pick
  2. Warehouse is shut down (fire, flood, system outage)
  3. Inventory reconciliation reveals the FC never actually had the stock
  4. Customer requests change of delivery address (different FC is now better)

Reallocation constraints:
  - Can only reallocate if order is not yet picked (status < PICK_STARTED)
  - Once picking has begun, the physical goods are in motion — can't recall
  - After reallocation: send new work order to new FC, cancel at old FC
  - Customer notification: "Your order is now shipping from a different location"
```

### The Allocation Decision Log

Every allocation decision must be logged. Why?

```
1. Debugging: "Why was this order sent to FC-Berlin instead of FC-Hamburg?"
2. Analytics: "Are we making good allocation decisions? Is FC-Berlin getting overloaded?"
3. ML training: If you're using ML for allocation, you need training data
4. Compliance: Some industries require a full audit trail of all decisions
5. Post-mortems: "Why did we break 500 promises last week?"

allocation_decisions table:
  order_id, order_item_id, fc_id, algorithm_version,
  decision_score, alternatives_considered, reason_code, decided_at
```

---

## The Three Systems — How They Work Together

```
                  CUSTOMER JOURNEY
                        │
         ┌──────────────▼──────────────┐
         │                             │
         │  1. "Do you have it?"        │ ← ATP Query (before order placed)
         │     "Can you deliver by X?"  │   "Show promise on product page"
         │                             │
         └──────────────┬──────────────┘
                        │ Yes, can promise
                        ▼
         ┌──────────────────────────────┐
         │  2. "I want to buy it"        │ ← OMS: Create Order
         │                              │   Status: PENDING
         │     Soft reserve inventory   │ ← ATP: Reserve inventory
         │     Status: CONFIRMED        │   OMS: Status → CONFIRMED
         └──────────────┬───────────────┘
                        │ Payment OK
                        ▼
         ┌──────────────────────────────┐
         │  3. "Which warehouse ships?" │ ← Allocation Engine
         │                              │   Scores all FCs
         │     Hard commit inventory    │ ← ATP: Convert reservation to commitment
         │     Status: ALLOCATED        │   OMS: Status → ALLOCATED
         └──────────────┬───────────────┘
                        │
                        ▼
         ┌──────────────────────────────┐
         │  4. Warehouse fulfills        │ ← WMS (not your domain)
         │     Status: SHIPPED           │   OMS: Status → SHIPPED
         └──────────────┬───────────────┘
                        │
                        ▼
         ┌──────────────────────────────┐
         │  5. Customer receives         │ ← OMS: Status → DELIVERED
         │     Promise: Kept or Broken?  │ ← ATP: Record promise outcome
         └──────────────────────────────┘
```

**Key integration points:**

| System | Owns | Calls |
|---|---|---|
| OMS | Order lifecycle, status, events | ATP, Allocation, Payment, WMS, Notifications |
| ATP | Inventory positions, reservations, promises | Inventory DB, OMS (to confirm/release) |
| Allocation | FC selection decisions, split logic | ATP (to reserve), WMS (to assign), OMS (to update) |

---

## The Ecosystem of Systems Around OMS

You will hear these terms in your interview. Know what each one is.

### ERP (Enterprise Resource Planning)
```
The master system for the whole business: finance, HR, procurement, supply chain.
Examples: SAP, Oracle ERP, Microsoft Dynamics.

Relationship to OMS:
  ERP is the "source of truth" for:
    - Product catalog and pricing
    - Supplier purchase orders (inbound supply for ATP)
    - Financial reconciliation (revenue from orders)

  OMS sends to ERP: "Order #12345 was paid — recognize this revenue"
  ERP sends to OMS: "PO #678 arriving at FC-London tomorrow — update ATP"
```

### WMS (Warehouse Management System)
```
The system that runs inside a warehouse/fulfillment center.
Controls: receiving goods, storing them, picking for orders, packing, shipping.
Examples: Manhattan Associates WMS, Blue Yonder, HighJump.

Relationship to OMS/Allocation:
  Allocation says: "FC-London should fulfill order #12345"
  → OMS sends work order to FC-London's WMS
  WMS says: "Picked and packed, here's the tracking number"
  → WMS sends fulfillment confirmation to OMS
  OMS updates status: ALLOCATED → SHIPPED
```

### TMS (Transportation Management System)
```
Manages the shipping and carrier relationships.
Selects which carrier (FedEx, UPS, DHL, own delivery fleet) for each shipment.
Examples: Oracle TMS, MercuryGate, project44.

Relationship to OMS:
  OMS: "I need to ship this package from FC-London to a Berlin address"
  TMS: Rates all available carriers, picks cheapest/fastest
  TMS: Generates shipping label, returns tracking number to OMS
  TMS: Receives carrier tracking events, relays to OMS
```

### PIM (Product Information Management)
```
The source of truth for product data:
  - Product name, description, images
  - SKU codes, barcodes (EAN, UPC)
  - Dimensions and weight (needed for shipping cost calculation)
  - Category, attributes (size, color, material)
Examples: Akeneo, Contentful, inRiver.

Relationship to OMS/ATP:
  OMS reads product data from PIM (product name in the order receipt)
  ATP reads weight/dimensions from PIM (affects shipping time calculation)
  Allocation reads special attributes from PIM (hazmat? temperature-sensitive? oversized?)
```

### 3PL (Third-Party Logistics)
```
An external company that handles warehousing and fulfillment FOR you.
Instead of owning your warehouse, you pay a 3PL to store and ship your products.
Examples: Flexport, Ryder, XPO Logistics, Amazon FBA.

Relationship to OMS:
  Exactly the same as an internal FC — OMS sends allocation decisions to 3PL,
  3PL sends back fulfillment confirmations.
  The difference: 3PL has its own systems (may require EDI, API, or file-based integration)
```

### MDM (Master Data Management)
```
Ensures that "the same thing" has one consistent identity across all systems.
A customer might be stored as:
  - "John Smith" in OMS
  - "J. Smith" in ERP
  - "SMITH, JOHN" in TMS

MDM reconciles these and assigns a single Master Customer ID.
Same for products: one SKU across all systems.
```

---

## Key Domain Metrics — What the Business Cares About

When you work in this domain, your engineering decisions are evaluated against these KPIs:

```
Order Fulfillment:
  ✦ Order Placement Rate: Orders created per second (system capacity)
  ✦ Order Error Rate: % of orders that fail validation or processing
  ✦ Order Cycle Time: Average time from order placed to order shipped

Promise Accuracy (ATP Team's North Star KPI):
  ✦ Promise Hit Rate: % of orders delivered by the promised date
                      Target: > 99.5%
  ✦ Promise Accuracy by FC: Which FCs are reliable vs which are breaking promises?
  ✦ Promise Accuracy by SKU: Which products are over-promised?
  ✦ ATP Cache Hit Rate: % of promise queries served from cache

Allocation Quality (Allocation Team's KPIs):
  ✦ Single-FC Fulfillment Rate: % of orders fulfilled from one FC (higher = better)
  ✦ Split Rate: % of orders requiring multiple FCs
  ✦ Allocation Accuracy: Did the FC we allocated to actually have the stock?
  ✦ Reallocation Rate: % of orders reallocated after initial assignment
  ✦ Average Allocation Score: How optimal are our FC selections?

Inventory Health:
  ✦ Inventory Accuracy: Physical count vs system count (target: >99%)
  ✦ Oversell Rate: % of orders where we promised but couldn't fulfill
  ✦ Stockout Rate: % of times a popular SKU showed "out of stock" unnecessarily
  ✦ Safety Stock Effectiveness: Are our buffers calibrated correctly?
```

---

## Common Domain Interview Questions — With Model Answers

**"What is the difference between OMS and ERP?"**
```
ERP is the company's master system — it owns financials, procurement, HR.
OMS is specialized for order processing — it focuses on the customer-facing
order lifecycle. They integrate: OMS sends fulfilled order revenue to ERP,
ERP sends incoming purchase orders to ATP for supply planning.
An ERP's order module is usually too slow and generic for high-volume e-commerce.
That's why dedicated OMS systems exist.
```

**"Why would you have a separate ATP system instead of just querying inventory?"**
```
Two reasons:
1. Performance: At checkout, you need <200ms response. If 50,000 users are on
   the product page simultaneously, you can't run complex inventory calculations
   for each one. ATP pre-computes the answers and serves from cache.

2. Accuracy: ATP doesn't just check on-hand inventory — it also considers
   soft reservations (items in other customers' carts), committed demand,
   inbound supply, and lead times. A raw inventory query would show 10 units
   as available, but ATP knows 8 are already soft-reserved.
```

**"What is a promise breach and how do you prevent it?"**
```
A promise breach is when we told a customer "delivery by Tuesday" but delivered Wednesday.
Causes:
  - ATP was based on inaccurate inventory data
  - Warehouse couldn't pick the item (shrinkage, damage)
  - Carrier delay (weather, volume surges)
  - We over-promised (used too aggressive lead times)

Prevention:
  - Safety stock buffer (don't promise the last unit)
  - Conservative lead time models (use P95 transit time, not average)
  - Real-time inventory reconciliation
  - Carrier SLA monitoring (if FedEx is late 5% of the time, add buffer)
  - Circuit breaker: if breach rate > threshold for an FC, stop promising from it

Recovery:
  - Proactive notification to customer: "Your order is delayed, new ETA is Thursday"
  - Compensation: discount voucher, free expedited shipping
  - Business escalation: if breach rate spikes, ops team investigates
```

**"How do you handle Black Friday — 100x normal order volume?"**
```
This is a system design question. My approach:

Pre-event preparation:
  - ATP: Pre-compute inventory positions and cache them aggressively
  - Allocation: Pre-assign inventory quotas per FC per product category
  - OMS: Scale horizontally (more API server instances, auto-scaling configured)
  - Database: Read replicas spun up, connection pools tuned

During the event:
  - Order queue: Instead of synchronous order creation, accept orders into a queue
    and process them (async validation, allocation). Customer gets "order received" immediately.
  - Rate limiting: Protect downstream services (payment, WMS) from being overwhelmed
  - Circuit breakers: If payment service is slow, queue the retry; don't fail the order
  - Shed load gracefully: If queue is full, show "system busy, try in 5 minutes"
    rather than crashing

Inventory:
  - Redis atomic decrements for inventory reservation (not DB locks)
  - Pre-allocated slots: "FC-East has 10,000 units reserved for this SKU today"
  - When quota exhausted: show "Sold Out" — no DB query needed
```

---

## Glossary — Every Term You Will Hear

```
SKU (Stock Keeping Unit)
  A unique identifier for a product variant.
  "Nike Air Max, Blue, Size 42" is one SKU.
  "Nike Air Max, Blue, Size 43" is a DIFFERENT SKU.
  All inventory tracking is at the SKU level.

UPC / EAN / GTIN
  Barcode standards. Each physical product has one.
  EAN-13 (European), UPC-A (US), GTIN (global standard).
  Used to identify products when receiving stock at a warehouse.

FC (Fulfillment Center)
  A warehouse optimized for e-commerce order fulfillment.
  Designed for high-speed pick-and-pack of individual items.
  Different from a distribution center (handles pallets for stores).

Pick, Pack, Ship
  The three stages inside a warehouse:
  Pick: Worker walks to shelf, takes items for an order
  Pack: Items placed in box with packing material, sealed
  Ship: Label printed, handed to carrier

Pick List
  The instruction sent to a warehouse worker:
  "Go to aisle 3, shelf B4, pick 2× SKU-ABC123"
  Generated by OMS/WMS when an order is allocated.

Carrier
  The company that transports the package: FedEx, UPS, DHL, Royal Mail.
  "Last Mile" carrier = the one that delivers to the customer's door.

Manifest
  The list of all packages handed to a carrier in one batch.
  "End-of-day manifest": all packages picked up by FedEx today.

PO (Purchase Order)
  A business buys goods from a supplier using a PO.
  PO #1234: "Order 500 units of SKU-ABC from Supplier X, arriving March 15."
  ATP uses inbound POs to calculate future inventory availability.

Shrinkage
  Inventory loss due to theft, damage, miscounts, or expiry.
  Typical retail shrinkage: 1-2% of inventory per year.
  ATP uses shrinkage rates to set safety stock levels.

Safety Stock
  A buffer quantity kept in inventory to absorb uncertainty.
  "Even though ATP says 10 units available, only promise 9 — keep 1 as safety stock."
  Protects against: demand spikes, supplier delays, inaccurate counts.

Reorder Point (ROP)
  The inventory level that triggers a new purchase order to the supplier.
  "When SKU-ABC drops below 50 units, order 200 more."
  Part of supply chain management, feeds into ATP's inbound schedule.

Backorder
  An order accepted when the item is currently out of stock.
  Customer agrees to wait. Fulfilled when new stock arrives.
  Different from cancellation: the order exists and is committed.

Pre-order
  An order placed before a product is released/available.
  Common in gaming, fashion, electronics.
  ATP must model this: "100 pre-orders committed, product arriving March 20."

Drop Shipping
  You take the order, but a supplier ships directly to the customer.
  You never touch the physical goods.
  Allocation variant: route to supplier's system instead of your FC.

EDI (Electronic Data Interchange)
  Old standard for B2B data exchange. Used by retailers and large suppliers.
  "Send your PO to us as an EDI 850 document."
  You will encounter EDI when integrating with legacy ERP/WMS systems.

SLA / SLO
  SLA: Service Level Agreement — a contract with the customer.
       "We guarantee 99.9% uptime. We guarantee delivery by promised date."
  SLO: Service Level Objective — an internal target.
       "We want 99.95% promise hit rate." (More ambitious than what we promise externally.)

On-Time In-Full (OTIF)
  Retail KPI: Was the delivery on time AND was it the full quantity ordered?
  A retailer receiving from a supplier might penalize if OTIF < 98%.
  In e-commerce: was the customer order delivered on the promised date with all items?

Carrier SLA
  The delivery time guaranteed by a carrier.
  "FedEx Express: Next Business Day"
  "FedEx Ground: 3-5 Business Days"
  ATP uses carrier SLA data to calculate promise dates.

Carrier Cut-off Time
  The deadline for handing packages to a carrier for same-day collection.
  "FedEx picks up from FC-London at 3 PM."
  Allocation must check: if it's 4 PM, this FC cannot ship same-day.

Returns / Reverse Logistics
  The process of a customer returning an item.
  Engineering challenges:
    - Generate return shipping label
    - Track return in transit
    - Inspect item upon receipt
    - Restock (if sellable) or dispose
    - Trigger refund

Refund
  Returning money to the customer.
  Partial refund: one item of a multi-item order returned.
  Payment service handles money, but OMS tracks the refund status.

3PL (Third-Party Logistics)
  External fulfillment provider. You don't own the warehouse — they do.
  You integrate with their API/EDI to send allocation decisions.

Cross-docking
  Goods arriving at a warehouse are immediately transferred to outbound
  shipments without being stored. Very efficient for fast-moving items.
  Reduces storage cost, speeds up fulfillment.
```


> Each case study follows the **5-Step Framework** from SD Session 1.
> Read the question. Pause. Then walk through it step by step — out loud.
> The ✦ markers show exactly what to say to impress the interviewer.

---

## Case Study 1: Design an Order Management System (OMS) 🛒

> **Question**: "Design an Order Management System that handles order placement,
> status tracking, and fulfillment across multiple channels (web, mobile, in-store)."

This is your **home domain** — you should be able to drive this conversation deeply.

---

### Step 1 — Clarify Requirements

**You say**: "Before I start designing, I have a few clarifying questions."

```
Functional:
  ✦ "What channels create orders? Web, mobile, POS, marketplace (Amazon)?"
  ✦ "What does an order lifecycle look like? Placed → Confirmed → Allocated
     → Shipped → Delivered? Are there cancellations and returns?"
  ✦ "Do we need to support split shipments? (One order, items from 3 warehouses)"
  ✦ "Is there a payment component, or is that a separate service we integrate with?"
  ✦ "Do we need real-time order tracking for customers?"

Non-functional:
  ✦ "What's the order volume? 100K orders/day? 10M?"
  ✦ "What's the peak multiplier? (Black Friday = 10x normal?)"
  ✦ "What's the SLA for order placement latency? 500ms p99?"
  ✦ "Is strong consistency required? (Can two customers buy the last item?)"
  ✦ "Global or single-region?"
```

**Assume the interviewer says**:
- 500K orders/day, 10x peak on sale days = 5M orders/day
- Web + mobile + marketplace channels
- Full lifecycle including cancellations and returns
- Payment is a separate service (just integrate)
- Strong consistency for inventory: overselling is NOT acceptable

---

### Step 2 — Estimate Scale

```
Normal:
  500K orders/day ÷ 86,400 = ~6 orders/sec (writes)
  Each order triggers: inventory check, payment, notification, fulfillment → ~30 events/sec

Peak (10x):
  60 orders/sec
  Status reads (customers checking): ~100x orders = 600 reads/sec normal, 6,000 peak

Storage per order:
  Order: ~2KB (header + line items + address)
  500K/day × 365 × 5 years × 2KB = ~1.8 TB/year
  → Feasible in a single RDBMS with partitioning, or NoSQL for scale

✦ "Given 60 writes/sec at peak, a single Postgres instance handles this easily.
   But the read pattern (customers polling status) at 6,000/sec benefits from caching."
```

---

### Step 3 — API & Data Model

**Core API:**
```
POST   /orders                    Create order
GET    /orders/{orderId}           Get order details + status
PATCH  /orders/{orderId}/cancel    Cancel order
POST   /orders/{orderId}/returns   Initiate return

POST   /orders/{orderId}/events    (Internal) Add lifecycle event
GET    /orders?customerId=X        Customer's order history
```

**Data Model:**
```sql
-- Orders table (the core)
orders (
  id           UUID         PRIMARY KEY,
  customer_id  BIGINT       NOT NULL,
  channel      VARCHAR(20)  NOT NULL,  -- WEB, MOBILE, POS, MARKETPLACE
  status       VARCHAR(30)  NOT NULL,  -- see FSM below
  total_amount DECIMAL(12,2),
  currency     CHAR(3),
  created_at   TIMESTAMP,
  updated_at   TIMESTAMP,
  version      INT          DEFAULT 1  -- for optimistic locking
)

-- Line items
order_items (
  id          UUID,
  order_id    UUID         REFERENCES orders(id),
  product_id  BIGINT,
  sku         VARCHAR(50),
  quantity    INT,
  unit_price  DECIMAL(10,2),
  warehouse_id BIGINT      -- set during allocation
)

-- Audit trail (append-only, never update)
order_events (
  id          UUID,
  order_id    UUID,
  event_type  VARCHAR(50),  -- ORDER_PLACED, PAYMENT_CONFIRMED, SHIPPED, etc.
  payload     JSONB,        -- flexible metadata per event type
  created_at  TIMESTAMP,
  created_by  VARCHAR(100)  -- user or service name
)
```

**Order Status FSM (Finite State Machine) — ✦ Always draw this:**
```
                    ┌──────────────────────────────────┐
                    │                                  ▼
  PENDING ──► CONFIRMED ──► ALLOCATED ──► SHIPPED ──► DELIVERED
     │            │              │
     ▼            ▼              ▼
  FAILED      CANCELLED      CANCELLED
                                 │
                                 └──► RETURN_REQUESTED ──► RETURNED
```

✦ **Say this**: "I model the order status as a Finite State Machine. Each transition is a valid business event. The system rejects invalid transitions — you can't go from PENDING to SHIPPED. This is enforced in the service layer, not just the database."

---

### Step 4 — High-Level Architecture

```
                          ┌───────────────────────────────────────┐
                          │         API Gateway / BFF             │
                          │  (auth, rate limit, routing)          │
                          └────────────────┬──────────────────────┘
                                           │
                    ┌──────────────────────┼────────────────────────┐
                    ▼                      ▼                        ▼
           ┌──────────────┐      ┌──────────────────┐    ┌──────────────────┐
           │  Order API   │      │  Order Query API │    │  Admin API       │
           │  (writes)    │      │  (reads/search)  │    │  (ops, override) │
           └──────┬───────┘      └────────┬─────────┘    └──────────────────┘
                  │                       │
                  ▼                       ▼
           ┌──────────────┐      ┌──────────────────┐
           │  Order DB    │      │  Redis Cache     │
           │  (Postgres)  │◄─────│  (order status,  │
           │              │      │   hot queries)   │
           └──────┬───────┘      └──────────────────┘
                  │
                  ▼
           ┌──────────────┐
           │ Event Bus    │  ← Kafka / SQS
           │ (Kafka)      │
           └──────┬───────┘
                  │
     ┌────────────┼────────────┬──────────────┐
     ▼            ▼            ▼              ▼
┌─────────┐ ┌──────────┐ ┌────────────┐ ┌──────────────┐
│Inventory│ │ Payment  │ │Fulfillment │ │Notification  │
│ Service │ │ Service  │ │ Service    │ │ Service      │
└─────────┘ └──────────┘ └────────────┘ └──────────────┘
```

---

### Step 5 — Deep Dives (What the Interviewer Will Ask)

**"How do you prevent two customers from buying the last item?"**
```
✦ "This is the inventory consistency problem. Two approaches:

Option A: Pessimistic Locking (SELECT FOR UPDATE)
  BEGIN;
  SELECT quantity FROM inventory WHERE sku = 'ABC' FOR UPDATE;
  -- Only one transaction holds this lock
  IF quantity >= requested THEN
    UPDATE inventory SET quantity = quantity - requested;
    INSERT INTO order_items ...;
  COMMIT;
  
  Pro: Simple, guaranteed no overselling
  Con: Lock contention at high volume, serializes throughput for hot SKUs

Option B: Optimistic Locking with version numbers
  -- Read: SELECT quantity, version FROM inventory WHERE sku = 'ABC'
  -- Write: UPDATE inventory SET quantity = q-1, version = version+1
  --        WHERE sku = 'ABC' AND version = {read_version}
  -- If 0 rows updated: conflict → retry
  
  Pro: No locks, high throughput for normal case
  Con: Retry storms during flash sales (all fail, all retry simultaneously)

Option C: Reservation model (what I'd recommend for OMS)
  1. On order placement: reserve inventory (soft lock in Redis with TTL)
  2. On payment confirmation: confirm reservation in DB
  3. On timeout or failure: release reservation

  Redis: SETNX inventory:reserved:{orderId}:{sku} {qty} EX 600
  
  Pro: Non-blocking, fast, handles payment delays
  Con: Reserved stock appears unavailable to other buyers during TTL"
```

**"How do you handle the Outbox pattern for reliability?"**
```
✦ "When an order is placed, I need to both save to the DB and publish to Kafka.
   These two operations are not atomic — Kafka could be down.
   
   Solution: The Outbox Pattern.
   1. In the same DB transaction as order creation, write to an 'outbox' table.
   2. A background poller reads unprocessed outbox events, publishes to Kafka.
   3. Marks event as published.
   
   This guarantees at-least-once delivery. Consumers must be idempotent."
```

---

## Case Study 2: Design an Order Promise System (ATP) 📦

> **Question**: "Design an Available-to-Promise (ATP) system. When a customer
> views a product or places an order, we must instantly and accurately promise
> a delivery date. This promise must be honored."

✦ **This is directly relevant to your target team.** Know this deeply.

---

### What is ATP?

Available-to-Promise answers: **"Can I promise this item, in this quantity, by this date, to this customer?"**

It requires:
1. Current inventory on-hand
2. Inbound supply (POs arriving from suppliers)
3. Existing demand (orders already placed)
4. Lead time (warehouse → carrier → customer)

```
ATP = (On-hand inventory) + (Inbound supply) - (Committed demand)
      ─────────────────────────────────────────────────────────────
                        Per SKU, per time bucket
```

---

### Step 1 — Clarify Requirements

```
✦ "Is this a hard promise (legally binding SLA) or a soft promise (best effort)?"
✦ "What's the required accuracy? Can we slightly overcommit (99% hit rate OK)?"
✦ "How many SKUs? 10K? 10M?"
✦ "How many promise queries per second? (Product page loads + checkout)"
✦ "What's the acceptable latency for a promise? 200ms? 500ms?"
✦ "Do we need to promise at SKU level, or also at bundle/variant level?"
✦ "How far out do we promise? 2 days? 30 days?"
✦ "Multiple fulfillment centers? Do we promise from the closest one?"
```

**Assume:**
- 50K SKUs, 5 fulfillment centers (FCs)
- 10,000 promise queries/sec (product pages + checkout)
- 500ms p99 latency requirement
- Hard promise: if we say "2 days," we must deliver in 2 days
- Accuracy target: <0.1% broken promises

---

### Step 2 — Estimate Scale

```
Promise reads: 10,000/sec
  → Cannot hit DB directly. Must use a pre-computed, cached "promise horizon" per SKU.

Inventory updates: When orders placed, received, cancelled
  → ~60 order/sec (from Case Study 1) × avg 2.5 items = 150 inventory events/sec
  → Manageable for real-time cache invalidation

ATP snapshot per SKU:
  50K SKUs × 5 FCs × 30 time buckets = 7.5M data points
  Each point ≈ 100 bytes → 750MB — fits in Redis cluster
```

---

### Step 3 — Data Model

```
Inventory Position (per SKU, per FC):
  ─────────────────────────────────────────────────────────
  on_hand_qty:       Int   -- physically present, not reserved
  reserved_qty:      Int   -- soft-reserved (orders in progress)
  committed_qty:     Int   -- hard-committed (payment confirmed)
  available_qty:     Int   -- on_hand - reserved - committed
  inbound_qty:       Int   -- POs arriving in next 30 days
  inbound_date:      Date  -- earliest PO arrival date

Promise Rule (per SKU or category):
  ─────────────────────────────────────────────────────────
  sku_id           VARCHAR
  fc_id            VARCHAR
  lead_time_days   INT    -- FC to customer shipping time
  safety_stock_pct FLOAT  -- 10% buffer, don't promise last unit
  promise_horizon  INT    -- max days out we can promise

Delivery Promise:
  ─────────────────────────────────────────────────────────
  promise_id      UUID
  order_id        UUID
  order_item_id   UUID
  sku             VARCHAR
  fc_id           VARCHAR
  promised_date   DATE
  committed_at    TIMESTAMP
  status          ENUM(ACTIVE, FULFILLED, BROKEN, CANCELLED)
```

---

### Step 4 — Architecture

```
Product Page / Checkout
        │
        ▼
  ┌─────────────────────────────┐
  │   Promise API               │
  │   POST /promise/query       │  ← "Can you promise SKU X, qty 2, by day Y?"
  └──────────────┬──────────────┘
                 │
                 ▼
  ┌─────────────────────────────┐
  │   ATP Calculation Engine    │
  │                             │
  │   1. Read ATP cache (Redis) │  ← Pre-computed horizon per SKU/FC
  │   2. Apply lead time        │
  │   3. Pick best FC           │
  │   4. Return promise date    │
  └──────────────┬──────────────┘
                 │
        ┌────────┴────────┐
        ▼                 ▼
  ┌──────────┐    ┌──────────────────┐
  │  Redis   │    │  Inventory DB    │
  │  ATP     │◄───│  (Postgres)      │
  │  Cache   │    │                  │
  └──────────┘    └──────────────────┘
                           ▲
                           │ inventory events
  ┌────────────────────────┤
  │  Event Consumer         │
  │  (Order placed/cancelled│
  │   PO received/updated)  │
  └─────────────────────────┘
```

---

### Step 5 — Deep Dives

**"How do you keep the ATP cache consistent with the source of truth?"**
```
✦ "I use an event-driven cache refresh pattern:

Every inventory-affecting event publishes to Kafka:
  ORDER_RESERVED    → decrease available_qty
  ORDER_CONFIRMED   → move reserved → committed
  ORDER_CANCELLED   → release reservation
  PO_RECEIVED       → increase on_hand
  PO_UPDATED        → update inbound supply

The ATP Consumer:
  1. Reads event from Kafka
  2. Recomputes ATP for affected SKU/FC
  3. Updates Redis: SET atp:{sku}:{fc} {serialized_atp} EX 3600

For the promise query:
  1. Read from Redis (cache miss → fallback to DB)
  2. Apply lead time + safety stock
  3. Return earliest date where ATP > 0

Eventual consistency gap: There's a ~100ms lag between order and cache update.
Mitigation: Safety stock buffer (promise 10% less than actual available).
This is the key trade-off: accuracy vs throughput."
```

**"What happens when your promised date is broken?"**
```
✦ "This requires a Promise Monitoring service:
  1. Async job runs daily: SELECT promises WHERE promised_date < NOW() AND status = ACTIVE
  2. For each: check actual fulfillment status
  3. If order not shipped by promised_date → trigger:
     a. Customer notification (apology + new ETA)
     b. Business alert (metric: promise_breach_rate)
     c. Compensation (voucher, upgrade shipping)
  
  Broken promises are tracked as a KPI.
  If breach_rate > 0.1%, pause promising for that SKU until replenished."
```

---

## Case Study 3: Design an Order Allocation Engine 🏭

> **Question**: "You have 1 million daily orders and 20 fulfillment centers (FCs)
> across the country. Design a system that decides which FC fulfills each order,
> minimizing cost and delivery time while respecting inventory constraints."

---

### Step 1 — Clarify Requirements

```
✦ "What are the optimization objectives? Cost? Speed? Both? Priority order?"
✦ "Do we allow split shipments? (Order fulfilled from 2 FCs)"
✦ "Can we reallocate an order after initial allocation? (FC runs out of stock)"
✦ "What's the SLA for allocation? Real-time at checkout, or async after order placed?"
✦ "Do FCs have capacity limits per day?"
✦ "Do carriers have cutoff times? (Ship before 3pm for same-day)"
✦ "Are there restricted items? (Hazmat shipped only from specific FCs)"
```

**Assume:**
- 1M orders/day (~12 orders/sec normal, 100/sec peak)
- 20 FCs across the US
- Optimize for: delivery speed first, cost second
- Split shipments allowed (up to 3 splits)
- Allocation must complete within 5 seconds of order placement
- Reallocation allowed until order is picked

---

### Step 2 — The Allocation Algorithm

This is the intellectual meat. Describe it clearly.

**Inputs per order:**
```
- Customer zip code (→ distance to each FC)
- Order items × quantities
- Required delivery date (from promise)
```

**FC Scoring (for each candidate FC):**
```
score(fc, order) = w1 × delivery_speed(fc, customer_zip)
                + w2 × inventory_availability(fc, items)
                + w3 × shipping_cost(fc, customer_zip, weight)
                + w4 × fc_workload_factor(fc)  ← don't overload one FC

where w1 + w2 + w3 + w4 = 1.0 (configurable weights)
```

**Algorithm (simplified):**
```
1. Filter: Remove FCs that can't fulfill ANY item in the order.
2. For single-FC fulfillment:
     Score all FCs that can fulfill the entire order.
     Pick the highest-scoring FC.
3. If no single FC can fulfill: try split fulfillment.
     Greedy: Pick the best FC for the most items.
     Assign remaining items to next best FC.
     Max 3 splits.
4. If still unfulfillable: hold order, trigger re-promise flow.
```

---

### Step 3 — Architecture

```
Order Placed
     │
     ▼
┌────────────────────────────────────────────────────┐
│               Allocation Service                   │
│                                                    │
│  1. Load inventory snapshot (Redis, all 20 FCs)    │
│  2. Load FC metadata (location, capacity, carrier  │
│     cutoffs, restricted items)                     │
│  3. Run allocation algorithm                       │
│  4. Reserve inventory (optimistic lock)            │
│  5. Publish AllocationDecision event               │
└────────────────────────────────────────────────────┘
         │                          │
         ▼                          ▼
┌──────────────────┐     ┌─────────────────────────┐
│  Inventory DB    │     │  Kafka: allocation.      │
│  (reserve qty)   │     │  decided                │
└──────────────────┘     └──────────┬──────────────┘
                                    │
                  ┌─────────────────┼──────────────────┐
                  ▼                 ▼                   ▼
          ┌────────────┐   ┌──────────────┐   ┌──────────────┐
          │ WMS        │   │  Promise     │   │ Notification │
          │ (Warehouse │   │  Service     │   │ Service      │
          │ Mgmt Sys)  │   │  (update     │   │ (confirm     │
          │            │   │  committed)  │   │  to cust.)   │
          └────────────┘   └──────────────┘   └──────────────┘
```

---

### Step 4 — Deep Dives

**"How do you handle reallocation when the allocated FC goes out of stock?"**
```
✦ "This is the most complex scenario. Three triggers can cause reallocation:

1. FC discovers a pick failure (item is physically missing, damage, etc.)
2. Inventory reconciliation finds a discrepancy
3. FC is temporarily shut down (weather, outage)

My approach:
  - AllocationMonitor service subscribes to inventory events
  - When available_qty at allocated FC < order_qty:
    - Check if order is still in PRE_PICK status (reallocatable)
    - Run allocation algorithm again, excluding current FC
    - If new FC found: update allocation, re-notify WMS
    - If no FC found: trigger customer notification + re-promise

Key constraint: Once a pick has started, we CANNOT reallocate.
  Use FC order status as the guard: only reallocate if status = PENDING_PICK."
```

**"How do you handle flash sales where 50,000 orders arrive in 60 seconds?"**
```
✦ "This is a thundering herd on inventory. My approach:

1. Pre-compute allocation slots before the sale:
   'Reserve 10,000 units at FC-West, 8,000 at FC-East, etc.'
   This is a capacity planning step, not real-time.

2. Use a queue-based allocation:
   Orders enter a priority queue (by order time).
   Allocator processes in FIFO order.
   First 10,000 orders get FC-West. Next 8,000 get FC-East.

3. Inventory is decremented atomically in Redis using Lua scripts:
   local qty = redis.call('GET', key)
   if tonumber(qty) >= requested then
     redis.call('DECRBY', key, requested)
     return 1
   end
   return 0  -- out of stock at this FC

4. If Redis says out-of-stock, check next FC immediately.
5. Unfulfillable orders trigger waitlist or backorder flow.

The key insight: Serialize writes per FC using Redis atomic ops.
This avoids DB lock contention and scales to any write rate."
```

---

## Case Study 4: Design a Real-Time Inventory Tracking System 📊

> **Question**: "Design a system that tracks inventory levels in real-time
> across 20 warehouses, 5 million SKUs, and serves 50,000 inventory
> queries/sec with <10ms latency."

---

### Step 1 — Clarify

```
✦ "What events change inventory? Receipts, picks, damages, returns, adjustments?"
✦ "Reads: current on-hand only, or also reservations, inbound?"
✦ "Write consistency: can we tolerate 1s lag, or must it be real-time?"
✦ "Is negative inventory allowed? (Backorders?)"
✦ "Audit requirement? Full history of every adjustment?"
```

---

### Step 3 — The Key Insight: CQRS + Event Sourcing

```
✦ "Standard CRUD at 50K reads/sec + high writes is an anti-pattern.
   I'd use CQRS (Command Query Responsibility Segregation):

Write side (Commands):
  - All inventory changes are EVENTS (never direct updates)
  - InventoryReceived, ItemPicked, DamageReported, ReturnProcessed
  - Events appended to an append-only log (Kafka or event store)
  - Current state = REPLAY of all events for a SKU/FC pair

Read side (Queries):
  - Materialized views maintained by event consumers
  - Current inventory = pre-computed, cached in Redis
  - Zero DB reads for 50K/sec queries (Redis handles it)

Benefits:
  - Full audit trail (every event logged)
  - Can recompute current state if Redis is wiped (replay events)
  - Read and write paths scale independently"
```

**Architecture:**
```
Pick/Receive/Adjust                       Customer/OMS Queries
       │                                          │
       ▼                                          ▼
┌─────────────────┐                    ┌──────────────────┐
│  Inventory      │                    │  Inventory       │
│  Command API    │                    │  Query API       │
└────────┬────────┘                    └────────┬─────────┘
         │ append event                          │
         ▼                                       ▼
┌─────────────────┐                    ┌──────────────────┐
│  Event Log      │                    │  Redis           │
│  (Kafka /       │ ──► Consumer ──►   │  inventory:      │
│  Event Store)   │    (materializes)  │  {sku}:{fc} = qty│
└─────────────────┘                    └──────────────────┘
         │ (also persisted)                      ▲
         ▼                                       │ rebuild on cache miss
┌─────────────────┐                    ┌──────────────────┐
│  Event Store DB │ ─────────────────► │  Postgres        │
│  (append-only)  │                    │  (materialized   │
└─────────────────┘                    │   inventory view)│
                                       └──────────────────┘
```

---

### Deep Dive: Preventing Overselling

```
✦ "The invariant we must enforce: available_qty never goes below 0
   (unless backorders are allowed).

Approach — Atomic Redis decrement with floor check:
  Lua script (atomic in Redis):
    local current = tonumber(redis.call('GET', key) or 0)
    if current >= requested then
      redis.call('DECRBY', key, requested)
      return current - requested  -- new qty
    else
      return -1  -- insufficient stock
    end

  Why Lua? Redis executes Lua scripts atomically.
  No two concurrent requests can both read qty=1 and both decrement.

If Redis fails (cache miss or unavailable):
  Fallback to DB with SELECT FOR UPDATE (serialized, slower but safe).

Distributed inventory (multiple DCs):
  Don't try to maintain a global atomic counter across DCs.
  Instead: allocate inventory quotas per DC (DC-East owns 5000 units,
  DC-West owns 3000). Each DC manages its own quota atomically.
  This avoids cross-DC coordination entirely."
```

---

## Case Study 5: Design an Order Notification System 🔔

> **Question**: "Design a system that sends order status notifications
> to customers via Email, SMS, and Push at every order lifecycle event.
> It must handle 2M notifications/day with guaranteed delivery."

---

### Step 1 — Clarify

```
✦ "What triggers a notification? All status changes, or selected ones?"
✦ "Which channels? Email, SMS, Push. In-app too?"
✦ "User preferences — can customers opt out of SMS but keep email?"
✦ "Guaranteed delivery — what does that mean? At-least-once? No duplicates?"
✦ "What's the SLA? Notification within X minutes of the event?"
✦ "Retries? How many? What's the dead letter queue strategy?"
```

---

### Architecture

```
Order Events (Kafka: order.status.changed)
        │
        ▼
┌──────────────────────────────────┐
│     Notification Dispatcher      │
│                                  │
│  1. Read order event             │
│  2. Load customer preferences    │
│  3. Load message template        │
│  4. Resolve channels for user    │
│  5. Enqueue per-channel jobs     │
└───────────┬──────────────────────┘
            │
  ┌─────────┼──────────┬──────────────┐
  ▼         ▼          ▼              ▼
[Email Q] [SMS Q] [Push Q]      [In-App Q]
  │         │          │
  ▼         ▼          ▼
SendGrid  Twilio    Firebase/APNs
```

**Key design decisions:**

**1. Idempotency — No Duplicate Notifications:**
```sql
-- Notification dedup table
notifications_sent (
  idempotency_key  VARCHAR  PRIMARY KEY,  -- hash(order_id + event_type + channel)
  sent_at          TIMESTAMP,
  status           ENUM(SENT, FAILED)
)

-- Before sending: check if already sent
-- After sending: record in this table
-- Retry safely: second attempt hits "already sent" → skip
```

**2. Retry with Dead Letter Queue:**
```
Attempt 1: Send → provider error
  Wait 30s → Attempt 2: Still fails
  Wait 2m  → Attempt 3: Still fails
  → Move to DLQ (Dead Letter Queue)
  → Alert on-call if DLQ depth > threshold
  → Manual investigation or re-process next day
```

**3. Template Engine:**
```
Template: "Your order {order_id} has been shipped via {carrier}.
           Track it here: {tracking_url}"
           
Localization: template_id + locale → localized template
              "order.shipped.en_US", "order.shipped.fr_FR"
```

---

## Case Study 6: Design a Distributed Job Scheduler ⏱️

> **Question**: "Design a job scheduler that runs millions of delayed and
> recurring jobs reliably — like sending a 'rate your order' email 3 days
> after delivery, or reconciling inventory every hour."

---

### Why This Matters for Your Domain

In an OMS/ATP context, you need schedulers for:
- Auto-cancel unpaid orders after 30 minutes
- Release expired inventory reservations
- Daily inventory reconciliation
- Send "your order is late" notifications
- Trigger re-promise flow for breached promises

---

### Step 1 — Clarify

```
✦ "Scale: how many jobs/day? 10M? 1B?"
✦ "Delay range: seconds to days? Or also months-years?"
✦ "Recurring jobs: cron-like (every hour) or interval-based (every 5 min)?"
✦ "Delivery guarantee: at-least-once, at-most-once, or exactly-once?"
✦ "Max latency between scheduled time and execution: 1s? 10s?"
✦ "Job payload size? Small (trigger only) or large (carry data)?"
```

---

### Architecture

```
Job Producer                         Job Workers
     │                                   │
     ▼                                   │
POST /jobs                               │
{                                        │
  "type": "send_review_email",           │
  "run_at": "2024-01-05T10:00:00Z",      │
  "payload": {"order_id": "abc123"},     │
  "max_retries": 3                       │
}                                        │
     │                                   │
     ▼                                   │
┌────────────────────────┐               │
│   Job Store (Postgres) │               │
│                        │               │
│  jobs (                │               │
│    id, type,           │               │
│    run_at,  ──────────►│ Poller        │
│    status,             │ (every sec)   │
│    payload,            │ SELECT jobs   │
│    attempts            │ WHERE run_at  │
│  )                     │ <= NOW()      │
└────────────────────────┘ AND status=   │
                           'PENDING'     │
                           LIMIT 100     │
                           FOR UPDATE    │
                           SKIP LOCKED   │
                                │        │
                                ▼        │
                          ┌──────────┐   │
                          │ Kafka /  │   │
                          │ SQS      │───┘
                          │ Queue    │
                          └──────────┘
```

**The critical SQL pattern — `SKIP LOCKED`:**
```sql
-- Multiple pollers can run in parallel — each gets different jobs
-- SKIP LOCKED: if a row is locked by another transaction, skip it
SELECT id, type, payload FROM jobs
WHERE run_at <= NOW()
  AND status = 'PENDING'
ORDER BY run_at ASC
LIMIT 100
FOR UPDATE SKIP LOCKED;

-- After picking up the job:
UPDATE jobs SET status = 'PROCESSING', started_at = NOW() WHERE id = $1;
```

**Handling stuck jobs (worker died mid-execution):**
```sql
-- Heartbeat: worker updates heartbeat_at every 30s
-- Watchdog: finds jobs with heartbeat_at > 2 minutes ago → reset to PENDING
UPDATE jobs
SET status = 'PENDING', attempts = attempts + 1
WHERE status = 'PROCESSING'
  AND heartbeat_at < NOW() - INTERVAL '2 minutes'
  AND attempts < max_retries;
```

---

## Case Study 7: Design an Activity Feed / Order Timeline 📜

> **Question**: "Design an order timeline/activity feed where customers and
> CS agents can see every event that happened to an order in real-time."

This is essentially an **event-sourced audit log** with a real-time feed.

---

### The Design

**Core insight**: Order events ARE the timeline. We already append to `order_events` table. The challenge is:
1. Making it queryable efficiently
2. Real-time push to open browser sessions
3. CS agents viewing many orders simultaneously

```
Events flowing in:
  ORDER_PLACED → PAYMENT_CONFIRMED → ALLOCATED → PICK_STARTED
  → SHIPPED → OUT_FOR_DELIVERY → DELIVERED

Architecture:
                             ┌─── WebSocket Hub ───┐
Order Events (Kafka)         │                     │
        │                    │  Customer's browser ◄┤
        ▼                    │  CS Agent browser   ◄┤
┌──────────────────┐         └─────────────────────┘
│ Feed Service     │──────►  SSE / WebSocket push
│                  │         (long-lived connections)
│ - Subscribe to   │
│   order events   │──────►  Redis Pub/Sub
│ - Fan out to     │         (fanout to connected clients)
│   connected      │
│   clients        │──────►  Postgres
└──────────────────┘         (persist all events for history)
```

**Real-time delivery via Server-Sent Events (SSE):**
```go
// Client subscribes to order timeline
GET /orders/{orderId}/stream
Accept: text/event-stream

// Server pushes events as they arrive:
data: {"event":"SHIPPED","carrier":"FedEx","tracking":"1234567890","ts":"..."}

data: {"event":"OUT_FOR_DELIVERY","eta":"2024-01-05T14:00:00Z","ts":"..."}
```

---

## Case Study 8: Design a Product / Order Search System 🔍

> **Question**: "Design a search system that lets customers search their
> order history by keyword, date range, product name, status, and lets
> CS agents search across ALL orders."

---

### Step 1 — Clarify

```
✦ "Full-text search on what fields? Order ID, product name, SKU, address?"
✦ "Faceted filtering? (Filter by status, date range, channel)"
✦ "Customer search vs CS agent search — different access scopes"
✦ "Scale: 100M orders to search over?"
✦ "Latency: 200ms p99?"
✦ "Freshness: how soon after order creation does it appear in search?"
```

---

### Architecture — CQRS with Elasticsearch

```
Order Created / Updated
        │
        ▼
   Kafka: order.events
        │
        ▼
┌──────────────────────────────┐
│  Search Indexer Service       │
│                               │
│  - Consumes order events      │
│  - Transforms to search doc   │
│  - Upserts to Elasticsearch   │
└──────────────────────────────┘
        │
        ▼
┌──────────────────────────────┐
│    Elasticsearch Cluster      │
│                               │
│  Index: orders                │
│  {                            │
│    order_id, customer_id,     │
│    status, channel,           │
│    created_at,                │
│    items: [{name, sku}],      │
│    shipping_address           │
│  }                            │
└──────────────────────────────┘
        ▲
        │ query
┌──────────────────────────────┐
│    Search API                 │
│                               │
│  GET /search/orders           │
│  ?q=blue+jacket               │
│  &status=DELIVERED            │
│  &from=2024-01-01             │
│  &size=20&page=2              │
└──────────────────────────────┘
```

**Elasticsearch query (what happens under the hood):**
```json
{
  "query": {
    "bool": {
      "must": {
        "multi_match": {
          "query": "blue jacket",
          "fields": ["items.name^3", "items.sku", "order_id"]
        }
      },
      "filter": [
        { "term":  { "customer_id": "123" }},
        { "term":  { "status": "DELIVERED" }},
        { "range": { "created_at": { "gte": "2024-01-01" }}}
      ]
    }
  },
  "sort": [{ "created_at": "desc" }],
  "from": 20, "size": 20
}
```

**CS Agent search (all orders, no customer_id filter):**
```
✦ "CS agents search across all customers — I'd use a separate index with
   stricter access controls, rate limiting per agent, and audit logging
   of every CS search query for compliance."
```

---

## Master Interview Cheat Sheet — Your Domain

```
You are interviewing for Order Management / Order Promises / Order Allocation.

Always lead with:
  ✦ "In an OMS context, the most critical invariant is: no overselling.
     I'll model the inventory as an atomic counter with reservations."

Key patterns to mention:
  ┌─────────────────────────────────────────────────────┐
  │  Domain Pattern       → Technical Solution          │
  ├─────────────────────────────────────────────────────┤
  │  No overselling        → Redis atomic DECRBY + Lua  │
  │  Order lifecycle       → FSM + event sourcing       │
  │  Reliable events       → Outbox pattern             │
  │  ATP accuracy          → Pre-computed cache + TTL   │
  │  Allocation            → Score + greedy + fallback  │
  │  Delayed jobs          → Postgres SKIP LOCKED       │
  │  Audit trail           → Append-only event log      │
  │  Search                → CQRS + Elasticsearch       │
  │  Real-time updates     → Kafka → SSE/WebSocket      │
  │  Flash sales           → Queue + pre-allocated slots│
  └─────────────────────────────────────────────────────┘

Trade-off phrases that impress interviewers:
  ✦ "I'd use optimistic locking for normal load, but switch to
     a reservation model during flash sales to avoid retry storms."
  ✦ "This is a CP vs AP trade-off. For inventory, I choose CP —
     I'd rather reject an order than oversell."
  ✦ "I'd start with a monolith here. The complexity of splitting
     OMS into microservices isn't justified until you hit scale
     or independent team ownership is needed."
  ✦ "The Outbox pattern adds complexity but it's the only way to
     guarantee that a DB write and a Kafka publish are atomic."
  ✦ "For the promise system, I accept eventual consistency of the
     ATP cache — a 100ms lag is fine if I use safety stock buffers."
```

---

## Session 6: The `any` Type, Interface Boxing & Mixed-Type Collections ✅

### The Question That Started It All

*"Can I have a mixed-type slice in Go?"*

Yes — `[]any` works. But the real lesson isn't the syntax. It's **what the runtime does** when you put a value into `any`, and **what it costs**.

### `any` = `interface{}` — Nothing More

```go
// In Go source: builtin/builtin.go
type any = interface{}  // type alias — identical at compiler level
```

`any` has zero methods → every type satisfies it → it accepts anything. But Go Proverb #7 warns: *"interface{} says nothing."*

### The `eface` Runtime Struct

When you assign a value to `any`, Go creates an `eface`:

```
runtime.eface (16 bytes)
┌──────────────────────┬──────────────────────┐
│  _type  *_type       │  data unsafe.Pointer │
│  (what type is it?)  │  (where is the data?)│
└──────────────────────┴──────────────────────┘
```

So `[]any` is really an array of 16-byte `eface` structs — NOT the actual values:

```go
mixed := []any{42, "hello", true}
// Runtime sees: [{intType, ptrTo42}, {stringType, ptrToHello}, {boolType, ptrToTrue}]
```

### Boxing: The Three Paths

**"Boxing"** = wrapping a concrete value into an interface value. The cost depends entirely on what type you're boxing:

```
                 Boxing a value into any
                        │
          ┌─────────────┴─────────────┐
          │                           │
    Pointer-shaped?             Value type?
    (*T, map, chan, func)       (int, string, struct, etc.)
          │                           │
          ▼                    Small int 0-255?
    ZERO COST                  ┌──────┴──────┐
    (pointer IS the value)    YES            NO
                               │              │
                               ▼              ▼
                         ZERO COST        HEAP ALLOC
                         (staticuint64s)  (mallocgc)
```

### Path 1: Direct Interface (Pointer Types) — FREE

Types whose values are already pointers store directly in `eface.data`:

```go
var a any = &User{Name: "sam"}  // *User → pointer value goes in data. FREE.
var b any = myMap               // map is *runtime.hmap. FREE.
var c any = myChan              // chan is *runtime.hchan. FREE.
var d any = myFunc              // func is *runtime.funcval. FREE.
```

### Path 2: `staticuint64s` (Small Values 0-255) — FREE

Go pre-allocates a static array of 256 values in the binary:

```go
// runtime/iface.go
var staticuint64s = [256]uint64{0, 1, 2, 3, ..., 255}
```

When you box a small value:

```go
var a any = 42     // data → &staticuint64s[42]. No allocation!
var b any = true   // data → &staticuint64s[1]. No allocation!
var c any = byte('A') // data → &staticuint64s[65]. No allocation!
```

**But `float64` doesn't benefit!** `float64(1.0)` has bits `0x3FF0000000000000` — way beyond 255. Even `var a any = 1.0` allocates.

### Path 3: Heap Allocation — The Expensive Path

Everything else gets heap-allocated via the `convT` family:

```go
var a any = 256         // int > 255 → convT64 → mallocgc(8 bytes)
var b any = 3.14        // float64 → convT64 → mallocgc(8 bytes)
var c any = "hello"     // string → convTstring → mallocgc(16 bytes, GC scans!)
var d any = []int{1,2}  // slice → convTslice → mallocgc(24 bytes, GC scans!)
var e any = BigStruct{} // struct → convT → mallocgc(sizeof, GC scans if has ptrs)
```

### The Complete Cost Table

```
🟢 FREE (zero allocation):
   *T, map, chan, func        → pointer-shaped, direct interface
   int/uint 0-255             → staticuint64s
   bool                       → 0 or 1, always in staticuint64s
   byte (uint8)               → full range covered

🟡 CHEAP (small heap alloc, no GC scanning):
   int/float64 > 255          → 8B alloc, but noscan (no pointers)
   struct{x, y int}           → alloc, but noscan

🔴 EXPENSIVE (heap alloc + GC scanning):
   string                     → 16B header copy, has pointer → GC scans
   slice                      → 24B header copy, has pointer → GC scans
   struct with pointer fields  → full copy, GC must scan all pointers
```

### The Equality Trap — Interview Favorite

```go
var a any = 42
var b any = 42
fmt.Println(a == b)  // true ✅ — ints are comparable

var c any = []int{1, 2}
var d any = []int{1, 2}
fmt.Println(c == d)  // 💥 PANIC — slices are NOT comparable
```

**Rule:** `==` on `any` values is a **runtime check**. If the underlying type isn't
comparable (slices, maps, funcs), it panics. The compiler can't catch this!

### JSON Unmarshalling into `any` — Know the Types

```go
var data any
json.Unmarshal([]byte(`{"age": 30, "scores": [95, 87]}`), &data)
// age is float64(30), NOT int!
// scores is []any{float64(95), float64(87)}, NOT []int!
```

**Trap:** JSON numbers are ALWAYS `float64`. JSON objects → `map[string]any`. Arrays → `[]any`.

### Generics Replaced Most `any` Uses (Go 1.18+)

```go
// ❌ Pre-generics — type-unsafe, boxing cost
func Contains(slice []any, target any) bool { ... }

// ✅ Post-generics — type-safe, ZERO boxing, inlineable
func Contains[T comparable](slice []T, target T) bool { ... }
```

The decision tree:
- Know the type? → Use concrete type
- Need to support multiple types? → Use generics with constraints
- Truly dynamic/unknown type? → Use `any` (rare in well-designed code)

### Performance Chain: Why Boxing Matters at Scale

```
Boxing → heap allocation → GC scanning → pointer indirection → cache misses → no inlining
         ↑                  ↑              ↑                     ↑
    convT/mallocgc    GC must trace    CPU follows pointer    Data not contiguous
                      all pointers     to read actual value   in memory
```

In a hot loop processing 1M items:
- `[]int` → contiguous memory, CPU cache-friendly, zero GC → **fast**
- `[]any` → scattered heap objects, GC scans every pointer → **5-10x slower**

### Interview Trap: `fmt.Sprintf` vs `strconv.Itoa`

```go
fmt.Sprintf("%d", n)   // n gets boxed into any (variadic ...any) → allocation
strconv.Itoa(n)        // no interface, no boxing → zero allocation
```

This is why `strconv` functions exist alongside `fmt` — for hot paths where boxing cost matters.

### 📖 Deep Dive

Full reference document: [`learnings/12_any_type_boxing.md`](learnings/12_any_type_boxing.md)
Covers: `convT` family internals, `staticuint64s` implementation, GC shape stenciling for generics, benchmark patterns.

### Check Questions

1. You have `var a any = myStruct` where `myStruct` has 3 string fields. How many allocations happen and why?
2. Why does `var a any = 1.0` allocate but `var a any = 0` does not?
3. You're building a high-throughput event pipeline processing 100k events/sec. Someone suggests using `[]any` to hold mixed event types. What's wrong and what do you suggest instead?


