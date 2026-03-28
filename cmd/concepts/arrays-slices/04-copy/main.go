// 04-copy demonstrates the copy built-in for independent slice duplication.
//
// Run:  go run .
//
// ============================================================
// COPY
// ============================================================
// copy(dst, src) copies min(len(dst), len(src)) elements.
// After copy, dst is completely independent from src —
// no shared backing array, no aliasing.
// This is the idiomatic way to detach a sub-slice from its parent.
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
	fmt.Printf("%s%s  Copy — Independent Slice Duplication     %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// copy(dst, src) — copies min(len(dst), len(src)) elements
	fmt.Printf("%s▸ Full Copy%s\n", cyan+bold, reset)
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	n := copy(dst, src)
	fmt.Printf("  src = %v\n", src)
	fmt.Printf("  dst := make([]int, %d); copy(dst, src)\n", len(src))
	fmt.Printf("  copied %s%d%s elements → dst = %s%v%s\n", magenta, n, reset, magenta, dst, reset)
	fmt.Printf("  %s✔ copy returns the number of elements copied: min(len(dst), len(src))%s\n\n", green, reset)

	// Now dst is independent from src
	fmt.Printf("%s▸ Independence After Copy%s\n", cyan+bold, reset)
	dst[0] = 99
	fmt.Printf("  dst[0] = 99\n")
	fmt.Printf("  src = %s%v%s  ← unchanged\n", magenta, src, reset)
	fmt.Printf("  dst = %s%v%s  ← independent copy\n", magenta, dst, reset)
	fmt.Printf("  %s✔ After copy, dst has its own backing array — no aliasing, no shared memory%s\n\n", green, reset)

	// Partial copy
	fmt.Printf("%s▸ Partial Copy (dst smaller than src)%s\n", cyan+bold, reset)
	partial := make([]int, 3)
	n = copy(partial, src)
	fmt.Printf("  partial := make([]int, 3); copy(partial, src)\n")
	fmt.Printf("  copied %s%d%s elements → partial = %s%v%s\n", magenta, n, reset, magenta, partial, reset)
	fmt.Printf("  %s✔ Only min(len(dst), len(src)) = min(3, 5) = 3 elements copied — no panic%s\n\n", green, reset)

	// Idiomatic clone pattern
	fmt.Printf("%s▸ Idiomatic Clone Pattern%s\n", cyan+bold, reset)
	original := []int{10, 20, 30}
	clone := make([]int, len(original))
	copy(clone, original)
	fmt.Printf("  clone := make([]int, len(original)); copy(clone, original)\n")
	fmt.Printf("  clone = %s%v%s\n", magenta, clone, reset)
	fmt.Printf("  %s✔ This is the standard Go pattern to detach a sub-slice from its parent%s\n", green, reset)
	fmt.Printf("  %s⚠ append([]int(nil), src...) also works but allocates via growslice — copy is explicit%s\n\n", yellow, reset)

	// Copy between overlapping slices
	fmt.Printf("%s▸ Overlapping Copy (Shift Elements)%s\n", cyan+bold, reset)
	data := []int{1, 2, 3, 4, 5}
	fmt.Printf("  before: %v\n", data)
	copy(data[1:], data[:4]) // shift right by 1
	fmt.Printf("  copy(data[1:], data[:4]) → %s%v%s\n", magenta, data, reset)
	fmt.Printf("  %s✔ copy handles overlapping src/dst correctly — uses memmove internally%s\n\n", green, reset)

	fmt.Printf("%s%s── Key Takeaways ──%s\n", bold, blue, reset)
	fmt.Printf("  %s✔ copy creates a fully independent duplicate — the safe way to detach slices%s\n", green, reset)
	fmt.Printf("  %s✔ Always pre-allocate dst with make() — copy does NOT grow the destination%s\n", green, reset)
	fmt.Printf("  %s⚠ copy silently copies fewer elements if dst is too small — no error, no panic%s\n", yellow, reset)
}
