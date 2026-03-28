// Switch in Go — demonstrates all switch variants including fallthrough.
package main

import "fmt"

func main() {
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
