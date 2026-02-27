package functions

import "testing"

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
		gotMin, gotMax := MinMaxSolution(tt.nums)
		if gotMin != tt.wantMin || gotMax != tt.wantMax {
			t.Errorf("MinMax(%v) = (%d,%d), want (%d,%d)", tt.nums, gotMin, gotMax, tt.wantMin, tt.wantMax)
		}
	}
}

func TestSum(t *testing.T) {
	if got := SumSolution(1, 2, 3); got != 6 {
		t.Errorf("Sum(1,2,3) = %d, want 6", got)
	}
	if got := SumSolution(); got != 0 {
		t.Errorf("Sum() = %d, want 0", got)
	}
	if got := SumSolution(10); got != 10 {
		t.Errorf("Sum(10) = %d, want 10", got)
	}
}

func TestApply(t *testing.T) {
	double := func(x int) int { return x * 2 }
	got := ApplySolution([]int{1, 2, 3}, double)
	want := []int{2, 4, 6}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("Apply double: got[%d]=%d, want %d", i, got[i], v)
		}
	}
}

func TestMakeAdder(t *testing.T) {
	add5 := MakeAdderSolution(5)
	if got := add5(3); got != 8 {
		t.Errorf("add5(3) = %d, want 8", got)
	}
	if got := add5(10); got != 15 {
		t.Errorf("add5(10) = %d, want 15", got)
	}
	add10 := MakeAdderSolution(10)
	if got := add10(0); got != 10 {
		t.Errorf("add10(0) = %d, want 10", got)
	}
}

func TestFibonacci(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 0}, {1, 1}, {2, 1}, {5, 5}, {10, 55},
	}
	for _, tt := range tests {
		if got := FibonacciSolution(tt.n); got != tt.want {
			t.Errorf("Fibonacci(%d) = %d, want %d", tt.n, got, tt.want)
		}
		if got := FibonacciMemoSolution(tt.n); got != tt.want {
			t.Errorf("FibonacciMemo(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}
