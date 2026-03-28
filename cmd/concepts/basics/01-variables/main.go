// Variables in Go — demonstrates declaration, initialization, and assignment.
//
// Every Go file starts with a package declaration.
// "main" package is special: it defines an executable program.
// Other packages are libraries.
//
// Imports bring in other packages. Unused imports are a COMPILE ERROR in Go.
package main

import "fmt"

func main() {
	// --- var keyword (explicit) ---
	var a int     // zero value: 0
	var b string  // zero value: ""
	var c bool    // zero value: false
	var d float64 // zero value: 0.0

	fmt.Println("Zero values:", a, b, c, d)

	// --- var with initializer ---
	var x int = 42
	var name string = "Gopher"
	fmt.Println(x, name)

	// --- Short variable declaration := (most common inside functions) ---
	// Type is INFERRED automatically
	age := 25
	pi := 3.14159
	isGoFun := true
	fmt.Println(age, pi, isGoFun)

	// --- Multiple assignment ---
	x1, y1 := 10, 20
	fmt.Println(x1, y1)

	// --- Swap values (Go's elegant way) ---
	x1, y1 = y1, x1
	fmt.Println("Swapped:", x1, y1)

	// --- Blank identifier _ (discard a value) ---
	_, second := 100, 200
	fmt.Println("Only second:", second)
}
