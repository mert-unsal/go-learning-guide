// Standalone demo: Pass by Reference with Pointers
//
// Go is ALWAYS pass-by-value. When you pass a pointer, the pointer itself
// is copied, but both copies point to the same memory — giving you
// reference semantics. This is identical to how Java passes object refs.
//
// Slices, maps, and channels are already reference types (they contain an
// internal pointer to the backing data), so you rarely need explicit
// pointers for them.
//
// Run:  go run .
package main

import "fmt"

// WITHOUT pointer: modifies a COPY, original unchanged
func incrementValue(n int) {
	n++ // modifies local copy only
}

// WITH pointer: modifies the ORIGINAL
func incrementPointer(n *int) {
	*n++ // dereferences and modifies the original
}

func main() {
	x := 10
	incrementValue(x)
	fmt.Println("After incrementValue:", x) // still 10

	incrementPointer(&x)
	fmt.Println("After incrementPointer:", x) // 11

	// Slices, maps, and channels are already reference types —
	// you don't need pointers for them in most cases.
}
