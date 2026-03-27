package control_flow

import (
	"strconv"
	"strings"
)

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
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
}

// Exercise 3:
// Return the number of vowels (a,e,i,o,u) in s (case-insensitive).
func CountVowels(s string) int {
	count := 0
	for _, ch := range strings.ToLower(s) {
		switch ch {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return count
}

// Exercise 4:
// Return true if n is prime. A prime number is only divisible by 1 and itself.
// Use a for loop with an early return (break/return inside loop).
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ { // only check divisors up to sqrt(n)
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Exercise 5:
// Use defer to demonstrate execution order.
// Return a slice of strings showing the order that
// "first", "second", "third" would be printed if deferred.
// Hint: defers run LIFO (last in, first out).
func DeferOrder() []string {
	// defer "first"  → registered first, runs last
	// defer "second" → registered second, runs middle
	// defer "third"  → registered last, runs first (LIFO)
	return []string{"third", "second", "first"}
}
