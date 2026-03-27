package arrays

import "strconv"

// ============================================================
// FizzBuzz — [E]
// ============================================================
// For each number from 1 to n:
// - "Fizz" if divisible by 3
// - "Buzz" if divisible by 5
// - "FizzBuzz" if divisible by both
// - the number itself otherwise

// FizzBuzz returns the FizzBuzz sequence up to n.
// Time: O(n)  Space: O(n)
func FizzBuzz(n int) []string {
	result := make([]string, n)
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			result[i-1] = "FizzBuzz"
		case i%3 == 0:
			result[i-1] = "Fizz"
		case i%5 == 0:
			result[i-1] = "Buzz"
		default:
			result[i-1] = strconv.Itoa(i)
		}
	}
	return result
}
