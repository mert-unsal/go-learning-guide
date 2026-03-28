package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================
// Semaphore — Buffered Channel as Bounded Concurrency Gate
// ============================================================
//
// The Problem:
//   You have 1000 tasks to run concurrently, but the downstream resource
//   (database, API, file system) can only handle 10 concurrent operations.
//   Launching 1000 goroutines overwhelms the resource. You need a gate
//   that limits how many goroutines are active at once.
//
// Real-world example:
//   A batch job that resizes 10,000 images. Each resize needs a file
//   handle. The OS limits open file descriptors to 1024. You use a
//   semaphore of 100 to keep file handle usage safely under the limit.
//
// The Pattern:
//   sem := make(chan struct{}, maxConcurrent)
//   // For each task:
//   sem <- struct{}{}    // acquire: blocks when buffer full (at capacity)
//   go func() {
//       defer func() { <-sem }()  // release: frees a slot
//       doWork()
//   }()
//
// Why channels work here:
//   A buffered channel IS a counting semaphore:
//   - Buffer capacity = maximum concurrent operations
//   - Send = acquire (blocks when full → max concurrency reached)
//   - Receive = release (frees one slot for the next goroutine)
//   - The runtime handles all synchronization — no mutexes needed
//
// How it works at the runtime level:
//   make(chan struct{}, 3) creates a channel with 3-slot buffer.
//
//   sem <- struct{}{} :
//     If buffer has space → value enqueued, goroutine continues (acquired)
//     If buffer full → goroutine parks in sendq (waits for a slot)
//
//   <-sem :
//     Dequeues one value from buffer → frees a slot
//     If a goroutine is parked in sendq → it's woken up (slot acquired)
//
//   ┌────────────────────────────────────────────────┐
//   │  Semaphore: make(chan struct{}, 3)             │
//   │                                                │
//   │  acquire ──► [■][■][■]  buffer full → block    │
//   │               ▲                                │
//   │  release ──► [■][■][_]  slot freed → unblock   │
//   │                                                │
//   │  Task 1: running   ■                           │
//   │  Task 2: running   ■                           │
//   │  Task 3: running   ■                           │
//   │  Task 4: waiting   ⌛ (blocked on acquire)     │
//   │  Task 5: waiting   ⌛ (blocked on acquire)     │
//   │                                                │
//   │  Task 1 finishes → release → Task 4 unblocks   │
//   └────────────────────────────────────────────────┘
//
// Compare with Exercise 7 (ProcessWithLimit):
//   The exercises package uses this exact pattern — a buffered channel
//   as a semaphore to limit concurrent processing. The key difference
//   is where you acquire: before launching the goroutine (controls
//   goroutine count) vs inside the goroutine (controls active work).
//   Acquiring BEFORE the goroutine launch is generally preferred —
//   it prevents goroutine pile-up when tasks arrive faster than they
//   complete.

func main() {
	const numTasks = 10
	const maxConcurrent = 3

	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	var currentRunning atomic.Int32
	var peakConcurrency atomic.Int32
	var mu sync.Mutex
	var log []string

	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		sem <- struct{}{} // acquire — blocks when 3 tasks already running
		go func(taskID int) {
			defer wg.Done()
			defer func() { <-sem }() // release — frees slot on exit

			// Track concurrency
			running := currentRunning.Add(1)
			// Update peak using compare-and-swap loop
			for {
				peak := peakConcurrency.Load()
				if running <= peak || peakConcurrency.CompareAndSwap(peak, running) {
					break
				}
			}

			// Simulate work
			time.Sleep(20 * time.Millisecond)

			mu.Lock()
			log = append(log, fmt.Sprintf("task %2d completed (concurrent: %d)", taskID, running))
			mu.Unlock()

			currentRunning.Add(-1)
		}(i)
	}

	wg.Wait()

	fmt.Printf("  %d tasks, max %d concurrent:\n", numTasks, maxConcurrent)
	for _, entry := range log {
		fmt.Printf("    %s\n", entry)
	}
	fmt.Printf("  Peak concurrency observed: %d (limit: %d)\n", peakConcurrency.Load(), maxConcurrent)
}
