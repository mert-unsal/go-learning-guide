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
	d := Dog{
		Animal: Animal{Name: "Rex"},
		Breed:  "Labrador",
	}

	// Access embedded fields directly (promoted)
	fmt.Println(d.Name)    // "Rex" — from Animal
	fmt.Println(d.Breed)   // "Labrador"
	fmt.Println(d.Speak()) // "Rex says: Woof!" — Dog's override

	// Still accessible via explicit path
	fmt.Println(d.Animal.Speak()) // "Rex makes a sound"

	c := Cat{Animal: Animal{Name: "Whiskers"}, Indoor: true}
	fmt.Println(c.Speak()) // "Whiskers makes a sound" — from Animal (no override)
	fmt.Println(c.Indoor)
}
