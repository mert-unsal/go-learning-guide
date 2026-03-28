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

func main() {
	x := 42
	p := &x // p is of type *int; it POINTS TO x

	fmt.Println("x =", x)   // 42
	fmt.Println("p =", p)   // memory address like 0xc0000b4010
	fmt.Println("*p =", *p) // 42 (dereference: value at the address)

	*p = 100       // modify the value through the pointer
	fmt.Println(x) // 100 — x was changed!

	// new() allocates and returns a pointer to a zeroed value
	q := new(int) // *int, zero value is 0
	*q = 55
	fmt.Println(*q) // 55

	// nil pointer: a pointer with no address
	var r *int     // nil
	fmt.Println(r) // <nil>
	// *r = 5        // PANIC: nil pointer dereference — never do this!

	// Safe nil check before dereferencing
	if r != nil {
		fmt.Println(*r)
	} else {
		fmt.Println("r is nil, can't dereference")
	}
}
