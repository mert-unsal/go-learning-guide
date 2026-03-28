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
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Error Wrapping — fmt.Errorf with %%w     %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s%s▸ Building the Error Chain%s\n", cyan, bold, reset)
	fmt.Printf("  %sCall chain: main → getUserFromDB(42) → queryDB(\"\")%s\n", dim, reset)
	fmt.Printf("  %sEach layer wraps with %%w, preserving the original cause%s\n\n", dim, reset)

	err := getUserFromDB(42)
	if err != nil {
		fmt.Printf("%s%s▸ Full Error Chain%s\n", cyan, bold, reset)
		fmt.Printf("  err = %s%v%s\n", magenta, err, reset)
		fmt.Printf("  %s✔ Each layer added context: \"getUserFromDB(id=42): queryDB: database error: empty query\"%s\n\n", green, reset)

		// errors.Is: checks if error (or any wrapped error) matches target
		fmt.Printf("%s%s▸ errors.Is — Walk the Chain to Match a Sentinel%s\n", cyan, bold, reset)
		if errors.Is(err, ErrDatabase) {
			fmt.Printf("  errors.Is(err, ErrDatabase) = %s%v%s\n", magenta, true, reset)
			fmt.Printf("  %s✔ Matched! errors.Is walks through every Unwrap() in the chain%s\n", green, reset)
		}
		fmt.Printf("  %s⚠ Direct == comparison would FAIL here: the outer error is not ErrDatabase%s\n", yellow, reset)
		fmt.Printf("  (err == ErrDatabase) = %s%v%s\n\n", magenta, err == ErrDatabase, reset)

		// errors.Unwrap: get the next error in the chain
		fmt.Printf("%s%s▸ errors.Unwrap — Peel One Layer%s\n", cyan, bold, reset)
		unwrapped := errors.Unwrap(err)
		fmt.Printf("  errors.Unwrap(err) = %s%v%s\n", magenta, unwrapped, reset)
		if unwrapped != nil {
			unwrapped2 := errors.Unwrap(unwrapped)
			fmt.Printf("  errors.Unwrap(↑)   = %s%v%s\n", magenta, unwrapped2, reset)
			fmt.Printf("  %s✔ Each Unwrap peels one %%w layer, revealing the inner error%s\n\n", green, reset)
		}
	}

	fmt.Printf("%s%s▸ %%w vs %%v — A Critical Distinction%s\n", cyan, bold, reset)
	fmt.Printf("  %s✔ %%w wraps: creates an Unwrap()-able chain (errors.Is/As can traverse it)%s\n", green, reset)
	fmt.Printf("  %s⚠ %%v formats: the original error is embedded as text only — chain is broken%s\n", yellow, reset)
	fmt.Printf("  %s✔ Rule: use %%w when callers need to inspect the cause; %%v when you want to hide it%s\n", green, reset)
}
