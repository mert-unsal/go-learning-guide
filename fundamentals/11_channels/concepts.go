// Package channels covers Go channels: buffered vs unbuffered,
// select, done patterns, fan-out/fan-in, and common idioms.
package channels

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================
// 1. CHANNEL BASICS
// ============================================================
// A channel is a typed conduit for communication between goroutines.
// make(chan T)       → unbuffered: sender blocks until receiver is ready
// make(chan T, n)    → buffered:   sender blocks only when buffer is full
// close(ch)         → signals no more values will be sent
// <-ch              → receive
// ch <- val         → send

func DemonstrateBasics() {
	// Unbuffered channel — synchronization point
	ch := make(chan int)

	go func() {
		ch <- 42 // blocks until someone receives
	}()

	val := <-ch // blocks until something is sent
	fmt.Println("Received:", val)

	// Buffered channel — asynchronous up to buffer size
	bch := make(chan string, 3)
	bch <- "one" // doesn't block (buffer not full)
	bch <- "two"
	bch <- "three"
	// bch <- "four" // would block — buffer is full

	fmt.Println(<-bch) // "one"
	fmt.Println(<-bch) // "two"
	fmt.Println(<-bch) // "three"
}

// ============================================================
// 2. RANGE OVER CHANNEL
// ============================================================
// range over a channel receives values until it's CLOSED.
// Always close channels from the SENDER side.

func generate(nums ...int) <-chan int { // returns a receive-only channel
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out) // signal that no more values will be sent
	}()
	return out
}

func DemonstrateRange() {
	for n := range generate(2, 3, 4, 5) {
		fmt.Print(n*n, " ") // 4 9 16 25
	}
	fmt.Println()
}

// ============================================================
// 3. SELECT — multiplex channels
// ============================================================
// select waits on multiple channel operations.
// Picks a random case if multiple are ready.

func DemonstrateSelect() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Millisecond)
		ch1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Millisecond)
		ch2 <- "two"
	}()

	// Receive from whichever is ready first
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received from ch1:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received from ch2:", msg2)
		}
	}
}

// ============================================================
// 4. DONE PATTERN (Cancellation)
// ============================================================
// Use a 'done' channel to signal goroutines to stop.

func DemonstrateDonePattern() {
	done := make(chan struct{})
	results := make(chan int)

	go func() {
		i := 0
		for {
			select {
			case <-done: // cancellation signal
				close(results)
				return
			case results <- i: // send work
				i++
			}
		}
	}()

	// Receive first 5 results
	for i := 0; i < 5; i++ {
		fmt.Print(<-results, " ")
	}
	fmt.Println()

	close(done) // signal goroutine to stop
	time.Sleep(1 * time.Millisecond)
}

// ============================================================
// 5. FAN-OUT / FAN-IN PATTERN
// ============================================================

// Fan-out: distribute work across multiple goroutines
// Fan-in: merge results from multiple goroutines into one channel

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	merged := make(chan int)

	// Start output goroutine for each input channel
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			merged <- n
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	// Start a goroutine to close merged once all inputs are done
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// ============================================================
// 6. TIMEOUT PATTERN
// ============================================================

func fetchWithTimeout(duration time.Duration) (string, error) {
	result := make(chan string, 1)

	go func() {
		// Simulate slow operation
		time.Sleep(duration * 2)
		result <- "data"
	}()

	select {
	case data := <-result:
		return data, nil
	case <-time.After(duration):
		return "", fmt.Errorf("timeout after %v", duration)
	}
}

func DemonstrateTimeout() {
	_, err := fetchWithTimeout(10 * time.Millisecond)
	if err != nil {
		fmt.Println("Timeout:", err)
	}
}

// ============================================================
// 7. PIPELINE PATTERN
// ============================================================

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func DemonstratePipeline() {
	// Pipeline: generate → square → square
	c := generate(2, 3, 4)
	sq1 := square(c)
	sq2 := square(sq1) // 2^4=16, 3^4=81, 4^4=256

	for n := range sq2 {
		fmt.Print(n, " ")
	}
	fmt.Println()
}

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Channel Basics ===")
	DemonstrateBasics()
	fmt.Println("\n=== Range Over Channel ===")
	DemonstrateRange()
	fmt.Println("\n=== Select ===")
	DemonstrateSelect()
	fmt.Println("\n=== Done Pattern ===")
	DemonstrateDonePattern()
	fmt.Println("\n=== Timeout ===")
	DemonstrateTimeout()
	fmt.Println("\n=== Pipeline ===")
	DemonstratePipeline()
}
