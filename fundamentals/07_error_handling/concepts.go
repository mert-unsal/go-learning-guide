// Package error_handling covers Go's error handling philosophy:
// the error interface, errors.New, fmt.Errorf with %w (wrapping),
// custom error types, panic, and recover.
package error_handling

import (
	"errors"
	"fmt"
)

// ============================================================
// 1. THE ERROR INTERFACE
// ============================================================
// error is a built-in interface: type error interface { Error() string }
// Convention: return error as the LAST return value.
// Convention: name the error return value 'err'.
// Convention: check errors IMMEDIATELY after calling a function.

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func DemonstrateBasicErrors() {
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

// ============================================================
// 2. CUSTOM ERROR TYPES
// ============================================================
// Implement the error interface with a struct for rich error info.

type NotFoundError struct {
	ID   int
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Name, e.ID)
}

type PermissionError struct {
	User   string
	Action string
}

func (e *PermissionError) Error() string {
	return fmt.Sprintf("user %q is not allowed to %s", e.User, e.Action)
}

func findUser(id int) error {
	if id == 0 {
		return &NotFoundError{ID: id, Name: "User"}
	}
	return nil
}

func DemonstrateCustomErrors() {
	err := findUser(0)
	if err != nil {
		// Type assert to extract details
		var nfe *NotFoundError
		if errors.As(err, &nfe) { // errors.As is the modern way
			fmt.Printf("Not found: ID=%d, Resource=%s\n", nfe.ID, nfe.Name)
		}
	}
}

// ============================================================
// 3. ERROR WRAPPING WITH fmt.Errorf and %w
// ============================================================
// Wrap errors to add context. Unwrap with errors.Is / errors.As.

var ErrDatabase = errors.New("database error") // sentinel error

func queryDB(query string) error {
	if query == "" {
		return fmt.Errorf("queryDB: %w: empty query", ErrDatabase)
	}
	return nil
}

func getUserFromDB(id int) error {
	err := queryDB("") // simulate error
	if err != nil {
		return fmt.Errorf("getUserFromDB(id=%d): %w", id, err)
	}
	return nil
}

func DemonstrateErrorWrapping() {
	err := getUserFromDB(42)
	if err != nil {
		fmt.Println("Error:", err) // full chain

		// errors.Is: checks if error (or any wrapped error) matches target
		if errors.Is(err, ErrDatabase) {
			fmt.Println("This is a database error")
		}

		// errors.Unwrap: get the next error in the chain
		fmt.Println("Unwrapped:", errors.Unwrap(err))
	}
}

// ============================================================
// 4. SENTINEL ERRORS (common pattern)
// ============================================================
// Predefined errors that callers can compare against.
// Convention: name them Err* or err* (if unexported).

var (
	ErrNotFound   = errors.New("not found")
	ErrPermission = errors.New("permission denied")
	ErrTimeout    = errors.New("timeout")
)

func fetchResource(id int, user string) error {
	if id > 1000 {
		return fmt.Errorf("fetchResource: %w", ErrNotFound)
	}
	if user == "guest" {
		return fmt.Errorf("fetchResource: %w", ErrPermission)
	}
	return nil
}

func DemonstrateSentinelErrors() {
	err := fetchResource(9999, "alice")
	switch {
	case errors.Is(err, ErrNotFound):
		fmt.Println("Resource not found")
	case errors.Is(err, ErrPermission):
		fmt.Println("Permission denied")
	case err != nil:
		fmt.Println("Other error:", err)
	default:
		fmt.Println("Success")
	}
}

// ============================================================
// 5. PANIC AND RECOVER
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

func mustPositive(n int) int {
	if n <= 0 {
		panic(fmt.Sprintf("expected positive number, got %d", n))
	}
	return n
}

func DemonstratePanicRecover() {
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

// ============================================================
// 6. BEST PRACTICES SUMMARY
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

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Basic Errors ===")
	DemonstrateBasicErrors()
	fmt.Println("\n=== Custom Errors ===")
	DemonstrateCustomErrors()
	fmt.Println("\n=== Error Wrapping ===")
	DemonstrateErrorWrapping()
	fmt.Println("\n=== Sentinel Errors ===")
	DemonstrateSentinelErrors()
	fmt.Println("\n=== Panic/Recover ===")
	DemonstratePanicRecover()
}
