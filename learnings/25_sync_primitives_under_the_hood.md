# 25 — sync Primitives Under the Hood

> **Go's sync package provides the building blocks for shared-memory concurrency.**
> Every primitive has specific internal behavior, performance characteristics,
> and production gotchas. This chapter goes under the hood on each one.

---

## Table of Contents

1. [sync.Mutex: The Spinning Lock](#1-syncmutex-the-spinning-lock)
2. [sync.RWMutex: Readers and Writers](#2-syncrwmutex-readers-and-writers)
3. [sync.WaitGroup: The Barrier](#3-syncwaitgroup-the-barrier)
4. [sync.Once: Exactly Once, Forever](#4-synconce-exactly-once-forever)
5. [sync.Pool: The Per-P Object Cache](#5-syncpool-the-per-p-object-cache)
6. [sync.Map: The Specialized Concurrent Map](#6-syncmap-the-specialized-concurrent-map)
7. [sync.Cond: Condition Variables](#7-synccond-condition-variables)
8. [sync/atomic: Lock-Free Primitives](#8-syncatomic-lock-free-primitives)
9. [Choosing the Right Primitive](#9-choosing-the-right-primitive)
10. [Production Patterns](#10-production-patterns)
11. [Performance Cost Table](#11-performance-cost-table)
12. [Quick Reference Card](#12-quick-reference-card)
13. [Further Reading](#13-further-reading)

---

## 1. sync.Mutex: The Spinning Lock

### The Internal State Machine

`sync.Mutex` is a single `int32` with bit-packed state:

```go
// sync/mutex.go
type Mutex struct {
    state int32
    sema  uint32
}
```

The `state` field packs three pieces of information:

```
┌─────────────────────────────────────────────────────┐
│  Mutex state bits                                   │
│                                                     │
│  bit 0: locked (1 = someone holds the lock)        │
│  bit 1: woken  (1 = a goroutine has been woken)    │
│  bit 2: starvation mode (1 = fairness mode active) │
│  bits 3+: waiter count (number of blocked waiters) │
│                                                     │
│  state = waiters<<3 | starvation<<2 | woken<<1 | locked │
└─────────────────────────────────────────────────────┘
```

### Two Modes: Normal and Starvation

**Normal mode (default):**
1. Goroutine tries `atomic.CompareAndSwapInt32(&m.state, 0, 1)` (fast path)
2. If CAS succeeds → lock acquired, zero syscalls
3. If CAS fails → **spin** for a few iterations (up to 4 spins on multi-core)
4. If spinning fails → park on semaphore (`runtime_SemacquireMutex`)

**Starvation mode (Go 1.9+):**
- Activates when a waiter has been blocked for > 1ms
- In starvation mode, the lock is handed directly to the longest-waiting
  goroutine (FIFO). No spinning, no barging
- Deactivates when the waiter is the last in the queue or waited < 1ms

```
┌──────────────────────────────────────────────────────┐
│  Lock acquisition flow                                │
│                                                       │
│  atomic CAS ──► success ──► LOCKED (0 syscalls)      │
│      │                                                │
│      ▼ fail                                           │
│  spin (≤4 iters) ──► CAS succeeds ──► LOCKED         │
│      │                                                │
│      ▼ still fail                                     │
│  semaphore park ──► woken ──► CAS ──► LOCKED          │
│      │                        │                       │
│      │                        ▼ fail (>1ms)           │
│      │                   STARVATION MODE               │
│      │                   (direct handoff, FIFO)        │
│      ▼                                                │
│  OS futex(2) on Linux / kevent on macOS               │
└──────────────────────────────────────────────────────┘
```

### Performance Characteristics

| Scenario | Cost |
|----------|------|
| Uncontended lock/unlock | ~15ns (atomic CAS + atomic store) |
| Contended, wins spin | ~50-200ns (few CAS retries) |
| Contended, parks on semaphore | ~1-10μs (context switch) |
| Starvation mode handoff | ~500ns (direct handoff, no spinning) |

### The Copy Trap

**sync.Mutex MUST NOT be copied after first use.** The `go vet` tool detects this:

```go
// ❌ BUG: copying a locked mutex
var mu sync.Mutex
mu.Lock()
mu2 := mu  // copies locked state — mu2 starts locked, can never be unlocked!
```

This applies to all sync types: Mutex, RWMutex, WaitGroup, Once, Cond, Pool, Map.

---

## 2. sync.RWMutex: Readers and Writers

```go
type RWMutex struct {
    w           Mutex        // held if there are pending writers
    writerSem   uint32       // semaphore for writers to wait for readers
    readerSem   uint32       // semaphore for readers to wait for writers
    readerCount atomic.Int32 // number of pending readers
    readerWait  atomic.Int32 // number of departing readers
}
```

### Read-Write Semantics

- **Multiple readers** can hold `RLock()` simultaneously
- **Writers** need exclusive access: `Lock()` blocks until all readers release
- **New readers** block when a writer is waiting (prevents writer starvation)

```
┌───────────────────────────────────────────────────────┐
│  RWMutex states                                       │
│                                                       │
│  readerCount > 0, no writer:  readers running freely  │
│  readerCount > 0, writer:     writer waits for readers│
│  readerCount = 0, writer:     writer has exclusive    │
│  writer waiting:               new readers block      │
│                                (prevents starvation)  │
└───────────────────────────────────────────────────────┘
```

### The Write-Starvation Prevention Trick

When `Lock()` is called, the writer subtracts `rwmutexMaxReaders` (1<<30) from
`readerCount`, making it negative. This signals to `RLock()` callers that a
writer is waiting — new readers park instead of proceeding.

### When to Use RWMutex vs Mutex

| Pattern | Use |
|---------|-----|
| Read-heavy (95%+ reads) | `sync.RWMutex` |
| Equal reads and writes | `sync.Mutex` (simpler, similar perf) |
| Write-heavy | `sync.Mutex` (RWMutex overhead not worth it) |
| Very short critical section | `sync.Mutex` (spinning is efficient) |
| Need concurrent readers | `sync.RWMutex` |

**Benchmark before switching.** On low-core machines, `RWMutex` overhead
(two atomics for RLock/RUnlock vs one for Mutex) can make it slower than
plain Mutex even for read-heavy workloads.

---

## 3. sync.WaitGroup: The Barrier

```go
type WaitGroup struct {
    noCopy noCopy
    state  atomic.Uint64 // packs counter + waiter count
    sema   uint32
}
```

### How It Works

- `Add(n)` increments the counter atomically
- `Done()` calls `Add(-1)`
- `Wait()` blocks until the counter reaches zero

The counter and waiter count are packed into a single `uint64` for atomicity:

```
┌─────────────────────────────────────────────────────┐
│  WaitGroup.state (uint64)                           │
│                                                     │
│  Upper 32 bits: counter (set by Add/Done)          │
│  Lower 32 bits: number of goroutines in Wait()     │
│                                                     │
│  When counter reaches 0 and waiters > 0:           │
│  → release all waiters via semaphore               │
└─────────────────────────────────────────────────────┘
```

### Common Bugs

```go
// ❌ BUG: Add inside goroutine (race with Wait)
var wg sync.WaitGroup
for i := 0; i < n; i++ {
    go func() {
        wg.Add(1)   // too late! Wait() might already have returned
        defer wg.Done()
        work()
    }()
}
wg.Wait()

// ✅ CORRECT: Add before launching goroutine
var wg sync.WaitGroup
for i := 0; i < n; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        work()
    }()
}
wg.Wait()
```

---

## 4. sync.Once: Exactly Once, Forever

```go
type Once struct {
    done atomic.Uint32
    m    Mutex
}
```

### The Fast Path / Slow Path Design

```go
func (o *Once) Do(f func()) {
    if o.done.Load() == 0 {  // fast path: atomic load (most calls)
        o.doSlow(f)
    }
}

func (o *Once) doSlow(f func()) {
    o.m.Lock()               // slow path: mutex (first call only)
    defer o.m.Unlock()
    if o.done.Load() == 0 {  // double-check under lock
        defer o.done.Store(1)
        f()
    }
}
```

The fast path is a single atomic load — O(1), no contention. Only the first
caller hits the slow path (mutex + function execution). All other callers
either see `done == 1` immediately, or block on the mutex and find `done == 1`
when they acquire it.

### OnceValue and OnceFunc (Go 1.21+)

```go
// sync.OnceValue — compute once, return cached value
getter := sync.OnceValue(func() *Config {
    return loadConfig()  // called once, result cached
})
cfg := getter()  // all calls return same *Config

// sync.OnceFunc — like Once.Do but as a function
init := sync.OnceFunc(func() {
    setupDatabase()
})
init()  // first call runs fn, subsequent calls are no-ops
```

These are syntactic sugar over `sync.Once` but with better ergonomics.

---

## 5. sync.Pool: The Per-P Object Cache

`sync.Pool` is a per-P (per-processor) cache with a two-generation GC strategy.

### Internal Structure

```
┌────────────────────────────────────────────────────────┐
│  sync.Pool internals                                    │
│                                                         │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  │
│  │  P0     │  │  P1     │  │  P2     │  │  P3     │  │
│  │ private │  │ private │  │ private │  │ private │  │
│  │ shared[]│  │ shared[]│  │ shared[]│  │ shared[]│  │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘  │
│       │            │            │            │         │
│       └────────────┴────────────┴────────────┘         │
│                    victim cache                         │
│              (previous generation)                      │
└────────────────────────────────────────────────────────┘
```

- **private:** one object per P, lock-free access (fastest path)
- **shared:** per-P list, other Ps can steal from it (work-stealing)
- **victim cache:** after GC, current pools move to victim. Victim is cleared on
  next GC. Objects survive exactly one GC cycle in victim

### Get/Put Flow

```
Get():
  1. Check local P's private → return if non-nil (fastest)
  2. Pop from local P's shared list
  3. Steal from other P's shared lists
  4. Check victim cache (previous GC generation)
  5. Call pool.New() (slowest — allocation)

Put(x):
  1. Store in local P's private (if empty)
  2. Push to local P's shared list
```

### When to Use and When NOT to Use

**Use sync.Pool for:**
- Temporary buffers (`bytes.Buffer`, `[]byte`)
- Objects allocated in hot paths (JSON encode state, HTTP buffers)
- Objects that are expensive to allocate but cheap to reset

**Do NOT use sync.Pool for:**
- Connection pools (connections disappear on GC!)
- Objects that must have a bounded count
- Objects that take significant time to construct
- Long-lived objects (Pool is for temporary reuse only)

### The encoding/json Connection

The `encoding/json` package uses `sync.Pool` for its `encodeState` buffers:

```go
// encoding/json/encode.go
var encodeStatePool sync.Pool

func newEncodeState() *encodeState {
    if v := encodeStatePool.Get(); v != nil {
        e := v.(*encodeState)
        e.Reset()
        return e
    }
    return new(encodeState)
}
```

This is why repeated `json.Marshal` calls have low allocation overhead.

---

## 6. sync.Map: The Specialized Concurrent Map

`sync.Map` is NOT a general-purpose concurrent map. It's optimized for
exactly two patterns:

1. **Write-once, read-many:** keys written once, read frequently
2. **Disjoint key sets:** goroutines access non-overlapping keys

### Internal Structure

```go
type Map struct {
    mu     Mutex
    read   atomic.Pointer[readOnly] // lock-free reads
    dirty  map[any]*entry            // requires mu
    misses int                       // count of reads that hit dirty
}
```

```
┌────────────────────────────────────────────────────────┐
│  sync.Map dual-store architecture                       │
│                                                         │
│  ┌──────────┐     ┌──────────┐                         │
│  │  read    │     │  dirty   │                         │
│  │ (atomic) │     │ (mutex)  │                         │
│  │          │     │          │                         │
│  │ key→*entry│     │ key→*entry│                        │
│  │ key→*entry│     │ key→*entry│                        │
│  │          │     │ key→*entry│ ← new keys here only   │
│  └──────────┘     └──────────┘                         │
│                                                         │
│  Load: check read (lock-free) → miss → check dirty     │
│  Store existing: CAS on *entry (lock-free)             │
│  Store new key: mu.Lock, add to dirty                  │
│  Promotion: after enough misses, dirty → read          │
└────────────────────────────────────────────────────────┘
```

### When sync.Map Beats map+Mutex

| Workload | sync.Map | map + RWMutex | Winner |
|----------|----------|---------------|--------|
| 99% reads, stable keys | Very fast | Fast | sync.Map |
| Many goroutines, disjoint keys | Fast | Contended | sync.Map |
| 50/50 read/write | Slower | Faster | map + Mutex |
| Few keys, high contention | Slower | Faster | map + Mutex |
| Need typed keys/values | No type safety | Type-safe | map + Mutex |

**Default choice:** `map[K]V` + `sync.RWMutex`. Only use `sync.Map` when
profiling shows lock contention is a bottleneck AND your access pattern
matches one of the two optimized patterns.

---

## 7. sync.Cond: Condition Variables

```go
type Cond struct {
    noCopy  noCopy
    L       Locker          // usually &mu (Mutex or RWMutex)
    notify  notifyList      // list of waiting goroutines
    checker copyChecker
}
```

### The Wait/Signal/Broadcast Pattern

```go
cond := sync.NewCond(&mu)

// Waiter goroutine:
mu.Lock()
for !condition() {   // MUST be a loop, not if
    cond.Wait()      // atomically: unlock mu + sleep + re-lock on wake
}
// condition is true, mu is locked
doWork()
mu.Unlock()

// Signaler goroutine:
mu.Lock()
updateCondition()
cond.Signal()    // wake one waiter
// or cond.Broadcast()  // wake ALL waiters
mu.Unlock()
```

### Why the Loop?

`Wait()` can experience **spurious wakeups** — the goroutine is woken even
though the condition hasn't changed. The loop re-checks the condition.
This is the same pattern as Java's `Object.wait()` and POSIX `pthread_cond_wait`.

### When to Use Cond

Rarely. Most Go programs use channels instead. Cond is useful when:
- Multiple goroutines need to wait for the same condition
- You need `Broadcast()` (wake all) — channels can't do this cleanly
- The condition is complex (not just "data available")

---

## 8. sync/atomic: Lock-Free Primitives

The `sync/atomic` package provides low-level atomic operations that map
directly to CPU instructions (LOCK prefix on x86, LDAR/STLR on ARM).

### Key Types (Go 1.19+)

```go
var counter atomic.Int64
counter.Add(1)           // atomic increment
counter.Load()           // atomic read
counter.Store(42)        // atomic write
counter.CompareAndSwap(42, 43) // CAS

var flag atomic.Bool
flag.Store(true)

var ptr atomic.Pointer[Config]
ptr.Store(&newConfig)
cfg := ptr.Load()
```

### atomic.Value — Lock-Free Config Reload

```go
var config atomic.Value // holds *Config

// Writer (rare):
config.Store(&Config{Host: "prod", Port: 443})

// Reader (frequent, from any goroutine):
cfg := config.Load().(*Config)
```

**Production pattern:** config reload without locks. Writer publishes a new
`*Config`, readers always see a consistent snapshot. No torn reads.

### When to Use Atomic vs Mutex

| Need | Use |
|------|-----|
| Simple counter | `atomic.Int64` |
| Simple flag | `atomic.Bool` |
| Config swap (write-rare, read-often) | `atomic.Value` or `atomic.Pointer[T]` |
| Multiple fields updated together | `sync.Mutex` (atomics can't do multi-field) |
| Complex invariants | `sync.Mutex` |
| Anything beyond simple read/write | `sync.Mutex` |

---

## 9. Choosing the Right Primitive

```
┌────────────────────────────────────────────────────────────┐
│  Decision tree: which sync primitive?                       │
│                                                             │
│  Need to protect shared data?                               │
│  ├── Simple counter/flag → atomic.Int64 / atomic.Bool      │
│  ├── Config reload (rare write, frequent read) → atomic.Value │
│  ├── Map with concurrent access:                            │
│  │   ├── Write-once-read-many → sync.Map                   │
│  │   ├── Disjoint keys → sync.Map                          │
│  │   └── General → map + sync.RWMutex                      │
│  ├── Read-heavy struct → sync.RWMutex                      │
│  └── General → sync.Mutex                                  │
│                                                             │
│  Need to coordinate goroutines?                             │
│  ├── Wait for N goroutines → sync.WaitGroup                │
│  ├── One-time init → sync.Once                             │
│  ├── Reusable buffers → sync.Pool                          │
│  ├── Wake all waiters → sync.Cond (Broadcast)              │
│  └── Communicate data → channels (not sync)                │
│                                                             │
│  Need to deduplicate concurrent calls?                      │
│  └── singleflight (golang.org/x/sync/singleflight)        │
└────────────────────────────────────────────────────────────┘
```

**Go proverb:** *"Channels orchestrate; mutexes serialize."*

Use channels for communication between goroutines. Use mutexes for protecting
shared data structures. Don't force one pattern where the other is natural.

---

## 10. Production Patterns

### Pattern 1: Singleflight (Thundering Herd Protection)

```go
import "golang.org/x/sync/singleflight"

var group singleflight.Group

func GetUser(id string) (*User, error) {
    v, err, _ := group.Do(id, func() (any, error) {
        return db.QueryUser(id) // only ONE db call per id
    })
    if err != nil {
        return nil, err
    }
    return v.(*User), nil
}
```

10 goroutines requesting the same user ID → 1 database query.
This is critical for cache miss scenarios at high concurrency.

### Pattern 2: sync.Pool for HTTP Handlers

```go
var bufPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

func handler(w http.ResponseWriter, r *http.Request) {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufPool.Put(buf)
    
    // use buf for response construction
    json.NewEncoder(buf).Encode(response)
    w.Write(buf.Bytes())
}
```

### Pattern 3: Graceful Shutdown with WaitGroup

```go
func main() {
    var wg sync.WaitGroup
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()
    
    // Start workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            worker(ctx)
        }()
    }
    
    <-ctx.Done()        // wait for signal
    cancel()            // cancel all workers
    wg.Wait()           // wait for all workers to finish
    log.Println("shutdown complete")
}
```

### Pattern 4: Lazy Initialization

```go
// Before Go 1.21:
var (
    dbOnce sync.Once
    db     *sql.DB
)
func GetDB() *sql.DB {
    dbOnce.Do(func() {
        db = connectDB()
    })
    return db
}

// Go 1.21+:
var getDB = sync.OnceValue(func() *sql.DB {
    return connectDB()
})
// Usage: db := getDB()
```

---

## 11. Performance Cost Table

| Operation | Cost | Notes |
|-----------|------|-------|
| `Mutex.Lock()` uncontended | ~15ns | Single CAS |
| `Mutex.Lock()` contended (spin wins) | ~50-200ns | Few CAS retries |
| `Mutex.Lock()` contended (park) | ~1-10μs | OS context switch |
| `RWMutex.RLock()` uncontended | ~20ns | Two atomics (add + load) |
| `RWMutex.Lock()` | ~15ns + reader drain | Must wait for all readers |
| `WaitGroup.Add(1)` | ~5ns | Atomic add |
| `WaitGroup.Wait()` | ~5ns if counter=0 | Atomic load fast path |
| `Once.Do()` (after first call) | ~1ns | Atomic load, branch |
| `Once.Do()` (first call) | ~15ns + fn | Mutex + atomic store |
| `Pool.Get()` (private hit) | ~5ns | Local P, no lock |
| `Pool.Get()` (shared/steal) | ~50-200ns | May steal from other P |
| `Pool.Get()` (New call) | varies | Allocation |
| `atomic.Int64.Add(1)` | ~5ns | CPU LOCK ADD instruction |
| `atomic.Value.Load()` | ~2ns | Atomic pointer load |
| `sync.Map.Load()` (read hit) | ~10ns | Atomic load, no lock |
| `sync.Map.Store()` (new key) | ~100ns | Mutex + dirty map |

---

## 12. Quick Reference Card

```
┌───────────────────────────────────────────────────────────────────┐
│  sync PACKAGE — QUICK REFERENCE                                  │
├───────────────────────────────────────────────────────────────────┤
│                                                                   │
│  MUTEX                                                            │
│  var mu sync.Mutex                                               │
│  mu.Lock()   / mu.Unlock()         Exclusive lock                │
│  mu.TryLock()                      Non-blocking (Go 1.18+)      │
│                                                                   │
│  RWMUTEX                                                          │
│  var mu sync.RWMutex                                             │
│  mu.RLock()  / mu.RUnlock()        Shared read lock              │
│  mu.Lock()   / mu.Unlock()         Exclusive write lock          │
│                                                                   │
│  WAITGROUP                                                        │
│  var wg sync.WaitGroup                                           │
│  wg.Add(n)                         Increment counter             │
│  wg.Done()                         Decrement (= Add(-1))        │
│  wg.Wait()                         Block until counter = 0      │
│                                                                   │
│  ONCE                                                             │
│  var once sync.Once                                              │
│  once.Do(fn)                       Run fn exactly once           │
│  sync.OnceValue(fn)                Once + cache return value     │
│  sync.OnceFunc(fn)                 Once as standalone function   │
│                                                                   │
│  POOL                                                             │
│  pool := sync.Pool{New: func() any { ... }}                     │
│  obj := pool.Get()                 Get or create                 │
│  pool.Put(obj)                     Return for reuse              │
│  Objects may be GC'd at any time — NOT a connection pool!        │
│                                                                   │
│  MAP                                                              │
│  var m sync.Map                                                  │
│  m.Store(key, val)                 Set                           │
│  val, ok := m.Load(key)           Get                           │
│  m.LoadOrStore(key, val)          Get or set                    │
│  m.Delete(key)                     Remove                        │
│  m.Range(func(k,v any) bool)     Iterate (not consistent)      │
│                                                                   │
│  COND                                                             │
│  cond := sync.NewCond(&mu)                                      │
│  cond.Wait()                       Unlock + sleep + re-lock     │
│  cond.Signal()                     Wake one waiter              │
│  cond.Broadcast()                  Wake all waiters             │
│                                                                   │
│  ATOMIC (sync/atomic)                                             │
│  var n atomic.Int64                                              │
│  n.Add(1) / n.Load() / n.Store(v) / n.CompareAndSwap(old,new)  │
│  var v atomic.Value                                              │
│  v.Store(x) / v.Load()            Lock-free value swap          │
│                                                                   │
│  RULES                                                            │
│  1. NEVER copy sync types after first use (go vet catches this) │
│  2. Always Add() BEFORE launching goroutine                      │
│  3. Always Unlock() in defer (panic safety)                      │
│  4. Cond.Wait() MUST be in a for loop (spurious wakeups)        │
│  5. Pool objects can vanish on GC — don't store connections      │
│  6. sync.Map: only for write-once-read-many or disjoint keys    │
│  7. Prefer channels for communication, mutexes for protection   │
└───────────────────────────────────────────────────────────────────┘
```

---

## 13. Further Reading

- **sync/mutex.go:** ~200 lines — read the state machine comments
- **sync/pool.go:** Per-P pooling with victim cache
- **sync/map.go:** Dual read/dirty store architecture
- **sync/once.go:** 30 lines — the simplest and most elegant sync primitive
- **Russ Cox, "Go Memory Model" (2022 revision):** Defines happens-before
  relationships for all sync operations
- **Bryan Mills, "sync.Map" talk at GopherCon 2017:** When and why to use it

---

## Companion Exercises

Practice these concepts:
→ [`exercises/stdlib/09_sync/`](../exercises/stdlib/09_sync/) — 12 exercises
covering Mutex, RWMutex, WaitGroup, Once, Pool, sync.Map, Cond, atomic,
singleflight, and production patterns.
