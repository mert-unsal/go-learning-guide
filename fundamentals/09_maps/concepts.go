// Package maps covers Go maps: creation, CRUD operations,
// existence checks, iteration, and common patterns for coding exams.
package maps

import (
	"fmt"
	"sort"
)

// ============================================================
// 1. MAP BASICS
// ============================================================
// A map is an unordered collection of key-value pairs.
// Maps are REFERENCE TYPES — passed by reference automatically.
// Key type must be COMPARABLE (no slices, maps, or funcs as keys).

func DemonstrateBasics() {
	// Create with make (PREFERRED for non-literal maps)
	m := make(map[string]int)

	// Map literal
	scores := map[string]int{
		"Alice": 95,
		"Bob":   87,
		"Carol": 92,
	}
	fmt.Println(scores)

	// Insert / Update
	m["go"] = 100
	m["python"] = 90
	m["go"] = 110 // update

	// Read
	fmt.Println("go:", m["go"]) // 110

	// Reading a missing key returns ZERO VALUE (no panic!)
	fmt.Println("rust:", m["rust"]) // 0 — not found, zero value

	// EXISTENCE CHECK — the comma-ok idiom (critical pattern!)
	val, ok := m["rust"]
	if ok {
		fmt.Println("Found:", val)
	} else {
		fmt.Println("rust not found, zero value:", val)
	}

	// Short form
	if v, ok := scores["Alice"]; ok {
		fmt.Println("Alice's score:", v)
	}

	// Delete
	delete(m, "python")
	fmt.Println("After delete:", m)

	// Length
	fmt.Println("Length:", len(scores))

	// nil map — READ is safe, WRITE causes panic!
	var nilMap map[string]int
	fmt.Println("nil map read:", nilMap["key"]) // 0, no panic
	// nilMap["key"] = 1 // PANIC: assignment to entry in nil map
}

// ============================================================
// 2. MAP ITERATION
// ============================================================
// Maps in Go have RANDOM iteration order.
// If you need ordered output, sort the keys first.

func DemonstrateIteration() {
	m := map[string]int{"c": 3, "a": 1, "b": 2}

	// Direct iteration — RANDOM ORDER
	fmt.Println("Random order:")
	for k, v := range m {
		fmt.Printf("  %s: %d\n", k, v)
	}

	// Ordered iteration: sort keys first
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

// ============================================================
// 3. COMMON MAP PATTERNS (Essential for LeetCode!)
// ============================================================

// Pattern 1: Frequency counter (character/number frequency)
func FrequencyCount(s string) map[rune]int {
	freq := make(map[rune]int)
	for _, ch := range s {
		freq[ch]++ // if key doesn't exist, starts at zero value 0, then +1
	}
	return freq
}

// Pattern 2: Set (use map[T]bool or map[T]struct{})
func UniqueElements(nums []int) []int {
	seen := make(map[int]struct{}) // struct{} uses NO memory (vs bool)
	result := []int{}
	for _, n := range nums {
		if _, ok := seen[n]; !ok {
			seen[n] = struct{}{}
			result = append(result, n)
		}
	}
	return result
}

// Pattern 3: Grouping (group items by a key)
func GroupByLength(words []string) map[int][]string {
	groups := make(map[int][]string)
	for _, w := range words {
		n := len(w)
		groups[n] = append(groups[n], w) // append to nil slice works!
	}
	return groups
}

// Pattern 4: Two-pass map (e.g., Two Sum)
func TwoSum(nums []int, target int) (int, int) {
	// Map: value → index
	seen := make(map[int]int)
	for i, n := range nums {
		complement := target - n
		if j, ok := seen[complement]; ok {
			return j, i
		}
		seen[n] = i
	}
	return -1, -1
}

// Pattern 5: Memoization (cache expensive results)
func fibMemo(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}
	if v, ok := memo[n]; ok {
		return v
	}
	result := fibMemo(n-1, memo) + fibMemo(n-2, memo)
	memo[n] = result
	return result
}

func Fibonacci(n int) int {
	return fibMemo(n, make(map[int]int))
}

func DemonstratePatterns() {
	// Frequency
	freq := FrequencyCount("hello world")
	fmt.Println("Frequency:", freq)

	// Unique
	fmt.Println("Unique:", UniqueElements([]int{1, 2, 2, 3, 3, 3, 4}))

	// Grouping
	groups := GroupByLength([]string{"go", "is", "fun", "and", "fast"})
	fmt.Println("Groups by length:", groups)

	// Two Sum
	i, j := TwoSum([]int{2, 7, 11, 15}, 9)
	fmt.Println("TwoSum indices:", i, j) // 0, 1

	// Fibonacci with memoization
	fmt.Println("Fib(10):", Fibonacci(10)) // 55
}

// ============================================================
// 4. MAP OF MAPS (2D map)
// ============================================================

func DemonstrateNestedMaps() {
	// Adjacency list for a graph
	graph := make(map[string]map[string]int)

	// Safe way to add nested map entries
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

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Map Basics ===")
	DemonstrateBasics()
	fmt.Println("\n=== Map Iteration ===")
	DemonstrateIteration()
	fmt.Println("\n=== Map Patterns ===")
	DemonstratePatterns()
	fmt.Println("\n=== Nested Maps ===")
	DemonstrateNestedMaps()
}
