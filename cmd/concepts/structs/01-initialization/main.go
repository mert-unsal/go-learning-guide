// Standalone demo: Struct Initialization Patterns
//
// A struct groups related data. It's Go's primary way to create custom
// data types (no classes in Go!).
//
// Under the hood: struct fields are laid out contiguously in memory with
// padding for alignment. Field order matters — grouping fields by size
// reduces padding and shrinks the struct. Inspect with unsafe.Sizeof/Alignof.
//
// Run:  go run .
package main

import "fmt"

type Point struct {
	X float64
	Y float64
}

func main() {
	// Named fields (PREFERRED — order independent, self-documenting)
	p1 := Point{X: 1.0, Y: 2.0}

	// Positional (avoid — fragile if fields are added/reordered)
	p2 := Point{3.0, 4.0}

	// Zero value struct
	var p3 Point // X=0, Y=0

	// Pointer to struct
	p4 := &Point{X: 5.0, Y: 6.0}

	fmt.Println(p1, p2, p3, *p4)

	// Accessing fields
	fmt.Println("X:", p1.X, "Y:", p1.Y)

	// Through a pointer — Go auto-dereferences
	fmt.Println("p4.X:", p4.X) // same as (*p4).X
}
