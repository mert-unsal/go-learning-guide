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
	nums := []int{1, 2, 3, 4, 5}

	// Pass an anonymous function (lambda)
	doubled := apply(nums, func(n int) int {
		return n * 2
	})
	fmt.Println(doubled) // [2 4 6 8 10]

	// Assign a function to a variable
	square := func(n int) int { return n * n }
	squared := apply(nums, square)
	fmt.Println(squared) // [1 4 9 16 25]
}
