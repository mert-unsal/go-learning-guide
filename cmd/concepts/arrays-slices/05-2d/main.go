// 05-2d demonstrates 2D slices and essential slice tricks for coding exams.
//
// Run:  go run .
//
// ============================================================
// 2D SLICES
// ============================================================
// Go doesn't have true multi-dimensional slices. A 2D slice is
// a slice of slices — each inner slice is independently allocated.
// This means rows can have different lengths (jagged arrays).
//
// ============================================================
// ESSENTIAL SLICE TRICKS (Critical for coding exams!)
// ============================================================
// These in-place operations avoid allocations and are O(n) or better.
// Know them cold for interviews.
package main

import "fmt"

func main() {
	// --- 2D Slices ---
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

	fmt.Println()

	// --- Slice Tricks ---
	sliceTricks()
}

func sliceTricks() {
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
