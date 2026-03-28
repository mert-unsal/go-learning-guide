package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================
// Worker Pool with Backpressure — N Goroutines, Shared Channel
// ============================================================
//
// The Problem:
//   You have a stream of jobs (HTTP requests, queue messages, file
//   processing tasks). Processing each job takes time. If you spawn
//   one goroutine per job, you'll exhaust memory at 100k jobs. If you
//   process sequentially, throughput is limited by single-core speed.
//   You need bounded parallelism with natural flow control.
//
// Real-world example:
//   An image processing pipeline: thumbnails arrive from an S3 event
//   queue at 1000/sec. You run 8 workers (one per CPU core). When all
//   8 are busy resizing images, the channel buffer fills up, and the
//   SQS consumer blocks — backpressure propagates to the source.
//
// The Pattern:
//   N goroutines all range over the SAME jobs channel:
//     jobs := make(chan Job, bufferSize)
//     for i := 0; i < N; i++ {
//         go worker(jobs)   // each worker: for job := range jobs { ... }
//     }
//     // producer sends to jobs, then close(jobs) when done
//
// Why channels work here:
//   1. Fan-out is free: multiple goroutines receiving from the same channel
//      is safe — the runtime ensures each value goes to exactly ONE receiver
//   2. Backpressure is automatic: when all workers are busy, the channel
//      buffer fills up, and send blocks — the producer naturally slows down
//   3. Clean shutdown: close(jobs) causes all workers' range loops to exit
//   4. No load balancer needed: the runtime's recvq (FIFO wait queue)
//      distributes work fairly to whichever worker is free next
//
// Under the hood — how work is distributed:
//   When a producer sends to the jobs channel:
//   - If a worker is waiting in recvq → direct handoff (no buffer copy)
//   - If no worker is waiting but buffer has space → enqueue in buffer
//   - If buffer is full → producer parks in sendq (BACKPRESSURE)
//
//   When a worker calls <-jobs:
//   - If buffer has data → dequeue from buffer
//   - If buffer empty but sender waiting → direct receive from sender
//   - If both empty → worker parks in recvq (waits for work)
//
//   ┌────────────────────────────────────────────────┐
//   │  Producer ──send──► [buffer size=4] ──recv──►  │
//   │                     [job][job][..][..]          │
//   │                          │    │    │            │
//   │                          ▼    ▼    ▼            │
//   │                        W0   W1   W2            │
//   │                                                │
//   │  Buffer full? Producer blocks. (backpressure)  │
//   │  All workers busy? Buffer absorbs burst.       │
//   │  Buffer + workers full? System at capacity.    │
//   └────────────────────────────────────────────────┘
//
// Backpressure flow:
//   producer too fast → buffer fills → send blocks → producer slows
//   workers finish → buffer drains → send unblocks → producer resumes
//   This is NATURAL flow control — no rate limiter or circuit breaker needed.

func main() {
	const numWorkers = 3
	const numJobs = 10
	const bufferSize = 2 // small buffer to make backpressure visible

	jobs := make(chan int, bufferSize)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var log []string

	// Start workers — each ranges over the shared jobs channel
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate work
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
				log = append(log, fmt.Sprintf("worker %d processed job %d", id, job))
				mu.Unlock()
			}
			// range exits when jobs is closed and buffer is drained
		}(w)
	}

	// Producer — sends jobs, blocks when buffer is full (backpressure)
	for j := 1; j <= numJobs; j++ {
		jobs <- j // blocks when buffer full AND all workers busy
	}
	close(jobs) // signal workers: no more jobs coming
	wg.Wait()

	fmt.Printf("  %d workers processed %d jobs (buffer=%d):\n", numWorkers, numJobs, bufferSize)
	for _, entry := range log {
		fmt.Printf("    %s\n", entry)
	}
}
