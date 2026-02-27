package arrays_slices

import (
	"reflect"
	"testing"
)

func TestReverseSlice(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{42}, []int{42}},
		{[]int{}, []int{}},
	}
	for _, tt := range tests {
		s := make([]int, len(tt.input))
		copy(s, tt.input)
		ReverseSlice(s)
		if !reflect.DeepEqual(s, tt.want) {
			t.Errorf("❌ ReverseSlice(%v) = %v, want %v  ← Hint: two-pointer swap", tt.input, s, tt.want)
		} else {
			t.Logf("✅ ReverseSlice(%v) = %v", tt.input, s)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{1, 1, 2, 3, 3, 4}, []int{1, 2, 3, 4}},
		{[]int{1, 1, 1}, []int{1}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{}, []int{}},
	}
	for _, tt := range tests {
		s := make([]int, len(tt.input))
		copy(s, tt.input)
		got := RemoveDuplicates(s)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("❌ RemoveDuplicates(%v) = %v, want %v  ← Hint: write-pointer pattern", tt.input, got, tt.want)
		} else {
			t.Logf("✅ RemoveDuplicates(%v) = %v", tt.input, got)
		}
	}
}

func TestMake2D(t *testing.T) {
	m := Make2D(3, 4)
	if len(m) != 3 {
		t.Fatalf("❌ Make2D rows = %d, want 3", len(m))
	}
	for i, row := range m {
		if len(row) != 4 {
			t.Errorf("❌ row[%d] len = %d, want 4", i, len(row))
		}
	}
	m[0][0] = 99
	if m[1][0] == 99 {
		t.Error("❌ rows share underlying array — they must be independent  ← Hint: allocate each row separately")
	} else {
		t.Logf("✅ Make2D(3,4) = %dx%d grid with independent rows", len(m), len(m[0]))
	}
}

func TestRotateLeft(t *testing.T) {
	tests := []struct {
		input []int
		k     int
		want  []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{3, 4, 5, 1, 2}},
		{[]int{1, 2, 3, 4, 5}, 0, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 5, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3}, 1, []int{2, 3, 1}},
	}
	for _, tt := range tests {
		s := make([]int, len(tt.input))
		copy(s, tt.input)
		RotateLeft(s, tt.k)
		if !reflect.DeepEqual(s, tt.want) {
			t.Errorf("❌ RotateLeft(%v, %d) = %v, want %v  ← Hint: three-reversal trick", tt.input, tt.k, s, tt.want)
		} else {
			t.Logf("✅ RotateLeft(%v, %d) = %v", tt.input, tt.k, s)
		}
	}
}

func TestFilter(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	got := Filter([]int{1, 2, 3, 4, 5, 6}, isEven)
	want := []int{2, 4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("❌ Filter(isEven, [1..6]) = %v, want %v", got, want)
	} else {
		t.Logf("✅ Filter(isEven, [1..6]) = %v", got)
	}

	isPositive := func(n int) bool { return n > 0 }
	got2 := Filter([]int{-2, -1, 0, 1, 2}, isPositive)
	want2 := []int{1, 2}
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("❌ Filter(isPositive) = %v, want %v", got2, want2)
	} else {
		t.Logf("✅ Filter(isPositive, [-2,-1,0,1,2]) = %v", got2)
	}
}

func TestMergeSorted(t *testing.T) {
	tests := []struct {
		a, b []int
		want []int
	}{
		{[]int{1, 3, 5}, []int{2, 4, 6}, []int{1, 2, 3, 4, 5, 6}},
		{[]int{1, 2, 3}, []int{}, []int{1, 2, 3}},
		{[]int{}, []int{4, 5}, []int{4, 5}},
		{[]int{1}, []int{1}, []int{1, 1}},
	}
	for _, tt := range tests {
		got := MergeSorted(tt.a, tt.b)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("❌ MergeSorted(%v, %v) = %v, want %v  ← Hint: two-pointer merge", tt.a, tt.b, got, tt.want)
		} else {
			t.Logf("✅ MergeSorted(%v, %v) = %v", tt.a, tt.b, got)
		}
	}
}
