// Package functions covers: function signatures, multiple returns,
// named returns, variadic functions, closures, and function types.
package functions

import "fmt"

// ============================================================
// 1. BASIC FUNCTION SYNTAX
// ============================================================
// func name(param type, param type) returnType { ... }

func add(a, b int) int { // when params share a type, only last needs it
	return a + b
}

// ============================================================
// 2. MULTIPLE RETURN VALUES
// ============================================================
// Go functions can return multiple values. This is the idiomatic
// way to return a result AND an error.

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil // nil means "no error"
}

// Named return values: declared in the signature
// They act like local variables, and 'return' alone returns them
func minMax(arr []int) (min, max int) {
	min, max = arr[0], arr[0]
	for _, v := range arr[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return // "naked return" - returns named values min and max
}

// ============================================================
// 3. VARIADIC FUNCTIONS
// ============================================================
// ...Type means the function accepts any number of that type
// The variadic parameter is treated as a SLICE inside the function

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func DemonstrateVariadic() {
	fmt.Println(sum(1, 2, 3))       // 6
	fmt.Println(sum(1, 2, 3, 4, 5)) // 15

	// Spread a slice with ...
	nums := []int{10, 20, 30}
	fmt.Println(sum(nums...)) // 60
}

// ============================================================
// 4. FUNCTIONS AS FIRST-CLASS VALUES
// ============================================================
// Functions are values in Go. They can be:
// - Assigned to variables
// - Passed as arguments
// - Returned from functions

func apply(nums []int, fn func(int) int) []int {
	result := make([]int, len(nums))
	for i, v := range nums {
		result[i] = fn(v)
	}
	return result
}

func DemonstrateFunctionValues() {
	nums := []int{1, 2, 3, 4, 5}

	// Pass an anonymous function (lambda)
	doubled := apply(nums, func(n int) int {
		return n * 2
	})
	fmt.Println(doubled) // [2 4 6 8 10]

	// Assign a function to a variable
	square := func(n int) int { return n * n }
	squared := apply(nums, square)
	fmt.Println(squared) // [1 4 9 16 25]
}

// ============================================================
// 5. CLOSURES
// ============================================================
// A closure is a function that CAPTURES variables from its surrounding scope.
// The captured variable is shared between the function and its environment.

// makeCounter returns a closure that increments a counter
func makeCounter() func() int {
	count := 0 // this variable is CAPTURED by the inner function
	return func() int {
		count++ // modifies the captured variable
		return count
	}
}

// makeAdder returns a closure that adds a fixed value
func makeAdder(x int) func(int) int {
	return func(y int) int {
		return x + y // x is captured from makeAdder's scope
	}
}

func DemonstrateClosures() {
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

// ============================================================
// 6. INIT FUNCTION
// ============================================================
// init() is a special function called automatically when the package
// is initialized. Rules:
// - No parameters, no return values
// - Can't be called manually
// - Multiple init() functions allowed per file/package
// - Called after all variable declarations are evaluated
// - Called before main()

var configLoaded bool

func init() {
	// This runs automatically when the package is first used
	configLoaded = true
	// Common uses: initialize package-level variables, register things
}

// ============================================================
// 7. DEFER WITH FUNCTIONS (Advanced)
// ============================================================

// measureTime demonstrates using defer to measure execution time
func measureTime(name string) func() {
	fmt.Printf("Starting %s\n", name)
	return func() {
		fmt.Printf("Done %s\n", name)
	}
}

func someExpensiveOperation() {
	defer measureTime("expensive operation")()
	// ^ Note: the () at the end calls the returned function at defer time
	// Do work here...
}

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Variadic ===")
	DemonstrateVariadic()
	fmt.Println("\n=== Function Values ===")
	DemonstrateFunctionValues()
	fmt.Println("\n=== Closures ===")
	DemonstrateClosures()
	fmt.Println("Config loaded by init():", configLoaded)
}
