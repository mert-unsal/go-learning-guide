// Defer Argument Evaluation — a very common interview question!
//
// RULE: Arguments are evaluated at DEFER time, not when the deferred call executes.
// RULE: Defer CAN modify named return values.
// This is useful for wrapping errors or ensuring consistent cleanup.
package main

import "fmt"

// deferWithNamedReturn demonstrates that defer can modify named return values.
// The deferred closure overwrites whatever was returned.
func deferWithNamedReturn() (result string) {
	defer func() {
		result = "modified by defer" // overwrites whatever was returned
	}()
	return "original" // this value gets overwritten by defer above
}

func main() {
	fmt.Println("--- Argument evaluation at defer time ---")
	x := 10
	defer fmt.Println("deferred x =", x) // captures x=10 RIGHT NOW
	x = 99
	fmt.Println("current x =", x) // prints 99
	// Output:
	// current x = 99
	// deferred x = 10  ← x was captured as 10, not 99

	fmt.Println("\n--- Named return modification ---")
	result := deferWithNamedReturn()
	fmt.Println("deferWithNamedReturn() =", result) // "modified by defer"
}
