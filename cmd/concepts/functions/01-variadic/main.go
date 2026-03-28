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
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Variadic Functions in Go               %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Calling with individual arguments ---
	fmt.Printf("%s▸ Calling with individual arguments%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ The compiler packs args into a []int behind the scenes%s\n", green, reset)

	r1 := sum(1, 2, 3)
	fmt.Printf("  sum(1, 2, 3)       = %s%d%s\n", magenta, r1, reset)

	r2 := sum(1, 2, 3, 4, 5)
	fmt.Printf("  sum(1, 2, 3, 4, 5) = %s%d%s\n", magenta, r2, reset)

	r3 := sum()
	fmt.Printf("  sum()              = %s%d%s  ← zero args is valid; range over nil slice is a no-op\n", magenta, r3, reset)

	// --- Spreading a slice ---
	fmt.Printf("\n%s▸ Spreading a slice with the ... operator%s\n", cyan+bold, reset)
	nums := []int{10, 20, 30}
	fmt.Printf("  nums := %s%v%s\n", magenta, nums, reset)

	r4 := sum(nums...)
	fmt.Printf("  sum(nums...) = %s%d%s\n", magenta, r4, reset)
	fmt.Printf("  %s✔ No copy is made — the same backing array is reused%s\n", green, reset)
	fmt.Printf("  %s⚠ Never mutate a variadic param when callers pass slices!%s\n", yellow, reset)

	// --- Under the hood ---
	fmt.Printf("\n%s▸ Under the hood%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Variadic param is a regular []T slice inside the function%s\n", green, reset)
	fmt.Printf("  %s✔ With individual args: compiler allocates []T on stack (or heap if it escapes)%s\n", green, reset)
	fmt.Printf("  %s✔ With slice spread (nums...): zero-copy, backing array shared directly%s\n", green, reset)
}
