// 03-append demonstrates the append built-in and capacity growth.
//
// Run:  go run .
//
// ============================================================
// APPEND
// ============================================================
// append returns a NEW slice header. Always reassign: s = append(s, v).
//
// Growth strategy (runtime/slice.go → growslice):
//   cap < 256  → double capacity (2x)
//   cap >= 256 → grow by ~1.25x + 192
// When cap is exceeded, Go allocates a NEW backing array —
// the old one is left for GC. This means previously-shared
// slices will NO LONGER see each other's mutations.
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
	fmt.Printf("%s%s  Append — Growth & Reallocation          %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// Append single element
	fmt.Printf("%s▸ Append Single Element%s\n", cyan+bold, reset)
	s := []int{1, 2, 3}
	fmt.Printf("  before: %v  len=%s%d%s  cap=%s%d%s\n", s, magenta, len(s), reset, magenta, cap(s), reset)
	s = append(s, 4)
	fmt.Printf("  after append(s, 4): %v  len=%s%d%s  cap=%s%d%s\n", s, magenta, len(s), reset, magenta, cap(s), reset)
	fmt.Printf("  %s✔ Cap doubled from 3→6 because cap < 256 triggers 2x growth (runtime.growslice)%s\n\n", green, reset)

	// Append multiple elements
	fmt.Printf("%s▸ Append Multiple Elements%s\n", cyan+bold, reset)
	s = append(s, 5, 6, 7)
	fmt.Printf("  append(s, 5, 6, 7): %v  len=%s%d%s  cap=%s%d%s\n", s, magenta, len(s), reset, magenta, cap(s), reset)
	fmt.Printf("  %s✔ Variadic append — compiler packs args into a single growslice call%s\n\n", green, reset)

	// Append another slice with ...
	fmt.Printf("%s▸ Append Another Slice (... Spread)%s\n", cyan+bold, reset)
	extra := []int{8, 9, 10}
	s = append(s, extra...)
	fmt.Printf("  append(s, extra...): %v\n", s)
	fmt.Printf("  len=%s%d%s  cap=%s%d%s\n", magenta, len(s), reset, magenta, cap(s), reset)
	fmt.Printf("  %s✔ The ... operator unpacks the slice — equivalent to append(s, 8, 9, 10)%s\n\n", green, reset)

	// Growth strategy demo
	fmt.Printf("%s▸ Capacity Growth Strategy%s\n", cyan+bold, reset)
	g := make([]int, 0)
	prevCap := cap(g)
	for i := 0; i < 20; i++ {
		g = append(g, i)
		if cap(g) != prevCap {
			fmt.Printf("  len=%s%2d%s  cap grew: %s%d → %d%s\n", magenta, len(g), reset, yellow, prevCap, cap(g), reset)
			prevCap = cap(g)
		}
	}
	fmt.Printf("  %s✔ Growth: 2x when cap < 256, then ~1.25x + 192 (see runtime/slice.go:growslice)%s\n\n", green, reset)

	// IMPORTANT: when cap is exceeded, Go allocates a NEW array
	fmt.Printf("%s▸ Reallocation Breaks Shared Backing Arrays%s\n", cyan+bold, reset)
	a := make([]int, 3, 3)
	b := a // b shares a's backing array
	fmt.Printf("  a := make([]int, 3, 3); b := a  — same backing array\n")
	a = append(a, 4) // cap exceeded! a gets a new backing array
	a[0] = 99
	fmt.Printf("  a = append(a, 4) — cap exceeded, new backing array allocated\n")
	fmt.Printf("  a[0] = %s%d%s, b[0] = %s%d%s\n", magenta, a[0], reset, magenta, b[0], reset)
	fmt.Printf("  %s⚠ After reallocation, a and b point to DIFFERENT arrays — mutations are isolated%s\n", yellow, reset)
	fmt.Printf("  %s✔ Always reassign: s = append(s, v) — append may return a new slice header%s\n\n", green, reset)

	fmt.Printf("%s%s── Key Takeaways ──%s\n", bold, blue, reset)
	fmt.Printf("  %s✔ append returns a NEW slice header — always reassign%s\n", green, reset)
	fmt.Printf("  %s✔ Pre-allocate with make([]T, 0, n) to avoid repeated growslice in hot paths%s\n", green, reset)
	fmt.Printf("  %s⚠ Reallocation silently breaks shared backing arrays — a subtle source of bugs%s\n", yellow, reset)
}
