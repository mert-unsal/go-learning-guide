// 02-custom-errors demonstrates custom error types in Go.
//
// Run:  go run .
//
// ============================================================
// CUSTOM ERROR TYPES
// ============================================================
// Implement the error interface with a struct for rich error info.
// Use errors.As (not type assertion) to extract details —
// errors.As works correctly through wrapped error chains.
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

// NotFoundError carries structured info about a missing resource.
type NotFoundError struct {
	ID   int
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Name, e.ID)
}

// PermissionError carries structured info about a denied action.
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

func checkPermission(user, action string) error {
	if user == "guest" {
		return &PermissionError{User: user, Action: action}
	}
	return nil
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Custom Errors — Structured Error Types  %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s%s▸ Why Custom Error Types?%s\n", cyan, bold, reset)
	fmt.Printf("  %sA plain string error loses structured data. Custom types carry fields%s\n", dim, reset)
	fmt.Printf("  %sthat callers can inspect programmatically with errors.As().%s\n\n", dim, reset)

	// NotFoundError demo
	fmt.Printf("%s%s▸ NotFoundError: findUser(0)%s\n", cyan, bold, reset)
	err := findUser(0)
	if err != nil {
		fmt.Printf("  err.Error() = %s%q%s\n", magenta, err.Error(), reset)

		// Type assert to extract details
		var nfe *NotFoundError
		if errors.As(err, &nfe) {
			fmt.Printf("  %s✔ errors.As matched *NotFoundError%s\n", green, reset)
			fmt.Printf("  nfe.ID   = %s%d%s\n", magenta, nfe.ID, reset)
			fmt.Printf("  nfe.Name = %s%q%s\n\n", magenta, nfe.Name, reset)
		}
	}

	// PermissionError demo
	fmt.Printf("%s%s▸ PermissionError: checkPermission(\"guest\", \"delete\")%s\n", cyan, bold, reset)
	err = checkPermission("guest", "delete")
	if err != nil {
		fmt.Printf("  err.Error() = %s%q%s\n", magenta, err.Error(), reset)

		var pe *PermissionError
		if errors.As(err, &pe) {
			fmt.Printf("  %s✔ errors.As matched *PermissionError%s\n", green, reset)
			fmt.Printf("  pe.User   = %s%q%s\n", magenta, pe.User, reset)
			fmt.Printf("  pe.Action = %s%q%s\n\n", magenta, pe.Action, reset)
		}
	}

	// Success path
	fmt.Printf("%s%s▸ Success Case: findUser(1)%s\n", cyan, bold, reset)
	err = findUser(1)
	if err == nil {
		fmt.Printf("  %s✔ No error — user found successfully%s\n\n", green, reset)
	}

	fmt.Printf("%s%s▸ errors.As vs Type Assertion%s\n", cyan, bold, reset)
	fmt.Printf("  %s⚠ Direct type assertion (err.(*NotFoundError)) breaks on wrapped errors%s\n", yellow, reset)
	fmt.Printf("  %s✔ errors.As walks the entire wrapped chain — always prefer it%s\n", green, reset)
	fmt.Printf("  %s✔ Use pointer receiver on Error() so the interface is satisfied by *T, not T%s\n", green, reset)
	fmt.Printf("  %s✔ Custom types let you carry domain context (IDs, field names, codes)%s\n", green, reset)
}
