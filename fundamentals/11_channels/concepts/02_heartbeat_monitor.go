package concepts

import (
	"fmt"
	"time"
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

	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return result
			}
			// Got heartbeat — service is alive, reset the inactivity timer
			result.Heartbeats++
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
			timer.Reset(timeout) // keep monitoring

		case <-done:
			return result
		}
	}
}

// DemonstrateHeartbeatMonitor shows the monitor detecting healthy heartbeats
// and then detecting a timeout when heartbeats stop.
func DemonstrateHeartbeatMonitor() {
	heartbeat := make(chan struct{})
	done := make(chan struct{})

	// Start monitor with 200ms timeout
	resultCh := make(chan HeartbeatResult, 1)
	go func() {
		resultCh <- HeartbeatMonitor(heartbeat, done, 200*time.Millisecond)
	}()

	// Send 3 heartbeats at 50ms intervals (well within 200ms timeout)
	for i := 0; i < 3; i++ {
		heartbeat <- struct{}{}
		fmt.Printf("  💓 heartbeat %d sent\n", i+1)
		time.Sleep(50 * time.Millisecond)
	}

	// Stop sending — let the timeout fire
	fmt.Println("  ... stopped sending heartbeats ...")
	time.Sleep(300 * time.Millisecond)

	// Shut down the monitor
	close(done)
	result := <-resultCh

	fmt.Printf("  Heartbeats received: %d\n", result.Heartbeats)
	fmt.Printf("  Timeouts detected:   %d\n", result.Timeouts)
}
