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
