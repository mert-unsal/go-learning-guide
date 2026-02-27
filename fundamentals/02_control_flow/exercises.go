package control_flow

import "strconv"

// ============================================================
// EXERCISES — 02 Control Flow
// ============================================================
// Implement each function. Run: go test ./fundamentals/02_control_flow/...

// Exercise 1:
// Return "Fizz" if n divisible by 3, "Buzz" if by 5,
// "FizzBuzz" if both, else the number as a string.
// Use switch, NOT if/else.
func FizzBuzzSwitch(n int) string {
	// TODO: implement using switch
	switch {
	case n%3 == 0 && n%5 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return strconv.Itoa(n)
	}

}

// Exercise 2:
// Sum all integers from 1 to n (inclusive) using a for loop.
func SumTo(n int) int {
	// TODO: implement
	return 0
}

// Exercise 3:
// Return the number of vowels (a,e,i,o,u) in s (case-insensitive).
func CountVowels(s string) int {
	// TODO: implement using for range
	return 0
}

// Exercise 4:
// Return true if n is prime. A prime number is only divisible by 1 and itself.
// Use a for loop with an early return (break/return inside loop).
func IsPrime(n int) bool {
	// TODO: implement
	return false
}

// Exercise 5:
// Use defer to demonstrate execution order.
// Return a slice of strings showing the order that
// "first", "second", "third" would be printed if deferred.
// Hint: defers run LIFO (last in, first out).
func DeferOrder() []string {
	// TODO: return []string{"third", "second", "first"}
	// (the order defers execute: last-deferred runs first)
	return nil
}
