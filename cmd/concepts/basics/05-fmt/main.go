// fmt Package — formatting verbs and output functions.
//
// Key formatting verbs:
//
//	%s  — string
//	%d  — integer (decimal)
//	%f  — float (default width)
//	%.2f — float with 2 decimal places
//	%e  — scientific notation
//	%t  — boolean
//	%T  — prints the TYPE of a variable
//	%v  — default format (works for anything)
//	%#v — Go-syntax representation
//	%q  — quoted string
//	%p  — pointer address
//	%b  — binary
//	%o  — octal
//	%x  — hexadecimal
//
// Stringer interface — any type with String() method works with %s and %v
package main

import "fmt"

func main() {
	name := "Gopher"
	age := 5
	pi := 3.14159

	// Printf: formatted printing
	fmt.Printf("Name: %s\n", name)       // string
	fmt.Printf("Age: %d\n", age)         // integer
	fmt.Printf("Pi: %.2f\n", pi)         // float with 2 decimal places
	fmt.Printf("Pi: %e\n", pi)           // scientific notation
	fmt.Printf("Bool: %t\n", true)       // boolean
	fmt.Printf("Type: %T\n", name)       // prints the TYPE of variable
	fmt.Printf("Value: %v\n", age)       // default format (works for anything)
	fmt.Printf("Go syntax: %#v\n", name) // Go-syntax representation

	// Sprintf: returns formatted string (doesn't print)
	msg := fmt.Sprintf("Hello, %s! You are %d years old.", name, age)
	fmt.Println(msg)

	// Stringer interface — any type with String() method works with %s and %v
}
