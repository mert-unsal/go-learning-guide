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

import (
	"fmt"
	"unsafe"
)

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
	fmt.Printf("%s%s  Arrays — Fixed-Size Value Types         %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// Declaration with size
	fmt.Printf("%s▸ Zero-Value Declaration%s\n", cyan+bold, reset)
	var a [5]int // zero value: [0 0 0 0 0]
	fmt.Printf("  var a [5]int → %v\n", a)
	fmt.Printf("  %s✔ All elements default to zero value of the element type%s\n", green, reset)
	fmt.Printf("  sizeof([5]int) = %s%d bytes%s (5 × 8 bytes on 64-bit)\n\n", magenta, unsafe.Sizeof(a), reset)

	// Array literal
	fmt.Printf("%s▸ Array Literal%s\n", cyan+bold, reset)
	b := [3]string{"Go", "Python", "Rust"}
	fmt.Printf("  b := [3]string{\"Go\", \"Python\", \"Rust\"}\n")
	fmt.Printf("  b[0] = %s%q%s\n", magenta, b[0], reset)
	fmt.Printf("  %s✔ Size is part of the type — [3]string and [4]string are DIFFERENT types%s\n\n", green, reset)

	// Let compiler count with ...
	fmt.Printf("%s▸ Compiler-Counted Array ([...])%s\n", cyan+bold, reset)
	c := [...]int{1, 2, 3, 4, 5} // size inferred as 5
	fmt.Printf("  c := [...]int{1,2,3,4,5}\n")
	fmt.Printf("  len(c) = %s%d%s — compiler infers the size from the literal\n", magenta, len(c), reset)
	fmt.Printf("  %s✔ [...] is syntactic sugar — the resulting type is [5]int, fully fixed at compile time%s\n\n", green, reset)

	// Arrays are values (copied, not referenced)
	fmt.Printf("%s▸ Arrays Are Values (Copy Semantics)%s\n", cyan+bold, reset)
	d := c // d is a COPY of c
	d[0] = 99
	fmt.Printf("  d := c; d[0] = 99\n")
	fmt.Printf("  c[0] = %s%d%s, d[0] = %s%d%s\n", magenta, c[0], reset, magenta, d[0], reset)
	fmt.Printf("  %s✔ Assignment copies ALL elements — no shared backing store%s\n", green, reset)
	fmt.Printf("  %s⚠ Large arrays (e.g. [1_000_000]int) are expensive to copy — pass by pointer or use slices%s\n\n", yellow, reset)

	// 2D array
	fmt.Printf("%s▸ 2D Array (Fixed Grid)%s\n", cyan+bold, reset)
	var grid [3][3]int
	grid[1][1] = 5
	fmt.Printf("  var grid [3][3]int; grid[1][1] = 5\n")
	for i, row := range grid {
		fmt.Printf("  row %d: %v\n", i, row)
	}
	fmt.Printf("  %s✔ Contiguous memory: all 9 ints laid out sequentially (good cache locality)%s\n", green, reset)
	fmt.Printf("  sizeof([3][3]int) = %s%d bytes%s\n\n", magenta, unsafe.Sizeof(grid), reset)

	fmt.Printf("%s%s── Key Takeaways ──%s\n", bold, blue, reset)
	fmt.Printf("  %s✔ Arrays are value types — size is part of the type%s\n", green, reset)
	fmt.Printf("  %s✔ Contiguous memory layout gives excellent cache performance%s\n", green, reset)
	fmt.Printf("  %s⚠ Rarely used directly in Go — slices are the idiomatic choice%s\n", yellow, reset)
}
