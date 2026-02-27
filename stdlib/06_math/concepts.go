// Package math_pkg demonstrates the math, math/rand, and math/big standard libraries.
// Run: go run stdlib/03_math/concepts.go
package math_pkg

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"
)

// ============================================================
// 1. THE math PACKAGE
// ============================================================
// The math package provides constants and functions for floating-point math.
// All functions operate on float64.

func DemonstrateMath() {
	// --- Constants ---
	fmt.Println("math.Pi     =", math.Pi)     // 3.141592653589793
	fmt.Println("math.E      =", math.E)      // 2.718281828459045
	fmt.Println("math.Sqrt2  =", math.Sqrt2)  // 1.4142135623730951
	fmt.Println("math.MaxInt =", math.MaxInt) // largest int value
	fmt.Println("math.MinInt =", math.MinInt) // smallest int value

	// Infinity and NaN
	inf := math.Inf(1) // positive infinity
	nan := math.NaN()
	fmt.Println("Inf(1):", inf, "IsInf:", math.IsInf(inf, 1))
	fmt.Println("NaN:   ", nan, "IsNaN:", math.IsNaN(nan))

	// --- Basic functions ---
	fmt.Println("Abs(-3.5)       =", math.Abs(-3.5))  // 3.5
	fmt.Println("Ceil(2.3)       =", math.Ceil(2.3))  // 3
	fmt.Println("Floor(2.9)      =", math.Floor(2.9)) // 2
	fmt.Println("Round(2.5)      =", math.Round(2.5)) // 3
	fmt.Println("Trunc(3.9)      =", math.Trunc(3.9)) // 3

	// --- Power & roots ---
	fmt.Println("Sqrt(16)        =", math.Sqrt(16))   // 4
	fmt.Println("Pow(2, 10)      =", math.Pow(2, 10)) // 1024
	fmt.Println("Cbrt(27)        =", math.Cbrt(27))   // 3

	// --- Logarithms ---
	fmt.Println("Log(math.E)     =", math.Log(math.E)) // 1 (natural log)
	fmt.Println("Log2(1024)      =", math.Log2(1024))  // 10
	fmt.Println("Log10(1000)     =", math.Log10(1000)) // 3

	// --- Min / Max ---
	fmt.Println("Max(3.0, 7.0)   =", math.Max(3.0, 7.0)) // 7
	fmt.Println("Min(3.0, 7.0)   =", math.Min(3.0, 7.0)) // 3

	// --- Trigonometry (radians) ---
	fmt.Println("Sin(Pi/2)       =", math.Sin(math.Pi/2)) // 1
	fmt.Println("Cos(0)          =", math.Cos(0))         // 1

	// --- Integer max/min helpers (common in competitive programming) ---
	const MaxInt = int(^uint(0) >> 1)
	const MinInt = -MaxInt - 1
	fmt.Printf("MaxInt: %d  MinInt: %d\n", MaxInt, MinInt)
}

// ============================================================
// 2. math/rand — PSEUDO-RANDOM NUMBERS
// ============================================================
// Go 1.20+ rand functions are automatically seeded.
// For older versions, you must call rand.Seed(time.Now().UnixNano()).
//
// Note: math/rand produces PSEUDO-random numbers (deterministic given a seed).
// For cryptographic randomness, use crypto/rand.

func DemonstrateRand() {
	// Seed for reproducibility in examples (in production, skip this)
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// Random int in [0, n)
	fmt.Println("Intn(100):", rng.Intn(100))

	// Random float in [0.0, 1.0)
	fmt.Printf("Float64: %.4f\n", rng.Float64())

	// Random int in range [min, max)
	min, max := 10, 50
	fmt.Println("Range [10,50):", min+rng.Intn(max-min))

	// Shuffle a slice
	s := []int{1, 2, 3, 4, 5}
	rng.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	fmt.Println("Shuffled:", s)

	// Random permutation
	perm := rng.Perm(5) // [0..4] in random order
	fmt.Println("Perm(5):", perm)
}

// ============================================================
// 3. math/big — ARBITRARY PRECISION INTEGERS
// ============================================================
// Use when values exceed int64 (~9.2 × 10^18).
// Common examples: large factorials, Fibonacci, cryptographic computations.
//
// big.Int methods mutate the receiver and return it for chaining.

func DemonstrateBigInt() {
	// --- Basic big.Int operations ---
	a := big.NewInt(1000000000) // 10^9
	b := big.NewInt(1000000000)
	product := new(big.Int).Mul(a, b)     // 10^18
	fmt.Println("10^9 × 10^9 =", product) // 1000000000000000000

	// Compare
	c := big.NewInt(42)
	d := big.NewInt(100)
	fmt.Println("42.Cmp(100):", c.Cmp(d)) // -1 (less than)

	// --- Factorial using big.Int ---
	fmt.Println("50! =", bigFactorial(50))

	// --- Fibonacci using big.Int ---
	fmt.Println("Fib(100) =", bigFibonacci(100))
}

// bigFactorial computes n! using arbitrary-precision arithmetic.
func bigFactorial(n int64) *big.Int {
	result := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}
	return result
}

// bigFibonacci computes the nth Fibonacci number using big.Int.
func bigFibonacci(n int) *big.Int {
	if n <= 0 {
		return big.NewInt(0)
	}
	if n == 1 {
		return big.NewInt(1)
	}
	prev2 := big.NewInt(0)
	prev1 := big.NewInt(1)
	cur := new(big.Int)
	for i := 2; i <= n; i++ {
		cur.Add(prev1, prev2)
		prev2.Set(prev1)
		prev1.Set(cur)
	}
	return prev1
}

// ============================================================
// 4. COMMON PATTERNS IN LEETCODE
// ============================================================

// IntMax returns the larger of two ints (math.Max only works on float64).
func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// IntMin returns the smaller of two ints.
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// IntAbs returns the absolute value of an int.
func IntAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Clamp returns x clamped to [lo, hi].
func Clamp(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
