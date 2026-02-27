package basics

// ============================================================
// EXERCISES — 01 Basics
// ============================================================
// Try to solve each exercise BEFORE looking at solutions.go!
//
// Instructions:
//   Implement each function body below.
//   Run: go test ./fundamentals/01_basics/... to check your work.

// Exercise 1:
// Write a function that takes a Celsius temperature and returns
// the Fahrenheit equivalent. Formula: F = C * 9/5 + 32
func CelsiusToFahrenheit(c float64) float64 {
	// TODO: implement
	return 0
}

// Exercise 2:
// Write a function that swaps two integers and returns them.
// Hint: Go supports multiple return values!
func SwapInts(a, b int) (int, int) {
	// TODO: implement
	return 0, 0
}

// Exercise 3:
// Write a function that takes a string and returns its length
// in CHARACTERS (not bytes). Strings can contain Unicode.
// Hint: convert to []rune
func CharacterCount(s string) int {
	// TODO: implement
	return 0
}

// Exercise 4:
// Create a function that returns the minimum and maximum
// values from three integers.
func MinMax(a, b, c int) (min, max int) {
	// TODO: implement
	return 0, 0
}

// Exercise 5:
// Using iota, define a type 'Direction' with constants:
// North=0, East=1, South=2, West=3
// Then write a function that returns the string name of a direction.
type Direction int

const (
	// TODO: define North, East, South, West using iota
	North Direction = iota
	East
	South
	West
)

func DirectionName(d Direction) string {
	// TODO: implement using a switch statement
	return ""
}
