// Standalone demo: Anonymous Structs
//
// Useful for one-off data grouping, test cases, etc.
// Anonymous structs are defined and used in place — no named type needed.
//
// The table-driven test pattern below is the idiomatic Go way to write
// tests. Each test case is a struct literal in a slice, iterated with
// range. This keeps tests concise, easy to extend, and gives you
// named fields that serve as documentation.
//
// Run:  go run .
package main

import "fmt"

func main() {
	// Define and initialize in one step
	point := struct {
		X, Y int
	}{X: 10, Y: 20}
	fmt.Println(point)

	// Very common in tests (table-driven test pattern):
	tests := []struct {
		input    int
		expected int
	}{
		{1, 1},
		{2, 4},
		{3, 9},
	}
	for _, tt := range tests {
		result := tt.input * tt.input
		if result == tt.expected {
			fmt.Printf("PASS: %d^2 = %d\n", tt.input, result)
		}
	}
}
