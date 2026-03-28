package channels

import (
	"context"
	"reflect"
	"runtime"
	"sort"
	"sync/atomic"
	"testing"
	"time"
)

func TestSumAsync(t *testing.T) {
	ch := make(chan int, 1)
	go SumAsync([]int{1, 2, 3, 4, 5}, ch)
	got := <-ch
	if got != 15 {
		t.Errorf("❌ SumAsync([1..5]) = %d, want 15  ← Hint: compute sum and send on ch", got)
	} else {
		t.Logf("✅ SumAsync([1..5]) = %d", got)
	}
}

func TestGenerate(t *testing.T) {
	ch := Generate(5)
	var result []int
	for v := range ch {
		result = append(result, v)
	}
	want := []int{1, 2, 3, 4, 5}
	match := len(result) == len(want)
	if match {
		for i, v := range want {
			if result[i] != v {
				match = false
				break
			}
		}
	}
	if !match {
		t.Errorf("❌ Generate(5) = %v, want %v  ← Hint: close the channel when done", result, want)
	} else {
		t.Logf("✅ Generate(5) = %v", result)
	}
}

func TestSquare(t *testing.T) {
	in := Generate(4)
	out := Square(in)
	want := []int{1, 4, 9, 16}
	for _, w := range want {
		got := <-out
		if got != w {
			t.Errorf("❌ Square: got %d, want %d", got, w)
		} else {
			t.Logf("✅ Square: %d²=%d", w, got)
		}
	}
}

func TestMerge(t *testing.T) {
	make123 := func() <-chan int {
		ch := make(chan int, 3)
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
		return ch
	}
	merged := Merge(make123(), make123())
	var result []int
	for v := range merged {
		result = append(result, v)
	}
	sort.Ints(result)
	want := []int{1, 1, 2, 2, 3, 3}
	match := reflect.DeepEqual(result, want)
	if !match {
		t.Errorf("❌ Merge = %v, want %v  ← Hint: fan-in with WaitGroup", result, want)
	} else {
		t.Logf("✅ Merge([1,2,3],[1,2,3]) = %v", result)
	}
}

func TestWithTimeout(t *testing.T) {
	fast := make(chan int, 1)
	fast <- 42
	v, ok := WithTimeout(fast, 100)
	if !ok || v != 42 {
		t.Errorf("❌ WithTimeout fast: got (%d,%v), want (42,true)", v, ok)
	} else {
		t.Logf("✅ WithTimeout received %d before timeout", v)
	}

	slow := make(chan int)
	v, ok = WithTimeout(slow, 10)
	if ok {
		t.Errorf("❌ WithTimeout slow: expected timeout, got value %d", v)
	} else {
		t.Logf("✅ WithTimeout correctly timed out after 10ms")
	}
}

// ============================================================
// Tests for Exercise 6: MergeN
// ============================================================

func TestMergeN(t *testing.T) {
	// Create 3 channels with known values
	makeCh := func(vals ...int) <-chan int {
		ch := make(chan int, len(vals))
		for _, v := range vals {
			ch <- v
		}
		close(ch)
		return ch
	}

	ch1 := makeCh(1, 2, 3)
	ch2 := makeCh(10, 20)
	ch3 := makeCh(100)

	merged := MergeN(ch1, ch2, ch3)
	var result []int
	for v := range merged {
		result = append(result, v)
	}

	sort.Ints(result)
	want := []int{1, 2, 3, 10, 20, 100}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("❌ MergeN = %v, want %v", result, want)
	} else {
		t.Logf("✅ MergeN merged 3 channels: %v", result)
	}
}

func TestMergeN_Empty(t *testing.T) {
	// Edge case: no channels
	merged := MergeN()
	var result []int
	for v := range merged {
		result = append(result, v)
	}
	if len(result) != 0 {
		t.Errorf("❌ MergeN() with no channels should return empty, got %v", result)
	} else {
		t.Logf("✅ MergeN with zero channels returns empty")
	}
}

// ============================================================
// Tests for Exercise 7: Semaphore
// ============================================================

func TestProcessWithLimit(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	double := func(x int) int { return x * 2 }

	results := ProcessWithLimit(items, 3, double)
	want := []int{2, 4, 6, 8, 10, 12, 14, 16}

	if !reflect.DeepEqual(results, want) {
		t.Errorf("❌ ProcessWithLimit = %v, want %v", results, want)
	} else {
		t.Logf("✅ ProcessWithLimit doubled all items: %v", results)
	}
}

func TestProcessWithLimit_ConcurrencyBound(t *testing.T) {
	// Verify that no more than maxConcurrent goroutines run at once
	var currentConcurrency int64
	var maxObserved int64
	maxConcurrent := 3

	items := make([]int, 20)
	for i := range items {
		items[i] = i
	}

	slowFn := func(x int) int {
		cur := atomic.AddInt64(&currentConcurrency, 1)
		// Track max observed concurrency
		for {
			old := atomic.LoadInt64(&maxObserved)
			if cur <= old || atomic.CompareAndSwapInt64(&maxObserved, old, cur) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond) // simulate work
		atomic.AddInt64(&currentConcurrency, -1)
		return x
	}

	ProcessWithLimit(items, maxConcurrent, slowFn)

	observed := atomic.LoadInt64(&maxObserved)
	if observed > int64(maxConcurrent) {
		t.Errorf("❌ Max concurrency was %d, expected ≤ %d", observed, maxConcurrent)
	} else {
		t.Logf("✅ Max concurrency observed: %d (limit: %d)", observed, maxConcurrent)
	}
}

