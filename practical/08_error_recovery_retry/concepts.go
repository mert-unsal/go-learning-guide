// Package error_recovery_retry teaches how Go handles the equivalent of
// Java's try-catch-finally and retry patterns.
//
// Coming from Java you're used to:
//
//	try {
//	    riskyOperation();
//	} catch (SpecificException e) {
//	    handleSpecific(e);
//	} catch (Exception e) {
//	    handleGeneral(e);
//	} finally {
//	    cleanup();
//	}
//
// Go has NO try-catch. This is intentional — not an omission.
// Go's philosophy: errors are values, not control flow.
//
// This module covers:
//  1. defer/recover — Go's "catch" equivalent (for panics ONLY)
//  2. The recovery middleware pattern (production HTTP/gRPC handlers)
//  3. Retry with backoff — what Java does with Spring Retry / Resilience4j
//  4. Retryable vs permanent errors — classifying failures
//  5. Context-aware retry — respecting cancellation and deadlines
package error_recovery_retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

// ============================================================
// 1. DEFER / RECOVER — Go's "try-catch" equivalent
// ============================================================
//
// Java mental model mapping:
//   try        → the function body itself
//   catch      → defer func() { if r := recover(); r != nil { ... } }()
//   finally    → defer cleanup()
//   throw      → panic(value)
//
// CRITICAL DIFFERENCE:
// In Java, you throw exceptions for EXPECTED failures (IOException, SQLException).
// In Go, panic is for PROGRAMMER ERRORS ONLY — nil dereference, index out of range,
// impossible states. Expected failures use error return values.
//
// Rule of thumb:
//   - Java throws ~80% exceptions, ~20% unchecked
//   - Go returns ~99% errors, ~1% panics
//
// Under the hood:
// - panic() unwinds the stack, running deferred functions in LIFO order
// - recover() in a deferred function captures the panic value and resumes normal flow
// - If no recover() catches it, the goroutine prints stack trace and crashes the process
// - recover() ONLY works in deferred functions — calling it directly is a no-op

func DemonstrateTryCatchEquivalent() {
	fmt.Println("=== Go's try-catch equivalent ===")

	// This is the closest pattern to try-catch in Go:
	result, err := safeCall(func() interface{} {
		// This simulates a "risky operation" that panics
		panic("something went terribly wrong")
	})

	if err != nil {
		fmt.Println("Caught:", err) // "Caught: recovered panic: something went terribly wrong"
	} else {
		fmt.Println("Result:", result)
	}

	// Compare to Java:
	//   try { riskyOp(); } catch (Exception e) { System.out.println("Caught: " + e); }
}

// safeCall is the "try-catch wrapper" — converts panics to errors.
// This is the canonical Go pattern for calling code that might panic.
//
// Under the hood: defer runs AFTER the function returns (or panics).
// recover() returns the value passed to panic(), or nil if no panic occurred.
// Named return values (result, err) let the deferred function modify them.
func safeCall(fn func() interface{}) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert panic value to error
			switch v := r.(type) {
			case error:
				err = fmt.Errorf("recovered panic: %w", v)
			case string:
				err = fmt.Errorf("recovered panic: %s", v)
			default:
				err = fmt.Errorf("recovered panic: %v", v)
			}
		}
	}()

	result = fn()
	return result, nil
}

// ============================================================
// 2. RECOVERY MIDDLEWARE — Production pattern
// ============================================================
//
// In Java, your framework (Spring, Vert.x) catches exceptions in handlers.
// In Go, you write recovery middleware yourself. Every production Go HTTP
// server should have this.
//
// This is the pattern used by gin, echo, chi, and every serious Go framework.
//
// WHY: If a handler panics and you DON'T recover, the entire process crashes.
// With recovery middleware, the panic is caught, logged, and a 500 is returned.
//
// IMPORTANT: recover() only catches panics in the SAME GOROUTINE.
// If your handler spawns goroutines, each one needs its own recovery.

// RecoverFunc represents a function that wraps another function with panic recovery.
// In production, this would wrap http.HandlerFunc.
type RecoverFunc func(fn func()) error

// DemonstrateRecoveryMiddleware shows the middleware pattern.
func DemonstrateRecoveryMiddleware() {
	fmt.Println("\n=== Recovery Middleware ===")

	// A "handler" that panics
	handler := func() {
		panic("nil pointer in handler")
	}

	// Wrap with recovery — in production this is middleware
	err := recoverMiddleware(handler)
	if err != nil {
		fmt.Println("Handler recovered:", err)
		// In production: log.Error, set HTTP 500, increment panic metric
	}
}

