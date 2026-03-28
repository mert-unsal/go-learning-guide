package main

import (
	"fmt"
	"sync"
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

var printMu sync.Mutex

// ============================================================
// Fan-Out / Fan-In — Parallel Processing Pipeline
// ============================================================
//
// The Problem:
//   You have a stream of work items that can be processed independently.
//   A single goroutine processes them sequentially, but each item takes
//   significant time (CPU work, I/O, network call). You need N workers
//   processing in parallel, then collect all results into a single stream.
//
// Real-world example:
//   An image processing service receives upload URLs. Each image needs
//   resizing, watermarking, and compression — CPU-heavy, independent work.
//   Fan-out to GOMAXPROCS workers, fan-in results to a single channel
//   that writes to the CDN upload queue.
//
// The Pattern:
//   Fan-out + Fan-in is two distinct patterns combined:
//
//   Fan-out: multiple goroutines reading from the SAME channel.
//     - One producer sends work to a shared channel
//     - N workers receive from that channel concurrently
//     - The runtime's recvq (receiver wait queue inside hchan) ensures
//       each value goes to EXACTLY one receiver — no duplicates, no locks
//     - The runtime picks a waiting goroutine from recvq in FIFO order
//       (but under contention, which goroutine wins is non-deterministic
//       from the caller's perspective)
//
//   Fan-in: multiple channels merged into one output channel.
//     - Each worker writes results to its own output channel
//     - A merger goroutine (one per input channel) forwards values to
//       a single shared output channel
//     - A sync.WaitGroup tracks when all merger goroutines are done
//     - A separate goroutine waits on the WaitGroup, then closes the
//       output channel — this is the ONLY safe way to close a channel
//       fed by multiple writers
//
//   The combination:
//
//   ┌──────────┐     ┌──────────┐     ┌──────────┐
//   │ Producer │────►│ shared   │     │ merged   │────► Consumer
//   │          │     │ work ch  │     │ output   │
//   └──────────┘     └────┬─────┘     └────▲─────┘
//                    ┌────┼────┐      ┌────┼────┐
//                    ▼    ▼    ▼      │    │    │
//                  ┌───┐┌───┐┌───┐ ┌───┐┌───┐┌───┐
//                  │W1 ││W2 ││W3 │ │ch1││ch2││ch3│
//                  └─┬─┘└─┬─┘└─┬─┘ └─▲─┘└─▲─┘└─▲─┘
//                    │    │    │     │    │    │
//                    └────┴────┴─────┘    │    │
//                    (each worker writes  │    │
//                     to its own channel) ─┘────┘
//
//   FAN-OUT (left side)           FAN-IN (right side)
//   Runtime distributes.          WaitGroup + close().
//
// Why channels work here:
//   1. Fan-out requires zero coordination code. Multiple goroutines
//      receiving from the same channel is safe by design — the runtime's
//      hchan mutex serializes access, and the recvq ensures exactly-once
//      delivery per value.
//   2. Fan-in uses one goroutine per input channel. Each independently
//      does for-range on its input and sends to the shared output.
//      WaitGroup counts completions, and a closer goroutine calls
//      close(out) after wg.Wait() — separating "close" from "send"
//      avoids the classic "send on closed channel" panic.
//   3. This pattern is the building block for parallel pipelines:
//      stage1 → fan-out to N workers → fan-in → stage2.

// fanOutFanIn demonstrates fan-out (1 producer → N workers reading from
// the same channel) and fan-in (N result channels → 1 merged output).
// Returns the collected results and which worker processed each item.
func fanOutFanIn(items []int, numWorkers int) []string {
	workerColors := []string{cyan, yellow, magenta}

	// Fan-out: one shared work channel, multiple workers reading from it
	work := make(chan int)
	go func() {
		for _, item := range items {
			printMu.Lock()
			fmt.Printf("  %s→%s Producer sends item %s%d%s to work channel\n",
				green, reset, magenta, item, reset)
			printMu.Unlock()
			work <- item
		}
		close(work) // producer closes — workers will exit range loop
		printMu.Lock()
		fmt.Printf("  %s✔ Producer done — work channel closed%s\n\n", green, reset)
		printMu.Unlock()
	}()

	// Each worker gets its own result channel
	workerChans := make([]<-chan string, numWorkers)
	for i := 0; i < numWorkers; i++ {
		ch := make(chan string)
		workerChans[i] = ch
		go func(id int, out chan<- string) {
			defer close(out) // worker closes its own output channel
			for item := range work {
				squared := item * item
				printMu.Lock()
				fmt.Printf("  %s%s[W%d]%s ← received %s%d%s → computed %s%d×%d = %d%s\n",
					bold, workerColors[id%len(workerColors)], id, reset,
					magenta, item, reset, dim, item, item, squared, reset)
				printMu.Unlock()
				// Simulate processing — square the value
				result := fmt.Sprintf("worker-%d processed %d → %d", id, item, squared)
				out <- result
			}
		}(i, ch)
	}

	// Fan-in: merge all worker channels into one output channel
	merged := fanIn(workerChans)

	// Collect all results
	printMu.Lock()
	fmt.Printf("%s▸ Consumer Collecting Results%s\n", cyan+bold, reset)
	printMu.Unlock()
	var results []string
	for r := range merged {
		printMu.Lock()
		fmt.Printf("  %s◆%s Consumer received: %s%s%s\n",
			blue+bold, reset, dim, r, reset)
		printMu.Unlock()
		results = append(results, r)
	}
	return results
}

// fanIn merges multiple read-only channels into a single output channel.
// One goroutine per input channel forwards values. A WaitGroup + closer
// goroutine ensures the output channel is closed exactly once after all
// inputs are drained.
func fanIn(channels []<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	// One forwarding goroutine per input channel
	for i, ch := range channels {
		wg.Add(1)
		go func(idx int, c <-chan string) {
			defer wg.Done()
			for v := range c {
				printMu.Lock()
				fmt.Printf("  %s⇉%s  Fan-in: ch[%s%d%s] → merged output\n",
					cyan, reset, magenta, idx, reset)
				printMu.Unlock()
				out <- v
			}
		}(i, ch)
	}

	// Closer goroutine: waits for all forwarders, then closes output.
	// This MUST be in a separate goroutine — if we did wg.Wait() in the
	// current goroutine, we'd block before anyone reads from out.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	numWorkers := 3

	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Fan-Out / Fan-In Pipeline                      %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Pipeline Architecture%s\n", cyan+bold, reset)
	fmt.Printf("  Producer → %s[shared work ch]%s → %s%d workers%s → %s[per-worker ch]%s → Fan-In → Consumer\n\n",
		dim, reset, magenta, numWorkers, reset, dim, reset)

	fmt.Printf("  Fan-out: 1 producer → %s%d%s workers %s(runtime distributes via recvq)%s\n",
		magenta, numWorkers, reset, dim, reset)
	fmt.Printf("  Fan-in:  %s%d%s result channels → 1 merged output %s(WaitGroup + close)%s\n\n",
		magenta, numWorkers, reset, dim, reset)

	fmt.Printf("%s▸ Data Flow%s\n", cyan+bold, reset)

	results := fanOutFanIn(items, numWorkers)

	fmt.Printf("\n%s▸ Final Results%s\n", cyan+bold, reset)
	for _, r := range results {
		fmt.Printf("    %s%s%s\n", dim, r, reset)
	}
	fmt.Printf("\n  Total results: %s%d%s (all %s%d%s items processed exactly once)\n",
		magenta, len(results), reset, magenta, len(items), reset)

	fmt.Printf("\n%s▸ Key Observations%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Fan-out: multiple goroutines safely read from one channel (exactly-once delivery)%s\n", green, reset)
	fmt.Printf("  %s✔ Fan-in: WaitGroup + closer goroutine is the ONLY safe way to close a multi-writer channel%s\n", green, reset)
	fmt.Printf("  %s✔ Each worker has its own output channel — no contention on result writes%s\n", green, reset)
	fmt.Printf("  %s⚠ Order is non-deterministic — items are distributed by the scheduler, not round-robin%s\n", yellow, reset)
}