// ============================================================
// Tests for Exercise 8: CloseDrain
// ============================================================

func TestSendAndClose_DrainBuffer(t *testing.T) {
	values := []int{10, 20, 30}
	ch := SendAndClose(values, 5)

	// Should receive all buffered values after close
	var received []int
	for v := range ch {
		received = append(received, v)
	}

	if !reflect.DeepEqual(received, values) {
		t.Errorf("❌ Drain after close = %v, want %v", received, values)
	} else {
		t.Logf("✅ Drained all %d values after close: %v", len(received), received)
	}
}

func TestSendAndClose_ZeroAfterDrain(t *testing.T) {
	ch := SendAndClose([]int{42}, 5)
	<-ch // drain the one value

	// Now channel is closed AND empty — should return zero, false
	v, ok := <-ch
	if ok {
		t.Errorf("❌ Expected ok=false after drain, got (%d, true)", v)
	} else if v != 0 {
		t.Errorf("❌ Expected zero value after drain, got %d", v)
	} else {
		t.Logf("✅ After drain: received (0, false) — correct close semantics")
	}
}

func TestSendAndClose_MultipleZeroReads(t *testing.T) {
	ch := SendAndClose([]int{}, 5) // empty, already closed

	// Reading from closed empty channel should return (0, false) repeatedly
	for i := 0; i < 3; i++ {
		v, ok := <-ch
		if ok || v != 0 {
			t.Errorf("❌ Read %d: got (%d, %v), want (0, false)", i, v, ok)
			return
		}
	}
	t.Logf("✅ Multiple reads from closed empty channel all return (0, false)")
}

// ============================================================
// Tests for Exercise 9: TrySend / TryReceive
// ============================================================

func TestTrySend(t *testing.T) {
	ch := make(chan int, 1)

	// Should succeed — buffer has space
	if !TrySend(ch, 42) {
		t.Errorf("❌ TrySend to empty buffered channel should succeed")
	} else {
		t.Logf("✅ TrySend succeeded on empty buffer")
	}

	// Should fail — buffer is full
	if TrySend(ch, 99) {
		t.Errorf("❌ TrySend to full buffered channel should fail")
	} else {
		t.Logf("✅ TrySend correctly returned false on full buffer")
	}
}

func TestTryReceive(t *testing.T) {
	ch := make(chan int, 1)

	// Should fail — nothing in channel
	v, ok := TryReceive(ch)
	if ok {
		t.Errorf("❌ TryReceive from empty channel should fail, got %d", v)
	} else {
		t.Logf("✅ TryReceive correctly returned false on empty channel")
	}

	// Put a value in, then try receive
	ch <- 42
	v, ok = TryReceive(ch)
	if !ok || v != 42 {
		t.Errorf("❌ TryReceive = (%d, %v), want (42, true)", v, ok)
	} else {
		t.Logf("✅ TryReceive got %d from buffered channel", v)
	}
}

// ============================================================
// Tests for Exercise 10: LeakCheck — SafeGenerator
// ============================================================

func TestSafeGenerator(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := SafeGenerator(ctx)

	// Read a few values
	var received []int
	for i := 0; i < 5; i++ {
		v := <-ch
		received = append(received, v)
	}

	want := []int{0, 1, 2, 3, 4}
	if !reflect.DeepEqual(received, want) {
		t.Errorf("❌ SafeGenerator first 5 = %v, want %v", received, want)
	} else {
		t.Logf("✅ SafeGenerator produced: %v", received)
	}

	cancel() // cancel the context

	// After cancel, the channel should eventually close
	// Drain any remaining values
	drained := 0
	for range ch {
		drained++
		if drained > 100 {
			t.Fatalf("❌ SafeGenerator didn't stop after cancel — goroutine leak!")
		}
	}
	t.Logf("✅ SafeGenerator stopped after cancel (drained %d extra values)", drained)
}

func TestSafeGenerator_NoLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	ch := SafeGenerator(ctx)
	<-ch // read one value
	cancel()

	// Drain to let goroutine exit
	for range ch {
	}

	// Give goroutine time to fully exit
	time.Sleep(50 * time.Millisecond)

	after := runtime.NumGoroutine()
	if after > before+1 {
		t.Errorf("❌ Goroutine leak! Before: %d, After: %d", before, after)
	} else {
		t.Logf("✅ No goroutine leak: before=%d, after=%d", before, after)
	}
}

// ============================================================
// Tests for Exercise 11: OrDone
// ============================================================

