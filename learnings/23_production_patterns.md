# 23 — Production Patterns: The Five Pillars of Go Architecture

> Five patterns you'll find in every serious Go codebase. These aren't
> theoretical — they're the actual architecture of production services
> at Google, Uber, Cloudflare, Stripe, and every well-run Go shop.

---

## Table of Contents

1. [Middleware — Function Composition](#1-middleware--function-composition)
2. [Functional Options — Self-Documenting Configuration](#2-functional-options--self-documenting-configuration)
3. [Graceful Shutdown — Production Lifecycle Management](#3-graceful-shutdown--production-lifecycle-management)
4. [Interface-Based Dependency Injection — The Go Way](#4-interface-based-dependency-injection--the-go-way)
5. [Table-Driven Tests — The Only Way to Test in Go](#5-table-driven-tests--the-only-way-to-test-in-go)
6. [How These Five Connect](#6-how-these-five-connect)

---

## 1. Middleware — Function Composition

Covered in depth in [Chapter 22](./22_middleware_pattern.md).

Quick recap — the signature that powers all Go HTTP frameworks:

```go
type Middleware func(http.Handler) http.Handler
```

Three features make it work: methods on function types (`HandlerFunc`),
structural typing (implicit interface satisfaction), closures (capture `next`).

---

## 2. Functional Options — Self-Documenting Configuration

### The Problem: Configuring Complex Structs

Go has no default parameter values, no named arguments, no method overloading.
How do you create flexible, readable constructors?

```go
// ❌ Attempt 1: Giant constructor — unreadable, fragile
func NewServer(addr string, port int, timeout time.Duration,
    maxConns int, logger *slog.Logger, tls bool) *Server

s := NewServer("localhost", 8080, 30*time.Second, 100, logger, true)
//                                                 ↑ what is 100?

// ❌ Attempt 2: Config struct — zero-value ambiguity
type Config struct {
    Timeout  time.Duration
    MaxConns int
}
// Did the user set MaxConns=0 intentionally or forget to set it?
// You can't tell. Both are the zero value.

// ❌ Attempt 3: Pointer fields to distinguish "not set" from "zero"
type Config struct {
    Timeout  *time.Duration  // nil = not set, &0 = set to zero
    MaxConns *int            // ugly, error-prone, annoying to use
}
```

### The Solution: Functional Options Pattern

```go
// Step 1: Define a function type for options
type Option func(*Server)

// Step 2: Create named option constructors
func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d        // closure captures 'd'
    }
}

func WithMaxConns(n int) Option {
    return func(s *Server) {
        s.maxConns = n       // closure captures 'n'
    }
}

func WithLogger(l *slog.Logger) Option {
    return func(s *Server) {
        s.logger = l         // closure captures 'l'
    }
}

func WithTLS(certFile, keyFile string) Option {
    return func(s *Server) {
        s.tls = true
        s.certFile = certFile
        s.keyFile = keyFile
    }
}
```

```go
// Step 3: Constructor with sensible defaults + variadic options
func NewServer(addr string, opts ...Option) *Server {
    // Start with sensible defaults
    s := &Server{
        addr:     addr,
        timeout:  30 * time.Second,
        maxConns: 100,
        logger:   slog.Default(),
    }

    // Apply each option — overrides only what the user specified
    for _, opt := range opts {
        opt(s)
    }

    return s
}
```

```go
// Usage — clean, readable, self-documenting:

// Override only what you need:
s := NewServer("localhost:8080",
    WithTimeout(60 * time.Second),
    WithMaxConns(500),
)

// Use all defaults (perfectly valid):
s := NewServer("localhost:8080")

// Full configuration:
s := NewServer("localhost:8080",
    WithTimeout(60 * time.Second),
    WithMaxConns(1000),
    WithLogger(customLogger),
    WithTLS("cert.pem", "key.pem"),
)
```

### Under the Hood — What the Compiler Sees

```
  NewServer("addr", WithTimeout(60s), WithMaxConns(500))

  Step 1: WithTimeout(60s) executes:
  ┌──────────────────────────────────────────────┐
  │ Returns a closure:                            │
  │   funcval {                                   │
  │     fn: anon_func_code_ptr                   │
  │     captured: d = 60s                        │  ← closure captures the value
  │   }                                           │
  └──────────────────────────────────────────────┘

  Step 2: WithMaxConns(500) executes:
  ┌──────────────────────────────────────────────┐
  │ Returns a closure:                            │
  │   funcval {                                   │
  │     fn: anon_func_code_ptr                   │
  │     captured: n = 500                        │
  │   }                                           │
  └──────────────────────────────────────────────┘

  Step 3: NewServer receives opts = []Option{closure1, closure2}

  Step 4: Loop applies each:
    opts[0](s)  → s.timeout = 60s     (closure1 runs, uses captured d)
    opts[1](s)  → s.maxConns = 500    (closure2 runs, uses captured n)
```

### Why Not a Config Struct?

```
  Config struct:                        Functional Options:
  ──────────────                        ────────────────────
  ✅ Familiar to Java/C# devs          ✅ Self-documenting API
  ✅ Easy to serialize (JSON/YAML)      ✅ No zero-value ambiguity
  ❌ Zero-value ambiguity               ✅ Extensible without breaking API
  ❌ All-or-nothing (must know all)     ✅ Only set what you need
  ❌ Adding field may break callers     ✅ Adding option never breaks callers
  ❌ Validation happens after creation  ✅ Can validate during construction
```

Adding a new option to the config struct requires checking every call site.
Adding a new `WithXxx()` function is purely additive — existing code doesn't change.

### Where You'll See This

```
  Library                          Example
  ───────                          ───────
  gRPC                             grpc.NewServer(grpc.MaxRecvMsgSize(1024))
  zap                              zap.NewProduction(zap.AddCaller())
  chi                              chi.NewRouter() + r.Use(middleware.WithValue(...))
  Google Cloud client               storage.NewClient(ctx, option.WithCredentialsFile("..."))
  Kubernetes client-go              rest.InClusterConfig() with overrides
  Your own services                 NewOrderService(WithDB(db), WithCache(redis))
```

### Advanced: Option with Validation

```go
func WithMaxConns(n int) Option {
    return func(s *Server) {
        if n <= 0 {
            n = 1  // enforce minimum
        }
        if n > 10000 {
            n = 10000  // enforce maximum
        }
        s.maxConns = n
    }
}
```

### Advanced: Option That Returns Error

Some libraries use an error-returning option for validation:

```go
type Option func(*Server) error

func NewServer(addr string, opts ...Option) (*Server, error) {
    s := &Server{addr: addr, timeout: 30 * time.Second}
    for _, opt := range opts {
        if err := opt(s); err != nil {
            return nil, fmt.Errorf("server option: %w", err)
        }
    }
    return s, nil
}

func WithPort(p int) Option {
    return func(s *Server) error {
        if p < 1 || p > 65535 {
            return fmt.Errorf("invalid port: %d", p)
        }
        s.port = p
        return nil
    }
}
```

---

## 3. Graceful Shutdown — Production Lifecycle Management

Every production Go service needs to handle shutdown correctly. In Kubernetes
and Cloud Run, your process receives SIGTERM and has ~30 seconds to clean up
before SIGKILL. If you don't handle it, in-flight requests get dropped.

### The Problem

```
  Without graceful shutdown:

  1. Kubernetes sends SIGTERM to your pod
  2. Go program exits immediately (os.Exit or main() returns)
  3. In-flight HTTP requests → connection reset errors
  4. Database transactions → left in unknown state
  5. Pub/Sub messages → not ACKed, will be redelivered (duplicate processing)
  6. Customer sees: "500 Internal Server Error" ← during a normal deployment!
```

### The Solution: Signal + Context + Ordered Drain

```go
func main() {
    // Phase 1: Create a context that cancels on SIGTERM/SIGINT
    ctx, stop := signal.NotifyContext(context.Background(),
        syscall.SIGTERM, syscall.SIGINT,
    )
    defer stop()

    // Phase 2: Initialize everything
    db := initDB()
    defer db.Close()

    srv := &http.Server{
        Addr:    ":8080",
        Handler: buildRouter(db),
    }

    // Phase 3: Start server in a goroutine
    go func() {
        slog.Info("server starting", "addr", srv.Addr)
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            slog.Error("server error", "err", err)
        }
    }()

    // Phase 4: Wait for shutdown signal
    <-ctx.Done()
    slog.Info("shutdown signal received")

    // Phase 5: Graceful shutdown with timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    // Stop accepting new requests, wait for in-flight to complete
    if err := srv.Shutdown(shutdownCtx); err != nil {
        slog.Error("shutdown error", "err", err)
    }

    // Phase 6: Close resources in reverse order
    // (DB connections, message queues, flush logs)
    db.Close()
    slog.Info("server stopped cleanly")
}
```

### The Shutdown Timeline

```
  Time 0s:   SIGTERM received
             ├── ctx.Done() fires
             ├── srv.Shutdown() called:
             │   ├── Close listener (no new connections)
             │   ├── Wait for in-flight requests to complete
             │   └── Returns when all done (or timeout)
             │
  Time 1-5s: In-flight requests finish naturally
             ├── Each request completes, response sent
             └── Connection count drops to 0
             │
  Time 5s:   srv.Shutdown() returns
             ├── db.Close() — close database connections
             ├── Flush logs
             └── main() returns — process exits with code 0
             │
  Time 30s:  (If still running) Kubernetes sends SIGKILL
             └── Process killed forcefully — avoid this!
```

### What `http.Server.Shutdown()` Does Under the Hood

```go
  func (srv *Server) Shutdown(ctx context.Context) error {
      // 1. Close the listener — no new connections accepted
      srv.listener.Close()

      // 2. Set a flag — new requests on existing connections get 503
      srv.inShutdown.Store(true)

      // 3. Wait for active connections to finish
      //    OR ctx deadline/cancel
      for {
          if srv.activeConns.Load() == 0 {
              return nil  // all done!
          }
          select {
          case <-ctx.Done():
              return ctx.Err()  // timeout — some requests didn't finish
          case <-time.After(500 * time.Millisecond):
              // poll again
          }
      }
  }
```

Notice: this uses `context` for the shutdown timeout — the same pattern
we discussed with `ctx.Done()` in select. It's channels all the way down.

### GCP Cloud Run Specifics

```
  Cloud Run lifecycle:
  1. Instance receives SIGTERM
  2. Default grace period: 10 seconds (configurable up to 3600s)
  3. SIGKILL after grace period

  Best practice:
  ├── Set Cloud Run timeout > your shutdown timeout
  │   e.g., Cloud Run: 30s, your code: 15s
  ├── Use signal.NotifyContext (not os.Signal channel)
  ├── Drain Pub/Sub consumers before closing DB
  └── Log "shutdown complete" so you can verify in Cloud Logging
```

### The Complete Shutdown Order

```
  Shutdown resources in REVERSE dependency order:

  Startup:                          Shutdown:
  1. Config                         6. Close DB connections
  2. Database                       5. Close message consumers
  3. Cache (Redis)                  4. Close cache connections
  4. Message consumers (Pub/Sub)    3. Flush metrics/logs
  5. HTTP server                    2. HTTP server.Shutdown()
  6. Start accepting requests       1. Stop accepting requests

  Rule: last started, first stopped.
  Why: HTTP server depends on DB — don't close DB while requests use it!
```

---

## 4. Interface-Based Dependency Injection — The Go Way

### The Problem: How to Wire Services Together

In Java, you'd use Spring's `@Autowired`. In C#, you'd use a DI container.
In Go? **Constructor injection. That's it.**

```go
// The "Go Way" — no framework, no container, no annotations

// Step 1: Define interfaces at the CONSUMER (not the provider)
type OrderRepository interface {
    Create(ctx context.Context, order *Order) error
    GetByID(ctx context.Context, id string) (*Order, error)
}

// Step 2: Implement the interface
type postgresOrderRepo struct {
    db *sql.DB
}

func NewPostgresOrderRepo(db *sql.DB) *postgresOrderRepo {
    return &postgresOrderRepo{db: db}
}

func (r *postgresOrderRepo) Create(ctx context.Context, order *Order) error {
    _, err := r.db.ExecContext(ctx,
        "INSERT INTO orders (id, customer_id, total) VALUES ($1, $2, $3)",
        order.ID, order.CustomerID, order.Total,
    )
    return err
}

// Step 3: Service ACCEPTS the interface (not the concrete type)
type OrderService struct {
    repo   OrderRepository   // interface, not *postgresOrderRepo
    logger *slog.Logger
}

func NewOrderService(repo OrderRepository, logger *slog.Logger) *OrderService {
    return &OrderService{repo: repo, logger: logger}
}

// Step 4: Wire everything together in main()
func main() {
    db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    repo := NewPostgresOrderRepo(db)         // concrete type
    svc := NewOrderService(repo, slog.Default()) // injected as interface
    handler := NewOrderHandler(svc)
    // ...
}
```

### Why No DI Framework?

```
  Java Spring:                        Go:
  ──────────                          ───
  @Service                            (no annotation)
  @Autowired                          (no annotation)
  @Qualifier("primary")               (no annotation)
  ComponentScan("com.app")            (no scanning)
  ApplicationContext                   main() function

  Spring's DI container:              Go's DI container:
  ├── Reflection-based                ├── YOUR main() function
  ├── Runtime wiring                  ├── Compile-time wiring
  ├── Errors at startup               ├── Errors at COMPILE time
  ├── 500+ classes in framework       ├── 0 framework classes
  └── "Magic" — hard to trace         └── Explicit — easy to trace
```

Go developers reject DI containers because:
1. **Compile-time safety** — if you forget to wire something, it won't compile
2. **Readability** — `main()` shows the entire object graph, no hidden magic
3. **Debuggability** — stack traces go through your code, not framework internals

### The Key Principle: Accept Interfaces, Return Structs

```go
// ✅ Accept interface — flexible, testable
func NewOrderService(repo OrderRepository) *OrderService

// ❌ Accept concrete type — coupled, hard to test
func NewOrderService(repo *PostgresOrderRepo) *OrderService

// ✅ Return concrete type — caller gets full type information
func NewPostgresOrderRepo(db *sql.DB) *postgresOrderRepo

// ❌ Return interface — hides useful methods, reduces type info
func NewPostgresOrderRepo(db *sql.DB) OrderRepository
```

Why return structs? Because Go interfaces are satisfied implicitly — the caller
can assign your `*postgresOrderRepo` to any interface it satisfies. Returning
an interface **restricts** what the caller can do. Return the concrete type and
let the caller decide which interface to use.

### Testing With Interfaces — No Mock Framework Needed

```go
// In test: create a simple fake that satisfies the interface
type fakeOrderRepo struct {
    orders map[string]*Order
}

func (f *fakeOrderRepo) Create(ctx context.Context, order *Order) error {
    f.orders[order.ID] = order
    return nil
}

func (f *fakeOrderRepo) GetByID(ctx context.Context, id string) (*Order, error) {
    o, ok := f.orders[id]
    if !ok {
        return nil, ErrNotFound
    }
    return o, nil
}

func TestOrderService_Create(t *testing.T) {
    repo := &fakeOrderRepo{orders: make(map[string]*Order)}
    svc := NewOrderService(repo, slog.Default())

    err := svc.CreateOrder(ctx, &Order{ID: "123", Total: 99.99})
    if err != nil {
        t.Fatal(err)
    }

    // Verify via the fake
    if repo.orders["123"].Total != 99.99 {
        t.Error("order not stored correctly")
    }
}
```

No `gomock`. No `testify/mock`. No code generation. Just a struct that
implements the interface. This is why Go interfaces are small — a 2-method
interface means a 2-method fake. A 20-method interface means nobody wants
to write the fake → **design pressure toward small interfaces**.

### Where Interfaces Are Defined — The Consumer Rule

```
  ❌ WRONG: define interface next to the implementation
  ├── repository/
  │   ├── order_repo.go          ← interface + implementation together
  │   └── postgres_order_repo.go

  ✅ RIGHT: define interface at the consumer
  ├── service/
  │   └── order_service.go       ← interface defined HERE (consumer)
  ├── repository/
  │   └── postgres_order_repo.go ← implementation, doesn't know about interface
```

The consumer defines what it needs. The implementation doesn't even import
the interface package. This is structural typing in action — the implementation
satisfies the interface without knowing it exists.

---

## 5. Table-Driven Tests — The Only Way to Test in Go

### The Pattern

```go
func TestParseStatus(t *testing.T) {
    tests := []struct {
        name   string
        input  string
        want   OrderStatus
        wantErr bool
    }{
        {name: "valid created", input: "created", want: StatusCreated},
        {name: "valid shipped", input: "shipped", want: StatusShipped},
        {name: "case insensitive", input: "CREATED", want: StatusCreated},
        {name: "empty string", input: "", wantErr: true},
        {name: "invalid status", input: "flying", wantErr: true},
        {name: "whitespace", input: " shipped ", want: StatusShipped},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseStatus(tt.input)
            if tt.wantErr {
                if err == nil {
                    t.Errorf("ParseStatus(%q) expected error, got %v", tt.input, got)
                }
                return
            }
            if err != nil {
                t.Errorf("ParseStatus(%q) unexpected error: %v", tt.input, err)
                return
            }
            if got != tt.want {
                t.Errorf("ParseStatus(%q) = %v, want %v", tt.input, got, tt.want)
            }
        })
    }
}
```

### Why This Pattern Dominates Go

```
  1. Adding a test case = adding ONE line to the table
     (no new function, no new file, no copy-paste)

  2. t.Run(tt.name, ...) creates SUBTESTS:
     → go test -run TestParseStatus/empty_string  ← run just one case
     → go test -v  ← see every case name in output

  3. t.Parallel() in the subtest:
     → all cases run concurrently (free speedup)

  4. The struct fields ARE the documentation:
     → name, input, want, wantErr — self-explanatory

  5. Failure messages include input:
     → "ParseStatus("flying") = ..." not just "got X, want Y"
```

### Structure for Complex Tests

```go
func TestCreateOrder(t *testing.T) {
    tests := []struct {
        name      string
        order     Order
        setupRepo func(*fakeOrderRepo)    // customize fake per test
        wantErr   error                    // specific error to check
    }{
        {
            name:  "success",
            order: Order{ID: "1", Total: 99.99},
            wantErr: nil,
        },
        {
            name:  "duplicate order",
            order: Order{ID: "1", Total: 99.99},
            setupRepo: func(r *fakeOrderRepo) {
                r.orders["1"] = &Order{ID: "1"}  // pre-populate
            },
            wantErr: ErrDuplicateOrder,
        },
        {
            name:  "zero total rejected",
            order: Order{ID: "2", Total: 0},
            wantErr: ErrInvalidTotal,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &fakeOrderRepo{orders: make(map[string]*Order)}
            if tt.setupRepo != nil {
                tt.setupRepo(repo)
            }
            svc := NewOrderService(repo, slog.Default())

            err := svc.CreateOrder(context.Background(), &tt.order)

            if !errors.Is(err, tt.wantErr) {
                t.Errorf("CreateOrder() error = %v, want %v", err, tt.wantErr)
            }
        })
    }
}
```

### Anti-Patterns to Avoid

```
  ❌ One assertion per test function (JUnit style):
     func TestParseStatus_Created(t *testing.T) { ... }
     func TestParseStatus_Shipped(t *testing.T) { ... }
     func TestParseStatus_Empty(t *testing.T) { ... }
     → 50 functions for 50 cases. Unmaintainable.

  ❌ Shared mutable state across test cases:
     repo := &fakeRepo{}
     for _, tt := range tests {
         // each test mutates repo → order-dependent, flaky
     }
     → Create fresh state inside the loop.

  ❌ Skipping t.Run (no subtests):
     for _, tt := range tests {
         got := Parse(tt.input)  // which case failed? Good luck.
     }
     → Always use t.Run(tt.name, ...) for identifiable failures.

  ❌ Testing implementation, not behavior:
     "function called repo.Create with these exact args"
     → Test the OUTPUT, not the internal calls.
```

---

## 6. How These Five Connect

```
  These five patterns form a complete architecture:

  ┌─────────────────────────────────────────────────────────────────┐
  │                     main()                                      │
  │                                                                 │
  │  // Functional Options → configure each component               │
  │  db := NewDB(WithConnPool(50), WithTimeout(5*time.Second))     │
  │  repo := NewPostgresRepo(db)                                    │
  │                                                                 │
  │  // Interface-Based DI → wire services via interfaces           │
  │  svc := NewOrderService(repo, WithLogger(logger))               │
  │  handler := NewOrderHandler(svc)                                │
  │                                                                 │
  │  // Middleware → wrap handler with cross-cutting concerns        │
  │  router := chi.NewRouter()                                      │
  │  router.Use(LoggingMiddleware, AuthMiddleware, RecoverMiddleware)│
  │  router.Mount("/orders", handler.Routes())                      │
  │                                                                 │
  │  // Graceful Shutdown → manage the lifecycle                    │
  │  srv := &http.Server{Handler: router}                           │
  │  go srv.ListenAndServe()                                        │
  │  <-ctx.Done()                                                   │
  │  srv.Shutdown(shutdownCtx)                                      │
  │  db.Close()                                                     │
  │                                                                 │
  │  // Table-Driven Tests → verify every component                 │
  │  // (in *_test.go files with fakes injected via interfaces)     │
  └─────────────────────────────────────────────────────────────────┘

  Functional Options  → how you BUILD components
  Interface-Based DI  → how you CONNECT components
  Middleware          → how you WRAP components with cross-cutting behavior
  Graceful Shutdown   → how you START and STOP the system
  Table-Driven Tests  → how you VERIFY everything works
```

> **Go Wisdom**: *"Clear is better than clever."*
> None of these patterns require a framework. They're just functions, interfaces,
> and structs composed together. The entire architecture is visible in `main()`.
> No hidden magic, no annotation processing, no runtime reflection.
> That's the Go way.
