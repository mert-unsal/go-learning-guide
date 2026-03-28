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
	// User and Point are completely different types.
	// They were never declared to implement anything.
	// They work here purely because they have the right method.
	printIt(User{Name: "Alice", Age: 30})
	printIt(Point{X: 3.5, Y: -1.2})
}
