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

func main() {
	src := []int{1, 2, 3, 4, 5}

	// copy(dst, src) — copies min(len(dst), len(src)) elements
	dst := make([]int, len(src))
	n := copy(dst, src)
	fmt.Println("Copied", n, "elements:", dst)

	// Now dst is independent from src
	dst[0] = 99
	fmt.Println("src:", src) // unchanged
	fmt.Println("dst:", dst)

	// Partial copy
	partial := make([]int, 3)
	copy(partial, src)
	fmt.Println("Partial:", partial) // [1 2 3]
}
