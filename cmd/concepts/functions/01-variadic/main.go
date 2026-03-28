// Standalone demo: Variadic Functions in Go
//
// ...Type means the function accepts any number of that type.
// The variadic parameter is treated as a SLICE inside the function.
//
// Under the hood: the compiler allocates a []T on the stack (or heap if
// it escapes) and packs the arguments into it. When you spread an existing
// slice with nums..., NO copy is made — the same backing array is reused.
// This is why you must never mutate a variadic param if callers pass slices.
//
// Run:  go run .
package main

import "fmt"

// sum accepts zero or more ints via the variadic ...int parameter.
// Inside the function, nums is a plain []int slice.
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	fmt.Println(sum(1, 2, 3))       // 6
	fmt.Println(sum(1, 2, 3, 4, 5)) // 15

	// Spread a slice with the ... operator.
	// This passes the slice's backing array directly — no allocation.
	nums := []int{10, 20, 30}
	fmt.Println(sum(nums...)) // 60
}
