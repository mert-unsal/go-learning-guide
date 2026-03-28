// 02-slices demonstrates Go slices — the primary dynamic collection type.
//
// Run:  go run .
//
// ============================================================
// SLICES — The primary collection type in Go
// ============================================================
// A slice is a DYNAMIC VIEW into an array.
// Internally: { pointer to array, length, capacity }
//   — 3-word struct (24 bytes on 64-bit): runtime represents this as
//     reflect.SliceHeader { Data uintptr; Len int; Cap int }.
// Slices are reference types — they share the underlying array.
//
// Key distinction:
//   nil slice  → var s []int        → s == nil is true,  len=0, cap=0
//   empty slice → s := []int{}      → s == nil is false, len=0, cap=0
//   Both work fine with append, range, and len.
package main

import "fmt"

func main() {
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
