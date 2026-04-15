# 04 Middleware — HTTP Middleware Patterns

> **Companion chapter:** [learnings/27_net_http_under_the_hood.md](../../../learnings/27_net_http_under_the_hood.md)

## Exercises

| # | Middleware | Pattern | Difficulty |
|---|----------|---------|------------|
| 1 | `TimingMiddleware` | Request duration measurement | ⭐⭐ |
| 2 | `CORSMiddleware` | CORS headers, OPTIONS preflight | ⭐⭐ |
| 3 | `RequestIDMiddleware` | Generate/propagate request IDs | ⭐⭐⭐ |
| 4 | `ContentTypeMiddleware` | Enforce Content-Type | ⭐⭐ |
| 5 | `MaxBytesMiddleware` | Limit request body size | ⭐⭐ |
| 6 | `BasicAuthMiddleware` | HTTP Basic Auth | ⭐⭐ |
| 7 | `CacheControlMiddleware` | Cache headers | ⭐ |
| 8 | `SecurityHeadersMiddleware` | Security headers (OWASP) | ⭐⭐ |
| 9 | `TimeoutMiddleware` | Per-request timeout | ⭐⭐⭐ |
| 10 | `MethodOverrideMiddleware` | X-HTTP-Method-Override | ⭐⭐ |
| 11 | `ConditionalMiddleware` | Apply middleware conditionally | ⭐⭐⭐ |
| 12 | `Chain` | Compose middlewares | ⭐⭐⭐ |

## How to Practice

```bash
go test -race -v ./exercises/advanced/04_middleware/
go test -race -run TestCORSMiddleware ./exercises/advanced/04_middleware/
```

## Key Insights

- **Middleware signature**: `func(http.Handler) http.Handler`
- **Onion model**: first middleware in chain is outermost (executes first)
- **Header order**: set response headers BEFORE `w.Write()` or `w.WriteHeader()`
- **httptest**: test middleware without a running server
- **Security headers**: X-Content-Type-Options, X-Frame-Options, HSTS are production essentials
