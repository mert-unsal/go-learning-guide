// Standalone demo: Constructor Functions & the Stringer Interface
//
// Go has no constructors. Convention: use a New* function that returns
// a pointer to a fully initialized struct. This is the composition root
// pattern — all setup lives in one explicit place.
//
// Implementing fmt.Stringer (the String() method) gives your type custom
// output when passed to fmt.Println, fmt.Sprintf, etc. The fmt package
// checks for Stringer via a type assertion at runtime.
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

// Stack is a LIFO data structure backed by a slice.
type Stack struct {
	items []int
}

// NewStack is the constructor — returns a ready-to-use *Stack.
func NewStack() *Stack {
	return &Stack{items: make([]int, 0)}
}

func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	n := len(s.items)
	item := s.items[n-1]
	s.items = s.items[:n-1]
	return item, true
}

func (s *Stack) Peek() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack) Len() int {
	return len(s.items)
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// String implements fmt.Stringer for custom printing.
func (s *Stack) String() string {
	return fmt.Sprintf("Stack%v", s.items)
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Structs: Constructor & Stringer        %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Constructor Pattern ---
	fmt.Printf("%s▸ Constructor Pattern (NewStack)%s\n", cyan+bold, reset)

	s := NewStack()
	fmt.Printf("  s := NewStack() → %s%v%s (addr: %s%p%s)\n", magenta, s, reset, magenta, s, reset)
	fmt.Printf("  %s✔ Go has no constructors — convention: NewType() returns *Type%s\n", green, reset)
	fmt.Printf("  %s✔ Returns pointer so all methods can modify the same instance%s\n\n", green, reset)

	// --- Push Operations ---
	fmt.Printf("%s▸ Push Operations%s\n", cyan+bold, reset)

	s.Push(1)
	s.Push(2)
	s.Push(3)
	fmt.Printf("  Push(1), Push(2), Push(3) → %s%v%s\n", magenta, s, reset)
	fmt.Printf("  %s✔ Backed by slice — append() handles growth via runtime.growslice%s\n\n", green, reset)

	// --- Stringer Interface ---
	fmt.Printf("%s▸ fmt.Stringer Interface%s\n", cyan+bold, reset)

	fmt.Printf("  fmt.Println(s) outputs: %s%v%s\n", magenta, s, reset)
	fmt.Printf("  %s✔ String() method satisfies fmt.Stringer — called automatically by fmt functions%s\n", green, reset)
	fmt.Printf("  %s✔ The fmt package uses type assertion to check for Stringer at runtime%s\n\n", green, reset)

	// --- Pop ---
	fmt.Printf("%s▸ Pop (LIFO)%s\n", cyan+bold, reset)

	if top, ok := s.Pop(); ok {
		fmt.Printf("  Pop() → value: %s%d%s, ok: %s%v%s\n", magenta, top, reset, magenta, ok, reset)
	}
	fmt.Printf("  Stack after Pop: %s%v%s, Len: %s%d%s\n", magenta, s, reset, magenta, s.Len(), reset)
	fmt.Printf("  %s✔ Returns (value, bool) — the comma-ok pattern prevents panics on empty stack%s\n\n", green, reset)

	// --- Peek ---
	fmt.Printf("%s▸ Peek & IsEmpty%s\n", cyan+bold, reset)

	if top, ok := s.Peek(); ok {
		fmt.Printf("  Peek() → %s%d%s (top of stack, not removed)\n", magenta, top, reset)
	}
	fmt.Printf("  IsEmpty() = %s%v%s\n", magenta, s.IsEmpty(), reset)
	fmt.Printf("  %s✔ All methods use pointer receiver (*Stack) — consistent and can modify state%s\n", green, reset)
}
