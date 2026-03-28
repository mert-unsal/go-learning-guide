// Package concepts contains standalone conceptual examples for channel patterns.
// Each file explains one production pattern: the problem, the solution, and why
// channels are the right tool.
package concepts

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================
// Worker Event Loop — for/select as a Multiplexed Event Handler
// ============================================================
//
// The Problem:
//   Background workers (Pub/Sub consumers, queue processors, batch jobs)
//   have no framework wrapping them. Unlike HTTP handlers where middleware
//   captures metrics automatically, a worker manages its own lifecycle.
//   It must handle multiple concerns in a single goroutine:
//     - Process incoming work from a job queue
//     - Run periodic maintenance (flush batches, report metrics)
//     - Shut down gracefully when told to stop
//
// Real-world example:
//   A Cloud Run worker pulling from Pub/Sub, batching database writes
//   every 10 seconds, and reporting queue lag metrics — all in one
//   for/select loop.
//
// The Pattern:
//   for/select with three channel cases:
//     - jobs channel: primary work arrives here
//     - ticker.C: periodic maintenance fires here
//     - done channel: shutdown signal arrives here
//
// Why channels work here:
//   The goroutine SLEEPS between events — zero CPU when idle.
//   select wakes it instantly when any channel has data.
//   Each channel represents a different event source, and select
//   multiplexes them into a single handler loop.
//
//   ┌─────────── for-select loop ───────────┐
//   │                                        │
//   │  JOBS ──────► process(job)             │  ← primary work
//   │                                        │
//   │  TICKER ────► flush + report           │  ← periodic maintenance
//   │                                        │
//   │  DONE ──────► drain remaining → exit   │  ← graceful shutdown
//   │                                        │
//   │  The goroutine SLEEPS between events.  │
//   │  Zero CPU when idle. Wakes instantly   │
//   │  when any channel has data.            │
//   └────────────────────────────────────────┘
//
// The drain loop on shutdown:
//   When the done signal arrives, remaining jobs in the channel buffer
//   should still be processed. The inner select+default loop drains
//   them without blocking: if the channel is empty, default fires
//   and the worker exits cleanly.

// Job represents a unit of work for the worker to process.
type Job struct {
	ID   int
	Data string
}

// WorkerEventLoop demonstrates the for/select event loop pattern.
// It processes jobs, runs periodic maintenance, and shuts down gracefully.
func WorkerEventLoop(jobs <-chan Job, done <-chan struct{}) []string {
	var processed []string
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	flushCount := 0

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				// Channel closed — no more jobs coming
				return processed
			}
			// Primary work: process the job
			result := fmt.Sprintf("processed:%s", job.Data)
			processed = append(processed, result)

		case <-ticker.C:
			// Periodic maintenance: flush batches, report metrics
			flushCount++
			_ = flushCount // in production: flush DB batch, report queue depth

		case <-done:
			// Graceful shutdown — drain remaining jobs in buffer
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return processed
					}
					result := fmt.Sprintf("drained:%s", job.Data)
					processed = append(processed, result)
				default:
					// buffer empty — exit cleanly
					return processed
				}
			}
		}
	}
}

// DemonstrateWorkerEventLoop shows the worker processing jobs, then
// shutting down gracefully and draining remaining work.
func DemonstrateWorkerEventLoop() {
	jobs := make(chan Job, 10)
	done := make(chan struct{})

	// Pre-load some jobs
	jobs <- Job{ID: 1, Data: "order-101"}
	jobs <- Job{ID: 2, Data: "order-102"}
	jobs <- Job{ID: 3, Data: "order-103"}

	var results []string
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		results = WorkerEventLoop(jobs, done)
	}()

	// Let worker process for a bit
	time.Sleep(50 * time.Millisecond)

	// Add more jobs then signal shutdown
	jobs <- Job{ID: 4, Data: "order-104"}
	jobs <- Job{ID: 5, Data: "order-105"}
	close(done) // signal graceful shutdown

	wg.Wait()

	for _, r := range results {
		fmt.Printf("  %s\n", r)
	}
	fmt.Printf("  Total: %d jobs handled\n", len(results))
}
