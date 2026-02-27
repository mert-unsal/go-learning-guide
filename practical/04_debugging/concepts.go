// Package debugging covers how to debug Go applications.
//
// ============================================================
// DEBUGGING GO APPLICATIONS — COMPLETE GUIDE
// ============================================================
//
// ─────────────────────────────────────────────────────────────
// 1. PRINT DEBUGGING (simplest)
// ─────────────────────────────────────────────────────────────
//   fmt.Println("value:", x)
//   fmt.Printf("type=%T value=%+v\n", x, x)  // %T=type, %+v=struct with fields
//   fmt.Printf("%#v\n", x)                    // Go-syntax representation
//
// ─────────────────────────────────────────────────────────────
// 2. DELVE DEBUGGER — the standard Go debugger
// ─────────────────────────────────────────────────────────────
//   Install:
//     go install github.com/go-delve/delve/cmd/dlv@latest
//
//   ── Start debugging ───────────────────────────────────────
//   dlv debug .                        ← compile + start debug session
//   dlv debug ./cmd/myapp/             ← debug specific package
//   dlv debug . -- --port=8080         ← pass args to your program
//   dlv exec ./myapp                   ← debug an existing binary
//   dlv test ./...                     ← debug tests
//   dlv attach <pid>                   ← attach to running process
//
//   ── Common Delve commands (inside dlv session) ────────────
//   break main.main        ← set breakpoint at function
//   break main.go:42       ← set breakpoint at file:line
//   b main.go:42           ← shorthand for break
//   condition 1 x > 10     ← conditional breakpoint
//
//   continue   (c)         ← run until next breakpoint
//   next       (n)         ← step over (next line)
//   step       (s)         ← step into function call
//   stepout    (so)        ← step out of current function
//
//   print x    (p x)       ← print variable
//   locals                 ← print all local variables
//   vars                   ← print all package-level variables
//   args                   ← print function arguments
//   goroutines             ← list all goroutines
//   goroutine 5            ← switch to goroutine 5
//   stack                  ← print call stack
//
//   list                   ← show current source location
//   list main.go:42        ← show source at line
//
//   set x = 100            ← change variable value at runtime
//
//   breakpoints (bp)       ← list all breakpoints
//   clear 1                ← remove breakpoint #1
//   clearall               ← remove all breakpoints
//
//   exit  (q)              ← quit debugger
//   help                   ← list all commands
//
// ─────────────────────────────────────────────────────────────
// 3. GOLAND / VSCODE DEBUGGER (GUI)
// ─────────────────────────────────────────────────────────────
//   GoLand:
//   • Click the gutter (left of line numbers) to set a breakpoint
//   • Click the green bug icon ▼ or Shift+F9 to start debugger
//   • F8 = step over, F7 = step into, Shift+F8 = step out
//   • Variables panel shows all locals automatically
//   • Evaluate Expression: Alt+F8 → type any expression
//
//   VSCode (needs Go extension + dlv):
//   • Add to .vscode/launch.json:
//     {
//       "version": "0.2.0",
//       "configurations": [
//         {
//           "name": "Launch",
//           "type": "go",
//           "request": "launch",
//           "mode": "auto",
//           "program": "${workspaceFolder}/cmd/myapp",
//           "env": { "APP_ENV": "development" },
//           "args": ["--port=8080"]
//         }
//       ]
//     }
//   • Press F5 to start debugging
//   • F10 = step over, F11 = step into, Shift+F11 = step out
//
// ─────────────────────────────────────────────────────────────
// 4. RACE DETECTOR — find concurrency bugs
// ─────────────────────────────────────────────────────────────
//   go run -race .
//   go test -race ./...
//   go build -race -o myapp .     ← builds race-instrumented binary
//
//   The race detector reports when two goroutines access the same
//   memory concurrently and at least one is a write.
//
// ─────────────────────────────────────────────────────────────
// 5. PPROF — CPU and memory profiling
// ─────────────────────────────────────────────────────────────
//   Add to your HTTP server (import _ "net/http/pprof"):
//
//     import (
//         "net/http"
//         _ "net/http/pprof"  // registers /debug/pprof/ routes
//     )
//
//     func main() {
//         go func() {
//             http.ListenAndServe("localhost:6060", nil)
//         }()
//         // ... rest of your app
//     }
//
//   Then:
//     go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
//     go tool pprof http://localhost:6060/debug/pprof/heap
//
//   Inside pprof:
//     top10     ← top 10 functions by CPU/memory
//     web       ← open flame graph in browser (needs graphviz)
//     list func ← show annotated source
//
// ─────────────────────────────────────────────────────────────
// 6. go vet — static analysis
// ─────────────────────────────────────────────────────────────
//   go vet ./...        ← catches common mistakes
//
//   Examples it catches:
//   • fmt.Printf format string mismatches
//   • Unreachable code
//   • Incorrect mutex usage
//   • Suspicious composite literals
//
// ─────────────────────────────────────────────────────────────
// 7. LOGGING FOR DEBUGGING (structured logging)
// ─────────────────────────────────────────────────────────────
//   Standard library:
//     log.Printf("user=%s action=%s", userID, action)
//     log.Fatalf("failed to start: %v", err)  ← logs + os.Exit(1)
//
//   slog (Go 1.21+ built-in structured logger):
//     slog.Info("server started", "port", 8080)
//     slog.Error("db error", "err", err, "query", q)
//     slog.Debug("request", "method", r.Method, "path", r.URL.Path)
//
//   Enable debug level:
//     slog.SetLogLoggerLevel(slog.LevelDebug)
//
//   JSON output (for log aggregators like Datadog, Splunk):
//     logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
//     logger.Info("processing", "jobID", 42, "status", "ok")
//
// ─────────────────────────────────────────────────────────────
// 8. COMMON DEBUGGING TIPS
// ─────────────────────────────────────────────────────────────
//   • panic prints a full goroutine stack trace → very useful
//   • Add recover() in deferred function to catch panics gracefully
//   • Use errors.Is / errors.As to inspect wrapped errors
//   • fmt.Errorf("context: %w", err) wraps errors with context
//   • Use os.Exit(1) NOT panic in main for fatal startup errors
//   • Table-driven tests catch edge cases before you need the debugger

package debugging
