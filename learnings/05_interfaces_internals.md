# Deep Dive: Go Interface Internals — iface, eface, itab & the Nil Trap

> Everything the runtime does when you assign a value to an interface,
> call a method through it, or compare it to nil.

---

## Table of Contents

0. [Why Interfaces? — The 6 Motivations](#0-why-interfaces--the-6-motivations)
1. [The Two Runtime Structs: `iface` and `eface`](#1-the-two-runtime-structs-iface-and-eface)
2. [The `itab` — The Method Dispatch Table](#2-the-itab--the-method-dispatch-table)
3. [Compile Time vs Runtime — Who Does What?](#3-compile-time-vs-runtime--who-does-what)
4. [Step-by-Step: Interface Assignment](#4-step-by-step-interface-assignment)
5. [Step-by-Step: Interface Method Call](#5-step-by-step-interface-method-call)
6. [The Three States of an Interface](#6-the-three-states-of-an-interface)
7. [The Nil Interface Trap](#7-the-nil-interface-trap)
8. [The Three Guards Against the Nil Trap](#8-the-three-guards-against-the-nil-trap)
9. [Type Assertions Under the Hood](#9-type-assertions-under-the-hood)
10. [Performance Implications](#10-performance-implications)
11. [Quick Reference Card](#11-quick-reference-card)

---

## 0. Why Interfaces? — The 6 Motivations

Before diving into runtime internals, it's essential to understand **why** interfaces
exist in Go and what problems they solve. Interfaces are not just a language feature —
they are the **only mechanism for abstraction and polymorphism** in Go.

### Key Prerequisite: Methods Don't Require Interfaces

A struct can have methods without any interface being involved:

```go
type User struct{ Name string }

func (u User) Greet() string { return "Hello, " + u.Name }

u := User{Name: "Mert"}
u.Greet() // works — no interface needed, direct call, inlineable
```

Methods belong to **types**. Interfaces define **contracts for consumers**. They are
separate concepts that the compiler connects when needed.

---

### Motivation 1: Polymorphism Without Inheritance

Go has no classes, no inheritance, no `extends`. Interfaces are the **only way** to write
a function that accepts multiple types based on shared behavior:

```go
func Process(r io.Reader) error {
    buf := make([]byte, 1024)
    _, err := r.Read(buf) // works for *os.File, *http.Response.Body,
    return err            // *bytes.Buffer, *gzip.Reader, *tls.Conn...
}
```

Without interfaces, you'd need a separate `ProcessFile()`, `ProcessHTTP()`,
`ProcessBuffer()` — one per type. Interfaces let you write one function that works
with any type sharing the same behavior.

**Under the hood:** The compiler verifies at the assignment site that the concrete type
has the required methods. At runtime, the `itab` dispatches the correct method pointer.
This is **structural typing** — no `implements` keyword, no registration.

---

### Motivation 2: Decoupling — Consumer Defines the Contract

This is Go's **most important interface design decision** and the biggest departure
from Java/C#. Understanding this deeply changes how you design entire systems.

#### The Java/C# Way: Producer Declares (Coupled)

In Java, the **producer** must know about and declare every interface it implements:

```java
// The interface lives in a shared/common package
public interface UserRepository {
    User getByID(int id);
}

// The producer MUST import and declare the interface
class PostgresRepo implements UserRepository {  // ← explicit declaration
    public User getByID(int id) { ... }
}
```

**The problem:** `PostgresRepo` is **coupled** to `UserRepository` at compile time.
If `UserRepository` adds a method, `PostgresRepo` breaks. If you want `PostgresRepo`
to satisfy a new interface from a different package, you must modify `PostgresRepo`'s
source code. The producer must know about every consumer upfront.

#### The Go Way: Consumer Defines (Decoupled)

In Go, the **consumer** defines what it needs. The producer is completely unaware:

```go
// ─── package postgres (the producer) ───────────────────────
// Knows NOTHING about any interface. Just a struct with methods.
package postgres

type UserRepo struct {
    db *sql.DB
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*User, error) {
    // query database...
}

func (r *UserRepo) Create(ctx context.Context, u *User) error {
    // insert into database...
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
    // delete from database...
}

// ─── package service (consumer A) ──────────────────────────
// Defines ONLY what it needs — 1 method out of 3
package service

type UserReader interface {
    GetByID(ctx context.Context, id int) (*User, error)
}

func NewUserService(reader UserReader) *UserService { ... }

// ─── package admin (consumer B) ────────────────────────────
// Defines what IT needs — different subset
package admin

type UserManager interface {
    GetByID(ctx context.Context, id int) (*User, error)
    Delete(ctx context.Context, id int) error
}

func NewAdminService(mgr UserManager) *AdminService { ... }
```

**Key observations:**
- `postgres.UserRepo` never imports `service` or `admin` packages
- `service` only asks for `GetByID` — it can't accidentally call `Delete`
- `admin` asks for `GetByID` + `Delete` — a different, broader contract
- Both interfaces are satisfied by the **same** `*UserRepo` without it knowing
- Each consumer gets the **narrowest interface** it needs (Interface Segregation)

#### The Import Graph: Dependencies Point Inward

This is the architectural consequence. Compare the import graphs:

```
JAVA — dependencies point outward (producer → consumer's interface):

  ┌────────────────┐         ┌────────────────────┐
  │  PostgresRepo  │────────▶│  UserRepository     │
  │  (producer)    │imports  │  (shared interface)  │
  └────────────────┘         └────────────────────┘
                                      ▲
  ┌────────────────┐                  │
  │  UserService   │──────────────────┘
  │  (consumer)    │imports
  └────────────────┘


GO — dependencies point inward (consumer defines what it needs):

  ┌────────────────┐
  │  PostgresRepo  │  ← knows nothing about anyone
  │  (producer)    │
  └────────────────┘
          ▲ (no import — compiler matches structurally)
          │
  ┌────────────────┐
  │  UserService   │  ← defines its own UserReader interface
  │  (consumer)    │
  └────────────────┘
```

In Go, there is **no shared interface package**. Each consumer defines exactly
what it needs. The producer doesn't import the consumer. The consumer doesn't
import the producer (in the interface definition — only at the wiring site in `main()`).

#### Why This Matters in Practice

**1. You can satisfy interfaces that don't exist yet.**

If someone writes a new package next year with an interface containing `GetByID`,
your `postgres.UserRepo` already satisfies it — without any modification. In Java,
you'd need to go back and add `implements NewInterface`.

**2. Each consumer gets the minimal contract (Interface Segregation).**

In Java, if `UserRepository` has 10 methods, every consumer sees all 10 — even if
it only uses 1. In Go, each consumer defines a 1-2 method interface with just what
it needs. This makes code easier to understand, test, and refactor.

**3. No "shared interface" package that everything depends on.**

In Java/C# enterprise projects, you often see a `common` or `contracts` package
that every module imports. Change one interface there → rebuild everything. In Go,
each consumer's interface is local and independent — changing one doesn't affect others.

**4. Swapping implementations requires zero changes in business logic.**

```go
// main.go — the composition root, the ONLY place that knows concrete types
func main() {
    var repo service.UserReader

    if cfg.UseCache {
        repo = redis.NewCachedRepo(postgres.NewUserRepo(db))  // stacked!
    } else {
        repo = postgres.NewUserRepo(db)
    }

    svc := service.NewUserService(repo)  // service doesn't care which one
}
```

The `service` package never changes. It doesn't know Redis exists. It doesn't
know PostgreSQL exists. It only knows `UserReader` — an interface it defined itself.

#### The Go Proverb

> *"Accept interfaces, return concrete types."*

This means:
- **Function parameters**: use interfaces (accept any type with matching behavior)
- **Function return values**: return the concrete type (give the caller full access)

```go
// ✅ Idiomatic Go
func NewUserService(store UserReader) *UserService { ... }
//                  ^^^^^^^^^^^         ^^^^^^^^^^^^
//                  interface (input)   concrete (output)

// ❌ Over-abstracted
func NewUserService(store UserReader) UserServiceInterface { ... }
//                                    ^^^^^^^^^^^^^^^^^^^^
//                  returning an interface hides the concrete type unnecessarily
```

Return concrete types so callers can access all methods, not just the interface subset.
Let the **next consumer** decide what interface to narrow it down to.

#### Comparison Table

| Aspect | Java/C# (Producer Declares) | Go (Consumer Defines) |
|--------|----------------------------|----------------------|
| Who defines the interface? | Shared package or the producer | The consumer, locally |
| Producer imports interface? | Yes — must declare `implements` | No — unaware of any interface |
| Adding a method to interface | Breaks all producers | Only affects consumers who defined it |
| Satisfying a new interface | Modify producer source code | Automatic — if methods match, it works |
| Shared interface package | Common pattern (`contracts`, `api`) | Not needed — no shared dependency |
| Minimal contracts per consumer | Hard — everyone sees the full interface | Natural — each consumer defines just what it needs |

---

### Motivation 3: Testability

This follows directly from Motivation 2. When a function depends on an interface,
you can inject a test double with zero frameworks:

```go
// Production
func NewService(store UserStore) *Service { ... }

// Test — 10 lines, no mocking library
type fakeStore struct {
    users map[int]*User
}

func (f *fakeStore) GetByID(_ context.Context, id int) (*User, error) {
    u, ok := f.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return u, nil
}

func TestService(t *testing.T) {
    svc := NewService(&fakeStore{users: map[int]*User{1: {Name: "Mert"}}})
    // test svc methods...
}
```

Without the interface, `NewService` would take `*sql.DB` or `*PostgresRepo` — and
your unit tests would need a real database.

---

### Motivation 4: Stdlib Integration Points

The standard library defines **small interfaces as extension hooks**. Implement them
and the entire Go ecosystem works with your type automatically:

| Interface | Method(s) | What You Unlock |
|-----------|-----------|-----------------|
| `fmt.Stringer` | `String() string` | Custom output in `fmt.Println`, `%v`, `%s` |
| `error` | `Error() string` | Entire error handling ecosystem (`errors.Is`, `errors.As`, `%w`) |
| `io.Reader` | `Read([]byte) (int, error)` | Every I/O function: `bufio`, `gzip`, `json.Decoder`, `io.Copy` |
| `io.Writer` | `Write([]byte) (int, error)` | `fmt.Fprintf`, `json.Encoder`, `http.ResponseWriter` |
| `sort.Interface` | `Len`, `Less`, `Swap` | `sort.Sort()` for any collection |
| `json.Marshaler` | `MarshalJSON() ([]byte, error)` | Custom JSON encoding |
| `json.Unmarshaler` | `UnmarshalJSON([]byte) error` | Custom JSON decoding |
| `http.Handler` | `ServeHTTP(w, r)` | HTTP server routing, middleware |
| `encoding.TextMarshaler` | `MarshalText() ([]byte, error)` | Used as map keys in JSON, YAML, etc. |

**Key insight:** `Error()` and `String()` are **not special built-in methods**. They are
ordinary methods that satisfy ordinary interfaces (`error` and `fmt.Stringer`). The
"magic" is that `fmt.Println` internally does a type assertion to `Stringer` — if your
type satisfies it, the custom output is used. No registration, no annotation.

**Your type can satisfy interfaces it has never heard of.** If someone writes a new
interface tomorrow with a `String() string` method, your type already satisfies it —
without changing a single line of your code. This is implicit satisfaction in action.

---

### Motivation 5: Behavioral Abstraction at Architecture Boundaries

In a production service, code is organized into **layers**. Each layer has a single
responsibility. The critical question is: **how do layers talk to each other?**

#### The Problem: Without Interfaces, Layers Are Welded Together

Imagine a simple order service with no interfaces:

```go
// ─── handler layer ─────────────────────────────
package handler

import "myapp/internal/service"  // ← imports concrete service
import "myapp/internal/postgres" // ← imports concrete repo (indirectly, through service)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
    svc := service.OrderService{Repo: &postgres.OrderRepo{DB: db}}
    //                                 ^^^^^^^^^^^^^^^^
    //     handler knows about postgres! If we switch to MongoDB,
    //     handler code changes too.
    svc.Create(order)
}

// ─── service layer ─────────────────────────────
package service

import "myapp/internal/postgres"  // ← welded to postgres

type OrderService struct {
    Repo *postgres.OrderRepo  // ← concrete type, not an interface
}

func (s *OrderService) Create(o Order) error {
    return s.Repo.Insert(o)  // directly calls postgres
}
```

**What's wrong here:**
- `service` imports `postgres` — business logic is coupled to infrastructure
- Changing the database means modifying `service` AND `handler`
- You can't test `service` without a running PostgreSQL instance
- Every package knows about every other package — a dependency web

#### The Solution: Interfaces as Seams Between Layers

Each layer defines an interface for what it **needs from the layer below**.
Each layer **never imports** the layer below — only the interface it defined:

```go
// ─── domain layer (innermost — no dependencies) ───────────
package domain

type Order struct {
    ID     string
    UserID string
    Amount float64
    Status string
}

// ─── service layer (defines what it needs from persistence) ─
package service

// The service defines this interface — NOT the repo package
type OrderRepository interface {
    Insert(ctx context.Context, o *domain.Order) error
    GetByID(ctx context.Context, id string) (*domain.Order, error)
    UpdateStatus(ctx context.Context, id, status string) error
}

// The service also defines this — what it needs from payment
type PaymentGateway interface {
    Charge(ctx context.Context, userID string, amount float64) (txID string, err error)
}

type OrderService struct {
    repo    OrderRepository  // interface — doesn't know if it's Postgres, Mongo, or a mock
    payment PaymentGateway   // interface — doesn't know if it's Stripe, PayPal, or a mock
}

func NewOrderService(repo OrderRepository, pay PaymentGateway) *OrderService {
    return &OrderService{repo: repo, payment: pay}
}

func (s *OrderService) PlaceOrder(ctx context.Context, o *domain.Order) error {
    txID, err := s.payment.Charge(ctx, o.UserID, o.Amount)
    if err != nil {
        return fmt.Errorf("payment failed: %w", err)
    }
    o.Status = "paid"
    if err := s.repo.Insert(ctx, o); err != nil {
        // In production: refund the payment here
        return fmt.Errorf("save order: %w", err)
    }
    return nil
}

// ─── postgres package (implements OrderRepository) ─────────
package postgres

// Knows NOTHING about the service package or its interface.
// Just a struct with methods that happen to match.
type OrderRepo struct {
    db *sql.DB
}

func (r *OrderRepo) Insert(ctx context.Context, o *domain.Order) error {
    _, err := r.db.ExecContext(ctx,
        "INSERT INTO orders (id, user_id, amount, status) VALUES ($1,$2,$3,$4)",
        o.ID, o.UserID, o.Amount, o.Status)
    return err
}

func (r *OrderRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) { ... }
func (r *OrderRepo) UpdateStatus(ctx context.Context, id, status string) error { ... }

// ─── stripe package (implements PaymentGateway) ────────────
package stripe

// Also knows NOTHING about the service package.
type Client struct {
    apiKey string
}

func (c *Client) Charge(ctx context.Context, userID string, amount float64) (string, error) {
    // call Stripe API...
    return txID, nil
}
```

#### The Import Graph: Clean Layer Separation

```
                    ┌────────────┐
                    │   domain   │  ← no imports (pure business entities)
                    └────────────┘
                          ▲
                          │ imports domain only
                    ┌────────────┐
                    │  service   │  ← defines interfaces, contains business logic
                    │            │     does NOT import postgres or stripe
                    └────────────┘
                          ▲
          ┌───────────────┼───────────────┐
          │               │               │
    ┌───────────┐  ┌────────────┐  ┌────────────┐
    │  postgres  │  │   stripe   │  │   handler   │
    │ (repo impl)│  │(pay impl)  │  │ (HTTP layer)│
    └───────────┘  └────────────┘  └────────────┘
          │               │               │
          └───────────────┼───────────────┘
                          │ all wired together in
                    ┌────────────┐
                    │   main()   │  ← the composition root
                    └────────────┘
```

**Notice:** `service` has **zero imports** from infrastructure packages. The arrows
point **inward** — outer layers depend on inner layers, never the reverse. This is
the **Dependency Inversion Principle** — and in Go, it happens naturally through
consumer-defined interfaces without any DI framework.

#### The Composition Root: `main()` Wires Everything

`main()` is the **only place** in the entire codebase that knows all concrete types:

```go
func main() {
    // --- infrastructure ---
    db, _ := sql.Open("postgres", cfg.DatabaseURL)
    repo := postgres.NewOrderRepo(db)
    pay := stripe.NewClient(cfg.StripeKey)

    // --- business logic ---
    svc := service.NewOrderService(repo, pay)
    //                             ^^^^  ^^^
    //     *postgres.OrderRepo satisfies service.OrderRepository ✓
    //     *stripe.Client satisfies service.PaymentGateway ✓
    //     compiler checks this HERE, at the wiring site

    // --- transport layer ---
    h := handler.NewOrderHandler(svc)
    http.ListenAndServe(":8080", h)
}
```

No DI container. No annotations. No reflection. Explicit constructor injection.
If the types don't match, you get a **compile-time error** right here in `main()`.

#### What This Enables in Practice

**1. Swap infrastructure without touching business logic.**

Migrating from PostgreSQL to MongoDB? Write a new `mongo.OrderRepo` with the same
methods. Change one line in `main()`. The `service` package never changes:

```go
// Before
repo := postgres.NewOrderRepo(db)

// After
repo := mongo.NewOrderRepo(client)

// service.NewOrderService(repo, pay) — unchanged, doesn't care
```

**2. Test business logic in complete isolation.**

```go
func TestPlaceOrder(t *testing.T) {
    // Fake repo — 5 lines, no database
    fakeRepo := &fakeOrderRepo{
        insertFn: func(ctx context.Context, o *domain.Order) error {
            return nil
        },
    }
    // Fake payment — 5 lines, no Stripe
    fakePay := &fakePayment{
        chargeFn: func(ctx context.Context, uid string, amt float64) (string, error) {
            return "tx-123", nil
        },
    }

    svc := service.NewOrderService(fakeRepo, fakePay)
    err := svc.PlaceOrder(ctx, &domain.Order{Amount: 99.99})
    if err != nil {
        t.Fatal(err)
    }
}
```

Tests run in milliseconds. No database. No network. No Stripe sandbox.

**3. Teams work on layers independently.**

- **Team A** works on the `postgres` package — implements new query optimizations
- **Team B** works on the `service` package — adds new business rules
- Neither team touches the other's code. The interface is the contract between them.
- As long as `postgres.OrderRepo` still has the matching methods, everything compiles.

**4. Feature flags and gradual migrations.**

```go
func main() {
    var repo service.OrderRepository

    if cfg.UseNewDB {
        repo = cockroachdb.NewOrderRepo(newDB)  // new implementation
    } else {
        repo = postgres.NewOrderRepo(oldDB)     // old implementation
    }

    svc := service.NewOrderService(repo, pay)   // service is unaware
}
```

Route 5% of traffic to CockroachDB, 95% to PostgreSQL. The service layer has
**no idea** this is happening. The interface makes both implementations interchangeable.

**5. Observability wrappers without modifying business logic.**

```go
// A logging wrapper that satisfies the same interface
type loggingRepo struct {
    next   service.OrderRepository  // wraps the real repo via interface
    logger *slog.Logger
}

func (r *loggingRepo) Insert(ctx context.Context, o *domain.Order) error {
    start := time.Now()
    err := r.next.Insert(ctx, o)  // delegate to real repo
    r.logger.Info("repo.Insert", "duration", time.Since(start), "err", err)
    return err
}

// In main() — wrap the real repo
repo := postgres.NewOrderRepo(db)
loggingRepo := &loggingRepo{next: repo, logger: logger}
svc := service.NewOrderService(loggingRepo, pay)  // service sees the same interface
```

No modifications to `postgres.OrderRepo`. No modifications to `service.OrderService`.
You added observability by stacking a wrapper — this is Motivation 6 (decoration)
working together with Motivation 5 (architecture boundaries).

#### The Layered Architecture Summary

| Layer | Defines | Imports | Knows About |
|-------|---------|---------|-------------|
| `domain` | Entities, value objects | Nothing | Nothing |
| `service` | Business logic + interfaces for dependencies | `domain` only | What behavior it needs, not who provides it |
| `postgres` | Repository implementation | `domain` + `database/sql` | How to store entities, not who uses them |
| `stripe` | Payment implementation | `domain` + Stripe SDK | How to charge, not who calls it |
| `handler` | HTTP transport | `domain` + `service` | How to parse requests, not how business works |
| `main()` | Wiring | Everything | All concrete types — the only place that does |

#### Why This Needs Interfaces (Not Just Structs)

If `OrderService` held `*postgres.OrderRepo` (a concrete type) instead of
`OrderRepository` (an interface):

- `service` would import `postgres` → coupled to infrastructure
- Testing requires a real database → slow, flaky tests
- Swapping to MongoDB → change `service` package → change `handler` package → ripple
- No logging wrapper possible without modifying `postgres.OrderRepo`
- Feature flags for gradual migration → impossible without `if/else` in service code

The interface is the **architectural boundary**. Without it, the layers collapse
into a monolithic dependency chain.

---

### Motivation 6: Composition Through Decoration (Wrapper Pattern)

Interfaces enable **stacking behavior in layers** — the decorator pattern. Each wrapper
adds one capability and delegates everything else to the inner interface.

**The `context.Context` example:**

```go
ctx := context.Background()                       // emptyCtx — does nothing
ctx  = context.WithValue(ctx, "userID", 42)        // valueCtx wraps emptyCtx
ctx  = context.WithCancel(ctx)                     // cancelCtx wraps valueCtx
ctx  = context.WithValue(ctx, "traceID", "abc")    // valueCtx wraps cancelCtx
```

Each wrapper struct holds the parent as the **`Context` interface**, not as a concrete type:

```go
type cancelCtx struct {
    Context             // ← parent as INTERFACE
    done    chan struct{}
    err     error
}
```

It only overrides the methods it cares about (`Done()`, `Err()`). Everything else
(`Deadline()`, `Value()`) delegates to the parent via the embedded interface.

**The same pattern in `io.Reader`:**

```go
file, _ := os.Open("data.gz")        // *os.File       → io.Reader
buf := bufio.NewReader(file)           // bufio.Reader   → wraps io.Reader
gz, _ := gzip.NewReader(buf)           // gzip.Reader    → wraps io.Reader
limited := io.LimitReader(gz, 1024)    // LimitedReader  → wraps io.Reader
```

**The same pattern in `http.Handler` middleware:**

```go
handler := myHandler{}                 // implements http.Handler
handler = loggingMiddleware(handler)    // wraps http.Handler → returns http.Handler
handler = authMiddleware(handler)       // wraps http.Handler → returns http.Handler
```

**Why this requires an interface:** If the parent field were a concrete type
(e.g., `*cancelCtx`), you couldn't wrap a `valueCtx` inside a `timerCtx` inside a
`cancelCtx`. The interface makes the layers **infinitely composable** — each wrapper
doesn't know or care what's inside it.

> See [12_context_interface_deep_dive.md](./12_context_interface_deep_dive.md) for a
> complete walkthrough of how `context.Context` uses this pattern.

---

### When NOT to Use Interfaces

Equally important — interfaces have costs. Don't use them reflexively:

| Anti-Pattern | Why It's Wrong |
|-------------|----------------|
| **Preemptive interfaces** — creating an interface before you have 2+ consumers | You're guessing at the contract. Wait until a real second use case appears |
| **Fat interfaces** — 8+ methods | *"The bigger the interface, the weaker the abstraction."* Likely a concrete type in disguise |
| **Interfaces in hot paths** — tight loops, high-frequency calls | Interface method calls go through `itab` indirection — **not inlineable** by the compiler |
| **Interfaces for the sake of "clean code"** | If there's only one implementation and you don't need testability, a concrete type is simpler and faster |
| **Using `interface{}` / `any` everywhere** | You lose type safety. Generics (Go 1.18+) solve most cases where `any` was previously needed |

**The rule:** Create an interface when you need **polymorphism, testability, or
composability**. If you're reaching for an interface with one implementation that you
don't need to test in isolation — you don't need an interface.

---

### Summary Table

| # | Motivation | Core Benefit | Key Example |
|---|-----------|--------------|-------------|
| 1 | Polymorphism without inheritance | One function, many types | `io.Reader` accepted everywhere |
| 2 | Consumer-defined contracts | Dependencies point inward, packages decoupled | Service defines `UserStore`, doesn't import repo |
| 3 | Testability | Inject fakes without frameworks | 10-line fake struct in test file |
| 4 | Stdlib integration hooks | Implement a method, unlock an ecosystem | `Error()` → error handling, `String()` → printing |
| 5 | Architecture boundary abstraction | Layers separated, swappable, independently deployable | Handler → Service → Repository pattern |
| 6 | Composition through decoration | Stack behaviors in layers, infinitely composable | `context.Context`, `io.Reader` chains, HTTP middleware |

---

*Now that we understand **why** interfaces exist, the following sections explain
**how they work under the hood** at the runtime level.*

---

## 1. The Two Runtime Structs: `iface` and `eface`

Go uses **two different structs** for interface values at runtime, depending on whether
the interface has methods or not.

### `iface` — Non-Empty Interfaces (has methods)

Used for interfaces like `io.Reader`, `fmt.Stringer`, `error`, or any custom interface
with at least one method.

```
runtime.iface (16 bytes on 64-bit)
┌──────────────────────┬──────────────────────┐
│  tab  *itab          │  data unsafe.Pointer │
│  (method dispatch    │  (pointer to the     │
│   table + type info) │   concrete value)    │
└──────────────────────┴──────────────────────┘
```

- **`tab`** — pointer to an `itab` struct (see Section 2). Contains type metadata AND
  a cached array of method function pointers for this specific (interface, concrete type) pair.
- **`data`** — pointer to the actual value. If you assigned `&User{Name: "Bob"}`, this
  points to that `User` on the heap.

**Source:** `runtime/runtime2.go`
```go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}
```

### `eface` — Empty Interface (`interface{}` / `any`)

Used when the interface has **zero methods**. Since there's no method table needed,
Go uses a simpler, lighter struct.

```
runtime.eface (16 bytes on 64-bit)
┌──────────────────────┬──────────────────────┐
│  _type  *_type       │  data unsafe.Pointer │
│  (type descriptor    │  (pointer to the     │
│   only, no methods)  │   concrete value)    │
└──────────────────────┴──────────────────────┘
```

- **`_type`** — pointer to a `_type` struct (type descriptor: size, kind, hash, GC bitmap).
  No method table — just "what type is this?"
- **`data`** — same as `iface`, points to the actual value.

**Source:** `runtime/runtime2.go`
```go
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
```

### Visual Comparison

```
  Non-empty interface (e.g., io.Writer)       Empty interface (interface{} / any)
  ─────────────────────────────────────       ─────────────────────────────────
  ┌────────┬────────┐                         ┌────────┬────────┐
  │  tab ──┼──► itab│                         │ _type ─┼──► _type struct
  │        │  {     │                         │        │   {size, kind,
  │        │   inter│ ◄─ interface type info   │        │    hash, gcdata}
  │        │   _type│ ◄─ concrete type info    │        │
  │        │   fun[]│ ◄─ method pointers      │        │
  │        │  }     │                         │        │
  ├────────┼────────┤                         ├────────┼────────┤
  │ data ──┼──► concrete value on heap        │ data ──┼──► concrete value on heap
  └────────┴────────┘                         └────────┴────────┘
```

**Key insight:** Both are exactly 16 bytes (two pointers). The difference is what the
first pointer points TO — an `itab` (with methods) vs a bare `_type` (no methods).

---

## 2. The `itab` — The Method Dispatch Table

The `itab` is the heart of interface method dispatch. It answers: "For this specific
(interface type, concrete type) pair, where are the method implementations?"

```
runtime.itab
┌─────────────────────────────────────────────────────────────┐
│  inter  *interfacetype    // describes the interface        │
│         {                                                   │
│           typ    _type    // interface's own type info       │
│           pkgpath *string // package path                   │
│           mhdr   []imethod // list of methods the           │
│                            // interface requires            │
│         }                                                   │
├─────────────────────────────────────────────────────────────┤
│  _type  *_type            // describes the concrete type    │
│         {                                                   │
│           size    uintptr   // how many bytes                │
│           kind    uint8     // struct? ptr? int? etc.        │
│           hash    uint32    // for fast type comparison      │
│           str     nameOff   // type name                    │
│           ...               // GC bitmap, alignment, etc.   │
│         }                                                   │
├─────────────────────────────────────────────────────────────┤
│  hash   uint32            // copy of _type.hash for fast    │
│                           // type switch lookups            │
├─────────────────────────────────────────────────────────────┤
│  fun    [1]uintptr        // VARIABLE SIZE array of method  │
│         fun[0] = &(*EmailNotifier).Notify                   │
│         fun[1] = &(*EmailNotifier).OtherMethod              │
│         ...one entry per interface method, in sorted order  │
└─────────────────────────────────────────────────────────────┘
```

### How `itab` Gets Built and Cached

The runtime maintains a **global hash table** called `runtime.itabTable`.

```
                    itabTable (global, process-wide)
                    ────────────────────────────────
                    hash key = hash(inter, _type)

 ┌───────────────────────────────────────────────────────────────┐
 │ (io.Writer, *os.File)       → itab{ fun: [Write] }           │
 │ (io.Reader, *bytes.Buffer)  → itab{ fun: [Read] }            │
 │ (Notifier, *EmailNotifier)  → itab{ fun: [Notify] }          │
 │ (fmt.Stringer, *User)       → itab{ fun: [String] }          │
 │ ...                                                           │
 └───────────────────────────────────────────────────────────────┘
```

**Lifecycle:**

1. **First assignment** of `*EmailNotifier` to `Notifier` → runtime calls `runtime.getitab()`
2. `getitab()` checks `itabTable` for key `(Notifier, *EmailNotifier)`
3. **Cache miss** → runtime builds the `itab`:
   - Walks the interface's method list (`inter.mhdr`)
   - Walks the concrete type's method list (`_type.methods`)
   - Both are sorted alphabetically — matched with a **single O(n+m) merge walk**
   - Stores each matched method's function pointer in `fun[]`
4. **Stores** the new `itab` in `itabTable`
5. **All future assignments** of `*EmailNotifier` to `Notifier` → cache hit → O(1) lookup

**Important:** The `itab` is built **once per (interface, concrete type) pair** for the
entire lifetime of the process. It's never rebuilt or garbage collected.

---

## 3. Compile Time vs Runtime — Who Does What?

There is no `compileTime.iface`. The compiler and runtime have distinct, complementary roles:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          COMPILE TIME                                   │
│  (go build — the compiler, cmd/compile)                                │
│                                                                         │
│  1. TYPE CHECKING                                                       │
│     Verifies your concrete type has ALL methods the interface requires. │
│     If *EmailNotifier is missing Notify(), you get a compile error.     │
│     This happens at the assignment site:                                │
│       var n Notifier = &EmailNotifier{}  // compiler checks here       │
│                                                                         │
│  2. CODE GENERATION                                                     │
│     Emits assembly instructions to:                                     │
│     - Call runtime.getitab() to get/build the itab at runtime           │
│     - Construct the iface struct {tab, data}                            │
│     - For method calls: load tab → load fun[i] → indirect CALL         │
│                                                                         │
│  3. ESCAPE ANALYSIS                                                     │
│     Determines if the concrete value must escape to heap.               │
│     Assigning to an interface almost always forces a heap allocation    │
│     because iface.data is an unsafe.Pointer (must outlive the stack).   │
│                                                                         │
│  4. STATIC ITAB OPTIMIZATION (link time)                                │
│     For statically known assignments, the linker can pre-build some     │
│     itab entries — avoiding the runtime.getitab() call entirely.        │
├─────────────────────────────────────────────────────────────────────────┤
│                           RUNTIME                                       │
│  (the running binary — runtime package)                                │
│                                                                         │
│  1. ITAB TABLE                                                          │
│     Global hash table of cached itab entries.                           │
│     Built lazily on first use of each (interface, type) pair.           │
│                                                                         │
│  2. METHOD DISPATCH                                                     │
│     Indirect function call through itab.fun[i].                         │
│     CPU loads the function pointer and CALLs it.                        │
│                                                                         │
│  3. TYPE ASSERTIONS / TYPE SWITCHES                                     │
│     Compares itab._type.hash against target type hash.                  │
│     Fast path: hash match → verify → extract data.                     │
│                                                                         │
│  4. INTERFACE EQUALITY                                                  │
│     i == nil checks: tab == nil AND data == nil.                        │
│     i == j checks: same itab + deep-equal data.                        │
└─────────────────────────────────────────────────────────────────────────┘
```

**Verification tools:**
```bash
go build -gcflags='-m'     # escape analysis — what escapes to heap
go build -gcflags='-m -m'  # verbose escape analysis with reasons
go build -gcflags='-S'     # assembly output — see iface construction & dispatch
```

---

## 4. Step-by-Step: Interface Assignment

Let's trace exactly what happens when you write:

```go
type Notifier interface {
    Notify(msg string) error
}

type EmailNotifier struct{ smtp string }
func (e *EmailNotifier) Notify(msg string) error { /* ... */ }

var n Notifier = &EmailNotifier{smtp: "mx.google.com"}
```

### Step 1 — Compile Time: Type Check

```
Compiler sees: var n Notifier = &EmailNotifier{...}

  Does *EmailNotifier have Notify(string) error?
  ├─ Scan *EmailNotifier's method set
  ├─ Found: func (e *EmailNotifier) Notify(msg string) error  ✅
  └─ Type check passes — code generation proceeds
```

### Step 2 — Compile Time: Escape Analysis

```
Compiler decides: &EmailNotifier{smtp: "mx.google.com"}
  ├─ This value is assigned to an interface (iface.data = unsafe.Pointer)
  ├─ unsafe.Pointer could escape the current stack frame
  └─ Decision: ALLOCATE ON HEAP ← this is the cost of interfaces
```

### Step 3 — Runtime: Allocate the Value

```
  Heap:
  ┌──────────────────────────────┐
  │ EmailNotifier{               │
  │   smtp: "mx.google.com"     │  ◄── 0xc0000a4000 (heap address)
  │ }                            │
  └──────────────────────────────┘
```

### Step 4 — Runtime: Get or Build the `itab`

```
  runtime.getitab(Notifier, *EmailNotifier)
  ├─ Hash lookup in itabTable
  ├─ First time? Cache miss → BUILD:
  │   ├─ Walk Notifier.methods: [Notify(string) error]
  │   ├─ Walk (*EmailNotifier).methods: [Notify(string) error]
  │   ├─ Match: fun[0] = address of (*EmailNotifier).Notify
  │   └─ Store in itabTable
  └─ Return: *itab at 0x4a7b20
```

### Step 5 — Runtime: Assemble the `iface`

```
  n (on the stack):
  ┌──────────────────────┬──────────────────────┐
  │ tab = 0x4a7b20       │ data = 0xc0000a4000  │
  │ (itab for            │ (heap-allocated      │
  │  Notifier,           │  EmailNotifier)      │
  │  *EmailNotifier)     │                      │
  └──────────────────────┴──────────────────────┘
            │                       │
            ▼                       ▼
      itab{                   EmailNotifier{
        inter: Notifier         smtp: "mx.google.com"
        _type: *EmailNotifier }
        fun[0]: 0x4823a0  ──► (*EmailNotifier).Notify
      }
```

---

## 5. Step-by-Step: Interface Method Call

Now let's trace `n.Notify("server down")`:

```go
n.Notify("server down")
```

### Step 1 — Load the `itab`

```
  CPU reads n.tab → 0x4a7b20 (pointer to itab)
```

### Step 2 — Load the Method Pointer

```
  CPU reads itab.fun[0] → 0x4823a0 (address of (*EmailNotifier).Notify)

  fun[] index is determined at COMPILE TIME — the compiler knows Notify is
  method #0 in the interface's sorted method list. This is a constant offset.
```

### Step 3 — Load the Data Pointer

```
  CPU reads n.data → 0xc0000a4000 (pointer to EmailNotifier value)
```

### Step 4 — Indirect Call

```
  CALL 0x4823a0(0xc0000a4000, "server down")
       ▲              ▲              ▲
       │              │              └─ argument
       │              └─ receiver (e *EmailNotifier)
       └─ (*EmailNotifier).Notify function address

  This is an INDIRECT call — the CPU doesn't know the target until it
  reads the function pointer from the itab. This means:
  ├─ The compiler CANNOT inline this call
  ├─ The CPU branch predictor may mispredict on first call
  └─ Subsequent calls to the same (interface, type) pair are predicted correctly
```

### Visual Summary: The Full Call Chain

```
  n.Notify("server down")

  ┌─── n ────────────┐
  │ tab ──► itab     │     ┌─ EmailNotifier ─┐
  │ data ──┼─────────┼────►│ smtp: "mx..."   │
  └────────┼─────────┘     └─────────────────┘
           │                       ▲
           ▼                       │
     itab.fun[0] ──► (*EmailNotifier).Notify(receiver, "server down")
                                   │              ▲
                                   └──────────────┘
                                   receiver = n.data
```

---

## 6. The Three States of an Interface

Every interface value is in one of exactly three states:

```go
var s Stringer         // State 1: (nil, nil)
s = User{Name: "Bob"}  // State 2: (*itab, 0xc000...)
var u *User
s = u                  // State 3: (*itab, nil)   ← THE TRAP
```

### State 1 — True Nil Interface

```
var s Stringer

s:
┌──────────────┬──────────────┐
│     nil      │     nil      │
└──────────────┴──────────────┘

s == nil → TRUE ✅   (both words are zero)
```

Both `tab` and `data` are zero. This is what `nil` means for an interface.

### State 2 — Normal, Non-Nil Interface

```
s = User{Name: "Bob", Age: 25}

s:
┌──────────────┬──────────────┐
│  *itab(User) │  0xc000a080  │──► User{Name:"Bob", Age:25}
└──────────────┴──────────────┘

s == nil  → FALSE ✅
s.String() → works fine ✅
```

### State 3 — The Nil Pointer Trap ⚠️

```go
var u *User   // u is nil
s = u         // assign nil pointer to interface
```

```
s:
┌──────────────┬──────────────┐
│  *itab(User) │     nil      │
└──────────────┴──────────────┘

s == nil    → FALSE ❌  (tab is NOT nil — it has type info!)
s.String()  → PANIC ❌  (data IS nil — method dereferences nil receiver)
```

---

## 7. The Nil Interface Trap

### Why Does This Happen?

When Go assigns a value to an interface, it **always** writes the type into the `tab` field —
**regardless of whether the value itself is nil**.

```go
var u *User   // concrete nil pointer — the type (*User) IS known
s = u         // Go writes: tab = itab(*User), data = nil
              // it has no way to "forget" the type just because the value is nil
```

Go cannot collapse this back to `(nil, nil)` because the type information was known
statically at the assignment site. The compiler stamped `*User` into the `itab` and moved on.

### How `== nil` Works on Interfaces

The `==` operator on an interface checks:

```
(tab == nil) AND (data == nil)
```

Only when **both are zero** is the interface considered nil. State 3 has a non-nil `tab`,
so `== nil` is `false` — even though the data pointer inside is nil.

### The Classic Bug

```go
// BAD
func findUser(found bool) Stringer {
    var u *User              // type: *User, value: nil
    if !found {
        return u             // returns iface{tab: *itab, data: nil} — NOT nil!
    }
    return &User{Name: "Alice"}
}

result := findUser(false)
if result == nil {           // FALSE — the nil check is BYPASSED
    fmt.Println("not found")
}
result.String()              // PANIC — data is nil, method dereferences it
```

### What the Panic Trace Looks Like

```
n.Notify("hello")                    // interface method call
  ↓
runtime loads n.tab.fun[0]           → address of (*EmailNotifier).Notify ← SUCCEEDS
  ↓
calls (*EmailNotifier).Notify(nil, "hello")   ← receiver is nil
  ↓
inside Notify: e.smtpClient.Send()   → dereferences nil receiver 'e'
  ↓
PANIC: nil pointer dereference        ← error surfaces DEEP INSIDE the method

Stack trace points at EmailNotifier.Notify, NOT at where you assigned the typed nil.
This makes production debugging harder — the root cause is far from the crash site.
```

---

## 8. The Three Guards Against the Nil Trap

### Guard 1 — Fix at the Source (always prefer this)

**Never return a typed nil variable through an interface — return bare `nil` directly.**

```go
// BAD — typed nil leaks in
func findUser(found bool) Stringer {
    var u *User
    if !found {
        return u   // iface{*itab, nil} — bug!
    }
    return &User{Name: "Alice"}
}

// GOOD — return untyped nil
func findUser(found bool) Stringer {
    if !found {
        return nil  // iface{nil, nil} — true nil, == nil check works
    }
    return &User{Name: "Alice"}
}
```

Bare `return nil` in a function returning an interface produces `iface{nil, nil}`.

**When to use:** Every time you write a function returning an interface type. Make this a habit.

### Guard 2 — Type Assertion (when you know the concrete type)

```go
func safeCall(s Stringer) {
    u, ok := s.(*User)   // step 1: extract the concrete *User
    if !ok || u == nil {  // step 2: check type match AND nil pointer
        fmt.Println("skipping — nil or wrong type")
        return
    }
    fmt.Println(u.String()) // safe — u is a non-nil *User
}
```

Under the hood for `s.(*User)` when `s = iface{tab: *itab(*User), data: nil}`:

```
1. Compare s.tab._type.hash == hash(*User) → YES → ok = true
2. Extract s.data → nil → u = nil
∴ ok is true (type matched) but u is nil (pointer is nil)
  You must check BOTH.
```

**When to use:** When you know the expected concrete type (e.g., inside a service you control).

### Guard 3 — Reflect (last resort, for truly dynamic types)

```go
func isTrulyNil(i any) bool {
    if i == nil {
        return true  // fast path: iface{nil, nil}
    }
    v := reflect.ValueOf(i)
    switch v.Kind() {
    case reflect.Ptr, reflect.Interface,
         reflect.Slice, reflect.Map,
         reflect.Chan, reflect.Func:
        return v.IsNil()  // inspects the data word directly
    }
    return false  // int, string, struct etc. can never be nil
}
```

**When to use:** Generic libraries receiving `any` where you can't know the type. Avoid in hot paths
(reflection allocates and is ~10x slower than a type assertion).

### Decision Tree

```
You have an interface value and need to know if it's safe to call methods
│
├─ Can you change the code that PRODUCES the interface value?
│   └─ YES → Guard 1: return nil, not a typed nil variable
│
└─ NO, you receive it from outside
    │
    ├─ Do you know the concrete type?
    │   └─ YES → Guard 2: v, ok := i.(*Type); ok && v != nil
    │
    └─ NO, it's truly dynamic / any
        └─ Guard 3: reflect.ValueOf(i).IsNil()
```

---

## 9. Type Assertions Under the Hood

Given any interface variable `i` (e.g., `var i Notifier = &EmailNotifier{}`):

### Single Type Assertion: `v, ok := i.(ConcreteType)`

```
// i is any interface value — could be Notifier, io.Reader, error, any, etc.
// The assertion extracts the concrete type stored inside.

var i Notifier = &EmailNotifier{smtp: "mx.google.com"}
v, ok := i.(*EmailNotifier)   // is the concrete type inside i a *EmailNotifier?

Under the hood:
1. Read i.tab._type.hash                    // O(1) — hash stored in itab
2. Compare with hash(*EmailNotifier)        // O(1) — compiler knows this at build time
3. Match? → ok = true, v = *(*EmailNotifier)(i.data)   // cast the data pointer
   No match? → ok = false, v = nil (*EmailNotifier zero value)
```

### Type Switch: `switch v := i.(type) { ... }`

```go
// i is the interface variable you're inspecting.
// The switch tests which concrete type is stored inside.
var i any = 42

switch v := i.(type) {
case int:       // compare i._type.hash with hash(int)     ← eface, not itab (empty interface)
case string:    // compare i._type.hash with hash(string)
case *User:     // compare i._type.hash with hash(*User)
default:        // no match
}
```

The compiler generates a sequence of hash comparisons (or a jump table for many cases).
The hash is copied into `itab.hash` specifically to make this fast — no pointer chase needed.

---

## 10. Performance Implications

### Interface Method Calls Cannot Be Inlined

```go
func directCall(e *EmailNotifier) { e.Notify("hi") }  // ← CAN be inlined
func ifaceCall(n Notifier)        { n.Notify("hi") }  // ← CANNOT be inlined
```

Why? Inlining requires knowing the target function at compile time. With interfaces,
the target lives in `itab.fun[0]` which is only known at runtime.

**Production impact:** In hot paths processing 100k+ operations/sec, interface dispatch
overhead adds up. Consider using concrete types in inner loops and interfaces at boundaries.

### Interface Assignment Forces Heap Allocation

```go
var n Notifier = &EmailNotifier{smtp: "mx.google.com"}
//                ↑ this value ESCAPES to heap
//                  because iface.data is unsafe.Pointer
```

Every value assigned to an interface must live on the heap (the `data` pointer must remain
valid). This creates GC pressure. Verify with:

```bash
go build -gcflags='-m' ./...
# output: "EmailNotifier{} escapes to heap"
```

### Cost Summary

```
┌─────────────────────┬────────────────────┬─────────────────────────────────┐
│ Operation           │ Cost               │ Why                             │
├─────────────────────┼────────────────────┼─────────────────────────────────┤
│ Interface assign    │ ~30-50ns           │ itab lookup + possible heap     │
│ (first time)        │                    │ alloc + GC bookkeeping          │
├─────────────────────┼────────────────────┼─────────────────────────────────┤
│ Interface assign    │ ~5-15ns            │ itab cached, just heap alloc    │
│ (cached itab)       │                    │                                 │
├─────────────────────┼────────────────────┼─────────────────────────────────┤
│ Interface method    │ ~2-5ns overhead    │ Indirect call (load fun ptr +   │
│ call                │ vs direct call     │ CALL) — prevents inlining       │
├─────────────────────┼────────────────────┼─────────────────────────────────┤
│ Type assertion      │ ~1-2ns             │ Hash comparison + pointer cast  │
├─────────────────────┼────────────────────┼─────────────────────────────────┤
│ == nil check        │ ~0.5ns             │ Two pointer comparisons         │
└─────────────────────┴────────────────────┴─────────────────────────────────┘
```

---

## 11. Quick Reference Card

```
NON-EMPTY INTERFACE (has methods)
─────────────────────────────────
runtime.iface { tab *itab, data unsafe.Pointer }
  └─ itab { inter *interfacetype, _type *_type, hash uint32, fun [N]uintptr }
     └─ fun[] = cached method pointers for (interface, concrete type) pair
     └─ cached globally in runtime.itabTable — built once, reused forever

EMPTY INTERFACE (interface{} / any)
───────────────────────────────────
runtime.eface { _type *_type, data unsafe.Pointer }
  └─ no method table, just type descriptor + data pointer
  └─ simpler because no methods to dispatch

NIL SEMANTICS
─────────────
  iface{nil, nil}    == nil  → TRUE   ← "true nil" interface
  iface{*itab, nil}  == nil  → FALSE  ← "typed nil" TRAP
  iface{*itab, ptr}  == nil  → FALSE  ← normal non-nil

COMPILE TIME vs RUNTIME
───────────────────────
  Compile: type checking, code generation, escape analysis, static itab optimization
  Runtime: itab caching, method dispatch (indirect call), type assertions, nil checks

TOOLS
─────
  go build -gcflags='-m'      # escape analysis: what goes to heap
  go build -gcflags='-m -m'   # verbose escape analysis with reasons
  go build -gcflags='-S'      # assembly output: see iface construction
  go test -race ./...         # race detector: always run
```

---

## One-Line Summary

> An interface is a two-word struct `{tab, data}` — `tab` points to the `itab`
> (method dispatch table cached per type pair), `data` points to the value.
> It's `nil` only when **both** words are zero. The compiler verifies types;
> the runtime dispatches methods. This is why interface calls can't be inlined.

