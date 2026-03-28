// Standalone demo: Closures
//
// A closure is a function that CAPTURES variables from its surrounding scope.
// The captured variable is shared between the function and its environment —
// it lives on the heap (the compiler's escape analysis moves it there).
//
// Key insight: closures close over VARIABLES, not VALUES. This is the root
// of the classic loop-capture gotcha. Since Go 1.22 the loop variable is
// per-iteration by default, but in earlier versions you need the shadow trick.
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

// makeCounter returns a closure that increments a counter.
// Each call to makeCounter allocates a fresh `count` on the heap,
// so every returned closure has its own independent state.
func makeCounter() func() int {
	count := 0 // this variable is CAPTURED by the inner function
	return func() int {
		count++ // modifies the captured variable
		return count
	}
}

// makeAdder returns a closure that adds a fixed value.
// x is captured from makeAdder's parameter scope.
func makeAdder(x int) func(int) int {
	return func(y int) int {
		return x + y // x is captured from makeAdder's scope
	}
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Closures in Go                         %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Stateful closure: counter ---
	fmt.Printf("%s▸ Stateful closure: makeCounter()%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Each call allocates a fresh 'count' on the heap — independent state%s\n", green, reset)

	counter := makeCounter()
	fmt.Printf("  counter() = %s%d%s  ← first call: count goes 0→1\n", magenta, counter(), reset)
	fmt.Printf("  counter() = %s%d%s  ← same closure, same captured variable, count goes 1→2\n", magenta, counter(), reset)
	fmt.Printf("  counter() = %s%d%s  ← count goes 2→3\n", magenta, counter(), reset)

	fmt.Printf("\n  %s✔ Creating a second counter — completely independent state%s\n", green, reset)
	counter2 := makeCounter()
	fmt.Printf("  counter2() = %s%d%s  ← new closure, new 'count' variable on heap\n", magenta, counter2(), reset)

	// --- Closure with captured parameter ---
	fmt.Printf("\n%s▸ Closure capturing a parameter: makeAdder(5)%s\n", cyan+bold, reset)
	add5 := makeAdder(5)
	fmt.Printf("  add5(3)  = %s%d%s  ← x=5 captured from makeAdder's parameter scope\n", magenta, add5(3), reset)
	fmt.Printf("  add5(10) = %s%d%s  ← same captured x=5, new y=10\n", magenta, add5(10), reset)
	fmt.Printf("  %s✔ 'x' lives on the heap because escape analysis detects the closure captures it%s\n", green, reset)

	// --- Classic loop-capture gotcha ---
	fmt.Printf("\n%s▸ Loop-capture gotcha (with fix)%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ Closures close over VARIABLES, not VALUES%s\n", yellow, reset)
	fmt.Printf("  %s⚠ Without shadowing, all closures would share the same 'i' → all print final value%s\n", yellow, reset)
	fmt.Printf("  %s✔ Fix: 'i := i' creates a new variable per iteration (Go 1.22+ does this automatically)%s\n", green, reset)

	funcs := make([]func(), 3)
	for i := 0; i < 3; i++ {
		i := i // shadow i with a new variable per iteration (FIX!)
		funcs[i] = func() {
			fmt.Printf("  funcs[%d]() = %s%d%s\n", i, magenta, i, reset)
		}
	}
	funcs[0]()
	funcs[1]()
	funcs[2]()

	// --- Under the hood ---
	fmt.Printf("\n%s▸ Under the hood%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Captured vars escape to heap (confirmed with: go build -gcflags='-m')%s\n", green, reset)
	fmt.Printf("  %s✔ Each closure is a runtime.funcval struct pointing to the captured variables%s\n", green, reset)
	fmt.Printf("  %s✔ Multiple closures from the same scope share the SAME heap variable%s\n", green, reset)
}
