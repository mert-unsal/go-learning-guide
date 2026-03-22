package arrays_slices

import (
	"reflect"
	"testing"
)

// ============================================================
// PART A — Algorithm Pattern Tests
// ============================================================

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

// ============================================================
// PART B — Slice Internals Tests
// ============================================================

func TestSafeDelete(t *testing.T) {
	original := []int{10, 20, 30, 40, 50}
	originalCopy := make([]int, len(original))
	copy(originalCopy, original)

	result := SafeDelete(original, 2)

	// Check the result is correct
	want := []int{10, 20, 40, 50}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("❌ SafeDelete([10,20,30,40,50], 2) = %v, want %v", result, want)
		return
	}

	// Check original is NOT modified — this is the key lesson
	if !reflect.DeepEqual(original, originalCopy) {
		t.Errorf("❌ Original was modified: %v, want %v  ← Hint: copy to a new backing array before deleting", original, originalCopy)
	} else {
		t.Logf("✅ SafeDelete preserved original and returned %v", result)
	}

	// Edge case: delete first element
	result2 := SafeDelete([]int{1, 2, 3}, 0)
	if !reflect.DeepEqual(result2, []int{2, 3}) {
		t.Errorf("❌ SafeDelete([1,2,3], 0) = %v, want [2,3]", result2)
	}

	// Edge case: delete last element
	result3 := SafeDelete([]int{1, 2, 3}, 2)
	if !reflect.DeepEqual(result3, []int{1, 2}) {
		t.Errorf("❌ SafeDelete([1,2,3], 2) = %v, want [1,2]", result3)
	}
}

func TestCopySlice(t *testing.T) {
	src := []int{10, 20, 30, 40, 50}
	dst := CopySlice(src)

	// Check elements match
	if !reflect.DeepEqual(dst, src) {
		t.Errorf("❌ CopySlice(%v) = %v, want same elements", src, dst)
		return
	}

	// Check independence — modifying dst must not affect src
	dst[0] = 999
	if src[0] == 999 {
		t.Error("❌ Modifying copy changed the original  ← Hint: must allocate new backing array with make(), then copy()")
	} else {
		t.Logf("✅ CopySlice produced independent copy of %v", src)
	}

	// Check len matches
	if len(dst) != len(src) {
		t.Errorf("❌ len mismatch: got %d, want %d", len(dst), len(src))
	}

	// Edge case: empty slice
	empty := CopySlice([]int{})
	if len(empty) != 0 {
		t.Errorf("❌ CopySlice([]) len = %d, want 0", len(empty))
	}
}

func TestNilVsEmpty(t *testing.T) {
	nilSlice, emptySlice := NilVsEmpty()

	if nilSlice != nil {
		t.Error("❌ First return value should be nil  ← Hint: var s []int (don't initialize)")
	} else {
		t.Log("✅ First return value is nil")
	}

	if emptySlice == nil {
		t.Error("❌ Second return value should be non-nil  ← Hint: use []int{} or make([]int, 0)")
	} else {
		t.Log("✅ Second return value is non-nil")
	}

	// Both should have len=0 and cap=0
	if len(nilSlice) != 0 || len(emptySlice) != 0 {
		t.Errorf("❌ Both slices should have len=0, got nil=%d, empty=%d", len(nilSlice), len(emptySlice))
	}

	// Both should work with append
	nilSlice = append(nilSlice, 1)
	emptySlice = append(emptySlice, 1)
	if nilSlice[0] != 1 || emptySlice[0] != 1 {
		t.Error("❌ append should work on both nil and empty slices")
	} else {
		t.Log("✅ Both nil and empty slices work with append")
	}
}

func TestExtractWithoutLeak(t *testing.T) {
	huge := make([]int, 1_000_000)
	for i := range huge {
		huge[i] = i
	}

	small := ExtractWithoutLeak(huge, 100, 105)

	// Check elements
	want := []int{100, 101, 102, 103, 104}
	if !reflect.DeepEqual(small, want) {
		t.Errorf("❌ ExtractWithoutLeak(huge, 100, 105) = %v, want %v", small, want)
		return
	}

	// Check cap equals len (no excess capacity referencing the huge array)
	if cap(small) != len(small) {
		t.Errorf("❌ cap(small) = %d, want %d  ← Hint: allocate with make([]int, to-from) so cap == len", cap(small), len(small))
		return
	}

	// Check independence — modifying small must not affect huge
	small[0] = -1
	if huge[100] == -1 {
		t.Error("❌ Modifying extracted slice changed the original  ← Hint: use copy() to detach from original backing array")
	} else {
		t.Logf("✅ ExtractWithoutLeak produced independent slice %v with cap=%d", want, len(want))
	}
}

func TestObserveGrowth(t *testing.T) {
	caps := ObserveGrowth(10)

	if len(caps) != 10 {
		t.Fatalf("❌ ObserveGrowth(10) returned %d entries, want 10", len(caps))
	}

	// First element should have cap >= 1 (grew from 0)
	if caps[0] < 1 {
		t.Errorf("❌ After first append, cap should be >= 1, got %d", caps[0])
	}

	// Caps should be non-decreasing
	for i := 1; i < len(caps); i++ {
		if caps[i] < caps[i-1] {
			t.Errorf("❌ Capacity decreased at index %d: %d < %d  ← capacity never shrinks", i, caps[i], caps[i-1])
		}
	}

	// Last cap should be >= 10 (must hold 10 elements)
	if caps[len(caps)-1] < 10 {
		t.Errorf("❌ Final cap = %d, must be >= 10", caps[len(caps)-1])
	}

	// Verify growth happened (cap should not equal len for all entries — growslice must have occurred)
	growthOccurred := false
	for i, c := range caps {
		if c > i+1 {
			growthOccurred = true
			break
		}
	}
	if !growthOccurred {
		t.Error("❌ No growth beyond len detected — are you pre-allocating?  ← Hint: start with a nil slice")
	} else {
		t.Logf("✅ ObserveGrowth(10) caps = %v", caps)
	}
}

func TestDetachSlice(t *testing.T) {
	original := []int{10, 20, 30, 40, 50}
	detached := DetachSlice(original)

	// Check elements are the same
	if !reflect.DeepEqual(detached, original) {
		t.Errorf("❌ DetachSlice(%v) = %v, want same elements", original, detached)
		return
	}

	// Check cap == len (the full slice expression must limit capacity)
	if cap(detached) != len(detached) {
		t.Errorf("❌ cap(detached) = %d, want %d  ← Hint: use s[0:len(s):len(s)] to limit cap", cap(detached), len(detached))
		return
	}

	// The key test: appending to detached must NOT corrupt original
	originalCopy := make([]int, len(original))
	copy(originalCopy, original)

	detached = append(detached, 99)

	if !reflect.DeepEqual(original, originalCopy) {
		t.Errorf("❌ Appending to detached slice corrupted original: %v  ← Hint: full slice expression s[low:high:max] limits cap so append triggers growslice", original)
	} else {
		t.Logf("✅ DetachSlice: append to detached did not corrupt original, cap limited to %d", cap(DetachSlice(original)))
	}
}
