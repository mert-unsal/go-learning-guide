// Package main demonstrates Go's implicit interface satisfaction.
//
// ============================================================
// THE MENTAL MODEL — READ THIS FIRST
// ============================================================
//
// In OOP languages (Java, C#) you think:
//   "These types ARE a Shape, so I'll declare them as implementing Shape."
//   → The PRODUCER declares the relationship upfront.
//   → Types are grouped by what they ARE (taxonomy/hierarchy).
//
// In Go you think:
//   "My function needs something it can call Area() on. I'll define that contract."
//   → The CONSUMER defines the interface, right where it needs it.
//   → Types are grouped by what they CAN DO (behavior), not what they are.
//
// This is the fundamental difference. Go has no class, no extends,
// no implements. A type satisfies an interface just by having the methods —
// it doesn't know about the interface and doesn't need to.
//
// KEY RULES:
//   1. Keep interfaces small (1-3 methods). The smaller, the more useful.
//   2. Define interfaces where they are CONSUMED, not where types are defined.
//   3. Don't create interfaces speculatively — create them when you need abstraction.
//   4. Accept interfaces, return concrete types (Rob Pike's guideline).
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

// ============================================================
// 1. IMPLICIT SATISFACTION — NO 'implements' KEYWORD
// ============================================================
// The type and the interface have ZERO connection at the type definition site.
// The compiler checks at the point of USE, not at the point of definition.

// Stringer is defined in the fmt package. We define our own here to illustrate.
// Any type with a String() string method satisfies it — it never needs to say so.
type Stringer interface {
	String() string
}

type User struct {
	Name string
	Age  int
}

// User satisfies Stringer. User doesn't know Stringer exists.
// Stringer doesn't know User exists. The compiler connects them at use time.
func (u User) String() string {
	return fmt.Sprintf("%s (age %d)", u.Name, u.Age)
}

type Point struct {
	X, Y float64
}

// Point also satisfies Stringer — completely unrelated type, same behavior.
func (p Point) String() string {
	return fmt.Sprintf("(%.1f, %.1f)", p.X, p.Y)
}

// printIt was written knowing nothing about User or Point.
// It works with both because they share the behavior, not a type hierarchy.
func printIt(s Stringer) {
	fmt.Println(s.String())
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Implicit Interface Satisfaction         %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ No 'implements' keyword — the compiler does structural matching%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ User has String() string → satisfies Stringer automatically%s\n", green, reset)
	fmt.Printf("  %s✔ Point has String() string → also satisfies Stringer%s\n", green, reset)
	fmt.Printf("  %s✔ Neither type knows Stringer exists — zero coupling at definition site%s\n\n", green, reset)

	u := User{Name: "Alice", Age: 30}
	p := Point{X: 3.5, Y: -1.2}

	fmt.Printf("%s▸ Passing unrelated types to printIt(s Stringer)%s\n", cyan+bold, reset)
	fmt.Printf("  User  → ")
	printIt(u)
	fmt.Printf("  Point → ")
	printIt(p)
	fmt.Println()

	// Compile-time proof: assign to interface variable
	fmt.Printf("%s▸ Runtime interface value (iface) internals%s\n", cyan+bold, reset)
	var s Stringer = u
	fmt.Printf("  var s Stringer = u  →  s is %s(type=User, data=0x...)%s\n", magenta, reset)
	fmt.Printf("  s.String() = %s%q%s\n", magenta, s.String(), reset)

	s = p
	fmt.Printf("  s = p               →  s is %s(type=Point, data=0x...)%s\n", magenta, reset)
	fmt.Printf("  s.String() = %s%q%s\n\n", magenta, s.String(), reset)

	fmt.Printf("  %s⚠ Under the hood: non-empty interface = runtime.iface{tab *itab, data unsafe.Pointer}%s\n", yellow, reset)
	fmt.Printf("  %s⚠ The itab holds method pointers — cached globally per (interface, concrete) pair%s\n", yellow, reset)
	fmt.Printf("  %s⚠ In Java/C# you'd write 'class User implements Stringer' — Go needs nothing%s\n", yellow, reset)
}
