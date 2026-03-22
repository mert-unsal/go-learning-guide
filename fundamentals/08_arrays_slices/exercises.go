package arrays_slices

// ============================================================
// EXERCISES — 08 Arrays & Slices
// ============================================================

// Exercise 1:
// ReverseSlice reverses a slice IN-PLACE using the two-pointer technique.
//
// LESSON: Two-pointer pattern — one pointer from each end, swap and move inward.
// This is O(n) time, O(1) space (no extra allocation).
//
//	[1, 2, 3, 4, 5]
//	 L           R   → swap(1,5)
//	    L     R      → swap(2,4)
//	       L=R       → stop
func ReverseSlice(s []int) {
	left, right := 0, len(s)-1
	for left < right {
		s[left], s[right] = s[right], s[left] // Go multiple assignment = elegant swap
		left++
		right--
	}
}

// Exercise 2:
// RemoveDuplicates removes duplicates from a SORTED slice in-place.
//
// LESSON: Write-pointer pattern (also called "slow/fast pointer").
// 'write' is where the next unique value goes. 'read' scans forward.
// Only write when the value changes.
//
//	[1, 1, 2, 3, 3, 4]
//	 w=0, r=0: write=1
//	 w=1, r=1: same as last, skip
//	 w=1, r=2: write=2
//	 ...
func RemoveDuplicates(s []int) []int {
	if len(s) == 0 {
		return s
	}
	write := 1 // s[0] is always unique, start writing at index 1
	for read := 1; read < len(s); read++ {
		if s[read] != s[write-1] { // different from last written value?
			s[write] = s[read]
			write++
		}
	}
	return s[:write] // re-slice to only the unique portion
}

// Exercise 3:
// Make2D allocates a rows×cols matrix of zeros.
//
// LESSON: You CANNOT do `var m [rows][cols]int` — array sizes must be compile-time constants.
// Use make() + a loop. Each row must be individually allocated.
//
// TRAP: `make([][]int, rows)` allocates the outer slice, but each inner slice is nil.
// You must allocate each row separately.
func Make2D(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}
	return matrix
}

// Exercise 4:
// RotateLeft rotates s LEFT by k positions IN-PLACE using the three-reversal trick.
//
// LESSON: Three-reversal trick — no extra space needed.
// [1,2,3,4,5], k=2 → want [3,4,5,1,2]
//
//	Step 1: reverse whole:   [5,4,3,2,1]
//	Step 2: reverse first n-k=3: [3,4,5,2,1]
//	Step 3: reverse last k=2:    [3,4,5,1,2] ✅
func RotateLeft(s []int, k int) {
	n := len(s)
	if n == 0 {
		return
	}
	k = k % n // handle k >= n (rotating by full length = no-op)
	if k == 0 {
		return
	}
	reverse(s, 0, n-1)
	reverseRange(s, 0, n-k-1)
	reverseRange(s, n-k, n-1)
}

func reverseRange(s []int, l, r int) {
	for l < r {
		s[l], s[r] = s[r], s[l]
		l++
		r--
	}
}

// Exercise 5:
// Filter returns a new slice containing only elements where fn(element) is true.
//
// LESSON: Higher-order functions. Go doesn't have built-in map/filter/reduce,
// but you can write them easily. Note we build a new slice — we don't modify original.
func Filter(s []int, fn func(int) bool) []int {
	result := []int{} // empty non-nil slice (not nil, which matters for JSON etc)
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Exercise 6:
// MergeSorted merges two sorted slices into one sorted slice.
//
// LESSON: Classic two-pointer merge — the foundation of merge sort.
// Compare front elements, take the smaller, advance that pointer.
// Then drain whichever slice still has elements.
func MergeSorted(a, b []int) []int {
	result := make([]int, 0, len(a)+len(b)) // pre-allocate exact size
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}
	// Drain remaining elements (at most one of these loops runs)
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)
	return result
}
