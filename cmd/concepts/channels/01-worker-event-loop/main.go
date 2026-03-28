// Package main demonstrates the Worker Pool + Event Loop pattern with real-time
// status monitoring. Watch 3 workers drain a job queue, each job taking 1-2 seconds.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// ANSI color codes — works in all modern terminals including Windows Terminal.
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
	white   = "\033[37m"
	bgBlue  = "\033[44m"
	bgGreen = "\033[42m"
	bgRed   = "\033[41m"
)

// workerColors assigns a unique color to each worker for visual tracking.
var workerColors = []string{cyan, yellow, magenta, green, blue}

// ============================================================
// Worker Pool with Event Loop — Live Queue Simulation
// ============================================================
//
// What this demonstrates:
//   3 workers share a single jobs channel. Each worker runs a for/select
//   event loop. Jobs take 1-2 seconds to process (simulated work).
//   A status ticker prints a live dashboard every 500ms so you can
//   watch the queue drain in real time.
//
// Architecture:
//
//   ┌─────────────────────────────────────────────────────────┐
//   │  main goroutine                                        │
//   │  ┌─────────────┐                                       │
//   │  │ Load 12 jobs │──► jobs channel (buffered, cap=20)   │
//   │  └─────────────┘         │                             │
//   │                          ▼                             │
//   │         ┌────────────────┼────────────────┐            │
//   │         ▼                ▼                ▼            │
//   │   ┌──────────┐    ┌──────────┐    ┌──────────┐        │
//   │   │ Worker 1 │    │ Worker 2 │    │ Worker 3 │        │
//   │   │ for/sel  │    │ for/sel  │    │ for/sel  │        │
//   │   └──────────┘    └──────────┘    └──────────┘        │
//   │         │                │                │            │
//   │         └────────────────┼────────────────┘            │
//   │                          ▼                             │
//   │                   completed count                      │
//   │                                                        │
//   │   status ticker (500ms) reads:                         │
//   │     • active workers (atomic counter)                  │
//   │     • queue depth (len(jobs))                          │
//   │     • completed count (atomic counter)                 │
//   │                                                        │
//   │   done channel ──► all workers drain + exit            │
//   └─────────────────────────────────────────────────────────┘
//
// Key runtime observations:
//   • 3 workers compete on one channel — the scheduler decides who
//     wakes up. This is safe: channel receive is goroutine-safe
//     (hchan has an internal mutex).
//   • Active count goes up to 3 max (3 workers, each holds 1 job).
//   • Queue depth drops in bursts of ~3 as workers grab jobs.
//   • On shutdown, workers drain remaining buffered jobs before exiting.

// Job represents a unit of work.
type Job struct {
	ID       int
	Data     string
	Duration time.Duration // simulated processing time
}

// worker runs a for/select event loop: pulls jobs, handles ticks, responds to shutdown.
func worker(id int, jobs <-chan Job, done <-chan struct{}, active, completed *atomic.Int32, wg *sync.WaitGroup) {
	defer wg.Done()
	c := workerColors[(id-1)%len(workerColors)]
	tag := fmt.Sprintf("%s[worker-%d]%s", c, id, reset)

	fmt.Printf("%s ▶ Online — waiting for jobs\n", tag)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	processJob := func(job Job, phase string) {
		active.Add(1)
		fmt.Printf("%s ⚙ %s Job #%02d %s(%s)%s — will take %s%v%s\n",
			tag, phase, job.ID, dim, job.Data, reset, bold, job.Duration, reset)

		time.Sleep(job.Duration)

		active.Add(-1)
		completed.Add(1)
		fmt.Printf("%s %s✔ Done%s  Job #%02d (%s) — completed: %s%d%s total\n",
			tag, green, reset, job.ID, job.Data, bold, completed.Load(), reset)
	}

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("%s %s✖ Channel closed — exiting%s\n", tag, red, reset)
				return
			}
			processJob(job, "Processing")

		case <-ticker.C:
			fmt.Printf("%s %s⏱ Heartbeat — I'm alive and waiting for work%s\n", tag, dim, reset)

		case <-done:
			fmt.Printf("%s %s🛑 Shutdown signal — draining remaining jobs...%s\n", tag, red+bold, reset)
			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						fmt.Printf("%s    %sChannel closed during drain — exiting%s\n", tag, red, reset)
						return
					}
					processJob(job, fmt.Sprintf("%sDraining%s", yellow, reset))
				default:
					fmt.Printf("%s    %sBuffer empty — exiting cleanly%s\n", tag, green, reset)
					return
				}
			}
		}
	}
}

