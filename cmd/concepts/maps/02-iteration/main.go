// Maps Iteration — standalone demonstration of map iteration order
// and the sorted-keys pattern.
//
// Run: go run ./cmd/concepts/maps/02-iteration
package main

import (
	"fmt"
	"sort"
)

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
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Maps: Iteration Order                  %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	m := map[string]int{"c": 3, "a": 1, "b": 2}

	// --- Random Iteration ---
	fmt.Printf("%s▸ Random Iteration (pass 1)%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ Go runtime randomizes iteration start position (since Go 1.12)%s\n", yellow, reset)
	fmt.Printf("  %s⚠ Never write code that depends on map ordering!%s\n", yellow, reset)

	// Direct iteration — RANDOM ORDER
	// Each run may produce a different order. Never rely on this.
	for k, v := range m {
		fmt.Printf("    %s%s%s: %s%d%s\n", magenta, k, reset, magenta, v, reset)
	}

	fmt.Printf("\n%s▸ Random Iteration (pass 2 — may differ!)%s\n", cyan+bold, reset)
	for k, v := range m {
		fmt.Printf("    %s%s%s: %s%d%s\n", magenta, k, reset, magenta, v, reset)
	}
	fmt.Printf("  %s✔ Each range loop picks a random starting bucket — order is non-deterministic%s\n\n", green, reset)

	// --- Sorted Keys Pattern ---
	fmt.Printf("%s▸ Sorted Keys Pattern (deterministic output)%s\n", cyan+bold, reset)

	// Ordered iteration: extract keys, sort, then iterate by key.
	// Pre-allocate the slice with cap = len(m) to avoid re-allocation.
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Printf("  %s✔ Extract keys → sort.Strings() → iterate by sorted key%s\n", green, reset)
	fmt.Printf("  %s✔ Pre-allocate: make([]string, 0, len(m)) avoids growslice reallocation%s\n", green, reset)
	for _, k := range keys {
		fmt.Printf("    %s%s%s: %s%d%s\n", magenta, k, reset, magenta, m[k], reset)
	}
	fmt.Printf("  %s✔ Output is now deterministic — safe for tests, serialization, and display%s\n", green, reset)
}
