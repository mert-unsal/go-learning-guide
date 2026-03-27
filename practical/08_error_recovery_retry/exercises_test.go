package error_recovery_retry

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================
// TESTS — 08 Error Recovery & Retry
// ============================================================

func TestSafeExecute_NoPanic(t *testing.T) {
	err := SafeExecute(func() {
		// does nothing — should not error
	})
	if err != nil {
		t.Errorf("❌ SafeExecute(no panic) = %v, want nil\n"+
			"   Hint: when fn() doesn't panic, recover() returns nil — return nil", err)
	} else {
		t.Log("✅ SafeExecute(no panic) = nil")
	}
}

func TestSafeExecute_StringPanic(t *testing.T) {
	err := SafeExecute(func() {
		panic("something broke")
	})
	if err == nil {
		t.Error("❌ SafeExecute(string panic) = nil, want error\n" +
			"   Hint: defer func() { if r := recover(); r != nil { ... } }()")
	} else if !strings.Contains(err.Error(), "something broke") {
		t.Errorf("❌ SafeExecute error = %q, should contain \"something broke\"\n"+
			"   Hint: convert the recover() value to an error", err)
	} else {
		t.Logf("✅ SafeExecute(string panic) caught: %v", err)
	}
}

func TestSafeExecute_ErrorPanic(t *testing.T) {
	originalErr := errors.New("original error")
	err := SafeExecute(func() {
		panic(originalErr)
	})
	if err == nil {
		t.Error("❌ SafeExecute(error panic) = nil, want error")
	} else if !strings.Contains(err.Error(), "original error") {
		t.Errorf("❌ SafeExecute error = %q, should contain \"original error\"", err)
	} else {
		t.Logf("✅ SafeExecute(error panic) caught: %v", err)
	}
}

func TestSafeGoroutine_NoPanic(t *testing.T) {
	ch := SafeGoroutine(func() {
		// success — no panic
	})

	select {
	case err := <-ch:
		if err != nil {
			t.Errorf("❌ SafeGoroutine(success) sent error: %v, want nil\n"+
				"   Hint: send nil on the channel when fn() completes successfully", err)
		} else {
			t.Log("✅ SafeGoroutine(success) = nil")
		}
	case <-time.After(2 * time.Second):
		t.Error("❌ SafeGoroutine(success) timed out — did you send a result on the channel?\n" +
			"   Hint: use a buffered channel (make(chan error, 1)) so the goroutine doesn't block")
	}
}

