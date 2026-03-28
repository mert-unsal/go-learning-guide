// Nested Maps — standalone demonstration of map-of-maps for
// representing graph adjacency lists and 2D data structures.
//
// Run: go run ./cmd/concepts/maps/04-nested-maps
package main

import "fmt"

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
	addEdge("A", "C", 4)
	addEdge("B", "C", 2)

	fmt.Println("Graph:", graph)
}
