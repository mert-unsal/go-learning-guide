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
	counter := makeCounter()
	fmt.Println(counter()) // 1
	fmt.Println(counter()) // 2
	fmt.Println(counter()) // 3

	// Each call to makeCounter creates a NEW counter
	counter2 := makeCounter()
	fmt.Println(counter2()) // 1 (independent from counter)

	add5 := makeAdder(5)
	fmt.Println(add5(3))  // 8
	fmt.Println(add5(10)) // 15

	// Classic closure gotcha in loops:
	// WRONG way — all closures capture the same variable i
	funcs := make([]func(), 3)
	for i := 0; i < 3; i++ {
		i := i // shadow i with a new variable per iteration (FIX!)
		funcs[i] = func() {
			fmt.Println(i)
		}
	}
	funcs[0]() // 0
	funcs[1]() // 1
	funcs[2]() // 2
}
