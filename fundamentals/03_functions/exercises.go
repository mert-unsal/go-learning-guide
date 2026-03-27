package functions

import "slices"

// ============================================================
// EXERCISES — 03 Functions
// ============================================================

// Exercise 1:
// Write a function that returns both the min and max of a slice.
// Return (0, 0) for an empty slice.
func MinMax(nums []int) (min, max int) {
	// TODO: implement with multiple return values
	if len(nums) == 0 {
		return 0, 0
	}
	return slices.Min(nums), slices.Max(nums)
}

// Exercise 2:
// Write a variadic function that sums any number of integers.
// Example: Sum(1, 2, 3) → 6
func Sum(nums ...int) int {
	// TODO: implement
	sum := 0
	for _, v := range nums {
		sum += v
	}
	return sum
}

// Exercise 3:
// Write a function Apply that takes a slice and a function,
// and returns a new slice with the function applied to each element.
// Example: Apply([]int{1,2,3}, func(x int) int { return x*2 }) → [2,4,6]
func Apply(nums []int, fn func(int) int) []int {
	// TODO: implement (higher-order function)
	result := make([]int, len(nums))
	for i, v := range nums {
		result[i] = fn(v)
	}
	return result
}

// Exercise 4:
// Write a function MakeAdder that returns a closure.
// The closure adds n to whatever value is passed.
// Example: add5 := MakeAdder(5); add5(3) → 8
func MakeAdder(n int) func(int) int {
	// TODO: implement (closure)
	sum := n
	return func(x int) int { return sum + x }
}

// Exercise 5:
// Write a recursive function that computes the nth Fibonacci number.
// fib(0)=0, fib(1)=1, fib(n)=fib(n-1)+fib(n-2)
// Then write a memoized version using a map.
func Fibonacci(n int) int {
	// TODO: implement recursive version
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return Fibonacci(n-1) + Fibonacci(n-2)
	}
}

func FibonacciMemo(n int) int {
	// TODO: implement with memoization (use a map inside or as parameter)
	memo := make(map[int]int)
	return fibMemo(n, memo)
}

func fibMemo(n int, memo map[int]int) int {
	// TODO: implement
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	}
	if val, ok := memo[n]; ok {
		return val
	}
	memo[n] = fibMemo(n-1, memo) + fibMemo(n-2, memo)
	return memo[n]
}
