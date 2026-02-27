// Package control_flow covers: if/else, switch, for loops, defer, and goto.
package control_flow

import "fmt"

// ============================================================
// 1. IF / ELSE
// ============================================================
func DemonstrateIfElse() {
	x := 42

	// Basic if/else
	if x > 0 {
		fmt.Println("positive")
	} else if x < 0 {
		fmt.Println("negative")
	} else {
		fmt.Println("zero")
	}

	// KEY FEATURE: if with initialization statement
	// The variable 'n' is scoped ONLY to this if/else block
	if n := 10; n%2 == 0 {
		fmt.Println(n, "is even")
	} else {
		fmt.Println(n, "is odd")
	}
	// fmt.Println(n) // COMPILE ERROR: n is not accessible here
}

// ============================================================
// 2. SWITCH
// ============================================================
func DemonstrateSwitch() {
	// Basic switch — no 'break' needed, each case breaks automatically
	day := "Monday"
	switch day {
	case "Saturday", "Sunday": // multiple values per case
		fmt.Println("Weekend!")
	case "Monday":
		fmt.Println("Start of the work week")
	default:
		fmt.Println("Weekday")
	}

	// Switch with no expression (acts like if/else chain)
	x := 15
	switch {
	case x < 0:
		fmt.Println("negative")
	case x == 0:
		fmt.Println("zero")
	case x > 0 && x < 10:
		fmt.Println("small positive")
	default:
		fmt.Println("large positive")
	}

	// Switch with initializer
	switch n := 42; {
	case n < 0:
		fmt.Println("negative")
	case n < 100:
		fmt.Println("less than 100")
	default:
		fmt.Println("100 or more")
	}

	// fallthrough: explicitly continue to next case (rare in practice)
	switch 2 {
	case 1:
		fmt.Println("one")
		fallthrough
	case 2:
		fmt.Println("two")
		fallthrough // executes next case even if condition doesn't match
	case 3:
		fmt.Println("three") // this WILL print
	case 4:
		fmt.Println("four") // this will NOT print (fallthrough stops)
	}
}

// ============================================================
// 3. FOR LOOPS
// ============================================================
// Go has ONLY ONE loop keyword: 'for'
// It replaces while, do-while, and for from other languages.

func DemonstrateForLoops() {
	// Classic C-style for loop
	for i := 0; i < 5; i++ {
		fmt.Print(i, " ") // 0 1 2 3 4
	}
	fmt.Println()

	// While-style: just condition (no init or post)
	n := 1
	for n < 100 {
		n *= 2
	}
	fmt.Println("n =", n) // 128

	// Infinite loop (exit with 'break')
	count := 0
	for {
		count++
		if count >= 3 {
			break
		}
	}
	fmt.Println("count =", count)

	// continue: skip to next iteration
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue // skip even numbers
		}
		fmt.Print(i, " ") // 1 3 5 7 9
	}
	fmt.Println()

	// range over slice
	nums := []int{10, 20, 30, 40}
	for index, value := range nums {
		fmt.Printf("nums[%d] = %d\n", index, value)
	}

	// range over string (iterates over RUNES, not bytes)
	for i, ch := range "Go!" {
		fmt.Printf("index=%d, char=%c\n", i, ch)
	}

	// range over map (order is RANDOM)
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for key, val := range m {
		fmt.Printf("%s: %d\n", key, val)
	}

	// Discard index with _
	sum := 0
	for _, v := range []int{1, 2, 3, 4, 5} {
		sum += v
	}
	fmt.Println("sum =", sum)

	// Go 1.22+: range over integer
	for i := range 5 { // 0, 1, 2, 3, 4
		fmt.Print(i, " ")
	}
	fmt.Println()

	// Labeled break/continue (for nested loops)
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outer // breaks the OUTER loop
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
}

// ============================================================
// 4. DEFER
// ============================================================
// defer delays execution of a function until the surrounding
// function RETURNS. Very useful for cleanup (closing files, etc.)

func DemonstrateDefer() {
	fmt.Println("start")
	defer fmt.Println("deferred 1") // runs last
	defer fmt.Println("deferred 2") // runs second-to-last
	defer fmt.Println("deferred 3") // runs first among defers
	fmt.Println("end")
	// Output:
	// start
	// end
	// deferred 3  ← LIFO order (stack)
	// deferred 2
	// deferred 1
}

// Typical defer usage: resource cleanup
func readFile(filename string) {
	// In real code:
	// f, err := os.Open(filename)
	// if err != nil { ... }
	// defer f.Close() ← guaranteed to run when function exits

	fmt.Printf("Opening %s\n", filename)
	defer fmt.Printf("Closing %s\n", filename) // runs when function returns
	fmt.Printf("Reading %s\n", filename)
}

// Defer with loop — IMPORTANT GOTCHA
func deferInLoop() {
	// BAD: all defers stack up and run when function exits, not each iteration
	for i := 0; i < 3; i++ {
		defer fmt.Println("defer in loop:", i) // Don't do this for resource cleanup!
	}
}

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== If/Else ===")
	DemonstrateIfElse()
	fmt.Println("\n=== Switch ===")
	DemonstrateSwitch()
	fmt.Println("\n=== For Loops ===")
	DemonstrateForLoops()
	fmt.Println("\n=== Defer ===")
	DemonstrateDefer()
}
