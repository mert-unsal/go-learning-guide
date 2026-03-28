// If/Else in Go — demonstrates basic conditionals and init statements.
package main

import "fmt"

func main() {
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
