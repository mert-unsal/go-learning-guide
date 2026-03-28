// Package main contains a standalone conceptual example for the First Response Wins pattern.
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// ============================================================
// First Response Wins — Fan-Out to Multiple Backends
// ============================================================
//
// The Problem:
//   You have multiple backends that can answer the same query (database
//   replicas, CDN edge nodes, redundant API providers). Latency varies
//   per backend. Instead of picking one and hoping it's fast, you query
//   ALL of them in parallel and take whichever responds first.
//
// Real-world example:
//   Google's approach to tail latency: send the same read to multiple
//   Bigtable replicas, use the first response, cancel the rest.
//   Also used in DNS resolution (query multiple nameservers simultaneously).
//
// The Pattern:
//   Fan-out with buffered channel + context cancellation:
//     - Launch one goroutine per backend
//     - All goroutines send results to the same buffered channel
//     - Main goroutine takes the first result and cancels the rest
//
// Why channels work here:
//   1. Buffered channel (capacity = number of backends) prevents goroutine
//      leaks. Even if all backends respond after the first one is taken,
//      their sends succeed (into the buffer) and goroutines exit cleanly.
//
//   2. context.WithCancel propagates cancellation to slow backends.
//      When the first result arrives, cancel() tells remaining backends
//      to abort their work (if they respect ctx).
//
//   3. select on resultCh + ctx.Done() gives a clean timeout mechanism:
//      if NO backend responds before the parent context deadline, the
//      caller gets an error instead of hanging forever.
//
//   ┌────────────┐
//   │ Backend A  │──► 50ms  ──► resultCh ──► WINNER (first in)
//   ├────────────┤
//   │ Backend B  │──► 120ms ──► buffer absorbs, goroutine exits cleanly
//   ├────────────┤
//   │ Backend C  │──► cancel() fired, backend aborts via ctx
//   └────────────┘
//
//   The buffered channel is critical: without it, losing goroutines
//   would block forever on send (the main goroutine already moved on
//   after receiving the first result).

// QueryResult represents a response from a backend.
type QueryResult struct {
	Backend  string
	Data     string
	Duration time.Duration
}

// QueryFastest fans out a query to multiple backends and returns
// the first successful response. Remaining backends are cancelled.
func QueryFastest(ctx context.Context, query string, backends []string, queryFn func(ctx context.Context, backend, query string) (string, error)) (QueryResult, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel remaining backends when first responds

	// Buffered channel: capacity = len(backends) to prevent goroutine leaks
	resultCh := make(chan QueryResult, len(backends))

	for _, b := range backends {
		go func(backend string) {
			start := time.Now()
			data, err := queryFn(ctx, backend, query)
			if err != nil {
				return // backend failed or was cancelled — exit quietly
			}
			select {
			case resultCh <- QueryResult{
				Backend:  backend,
				Data:     data,
				Duration: time.Since(start),
			}:
			default:
				// another result already sent — discard
			}
		}(b)
	}

	// Wait for the first result or context cancellation
	select {
	case r := <-resultCh:
		return r, nil // first response wins
	case <-ctx.Done():
		return QueryResult{}, ctx.Err()
	}
}

func main() {
	backends := []string{"replica-us", "replica-eu", "replica-asia"}

	// Simulate backends with different latencies
	queryFn := func(ctx context.Context, backend, query string) (string, error) {
		latencies := map[string]time.Duration{
			"replica-us":   50 * time.Millisecond,
			"replica-eu":   150 * time.Millisecond,
			"replica-asia": 80 * time.Millisecond,
		}
		latency := latencies[backend]
		// Add some jitter
		jitter := time.Duration(rand.Intn(20)) * time.Millisecond
		latency += jitter

		select {
		case <-time.After(latency):
			return fmt.Sprintf("data-from-%s", backend), nil
		case <-ctx.Done():
			fmt.Printf("  ⏹ %s cancelled (was going to take %v)\n", backend, latency)
			return "", ctx.Err()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result, err := QueryFastest(ctx, "SELECT * FROM users", backends, queryFn)
	if err != nil {
		fmt.Printf("  ❌ All backends failed: %v\n", err)
		return
	}

	// Give cancelled backends time to print their cancellation
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("  🏆 Winner: %s responded in %v with: %s\n",
		result.Backend, result.Duration, result.Data)
}
