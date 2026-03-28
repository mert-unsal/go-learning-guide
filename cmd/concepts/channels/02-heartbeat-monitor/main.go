// Package main contains a standalone conceptual example for the Heartbeat Monitor pattern.
package main

import (
	"fmt"
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
// Heartbeat / Health Monitor — Timer Reset Pattern
// ============================================================
//
// The Problem:
//   You run a critical background service (database replica, worker process,
//   microservice dependency). You need to detect when it becomes unresponsive.
//   The service sends periodic heartbeat signals. If no heartbeat arrives
//   within a timeout window, the service is considered dead and ops should
//   be alerted.
//
// Real-world example:
//   A Kubernetes sidecar monitoring a main container's health endpoint.
//   If the container stops responding for 5 seconds, the sidecar reports
//   it as unhealthy, triggering a pod restart.
//
// The Pattern:
//   for/select with timer.Reset — the inactivity detection pattern:
//     - heartbeat channel: each heartbeat resets the timer
//     - timer.C: fires when no heartbeat arrives within the timeout window
//     - done channel: clean shutdown
//
// Why channels work here:
//   This combines two patterns you already know:
//   1. Timer reset (from the DualTimeoutWorker exercise) — reset on activity
//   2. for/select event loop — multiplex heartbeats, timeouts, and shutdown
//
//   The timer.Stop()/drain/Reset dance before resetting is essential:
//   if a heartbeat arrives after the timer fires but before select runs,
//   the stale timer value must be drained to avoid a false timeout on
//   the next iteration.
//
//   Timeline:
//
//   ──heartbeat────heartbeat────────────────heartbeat──────────
//        │              │            │           │
//        ▼              ▼            ▼           ▼
//      reset          reset       TIMEOUT!     reset
//      timer          timer       alert ops    timer

// HeartbeatResult represents what the monitor observed.
type HeartbeatResult struct {
	Heartbeats int
	Timeouts   int
}

// HeartbeatMonitor watches a heartbeat channel and detects when the
// service stops sending heartbeats within the timeout window.
// Returns after the done channel is closed.
func HeartbeatMonitor(heartbeat <-chan struct{}, done <-chan struct{}, timeout time.Duration) HeartbeatResult {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	result := HeartbeatResult{}
	fmt.Printf("  %s[monitor]%s started — timeout window: %s%v%s\n", cyan+bold, reset, magenta, timeout, reset)

	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				fmt.Printf("  %s[monitor]%s heartbeat channel closed — exiting\n", cyan+bold, reset)
				return result
			}
			// Got heartbeat — service is alive, reset the inactivity timer
			result.Heartbeats++
			fmt.Printf("  %s[monitor]%s %s✔ heartbeat #%d received%s — timer reset to %s%v%s\n",
				cyan+bold, reset, green, result.Heartbeats, reset, magenta, timeout, reset)
			if !timer.Stop() {
				select {
				case <-timer.C: // drain if already fired
				default:
				}
			}
			timer.Reset(timeout)

		case <-timer.C:
			// No heartbeat within timeout — service is unresponsive
			result.Timeouts++
			fmt.Printf("  %s[monitor]%s %s⚠ TIMEOUT #%d%s — no heartbeat for %s%v%s, service may be dead!\n",
				cyan+bold, reset, red+bold, result.Timeouts, reset, magenta, timeout, reset)
			timer.Reset(timeout) // keep monitoring

		case <-done:
			fmt.Printf("  %s[monitor]%s %sdone signal received%s — shutting down\n", cyan+bold, reset, cyan, reset)
			return result
		}
	}
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Heartbeat Monitor — Timer Reset Pattern        %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Pattern Overview%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Monitor watches a heartbeat channel with a rolling timeout%s\n", green, reset)
	fmt.Printf("  %s✔ Each heartbeat resets the timer — service is alive%s\n", green, reset)
	fmt.Printf("  %s✔ Timer fires if no heartbeat within window — service is dead%s\n", green, reset)
	fmt.Printf("  %s⚠ timer.Stop()/drain/Reset dance avoids stale timer events%s\n\n", yellow, reset)

	heartbeat := make(chan struct{})
	done := make(chan struct{})

	// Start monitor with 200ms timeout
	fmt.Printf("%s▸ Starting Monitor%s\n", cyan+bold, reset)
	fmt.Printf("  %sChannels: heartbeat (unbuffered), done (unbuffered), resultCh (buffered 1)%s\n", dim, reset)
	resultCh := make(chan HeartbeatResult, 1)
	go func() {
		resultCh <- HeartbeatMonitor(heartbeat, done, 200*time.Millisecond)
	}()

	// Send 3 heartbeats at 50ms intervals (well within 200ms timeout)
	fmt.Printf("\n%s▸ Sending Heartbeats (50ms interval, well within 200ms timeout)%s\n", cyan+bold, reset)
	for i := 0; i < 3; i++ {
		heartbeat <- struct{}{}
		fmt.Printf("  %s💓 heartbeat %d sent%s — blocks until monitor receives (unbuffered chan)\n", green+bold, i+1, reset)
		time.Sleep(50 * time.Millisecond)
	}

	// Stop sending — let the timeout fire
	fmt.Printf("\n%s▸ Simulating Service Failure%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ Stopped sending heartbeats — waiting 300ms to trigger timeout%s\n", yellow, reset)
	time.Sleep(300 * time.Millisecond)

	// Shut down the monitor
	fmt.Printf("\n%s▸ Graceful Shutdown%s\n", cyan+bold, reset)
	fmt.Printf("  %sClosing done channel to signal monitor goroutine to exit%s\n", dim, reset)
	close(done)
	result := <-resultCh

	fmt.Printf("\n%s▸ Results%s\n", cyan+bold, reset)
	fmt.Printf("  Heartbeats received: %s%d%s\n", magenta, result.Heartbeats, reset)
	fmt.Printf("  Timeouts detected:   %s%d%s\n", magenta, result.Timeouts, reset)
	fmt.Printf("\n  %s✔ The monitor correctly detected liveness and failure states%s\n", green, reset)
	fmt.Printf("  %s⚠ In production: timeouts trigger alerts, restarts, or failover%s\n", yellow, reset)
}
