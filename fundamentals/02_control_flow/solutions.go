package control_flow

import "strconv"

// ============================================================
// SOLUTIONS — 02 Control Flow
// ============================================================
// Only look here AFTER trying the exercises yourself!

func FizzBuzzSwitchSolution(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return strconv.Itoa(n)
	}
}

func SumToSolution(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
	// Math shortcut: return n * (n + 1) / 2
}

func CountVowelsSolution(s string) int {
	count := 0
	for _, ch := range s {
		switch ch {
		case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
			count++
		}
	}
	return count
}

func IsPrimeSolution(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ { // only check up to sqrt(n)
		if n%i == 0 {
			return false
		}
	}
	return true
}

func DeferOrderSolution() []string {
	// Defers run LIFO (stack order)
	// defer "first"  → runs 3rd
	// defer "second" → runs 2nd
	// defer "third"  → runs 1st
	return []string{"third", "second", "first"}
}