func recoverMiddleware(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("handler panic: %v", r)
		}
	}()
	fn()
	return nil
}

// ============================================================
// 3. SIMPLE RETRY — Fixed attempts
// ============================================================
//
// Java equivalent: Spring's @Retryable or Resilience4j retry.
// Go: you write it yourself. It's ~10 lines. No framework needed.
//
// The simplest retry: try N times, return last error if all fail.

func DemonstrateSimpleRetry() {
	fmt.Println("\n=== Simple Retry ===")

	attempts := 0
	err := retryFixed(3, 100*time.Millisecond, func() error {
		attempts++
		if attempts < 3 {
			return fmt.Errorf("attempt %d failed", attempts)
		}
		return nil // succeeds on 3rd attempt
	})

	if err != nil {
		fmt.Println("All retries failed:", err)
	} else {
		fmt.Printf("Succeeded after %d attempts\n", attempts)
	}
}

// retryFixed retries fn up to maxAttempts times with a fixed delay between attempts.
func retryFixed(maxAttempts int, delay time.Duration, fn func() error) error {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if i < maxAttempts-1 {
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("after %d attempts: %w", maxAttempts, lastErr)
}

// ============================================================
// 4. EXPONENTIAL BACKOFF — Production-grade retry
// ============================================================
//
// Fixed delay is naive. In production, you want exponential backoff with jitter.
//
// Why? If your service has 1000 clients and the DB goes down, fixed retry
// means all 1000 hit the DB at the exact same intervals → thundering herd.
//
// Exponential backoff: delay = baseDelay * 2^attempt
// Jitter: add randomness so clients don't retry in lockstep
//
// Formula: delay = min(baseDelay * 2^attempt + random(0, jitter), maxDelay)
//
// This is the same algorithm used by AWS SDK, gRPC, Google Cloud client libraries.

// BackoffConfig configures exponential backoff retry.
type BackoffConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Multiplier  float64 // typically 2.0
}

// DefaultBackoff returns production-reasonable defaults.
func DefaultBackoff() BackoffConfig {
	return BackoffConfig{
		MaxAttempts: 5,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    10 * time.Second,
		Multiplier:  2.0,
	}
}

func DemonstrateExponentialBackoff() {
	fmt.Println("\n=== Exponential Backoff ===")

	cfg := DefaultBackoff()
	cfg.MaxAttempts = 4

	attempt := 0
	err := retryWithBackoff(cfg, func() error {
		attempt++
		fmt.Printf("  Attempt %d at %v\n", attempt, time.Now().Format("15:04:05.000"))
		if attempt < 3 {
			return fmt.Errorf("still failing")
		}
		return nil
	})

	if err != nil {
		fmt.Println("Failed:", err)
	} else {
		fmt.Println("Succeeded!")
	}
}

// retryWithBackoff retries with exponential backoff and jitter.
func retryWithBackoff(cfg BackoffConfig, fn func() error) error {
	var lastErr error
	for i := 0; i < cfg.MaxAttempts; i++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if i < cfg.MaxAttempts-1 {
			// Calculate delay: baseDelay * multiplier^attempt
			delay := time.Duration(float64(cfg.BaseDelay) * math.Pow(cfg.Multiplier, float64(i)))
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}

			// Add jitter: random 0-50% of delay
			jitter := time.Duration(rand.Int64N(int64(delay) / 2))
			time.Sleep(delay + jitter)
		}
	}
	return fmt.Errorf("after %d attempts: %w", cfg.MaxAttempts, lastErr)
}

// ============================================================
// 5. RETRYABLE vs PERMANENT ERRORS
// ============================================================
//
// Not all errors should be retried. In Java, you'd catch specific exception types.
// In Go, you classify errors using sentinel errors or custom error types.
//
// Retryable:  network timeout, 503 Service Unavailable, connection reset
// Permanent:  400 Bad Request, 404 Not Found, validation error, auth failure
//
// Pattern: wrap errors with a "retryable" marker.

// PermanentError wraps an error to indicate it should NOT be retried.
// This is the same pattern used by AWS SDK Go v2.
type PermanentError struct {
	Err error
}

