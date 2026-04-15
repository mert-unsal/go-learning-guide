package middleware

// ============================================================
// EXERCISES -- 04 middleware: HTTP Middleware Patterns
// ============================================================
// 12 exercises covering production middleware patterns.
// Focus: timing, CORS, request ID, compression, rate limiting.

import (
	"net/http"
	"time"
)

// Middleware is the standard middleware type signature.
type Middleware func(http.Handler) http.Handler

// ────────────────────────────────────────────────────────────
// Exercise 1: TimingMiddleware -- measure request duration
// ────────────────────────────────────────────────────────────
// Record request duration. Call onDone(duration) after the handler returns.
// Use time.Since(start) to measure.

func TimingMiddleware(onDone func(time.Duration)) Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: wrap with timing
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: CORSMiddleware -- set CORS headers
// ────────────────────────────────────────────────────────────
// Set Access-Control-Allow-Origin to the given origin.
// Set Access-Control-Allow-Methods to "GET, POST, PUT, DELETE, OPTIONS".
// Set Access-Control-Allow-Headers to "Content-Type, Authorization".
// For OPTIONS requests, return 204 without calling next.

func CORSMiddleware(allowedOrigin string) Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: wrap with CORS headers
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: RequestIDMiddleware -- generate/propagate request IDs
// ────────────────────────────────────────────────────────────
// If X-Request-ID header exists, use it. Otherwise, call generateID().
// Set the ID on both request context and response header.
// Store in context using RequestIDKey.

type contextKey string

const RequestIDKey contextKey = "request-id"

func RequestIDMiddleware(generateID func() string) Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: read/generate ID, set on context and response header
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: ContentTypeMiddleware -- enforce Content-Type
// ────────────────────────────────────────────────────────────
// For POST/PUT/PATCH requests, check Content-Type header.
// If it doesn't match required, return 415 Unsupported Media Type.
// GET/DELETE/OPTIONS pass through without checking.

func ContentTypeMiddleware(required string) Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: check Content-Type for write methods
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: MaxBytesMiddleware -- limit request body size
// ────────────────────────────────────────────────────────────
// Wrap r.Body with http.MaxBytesReader to limit body size.
// If the body exceeds maxBytes, the reader returns an error.

func MaxBytesMiddleware(maxBytes int64) Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: BasicAuthMiddleware -- HTTP Basic Auth
// ────────────────────────────────────────────────────────────
// Check Basic Auth credentials using r.BasicAuth().
// If invalid, set WWW-Authenticate header and return 401.

func BasicAuthMiddleware(username, password string) Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: check r.BasicAuth()
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 7: CacheControlMiddleware -- set cache headers
// ────────────────────────────────────────────────────────────
// Set Cache-Control header to the given value (e.g., "max-age=3600").

func CacheControlMiddleware(value string) Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: w.Header().Set("Cache-Control", value)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: SecurityHeadersMiddleware -- set security headers
// ────────────────────────────────────────────────────────────
// Set the following headers:
//   X-Content-Type-Options: nosniff
//   X-Frame-Options: DENY
//   X-XSS-Protection: 1; mode=block
//   Strict-Transport-Security: max-age=31536000; includeSubDomains

func SecurityHeadersMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: set all security headers
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 9: TimeoutMiddleware -- per-request timeout
// ────────────────────────────────────────────────────────────
// Wrap the request context with a timeout. If the handler doesn't
// finish in time, the context is cancelled.
// Note: This doesn't stop the handler — it only cancels the context.

func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: use http.TimeoutHandler or context.WithTimeout
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 10: MethodOverrideMiddleware -- X-HTTP-Method-Override
// ────────────────────────────────────────────────────────────
// If the request has X-HTTP-Method-Override header, change r.Method.
// This allows PUT/DELETE via POST (useful for HTML forms).

func MethodOverrideMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: check header, override r.Method
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 11: ConditionalMiddleware -- apply middleware conditionally
// ────────────────────────────────────────────────────────────
// Only apply the wrapped middleware if the condition function returns true.
// Otherwise, pass through to next directly.

func ConditionalMiddleware(condition func(*http.Request) bool, mw Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		return next // TODO: if condition(r) { mw(next).ServeHTTP(w, r) } else { next.ServeHTTP(w, r) }
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: ChainMiddleware -- compose middlewares
// ────────────────────────────────────────────────────────────
// Apply middlewares in order. First middleware is outermost.
// Same as the one in stdlib/10_net_http but using the Middleware type.

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	return handler // TODO: iterate in reverse, wrapping handler
}
