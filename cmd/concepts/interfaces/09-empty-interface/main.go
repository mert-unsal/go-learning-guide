// Package main demonstrates the empty interface (any / interface{}).
//
// ============================================================
// 8. THE EMPTY INTERFACE — any
// ============================================================
// interface{} (aliased as 'any' since Go 1.18) has zero methods.
// Every type satisfies it. Use it only when you truly cannot know
// the type at compile time (e.g., JSON decoding, fmt internals).
// It throws away compile-time type safety — use sparingly.
//
// Under the hood:
//   - any is represented as runtime.eface = {_type *_type, data unsafe.Pointer}
//   - When you assign a value to any, the compiler boxes it:
//     the value is copied to the heap, and eface.data points to the copy.
//   - Small values (≤ pointer size) may be stored directly in the data word
//     as a compiler optimization.
//   - Every assignment to any is a potential heap allocation (escape analysis).
//
// When to use any:
//   ✅ JSON/YAML/TOML decoding into unknown structures
//   ✅ Generic containers before Go 1.18 generics
//   ✅ fmt.Println-style variadic functions
//   ✅ Bridging with untyped external data (e.g., plugin systems)
//
// When NOT to use any:
//   ❌ When you know the type — use the concrete type
//   ❌ When a small interface would work — define the behavior
//   ❌ As a lazy substitute for generics (Go 1.18+)
//   ❌ For dependency injection — use specific interfaces
//
// ============================================================
// SUMMARY: Go interfaces vs OOP interfaces
// ============================================================
//
//  OOP                            Go
//  ─────────────────────────────────────────────────────────
//  declared by the PRODUCER       defined by the CONSUMER
//  "class X implements I"         type just has the methods
//  big, upfront hierarchies       small, composed on demand
//  groups by WHAT THINGS ARE      groups by WHAT THINGS DO
//  coupling at definition time    coupling resolved at compile time, at use site
//  abstract classes, inheritance  no inheritance, only behavior contracts
//
// The practical result: in Go you can write an interface for a type
// defined in a package you don't own, and it will just work.
// No forking, no wrapping, no adapter boilerplate.
package main

import "fmt"

func printAnything(v any) {
	fmt.Printf("type: %-12T value: %v\n", v, v)
}

func main() {
	printAnything(42)
	printAnything("hello")
	printAnything(true)
	printAnything([]int{1, 2, 3})
	printAnything(nil)
}
