// 01-arrays demonstrates Go arrays: fixed-size, value-semantic collections.
//
// Run:  go run .
//
// ============================================================
// ARRAYS
// ============================================================
// Arrays in Go have a FIXED size. The size is part of the type.
// [3]int and [5]int are DIFFERENT types — you cannot assign one to the other.
// Arrays are VALUES — copying an array copies all elements (no shared backing store).
// Under the hood, the compiler allocates contiguous memory of size * element_size bytes.
package main

import "fmt"

func main() {
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
