package control_flow

// ============================================================
// EXERCISES — 02 Control Flow
// ============================================================
// Implement each function. Run: go test . -v  (from this folder)
//                          or: go test ./fundamentals/02_control_flow/... -v  (from project root)

// Exercise 1:
// Return "Fizz" if n divisible by 3, "Buzz" if by 5,
// "FizzBuzz" if both, else the number as a string.
// Use switch, NOT if/else.
func FizzBuzzSwitch(n int) string {
	return ""
}

// Exercise 2:
// Sum all integers from 1 to n (inclusive) using a for loop.
func SumTo(n int) int {
	return 0
}

// Exercise 3:
// Return the number of vowels (a,e,i,o,u) in s (case-insensitive).
func CountVowels(s string) int {
	return 0
}

// Exercise 4:
// Return true if n is prime. A prime number is only divisible by 1 and itself.
// Use a for loop with an early return (break/return inside loop).
func IsPrime(n int) bool {
	return false
}

// Exercise 5:
// Use defer to demonstrate execution order.
// Return a slice of strings showing the order that
// "first", "second", "third" would be printed if deferred.
// Hint: defers run LIFO (last in, first out).
func DeferOrder() []string {
	return nil
}
