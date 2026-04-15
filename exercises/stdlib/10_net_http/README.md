# 10 net/http — Handlers, Middleware & Servers

> **Companion chapter:** [learnings/27_net_http_under_the_hood.md](../../../learnings/27_net_http_under_the_hood.md)

## Exercises

| # | Function / Type | Concepts | Difficulty |
|---|----------------|----------|------------|
| 1 | `HelloHandler` | `http.Handler` interface, `ServeHTTP` | ⭐ |
| 2 | `GreetHandler` | Query parameters, `r.URL.Query()` | ⭐ |
| 3 | `EchoHandler` | Request body, `io.ReadAll`, `defer Close` | ⭐ |
| 4 | `MethodRouter` | `r.Method`, method dispatch, 405 | ⭐⭐ |
| 5 | `JSONResponse` | `Content-Type` header, JSON response | ⭐⭐ |
| 6 | `StatusHandler` | `http.StatusText`, custom status codes | ⭐⭐ |
| 7 | `SetupMux` | `http.ServeMux`, `Handle` vs `HandleFunc` | ⭐⭐ |
| 8 | `LoggingMiddleware` | Middleware pattern, `func(http.Handler) http.Handler` | ⭐⭐⭐ |
| 9 | `RecoveryMiddleware` | `defer`/`recover`, panic-safe handlers | ⭐⭐⭐ |
| 10 | `AuthMiddleware` | Authorization header, 401 responses | ⭐⭐⭐ |
| 11 | `ChainMiddleware` | Middleware composition, onion model | ⭐⭐⭐ |
| 12 | `HeadersHandler` | Read/write headers, `X-Request-ID` pattern | ⭐⭐ |

## How to Practice

```bash
# Run all tests (they all fail initially)
go test -race -v ./exercises/stdlib/10_net_http/

# Run a specific test
go test -race -run TestLoggingMiddleware ./exercises/stdlib/10_net_http/
```

## Key Insights

- **`http.Handler`** is a single-method interface: `ServeHTTP(w, r)`. That's it
- **`http.HandlerFunc`** is the adapter: any `func(w, r)` becomes a Handler
- **Middleware pattern**: `func(next http.Handler) http.Handler` — the onion model
- **`httptest.NewRecorder()`** + **`httptest.NewRequest()`** let you test handlers without a running server
- **Header ordering matters**: set headers BEFORE `w.Write()` or `w.WriteHeader()`
- **`http.ServeMux`** (Go 1.22+) supports method-based routing: `mux.HandleFunc("GET /hello", ...)`
