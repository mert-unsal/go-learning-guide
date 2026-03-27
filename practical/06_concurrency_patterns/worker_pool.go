package concurrency_patterns

import (
	"fmt"
	"sync"
	"time"
)

// Job represents the work to be done.
type Job struct {
	ID       int
	Payload  string
	ResultCh chan<- Result // Channel to send result back to caller
}

// Result represents the outcome of a job.
type Result struct {
	JobID int
	Data  string
	Err   error
}

// WorkerPool manages a pool of workers to process jobs.
type WorkerPool struct {
	numWorkers int
	jobQueue   chan Job
	wg         sync.WaitGroup
	quit       chan struct{}
}

// NewWorkerPool creates a new worker pool with the given number of workers and queue size.
func NewWorkerPool(numWorkers, queueSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobQueue:   make(chan Job, queueSize),
		quit:       make(chan struct{}),
	}
}

// Start launches the workers.
func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker processes jobs from the queue until stopped.
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	// fmt.Printf("Worker %d started\n", id)

	for {
		select {
		case job, ok := <-wp.jobQueue:
			if !ok {
				// Queue closed and drained
				return
			}
			// Process the job
			result := process(job)
			
			// Send result if a channel was provided
			if job.ResultCh != nil {
				job.ResultCh <- result
			}

		case <-wp.quit:
			// Received shutdown signal
			// fmt.Printf("Worker %d stopping\n", id)
			return
		}
	}
}

// process simulates job processing.
func process(j Job) Result {
	// Simulate work
	time.Sleep(10 * time.Millisecond) // Fast for tests
	return Result{
		JobID: j.ID,
		Data:  fmt.Sprintf("Processed %s", j.Payload),
		Err:   nil,
	}
}

// AddJob adds a job to the queue. Returns error if pool is stopped or full (if using try-send).
// This simple version blocks if queue is full.
func (wp *WorkerPool) AddJob(j Job) {
	wp.jobQueue <- j
}

// Stop initiates a "hard" shutdown.
// Workers stop after finishing their current task, ignoring pending jobs.
func (wp *WorkerPool) Stop() {
	close(wp.quit)
	wp.wg.Wait()
}

// StopAndDrain initiates a "graceful" shutdown.
// Workers process all jobs in the queue before exiting.
func (wp *WorkerPool) StopAndDrain() {
	close(wp.jobQueue)
	wp.wg.Wait()
}
