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

func main() {
	err := findUser(0)
	if err != nil {
		// Type assert to extract details
		var nfe *NotFoundError
		if errors.As(err, &nfe) { // errors.As is the modern way
			fmt.Printf("Not found: ID=%d, Resource=%s\n", nfe.ID, nfe.Name)
		}
	}
}
