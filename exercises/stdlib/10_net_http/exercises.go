package net_http

// ============================================================
// EXERCISES -- 10 net/http: Handlers, Middleware & Servers
// ============================================================
// 12 exercises covering Go's net/http at production depth.
// Focus: Handler interface, ServeMux, middleware, testing.

import (
	"net/http"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: HelloHandler -- basic http.Handler
// ────────────────────────────────────────────────────────────
// Write a handler that responds with "Hello, World!" and status 200.
// The Handler interface: ServeHTTP(w http.ResponseWriter, r *http.Request)

type HelloHandler struct{}

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: w.WriteHeader(200); w.Write([]byte("Hello, World!"))
}

// ────────────────────────────────────────────────────────────
// Exercise 2: GreetHandler -- read URL query parameter
// ────────────────────────────────────────────────────────────
// GET /greet?name=Go → "Hello, Go!"
// If name is missing, respond with "Hello, stranger!"

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: r.URL.Query().Get("name")
}

// ────────────────────────────────────────────────────────────
// Exercise 3: EchoHandler -- read and echo request body
// ────────────────────────────────────────────────────────────
// Read the entire request body and write it back as the response.
// Always close the body.

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: io.ReadAll(r.Body), defer r.Body.Close(), w.Write(body)
	_ = r
}

// ────────────────────────────────────────────────────────────
// Exercise 4: MethodRouter -- dispatch by HTTP method
// ────────────────────────────────────────────────────────────
// GET → respond "GET OK"
// POST → respond "POST OK"
// Everything else → 405 Method Not Allowed

func MethodRouter(w http.ResponseWriter, r *http.Request) {
	// TODO: switch r.Method { case http.MethodGet: ... }
}

// ────────────────────────────────────────────────────────────
// Exercise 5: JSONResponse -- set Content-Type and write JSON
// ────────────────────────────────────────────────────────────
// Respond with {"status":"ok"} and Content-Type: application/json

func JSONResponse(w http.ResponseWriter, r *http.Request) {
	// TODO: w.Header().Set("Content-Type", "application/json")
	// Then write the JSON bytes
}

// ────────────────────────────────────────────────────────────
// Exercise 6: StatusHandler -- custom status codes
// ────────────────────────────────────────────────────────────
// Read ?code=404 from query, respond with that status code.
// Body should be the status text: "Not Found", "OK", etc.
// Invalid code defaults to 200.

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: strconv.Atoi(r.URL.Query().Get("code"))
	// http.StatusText(code) returns the text for a status code
}

// ────────────────────────────────────────────────────────────
// Exercise 7: SetupMux -- configure a ServeMux with routes
// ────────────────────────────────────────────────────────────
// Register these routes on the mux:
//   GET /hello  → HelloHandler
//   GET /greet  → GreetHandler (use http.HandlerFunc adapter)
//   POST /echo  → EchoHandler
// Return the configured mux.

func SetupMux() *http.ServeMux {
	mux := http.NewServeMux()
	// TODO: mux.Handle("/hello", HelloHandler{})
	// mux.HandleFunc("/greet", GreetHandler)
	return mux
}

// ────────────────────────────────────────────────────────────
// Exercise 8: LoggingMiddleware -- wrap handler with logging
// ────────────────────────────────────────────────────────────
// Return a handler that:
//  1. Calls logger(r.Method, r.URL.Path) before the request
//  2. Calls next.ServeHTTP(w, r)
// The logger function is injected for testability.

func LoggingMiddleware(next http.Handler, logger func(method, path string)) http.Handler {
	// TODO: return http.HandlerFunc(func(w, r) { logger(...); next.ServeHTTP(w, r) })
	return next
}

// ────────────────────────────────────────────────────────────
// Exercise 9: RecoveryMiddleware -- catch panics
// ────────────────────────────────────────────────────────────
// Wrap the handler so that if it panics, the middleware:
//  1. Recovers the panic
//  2. Responds with 500 Internal Server Error
//  3. Calls onPanic(recovered) for observability

func RecoveryMiddleware(next http.Handler, onPanic func(v interface{})) http.Handler {
	// TODO: defer func() { if r := recover(); r != nil { onPanic(r); w.WriteHeader(500) } }()
	return next
}

// ────────────────────────────────────────────────────────────
// Exercise 10: AuthMiddleware -- check for Authorization header
// ────────────────────────────────────────────────────────────
// If r.Header.Get("Authorization") == expectedToken, call next.
// Otherwise respond with 401 Unauthorized.

func AuthMiddleware(next http.Handler, expectedToken string) http.Handler {
	// TODO: check header, respond 401 if missing/wrong
	return next
}

// ────────────────────────────────────────────────────────────
// Exercise 11: ChainMiddleware -- compose middleware functions
// ────────────────────────────────────────────────────────────
// Apply middlewares in order: first in the list wraps outermost.
// ChainMiddleware(handler, mw1, mw2, mw3) → mw1(mw2(mw3(handler)))
// This is the standard "onion" pattern.

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
	// TODO: iterate backwards, wrapping handler each time
	return handler
}

// ────────────────────────────────────────────────────────────
// Exercise 12: HeadersHandler -- read and set custom headers
// ────────────────────────────────────────────────────────────
// Read X-Request-ID from request headers.
// Set X-Request-ID on response (echo it back).
// Set X-Powered-By: "go-learning-guide" on response.
// Body: the request ID value (or "none" if header missing).

func HeadersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: r.Header.Get("X-Request-ID"), w.Header().Set(...)
}
