// Nested Maps — standalone demonstration of map-of-maps for
// representing graph adjacency lists and 2D data structures.
//
// Run: go run ./cmd/concepts/maps/04-nested-maps
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

// ============================================================
// MAP OF MAPS (2D map)
// ============================================================
// Nested maps are used for adjacency lists, 2D grids, and
// multi-dimensional lookups. The critical pattern: always check
// if the inner map is nil before writing to it.
//
// graph[from][to] = weight  ← panics if graph[from] is nil!
//
// Always initialize inner maps before use. The addEdge closure
// below shows the safe pattern.

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Maps: Nested Maps (Map of Maps)        %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Building a Graph ---
	fmt.Printf("%s▸ Graph Adjacency List (map[string]map[string]int)%s\n", cyan+bold, reset)

	// Adjacency list for a weighted directed graph
	graph := make(map[string]map[string]int)

	// Safe way to add nested map entries:
	// Check if the inner map exists; if nil, initialize it first.
	// This pattern applies to any map[K1]map[K2]V structure.
	addEdge := func(from, to string, weight int) {
		if graph[from] == nil {
			graph[from] = make(map[string]int)
		}
		graph[from][to] = weight
	}

	addEdge("A", "B", 1)
	fmt.Printf("  addEdge(\"A\",\"B\",1) → graph[\"A\"] was nil, initialized inner map first\n")
	addEdge("A", "C", 4)
	fmt.Printf("  addEdge(\"A\",\"C\",4) → graph[\"A\"] already exists, just set key\n")
	addEdge("B", "C", 2)
	fmt.Printf("  addEdge(\"B\",\"C\",2) → graph[\"B\"] was nil, initialized inner map first\n")
	fmt.Printf("  %s✔ Always check if inner map is nil before writing — prevents panic%s\n\n", green, reset)

	// --- nil Inner Map Danger ---
	fmt.Printf("%s▸ nil Inner Map Gotcha%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ graph[\"X\"][\"Y\"] = 1 panics if graph[\"X\"] is nil!%s\n", yellow, reset)
	fmt.Printf("  %s⚠ The outer map returns a nil map for missing keys — writing to it panics%s\n", yellow, reset)
	fmt.Printf("  %s✔ Safe pattern: if graph[from] == nil { graph[from] = make(map[string]int) }%s\n\n", green, reset)

	// --- Traversal ---
	fmt.Printf("%s▸ Graph Traversal%s\n", cyan+bold, reset)
	for from, edges := range graph {
		for to, weight := range edges {
			fmt.Printf("    %s%s%s → %s%s%s (weight: %s%d%s)\n", magenta, from, reset, magenta, to, reset, magenta, weight, reset)
		}
	}
	fmt.Printf("  %s✔ Full graph: %s%v%s\n", green, magenta, graph, reset)
	fmt.Printf("  %s⚠ Iteration order is random for both outer and inner maps%s\n", yellow, reset)
}
