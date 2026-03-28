// Maps Basics — standalone demonstration of map creation, CRUD,
// existence checks, and nil map behavior in Go.
//
// Run: go run ./cmd/concepts/maps/01-basics
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
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Maps: Basics (CRUD & nil Behavior)     %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Creation ---
	fmt.Printf("%s▸ Map Creation%s\n", cyan+bold, reset)

	// Create with make (PREFERRED for non-literal maps)
	m := make(map[string]int)
	fmt.Printf("  %s✔ make(map[string]int) creates an empty, initialized map%s\n", green, reset)

	// Map literal
	scores := map[string]int{
		"Alice": 95,
		"Bob":   87,
		"Carol": 92,
	}
	fmt.Printf("  %s✔ Map literal: %s%v%s\n", green, magenta, scores, reset)
	fmt.Printf("  %s✔ Under the hood: backed by runtime.hmap — hash table with 8-entry buckets%s\n\n", green, reset)

	// --- CRUD ---
	fmt.Printf("%s▸ CRUD Operations%s\n", cyan+bold, reset)

	// Insert / Update
	m["go"] = 100
	m["python"] = 90
	fmt.Printf("  Insert: m[\"go\"] = %s%d%s, m[\"python\"] = %s%d%s\n", magenta, m["go"], reset, magenta, m["python"], reset)

	m["go"] = 110 // update
	fmt.Printf("  Update: m[\"go\"] = %s%d%s (same syntax — key exists so it overwrites)\n", magenta, m["go"], reset)

	// Read
	fmt.Printf("  Read:   m[\"go\"] = %s%d%s\n", magenta, m["go"], reset)

	// Reading a missing key returns ZERO VALUE (no panic!)
	// This is a deliberate design choice — maps always return a valid
	// value, so you never get a nil pointer dereference from a read.
	fmt.Printf("  Missing: m[\"rust\"] = %s%d%s\n", magenta, m["rust"], reset)
	fmt.Printf("  %s✔ Reading a missing key returns zero value (0 for int) — no panic%s\n", green, reset)
	fmt.Printf("  %s⚠ Can't distinguish \"key with value 0\" from \"absent key\" without comma-ok%s\n\n", yellow, reset)

	// --- Comma-ok ---
	fmt.Printf("%s▸ Existence Check (comma-ok idiom)%s\n", cyan+bold, reset)

	// EXISTENCE CHECK — the comma-ok idiom (critical pattern!)
	// The second return value is a bool indicating presence.
	// This is Go's answer to "how do I distinguish zero value from absent?"
	val, ok := m["rust"]
	fmt.Printf("  val, ok := m[\"rust\"] → val = %s%d%s, ok = %s%v%s\n", magenta, val, reset, magenta, ok, reset)
	if ok {
		fmt.Printf("  Found: %d\n", val)
	} else {
		fmt.Printf("  %s✔ ok=false confirms the key is absent, not just zero-valued%s\n", green, reset)
	}

	// Short form — idiomatic Go: declare and check in one statement
	if v, ok := scores["Alice"]; ok {
		fmt.Printf("  Short form: scores[\"Alice\"] → v = %s%d%s, ok = %s%v%s\n", magenta, v, reset, magenta, ok, reset)
	}
	fmt.Println()

	// --- Delete ---
	fmt.Printf("%s▸ Delete%s\n", cyan+bold, reset)
	fmt.Printf("  Before: %s%v%s\n", magenta, m, reset)

	// Delete — deleting a missing key is a no-op (no panic)
	delete(m, "python")
	fmt.Printf("  After delete(m, \"python\"): %s%v%s\n", magenta, m, reset)
	fmt.Printf("  %s✔ Deleting a missing key is a no-op — no panic%s\n\n", green, reset)

	// --- Length ---
	fmt.Printf("%s▸ Length%s\n", cyan+bold, reset)
	fmt.Printf("  len(scores) = %s%d%s\n\n", magenta, len(scores), reset)

	// --- nil map ---
	fmt.Printf("%s▸ nil Map Behavior%s\n", cyan+bold, reset)

	// nil map — READ is safe, WRITE causes panic!
	// A nil map has no underlying hash table allocated. Reads return
	// zero values safely, but writes trigger: "assignment to entry in nil map"
	var nilMap map[string]int
	fmt.Printf("  var nilMap map[string]int → nilMap == nil: %s%v%s\n", magenta, nilMap == nil, reset)
	fmt.Printf("  nilMap[\"key\"] = %s%d%s (read is safe — returns zero value)\n", magenta, nilMap["key"], reset)
	fmt.Printf("  %s✔ Read from nil map: safe, returns zero value%s\n", green, reset)
	fmt.Printf("  %s⚠ WRITE to nil map PANICS: \"assignment to entry in nil map\"%s\n", yellow, reset)
	fmt.Printf("  %s⚠ Always initialize with make() or a literal before writing!%s\n", yellow, reset)
	// nilMap["key"] = 1 // PANIC: assignment to entry in nil map
}
