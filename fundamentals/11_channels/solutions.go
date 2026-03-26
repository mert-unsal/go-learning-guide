package channels

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================
// SOLUTIONS — 11 Channels
// ============================================================
func SumAsyncSolution(nums []int, ch chan<- int) {
sum := 0
for _, v := range nums {
sum += v
}
ch <- sum // send result — caller receives it
}
func GenerateSolution(n int) <-chan int {
ch := make(chan int)
go func() {
defer close(ch) // always close from the sender side
for i := 1; i <= n; i++ {
ch <- i
}
}()
return ch
}
func SquareSolution(in <-chan int) <-chan int {
out := make(chan int)
go func() {
defer close(out)
for v := range in { // range exits when in is closed
out <- v * v
}
}()
return out
}
func MergeSolution(a, b <-chan int) <-chan int {
out := make(chan int)
var wg sync.WaitGroup
forward := func(ch <-chan int) {
defer wg.Done()
for v := range ch {
out <- v
}
}
wg.Add(2)
go forward(a)
go forward(b)
// Close output when both inputs are drained
go func() {
wg.Wait()
close(out)
}()
return out
}
func WithTimeoutSolution(ch <-chan int, maxWaitMs int) (int, bool) {
	select {
	case v := <-ch:
		return v, true
	case <-time.After(time.Duration(maxWaitMs) * time.Millisecond):
		return 0, false
	}
}

// ============================================================
// Solution 6: MergeN — Nil-Channel Select Pattern
// ============================================================
// Key insight: setting a channel to nil in a select makes that case
// block forever — effectively disabling it. This is the idiomatic Go
// pattern for dynamically controlling select behavior.
func MergeNSolution(channels ...<-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		alive := len(channels)
		for alive > 0 {
			for i, ch := range channels {
				if ch == nil {
					continue
				}
				select {
				case v, ok := <-ch:
					if !ok {
						channels[i] = nil
						alive--
					} else {
						out <- v
					}
				default:
				}
			}
			if alive > 0 {
				time.Sleep(time.Microsecond)
			}
		}
	}()
	return out
}

// ============================================================
// Solution 7: Semaphore — Bounded Concurrency
// ============================================================
// The buffered channel acts as a counting semaphore:
//   - sem <- struct{}{} = acquire (blocks when buffer full = limit reached)
//   - <-sem = release (frees a slot for another goroutine)
//
// Performance note: for very high throughput, consider sync.Pool or
// errgroup.SetLimit() instead. Channel semaphores are ~50-100ns per
// acquire/release due to the hchan mutex.
func ProcessWithLimitSolution(items []int, maxConcurrent int, fn func(int) int) []int {
	results := make([]int, len(items))
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx, val int) {
			defer wg.Done()
			defer func() { <-sem }()
			results[idx] = fn(val)
		}(i, item)
	}

	wg.Wait()
	return results
}

// ============================================================
// Solution 8: CloseDrain — Close Semantics
// ============================================================
// Key insight from runtime: close(ch) sets hchan.closed = 1 but does NOT
// flush the buffer. Existing buffered values remain readable (ok=true).
// Only when both closed=1 AND qcount=0 does receive return (zero, false).
func SendAndCloseSolution(values []int, bufSize int) <-chan int {
	ch := make(chan int, bufSize)
	for _, v := range values {
		ch <- v
	}
	close(ch)
	return ch
}

// ============================================================
// Solution 9: TrySend / TryReceive
// ============================================================
// select with default is the non-blocking idiom:
//   - If the channel op can proceed immediately → that case runs
//   - If not → default runs immediately, no blocking
//
// Under the hood: selectgo() reaches Phase 3 (poll), finds nothing ready,
// sees default exists → returns default without entering Phase 4 (park).
func TrySendSolution(ch chan<- int, val int) bool {
	select {
	case ch <- val:
		return true
	default:
		return false
	}
}

func TryReceiveSolution(ch <-chan int) (int, bool) {
	select {
	case v := <-ch:
		return v, true
	default:
		return 0, false
	}
}

// ============================================================
// Solution 10: SafeGenerator — Context-Aware Producer
// ============================================================
// The select checks ctx.Done() alongside the channel send.
// When ctx is cancelled, ctx.Done() becomes ready → goroutine exits.
// defer close(ch) ensures the consumer's range loop terminates.
//
// Without context: goroutine blocks on ch <- i forever if nobody reads.
// With context: cancel() wakes the goroutine via ctx.Done() channel.
func SafeGeneratorSolution(ctx context.Context) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
				i++
			}
		}
	}()
	return ch
}

// ============================================================
// Solution 11: OrDone — Cancellation-Aware Wrapper
// ============================================================
// The double-select pattern is critical:
//   1. Outer select: wait for either ctx.Done() or a value from in
//   2. Inner select: when forwarding to out, also check ctx.Done()
//
// Why inner select? Without it, if out's consumer stops reading,
// the goroutine blocks on `out <- v` and never checks ctx.Done().
// The inner select ensures cancellation is respected at EVERY blocking point.
func OrDoneSolution(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

// ============================================================
// Solution Bonus: ChannelCounter vs AtomicCounter
// ============================================================
// ChannelCounter: each increment is a channel send → hchan mutex + memcopy
//   ~50-100ns per increment. Correct but slow for this use case.
//
// AtomicCounter: each increment is atomic.AddInt64 → single LOCK XADD instruction
//   ~5-15ns per increment. 5-10x faster for simple counters.
//
// Lesson: "Channels orchestrate; mutexes serialize."
// Use channels for goroutine communication. Use atomics for shared counters.
func ChannelCounterSolution(workers, incrementsPerWorker int) int {
	ch := make(chan struct{}, 100)
	done := make(chan int)

	go func() {
		count := 0
		for range ch {
			count++
		}
		done <- count
	}()

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerWorker; j++ {
				ch <- struct{}{}
			}
		}()
	}

	wg.Wait()
	close(ch)
	return <-done
}

func AtomicCounterSolution(workers, incrementsPerWorker int) int64 {
	var count int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerWorker; j++ {
				atomic.AddInt64(&count, 1)
			}
		}()
	}

	wg.Wait()
	return count
}