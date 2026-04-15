package context_exercises

import (
	"context"
	"time"
)

// ============================================================
// EXERCISES — 08 Context: Cancellation, Timeouts & Values
// ============================================================
//
// context.Context is a 4-method interface: Deadline, Done, Err, Value.
// It forms a linked-list chain via the decorator pattern:
//   Background → WithValue → WithCancel → WithValue → ...
//
// These exercises test your understanding of:
//   - Context value storage and O(n) lookup (§2-3 of chapter 19)
//   - Cancellation propagation: parent→child, not child→parent (§3)
//   - Timeout/deadline behavior (§3)
//   - Production patterns: select+Done, early bailout, typed keys (§7)
//   - The Context interface itself (§1)
//
// Exercises 1-3:   Context values — typed keys, chain lookup
// Exercises 4-6:   Cancellation — propagation rules
// Exercises 7-8:   Timeouts and deadlines
// Exercises 9-12:  Production patterns
// ============================================================

// ── Types used across exercises ──

// RequestKey is a private type for context keys to avoid collisions.
// Using string keys like context.WithValue(ctx, "userID", 42) risks
// collisions between packages. A private type is uncollisionable.
type RequestKey struct{}
type TraceKey struct{}

// ── Exercise 1 ──

// WithRequestID returns a new context with the given request ID stored
// under RequestKey{}. GetRequestID retrieves it.
//
// KEY PATTERN: Always use a private type (not string) as context key.
// Two packages both using "userID" as a string key would collide.
// A private type is unique to your package — impossible to collide.
//
// See learnings/19 §7 — "Don't store context in structs"
func WithRequestID(ctx context.Context, id string) context.Context {
	// TODO: return context.WithValue(ctx, RequestKey{}, id)
	return ctx
}

func GetRequestID(ctx context.Context) string {
	// TODO: retrieve value for RequestKey{}, type-assert to string
	return ""
}

// ── Exercise 2 ──

// ChainValues builds a context chain by adding each key-value pair from
// the pairs slice (in order) onto context.Background().
// Returns the final context.
//
// INSIGHT: Each WithValue creates a new node in a linked list.
// 5 calls = 5 nodes. Value lookup walks the chain — O(n).
// This is why you keep context values shallow.
//
// See learnings/19 §3 — "How the Chain Works"
func ChainValues(pairs [][2]string) context.Context {
	// TODO: start with context.Background(), loop over pairs,
	//       add each with context.WithValue (use the string key directly here)
	return context.Background()
}

// LookupAll retrieves all values for the given keys from the context.
// Returns a map of key→value for keys that exist. Missing keys are omitted.
func LookupAll(ctx context.Context, keys []string) map[string]string {
	// TODO: for each key, ctx.Value(key), type-assert to string, add to map if found
	return nil
}

// ── Exercise 3 ──

// CancelAndCheck creates a cancellable context from Background(),
// cancels it, and returns ctx.Err().
//
// INSIGHT: After cancel(), ctx.Done() channel is closed and
// ctx.Err() returns context.Canceled.
//
// See learnings/19 §3 — "What Happens When You Call cancel()"
func CancelAndCheck() error {
	// TODO: ctx, cancel := context.WithCancel(context.Background())
	//       cancel()
	//       return ctx.Err()
	return nil
}

// ── Exercise 4 ──

// ParentCancelsChild creates a parent context with cancel, derives a
// child context (also with cancel), cancels the PARENT, and returns
// both errors. Both should be context.Canceled.
//
// RULE: Cancellation propagates DOWNWARD — parent→child.
// When parent is cancelled, all children are also cancelled.
//
// See learnings/19 §7 — "Cancellation is One-Way"
func ParentCancelsChild() (parentErr, childErr error) {
	// TODO: parent, parentCancel := context.WithCancel(context.Background())
	//       child, _ := context.WithCancel(parent)
	//       parentCancel()
	//       return parent.Err(), child.Err()
	return nil, nil
}

// ── Exercise 5 ──

// ChildDoesNotCancelParent creates a parent and child context,
// cancels only the CHILD, and returns both errors.
// Parent should be nil (still alive), child should be Canceled.
//
// RULE: Cancellation does NOT propagate UPWARD — child→parent never happens.
// This is critical for request trees where subrequests may fail independently.
func ChildDoesNotCancelParent() (parentErr, childErr error) {
	// TODO: parent, _ := context.WithCancel(context.Background())
	//       child, childCancel := context.WithCancel(parent)
	//       childCancel()
	//       return parent.Err(), child.Err()
	return nil, nil
}

// ── Exercise 6 ──

