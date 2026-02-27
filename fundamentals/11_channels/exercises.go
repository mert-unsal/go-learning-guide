package channels
// ============================================================
// EXERCISES — 11 Channels
// ============================================================
// Exercise 1:
// Sum sends the sum of nums down the channel ch.
// The caller creates the channel and reads the result.
// This is the basic "goroutine returns a value via channel" pattern.
func SumAsync(nums []int, ch chan<- int) {
// TODO: compute sum and send to ch
}
// Exercise 2:
// Generate returns a receive-only channel that sends
// the integers 1..n then closes.
func Generate(n int) <-chan int {
// TODO: make channel, launch goroutine that sends 1..n then closes, return ch
return nil
}
// Exercise 3:
// Pipeline: square each value coming from 'in' and send to returned channel.
// The output channel closes when 'in' closes.
func Square(in <-chan int) <-chan int {
// TODO: make out channel, goroutine that ranges over in, squares, sends
return nil
}
// Exercise 4:
// Merge two channels into one. The output channel closes when BOTH inputs close.
// Use sync.WaitGroup inside.
func Merge(a, b <-chan int) <-chan int {
// TODO: fan-in pattern
return nil
}
// Exercise 5:
// WithTimeout receives from ch but returns (0, false) if nothing
// arrives within maxWaitMs milliseconds.
// Use select with time.After.
func WithTimeout(ch <-chan int, maxWaitMs int) (int, bool) {
// TODO: select { case v := <-ch: return v, true; case <-time.After(...): return 0, false }
return 0, false
}