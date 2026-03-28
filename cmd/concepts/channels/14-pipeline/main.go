package main

import "fmt"

// ============================================================
// Pipeline — Chain of Stages Connected by Channels
// ============================================================
//
// The Problem:
//   You have a multi-step data transformation: parse → validate → enrich
//   → write. Each step has different latency characteristics. Running
//   them sequentially wastes time — step 2 sits idle while step 1
//   processes the next item. You want all stages running concurrently,
//   each processing a different item simultaneously.
//
// Real-world example:
//   An ETL pipeline: read rows from CSV (I/O-bound) → parse and validate
//   (CPU-bound) → enrich with API calls (network-bound) → batch-insert
//   to database (I/O-bound). Each stage runs as a goroutine; channels
//   connect them. While stage 1 reads row N+1, stage 2 validates row N.
//
// The Pattern:
//   Each stage is a function with this signature:
//     func stage(in <-chan T) <-chan U
//   It launches a goroutine that:
//     1. Reads from in with for-range (blocks until data or close)
//     2. Transforms the value
//     3. Sends the result to out
//     4. When in closes (range exits), closes out
//
//   Stages compose by nesting:
//     result := stage3(stage2(stage1(source)))
//
//   ┌──────────┐    ch1    ┌──────────┐    ch2    ┌──────────┐    ch3    ┌──────────┐
//   │ generate │─────────►│  double  │─────────►│  addTen  │─────────►│ collect  │
//   │  1..5    │          │  x * 2   │          │  x + 10  │          │ results  │
//   └──────────┘          └──────────┘          └──────────┘          └──────────┘
//        │                      │                     │                     │
//     closes ch1 ──► range    closes ch2 ──► range   closes ch3 ──► range exits
//     when done      exits    when done      exits   when done
//
// Key rules:
//   1. The stage that CREATES a channel is responsible for CLOSING it.
//      generate creates ch1 and closes it. double creates ch2 and closes
//      it. This prevents "who closes?" confusion in concurrent code.
//
//   2. Cancellation propagates naturally via backpressure: if the consumer
//      stops reading, channels fill up, sends block, and all upstream
//      goroutines park. They consume zero CPU — the runtime removes them
//      from the scheduler's run queue. For eager cancellation (don't wait
//      for buffers to fill), pass a context.Context and select on ctx.Done()
//      alongside every send and receive.
//
//   3. Each stage runs in its own goroutine. With unbuffered channels,
//      at most one item is "in flight" between stages. Add buffering
//      (make(chan T, n)) if stages have variable processing times and
//      you want to absorb bursts — but measure first, don't buffer
//      speculatively.
//
// Why channels work here:
//   Channels provide the three things a pipeline needs:
//     - Data transfer: values flow from stage to stage
//     - Synchronization: send blocks until receiver is ready (backpressure)
//     - Completion signal: close(ch) tells downstream "no more data"
//   No mutexes, no condition variables, no manual signaling needed.
//   The for-range loop over a channel is the idiomatic "read until done"
//   construct — it exits when the channel is closed AND drained.

// generate produces integers from start to end (inclusive) on a channel.
// It creates and owns the output channel, closing it when all values
// are sent. This is the source stage of the pipeline.
func generate(start, end int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i <= end; i++ {
			out <- i
		}
	}()
	return out
}

// double reads integers from in, multiplies each by 2, and sends the
// result to a new output channel. Closes output when input is exhausted.
func double(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * 2
		}
	}()
	return out
}

// addTen reads integers from in, adds 10 to each, and sends the result
// to a new output channel. Closes output when input is exhausted.
func addTen(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v + 10
		}
	}()
	return out
}

func main() {
	// Compose the pipeline: each stage's output feeds the next stage's input.
	// This reads right-to-left: generate → double → addTen
	pipeline := addTen(double(generate(1, 5)))

	fmt.Println("  Pipeline: generate(1..5) → double(x*2) → addTen(x+10)")
	fmt.Println()

	// Collect and display results
	var results []int
	for v := range pipeline {
		results = append(results, v)
	}

	for _, original := range []int{1, 2, 3, 4, 5} {
		doubled := original * 2
		final := doubled + 10
		fmt.Printf("    %d → ×2 = %d → +10 = %d\n", original, doubled, final)
	}

	fmt.Printf("\n  Pipeline output: %v\n", results)
	fmt.Printf("  All %d items flowed through 3 concurrent stages\n", len(results))
}
