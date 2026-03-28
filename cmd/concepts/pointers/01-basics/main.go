// Standalone demo: Pointer Basics
//
// &x  → "address of x" → returns a *T (pointer to T)
// *p  → "dereference p" → returns the value at the address
// new(T) → allocates zeroed T, returns *T
//
// Under the hood: Go pointers are raw memory addresses (like C) but with
// no arithmetic. The GC tracks every live pointer to keep referenced objects
// alive. The escape analysis (`go build -gcflags='-m'`) decides whether the
// pointed-to value lives on the stack or must be promoted to the heap.
//
// Run:  go run .
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
	fmt.Printf("%s%s  Pointers: Basics (&, *, new)           %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Address & Dereference ---
	fmt.Printf("%s▸ Address-of (&) and Dereference (*)%s\n", cyan+bold, reset)

	x := 42
	p := &x // p is of type *int; it POINTS TO x

	fmt.Printf("  x  = %s%d%s\n", magenta, x, reset)
	fmt.Printf("  &x = %s%p%s  (address of x — a raw memory location)\n", magenta, &x, reset)
	fmt.Printf("  p  = %s%p%s  (p holds the same address)\n", magenta, p, reset)
	fmt.Printf("  *p = %s%d%s  (dereference: follow the pointer to read the value)\n", magenta, *p, reset)
	fmt.Printf("  %s✔ &x and p hold the same address — they point to the same memory%s\n\n", green, reset)

	// --- Modify via Pointer ---
	fmt.Printf("%s▸ Modify Through a Pointer%s\n", cyan+bold, reset)
	fmt.Printf("  Before: x = %s%d%s (at %s%p%s)\n", magenta, x, reset, magenta, &x, reset)

	*p = 100 // modify the value through the pointer
	fmt.Printf("  After *p = 100: x = %s%d%s\n", magenta, x, reset)
	fmt.Printf("  %s✔ *p = 100 writes to the address p holds — x and *p are the same memory%s\n\n", green, reset)

	// --- new() ---
	fmt.Printf("%s▸ new() — Allocate & Return Pointer%s\n", cyan+bold, reset)

	// new() allocates and returns a pointer to a zeroed value
	q := new(int) // *int, zero value is 0
	fmt.Printf("  q  = new(int) → %s%p%s (points to heap-allocated int)\n", magenta, q, reset)
	fmt.Printf("  *q = %s%d%s (zero value — new() always zeroes memory)\n", magenta, *q, reset)
	*q = 55
	fmt.Printf("  After *q = 55: *q = %s%d%s\n", magenta, *q, reset)
	fmt.Printf("  %s✔ new(T) allocates zeroed memory for T and returns *T%s\n", green, reset)
	fmt.Printf("  %s✔ Escape analysis decides if this lands on stack or heap%s\n\n", green, reset)

	// --- nil Pointer ---
	fmt.Printf("%s▸ nil Pointer%s\n", cyan+bold, reset)

	// nil pointer: a pointer with no address
	var r *int // nil
	fmt.Printf("  var r *int → r = %s%v%s\n", magenta, r, reset)
	fmt.Printf("  %s⚠ Dereferencing nil (*r) causes PANIC: nil pointer dereference%s\n", yellow, reset)

	// Safe nil check before dereferencing
	if r != nil {
		fmt.Printf("  *r = %d\n", *r)
	} else {
		fmt.Printf("  %s✔ Always check r != nil before dereferencing%s\n", green, reset)
	}
	// *r = 5        // PANIC: nil pointer dereference — never do this!
}
