// Maps Basics — standalone demonstration of map creation, CRUD,
// existence checks, and nil map behavior in Go.
//
// Run: go run ./cmd/concepts/maps/01-basics
package main

import "fmt"

// ============================================================
// MAP BASICS
// ============================================================
// A map is an unordered collection of key-value pairs.
// Maps are REFERENCE TYPES — passed by reference automatically.
// Key type must be COMPARABLE (no slices, maps, or funcs as keys).
//
// Under the hood: maps are hash tables backed by runtime.hmap.
// The runtime uses a hash function specific to the key type and
// organizes data into buckets (each bucket holds 8 key-value pairs).
// Growth triggers at a load factor of 6.5 — the map doubles its
// bucket count and incrementally rehashes (amortized O(1) insert).

func main() {
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
	// This is a deliberate design choice — maps always return a valid
	// value, so you never get a nil pointer dereference from a read.
	fmt.Println("rust:", m["rust"]) // 0 — not found, zero value

	// EXISTENCE CHECK — the comma-ok idiom (critical pattern!)
	// The second return value is a bool indicating presence.
	// This is Go's answer to "how do I distinguish zero value from absent?"
	val, ok := m["rust"]
	if ok {
		fmt.Println("Found:", val)
	} else {
		fmt.Println("rust not found, zero value:", val)
	}

	// Short form — idiomatic Go: declare and check in one statement
	if v, ok := scores["Alice"]; ok {
		fmt.Println("Alice's score:", v)
	}

	// Delete — deleting a missing key is a no-op (no panic)
	delete(m, "python")
	fmt.Println("After delete:", m)

	// Length
	fmt.Println("Length:", len(scores))

	// nil map — READ is safe, WRITE causes panic!
	// A nil map has no underlying hash table allocated. Reads return
	// zero values safely, but writes trigger: "assignment to entry in nil map"
	var nilMap map[string]int
	fmt.Println("nil map read:", nilMap["key"]) // 0, no panic
	// nilMap["key"] = 1 // PANIC: assignment to entry in nil map
}
