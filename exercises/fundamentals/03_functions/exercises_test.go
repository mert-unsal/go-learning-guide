package functions

import (
	"fmt"
	"testing"
)

func TestMinMax(t *testing.T) {
	tests := []struct {
		nums    []int
		wantMin int
		wantMax int
	}{
		{[]int{3, 1, 4, 1, 5, 9}, 1, 9},
		{[]int{-5, 0, 5}, -5, 5},
		{[]int{42}, 42, 42},
		{[]int{}, 0, 0},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("MinMax(%v)", tt.nums), func(t *testing.T) {
			gotMin, gotMax := MinMax(tt.nums)
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("❌ MinMax(%v) = (%d,%d), want (%d,%d)", tt.nums, gotMin, gotMax, tt.wantMin, tt.wantMax)
			} else {
				t.Logf("✅ MinMax(%v) = (%d,%d)", tt.nums, gotMin, gotMax)
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		args []int
		want int
	}{
		{[]int{1, 2, 3}, 6},
		{[]int{}, 0},
		{[]int{10}, 10},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Sum(%v)", tt.args), func(t *testing.T) {
			got := Sum(tt.args...)
			if got != tt.want {
				t.Errorf("❌ Sum(%v) = %d, want %d", tt.args, got, tt.want)
			} else {
				t.Logf("✅ Sum(%v) = %d", tt.args, got)
			}
		})
	}
}

func TestApply(t *testing.T) {
	double := func(x int) int { return x * 2 }
	got := Apply([]int{1, 2, 3}, double)
	want := []int{2, 4, 6}
	match := len(got) == len(want)
	if match {
		for i, v := range want {
			if got[i] != v {
				match = false
				break
			}
		}
	}
	if !match {
		t.Errorf("❌ Apply(double, [1,2,3]) = %v, want %v", got, want)
	} else {
		t.Logf("✅ Apply(double, [1,2,3]) = %v", got)
	}
}

func TestMakeAdder(t *testing.T) {
	add5 := MakeAdder(5)
	tests := []struct{ in, want int }{{3, 8}, {10, 15}, {0, 5}}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("MakeAdder(5)(%d)", tt.in), func(t *testing.T) {
			got := add5(tt.in)
			if got != tt.want {
				t.Errorf("❌ MakeAdder(5)(%d) = %d, want %d", tt.in, got, tt.want)
			} else {
				t.Logf("✅ MakeAdder(5)(%d) = %d", tt.in, got)
			}
		})
	}
}

func TestFibonacci(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 0}, {1, 1}, {2, 1}, {5, 5}, {10, 55},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Fibonacci(%d)", tt.n), func(t *testing.T) {
			got := Fibonacci(tt.n)
			if got != tt.want {
				t.Errorf("❌ Fibonacci(%d) = %d, want %d", tt.n, got, tt.want)
			} else {
				t.Logf("✅ Fibonacci(%d) = %d", tt.n, got)
			}
		})
	}
}

func TestFibonacciMemo(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 0}, {1, 1}, {2, 1}, {5, 5}, {10, 55},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("FibonacciMemo(%d)", tt.n), func(t *testing.T) {
			got := FibonacciMemo(tt.n)
			if got != tt.want {
				t.Errorf("❌ FibonacciMemo(%d) = %d, want %d", tt.n, got, tt.want)
			} else {
				t.Logf("✅ FibonacciMemo(%d) = %d", tt.n, got)
			}
		})
	}
}
