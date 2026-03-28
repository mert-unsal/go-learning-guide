// Maps Iteration — standalone demonstration of map iteration order
// and the sorted-keys pattern.
//
// Run: go run ./cmd/concepts/maps/02-iteration
package main

import (
	"fmt"
	"sort"
)

// ============================================================
// MAP ITERATION
// ============================================================
// Maps in Go have RANDOM iteration order.
// The runtime intentionally randomizes iteration starting position
// (since Go 1.12) to prevent code from depending on order.
//
// If you need deterministic output — e.g., for tests, serialization,
// or user-facing display — sort the keys first.

func main() {
	m := map[string]int{"c": 3, "a": 1, "b": 2}

	// Direct iteration — RANDOM ORDER
	// Each run may produce a different order. Never rely on this.
	fmt.Println("Random order:")
	for k, v := range m {
		fmt.Printf("  %s: %d\n", k, v)
	}

	// Ordered iteration: extract keys, sort, then iterate by key.
	// Pre-allocate the slice with cap = len(m) to avoid re-allocation.
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("Sorted order:")
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, m[k])
	}
}
