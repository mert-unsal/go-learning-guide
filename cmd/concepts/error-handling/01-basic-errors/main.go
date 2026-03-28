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

// divide returns an error when b is zero, following Go's convention
// of (result, error) return pairs.
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	// Always handle errors!
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", result)

	// Error case
	_, err = divide(5, 0)
	if err != nil {
		fmt.Println("Got error:", err)
	}
}
