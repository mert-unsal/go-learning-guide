// Math in Go — demonstrates the math, math/rand, and math/big standard libraries.
//
// Topics:
//   - math constants and functions (float64 operations)
//   - math/rand pseudo-random numbers
//   - math/big arbitrary precision integers
//   - Common integer helpers for competitive programming
//
// Run: go run cmd/concepts/stdlib/06-math/main.go
package main

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Math, Rand & Big                        %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	demonstrateMath()
	demonstrateRand()
	demonstrateBigInt()
	demonstrateIntHelpers()
}

// ============================================================
// 1. THE math PACKAGE
// ============================================================
// The math package provides constants and functions for floating-point math.
// All functions operate on float64.

func demonstrateMath() {
	fmt.Printf("%s▸ 1. math Package — Constants & Functions%s\n", cyan+bold, reset)

	// --- Constants ---
	fmt.Printf("\n  %s✔ Constants%s\n", green, reset)
	fmt.Println("  math.Pi     =", math.Pi)     // 3.141592653589793
	fmt.Println("  math.E      =", math.E)      // 2.718281828459045
	fmt.Println("  math.Sqrt2  =", math.Sqrt2)  // 1.4142135623730951
	fmt.Println("  math.MaxInt =", math.MaxInt) // largest int value
	fmt.Println("  math.MinInt =", math.MinInt) // smallest int value

	// Infinity and NaN
	fmt.Printf("\n  %s✔ Infinity & NaN%s\n", green, reset)
	inf := math.Inf(1) // positive infinity
	nan := math.NaN()
	fmt.Println("  Inf(1):", inf, "IsInf:", math.IsInf(inf, 1))
	fmt.Println("  NaN:   ", nan, "IsNaN:", math.IsNaN(nan))

	// --- Basic functions ---
	fmt.Printf("\n  %s✔ Rounding & Absolute Value%s\n", green, reset)
	fmt.Println("  Abs(-3.5)       =", math.Abs(-3.5))  // 3.5
	fmt.Println("  Ceil(2.3)       =", math.Ceil(2.3))  // 3
	fmt.Println("  Floor(2.9)      =", math.Floor(2.9)) // 2
	fmt.Println("  Round(2.5)      =", math.Round(2.5)) // 3
	fmt.Println("  Trunc(3.9)      =", math.Trunc(3.9)) // 3

	// --- Power & roots ---
	fmt.Printf("\n  %s✔ Powers & Roots%s\n", green, reset)
	fmt.Println("  Sqrt(16)        =", math.Sqrt(16))   // 4
	fmt.Println("  Pow(2, 10)      =", math.Pow(2, 10)) // 1024
	fmt.Println("  Cbrt(27)        =", math.Cbrt(27))   // 3

	// --- Logarithms ---
	fmt.Printf("\n  %s✔ Logarithms%s\n", green, reset)
	fmt.Println("  Log(math.E)     =", math.Log(math.E)) // 1 (natural log)
	fmt.Println("  Log2(1024)      =", math.Log2(1024))  // 10
	fmt.Println("  Log10(1000)     =", math.Log10(1000)) // 3

	// --- Min / Max ---
	fmt.Printf("\n  %s✔ Min / Max%s\n", green, reset)
	fmt.Println("  Max(3.0, 7.0)   =", math.Max(3.0, 7.0)) // 7
	fmt.Println("  Min(3.0, 7.0)   =", math.Min(3.0, 7.0)) // 3

	// --- Trigonometry (radians) ---
	fmt.Printf("\n  %s✔ Trigonometry (radians)%s\n", green, reset)
	fmt.Println("  Sin(Pi/2)       =", math.Sin(math.Pi/2)) // 1
	fmt.Println("  Cos(0)          =", math.Cos(0))         // 1

	// --- Integer max/min helpers (common in competitive programming) ---
	fmt.Printf("\n  %s✔ Integer Max/Min (manual — math.Max is float64 only)%s\n", green, reset)
	const MaxInt = int(^uint(0) >> 1)
	const MinInt = -MaxInt - 1
	fmt.Printf("  MaxInt: %d  MinInt: %d\n", MaxInt, MinInt)
	fmt.Println()
}

// ============================================================
// 2. math/rand — PSEUDO-RANDOM NUMBERS
// ============================================================
// Go 1.20+ rand functions are automatically seeded.
// For older versions, you must call rand.Seed(time.Now().UnixNano()).
//
// Note: math/rand produces PSEUDO-random numbers (deterministic given a seed).
// For cryptographic randomness, use crypto/rand.

func demonstrateRand() {
	fmt.Printf("%s▸ 2. math/rand — Pseudo-Random Numbers%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ math/rand is PSEUDO-random — use crypto/rand for security%s\n", yellow, reset)

	// Seed for reproducibility in examples (in production, skip this)
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// Random int in [0, n)
	fmt.Println("  Intn(100):", rng.Intn(100))

	// Random float in [0.0, 1.0)
	fmt.Printf("  Float64: %.4f\n", rng.Float64())

	// Random int in range [min, max)
	min, max := 10, 50
	fmt.Println("  Range [10,50):", min+rng.Intn(max-min))

	// Shuffle a slice
	s := []int{1, 2, 3, 4, 5}
	rng.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	fmt.Println("  Shuffled:", s)

	// Random permutation
	perm := rng.Perm(5) // [0..4] in random order
	fmt.Println("  Perm(5):", perm)
	fmt.Println()
}

// ============================================================
// 3. math/big — ARBITRARY PRECISION INTEGERS
// ============================================================
// Use when values exceed int64 (~9.2 × 10^18).
// Common examples: large factorials, Fibonacci, cryptographic computations.
//
// big.Int methods mutate the receiver and return it for chaining.

func demonstrateBigInt() {
	fmt.Printf("%s▸ 3. math/big — Arbitrary Precision%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Use when values exceed int64 (~9.2 × 10^18)%s\n", green, reset)

	// --- Basic big.Int operations ---
	a := big.NewInt(1000000000) // 10^9
	b := big.NewInt(1000000000)
	product := new(big.Int).Mul(a, b)     // 10^18
	fmt.Println("  10^9 × 10^9 =", product) // 1000000000000000000

	// Compare
	c := big.NewInt(42)
	d := big.NewInt(100)
	fmt.Println("  42.Cmp(100):", c.Cmp(d)) // -1 (less than)

	// --- Factorial using big.Int ---
	fmt.Println("  50! =", bigFactorial(50))

	// --- Fibonacci using big.Int ---
	fmt.Println("  Fib(100) =", bigFibonacci(100))
	fmt.Println()
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

// intMax returns the larger of two ints (math.Max only works on float64).
func intMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// intMin returns the smaller of two ints.
func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// intAbs returns the absolute value of an int.
func intAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// clamp returns x clamped to [lo, hi].
func clamp(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func demonstrateIntHelpers() {
	fmt.Printf("%s▸ 4. Integer Helpers (LeetCode Patterns)%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ math.Max/Min only work on float64 — use helpers for int%s\n", green, reset)

	fmt.Println("  intMax(3, 7)  =", intMax(3, 7))
	fmt.Println("  intMin(3, 7)  =", intMin(3, 7))
	fmt.Println("  intAbs(-5)    =", intAbs(-5))
	fmt.Println("  clamp(15, 0, 10) =", clamp(15, 0, 10))
	fmt.Println("  clamp(-3, 0, 10) =", clamp(-3, 0, 10))
	fmt.Println("  clamp(5, 0, 10)  =", clamp(5, 0, 10))
}
