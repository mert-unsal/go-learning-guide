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

type Person struct {
	Name string
	Age  int
}

func birthday(p *Person) {
	p.Age++ // Go auto-dereferences: (*p).Age++
}

func main() {
	alice := Person{Name: "Alice", Age: 30}
	birthday(&alice)
	fmt.Println(alice.Age) // 31

	// Creating a struct with new()
	bob := new(Person)
	bob.Name = "Bob"
	bob.Age = 25
	birthday(bob)        // bob is already a *Person
	fmt.Println(bob.Age) // 26

	// Pointer literal — common pattern
	charlie := &Person{Name: "Charlie", Age: 20}
	fmt.Println(charlie) // &{Charlie 20}
}