func main() {
	fmt.Printf("%s%s═══════════════════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Worker Pool + Event Loop — Live Queue Drain Simulation      %s\n", bold, blue, reset)
	fmt.Printf("%s%s═══════════════════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	const (
		numWorkers = 3
		numJobs    = 12
	)

	jobs := make(chan Job, 20)
	done := make(chan struct{})
	var active, completed atomic.Int32

	// Load all jobs upfront — each takes 1-2 seconds
	fmt.Printf("%s[main]%s Loading %s%d jobs%s into channel (cap=%d)...\n", white+bold, reset, green+bold, numJobs, reset, cap(jobs))
	for i := 1; i <= numJobs; i++ {
		duration := time.Duration(1000+rand.Intn(1000)) * time.Millisecond
		jobs <- Job{
			ID:       i,
			Data:     fmt.Sprintf("order-%03d", 100+i),
			Duration: duration,
		}
	}
	fmt.Printf("%s[main]%s Queue loaded: %s%d/%d%s slots used\n\n", white+bold, reset, yellow, len(jobs), cap(jobs), reset)

	// Status ticker — live dashboard every 500ms
	statusDone := make(chan struct{})
	go func() {
		tick := time.NewTicker(500 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				a := active.Load()
				c := completed.Load()
				q := len(jobs)

				// Build color-coded progress bar
				bar := ""
				for i := 0; i < int(a); i++ {
					bar += fmt.Sprintf("%s█%s", green, reset)
				}
				for i := int(a); i < numWorkers; i++ {
					bar += fmt.Sprintf("%s░%s", dim, reset)
				}

				// Color queue count: green=low, yellow=medium, red=high
				qColor := green
				if q > 6 {
					qColor = red
				} else if q > 3 {
					qColor = yellow
				}

				fmt.Printf("%s[status]%s Active: %s%d/%d%s [%s]  Queue: %s%-2d%s  Completed: %s%d/%d%s\n",
					dim, reset, bold, a, numWorkers, reset, bar, qColor, q, reset, green+bold, c, numJobs, reset)
			case <-statusDone:
				return
			}
		}
	}()

	// Launch worker pool
	var wg sync.WaitGroup
	fmt.Printf("%s[main]%s Launching %s%d workers%s...\n\n", white+bold, reset, cyan+bold, numWorkers, reset)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, done, &active, &completed, &wg)
	}

	// Wait until we have a few jobs left to show the drain behavior
	for completed.Load() < int32(numJobs-3) {
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println()
	fmt.Printf("%s%s[main] ═══════════════════════════════════════════════════════%s\n", bold, red, reset)
	fmt.Printf("%s[main]%s %d/%d jobs completed — %ssending shutdown signal%s\n", white+bold, reset, completed.Load(), numJobs, red+bold, reset)
	fmt.Printf("%s[main]%s Remaining buffered jobs will be drained by workers\n", white+bold, reset)
	fmt.Printf("%s%s[main] ═══════════════════════════════════════════════════════%s\n", bold, red, reset)
	fmt.Println()
	close(done)

	wg.Wait()
	close(statusDone)

	fmt.Println()
	fmt.Printf("%s%s═══════════════════════════════════════════════════════════════%s\n", bold, green, reset)
	fmt.Printf("%s%s  ✔ All done: %d/%d jobs processed by %d workers              %s\n", bold, green, completed.Load(), numJobs, numWorkers, reset)
	fmt.Printf("%s%s═══════════════════════════════════════════════════════════════%s\n", bold, green, reset)
}
