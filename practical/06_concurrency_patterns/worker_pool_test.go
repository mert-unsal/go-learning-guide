package concurrency_patterns

import (
	"fmt"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	numWorkers := 3
	queueSize := 5
	pool := NewWorkerPool(numWorkers, queueSize)
	pool.Start()

	// Create a channel to collect results
	results := make(chan Result, 10)

	// Send 10 jobs
	for i := 1; i <= 10; i++ {
		job := Job{
			ID:       i,
			Payload:  fmt.Sprintf("Task-%d", i),
			ResultCh: results,
		}
		pool.AddJob(job)
	}

	// Close the results channel after receiving all expected results, or check count
	// Here we just read 10 results
	for i := 1; i <= 10; i++ {
		select {
		case res := <-results:
			if res.Err != nil {
				t.Errorf("Job %d failed: %v", res.JobID, res.Err)
			}
			// t.Logf("Got result: %s", res.Data)
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for results")
		}
	}

	// Test graceful shutdown (drain)
	pool.StopAndDrain()
}

func TestWorkerPool_Stop(t *testing.T) {
	pool := NewWorkerPool(2, 2)
	pool.Start()
	
	// Add a job
	pool.AddJob(Job{ID: 1, Payload: "Quick"})
	
	// Stop immediately
	pool.Stop() 
	// If this hangs, the test fails (timeout)
}
