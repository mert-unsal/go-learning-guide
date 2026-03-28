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

func printAnything(v any) {
	typeName := fmt.Sprintf("%T", v)
	fmt.Printf("  type: %s%-12s%s  value: %s%v%s\n", magenta, typeName, reset, magenta, v, reset)
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Empty Interface (any) — runtime.eface   %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ any (interface{}) has ZERO methods — every type satisfies it%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Represented as runtime.eface = {_type *_type, data unsafe.Pointer}%s\n", green, reset)
	fmt.Printf("  %s✔ Only 2 words (16 bytes on 64-bit) — simpler than iface (no itab)%s\n", green, reset)
	fmt.Printf("  %s⚠ Boxing cost: assigning to any may copy value to heap (escape analysis)%s\n\n", yellow, reset)

	fmt.Printf("%s▸ Passing different types to printAnything(v any)%s\n", cyan+bold, reset)
	printAnything(42)
	printAnything("hello")
	printAnything(true)
	printAnything([]int{1, 2, 3})
	printAnything(nil)

	fmt.Printf("\n%s▸ Boxing/unboxing — what happens under the hood%s\n", cyan+bold, reset)
	var x any = 42
	fmt.Printf("  var x any = 42  → eface = %s(_type=int, data→42)%s\n", magenta, reset)
	fmt.Printf("  %s⚠ 42 is boxed: copied to heap, eface.data points to the copy%s\n", yellow, reset)

	// Unboxing via type assertion
	n, ok := x.(int)
	fmt.Printf("  n, ok := x.(int) → n=%s%d%s, ok=%s%v%s  (unboxing — itab type check)\n", magenta, n, reset, magenta, ok, reset)

	_, ok = x.(string)
	fmt.Printf("  _, ok := x.(string) → ok=%s%v%s  (type mismatch — no panic with comma-ok)\n\n", magenta, ok, reset)

	fmt.Printf("%s▸ When to use any — and when NOT to%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ JSON decoding into unknown structures%s\n", green, reset)
	fmt.Printf("  %s✔ fmt.Println-style variadic functions%s\n", green, reset)
	fmt.Printf("  %s⚠ NOT when you know the type — use concrete types%s\n", yellow, reset)
	fmt.Printf("  %s⚠ NOT as a lazy substitute for generics (Go 1.18+)%s\n", yellow, reset)
	fmt.Printf("  %s⚠ \"interface{} says nothing\" — Go Proverbs%s\n", yellow, reset)
}
