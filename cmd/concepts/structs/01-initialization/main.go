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

type Point struct {
	X float64
	Y float64
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Structs: Initialization Patterns       %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Named Fields ---
	fmt.Printf("%s▸ Named Fields (preferred)%s\n", cyan+bold, reset)

	// Named fields (PREFERRED — order independent, self-documenting)
	p1 := Point{X: 1.0, Y: 2.0}
	fmt.Printf("  p1 := Point{X: 1.0, Y: 2.0} → %s%+v%s\n", magenta, p1, reset)
	fmt.Printf("  %s✔ Order-independent, self-documenting — won't break if fields are reordered%s\n\n", green, reset)

	// --- Positional ---
	fmt.Printf("%s▸ Positional Init (avoid)%s\n", cyan+bold, reset)

	// Positional (avoid — fragile if fields are added/reordered)
	p2 := Point{3.0, 4.0}
	fmt.Printf("  p2 := Point{3.0, 4.0} → %s%+v%s\n", magenta, p2, reset)
	fmt.Printf("  %s⚠ Positional init is fragile — breaks silently if fields are added/reordered%s\n\n", yellow, reset)

	// --- Zero Value ---
	fmt.Printf("%s▸ Zero Value Struct%s\n", cyan+bold, reset)

	// Zero value struct
	var p3 Point // X=0, Y=0
	fmt.Printf("  var p3 Point → %s%+v%s\n", magenta, p3, reset)
	fmt.Printf("  %s✔ All fields set to their zero values (0.0 for float64)%s\n", green, reset)
	fmt.Printf("  %s✔ \"Make the zero value useful\" — Go proverb%s\n\n", green, reset)

	// --- Pointer to Struct ---
	fmt.Printf("%s▸ Pointer to Struct (&Type{})%s\n", cyan+bold, reset)

	// Pointer to struct
	p4 := &Point{X: 5.0, Y: 6.0}
	fmt.Printf("  p4 := &Point{X: 5.0, Y: 6.0}\n")
	fmt.Printf("    addr: %s%p%s, value: %s%+v%s\n", magenta, p4, reset, magenta, *p4, reset)
	fmt.Printf("  %s✔ &Type{} creates struct and returns pointer in one step%s\n\n", green, reset)

	// --- Field Access ---
	fmt.Printf("%s▸ Field Access%s\n", cyan+bold, reset)

	// Accessing fields
	fmt.Printf("  p1.X = %s%v%s, p1.Y = %s%v%s\n", magenta, p1.X, reset, magenta, p1.Y, reset)

	// Through a pointer — Go auto-dereferences
	fmt.Printf("  p4.X = %s%v%s (auto-deref: p4.X is sugar for (*p4).X)\n", magenta, p4.X, reset)
	fmt.Printf("  %s✔ Go auto-dereferences struct pointers — no explicit (*p4).X needed%s\n\n", green, reset)

	// --- Memory Layout ---
	fmt.Printf("%s▸ Memory Layout Awareness%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Fields are contiguous in memory with padding for alignment%s\n", green, reset)
	fmt.Printf("  %s✔ Field order matters: grouping by size reduces padding, shrinks struct%s\n", green, reset)
	fmt.Printf("  %s✔ Inspect with unsafe.Sizeof(), unsafe.Alignof(), unsafe.Offsetof()%s\n", green, reset)
}
