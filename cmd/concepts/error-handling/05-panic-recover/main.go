// 05-panic-recover demonstrates panic and recover in Go.
//
// Run:  go run .
//
// ============================================================
// PANIC AND RECOVER
// ============================================================
// panic: terminates the current goroutine, unwinds the stack
// recover: inside a deferred function, catches a panic
//
// USE panic for:
//   - Programmer errors (bugs): index out of range, nil dereference
//   - Truly unrecoverable situations
//
// DON'T use panic for:
//   - Normal error handling (use error return values instead)
//   - Expected failures (file not found, validation errors)
//
// ============================================================
// BEST PRACTICES SUMMARY
// ============================================================
// ✅ Return errors, don't ignore them
// ✅ Add context with fmt.Errorf("doing X: %w", err)
// ✅ Use errors.Is for comparison (works with wrapping)
// ✅ Use errors.As to extract error types (works with wrapping)
// ✅ Define sentinel errors for expected failure modes
// ✅ Use custom error types for structured error data
// ❌ Don't use panic for expected errors
// ❌ Don't ignore returned errors with _
// ❌ Don't use log.Fatal in library code (only in main)
package main

import "fmt"

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

// safeDivide recovers from any panic (e.g., integer division by zero)
// and converts it into an error return value.
func safeDivide(a, b int) (result int, err error) {
	// Recover from any panic in this function
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	// This will panic if b == 0 (integer division by zero)
	result = a / b
	return result, nil
}

// mustPositive panics if n is not positive — use only in init/setup code.
func mustPositive(n int) int {
	if n <= 0 {
		panic(fmt.Sprintf("expected positive number, got %d", n))
	}
	return n
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Panic & Recover — Last Resort Handling  %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s%s▸ When Does Panic Happen?%s\n", cyan, bold, reset)
	fmt.Printf("  %spanic unwinds the stack, running deferred functions, then crashes.%s\n", dim, reset)
	fmt.Printf("  %srecover() inside a defer catches the panic and returns the panic value.%s\n\n", dim, reset)

	// Recover from panic — division by zero
	fmt.Printf("%s%s▸ safeDivide(10, 0) — Panic Recovered as Error%s\n", cyan, bold, reset)
	result, err := safeDivide(10, 0)
	fmt.Printf("  result = %s%d%s\n", magenta, result, reset)
	fmt.Printf("  err    = %s%v%s\n", magenta, err, reset)
	if err != nil {
		fmt.Printf("  %s✖ Integer division by zero triggered a runtime panic%s\n", red, reset)
		fmt.Printf("  %s✔ defer+recover caught it and converted to an error return%s\n\n", green, reset)
	}

	// Successful division
	fmt.Printf("%s%s▸ safeDivide(10, 2) — No Panic%s\n", cyan, bold, reset)
	result, err = safeDivide(10, 2)
	fmt.Printf("  result = %s%d%s\n", magenta, result, reset)
	fmt.Printf("  err    = %s%v%s\n", magenta, err, reset)
	fmt.Printf("  %s✔ No panic — normal return path%s\n\n", green, reset)

	// mustPositive — the "must" pattern
	fmt.Printf("%s%s▸ mustPositive(5) — The \"must\" Convention%s\n", cyan, bold, reset)
	n := mustPositive(5)
	fmt.Printf("  n = %s%d%s\n", magenta, n, reset)
	fmt.Printf("  %s✔ Value is positive — no panic%s\n", green, reset)
	fmt.Printf("  %s⚠ mustPositive(-5) would panic — \"must\" functions crash on invalid input%s\n\n", yellow, reset)

	// Demonstrate recover from mustPositive
	fmt.Printf("%s%s▸ Recovering from mustPositive(-1)%s\n", cyan, bold, reset)
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  %s✖ Panic caught: %v%s\n", red, r, reset)
				fmt.Printf("  %s✔ recover() returned the panic value as interface{}%s\n\n", green, reset)
			}
		}()
		_ = mustPositive(-1)
	}()

	fmt.Printf("%s%s▸ Panic/Recover Rules%s\n", cyan, bold, reset)
	fmt.Printf("  %s✔ USE panic for programmer bugs: nil deref, out-of-range, impossible state%s\n", green, reset)
	fmt.Printf("  %s✔ USE \"must\" helpers in init()/main() setup — fail fast on misconfiguration%s\n", green, reset)
	fmt.Printf("  %s✖ DON'T use panic for expected errors (file not found, bad user input)%s\n", red, reset)
	fmt.Printf("  %s✖ DON'T let panics cross API boundaries — recover at goroutine/handler edges%s\n", red, reset)
	fmt.Printf("  %s⚠ In HTTP servers, add recovery middleware: defer+recover in every handler%s\n", yellow, reset)
	fmt.Printf("  %s⚠ recover() only works in the SAME goroutine — a panic in a child goroutine\n    cannot be caught by the parent%s\n", yellow, reset)
}
