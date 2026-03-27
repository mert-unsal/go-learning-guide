# 22 — The Middleware Pattern: Function Types, Closures & Interface Composition

> How Go's three unique features — methods on function types, structural typing,
> and closures — combine to create the HTTP middleware ecosystem with zero frameworks.
> This chapter connects Chapters 06 (closures), 08 (interfaces), and 21 (design philosophy).

---

## Table of Contents

1. [Plain Function vs Function Type](#1-plain-function-vs-function-type)
2. [Methods on Function Types — Go's Secret Weapon](#2-methods-on-function-types--gos-secret-weapon)
3. [The Middleware Pattern — Under the Hood](#3-the-middleware-pattern--under-the-hood)
4. [How Closures Work in Middleware](#4-how-closures-work-in-middleware)
5. [Without These Features — The Java Way](#5-without-these-features--the-java-way)
6. [ResponseWriter Wrapping — Struct Embedding in Action](#6-responsewriter-wrapping--struct-embedding-in-action)
7. [Closures Across Languages — Go vs Java vs JavaScript](#7-closures-across-languages--go-vs-java-vs-javascript)
8. [Production Middleware Stack](#8-production-middleware-stack)
9. [The Three Features That Make It All Work](#9-the-three-features-that-make-it-all-work)

---

## 1. Plain Function vs Function Type

This distinction is fundamental and often overlooked.

### Plain Function

A plain function is just a named block of code:

```go
func greet(name string) string {
    return "hello " + name
}
```

Under the hood, `greet` is a fixed symbol in the binary. The compiler knows its
address at compile time. You can call it, but its **type** is implicit.

```
  Binary:
  ┌──────────────────────────────────────┐
  │  .text section                       │
  │                                      │
  │  greet:                              │
  │    0x4a2100: MOVQ ...               │  ← fixed address in binary
  │    0x4a2108: CALL runtime.concatstr │
  │    0x4a2110: RET                    │
  └──────────────────────────────────────┘
```

### Function Type

A **function type** is a named type whose underlying type is a function signature:

```go
type Greeter func(name string) string
```

Now `Greeter` is a first-class type — just like `int` or `string`. You can:
- Declare variables of type `Greeter`
- Pass `Greeter` as a function parameter
- Return `Greeter` from a function
- **Define methods on `Greeter`** ← this is the key

```go
var g Greeter = greet   // assign plain function to function type variable
g("world")              // calls greet("world")
```

### Under the Hood — The funcval

When you assign a function to a variable, Go creates a `funcval` on the heap:

```
  g := greet

  Stack:                          Heap:
  ┌────────────────┐             ┌──────────────────────┐
  │ g (8 bytes)    │             │ funcval              │
  │ ptr ───────────┼────────────→│ fn: 0x4a2100 (greet) │
  └────────────────┘             └──────────────────────┘

  Variable 'g' is a pointer to a funcval struct.
  funcval contains a pointer to the actual function code.
  This indirection is what enables function values to exist at all.
```

**Why the indirection?** Because a function variable might hold different functions
at runtime. The funcval tells the runtime which code to jump to.

### Function Type vs Plain Function — The Key Difference

```go
// Plain function — fixed, no methods possible
func add(a, b int) int { return a + b }

// Function type — a TYPE that can have methods
type MathOp func(a, b int) int

// Now MathOp can have methods!
func (op MathOp) RunAndLog(a, b int) int {
    result := op(a, b)       // call the function itself
    fmt.Println("result:", result)
    return result
}

// Usage:
var op MathOp = add           // assign plain function to typed variable
op.RunAndLog(3, 4)            // calls add(3, 4) AND logs it
```

```
  Plain function:      has NO methods, is just code
  Function type:       IS a type, CAN have methods, CAN satisfy interfaces

  This is not possible in most languages:
  ❌ Java: functions can't have methods (they're not objects)
  ❌ Python: functions are objects but you don't typically add methods
  ❌ JavaScript: functions are objects but no static type system
  ✅ Go: function types are types, types can have methods → functions satisfy interfaces
```

---

## 2. Methods on Function Types — Go's Secret Weapon

This is where Go unlocks something no other mainstream language has:

### The http.HandlerFunc Pattern

```go
// From net/http/server.go — the actual stdlib code

// Step 1: Define a function type
type HandlerFunc func(ResponseWriter, *Request)

// Step 2: Give it a method that calls itself
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)    // f IS a function — call it
}

// Step 3: The Handler interface requires ServeHTTP
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// Result: HandlerFunc satisfies Handler interface!
```

Let's trace what happens step by step:

```
  Step 1: You write a plain function
  ───────────────────────────────────
  func myHandler(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("hello"))
  }

  Step 2: Convert to HandlerFunc (type conversion, not a function call)
  ────────────────────────────────────────────────────────────────────
  h := http.HandlerFunc(myHandler)

  What 'h' looks like in memory:
  ┌───────────────────────────────────────────────┐
  │ h (HandlerFunc)                                │
  │                                                │
  │ Underlying: func(ResponseWriter, *Request)     │
  │ Points to:  myHandler code at 0x4a3000        │
  │                                                │
  │ Has method: ServeHTTP(w, r) { h(w, r) }       │
  │ Therefore:  satisfies http.Handler interface    │
  └───────────────────────────────────────────────┘

  Step 3: Use as http.Handler anywhere
  ─────────────────────────────────────
  http.ListenAndServe(":8080", h)   // expects Handler, gets HandlerFunc ✅

  Step 4: When a request arrives
  ──────────────────────────────
  h.ServeHTTP(w, r)    // calls HandlerFunc.ServeHTTP
      → f(w, r)         // which calls myHandler(w, r)
          → w.Write([]byte("hello"))
```

### The Elegance

This pattern turns **any compatible function into an interface implementation**
with zero boilerplate. No anonymous classes. No wrapper objects. No decorators.
Just a type conversion:

```go
// This one-liner makes a function into an http.Handler:
http.HandlerFunc(myFunc)

// Java equivalent — minimum 3 lines (anonymous class):
new HttpHandler() {
    @Override
    public void handle(HttpExchange exchange) {
        myFunc(exchange);
    }
};
```

---

## 3. The Middleware Pattern — Under the Hood

### The Signature

```go
type Middleware func(http.Handler) http.Handler
//                   ↑ takes a handler    ↑ returns a NEW handler
//                   (the "next" in chain)  (wraps "next" with extra behavior)
```

### A Complete Middleware — Step by Step

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()                          // PRE-PROCESSING
        next.ServeHTTP(w, r)                         // DELEGATE
        slog.Info("request", "dur", time.Since(start)) // POST-PROCESSING
    })
}
```

Let's break down every piece:

```
  LoggingMiddleware(next http.Handler) http.Handler
  │                  │                  │
  │                  │                  └── Returns: a NEW Handler
  │                  └── Takes: the NEXT handler in the chain
  └── Function name

  return http.HandlerFunc(func(w, r) { ... })
  │      │                │
  │      │                └── Anonymous function (closure)
  │      └── Type conversion: func → HandlerFunc → Handler
  └── Returns the wrapped handler

  Inside the anonymous function:
  ┌──────────────────────────────────────────────────┐
  │  func(w, r) {                                    │
  │      start := time.Now()                         │  ← runs BEFORE next
  │      next.ServeHTTP(w, r)                        │  ← calls next handler
  │      slog.Info(...)                              │  ← runs AFTER next
  │  }                                               │
  │                                                  │
  │  This function CAPTURES 'next' from the outer    │
  │  function's scope. 'next' is a free variable.   │
  │  This makes it a CLOSURE.                        │
  └──────────────────────────────────────────────────┘
```

### Chaining — What Actually Happens

```go
handler := LoggingMiddleware(AuthMiddleware(rateLimitMiddleware(myHandler)))
```

Let's unwind this from inside out:

```
  Step 1: rateLimitMiddleware(myHandler)
  ┌─────────────────────────────────────────────────┐
  │ Returns HandlerFunc that:                        │
  │   1. Checks rate limit                           │
  │   2. Calls myHandler.ServeHTTP(w, r)            │
  │                                                  │
  │ Captures: next = myHandler                       │
  │ Let's call this: rlHandler                       │
  └─────────────────────────────────────────────────┘

  Step 2: AuthMiddleware(rlHandler)
  ┌─────────────────────────────────────────────────┐
  │ Returns HandlerFunc that:                        │
  │   1. Validates JWT token                         │
  │   2. Calls rlHandler.ServeHTTP(w, r)            │
  │                                                  │
  │ Captures: next = rlHandler                       │
  │ Let's call this: authHandler                     │
  └─────────────────────────────────────────────────┘

  Step 3: LoggingMiddleware(authHandler)
  ┌─────────────────────────────────────────────────┐
  │ Returns HandlerFunc that:                        │
  │   1. Starts timer                                │
  │   2. Calls authHandler.ServeHTTP(w, r)          │
  │   3. Logs duration                               │
  │                                                  │
  │ Captures: next = authHandler                     │
  │ This is the FINAL handler we register            │
  └─────────────────────────────────────────────────┘
```

### The Call Stack When a Request Arrives

```
  HTTP Request → handler.ServeHTTP(w, r)
  │
  ├── LoggingMiddleware's closure:
  │   │  start := time.Now()
  │   │
  │   ├── AuthMiddleware's closure:        (next = authHandler)
  │   │   │  token := r.Header.Get("Authorization")
  │   │   │  if invalid → w.WriteHeader(401); return  ← SHORT CIRCUIT
  │   │   │
  │   │   ├── RateLimitMiddleware's closure:  (next = rlHandler)
  │   │   │   │  if over limit → w.WriteHeader(429); return  ← SHORT CIRCUIT
  │   │   │   │
  │   │   │   ├── myHandler:                    (actual business logic)
  │   │   │   │     w.Write([]byte(`{"orders": [...]}`))
  │   │   │   │
  │   │   │   └── return to RateLimit (nothing to do after)
  │   │   │
  │   │   └── return to Auth (nothing to do after)
  │   │
  │   └── slog.Info("request", "dur", time.Since(start))
  │
  └── Response sent to client

  KEY INSIGHT: Any middleware can SHORT-CIRCUIT the chain by NOT calling
  next.ServeHTTP(w, r). Auth returns 401 without ever reaching your handler.
  This is the power of the onion model.
```

---

## 4. How Closures Work in Middleware

This is where Chapter 06 (closures) meets production code.

### What the Closure Captures

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // This function captures 'next' from the outer scope
        // 'next' is a FREE VARIABLE — not defined in this function
        next.ServeHTTP(w, r)
    })
}
```

```
  Memory layout after AuthMiddleware(someHandler) returns:

  Stack (AuthMiddleware frame):        Heap:
  ┌─────────────────────┐            ┌──────────────────────────────┐
  │ next: someHandler   │─ escapes →│ funcval (closure)            │
  │ return value: ──────┼───────────→│ fn: 0x4a5000 (anon func)   │
  └─────────────────────┘            │ next: someHandler (captured)│
  ↑ this frame is freed              └──────────────────────────────┘
    after return                       ↑ lives on heap because the
                                         returned HandlerFunc references it
```

The compiler's escape analysis sees that `next` is captured by the returned
closure → `next` must live on the heap. The `funcval` struct on the heap
holds both the function pointer AND the captured variable.

### Closures with Configuration

Middleware often needs configuration. Closures capture that too:

```go
func TimeoutMiddleware(maxDuration time.Duration) func(http.Handler) http.Handler {
    //                   ↑ configuration captured by closure
    return func(next http.Handler) http.Handler {
        //        ↑ 'next' also captured
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), maxDuration)
            //                                              ↑ uses captured config
            defer cancel()
            next.ServeHTTP(w, r.WithContext(ctx))
            //   ↑ uses captured 'next'
        })
    }
}

// Usage:
r.Use(TimeoutMiddleware(30 * time.Second))
```

```
  funcval for the innermost closure captures TWO variables:

  ┌──────────────────────────────────────┐
  │ funcval                              │
  │   fn:          anon func code ptr   │
  │   maxDuration: 30s (captured)       │  ← from outer outer scope
  │   next:        someHandler (captured)│  ← from outer scope
  └──────────────────────────────────────┘

  Three levels of nesting, two captured variables.
  The compiler flattens all captures into one funcval on the heap.
```

---

## 5. Without These Features — The Java Way

To truly appreciate Go's design, let's see what middleware looks like **without**
methods on function types, structural typing, and closures.

### Java — The Full Boilerplate

```java
// Step 1: Define the Handler interface (like Go — similar)
public interface Handler {
    void handle(Request req, Response resp);
}

// Step 2: Define middleware — needs an ABSTRACT CLASS
public abstract class Middleware implements Handler {
    protected Handler next;

    public Middleware(Handler next) {
        this.next = next;
    }
}

// Step 3: Implement a concrete middleware — needs a CLASS
public class LoggingMiddleware extends Middleware {
    public LoggingMiddleware(Handler next) {
        super(next);
    }

    @Override
    public void handle(Request req, Response resp) {
        long start = System.currentTimeMillis();
        next.handle(req, resp);  // delegate
        logger.info("Duration: " + (System.currentTimeMillis() - start));
    }
}

// Step 4: Implement your handler — another CLASS
public class OrderHandler implements Handler {
    @Override
    public void handle(Request req, Response resp) {
        resp.write("{\"orders\": []}");
    }
}

// Step 5: Chain them together
Handler handler = new LoggingMiddleware(
    new AuthMiddleware(
        new RateLimitMiddleware(
            new OrderHandler()
        )
    )
);
```

Count the pieces:
```
  Java needs:                           Go needs:
  ├── Handler interface           1     ├── Handler interface         1
  ├── Middleware abstract class    1     └── (that's it)
  ├── LoggingMiddleware class      1
  ├── AuthMiddleware class         1     Functions:
  ├── RateLimitMiddleware class    1     ├── LoggingMiddleware         1
  ├── OrderHandler class           1     ├── AuthMiddleware            1
  └── Each needs constructor       5     ├── RateLimitMiddleware       1
                                         └── orderHandler              1
  Total: 6 types + 5 constructors        Total: 1 type + 4 functions

  Why the difference?
  ├── Go: HandlerFunc eliminates the need for abstract Middleware class
  ├── Go: structural typing eliminates "implements Handler" everywhere
  ├── Go: closures replace the need for concrete middleware classes
  └── Go: no constructors needed (closures capture state directly)
```

### What If Go Didn't Have These Features?

```go
// WITHOUT methods on function types — need wrapper struct for every handler:
type myHandler struct{}
func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello"))
}

// WITHOUT closures — need struct to hold 'next':
type loggingMiddleware struct {
    next http.Handler
}
func (m *loggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    m.next.ServeHTTP(w, r)
    slog.Info("request", "dur", time.Since(start))
}

// Chaining:
handler := &loggingMiddleware{
    next: &authMiddleware{
        next: &rateLimitMiddleware{
            next: &myHandler{},
        },
    },
}

// This is essentially the Java approach — in Go!
// It works, but it's the OPPOSITE of idiomatic Go.
```

---

## 6. ResponseWriter Wrapping — Struct Embedding in Action

A common middleware need: capture the HTTP status code that the handler wrote.
Problem: `http.ResponseWriter` doesn't expose it after `WriteHeader()` is called.

```go
// Step 1: Create a wrapper that EMBEDS ResponseWriter
type statusCapture struct {
    http.ResponseWriter          // embedded — delegates all methods
    statusCode int
}

// Step 2: Override only WriteHeader to intercept the status
func (sc *statusCapture) WriteHeader(code int) {
    sc.statusCode = code
    sc.ResponseWriter.WriteHeader(code)  // delegate to original
}
```

```
  How struct embedding works here:

  http.ResponseWriter interface has 3 methods:
    Header() http.Header
    Write([]byte) (int, error)
    WriteHeader(statusCode int)

  statusCapture:
  ┌──────────────────────────────────────────────────┐
  │ http.ResponseWriter (embedded)                    │
  │   → Header()      ← delegated (auto-promoted)   │
  │   → Write()       ← delegated (auto-promoted)   │
  │   → WriteHeader() ← OVERRIDDEN by our method    │
  │                                                  │
  │ statusCode int    ← our extra field              │
  └──────────────────────────────────────────────────┘

  The struct satisfies http.ResponseWriter because:
  - Header() and Write() are promoted from the embedded field
  - WriteHeader() is provided by our explicit method
  All three methods exist → interface satisfied (structural typing!)
```

Usage in middleware:

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sc := &statusCapture{ResponseWriter: w, statusCode: 200}
        start := time.Now()

        next.ServeHTTP(sc, r)  // pass OUR wrapper, not original w
        //                ↑ downstream writes go through our wrapper

        slog.Info("request",
            "method", r.Method,
            "status", sc.statusCode,      // we captured it!
            "dur",    time.Since(start),
        )
    })
}
```

---

## 7. Closures Across Languages — Go vs Java vs JavaScript

### Go Closures — Capture by Reference, Full Mutation

```go
count := 0
increment := func() int {
    count++          // mutates captured variable directly
    return count
}
increment() // count = 1
increment() // count = 2
```

Under the hood: `count` escapes to heap. The closure's funcval holds a
pointer to `count`. Multiple closures sharing the same variable see each
other's mutations.

### Java "Closures" (Lambdas) — Effectively Final Only

```java
int count = 0;
Runnable fn = () -> count++;    // ❌ COMPILE ERROR
// "local variables referenced from a lambda must be final or effectively final"

// Workaround 1: use a mutable wrapper
AtomicInteger count = new AtomicInteger(0);
Runnable fn = () -> count.incrementAndGet();    // ✅

// Workaround 2: use a single-element array (ugly but common)
int[] count = {0};
Runnable fn = () -> count[0]++;                 // ✅
```

Why Java restricts this: Java lambdas capture the **value** of the variable,
not a reference to it. If the variable were mutable, the lambda and the
enclosing scope would have **different copies** — confusing and error-prone.
Java chose safety over flexibility.

### JavaScript Closures — Capture by Reference, Like Go

```javascript
let count = 0;
const increment = () => ++count;    // ✅ captures by reference
increment(); // count = 1
increment(); // count = 2
```

JavaScript closures work like Go's — capture by reference, full mutation.
But JavaScript has no type system to catch misuse at compile time, and
the event loop makes concurrent access to captured variables rare (single-threaded).

### The Comparison Table

```
  Feature              Go              Java 8+           JavaScript
  ──────               ──              ───────            ──────────
  Capture mechanism    By reference    By value (final)   By reference
  Mutation allowed     ✅ Yes           ❌ Only wrappers    ✅ Yes
  Methods on func      ✅ Yes           ❌ No               ❌ No
  Structural typing    ✅ Yes           ❌ No (nominal)     N/A (duck typed)
  Concurrency safe     ❌ Manual        ❌ Manual           N/A (single thread)
  Middleware from       3 features      abstract class     app.use(fn) works
    closures alone?     = zero boiler   + annotations      but no type safety
```

---

## 8. Production Middleware Stack

### Real-World Ordering (chi router)

```go
r := chi.NewRouter()

// Layer 1 — Infrastructure (outermost, runs first)
r.Use(middleware.RequestID)              // assign unique ID
r.Use(middleware.RealIP)                 // extract X-Forwarded-For
r.Use(middleware.Logger)                 // log every request
r.Use(middleware.Recoverer)              // defer+recover → 500 not crash
r.Use(middleware.Timeout(30*time.Second))// kill slow requests

// Layer 2 — Security
r.Use(corsMiddleware)                    // CORS headers
r.Use(authMiddleware)                    // validate JWT, set user in ctx

// Layer 3 — Rate Limiting
r.Use(rateLimitMiddleware)               // per-user token bucket

// Layer 4 — Routes (innermost)
r.Route("/orders", func(r chi.Router) {
    r.Get("/", listOrders)
    r.Post("/", createOrder)
    r.Route("/{id}", func(r chi.Router) {
        r.Get("/", getOrder)
        r.Patch("/cancel", cancelOrder)
    })
})

// Admin routes — ADDITIONAL middleware
r.Route("/admin", func(r chi.Router) {
    r.Use(adminOnlyMiddleware)           // on top of all above
    r.Get("/users", listUsers)
})
```

### How chi.Use() Works Internally

```go
// Simplified chi implementation
type Mux struct {
    middlewares []func(http.Handler) http.Handler
    // ...
}

func (mx *Mux) Use(middlewares ...func(http.Handler) http.Handler) {
    mx.middlewares = append(mx.middlewares, middlewares...)
}

// When a route is matched, chi chains all middlewares:
func chain(middlewares []func(http.Handler) http.Handler, endpoint http.Handler) http.Handler {
    // Wrap from inside out
    handler := endpoint
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}
```

That's the entire middleware engine — a slice of functions and a loop.
No reflection, no dependency injection container, no annotation processor.

---

## 9. The Three Features That Make It All Work

```
  Feature 1: Methods on Function Types
  ─────────────────────────────────────
  HandlerFunc has a ServeHTTP method → function satisfies Handler interface
  Without this: need a wrapper struct for every handler (Java pattern)

  Feature 2: Structural Typing (Implicit Interfaces)
  ──────────────────────────────────────────────────
  Any type with ServeHTTP(w, r) IS a Handler — no "implements" keyword
  Without this: every middleware must explicitly declare "implements Handler"

  Feature 3: Closures (Capture by Reference)
  ──────────────────────────────────────────
  Middleware returns a closure that captures 'next' handler + config
  Without this: need struct fields to hold 'next' and configuration

  Remove any ONE of these:
  ┌─────────────────────────┬────────────────────────────────────┐
  │ Missing Feature         │ What You'd Need Instead            │
  ├─────────────────────────┼────────────────────────────────────┤
  │ Methods on func types   │ Wrapper struct per handler         │
  │ Structural typing       │ "implements Handler" everywhere    │
  │ Closures                │ Struct fields for captured state   │
  │ All three               │ Java-style abstract class pattern  │
  └─────────────────────────┴────────────────────────────────────┘
```

### Connection to Go Design Philosophy (Chapter 21)

The middleware pattern is the **convergence point** of Go's core design decisions:

```
  Interface (1-method, implicit)
       │
       ├── enables HandlerFunc adapter pattern
       │
  Function types (can have methods)
       │
       ├── eliminates wrapper struct boilerplate
       │
  Closures (capture by reference)
       │
       ├── replaces constructor + fields with captured variables
       │
  Struct embedding
       │
       └── enables ResponseWriter wrapping with minimal code

  Result: the entire HTTP middleware ecosystem — chi, echo, gin, stdlib —
  built on func(Handler) Handler. No framework needed. No magic.
  Just types and functions.
```

> **Go Wisdom**: *"Clear is better than clever."*
> The middleware pattern is not clever. It's three features composing naturally.
> That's the point — Go's simplicity isn't a limitation, it's a design decision
> that makes patterns like this emerge without framework support.
