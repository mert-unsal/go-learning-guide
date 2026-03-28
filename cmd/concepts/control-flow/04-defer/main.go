// Defer in Go — delays execution until the surrounding function returns.
// Very useful for cleanup (closing files, unlocking mutexes, etc.)
//
// KEY RULES:
//  1. Deferred calls run in LIFO order (last deferred = first to run)
//  2. Arguments to deferred functions are evaluated IMMEDIATELY (at defer time)
//  3. Deferred functions CAN read and modify named return values
//  4. Do NOT defer inside loops for resource cleanup — they stack up
package main

import "fmt"

// readFile — typical defer usage: resource cleanup
func readFile(filename string) {
	// In real code:
	// f, err := os.Open(filename)
	// if err != nil { ... }
	// defer f.Close() ← guaranteed to run when function exits

	fmt.Printf("Opening %s\n", filename)
	defer fmt.Printf("Closing %s\n", filename) // runs when function returns
	fmt.Printf("Reading %s\n", filename)
	// Output:
	// Opening report.txt
	// Reading report.txt
	// Closing report.txt   ← defer ran after the function body finished
}

// deferInLoop — GOTCHA: defer inside a loop — do NOT do this for resources!
// All defers stack up and run only when the FUNCTION exits, not each iteration.
func deferInLoop() {
	fmt.Println("--- Defer in loop (GOTCHA) ---")
	for i := 0; i < 3; i++ {
		// BAD for file handles / DB connections — they stay open until function returns!
		defer fmt.Println("defer in loop:", i)
	}
	// Output (after function returns):
	// defer in loop: 2
	// defer in loop: 1
	// defer in loop: 0
	//
	// FIX: use an anonymous function or a helper function instead:
	// for i := 0; i < 3; i++ {
	//     func(i int) {
	//         // open resource, defer close, use resource — all scoped here
	//     }(i)
	// }
}

func main() {
	fmt.Println("--- LIFO order ---")
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

	fmt.Println("\n--- Resource cleanup pattern ---")
	readFile("report.txt")

	fmt.Println()
	deferInLoop()
}
