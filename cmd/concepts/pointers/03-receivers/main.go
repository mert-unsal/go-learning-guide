// Standalone demo: Value vs Pointer Receivers
//
// Methods can have either value receivers (T) or pointer receivers (*T).
// Use pointer receivers when:
//   a) The method needs to MODIFY the receiver
//   b) The struct is large (avoids copying)
//
// Under the hood: Go auto-inserts &/dereference as needed. If you call
// c.Increment() on a value, the compiler rewrites it to (&c).Increment().
// This only works when the value is addressable (variables, struct fields,
// slice elements — NOT map values or function returns).
//
// Rule: if ANY method has a pointer receiver, ALL methods should have
// pointer receivers for consistency. Mixing confuses interface satisfaction.
//
// Run:  go run .
package main

import "fmt"

// Counter is a simple stateful type that demonstrates the difference
// between value and pointer receivers.
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

func main() {
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
