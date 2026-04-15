# 27 — net/http Under the Hood: Server, ServeMux & the Handler Chain

> **Companion exercises:** [exercises/stdlib/10_net_http](../exercises/stdlib/10_net_http/)

---

## Table of Contents

1. [The Handler Interface — Go's HTTP Foundation](#1-the-handler-interface)
2. [Server Internals — From Listen to Response](#2-server-internals)
3. [ServeMux — Pattern Matching & Routing](#3-servemux)
4. [The Middleware Pattern — func(Handler) Handler](#4-the-middleware-pattern)
5. [ResponseWriter — What Happens When You Write](#5-responsewriter)
6. [Request Lifecycle — The Full Journey](#6-request-lifecycle)
7. [Connection Handling — Keep-Alive, HTTP/2, Hijack](#7-connection-handling)
8. [httptest — Testing Without a Network](#8-httptest)
9. [Timeouts — The Production Essentials](#9-timeouts)
10. [Graceful Shutdown](#10-graceful-shutdown)
11. [Cost Table](#11-cost-table)
12. [Quick Reference Card](#12-quick-reference-card)
13. [Further Reading](#13-further-reading)

---

## 1. The Handler Interface

The entire `net/http` package revolves around one interface:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

One method. That's it. Every HTTP handler, middleware, mux, and framework
in Go implements or wraps this single interface. This is the "accept
interfaces, return structs" philosophy in its purest form.

### HandlerFunc — The Adapter

```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

`HandlerFunc` is a type conversion, not a wrapper. It makes any plain function
satisfy the `Handler` interface. This is a Go pattern you'll see everywhere:
**method on function type** adapts a function to an interface.

### Under the Hood

When you write `http.HandleFunc("/hello", myFunc)`:
1. `myFunc` is converted to `HandlerFunc(myFunc)` 
2. `HandlerFunc` satisfies `Handler` via its `ServeHTTP` method
3. The mux stores a `Handler` (not a function)

Zero allocation. Zero overhead. The function pointer is the handler.

---

## 2. Server Internals

### The Server Struct

```go
type Server struct {
    Addr      string        // ":8080"
    Handler   Handler       // usually a ServeMux
    TLSConfig *tls.Config
    
    ReadTimeout    time.Duration  // time to read entire request
    WriteTimeout   time.Duration  // time to write response
    IdleTimeout    time.Duration  // keep-alive idle timeout
    MaxHeaderBytes int            // default 1MB
    
    // ... internal fields
}
```

### What ListenAndServe Does

```go
func (srv *Server) ListenAndServe() error {
    ln, err := net.Listen("tcp", srv.Addr)  // 1. bind socket
    if err != nil {
        return err
    }
    return srv.Serve(ln)                     // 2. accept loop
}
```

### The Accept Loop (runtime/server.go)

```
┌─────────────────────────────────────────────┐
│              srv.Serve(listener)              │
├─────────────────────────────────────────────┤
│  for {                                       │
│    conn, err := listener.Accept()  ←─────── blocks on epoll/kqueue
│    if err != nil { ... }                     │
│    go srv.newConn(conn).serve(ctx) ←──────── one goroutine per conn
│  }                                           │
└─────────────────────────────────────────────┘
```

**One goroutine per connection.** This is the Go model:
- No thread pool sizing decisions
- No async/await coloring
- The runtime scheduler (GMP) handles thousands of goroutines efficiently
- Each goroutine starts at ~2KB stack (grows as needed)

### Cost of 10,000 Concurrent Connections

```
10,000 connections × ~8KB goroutine stack = ~80MB
(plus per-connection buffers: ~4KB read + 4KB write = ~80MB more)
Total: ~160MB for 10K concurrent connections
```

Compare with Java's thread-per-request model: 10K threads × 1MB stack = 10GB.
This is why Go servers handle high concurrency with low memory.

### The Connection Serve Loop

Each connection goroutine runs this loop:

```
┌─────────────────────────────────────────────┐
│           conn.serve(ctx)                    │
├─────────────────────────────────────────────┤
│  for {                                       │
│    req, err := conn.readRequest()            │
│    if err != nil { break }                   │
│                                              │
│    w := &response{conn: c, req: req}         │
│    serverHandler{srv}.ServeHTTP(w, req)  ←── calls YOUR handler
│                                              │
│    w.finishRequest()                         │
│    if !w.shouldReuseConnection() { break }   │
│  }                                           │
│  conn.close()                                │
└─────────────────────────────────────────────┘
```

For HTTP/1.1 keep-alive, the loop handles multiple requests on one TCP
connection. For HTTP/2, a different multiplexed path is used.

---

## 3. ServeMux

### The Default Mux vs Custom Mux

```go
// Default (global) mux — AVOID in production (it's a global variable)
http.HandleFunc("/hello", handler)
http.ListenAndServe(":8080", nil)  // nil = use DefaultServeMux

// Custom mux — ALWAYS use this
mux := http.NewServeMux()
mux.HandleFunc("/hello", handler)
http.ListenAndServe(":8080", mux)
```

**Why avoid DefaultServeMux:** It's a package-level global. Any imported
package can register routes on it. This is a security risk (third-party
dependencies could add debug endpoints to your server).

### Pattern Matching (Pre Go 1.22)

Old-style patterns were simple path prefixes:
- `/hello` matches exactly `/hello`
- `/hello/` matches `/hello/` and any subpath (`/hello/world`, `/hello/a/b`)
- Longest match wins

### Enhanced Routing (Go 1.22+)

Go 1.22 added method-based routing and path parameters:

```go
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("GET /files/{path...}", serveFile)  // wildcard
```

Key features:
- `{id}` captures a path segment: `r.PathValue("id")`
- `{path...}` captures remaining path (greedy)
- `GET /users/{id}` only matches GET requests
- More specific patterns take precedence

### Pattern Priority (Go 1.22+)

```
1. "GET /users/admin"      — exact method + exact path (highest)
2. "GET /users/{id}"       — exact method + path param
3. "/users/{id}"           — any method + path param
4. "GET /users/"           — exact method + prefix
5. "/users/"               — any method + prefix (lowest)
```

### Internal Data Structure

`ServeMux` uses a tree of patterns (Go 1.22+) for O(path-depth) lookup,
replacing the old sorted slice of patterns. The tree supports:
- Method-specific routing at each node
- Path parameter capture via `{name}` nodes
- Wildcard matching via `{name...}` leaf nodes

---

## 4. The Middleware Pattern

### The Signature

```go
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // before
        next.ServeHTTP(w, r)
        // after
    })
}
```

This is the **decorator pattern** applied to `http.Handler`. Each middleware:
1. Wraps the next handler
2. Can modify the request (before)
3. Calls the next handler
4. Can modify the response or do cleanup (after)

### The Onion Model

```
Request →  [Logging → [Auth → [Recovery → [Handler] ] ] ] → Response
```

When you chain middleware: `Logging(Auth(Recovery(handler)))`, the execution
order is:

```
Logging-before
  Auth-before
    Recovery-before
      Handler
    Recovery-after
  Auth-after
Logging-after
```

### Composing Middleware

```go
type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, mws ...Middleware) http.Handler {
    for i := len(mws) - 1; i >= 0; i-- {
        handler = mws[i](handler)
    }
    return handler
}
```

Iterate in **reverse** so the first middleware in the list is the outermost
wrapper (first to execute).

### Common Production Middleware

| Middleware | Purpose | Typical Position |
|-----------|---------|-----------------|
| Recovery | Catch panics, return 500 | Outermost |
| Logging | Log method, path, duration, status | Near outer |
| Request ID | Generate/propagate X-Request-ID | Early |
| CORS | Set Access-Control headers | Before auth |
| Auth | Validate tokens, set user in context | Before handlers |
| Rate Limit | Per-client request throttling | After auth |
| Timeout | `context.WithTimeout` on requests | Inner |
| Compression | gzip response body | Near inner |

### ResponseWriter Wrapping

To capture the status code in logging middleware, you must wrap
`ResponseWriter`:

```go
type statusRecorder struct {
    http.ResponseWriter
    status int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.status = code
    r.ResponseWriter.WriteHeader(code)
}
```

This is necessary because `http.ResponseWriter` doesn't expose the status
code after `WriteHeader` is called.

---

## 5. ResponseWriter

### The Interface

```go
type ResponseWriter interface {
    Header() http.Header           // get response headers map
    Write([]byte) (int, error)     // write body bytes
    WriteHeader(statusCode int)    // send status code
}
```

### The Header-First Rule

```
w.Header().Set("Content-Type", "application/json")  ← 1. Set headers
w.WriteHeader(201)                                    ← 2. Send status
w.Write(body)                                         ← 3. Write body
```

**Critical:** Headers must be set BEFORE the first `Write()` call. Once
bytes are written, the status line and headers are already sent on the wire.
This is because HTTP uses a streaming protocol: headers come first.

If you call `Write()` without calling `WriteHeader()`, Go implicitly sends
`200 OK` as the status. This is a common source of bugs:

```go
// BUG: trying to set 400 AFTER writing
w.Write([]byte("bad request"))
w.WriteHeader(400)  // too late — 200 already sent!
```

### Flusher, Hijacker, Pusher

`ResponseWriter` can be type-asserted to additional interfaces:

```go
if flusher, ok := w.(http.Flusher); ok {
    flusher.Flush()  // force-send buffered data (SSE, streaming)
}

if hijacker, ok := w.(http.Hijacker); ok {
    conn, buf, _ := hijacker.Hijack()  // take over TCP connection (WebSocket)
}
```

**Caution:** If your middleware wraps `ResponseWriter`, these type assertions
break unless you also implement the optional interfaces.

---

## 6. Request Lifecycle

```
Client                    Server
  │                         │
  ├─── TCP Connect ────────►│  listener.Accept()
  │                         │  go conn.serve(ctx)
  ├─── HTTP Request ───────►│  conn.readRequest()
  │                         │  ↓
  │                         │  mux.ServeHTTP(w, req)
  │                         │    ↓ middleware chain
  │                         │    ↓ handler
  │                         │    ↓ w.Write(response)
  │◄── HTTP Response ───────┤  w.finishRequest()
  │                         │  (keep-alive: loop back to readRequest)
  │                         │
  ├─── TCP Close ──────────►│  conn.close()
```

### Context Propagation

Every `*http.Request` carries a context:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // cancelled when client disconnects
    
    // Pass to downstream calls
    result, err := db.QueryContext(ctx, "SELECT ...")
}
```

The request context is cancelled when:
1. The client closes the connection
2. The server's `ReadTimeout` or `WriteTimeout` fires
3. You manually cancel it via middleware

---

## 7. Connection Handling

### Keep-Alive (HTTP/1.1)

By default, Go reuses TCP connections for multiple requests:
- `Connection: keep-alive` is the default for HTTP/1.1
- `Server.IdleTimeout` controls how long idle connections stay open
- After the timeout, the connection is closed

### HTTP/2

Go's `net/http` supports HTTP/2 automatically when using TLS:

```go
srv := &http.Server{Addr: ":443", Handler: mux}
srv.ListenAndServeTLS("cert.pem", "key.pem")  // HTTP/2 enabled
```

HTTP/2 features:
- **Multiplexing**: multiple requests on one TCP connection
- **Server push**: `http.Pusher` interface (deprecated in practice)
- **Header compression**: HPACK encoding

### Connection Hijacking

For WebSocket or custom protocols:

```go
hijacker := w.(http.Hijacker)
conn, bufrw, _ := hijacker.Hijack()
// Now you own the raw TCP connection
defer conn.Close()
```

After hijacking, the server no longer manages the connection. You're
responsible for reading, writing, and closing.

---

## 8. httptest

### Testing Without a Server

```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/hello?name=Go", nil)
    w := httptest.NewRecorder()
    
    MyHandler(w, req)
    
    resp := w.Result()
    body, _ := io.ReadAll(resp.Body)
    
    if resp.StatusCode != 200 { t.Fatal(...) }
    if string(body) != "Hello, Go!" { t.Fatal(...) }
}
```

`httptest.NewRecorder()` returns a `*ResponseRecorder` that implements
`ResponseWriter` and buffers the response for inspection. No TCP, no port,
no network. Tests run in microseconds.

### Testing With a Real Server

```go
func TestAPI(t *testing.T) {
    srv := httptest.NewServer(myHandler)
    defer srv.Close()
    
    resp, _ := http.Get(srv.URL + "/hello")
    // srv.URL is like "http://127.0.0.1:54321"
}
```

`httptest.NewServer` starts a real HTTP server on a random port. Useful for
integration tests that need real HTTP behavior (redirects, cookies, TLS).

### TLS Testing

```go
srv := httptest.NewTLSServer(myHandler)
client := srv.Client()  // pre-configured to trust the test CA
resp, _ := client.Get(srv.URL + "/secure")
```

---

## 9. Timeouts

### The Three Timeouts

```go
srv := &http.Server{
    ReadTimeout:  5 * time.Second,   // entire request read
    WriteTimeout: 10 * time.Second,  // entire response write
    IdleTimeout:  120 * time.Second, // between keep-alive requests
}
```

```
┌──────────────────────────────────────────────────────────────┐
│  ReadTimeout              WriteTimeout                        │
│  ├──────────────┤         ├──────────────────────────────────┤│
│  │  Read Headers │         │  Handler execution + Response    ││
│  │  Read Body    │         │  write                           ││
│  │              │         │                                   ││
│  Accept ────────► Handler Start ─────────────► Response Done  │
│                                                               │
│              IdleTimeout                                      │
│              ├──────────────────────────────────────────┤      │
│              Response Done ────────────► Next Request    │      │
└──────────────────────────────────────────────────────────────┘
```

### Per-Request Timeout via Context

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()
    
    result, err := slowService(ctx)
    if err == context.DeadlineExceeded {
        http.Error(w, "timeout", http.StatusGatewayTimeout)
        return
    }
}
```

### Production Default Recommendations

| Timeout | API Server | File Upload | Proxy |
|---------|-----------|-------------|-------|
| ReadTimeout | 5s | 60s | 10s |
| WriteTimeout | 10s | 60s | 30s |
| IdleTimeout | 120s | 120s | 90s |

**Never run without timeouts in production.** A server with no ReadTimeout
is vulnerable to Slowloris attacks (slow clients holding connections open
indefinitely).

---

## 10. Graceful Shutdown

```go
func main() {
    srv := &http.Server{Addr: ":8080", Handler: mux}
    
    // Start server in background
    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // Wait for interrupt signal
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()
    <-ctx.Done()
    
    // Graceful shutdown with timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Printf("shutdown error: %v", err)
    }
}
```

### What Shutdown Does

1. **Closes listeners** — no new connections accepted
2. **Waits for in-flight requests** — currently running handlers finish
3. **Closes idle connections** — keep-alive connections with no active request
4. **Respects the context deadline** — if timeout fires, remaining connections are forcefully closed

### The Shutdown Sequence

```
SIGTERM received
  ↓
srv.Shutdown(ctx) called
  ↓
Listener closed ← no new connections
  ↓
Wait for active handlers to finish
  ↓
Close idle connections
  ↓
Return nil (or ctx.Err() if timeout)
```

---

## 11. Cost Table

| Operation | Cost | Notes |
|-----------|------|-------|
| Handler interface dispatch | ~1ns | Virtual method call |
| Middleware layer (no-op) | ~5-10ns | Function call + closure |
| `httptest.NewRecorder()` | ~200ns | Allocates buffer |
| `httptest.NewRequest()` | ~500ns | Parses URL |
| Real TCP connection | ~50-100μs | Network round trip |
| JSON marshal + write | ~1-10μs | Depends on payload size |
| TLS handshake | ~1-5ms | Certificate verification |
| Full request cycle (local) | ~10-100μs | Without network latency |

---

## 12. Quick Reference Card

```text
┌─────────────────────────────────────────────────────────────────┐
│                     NET/HTTP CHEAT SHEET                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Handler interface:                                             │
│    type Handler interface { ServeHTTP(w, r) }                   │
│    http.HandlerFunc(fn)  ← adapter for plain functions          │
│                                                                 │
│  Routing (Go 1.22+):                                            │
│    mux.HandleFunc("GET /users/{id}", fn)                        │
│    id := r.PathValue("id")                                      │
│                                                                 │
│  Middleware:                                                     │
│    func MW(next http.Handler) http.Handler {                    │
│      return http.HandlerFunc(func(w, r) {                       │
│        /* before */ next.ServeHTTP(w, r) /* after */             │
│      })                                                         │
│    }                                                            │
│                                                                 │
│  Testing:                                                       │
│    w := httptest.NewRecorder()                                  │
│    r := httptest.NewRequest("GET", "/path", body)               │
│    handler.ServeHTTP(w, r)                                      │
│    resp := w.Result()                                           │
│                                                                 │
│  Timeouts (MANDATORY in production):                            │
│    ReadTimeout, WriteTimeout, IdleTimeout                       │
│                                                                 │
│  Graceful shutdown:                                             │
│    srv.Shutdown(ctx) — waits for in-flight, closes idle         │
│                                                                 │
│  Header order: Set headers → WriteHeader → Write body           │
│                                                                 │
│  AVOID: http.DefaultServeMux (global, security risk)            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 13. Further Reading

- [`net/http` package source](https://cs.opensource.google/go/go/+/master:src/net/http/) — read `server.go` first
- [ServeMux Routing Enhancements (Go 1.22)](https://go.dev/blog/routing-enhancements) — official blog post
- [HTTP/2 in Go](https://pkg.go.dev/golang.org/x/net/http2) — the `x/net/http2` package
- [Timeout guide by Cloudflare](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/) — production timeout patterns
- [Middleware best practices](https://pkg.go.dev/github.com/justinas/alice) — `alice` for middleware chaining
- [Go HTTP server architecture (talk by Brad Fitzpatrick)](https://www.youtube.com/watch?v=jJZ-1VBR070) — from the creator of `net/http`

---

> **Next:** [28 — reflect Under the Hood](28_reflect_under_the_hood.md)
>
> **Companion exercises:** [exercises/stdlib/10_net_http](../exercises/stdlib/10_net_http/)
