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
	// Recover from panic
	result, err := safeDivide(10, 0)
	fmt.Println("safeDivide(10, 0):", result, err)

	result, err = safeDivide(10, 2)
	fmt.Println("safeDivide(10, 2):", result, err)

	// mustPositive — use in initialization code only
	// n := mustPositive(-5) // would panic
	n := mustPositive(5)
	fmt.Println("mustPositive(5):", n)
}
