package channels

// ============================================================
// EXERCISES — 11 Channels
// ============================================================

// Exercise 1:
// SumAsync sends the sum of nums down the channel ch.
// The caller creates the channel and reads the result.
// This is the basic "goroutine returns a value via channel" pattern.
func SumAsync(nums []int, ch chan<- int) {
	// TODO: compute sum, send it on ch
	panic("not implemented")
}

// Exercise 2:
// Generate returns a receive-only channel that sends
// the integers 1..n then closes.
func Generate(n int) <-chan int {
	// TODO: create channel, launch goroutine to send 1..n, close channel
	panic("not implemented")
}

// Exercise 3:
// Pipeline: square each value coming from 'in' and send to returned channel.
// The output channel closes when 'in' closes.
func Square(in <-chan int) <-chan int {
	// TODO: create output channel, range over in, send v*v, close output
	panic("not implemented")
}

// Exercise 4:
// Merge two channels into one. The output channel closes when BOTH inputs close.
// Hint: use sync.WaitGroup inside.
func Merge(a, b <-chan int) <-chan int {
	// TODO: fan-in pattern — one goroutine per input, WaitGroup to close output
	panic("not implemented")
}

// Exercise 5:
// WithTimeout receives from ch but returns (0, false) if nothing
// arrives within maxWaitMs milliseconds.
// Use select with time.After.
func WithTimeout(ch <-chan int, maxWaitMs int) (int, bool) {
	// TODO: select between ch receive and time.After timeout
	panic("not implemented")
}
