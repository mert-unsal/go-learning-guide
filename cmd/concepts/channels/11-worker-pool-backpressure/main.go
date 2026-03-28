package main

import (
	"fmt"
	"sync"
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

	workerColors := []string{cyan, yellow, magenta}

	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Worker Pool with Backpressure                  %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Configuration%s\n", cyan+bold, reset)
	fmt.Printf("  Workers: %s%d%s   Jobs: %s%d%s   Buffer: %s%d%s\n\n",
		magenta, numWorkers, reset, magenta, numJobs, reset, magenta, bufferSize, reset)

	jobs := make(chan int, bufferSize)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var log []string

	// Start workers — each ranges over the shared jobs channel
	fmt.Printf("%s▸ Starting Workers%s\n", cyan+bold, reset)
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		fmt.Printf("  %s%s[W%d]%s ready — listening on shared jobs channel\n",
			bold, workerColors[w], w, reset)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				mu.Lock()
				fmt.Printf("  %s%s[W%d]%s ← job %s%d%s\n",
					bold, workerColors[id], id, reset, magenta, job, reset)
				mu.Unlock()
				// Simulate work
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
				log = append(log, fmt.Sprintf("worker %d processed job %d", id, job))
				fmt.Printf("  %s%s[W%d]%s %s✔%s job %s%d%s complete\n",
					bold, workerColors[id], id, reset, green, reset, magenta, job, reset)
				mu.Unlock()
			}
			// range exits when jobs is closed and buffer is drained
		}(w)
	}

	// Producer — sends jobs, blocks when buffer is full (backpressure)
	fmt.Printf("\n%s▸ Producer Sending Jobs%s\n", cyan+bold, reset)
	for j := 1; j <= numJobs; j++ {
		mu.Lock()
		if len(jobs) >= bufferSize {
			fmt.Printf("  %s⚠ Buffer full (%d/%d) — producer blocks until a worker reads (backpressure!)%s\n",
				yellow, len(jobs), bufferSize, reset)
		}
		mu.Unlock()
		jobs <- j // blocks when buffer full AND all workers busy
		mu.Lock()
		fmt.Printf("  %s→%s sent job %s%d%s  %s[buffer: %d/%d]%s\n",
			green, reset, magenta, j, reset, dim, len(jobs), bufferSize, reset)
		mu.Unlock()
	}
	close(jobs) // signal workers: no more jobs coming
	mu.Lock()
	fmt.Printf("  %s✔ All jobs sent, channel closed — workers drain remaining%s\n", green, reset)
	mu.Unlock()
	wg.Wait()

	fmt.Printf("\n%s▸ Processing Log%s\n", cyan+bold, reset)
	fmt.Printf("  %s%d%s workers processed %s%d%s jobs (buffer=%s%d%s):\n",
		magenta, numWorkers, reset, magenta, numJobs, reset, magenta, bufferSize, reset)
	for _, entry := range log {
		fmt.Printf("    %s%s%s\n", dim, entry, reset)
	}

	fmt.Printf("\n%s▸ Key Observations%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Bounded parallelism: only %d goroutines, regardless of job count%s\n", green, numWorkers, reset)
	fmt.Printf("  %s✔ Backpressure is automatic: buffer=%d fills fast, producer blocks%s\n", green, bufferSize, reset)
	fmt.Printf("  %s✔ Fair distribution: runtime recvq hands work to the next free worker%s\n", green, reset)
	fmt.Printf("  %s✔ Clean shutdown: close(jobs) → range exits → wg.Done() → Wait() returns%s\n", green, reset)
	fmt.Printf("  %s⚠ Work is NOT evenly split — fast workers steal more from the queue%s\n", yellow, reset)
}
