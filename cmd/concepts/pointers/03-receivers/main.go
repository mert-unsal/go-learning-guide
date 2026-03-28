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
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Pointers: Value vs Pointer Receivers   %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	c := Counter{count: 0}
	fmt.Printf("%s▸ Pointer Receiver — Can Modify%s\n", cyan+bold, reset)
	fmt.Printf("  Counter addr: %s%p%s, count: %s%d%s\n", magenta, &c, reset, magenta, c.count, reset)

	c.Increment() // Go automatically takes address: (&c).Increment()
	c.Increment()
	c.Increment()
	fmt.Printf("  After 3x Increment(): count = %s%d%s\n", magenta, c.ValueGet(), reset)
	fmt.Printf("  %s✔ Pointer receiver (*Counter) gets &c — modifies the original struct%s\n", green, reset)
	fmt.Printf("  %s✔ Go auto-inserts & : c.Increment() becomes (&c).Increment()%s\n\n", green, reset)

	// --- Value Receiver ---
	fmt.Printf("%s▸ Value Receiver — Gets a Copy%s\n", cyan+bold, reset)
	fmt.Printf("  c.ValueGet() = %s%d%s\n", magenta, c.ValueGet(), reset)
	fmt.Printf("  %s✔ ValueGet() receives a copy of Counter — cannot modify original%s\n", green, reset)
	fmt.Printf("  %s✔ Safe for read-only access; the copy is cheap for small structs%s\n\n", green, reset)

	// --- Reset ---
	fmt.Printf("%s▸ Reset (Pointer Receiver)%s\n", cyan+bold, reset)
	c.Reset()
	fmt.Printf("  After Reset(): count = %s%d%s\n", magenta, c.ValueGet(), reset)
	fmt.Printf("  %s✔ Pointer receiver sets c.count = 0 on the original struct%s\n\n", green, reset)

	// --- Method Set Rules ---
	fmt.Printf("%s▸ Method Set Rules%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Value (T) method set: only value receivers%s\n", green, reset)
	fmt.Printf("  %s✔ Pointer (*T) method set: value + pointer receivers%s\n", green, reset)
	fmt.Printf("  %s⚠ If ANY method uses pointer receiver, ALL should for consistency%s\n", yellow, reset)
	fmt.Printf("  %s⚠ A value stored in an interface is NOT addressable — pointer methods won't be in its method set%s\n", yellow, reset)

	// Rule: if ANY method has a pointer receiver,
	// ALL methods should have pointer receivers (consistency).
}
