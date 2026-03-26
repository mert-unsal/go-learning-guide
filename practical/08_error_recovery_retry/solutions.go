package error_recovery_retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

// ============================================================
// SOLUTIONS — 08 Error Recovery & Retry
// ============================================================

// SafeExecuteSolution converts panics to errors using defer + recover.
func SafeExecuteSolution(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = fmt.Errorf("panic: %w", v)
			default:
				err = fmt.Errorf("panic: %v", v)
			}
		}
	}()
	fn()
	return nil
}

// SafeGoroutineSolution launches fn with recovery, reports via buffered channel.
func SafeGoroutineSolution(fn func()) <-chan error {
	ch := make(chan error, 1) // buffered — goroutine never blocks
	go func() {
		defer func() {
			if r := recover(); r != nil {
				ch <- fmt.Errorf("goroutine panic: %v", r)
				return
			}
			ch <- nil
		}()
		fn()
	}()
	return ch
}

// RetrySolution retries fn up to maxAttempts with fixed delay.
func RetrySolution(maxAttempts int, delay time.Duration, fn func() error) error {
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
	return lastErr
}

// RetryWithBackoffSolution retries with exponential backoff.
func RetryWithBackoffSolution(cfg BackoffOptions, fn func() error) error {
	var lastErr error
	for i := 0; i < cfg.MaxAttempts; i++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if i < cfg.MaxAttempts-1 {
			delay := time.Duration(float64(cfg.BaseDelay) * math.Pow(2, float64(i)))
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
			time.Sleep(delay)
		}
	}
	return lastErr
}

// RetryClassifiedSolution stops retrying on permanent errors.
func RetryClassifiedSolution(maxAttempts int, delay time.Duration, fn func() error) error {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if IsPermanentError(lastErr) {
			return lastErr
		}
		if i < maxAttempts-1 {
			time.Sleep(delay)
		}
	}
	return lastErr
}

// RetryWithContextSolution respects context cancellation during retries.
func RetryWithContextSolution(ctx context.Context, maxAttempts int, delay time.Duration, fn func() error) error {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if i < maxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
	}
	return lastErr
}
