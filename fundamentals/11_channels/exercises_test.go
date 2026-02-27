package channels

import (
	"reflect"
	"sort"
	"testing"
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
