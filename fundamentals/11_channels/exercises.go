package channels

import (
	"sync"
	"time"
)

// ============================================================
// EXERCISES — 11 Channels
// ============================================================

// Exercise 1:
// Sum sends the sum of nums down the channel ch.
// The caller creates the channel and reads the result.
// This is the basic "goroutine returns a value via channel" pattern.
func SumAsync(nums []int, ch chan<- int) {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	ch <- sum
}

// Exercise 2:
// Generate returns a receive-only channel that sends
// the integers 1..n then closes.
func Generate(n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= n; i++ {
			ch <- i
		}
	}()
	return ch
}

// Exercise 3:
// Pipeline: square each value coming from 'in' and send to returned channel.
// The output channel closes when 'in' closes.
func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * v
		}
	}()
	return out
}

// Exercise 4:
// Merge two channels into one. The output channel closes when BOTH inputs close.
// Use sync.WaitGroup inside.
func Merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	output := func(c <-chan int) {
		defer wg.Done()
		for v := range c {
			out <- v
		}
	}

	go output(a)
	go output(b)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Exercise 5:
// WithTimeout receives from ch but returns (0, false) if nothing
// arrives within maxWaitMs milliseconds.
// Use select with time.After.
func WithTimeout(ch <-chan int, maxWaitMs int) (int, bool) {
	select {
	case v := <-ch:
		return v, true
	case <-time.After(time.Duration(maxWaitMs) * time.Millisecond):
		return 0, false
	}
}
