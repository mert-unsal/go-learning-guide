package concepts

import (
	"fmt"
	"time"
)

// ============================================================
// Token Bucket Rate Limiter — Buffered Channel as Semaphore
// ============================================================
//
// The Problem:
//   Your API calls a downstream service (database, payment provider,
//   external API). That service can handle 100 requests/second. If you
//   send 500/sec, it crashes or throttles you. You need a gate that says:
//   "only 100 requests per second pass through."
//
// Real-world example:
//   API gateway rate limiting per client: each API key gets 100 req/sec.
//   Exceeding the limit returns HTTP 429 Too Many Requests immediately,
//   protecting both your infrastructure and the downstream service.
//
// The Pattern:
//   Token bucket using a buffered channel:
//     - Buffer capacity = maximum burst size (tokens the bucket holds)
//     - Each struct{}{} in the buffer = one token
//     - Taking a token = receiving from the channel (TryAcquire)
//     - Refilling = sending to the channel periodically (refill goroutine)
//
// Why channels work here:
//   1. Atomic take: <-tokens is goroutine-safe across all request handlers
//   2. Bounded capacity: buffer size IS the rate limit — enforced by the runtime
//   3. select+default gives non-blocking TryAcquire for free
//   4. The refill goroutine uses select+default to discard tokens when full
//      (the bucket never overflows past capacity)
//
// How the token flow works:
//
//   Refill goroutine                    Request handlers
//   ┌──────────────┐                    ┌──────────────┐
//   │ every 10ms:  │                    │ TryAcquire() │
//   │ tokens <- {} ├─── buffer(100) ────┤ <-tokens     │
//   │              │   [tok][tok][..]   │              │
//   │ bucket full? │                    │ empty?       │
//   │ └─ discard   │                    │ └─ reject    │
//   └──────────────┘                    └──────────────┘
//
// Timeline at 100 tokens/sec:
//   t=0.00s  Bucket: [100 tokens]
//            50 requests arrive → 50 take tokens → Bucket: [50]
//   t=0.01s  Refill adds 1 token → Bucket: [51]
//   t=0.02s  200 requests arrive → 51 served, 149 rejected
//   t=0.03s  Refill adds 1 → Bucket: [1]
//            ...steady state: ~100 requests served per second

// TokenBucket implements a rate limiter using a buffered channel.
type TokenBucket struct {
	tokens chan struct{}
	stop   chan struct{}
}

// NewTokenBucket creates a rate limiter that allows 'rate' operations per second.
// The bucket starts full, allowing an initial burst up to 'rate' requests.
func NewTokenBucket(rate int) *TokenBucket {
	tb := &TokenBucket{
		tokens: make(chan struct{}, rate),
		stop:   make(chan struct{}),
	}

	// Fill with initial tokens — allows burst up to 'rate'
	for i := 0; i < rate; i++ {
		tb.tokens <- struct{}{}
	}

	// Refill goroutine: adds one token every (1/rate) seconds
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rate))
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case tb.tokens <- struct{}{}:
					// token added to bucket
				default:
					// bucket already full, discard token
				}
			case <-tb.stop:
				return
			}
		}
	}()

	return tb
}

// TryAcquire attempts to take a token without blocking.
// Returns true if a token was available (proceed with request).
// Returns false if the bucket is empty (reject request, e.g., HTTP 429).
func (tb *TokenBucket) TryAcquire() bool {
	select {
	case <-tb.tokens:
		return true
	default:
		return false
	}
}

// Stop shuts down the refill goroutine.
func (tb *TokenBucket) Stop() {
	close(tb.stop)
}

// DemonstrateTokenBucketRateLimiter shows the rate limiter allowing an initial
// burst, rejecting excess requests, then allowing more after refill.
func DemonstrateTokenBucketRateLimiter() {
	limiter := NewTokenBucket(5) // 5 requests per second
	defer limiter.Stop()

	// Simulate 8 rapid requests — only first 5 should succeed (initial burst)
	fmt.Println("  Sending 8 rapid requests (rate limit: 5/sec):")
	for i := 1; i <= 8; i++ {
		if limiter.TryAcquire() {
			fmt.Printf("    Request %d: ✅ allowed\n", i)
		} else {
			fmt.Printf("    Request %d: ❌ rate limited\n", i)
		}
	}

	// Wait for tokens to refill
	fmt.Println("  ...waiting 1 second for refill...")
	time.Sleep(1 * time.Second)

	// Now tokens are available again
	fmt.Println("  After refill:")
	for i := 1; i <= 3; i++ {
		if limiter.TryAcquire() {
			fmt.Printf("    Request %d: ✅ allowed\n", i)
		}
	}
}
