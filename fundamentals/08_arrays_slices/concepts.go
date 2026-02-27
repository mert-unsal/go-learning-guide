// Package arrays_slices covers arrays vs slices, slice internals,
// append, copy, 2D slices, and essential slice tricks for coding exams.
package arrays_slices

import "fmt"

// ============================================================
// 1. ARRAYS
// ============================================================
// Arrays in Go have a FIXED size. The size is part of the type.
// [3]int and [5]int are DIFFERENT types.
// Arrays are VALUES — copying an array copies all elements.

func DemonstrateArrays() {
	// Declaration with size
	var a [5]int // zero value: [0 0 0 0 0]
	fmt.Println(a)

	// Array literal
	b := [3]string{"Go", "Python", "Rust"}
	fmt.Println(b[0]) // "Go"

	// Let compiler count with ...
	c := [...]int{1, 2, 3, 4, 5} // size inferred as 5
	fmt.Println(len(c))          // 5

	// Arrays are values (copied, not referenced)
	d := c // d is a COPY of c
	d[0] = 99
	fmt.Println(c[0], d[0]) // 1, 99 — c is unchanged

	// 2D array
	var grid [3][3]int
	grid[1][1] = 5
	fmt.Println(grid)
}

// ============================================================
// 2. SLICES — The primary collection type in Go
// ============================================================
// A slice is a DYNAMIC VIEW into an array.
// Internally: { pointer to array, length, capacity }
// Slices are reference types — they share underlying array.

func DemonstrateSlices() {
	// From literal
	s := []int{1, 2, 3, 4, 5}
	fmt.Println(s, "len:", len(s), "cap:", cap(s))

	// make([]T, length, capacity) — allocates with specific size
	s2 := make([]int, 3)     // [0 0 0], len=3, cap=3
	s3 := make([]int, 3, 10) // [0 0 0], len=3, cap=10
	fmt.Println(s2, s3)

	// nil slice vs empty slice
	var nilSlice []int             // nil, len=0, cap=0
	emptySlice := []int{}          // not nil, len=0, cap=0
	fmt.Println(nilSlice == nil)   // true
	fmt.Println(emptySlice == nil) // false
	// Both work fine with append, range, len

	// Slicing — s[low:high] — includes low, excludes high
	fmt.Println(s[1:3]) // [2 3]
	fmt.Println(s[:3])  // [1 2 3]
	fmt.Println(s[2:])  // [3 4 5]
	fmt.Println(s[:])   // [1 2 3 4 5]

	// IMPORTANT: slices share underlying array!
	a := []int{1, 2, 3, 4, 5}
	b := a[1:3] // b = [2, 3], shares array with a
	b[0] = 99
	fmt.Println(a) // [1 99 3 4 5] — a was modified through b!
}

// ============================================================
// 3. APPEND
// ============================================================

func DemonstrateAppend() {
	s := []int{1, 2, 3}

	// Append single element
	s = append(s, 4)

	// Append multiple elements
	s = append(s, 5, 6, 7)

	// Append another slice with ...
	extra := []int{8, 9, 10}
	s = append(s, extra...)

	fmt.Println(s) // [1 2 3 4 5 6 7 8 9 10]

	// IMPORTANT: when cap is exceeded, Go allocates a NEW array
	// The original backing array is unchanged
	a := make([]int, 3, 3)
	b := a           // b shares a's backing array
	a = append(a, 4) // cap exceeded! a gets a new backing array
	a[0] = 99
	fmt.Println(a[0], b[0]) // 99, 0 — they no longer share storage
}

// ============================================================
// 4. COPY
// ============================================================

func DemonstrateCopy() {
	src := []int{1, 2, 3, 4, 5}

	// copy(dst, src) — copies min(len(dst), len(src)) elements
	dst := make([]int, len(src))
	n := copy(dst, src)
	fmt.Println("Copied", n, "elements:", dst)

	// Now dst is independent from src
	dst[0] = 99
	fmt.Println("src:", src) // unchanged
	fmt.Println("dst:", dst)

	// Partial copy
	partial := make([]int, 3)
	copy(partial, src)
	fmt.Println("Partial:", partial) // [1 2 3]
}

// ============================================================
// 5. ESSENTIAL SLICE TRICKS (Critical for coding exams!)
// ============================================================

func SliceTricks() {
	// --- Delete element at index i (maintain order) ---
	s := []int{1, 2, 3, 4, 5}
	i := 2 // delete index 2 (value 3)
	s = append(s[:i], s[i+1:]...)
	fmt.Println("After delete:", s) // [1 2 4 5]

	// --- Delete element at index i (don't maintain order) ---
	s = []int{1, 2, 3, 4, 5}
	s[i] = s[len(s)-1]
	s = s[:len(s)-1]
	fmt.Println("After fast delete:", s)

	// --- Insert at index i ---
	s = []int{1, 2, 4, 5}
	i = 2
	s = append(s[:i+1], s[i:]...)
	s[i] = 3
	fmt.Println("After insert:", s) // [1 2 3 4 5]

	// --- Reverse a slice ---
	s = []int{1, 2, 3, 4, 5}
	for l, r := 0, len(s)-1; l < r; l, r = l+1, r-1 {
		s[l], s[r] = s[r], s[l]
	}
	fmt.Println("Reversed:", s) // [5 4 3 2 1]

	// --- Contains (no built-in in Go < 1.21) ---
	target := 3
	found := false
	for _, v := range s {
		if v == target {
			found = true
			break
		}
	}
	fmt.Printf("Contains %d: %v\n", target, found)

	// --- Remove duplicates ---
	withDups := []int{1, 2, 2, 3, 3, 3, 4}
	seen := make(map[int]bool)
	unique := make([]int, 0)
	for _, v := range withDups {
		if !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}
	fmt.Println("Unique:", unique) // [1 2 3 4]

	// --- Filter (keep elements matching condition) ---
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := nums[:0] // reuse backing array
	for _, v := range nums {
		if v%2 == 0 {
			evens = append(evens, v)
		}
	}
	fmt.Println("Evens:", evens)
}

// ============================================================
// 6. 2D SLICES
// ============================================================

func Demonstrate2D() {
	// Create a 3x3 matrix
	rows, cols := 3, 3
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}

	// Fill with values
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			matrix[i][j] = i*cols + j
		}
	}

	// Print matrix
	for _, row := range matrix {
		fmt.Println(row)
	}
}

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Arrays ===")
	DemonstrateArrays()
	fmt.Println("\n=== Slices ===")
	DemonstrateSlices()
	fmt.Println("\n=== Append ===")
	DemonstrateAppend()
	fmt.Println("\n=== Copy ===")
	DemonstrateCopy()
	fmt.Println("\n=== Slice Tricks ===")
	SliceTricks()
	fmt.Println("\n=== 2D Slices ===")
	Demonstrate2D()
}
