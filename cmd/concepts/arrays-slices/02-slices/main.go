// 02-slices demonstrates Go slices — the primary dynamic collection type.
//
// Run:  go run .
//
// ============================================================
// SLICES — The primary collection type in Go
// ============================================================
// A slice is a DYNAMIC VIEW into an array.
// Internally: { pointer to array, length, capacity }
//   — 3-word struct (24 bytes on 64-bit): runtime represents this as
//     reflect.SliceHeader { Data uintptr; Len int; Cap int }.
// Slices are reference types — they share the underlying array.
//
// Key distinction:
//   nil slice  → var s []int        → s == nil is true,  len=0, cap=0
//   empty slice → s := []int{}      → s == nil is false, len=0, cap=0
//   Both work fine with append, range, and len.
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
	fmt.Printf("%s%s  Slices — Dynamic Views Into Arrays      %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// From literal
	fmt.Printf("%s▸ Slice Literal%s\n", cyan+bold, reset)
	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("  s := []int{1,2,3,4,5}\n")
	fmt.Printf("  values=%v  len=%s%d%s  cap=%s%d%s\n", s, magenta, len(s), reset, magenta, cap(s), reset)
	fmt.Printf("  %s✔ Slice header is a 3-word struct: {pointer, len, cap} — 24 bytes on 64-bit%s\n\n", green, reset)

	// make([]T, length, capacity) — allocates with specific size
	fmt.Printf("%s▸ make() — Preallocate With Specific Size%s\n", cyan+bold, reset)
	s2 := make([]int, 3)     // [0 0 0], len=3, cap=3
	s3 := make([]int, 3, 10) // [0 0 0], len=3, cap=10
	fmt.Printf("  make([]int, 3)     → %v  len=%s%d%s  cap=%s%d%s\n", s2, magenta, len(s2), reset, magenta, cap(s2), reset)
	fmt.Printf("  make([]int, 3, 10) → %v  len=%s%d%s  cap=%s%d%s\n", s3, magenta, len(s3), reset, magenta, cap(s3), reset)
	fmt.Printf("  %s✔ Pre-allocating cap avoids repeated growslice calls in hot paths%s\n\n", green, reset)

	// nil slice vs empty slice
	fmt.Printf("%s▸ nil Slice vs Empty Slice%s\n", cyan+bold, reset)
	var nilSlice []int    // nil, len=0, cap=0
	emptySlice := []int{} // not nil, len=0, cap=0
	fmt.Printf("  var nilSlice []int    → nilSlice == nil: %s%t%s  len=%s%d%s  cap=%s%d%s\n",
		magenta, nilSlice == nil, reset, magenta, len(nilSlice), reset, magenta, cap(nilSlice), reset)
	fmt.Printf("  emptySlice := []int{} → emptySlice == nil: %s%t%s  len=%s%d%s  cap=%s%d%s\n",
		magenta, emptySlice == nil, reset, magenta, len(emptySlice), reset, magenta, cap(emptySlice), reset)
	fmt.Printf("  %s✔ Both work with append, range, len — prefer nil slice (var s []int) as zero value%s\n", green, reset)
	fmt.Printf("  %s⚠ JSON: nil slice marshals to \"null\", empty slice to \"[]\" — matters for APIs!%s\n\n", yellow, reset)

	// Slicing — s[low:high] — includes low, excludes high
	fmt.Printf("%s▸ Slice Expressions — s[low:high]%s\n", cyan+bold, reset)
	fmt.Printf("  s = %v\n", s)
	fmt.Printf("  s[1:3] = %s%v%s  — elements at index 1,2 (high is exclusive)\n", magenta, s[1:3], reset)
	fmt.Printf("  s[:3]  = %s%v%s  — first 3 elements\n", magenta, s[:3], reset)
	fmt.Printf("  s[2:]  = %s%v%s  — from index 2 to end\n", magenta, s[2:], reset)
	fmt.Printf("  s[:]   = %s%v%s  — full slice (same backing array)\n", magenta, s[:], reset)
	fmt.Printf("  %s✔ New slice shares the same backing array — cap is measured from low to end of backing array%s\n\n", green, reset)

	// IMPORTANT: slices share underlying array!
	fmt.Printf("%s▸ Shared Backing Array (Aliasing Trap)%s\n", cyan+bold, reset)
	a := []int{1, 2, 3, 4, 5}
	b := a[1:3] // b = [2, 3], shares array with a
	fmt.Printf("  a := []int{1,2,3,4,5}\n")
	fmt.Printf("  b := a[1:3] → b = %v  (shares a's backing array)\n", b)
	b[0] = 99
	fmt.Printf("  b[0] = 99   → a is now %s%v%s\n", magenta, a, reset)
	fmt.Printf("  %s⚠ Mutating b changed a! Slices share memory until a new allocation occurs%s\n", yellow, reset)
	fmt.Printf("  %s✔ Use full slice expression a[1:3:3] to limit cap and prevent accidental aliasing%s\n\n", green, reset)

	fmt.Printf("%s%s── Key Takeaways ──%s\n", bold, blue, reset)
	fmt.Printf("  %s✔ Slice = {ptr, len, cap} — a lightweight view, not a copy%s\n", green, reset)
	fmt.Printf("  %s✔ make([]T, 0, expectedSize) to pre-allocate and avoid GC pressure%s\n", green, reset)
	fmt.Printf("  %s⚠ Sub-slices share memory — copy() or full-slice-expr to detach%s\n", yellow, reset)
}
