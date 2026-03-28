// Standalone demo: Pass by Reference with Pointers
//
// Go is ALWAYS pass-by-value. When you pass a pointer, the pointer itself
// is copied, but both copies point to the same memory — giving you
// reference semantics. This is identical to how Java passes object refs.
//
// Slices, maps, and channels are already reference types (they contain an
// internal pointer to the backing data), so you rarely need explicit
// pointers for them.
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

// WITHOUT pointer: modifies a COPY, original unchanged
func incrementValue(n int) {
	n++ // modifies local copy only
}

// WITH pointer: modifies the ORIGINAL
func incrementPointer(n *int) {
	*n++ // dereferences and modifies the original
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Pointers: Pass by Reference            %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	x := 10

	// --- Pass by Value ---
	fmt.Printf("%s▸ Pass by Value (copy)%s\n", cyan+bold, reset)
	fmt.Printf("  x before incrementValue: %s%d%s (addr: %s%p%s)\n", magenta, x, reset, magenta, &x, reset)
	incrementValue(x)
	fmt.Printf("  x after  incrementValue: %s%d%s (addr: %s%p%s)\n", magenta, x, reset, magenta, &x, reset)
	fmt.Printf("  %s⚠ Function received a COPY of x — modifying it doesn't affect original%s\n", yellow, reset)
	fmt.Printf("  %s✔ Go is always pass-by-value. The int is copied into the function's stack frame%s\n\n", green, reset)

	// --- Pass by Pointer ---
	fmt.Printf("%s▸ Pass by Pointer (reference semantics)%s\n", cyan+bold, reset)
	fmt.Printf("  x before incrementPointer: %s%d%s (addr: %s%p%s)\n", magenta, x, reset, magenta, &x, reset)
	fmt.Printf("  Passing &x = %s%p%s to incrementPointer\n", magenta, &x, reset)
	incrementPointer(&x)
	fmt.Printf("  x after  incrementPointer: %s%d%s (addr: %s%p%s)\n", magenta, x, reset, magenta, &x, reset)
	fmt.Printf("  %s✔ The pointer itself was copied, but both copies point to the same memory%s\n", green, reset)
	fmt.Printf("  %s✔ *n++ dereferences the pointer and increments the value at that address%s\n\n", green, reset)

	// --- Reference Types ---
	fmt.Printf("%s▸ Built-in Reference Types%s\n", cyan+bold, reset)
	// Slices, maps, and channels are already reference types —
	// you don't need pointers for them in most cases.
	fmt.Printf("  %s✔ Slices, maps, and channels contain internal pointers already%s\n", green, reset)
	fmt.Printf("  %s✔ Passing a slice copies the header {ptr, len, cap} but shares the backing array%s\n", green, reset)
	fmt.Printf("  %s✔ Passing a map copies the *hmap pointer — both references see the same data%s\n", green, reset)
	fmt.Printf("  %s⚠ You rarely need *[]T or *map[K]V — only when you need to reassign the entire value%s\n", yellow, reset)
}
