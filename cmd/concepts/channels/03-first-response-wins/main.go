// Package main contains a standalone conceptual example for the First Response Wins pattern.
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
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

// backendColor returns a unique ANSI color tag for each backend.
func backendColor(backend string) string {
	switch backend {
	case "replica-us":
		return cyan
	case "replica-eu":
		return yellow
	case "replica-asia":
		return magenta
	default:
		return dim
	}
}

// QueryFastest fans out a query to multiple backends and returns
// the first successful response. Remaining backends are cancelled.
func QueryFastest(ctx context.Context, query string, backends []string, queryFn func(ctx context.Context, backend, query string) (string, error)) (QueryResult, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel remaining backends when first responds

	// Buffered channel: capacity = len(backends) to prevent goroutine leaks
	resultCh := make(chan QueryResult, len(backends))
	fmt.Printf("  %sresultCh buffer capacity = %s%d%s — prevents goroutine leaks if all respond%s\n",
		dim, magenta, len(backends), dim, reset)

	for _, b := range backends {
		go func(backend string) {
			clr := backendColor(backend)
			fmt.Printf("  %s[%s]%s 🏁 racing — query sent\n", clr+bold, backend, reset)
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
				fmt.Printf("  %s[%s]%s %s✔ responded in %v%s\n",
					clr+bold, backend, reset, green, time.Since(start), reset)
			default:
				// another result already sent — discard
				fmt.Printf("  %s[%s]%s %s… result discarded (winner already chosen)%s\n",
					clr+bold, backend, reset, dim, reset)
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
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  First Response Wins — Fan-Out Pattern           %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Pattern Overview%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Query ALL backends in parallel — take the fastest response%s\n", green, reset)
	fmt.Printf("  %s✔ context.WithCancel cancels slow backends after winner is chosen%s\n", green, reset)
	fmt.Printf("  %s✔ Buffered channel (cap=N) prevents goroutine leaks%s\n", green, reset)
	fmt.Printf("  %s⚠ Without buffered chan, losing goroutines block forever on send%s\n\n", yellow, reset)

	backends := []string{"replica-us", "replica-eu", "replica-asia"}
	fmt.Printf("%s▸ Backend Latencies%s\n", cyan+bold, reset)
	fmt.Printf("  %s[replica-us]%s    ~50ms  + jitter\n", cyan+bold, reset)
	fmt.Printf("  %s[replica-eu]%s    ~150ms + jitter\n", yellow+bold, reset)
	fmt.Printf("  %s[replica-asia]%s  ~80ms  + jitter\n\n", magenta+bold, reset)

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

		clr := backendColor(backend)
		select {
		case <-time.After(latency):
			return fmt.Sprintf("data-from-%s", backend), nil
		case <-ctx.Done():
			fmt.Printf("  %s[%s]%s %s⏹ cancelled (was going to take %v)%s\n",
				clr+bold, backend, reset, red, latency, reset)
			return "", ctx.Err()
		}
	}

	fmt.Printf("%s▸ Racing Backends%s\n", cyan+bold, reset)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result, err := QueryFastest(ctx, "SELECT * FROM users", backends, queryFn)
	if err != nil {
		fmt.Printf("  %s❌ All backends failed: %v%s\n", red+bold, err, reset)
		return
	}

	// Give cancelled backends time to print their cancellation
	time.Sleep(200 * time.Millisecond)

	clr := backendColor(result.Backend)
	fmt.Printf("\n%s▸ Result%s\n", cyan+bold, reset)
	fmt.Printf("  %s🏆 Winner: %s[%s]%s responded in %s%v%s\n",
		green+bold, clr+bold, result.Backend, reset, magenta, result.Duration, reset)
	fmt.Printf("  %sData: %s%s%s\n", dim, magenta, result.Data, reset)
	fmt.Printf("\n  %s✔ Slow backends were cancelled via context — no wasted work%s\n", green, reset)
	fmt.Printf("  %s⚠ This is Google's approach to tail latency (hedged requests)%s\n", yellow, reset)
}
