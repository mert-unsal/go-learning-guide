// 03-error-wrapping demonstrates error wrapping with fmt.Errorf and %w.
//
// Run:  go run .
//
// ============================================================
// ERROR WRAPPING WITH fmt.Errorf AND %w
// ============================================================
// Wrap errors to add context while preserving the original cause.
// Unwrap with errors.Is (value comparison) and errors.As (type extraction).
// errors.Unwrap returns the next error in the chain.
//
// Key insight: %w creates a chain — errors.Is walks the entire chain,
// so you can match a sentinel even through multiple layers of wrapping.
package main

import (
	"errors"
	"fmt"
)

// ErrDatabase is a sentinel error representing a database failure.
var ErrDatabase = errors.New("database error")

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

func main() {
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
