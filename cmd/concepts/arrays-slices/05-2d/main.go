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

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  2D Slices & Essential Slice Tricks       %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- 2D Slices ---
	fmt.Printf("%s▸ 2D Slice — Slice of Slices%s\n", cyan+bold, reset)
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
	fmt.Printf("  make([][]int, %d) — each row allocated independently:\n", rows)
	for i, row := range matrix {
		fmt.Printf("    row %d: %s%v%s\n", i, magenta, row, reset)
	}
	fmt.Printf("  %s✔ Each row is a separate heap allocation — rows can have different lengths (jagged)%s\n", green, reset)
	fmt.Printf("  %s⚠ Not contiguous memory like [3][3]int — worse cache locality for large matrices%s\n\n", yellow, reset)

	// --- Slice Tricks ---
	fmt.Printf("%s%s── Essential Slice Tricks (Interview Critical) ──%s\n\n", bold, blue, reset)
	sliceTricks()
}

func sliceTricks() {
	// --- Delete element at index i (maintain order) ---
	fmt.Printf("%s▸ Delete At Index (Order-Preserving) — O(n)%s\n", cyan+bold, reset)
	s := []int{1, 2, 3, 4, 5}
	i := 2 // delete index 2 (value 3)
	fmt.Printf("  before: %v — delete index %d (value %d)\n", s, i, s[i])
	s = append(s[:i], s[i+1:]...)
	fmt.Printf("  s = append(s[:i], s[i+1:]...) → %s%v%s\n", magenta, s, reset)
	fmt.Printf("  %s✔ Shifts elements left — O(n) but preserves order%s\n\n", green, reset)

	// --- Delete element at index i (don't maintain order) ---
	fmt.Printf("%s▸ Delete At Index (Swap-With-Last) — O(1)%s\n", cyan+bold, reset)
	s = []int{1, 2, 3, 4, 5}
	fmt.Printf("  before: %v — delete index %d\n", s, i)
	s[i] = s[len(s)-1]
	s = s[:len(s)-1]
	fmt.Printf("  s[i] = s[len-1]; s = s[:len-1] → %s%v%s\n", magenta, s, reset)
	fmt.Printf("  %s✔ O(1) delete — swap with last element, shrink. Use when order doesn't matter%s\n\n", green, reset)

	// --- Insert at index i ---
	fmt.Printf("%s▸ Insert At Index — O(n)%s\n", cyan+bold, reset)
	s = []int{1, 2, 4, 5}
	i = 2
	fmt.Printf("  before: %v — insert 3 at index %d\n", s, i)
	s = append(s[:i+1], s[i:]...)
	s[i] = 3
	fmt.Printf("  s = append(s[:i+1], s[i:]...); s[i] = 3 → %s%v%s\n", magenta, s, reset)
	fmt.Printf("  %s✔ Shifts elements right to make room — O(n) due to copy%s\n\n", green, reset)

	// --- Reverse a slice ---
	fmt.Printf("%s▸ Reverse In-Place — O(n)%s\n", cyan+bold, reset)
	s = []int{1, 2, 3, 4, 5}
	fmt.Printf("  before: %v\n", s)
	for l, r := 0, len(s)-1; l < r; l, r = l+1, r-1 {
		s[l], s[r] = s[r], s[l]
	}
	fmt.Printf("  two-pointer swap → %s%v%s\n", magenta, s, reset)
	fmt.Printf("  %s✔ Classic two-pointer pattern: O(n) time, O(1) space, zero allocations%s\n\n", green, reset)

	// --- Contains (no built-in in Go < 1.21) ---
	fmt.Printf("%s▸ Contains (Linear Scan) — O(n)%s\n", cyan+bold, reset)
	target := 3
	found := false
	for _, v := range s {
		if v == target {
			found = true
			break
		}
	}
	fmt.Printf("  s = %v — contains %d? %s%v%s\n", s, target, magenta, found, reset)
	fmt.Printf("  %s✔ Go 1.21+ has slices.Contains(); before that, manual loop is idiomatic%s\n\n", green, reset)

	// --- Remove duplicates ---
	fmt.Printf("%s▸ Remove Duplicates (Map-Based) — O(n)%s\n", cyan+bold, reset)
	withDups := []int{1, 2, 2, 3, 3, 3, 4}
	fmt.Printf("  input: %v\n", withDups)
	seen := make(map[int]bool)
	unique := make([]int, 0)
	for _, v := range withDups {
		if !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}
	fmt.Printf("  unique: %s%v%s\n", magenta, unique, reset)
	fmt.Printf("  %s✔ map[int]bool tracks seen values — O(n) time, O(n) space%s\n", green, reset)
	fmt.Printf("  %s⚠ For sorted input, can deduplicate in-place with O(1) space (two-pointer)%s\n\n", yellow, reset)

	// --- Filter (keep elements matching condition) ---
	fmt.Printf("%s▸ Filter In-Place (Zero-Alloc) — O(n)%s\n", cyan+bold, reset)
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("  input: %v — keep evens\n", nums)
	evens := nums[:0] // reuse backing array
	for _, v := range nums {
		if v%2 == 0 {
			evens = append(evens, v)
		}
	}
	fmt.Printf("  evens := nums[:0] + append → %s%v%s\n", magenta, evens, reset)
	fmt.Printf("  %s✔ nums[:0] reuses the backing array — zero heap allocations%s\n", green, reset)
	fmt.Printf("  %s⚠ This mutates the original slice's backing array — make a copy first if you need both%s\n\n", yellow, reset)

	fmt.Printf("%s%s── Key Takeaways ──%s\n", bold, blue, reset)
	fmt.Printf("  %s✔ Know these tricks cold — they appear in 80%%+ of array/slice interview problems%s\n", green, reset)
	fmt.Printf("  %s✔ Prefer O(1) swap-delete when order doesn't matter%s\n", green, reset)
	fmt.Printf("  %s✔ nums[:0] filter trick avoids allocations — critical for hot-path code%s\n", green, reset)
}