// TimeoutExpired creates a context with the given timeout, waits for
// it to expire, and returns ctx.Err().
//
// INSIGHT: context.WithTimeout internally creates a timerCtx that
// fires after the duration. Once expired, Done() closes and
// Err() returns context.DeadlineExceeded (not Canceled).
//
// Always defer cancel() even with timeouts — it releases the timer early.
func TimeoutExpired(timeout time.Duration) error {
	// TODO: ctx, cancel := context.WithTimeout(context.Background(), timeout)
	//       defer cancel()
	//       <-ctx.Done()    // block until timeout fires
	//       return ctx.Err()
	return nil
}

// ── Exercise 7 ──

// CheckDeadline creates a context with the given timeout and returns
// the deadline and whether one is set.
//
// INSIGHT: WithTimeout(ctx, d) is shorthand for WithDeadline(ctx, time.Now().Add(d)).
// The timerCtx overrides Deadline() to return the actual deadline time.
// Background() and cancelCtx do NOT have deadlines — they return (zero, false).
func CheckDeadline(timeout time.Duration) (deadline time.Time, ok bool) {
	// TODO: ctx, cancel := context.WithTimeout(context.Background(), timeout)
	//       defer cancel()
	//       return ctx.Deadline()
	return time.Time{}, false
}

// ── Exercise 8 ──

// ProcessItems iterates over items and squares each one.
// Before processing each item, check ctx.Err() — if the context
// is cancelled, stop immediately and return what was processed so far.
//
// PRODUCTION PATTERN: Every long-running loop should check ctx.Err()
// before expensive operations. This enables graceful cancellation
// without goroutine leaks.
//
// See learnings/19 §7 — "Always Check ctx.Err() Before Expensive Operations"
func ProcessItems(ctx context.Context, items []int) []int {
	// TODO: for each item, check ctx.Err() first, then append item*item
	return nil
}

// ── Exercise 9 ──

// SelectWithContext demonstrates the fundamental context+channel select
// pattern. It reads one value from ch or returns an error if ctx is
// done first.
//
// THE PATTERN:
//   select {
//   case <-ctx.Done():
//       return zero, ctx.Err()
//   case val := <-ch:
//       return val, nil
//   }
//
// This is THE most common context usage in production Go code.
// Every database call, HTTP call, and RPC call uses this pattern internally.
func SelectWithContext(ctx context.Context, ch <-chan string) (string, error) {
	// TODO: use select with ctx.Done() and ch
	return "", nil
}

// ── Exercise 10 ──

// FirstResult launches a goroutine for each task function, returns the
// first result, and cancels the context to stop remaining goroutines.
//
// PATTERN: Fan-out with cancellation. Launch N goroutines, first one to
// finish sends its result, cancel() cleans up the rest.
// This is used in: first-response-wins load balancing, speculative execution,
// racing multiple data sources.
//
// IMPORTANT: Use a buffered channel (cap = len(tasks)) so goroutines
// that finish after the first don't block forever on send.
func FirstResult(tasks []func() int) int {
	// TODO: ctx, cancel := context.WithCancel(context.Background())
	//       defer cancel()
	//       ch := make(chan int, len(tasks))
	//       launch goroutines, return first result from ch
	return 0
}

// ── Exercise 11 ──

// NestedTimeout creates a parent context with timeout1 and a child
// context with timeout2. Returns the effective deadline of the child.
//
// RULE: The effective deadline is always the EARLIER of parent and child.
// If parent expires in 5s and child in 10s, the child's effective
// deadline is 5s (limited by parent). WithTimeout won't extend past parent.
//
// See learnings/19 §3 — timerCtx inherits parent's deadline if shorter
func NestedTimeout(timeout1, timeout2 time.Duration) time.Duration {
	// TODO: parent, cancel1 := context.WithTimeout(Background(), timeout1)
	//       child, cancel2 := context.WithTimeout(parent, timeout2)
	//       defer cancel1(); defer cancel2()
	//       get child's deadline, return time.Until(deadline)
	return 0
}

// ── Exercise 12 ──

// AlwaysCancelled returns a context.Context implementation that is
// always in the cancelled state. Done() returns a closed channel,
// Err() returns context.Canceled, Deadline() returns zero/false,
// Value() returns nil.
//
// WHY: This is a test double. In tests, you pass this to any function
// that accepts context.Context to immediately trigger the cancellation path.
// No need for real timeouts or cancel functions — instant, deterministic.
//
// See learnings/19 §4 — "Testing Becomes Rigid" (struct) vs this approach
func AlwaysCancelled() context.Context {
	// TODO: define a type that implements the 4 Context methods
	//       Done() returns a pre-closed channel
	//       Err() returns context.Canceled
	return context.Background()
}
