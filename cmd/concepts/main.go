package main

import (
	"fmt"
	"os"
	"strings"

	"gointerviewprep/fundamentals/11_channels/concepts"
)

// Run individual channel concept demos or all of them.
//
// Usage:
//
//	go run ./cmd/concepts              # list all available concepts
//	go run ./cmd/concepts all          # run all demos
//	go run ./cmd/concepts 01           # run concept 01 (worker event loop)
//	go run ./cmd/concepts 04 09 16     # run specific concepts

var demos = []struct {
	id   string
	name string
	fn   func()
}{
	{"01", "Worker Event Loop", concepts.DemonstrateWorkerEventLoop},
	{"02", "Heartbeat Monitor", concepts.DemonstrateHeartbeatMonitor},
	{"03", "First Response Wins", concepts.DemonstrateFirstResponseWins},
	{"04", "Token Bucket Rate Limiter", concepts.DemonstrateTokenBucketRateLimiter},
	{"05", "Buffered Channel Lifecycle", concepts.DemonstrateBufferedChannelLifecycle},
	{"06", "Unbuffered Direct Transfer", concepts.DemonstrateUnbufferedDirectTransfer},
	{"07", "Sender Blocks, Receiver Wakes", concepts.DemonstrateSenderBlocksReceiverWakes},
	{"08", "Select Multi-Channel", concepts.DemonstrateSelectMultiChannel},
	{"09", "Close Wakes Receivers", concepts.DemonstrateCloseWakesReceivers},
	{"10", "Done / Cancellation", concepts.DemonstrateDoneCancellation},
	{"11", "Worker Pool Backpressure", concepts.DemonstrateWorkerPoolBackpressure},
	{"12", "Semaphore Bounded Concurrency", concepts.DemonstrateSemaphoreBoundedConcurrency},
	{"13", "Fan-Out / Fan-In", concepts.DemonstrateFanOutFanIn},
	{"14", "Pipeline", concepts.DemonstratePipeline},
	{"15", "Timeout with Context", concepts.DemonstrateTimeoutWithContext},
	{"16", "Or-Done (Double Select)", concepts.DemonstrateOrDone},
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Channel Concepts — 16 runnable demos")
		fmt.Println("Usage: go run ./cmd/concepts [id...] or 'all'")
		fmt.Println()
		for _, d := range demos {
			fmt.Printf("  %s  %s\n", d.id, d.name)
		}
		return
	}

	runAll := len(args) == 1 && strings.ToLower(args[0]) == "all"

	for _, d := range demos {
		if runAll || contains(args, d.id) {
			fmt.Printf("\n══════ %s: %s ══════\n\n", d.id, d.name)
			d.fn()
		}
	}
}

func contains(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}