func TestSafeGoroutine_WithPanic(t *testing.T) {
	ch := SafeGoroutine(func() {
		panic("goroutine exploded")
	})

	select {
	case err := <-ch:
		if err == nil {
			t.Error("❌ SafeGoroutine(panic) sent nil, want error\n" +
				"   Hint: recover() in the goroutine and send the panic as an error")
		} else if !strings.Contains(err.Error(), "goroutine exploded") {
			t.Errorf("❌ SafeGoroutine(panic) error = %q, should contain \"goroutine exploded\"", err)
		} else {
			t.Logf("✅ SafeGoroutine(panic) caught: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("❌ SafeGoroutine(panic) timed out")
	}
}

func TestRetry_SucceedsFirstTry(t *testing.T) {
	var count int32
	err := Retry(3, 10*time.Millisecond, func() error {
		atomic.AddInt32(&count, 1)
		return nil
	})
	if err != nil {
		t.Errorf("❌ Retry(succeeds first) = %v, want nil", err)
	} else if atomic.LoadInt32(&count) != 1 {
		t.Errorf("❌ Called %d times, want 1 — stop retrying after success!\n"+
			"   Hint: return nil immediately when fn() returns nil", atomic.LoadInt32(&count))
	} else {
		t.Log("✅ Retry succeeded on first attempt")
	}
}

func TestRetry_SucceedsOnThirdTry(t *testing.T) {
	var count int32
	err := Retry(5, 10*time.Millisecond, func() error {
		n := atomic.AddInt32(&count, 1)
		if n < 3 {
			return fmt.Errorf("attempt %d failed", n)
		}
		return nil
	})
	if err != nil {
		t.Errorf("❌ Retry(succeeds 3rd) = %v, want nil", err)
	} else if atomic.LoadInt32(&count) != 3 {
		t.Errorf("❌ Called %d times, want 3", atomic.LoadInt32(&count))
	} else {
		t.Log("✅ Retry succeeded on third attempt")
	}
}

func TestRetry_AllAttemptsFail(t *testing.T) {
	var count int32
	err := Retry(3, 10*time.Millisecond, func() error {
		atomic.AddInt32(&count, 1)
		return fmt.Errorf("always fails")
	})
	if err == nil {
		t.Error("❌ Retry(all fail) = nil, want error\n" +
			"   Hint: return the last error after all attempts are exhausted")
	} else if atomic.LoadInt32(&count) != 3 {
		t.Errorf("❌ Called %d times, want 3 — should try maxAttempts times", atomic.LoadInt32(&count))
	} else {
		t.Logf("✅ Retry(all fail) correctly returned: %v", err)
	}
}

func TestRetryWithBackoff_ExponentialDelay(t *testing.T) {
	var timestamps []time.Time
	cfg := BackoffOptions{
		MaxAttempts: 4,
		BaseDelay:   50 * time.Millisecond,
		MaxDelay:    1 * time.Second,
	}

	_ = RetryWithBackoff(cfg, func() error {
		timestamps = append(timestamps, time.Now())
		return fmt.Errorf("fail")
	})

	if len(timestamps) != 4 {
		t.Fatalf("❌ Expected 4 attempts, got %d", len(timestamps))
	}

	// Verify delays are roughly exponential: 50ms, 100ms, 200ms
	for i := 1; i < len(timestamps); i++ {
		gap := timestamps[i].Sub(timestamps[i-1])
		expectedMin := cfg.BaseDelay * time.Duration(1<<(i-1)) * 7 / 10 // 70% of expected (allow jitter)
		if gap < expectedMin {
			t.Errorf("❌ Gap %d→%d = %v, expected at least ~%v (exponential backoff)\n"+
				"   Hint: delay = baseDelay * 2^attempt, i.e. 50ms, 100ms, 200ms...",
				i-1, i, gap, expectedMin)
		} else {
			t.Logf("✅ Gap %d→%d = %v (exponential ✓)", i-1, i, gap)
		}
	}
}

func TestRetryWithBackoff_CapsAtMaxDelay(t *testing.T) {
	var timestamps []time.Time
	cfg := BackoffOptions{
		MaxAttempts: 5,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    150 * time.Millisecond, // cap kicks in at attempt 2 (200ms > 150ms)
	}

	_ = RetryWithBackoff(cfg, func() error {
		timestamps = append(timestamps, time.Now())
		return fmt.Errorf("fail")
	})

	if len(timestamps) < 4 {
		t.Fatalf("❌ Expected at least 4 attempts, got %d", len(timestamps))
	}

	// Later gaps should not exceed maxDelay + tolerance
	for i := 2; i < len(timestamps); i++ {
		gap := timestamps[i].Sub(timestamps[i-1])
		if gap > cfg.MaxDelay+80*time.Millisecond { // generous tolerance for scheduling
			t.Errorf("❌ Gap %d→%d = %v exceeds maxDelay %v\n"+
				"   Hint: cap the calculated delay at cfg.MaxDelay",
				i-1, i, gap, cfg.MaxDelay)
		} else {
			t.Logf("✅ Gap %d→%d = %v (capped at maxDelay ✓)", i-1, i, gap)
		}
	}
}

func TestRetryClassified_RetriesTransient(t *testing.T) {
	var count int32
	err := RetryClassified(5, 10*time.Millisecond, func() error {
		n := atomic.AddInt32(&count, 1)
		if n < 3 {
			return fmt.Errorf("transient error %d", n) // retryable
		}
		return nil
	})
	if err != nil {
		t.Errorf("❌ RetryClassified(transient→success) = %v, want nil", err)
	} else if atomic.LoadInt32(&count) != 3 {
		t.Errorf("❌ Called %d times, want 3", atomic.LoadInt32(&count))
	} else {
		t.Log("✅ RetryClassified retried transient errors and succeeded")
	}
}

func TestRetryClassified_StopsOnPermanent(t *testing.T) {
	var count int32
	err := RetryClassified(5, 10*time.Millisecond, func() error {
		n := atomic.AddInt32(&count, 1)
		if n == 1 {
			return fmt.Errorf("transient") // retryable
		}
		return MarkPermanent(fmt.Errorf("bad request")) // permanent — stop!
	})

	if err == nil {
		t.Error("❌ RetryClassified(permanent) = nil, want error")
	} else if atomic.LoadInt32(&count) > 2 {
		t.Errorf("❌ Called %d times, want 2 — should stop on permanent error!\n"+
			"   Hint: check IsPermanentError(err) and return immediately if true",
			atomic.LoadInt32(&count))
	} else if !strings.Contains(err.Error(), "bad request") {
		t.Errorf("❌ Error = %q, should contain \"bad request\"", err)
	} else {
		t.Logf("✅ RetryClassified stopped on permanent error after %d attempts: %v",
			atomic.LoadInt32(&count), err)
	}
}

func TestRetryWithContext_RespectsCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	var count int32
	start := time.Now()
	err := RetryWithContext(ctx, 100, 100*time.Millisecond, func() error {
		atomic.AddInt32(&count, 1)
		return fmt.Errorf("always fails")
	})

	elapsed := time.Since(start)
	if err == nil {
		t.Error("❌ RetryWithContext(cancelled) = nil, want error")
	} else if elapsed > 500*time.Millisecond {
		t.Errorf("❌ Took %v — didn't respect context timeout!\n"+
			"   Hint: use select { case <-ctx.Done(): ... case <-time.After(delay): ... }",
			elapsed)
	} else if atomic.LoadInt32(&count) > 5 {
		t.Errorf("❌ %d attempts in 150ms is suspicious — are you actually sleeping?",
			atomic.LoadInt32(&count))
	} else {
		t.Logf("✅ RetryWithContext stopped after %v (%d attempts): %v",
			elapsed, atomic.LoadInt32(&count), err)
	}
}

func TestRetryWithContext_SucceedsBeforeTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var count int32
	err := RetryWithContext(ctx, 5, 10*time.Millisecond, func() error {
		n := atomic.AddInt32(&count, 1)
		if n < 3 {
			return fmt.Errorf("not yet")
		}
		return nil
	})

	if err != nil {
		t.Errorf("❌ RetryWithContext(succeeds) = %v, want nil", err)
	} else if atomic.LoadInt32(&count) != 3 {
		t.Errorf("❌ Called %d times, want 3", atomic.LoadInt32(&count))
	} else {
		t.Log("✅ RetryWithContext succeeded before timeout")
	}
}

func TestRetryWithContext_AlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	var count int32
	err := RetryWithContext(ctx, 5, 10*time.Millisecond, func() error {
		atomic.AddInt32(&count, 1)
		return fmt.Errorf("should not run")
	})

	if err == nil {
		t.Error("❌ RetryWithContext(pre-cancelled) = nil, want error")
	} else if atomic.LoadInt32(&count) > 0 {
		t.Errorf("❌ fn was called %d times on already-cancelled context\n"+
			"   Hint: check ctx.Err() BEFORE calling fn",
			atomic.LoadInt32(&count))
	} else {
		t.Logf("✅ RetryWithContext(pre-cancelled) returned immediately: %v", err)
	}
}
