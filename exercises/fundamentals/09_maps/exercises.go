package maps

// ============================================================
// EXERCISES — 09 Maps
// ============================================================
//
// Go maps are hash tables backed by runtime.hmap: bucket arrays,
// tophash optimization, load factor 6.5, incremental evacuation.
//
// These exercises test your understanding of:
//   - Basic map patterns: frequency, grouping, seen-set (§1-6)
//   - Go-specific map behavior: nil safety, comma-ok, iteration order (§7-12)
//
// Exercises 1-6:  Map algorithms — frequency, grouping, sets
// Exercises 7-12: Go map internals — nil maps, comma-ok, merge, invert
//
// Deep dive: learnings/02_maps_buckets_and_growth.md
// ============================================================

// Exercise 1:
// CharFrequency counts the frequency of each character in a string.
//
// LESSON: `range` over a string gives (byte index, rune) — not (index, byte).
// Runes handle Unicode correctly. "hello" → map[rune]int{'h':1, 'e':1, 'l':2, 'o':1}
func CharFrequency(s string) map[rune]int {
	// TODO: iterate runes, count each in a map
	return nil
}

// Exercise 2:
// GroupByFirstChar groups words by their first character.
//
// LESSON: Map of slices — the value is itself a slice. You must append to it.
// The zero value of a slice is nil, and append(nil, x) works fine.
func GroupByFirstChar(words []string) map[byte][]string {
	// TODO: group words by w[0] (first byte — fine for ASCII)
	return nil
}

// Exercise 3:
// TopTwoFrequent returns the two most frequent elements.
//
// LESSON: Two-pass approach: first build a frequency map, then scan for top-2.
// Alternative: use a heap (see patterns/). For exactly top-2, scanning is simpler.
func TopTwoFrequent(nums []int) []int {
	// TODO: build frequency map, then find top two by count
	return nil
}

// Exercise 4:
// IsAnagram checks if two strings are anagrams.
//
// LESSON: Classic map technique — count up for s, count down for t.
// If every count returns to 0, they use the same characters the same number of times.
// O(n) time, O(k) space where k = number of unique characters.
func IsAnagram(s, t string) bool {
	// TODO: count chars in s, decrement for t, check all zero
	return false
}

// Exercise 5:
// FirstDuplicate returns the first integer that appears more than once.
//
// LESSON: "Seen set" pattern — use map[T]bool (or map[T]struct{}) as a set.
// map[T]struct{} uses slightly less memory (empty struct = 0 bytes),
// but map[T]bool is more readable.
func FirstDuplicate(nums []int) int {
	// TODO: use a seen set, return first repeated element (or -1)
	return 0
}

// Exercise 6:
// WordCount counts occurrences of each word in a sentence.
//
// LESSON: strings.Fields splits on any whitespace (spaces, tabs, newlines).
// Prefer it over strings.Split(s, " ") which doesn't handle multiple spaces.
func WordCount(sentence string) map[string]int {
	// TODO: split sentence into words, count each
	return nil
}

// ── Go Map Internals — Exercises 7-12 ──

// Exercise 7:
// NilMapRead safely reads a key from a map that might be nil.
// Returns the value and whether the key exists (comma-ok pattern).
//
// KEY INSIGHT: Reading from a nil map is safe — returns zero value.
// WRITING to a nil map panics! This is a common Go gotcha.
//
//   var m map[string]int  // nil
//   v := m["key"]         // v = 0, no panic
//   m["key"] = 1          // PANIC: assignment to entry in nil map
//
// The comma-ok pattern distinguishes "key missing" from "key exists with zero value":
//   v, ok := m["key"]     // ok=false means key missing, ok=true means key exists
func NilMapRead(m map[string]int, key string) (value int, exists bool) {
	// TODO: use comma-ok pattern: v, ok := m[key]
	return 0, false
}

// Exercise 8:
// InvertMap swaps keys and values. Since multiple keys may map to the
// same value, the result is map[int][]string.
//
// LESSON: Map inversion is a common transform. When values aren't unique,
// the inverted map must hold slices. append(nil, x) works — no need to
// pre-initialize the slice.
//
// Example: {"a":1, "b":2, "c":1} → {1:["a","c"], 2:["b"]}
func InvertMap(m map[string]int) map[int][]string {
	// TODO: for each k,v → append k to result[v]
	return nil
}

// Exercise 9:
// MergeMaps merges b into a, using the resolve function for key collisions.
// Returns a new map (does not modify a or b).
//
// LESSON: Go maps don't have a built-in merge. You iterate and decide.
// The resolve function receives (valueFromA, valueFromB) and returns the winner.
//
// Example: MergeMaps({"x":1}, {"x":2, "y":3}, func(a,b int) int { return a+b })
//          → {"x":3, "y":3}
func MergeMaps(a, b map[string]int, resolve func(int, int) int) map[string]int {
	// TODO: copy a into result, then merge b using resolve for collisions
	return nil
}

// Exercise 10:
// SetDifference returns keys that are in set a but NOT in set b.
// Both maps represent sets (the bool value doesn't matter — only key existence).
// Return the result sorted alphabetically.
//
// LESSON: map[T]bool is the idiomatic Go "set". Some use map[T]struct{}
// for zero-byte values, but map[T]bool is more readable.
// Checking membership: if set[key] { ... } — clean and idiomatic.
//
// Example: SetDifference({"go":true, "java":true, "python":true},
//                        {"java":true}) → ["go", "python"]
func SetDifference(a, b map[string]bool) []string {
	// TODO: iterate a, skip keys in b, collect and sort the rest
	return nil
}

// Exercise 11:
// UniqueValues returns all unique values from a map, sorted.
//
// LESSON: Extracting values requires iteration (no built-in Values() method).
// Use a set (map[int]bool) to deduplicate, then sort for deterministic output.
// Map iteration order is RANDOM in Go — you must sort for consistent results.
//
// Example: UniqueValues({"a":3, "b":1, "c":3, "d":2}) → [1, 2, 3]
func UniqueValues(m map[string]int) []int {
	// TODO: collect values into a set, convert to slice, sort
	return nil
}

// Exercise 12:
// MapEqual reports whether two maps have the same keys and values.
//
// KEY INSIGHT: Go maps cannot be compared with == (compile error, except nil).
// You must compare manually: same length, then check each key-value pair.
//
// This is why maps.Equal exists in Go 1.21+ (golang.org/x/exp/maps before that).
// Implement it manually to understand the check.
//
// Example: MapEqual({"a":1, "b":2}, {"b":2, "a":1}) → true
func MapEqual(a, b map[string]int) bool {
	// TODO: check length, then iterate a and verify each key exists in b with same value
	return false
}

