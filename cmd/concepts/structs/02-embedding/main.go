// Standalone demo: Struct Embedding (Composition over Inheritance)
//
// Go doesn't have inheritance. Instead, use EMBEDDING to compose types.
// An embedded struct's fields and methods are PROMOTED to the outer struct.
//
// Under the hood: embedding is syntactic sugar. The compiler generates
// forwarding methods and field accessors. There is NO vtable or dynamic
// dispatch — method resolution is fully static at compile time. If Dog
// embeds Animal and overrides Speak(), calling d.Speak() goes directly
// to Dog.Speak — there's no "super" chain. You can still reach the
// original via d.Animal.Speak().
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

type Animal struct {
	Name string
}

func (a Animal) Speak() string {
	return a.Name + " makes a sound"
}

type Dog struct {
	Animal // embedded (anonymous field) — NOT Animal Animal
	Breed  string
}

func (d Dog) Speak() string { // Dog can OVERRIDE Animal's method
	return d.Name + " says: Woof!"
}

type Cat struct {
	Animal
	Indoor bool
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Structs: Embedding (Composition)       %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	d := Dog{
		Animal: Animal{Name: "Rex"},
		Breed:  "Labrador",
	}

	// --- Field Promotion ---
	fmt.Printf("%s▸ Field Promotion (embedded fields accessible directly)%s\n", cyan+bold, reset)

	// Access embedded fields directly (promoted)
	fmt.Printf("  d.Name  = %s%q%s (promoted from Animal — no d.Animal.Name needed)\n", magenta, d.Name, reset)
	fmt.Printf("  d.Breed = %s%q%s (Dog's own field)\n", magenta, d.Breed, reset)
	fmt.Printf("  %s✔ Embedded struct fields are \"promoted\" — accessible as if they belong to outer struct%s\n\n", green, reset)

	// --- Method Override ---
	fmt.Printf("%s▸ Method Override%s\n", cyan+bold, reset)

	fmt.Printf("  d.Speak()        = %s%q%s\n", magenta, d.Speak(), reset)
	fmt.Printf("  %s✔ Dog.Speak() overrides Animal.Speak() — resolved statically at compile time%s\n", green, reset)
	fmt.Printf("  %s✔ No vtable, no dynamic dispatch — the compiler picks Dog.Speak directly%s\n\n", green, reset)

	// --- Explicit Access ---
	fmt.Printf("%s▸ Explicit Access to Embedded Type%s\n", cyan+bold, reset)

	// Still accessible via explicit path
	fmt.Printf("  d.Animal.Speak() = %s%q%s\n", magenta, d.Animal.Speak(), reset)
	fmt.Printf("  %s✔ Original method still accessible via d.Animal.Speak()%s\n", green, reset)
	fmt.Printf("  %s⚠ This is NOT inheritance — Go has no \"super\" keyword. Use explicit embedding path%s\n\n", yellow, reset)

	// --- Cat (no override) ---
	fmt.Printf("%s▸ Cat — No Override (uses Animal's method)%s\n", cyan+bold, reset)

	c := Cat{Animal: Animal{Name: "Whiskers"}, Indoor: true}
	fmt.Printf("  c.Speak()  = %s%q%s\n", magenta, c.Speak(), reset)
	fmt.Printf("  c.Indoor   = %s%v%s\n", magenta, c.Indoor, reset)
	fmt.Printf("  %s✔ Cat doesn't override Speak() — the promoted Animal.Speak() is used%s\n\n", green, reset)

	// --- Key Takeaways ---
	fmt.Printf("%s▸ Composition vs Inheritance%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Embedding = composition, not inheritance. No \"is-a\" relationship%s\n", green, reset)
	fmt.Printf("  %s✔ Compiler generates forwarding methods — zero runtime overhead%s\n", green, reset)
	fmt.Printf("  %s⚠ If two embedded types have the same method, you get ambiguity — must disambiguate%s\n", yellow, reset)
}