func (e *PermanentError) Error() string { return e.Err.Error() }
func (e *PermanentError) Unwrap() error { return e.Err }

// Permanent marks an error as non-retryable.
func Permanent(err error) error {
	return &PermanentError{Err: err}
}

// IsPermanent checks if an error (or any wrapped error) is permanent.
func IsPermanent(err error) bool {
	var pe *PermanentError
	return errors.As(err, &pe)
}

func DemonstrateRetryableErrors() {
	fmt.Println("\n=== Retryable vs Permanent Errors ===")

	// Retry-aware function: stops immediately on permanent errors
	attempt := 0
	err := retrySmartly(5, 50*time.Millisecond, func() error {
		attempt++
		if attempt == 1 {
			return fmt.Errorf("network timeout") // retryable
		}
		// On 2nd attempt, permanent failure — don't retry
		return Permanent(fmt.Errorf("invalid API key"))
	})

	fmt.Printf("Stopped after %d attempts: %v\n", attempt, err)
}

// retrySmartly stops retrying on permanent errors.
func retrySmartly(maxAttempts int, delay time.Duration, fn func() error) error {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if IsPermanent(lastErr) {
			return lastErr // don't retry permanent errors
		}
		if i < maxAttempts-1 {
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("after %d attempts: %w", maxAttempts, lastErr)
}

// ============================================================
// 6. CONTEXT-AWARE RETRY — Respecting deadlines
// ============================================================
//
// In production, retries MUST respect context cancellation and deadlines.
// If the HTTP request's context is cancelled (client disconnected), stop retrying.
// If the context has a deadline, don't retry past it.
//
// This is the equivalent of Java's Future.cancel() or CompletableFuture with timeout.

func DemonstrateContextRetry() {
	fmt.Println("\n=== Context-Aware Retry ===")

	// Context with 500ms deadline — retries stop when deadline expires
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	attempt := 0
	err := retryWithContext(ctx, 10, 200*time.Millisecond, func() error {
		attempt++
		return fmt.Errorf("still failing (attempt %d)", attempt)
	})

	fmt.Printf("Stopped after %d attempts: %v\n", attempt, err)
	// Output: Stopped after ~3 attempts (context deadline exceeded)
}

// retryWithContext retries until success, maxAttempts, or context cancellation.
func retryWithContext(ctx context.Context, maxAttempts int, delay time.Duration, fn func() error) error {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		// Check if context is already done before attempting
		if ctx.Err() != nil {
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if i < maxAttempts-1 {
			// Use select to sleep OR cancel — whichever comes first
			select {
			case <-ctx.Done():
				return fmt.Errorf("retry cancelled during backoff: %w", ctx.Err())
			case <-time.After(delay):
				// continue to next attempt
			}
		}
	}
	return fmt.Errorf("after %d attempts: %w", maxAttempts, lastErr)
}

// ============================================================
// COMPARISON CHEAT SHEET: Java → Go
// ============================================================
//
// ┌────────────────────────┬──────────────────────────────────────────┐
// │ Java                   │ Go                                       │
// ├────────────────────────┼──────────────────────────────────────────┤
// │ try { ... }            │ (just write the code)                    │
// │ catch (Exception e)    │ defer func() { recover() }()             │
// │ finally { ... }        │ defer cleanup()                          │
// │ throw new XException() │ panic(value) (RARE — only for bugs)      │
// │ throws IOException     │ func f() error (return value, not throw) │
// │ e.getMessage()         │ err.Error()                              │
// │ instanceof             │ errors.As(err, &target)                  │
// │ e == SomeErr           │ errors.Is(err, sentinel)                 │
// │ try-with-resources     │ defer file.Close()                       │
// │ @Retryable             │ retryWithBackoff(cfg, fn) (write it)     │
// │ CircuitBreaker         │ github.com/sony/gobreaker (or write it)  │
// └────────────────────────┴──────────────────────────────────────────┘
//
// KEY INSIGHT: Go makes error handling EXPLICIT at every call site.
// This feels verbose coming from Java, but it makes control flow obvious.
// You never wonder "can this throw?" — if it returns error, handle it.

// RunAll runs all demonstrations.
func RunAll() {
	DemonstrateTryCatchEquivalent()
	DemonstrateRecoveryMiddleware()
	DemonstrateSimpleRetry()
	DemonstrateExponentialBackoff()
	DemonstrateRetryableErrors()
	DemonstrateContextRetry()
}
