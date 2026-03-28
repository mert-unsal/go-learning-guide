// Standalone demo: Pointers to Structs
//
// Go auto-dereferences struct pointers: p.Field is syntactic sugar for
// (*p).Field. Combined with the &Type{} literal pattern, this makes
// pointer-based APIs feel natural.
//
// ============================================================
// WHEN TO USE POINTERS (Rules of Thumb)
// ============================================================
// USE a pointer (*T) when:
//   1. You need to modify the value
//   2. The struct is large (performance — avoids copying)
//   3. You need to represent "absence" (nil)
//   4. Consistency: if one method uses *T, all should
//
// USE a value (T) when:
//   1. The type is small (int, bool, small struct)
//   2. You don't need to modify it
//   3. You want immutability guarantees
//   4. The type is a map, slice, or channel (already reference types)
//
// Run:  go run .
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

type Person struct {
	Name string
	Age  int
}

func birthday(p *Person) {
	p.Age++ // Go auto-dereferences: (*p).Age++
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Pointers: Struct Pointers              %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Auto-Dereference ---
	fmt.Printf("%s▸ Auto-Dereference (p.Field vs (*p).Field)%s\n", cyan+bold, reset)

	alice := Person{Name: "Alice", Age: 30}
	fmt.Printf("  alice (value): %s%+v%s at addr %s%p%s\n", magenta, alice, reset, magenta, &alice, reset)

	birthday(&alice)
	fmt.Printf("  After birthday(&alice): Age = %s%d%s\n", magenta, alice.Age, reset)
	fmt.Printf("  %s✔ p.Age is syntactic sugar for (*p).Age — Go auto-dereferences struct pointers%s\n\n", green, reset)

	// --- new() ---
	fmt.Printf("%s▸ Creating with new()%s\n", cyan+bold, reset)

	// Creating a struct with new()
	bob := new(Person)
	fmt.Printf("  bob = new(Person) → addr: %s%p%s, value: %s%+v%s\n", magenta, bob, reset, magenta, *bob, reset)
	fmt.Printf("  %s✔ new(Person) returns *Person with all fields zeroed (\"\" and 0)%s\n", green, reset)

	bob.Name = "Bob"
	bob.Age = 25
	birthday(bob) // bob is already a *Person
	fmt.Printf("  After setting fields + birthday: %s%+v%s\n", magenta, *bob, reset)
	fmt.Printf("  %s✔ bob is already *Person — no need to pass &bob%s\n\n", green, reset)

	// --- Pointer Literal ---
	fmt.Printf("%s▸ Pointer Literal Pattern (&Type{})%s\n", cyan+bold, reset)

	// Pointer literal — common pattern
	charlie := &Person{Name: "Charlie", Age: 20}
	fmt.Printf("  charlie := &Person{...} → addr: %s%p%s, value: %s%+v%s\n", magenta, charlie, reset, magenta, *charlie, reset)
	fmt.Printf("  %s✔ &Type{} is the idiomatic way to create a pointer to a struct literal%s\n", green, reset)
	fmt.Printf("  %s✔ Equivalent to: tmp := Type{...}; ptr := &tmp — but more concise%s\n\n", green, reset)

	// --- When to Use Pointers ---
	fmt.Printf("%s▸ When to Use Pointers vs Values%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Use *T when: method modifies receiver, struct is large, need nil semantics%s\n", green, reset)
	fmt.Printf("  %s✔ Use T when: struct is small, read-only, want immutability guarantees%s\n", green, reset)
	fmt.Printf("  %s⚠ Maps, slices, channels are already reference types — rarely need pointers to them%s\n", yellow, reset)
}
