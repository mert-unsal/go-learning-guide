// Standalone demo: Functions as First-Class Values
//
// Functions are values in Go. They can be:
//   - Assigned to variables
//   - Passed as arguments
//   - Returned from functions
//
// Under the hood each function value is a pointer-sized closure struct
// (runtime.funcval). When the function captures no variables, the compiler
// can point all references at a single static funcval — zero allocation.
// Captured variables cause the funcval (and the variables) to escape to heap.
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

// apply takes a slice and a transformation function, returning a new slice
// with fn applied to every element. This is the "map" higher-order pattern.
func apply(nums []int, fn func(int) int) []int {
	result := make([]int, len(nums))
	for i, v := range nums {
		result[i] = fn(v)
	}
	return result
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Functions as First-Class Values         %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	nums := []int{1, 2, 3, 4, 5}
	fmt.Printf("  input: %s%v%s\n\n", magenta, nums, reset)

	// --- Anonymous function (lambda) ---
	fmt.Printf("%s▸ Passing an anonymous function (lambda)%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Functions are values — you can pass them inline just like an int or string%s\n", green, reset)
	doubled := apply(nums, func(n int) int {
		return n * 2
	})
	fmt.Printf("  apply(nums, func(n) n*2) = %s%v%s\n", magenta, doubled, reset)
	fmt.Printf("  %s✔ No capture → compiler uses a single static runtime.funcval (zero alloc)%s\n", green, reset)

	// --- Assigned function variable ---
	fmt.Printf("\n%s▸ Assigning a function to a variable%s\n", cyan+bold, reset)
	square := func(n int) int { return n * n }
	squared := apply(nums, square)
	fmt.Printf("  square := func(n) n*n\n")
	fmt.Printf("  apply(nums, square)      = %s%v%s\n", magenta, squared, reset)
	fmt.Printf("  %s✔ 'square' holds a function value — its type is func(int) int%s\n", green, reset)

	// --- Higher-order pattern ---
	fmt.Printf("\n%s▸ The higher-order 'map' pattern%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ apply() is a generic transformer: takes data + strategy, returns new data%s\n", green, reset)
	fmt.Printf("  %s✔ This is the foundation of functional patterns in Go%s\n", green, reset)
	fmt.Printf("  %s⚠ Under the hood: each func value is a pointer-sized closure struct (runtime.funcval)%s\n", yellow, reset)
	fmt.Printf("  %s⚠ Captured variables cause the funcval to escape to heap — watch in hot paths%s\n", yellow, reset)
}
