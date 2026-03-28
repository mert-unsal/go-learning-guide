// 01-basic-errors demonstrates Go's fundamental error handling pattern.
//
// Run:  go run .
//
// ============================================================
// THE ERROR INTERFACE
// ============================================================
// error is a built-in interface: type error interface { Error() string }
// Convention: return error as the LAST return value.
// Convention: name the error return value 'err'.
// Convention: check errors IMMEDIATELY after calling a function.
package main

import (
	"errors"
	"fmt"
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

// divide returns an error when b is zero, following Go's convention
// of (result, error) return pairs.
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Basic Errors — The Error Interface      %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s%s▸ The (result, error) Return Pattern%s\n", cyan, bold, reset)
	fmt.Printf("  %sGo convention: always return error as the last value%s\n", dim, reset)
	fmt.Printf("  %sCallers MUST check err != nil immediately after the call%s\n\n", dim, reset)

	// Always handle errors!
	fmt.Printf("%s%s▸ Successful Division: divide(10, 2)%s\n", cyan, bold, reset)
	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("  %s✖ Error: %v%s\n", red, err, reset)
		return
	}
	fmt.Printf("  result = %s%v%s\n", magenta, result, reset)
	fmt.Printf("  err    = %s%v%s\n", magenta, err, reset)
	fmt.Printf("  %s✔ err is nil — no error occurred%s\n\n", green, reset)

	// Error case
	fmt.Printf("%s%s▸ Error Case: divide(5, 0)%s\n", cyan, bold, reset)
	_, err = divide(5, 0)
	if err != nil {
		fmt.Printf("  err = %s%v%s\n", magenta, err, reset)
		fmt.Printf("  %s✖ errors.New(\"division by zero\") returned%s\n", red, reset)
		fmt.Printf("  %s✔ We checked err != nil and handled it gracefully%s\n\n", green, reset)
	}

	fmt.Printf("%s%s▸ Under the Hood%s\n", cyan, bold, reset)
	fmt.Printf("  %s✔ error is a built-in interface: type error interface { Error() string }%s\n", green, reset)
	fmt.Printf("  %s✔ errors.New() returns a *errors.errorString (unexported struct with pointer receiver)%s\n", green, reset)
	fmt.Printf("  %s✔ Pointer receiver means each errors.New() call creates a distinct error value%s\n", green, reset)
	fmt.Printf("  %s⚠ Never compare errors with == on wrapped errors — use errors.Is() instead%s\n", yellow, reset)
}