func TestOrDone_ForwardsValues(t *testing.T) {
	ctx := context.Background()
	in := make(chan int, 5)
	in <- 10
	in <- 20
	in <- 30
	close(in)

	out := OrDone(ctx, in)
	var result []int
	for v := range out {
		result = append(result, v)
	}

	want := []int{10, 20, 30}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("❌ OrDone = %v, want %v", result, want)
	} else {
		t.Logf("✅ OrDone forwarded all values: %v", result)
	}
}

func TestOrDone_RespectsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Channel that never closes — without OrDone, we'd block forever
	in := make(chan int)

	out := OrDone(ctx, in)

	// Send one value
	go func() { in <- 42 }()
	v := <-out
	if v != 42 {
		t.Errorf("❌ OrDone first value = %d, want 42", v)
	}

	// Cancel context — out should close even though in never closes
	cancel()

	// out should close promptly
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
	select {
	case _, ok := <-out:
		if ok {
			t.Logf("⚠️  OrDone sent an extra value before closing")
		} else {
			t.Logf("✅ OrDone closed output after context cancel")
		}
	case <-timer.C:
		t.Errorf("❌ OrDone didn't close within 500ms after cancel — possible goroutine leak")
	}
}

// ============================================================
// Tests for Exercise 12: DualTimeoutWorker
// ============================================================

func TestDualTimeoutWorker_ProcessesAll(t *testing.T) {
	// All values arrive quickly — worker should process all before any timeout
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	double := func(x int) int { return x * 2 }
	results := DualTimeoutWorker(ch, double, 500*time.Millisecond, 5*time.Second)

	want := []int{2, 4, 6}
	if !reflect.DeepEqual(results, want) {
		t.Errorf("❌ DualTimeoutWorker = %v, want %v", results, want)
	} else {
		t.Logf("✅ Processed all values: %v", results)
	}
}

func TestDualTimeoutWorker_InactivityTimeout(t *testing.T) {
	// Channel stays open but no values sent — should exit via inactivity
	ch := make(chan int)

	identity := func(x int) int { return x }
	start := time.Now()
	results := DualTimeoutWorker(ch, identity, 100*time.Millisecond, 5*time.Second)
	elapsed := time.Since(start)

	if len(results) != 0 {
		t.Errorf("❌ Expected empty results, got %v", results)
	}
	if elapsed > 500*time.Millisecond {
		t.Errorf("❌ Took %v — should have exited via inactivity (~100ms), not hard deadline", elapsed)
	} else {
		t.Logf("✅ Inactivity timeout fired after %v (no values received)", elapsed)
	}
}

func TestDualTimeoutWorker_HardDeadline(t *testing.T) {
	// Values keep coming forever — should exit via hard deadline
	ch := make(chan int)
	go func() {
		i := 0
		for {
			ch <- i // never stops sending
			i++
			time.Sleep(10 * time.Millisecond)
		}
	}()

	identity := func(x int) int { return x }
	start := time.Now()
	results := DualTimeoutWorker(ch, identity, 500*time.Millisecond, 200*time.Millisecond)
	elapsed := time.Since(start)

	if elapsed > 500*time.Millisecond {
		t.Errorf("❌ Took %v — hard deadline (200ms) should have stopped it sooner", elapsed)
	}
	if len(results) == 0 {
		t.Errorf("❌ Should have processed some values before deadline")
	} else {
		t.Logf("✅ Hard deadline fired after %v — processed %d values: %v", elapsed, len(results), results)
	}
}

func TestDualTimeoutWorker_InactivityResetsOnMessage(t *testing.T) {
	// Send values spaced at 50ms with 100ms inactivity timeout
	// Values should keep the timer alive, then inactivity fires after last value
	ch := make(chan int)
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
			time.Sleep(50 * time.Millisecond) // well within 100ms inactivity
		}
		// stop sending — inactivity should fire ~100ms after last send
	}()

	identity := func(x int) int { return x }
	start := time.Now()
	results := DualTimeoutWorker(ch, identity, 100*time.Millisecond, 5*time.Second)
	elapsed := time.Since(start)

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(results, want) {
		t.Errorf("❌ DualTimeoutWorker = %v, want %v", results, want)
	}
	// Should take ~150ms (3 sends at 50ms spacing) + ~100ms inactivity = ~250ms
	if elapsed > 1*time.Second {
		t.Errorf("❌ Took %v — inactivity timer wasn't resetting properly", elapsed)
	} else {
		t.Logf("✅ Timer reset correctly — processed %v in %v, then inactivity fired", results, elapsed)
	}
}

// ============================================================
// Tests for Bonus: ChannelCounter vs AtomicCounter
// ============================================================

func TestChannelCounter(t *testing.T) {
	got := ChannelCounter(10, 100)
	want := 1000
	if got != want {
		t.Errorf("❌ ChannelCounter(10, 100) = %d, want %d", got, want)
	} else {
		t.Logf("✅ ChannelCounter: %d increments correct", got)
	}
}

func TestAtomicCounter(t *testing.T) {
	got := AtomicCounter(10, 100)
	want := int64(1000)
	if got != want {
		t.Errorf("❌ AtomicCounter(10, 100) = %d, want %d", got, want)
	} else {
		t.Logf("✅ AtomicCounter: %d increments correct", got)
	}
}
