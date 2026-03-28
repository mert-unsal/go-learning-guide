// 03-append demonstrates the append built-in and capacity growth.
//
// Run:  go run .
//
// ============================================================
// APPEND
// ============================================================
// append returns a NEW slice header. Always reassign: s = append(s, v).
//
// Growth strategy (runtime/slice.go → growslice):
//   cap < 256  → double capacity (2x)
//   cap >= 256 → grow by ~1.25x + 192
// When cap is exceeded, Go allocates a NEW backing array —
// the old one is left for GC. This means previously-shared
// slices will NO LONGER see each other's mutations.
package main

import "fmt"

func main() {
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
