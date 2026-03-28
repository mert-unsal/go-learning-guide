// Package channels covers Go channels: buffered vs unbuffered,
// select, done patterns, fan-out/fan-in, and common idioms.
package channels

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================
// 1. CHANNEL BASICS
// ============================================================
// A channel is a typed conduit for communication between goroutines.
// make(chan T)       → unbuffered: sender blocks until receiver is ready
// make(chan T, n)    → buffered:   sender blocks only when buffer is full
// close(ch)         → signals no more values will be sent
// <-ch              → receive
// ch <- val         → send

func DemonstrateBasics() {
	// Unbuffered channel — synchronization point
	ch := make(chan int)

	go func() {
		ch <- 42 // blocks until someone receives
	}()

	val := <-ch // blocks until something is sent
	fmt.Println("Received:", val)

	// Buffered channel — asynchronous up to buffer size
	bch := make(chan string, 3)
	bch <- "one" // doesn't block (buffer not full)
	bch <- "two"
	bch <- "three"
	// bch <- "four" // would block — buffer is full

	fmt.Println(<-bch) // "one"
	fmt.Println(<-bch) // "two"
	fmt.Println(<-bch) // "three"
}

// ============================================================
// 2. RANGE OVER CHANNEL
// ============================================================
// range over a channel receives values until it's CLOSED.
// Always close channels from the SENDER side.

func generate(nums ...int) <-chan int { // returns a receive-only channel
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out) // signal that no more values will be sent
	}()
	return out
}

func DemonstrateRange() {
	for n := range generate(2, 3, 4, 5) {
		fmt.Print(n*n, " ") // 4 9 16 25
	}
	fmt.Println()
}

// ============================================================
// 3. SELECT — multiplex channels
// ============================================================
// select waits on multiple channel operations.
// Picks a random case if multiple are ready.

func DemonstrateSelect() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Millisecond)
		ch1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Millisecond)
		ch2 <- "two"
	}()

	// Receive from whichever is ready first
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received from ch1:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received from ch2:", msg2)
		}
	}
}

// ============================================================
// 4. DONE PATTERN (Cancellation)
// ============================================================
// Use a 'done' channel to signal goroutines to stop.

func DemonstrateDonePattern() {
	done := make(chan struct{})
	results := make(chan int)

	go func() {
		i := 0
		for {
			select {
			case <-done: // cancellation signal
				close(results)
				return
			case results <- i: // send work
				i++
			}
		}
	}()

	// Receive first 5 results
	for i := 0; i < 5; i++ {
		fmt.Print(<-results, " ")
	}
	fmt.Println()

	close(done) // signal goroutine to stop
	time.Sleep(1 * time.Millisecond)
}

// ============================================================
// 5. FAN-OUT / FAN-IN PATTERN
// ============================================================

// Fan-out: distribute work across multiple goroutines
// Fan-in: merge results from multiple goroutines into one channel

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	merged := make(chan int)

	// Start output goroutine for each input channel
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			merged <- n
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	// Start a goroutine to close merged once all inputs are done
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// ============================================================
// 6. TIMEOUT PATTERN
// ============================================================

func fetchWithTimeout(duration time.Duration) (string, error) {
	result := make(chan string, 1)

	go func() {
		// Simulate slow operation
		time.Sleep(duration * 2)
		result <- "data"
	}()

	select {
	case data := <-result:
		return data, nil
	case <-time.After(duration):
		return "", fmt.Errorf("timeout after %v", duration)
	}
}

func DemonstrateTimeout() {
	_, err := fetchWithTimeout(10 * time.Millisecond)
	if err != nil {
		fmt.Println("Timeout:", err)
	}
}

// ============================================================
// 7. PIPELINE PATTERN
// ============================================================

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func DemonstratePipeline() {
	// Pipeline: generate → square → square
	c := generate(2, 3, 4)
	sq1 := square(c)
	sq2 := square(sq1) // 2^4=16, 3^4=81, 4^4=256

	for n := range sq2 {
		fmt.Print(n, " ")
	}
	fmt.Println()
}

// ============================================================
// 8. TOKEN BUCKET RATE LIMITER
// ============================================================
// A buffered channel naturally models a token bucket:
//   - Buffer capacity = maximum burst size (how many tokens the bucket holds)
//   - Each struct{}{} in the buffer = one token
//   - Taking a token = receiving from the channel
//   - Refilling = sending to the channel periodically
//
// The problem this solves:
//   Your API calls a downstream service that handles 100 req/sec.
//   If you send 500/sec, it crashes. The rate limiter acts as a gate:
//   only requests that acquire a token pass through.
//
// Why a channel works here:
//   - Atomic take: <-tokens is goroutine-safe, no mutex needed
//   - Bounded capacity: buffer size IS the rate limit
//   - select+default gives non-blocking tryAcquire for free
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

func DemonstrateRateLimiter() {
	limiter := NewTokenBucket(5) // 5 requests per second
	defer limiter.Stop()

	// Simulate 8 rapid requests — only first 5 should succeed (initial burst)
	for i := 1; i <= 8; i++ {
		if limiter.TryAcquire() {
			fmt.Printf("  Request %d: ✅ allowed\n", i)
		} else {
			fmt.Printf("  Request %d: ❌ rate limited\n", i)
		}
	}

	// Wait for tokens to refill
	fmt.Println("  ...waiting 1 second for refill...")
	time.Sleep(1 * time.Second)

	// Now tokens are available again
	if limiter.TryAcquire() {
		fmt.Println("  After refill: ✅ allowed")
	}
}

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Channel Basics ===")
	DemonstrateBasics()
	fmt.Println("\n=== Range Over Channel ===")
	DemonstrateRange()
	fmt.Println("\n=== Select ===")
	DemonstrateSelect()
	fmt.Println("\n=== Done Pattern ===")
	DemonstrateDonePattern()
	fmt.Println("\n=== Timeout ===")
	DemonstrateTimeout()
	fmt.Println("\n=== Pipeline ===")
	DemonstratePipeline()
	fmt.Println("\n=== Token Bucket Rate Limiter ===")
	DemonstrateRateLimiter()
}
