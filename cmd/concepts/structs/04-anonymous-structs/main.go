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
	fmt.Printf("%s%s  Structs: Anonymous Structs             %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- One-off Struct ---
	fmt.Printf("%s▸ One-off Anonymous Struct%s\n", cyan+bold, reset)

	// Define and initialize in one step
	point := struct {
		X, Y int
	}{X: 10, Y: 20}
	fmt.Printf("  point := struct{X,Y int}{10,20} → %s%+v%s\n", magenta, point, reset)
	fmt.Printf("  %s✔ No named type needed — defined and used in place%s\n", green, reset)
	fmt.Printf("  %s✔ Useful for one-off data grouping, API responses, config bundles%s\n\n", green, reset)

	// --- Table-Driven Tests ---
	fmt.Printf("%s▸ Table-Driven Test Pattern (idiomatic Go)%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Each test case is an anonymous struct in a slice%s\n", green, reset)
	fmt.Printf("  %s✔ Named fields serve as documentation — self-describing test cases%s\n\n", green, reset)

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
			fmt.Printf("  %s✔ PASS%s: %s%d%s² = %s%d%s\n", green, reset, magenta, tt.input, reset, magenta, result, reset)
		} else {
			fmt.Printf("  %s✘ FAIL%s: %d² = %d (expected %d)\n", red, reset, tt.input, result, tt.expected)
		}
	}

	fmt.Printf("\n%s▸ Why Anonymous Structs?%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Avoid polluting package namespace with single-use types%s\n", green, reset)
	fmt.Printf("  %s✔ Table-driven tests: easy to add cases — just append a struct literal%s\n", green, reset)
	fmt.Printf("  %s✔ JSON decoding: decode into anonymous struct when you only need a few fields%s\n", green, reset)
	fmt.Printf("  %s⚠ If you use the same shape in multiple places, extract a named type instead%s\n", yellow, reset)
}
