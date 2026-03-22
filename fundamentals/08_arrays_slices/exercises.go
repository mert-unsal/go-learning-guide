package arrays_slices

// ============================================================
// EXERCISES — 08 Arrays & Slices
// ============================================================
//
// Exercises 1-6:  Algorithm patterns (reverse, dedup, 2D, rotate, filter, merge)
// Exercises 7-12: Slice internals (backing array, copy, nil vs empty, leaks, append growth, full slice expression)

// ============================================================
// PART A — Algorithm Patterns
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
	// TODO: implement two-pointer reverse
	panic("not implemented")
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
	// TODO: implement write-pointer deduplication on sorted input
	panic("not implemented")
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
	// TODO: allocate outer slice, then each inner row
	panic("not implemented")
}

// Exercise 4:
// RotateLeft rotates s LEFT by k positions IN-PLACE using the three-reversal trick.
//
// LESSON: Three-reversal trick — no extra space needed.
// [1,2,3,4,5], k=2 → want [3,4,5,1,2]
//
//	Step 1: reverse whole:       [5,4,3,2,1]
//	Step 2: reverse first n-k=3: [3,4,5,2,1]
//	Step 3: reverse last k=2:    [3,4,5,1,2] ✅
//
// Hint: write a helper func reverseRange(s []int, l, r int)
func RotateLeft(s []int, k int) {
	// TODO: implement three-reversal trick
	panic("not implemented")
}

// Exercise 5:
// Filter returns a new slice containing only elements where fn(element) is true.
//
// LESSON: Higher-order functions. Go doesn't have built-in map/filter/reduce,
// but you can write them easily. Note we build a new slice — we don't modify original.
func Filter(s []int, fn func(int) bool) []int {
	// TODO: build and return a new slice with matching elements
	panic("not implemented")
}

// Exercise 6:
// MergeSorted merges two sorted slices into one sorted slice.
//
// LESSON: Classic two-pointer merge — the foundation of merge sort.
// Compare front elements, take the smaller, advance that pointer.
// Then drain whichever slice still has elements.
func MergeSorted(a, b []int) []int {
	// TODO: implement two-pointer merge
	// Hint: pre-allocate with make([]int, 0, len(a)+len(b))
	panic("not implemented")
}

// ============================================================
// PART B — Slice Internals
// ============================================================

// Exercise 7:
// SafeDelete deletes the element at index i from s (order-preserving) and returns
// the shortened slice WITHOUT mutating the original backing array that the caller holds.
//
// LESSON: The standard delete trick `append(s[:i], s[i+1:]...)` mutates the original
// backing array. If another variable shares that backing array, it sees corrupted data.
// You must copy first to get an independent backing array, then delete from the copy.
//
// Example:
//
//	original := []int{10, 20, 30, 40, 50}
//	result := SafeDelete(original, 2)
//	// result   = [10, 20, 40, 50]
//	// original = [10, 20, 30, 40, 50]  ← must be unchanged!
func SafeDelete(s []int, i int) []int {
	// TODO: create an independent copy, then delete index i from the copy
	panic("not implemented")
}

// Exercise 8:
// CopySlice creates a completely independent copy of src using the copy() built-in.
// The returned slice must have the same length and elements as src, but a different
// backing array.
//
// LESSON: copy() copies min(len(dst), len(src)) elements. You must allocate dst
// with the right length BEFORE calling copy. copy() does NOT allocate for you.
//
// Requirements:
//   - Must use the copy() built-in (not append)
//   - Returned slice must have same len as src
//   - Modifying the returned slice must NOT affect src
func CopySlice(src []int) []int {
	// TODO: allocate a new slice with make(), then use copy()
	panic("not implemented")
}

// Exercise 9:
// NilVsEmpty returns two slices: one nil and one empty (non-nil).
// This exercise tests your understanding of the difference.
//
// LESSON: A nil slice and an empty slice behave identically for len(), cap(),
// append(), and range — but they differ in nil comparison and JSON marshaling.
//
//	var s []int       → nil slice:  s == nil is true,  json: "null"
//	s := []int{}      → empty slice: s == nil is false, json: "[]"
//	s := make([]int,0)→ empty slice: s == nil is false, json: "[]"
//
// Returns: (nilSlice, emptySlice)
func NilVsEmpty() ([]int, []int) {
	// TODO: return a nil slice and an empty (non-nil) slice
	panic("not implemented")
}

// Exercise 10:
// ExtractWithoutLeak takes a large slice and an index range [from, to),
// and returns a NEW slice containing only those elements.
// The returned slice must NOT hold a reference to the original's backing array.
//
// LESSON: s[from:to] creates a sub-slice that shares the original backing array.
// If the original is large (e.g., 1 million elements) and the sub-slice is small
// (e.g., 3 elements), the entire original array stays in memory because the sub-slice
// holds a pointer into it. This is a memory leak.
//
// Fix: copy the elements into a new, independent slice.
//
// Example:
//
//	huge := make([]int, 1_000_000)
//	small := ExtractWithoutLeak(huge, 0, 3)
//	// small has len=3, cap=3, independent backing array
//	// huge can now be garbage collected
func ExtractWithoutLeak(s []int, from, to int) []int {
	// TODO: extract s[from:to] into a new slice with its own backing array
	panic("not implemented")
}

// Exercise 11:
// ObserveGrowth appends n elements (values 0 to n-1) to an initially empty slice
// and returns a slice of every capacity value observed after each append.
//
// LESSON: append uses growslice when len==cap. The growth strategy is:
//   - 2x when cap < 256
//   - ~1.25x when cap >= 256
//   - Then rounded up to memory allocator size classes
//
// Understanding this helps you decide when to pre-allocate with make([]T, 0, n).
//
// Example:
//
//	ObserveGrowth(5) might return [1, 2, 4, 4, 8]
//	// append 0: cap goes 0→1
//	// append 1: cap goes 1→2
//	// append 2: cap goes 2→4
//	// append 3: cap stays 4
//	// append 4: cap goes 4→8
func ObserveGrowth(n int) []int {
	// TODO: start with a nil slice, append values 0..n-1,
	// record cap(s) after each append into a result slice
	panic("not implemented")
}

// Exercise 12:
// DetachSlice takes a slice and returns a new slice with the same elements
// but limited capacity, so that appending to the returned slice cannot
// accidentally overwrite elements in the original backing array.
//
// LESSON: The full slice expression s[low:high:max] sets cap = max - low.
// This is the Go idiom for "detaching" a sub-slice so future appends
// trigger growslice instead of overwriting shared backing data.
//
// You MUST use the full slice expression (three-index slice) in your solution.
//
// Example:
//
//	original := []int{10, 20, 30, 40, 50}
//	detached := DetachSlice(original)
//	detached = append(detached, 99)
//	// original is still [10, 20, 30, 40, 50] — not corrupted
func DetachSlice(s []int) []int {
	// TODO: return s with cap limited to len using full slice expression
	// Hint: s[0:len(s):len(s)]
	panic("not implemented")
}
