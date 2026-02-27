package functions

// ============================================================
// SOLUTIONS — 03 Functions
// ============================================================

func MinMaxSolution(nums []int) (min, max int) {
	if len(nums) == 0 {
		return 0, 0
	}
	min, max = nums[0], nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return // naked return — returns named values min and max
}

func SumSolution(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func ApplySolution(nums []int, fn func(int) int) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = fn(n)
	}
	return result
}

func MakeAdderSolution(n int) func(int) int {
	// n is captured in the closure — it persists between calls
	return func(x int) int {
		return x + n
	}
}

func FibonacciSolution(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciSolution(n-1) + FibonacciSolution(n-2)
}

func FibonacciMemoSolution(n int) int {
	memo := make(map[int]int)
	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		if v, ok := memo[n]; ok {
			return v // cache hit
		}
		memo[n] = fib(n-1) + fib(n-2) // compute and cache
		return memo[n]
	}
	return fib(n)
}
