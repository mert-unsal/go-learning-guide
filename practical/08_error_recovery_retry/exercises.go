package error_recovery_retry

import (
	"context"
	"errors"
	"time"
)

// ============================================================
// EXERCISES — 08 Error Recovery & Retry
// ============================================================

// Exercise 1: SafeExecute — Go's try-catch equivalent
//
// LESSON: Convert panics to errors using defer + recover.
// This is the pattern you'll use anywhere you call code that might panic
// (e.g., third-party libraries, reflection, type assertions).
//
// Requirements:
//   - Run fn()
//   - If fn panics, recover and return the panic value as an error
//   - If fn succeeds, return nil
//   - Handle both string panics and error panics
func SafeExecute(fn func()) error {
	// TODO: implement with defer + recover
	panic("not implemented")
}

// Exercise 2: SafeGoroutine — Recovery in spawned goroutines
//
// LESSON: recover() only works in the SAME goroutine that panicked.
// If you spawn goroutines, each one needs its own recovery.
// This is a critical production pattern — an unrecovered goroutine panic
// crashes the ENTIRE process.
//
// Requirements:
//   - Launch fn in a new goroutine with panic recovery
//   - If fn panics, send the panic value (as error) to the returned channel
//   - If fn succeeds, send nil to the channel
//   - The channel should be buffered (size 1) so the goroutine doesn't leak
func SafeGoroutine(fn func()) <-chan error {
	// TODO: launch goroutine with recovery, report result via channel
	panic("not implemented")
}

// Exercise 3: Retry — Simple retry with fixed delay
//
// LESSON: The simplest retry pattern. Try up to maxAttempts times.
// Return nil on success, last error after all attempts exhausted.
//
// Requirements:
//   - Call fn up to maxAttempts times
//   - If fn returns nil, return nil immediately (don't keep retrying)
//   - Sleep 'delay' between attempts (but not after the last failed attempt)
//   - If all attempts fail, return the last error
func Retry(maxAttempts int, delay time.Duration, fn func() error) error {
	// TODO: implement retry loop
	panic("not implemented")
}

// Exercise 4: RetryWithBackoff — Exponential backoff retry
//
// LESSON: In production, fixed delay causes thundering herd.
// Use exponential backoff: delay doubles each attempt.
//
// Formula: delay = baseDelay * 2^attempt (capped at maxDelay)
//
// Requirements:
//   - Retry up to cfg.MaxAttempts times
//   - Delay between attempts: baseDelay * 2^attemptIndex (0-indexed)
//   - Cap delay at cfg.MaxDelay
//   - Return nil on first success, last error after exhaustion
func RetryWithBackoff(cfg BackoffOptions, fn func() error) error {
	// TODO: implement exponential backoff retry
	panic("not implemented")
}

// BackoffOptions configures the backoff retry behavior.
type BackoffOptions struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// Exercise 5: RetryClassified — Stop retrying on permanent errors
//
// LESSON: Not all errors are retryable. A 400 Bad Request should NOT be retried.
// A 503 Service Unavailable should. Classify errors to avoid wasting resources.
//
// Use the IsPermanentError() helper below to check.
//
// Requirements:
//   - Retry up to maxAttempts with fixed delay
//   - If fn returns nil, return nil immediately
//   - If fn returns a permanent error (IsPermanentError), return it immediately
//   - Otherwise retry. If all attempts fail, return the last error.
func RetryClassified(maxAttempts int, delay time.Duration, fn func() error) error {
	// TODO: implement retry that stops on permanent errors
	panic("not implemented")
}

// RetryPermanentError marks an error as non-retryable.
// Wrap any error with this to signal "don't retry".
type RetryPermanentError struct {
	Err error
}

func (e *RetryPermanentError) Error() string { return e.Err.Error() }
func (e *RetryPermanentError) Unwrap() error { return e.Err }

// MarkPermanent wraps err to indicate it should not be retried.
func MarkPermanent(err error) error {
	return &RetryPermanentError{Err: err}
}

// IsPermanentError checks whether err (or any wrapped error) is permanent.
func IsPermanentError(err error) bool {
	var pe *RetryPermanentError
	return errors.As(err, &pe)
}

// Exercise 6: RetryWithContext — Context-aware retry
//
// LESSON: Production retries MUST respect context cancellation.
// If the caller's context is cancelled (HTTP client disconnected, shutdown signal),
// stop retrying immediately. Use select + time.After to make sleep cancellable.
//
// Requirements:
//   - Check ctx.Err() before each attempt — if cancelled, return ctx.Err()
//   - Call fn. If nil, return nil
//   - Between attempts, use select to either sleep (time.After) or detect ctx.Done()
//   - If context expires during sleep, return ctx.Err()
//   - If all attempts fail, return the last error
func RetryWithContext(ctx context.Context, maxAttempts int, delay time.Duration, fn func() error) error {
	// TODO: implement context-aware retry
	panic("not implemented")
}
