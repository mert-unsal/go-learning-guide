// Package pointers covers Go pointers: address-of, dereferencing,
// pointer receivers, and when to use pointers vs values.
package pointers

import "fmt"

// ============================================================
// 1. BASICS: address-of & dereferencing
// ============================================================
// &x  → "address of x" → returns a *T (pointer to T)
// *p  → "dereference p" → returns the value at the address
// new(T) → allocates zeroed T, returns *T

func DemonstrateBasics() {
	x := 42
	p := &x // p is of type *int; it POINTS TO x

	fmt.Println("x =", x)   // 42
	fmt.Println("p =", p)   // memory address like 0xc0000b4010
	fmt.Println("*p =", *p) // 42 (dereference: value at the address)

	*p = 100       // modify the value through the pointer
	fmt.Println(x) // 100 — x was changed!

	// new() allocates and returns a pointer to a zeroed value
	q := new(int) // *int, zero value is 0
	*q = 55
	fmt.Println(*q) // 55

	// nil pointer: a pointer with no address
	var r *int     // nil
	fmt.Println(r) // <nil>
	// *r = 5        // PANIC: nil pointer dereference — never do this!

	// Safe nil check before dereferencing
	if r != nil {
		fmt.Println(*r)
	} else {
		fmt.Println("r is nil, can't dereference")
	}
}

// ============================================================
// 2. WHY POINTERS? Pass by reference
// ============================================================

// WITHOUT pointer: modifies a COPY, original unchanged
func incrementValue(n int) {
	n++ // modifies local copy only
}

// WITH pointer: modifies the ORIGINAL
func incrementPointer(n *int) {
	*n++ // dereferences and modifies the original
}

func DemonstratePassByReference() {
	x := 10
	incrementValue(x)
	fmt.Println("After incrementValue:", x) // still 10

	incrementPointer(&x)
	fmt.Println("After incrementPointer:", x) // 11

	// Slices, maps, and channels are already reference types —
	// you don't need pointers for them in most cases.
}

// ============================================================
// 3. POINTER RECEIVERS ON METHODS
// ============================================================
// Methods can have either value receivers (T) or pointer receivers (*T).
// Use pointer receivers when:
//   a) The method needs to MODIFY the receiver
//   b) The struct is large (avoids copying)

type Counter struct {
	count int
}

// Value receiver — gets a COPY; cannot modify original
func (c Counter) ValueGet() int {
	return c.count
}

// Pointer receiver — gets a POINTER; CAN modify original
func (c *Counter) Increment() {
	c.count++
}

func (c *Counter) Reset() {
	c.count = 0
}

func DemonstrateReceivers() {
	c := Counter{count: 0}
	c.Increment() // Go automatically takes address: (&c).Increment()
	c.Increment()
	c.Increment()
	fmt.Println("Count:", c.ValueGet()) // 3
	c.Reset()
	fmt.Println("After reset:", c.ValueGet()) // 0

	// Rule: if ANY method has a pointer receiver,
	// ALL methods should have pointer receivers (consistency).
}

// ============================================================
// 4. POINTERS TO STRUCTS
// ============================================================

type Person struct {
	Name string
	Age  int
}

func birthday(p *Person) {
	p.Age++ // Go auto-dereferences: (*p).Age++
}

func DemonstrateStructPointers() {
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

// ============================================================
// 5. WHEN TO USE POINTERS (Rules of Thumb)
// ============================================================
// USE a pointer (*T) when:
//   1. You need to modify the value
//   2. The struct is large (performance)
//   3. You need to represent "absence" (nil)
//   4. Consistency: if one method uses *T, all should
//
// USE a value (T) when:
//   1. The type is small (int, bool, small struct)
//   2. You don't need to modify it
//   3. You want immutability guarantees
//   4. The type is a map, slice, or channel (already reference types)

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Pointer Basics ===")
	DemonstrateBasics()
	fmt.Println("\n=== Pass by Reference ===")
	DemonstratePassByReference()
	fmt.Println("\n=== Pointer Receivers ===")
	DemonstrateReceivers()
	fmt.Println("\n=== Struct Pointers ===")
	DemonstrateStructPointers()
}
