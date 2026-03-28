// 04-sentinel-errors demonstrates predefined sentinel errors in Go.
//
// Run:  go run .
//
// ============================================================
// SENTINEL ERRORS (common pattern)
// ============================================================
// Predefined errors that callers can compare against.
// Convention: name them Err* (exported) or err* (unexported).
// Always compare with errors.Is — it walks the wrapped chain,
// unlike direct == which would fail on wrapped errors.
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

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Sentinel Errors — Predefined Error Values%s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s%s▸ Defined Sentinels%s\n", cyan, bold, reset)
	fmt.Printf("  ErrNotFound   = %s%v%s\n", magenta, ErrNotFound, reset)
	fmt.Printf("  ErrPermission = %s%v%s\n", magenta, ErrPermission, reset)
	fmt.Printf("  ErrTimeout    = %s%v%s\n", magenta, ErrTimeout, reset)
	fmt.Printf("  %s✔ Convention: package-level var named Err* (exported) or err* (unexported)%s\n\n", green, reset)

	// Case 1: Not found
	fmt.Printf("%s%s▸ Case 1: fetchResource(9999, \"alice\") — Not Found%s\n", cyan, bold, reset)
	err := fetchResource(9999, "alice")
	fmt.Printf("  err = %s%v%s\n", magenta, err, reset)
	switch {
	case errors.Is(err, ErrNotFound):
		fmt.Printf("  errors.Is(err, ErrNotFound) = %s%v%s\n", magenta, true, reset)
		fmt.Printf("  %s✖ Resource not found%s\n\n", red, reset)
	case errors.Is(err, ErrPermission):
		fmt.Printf("  %s✖ Permission denied%s\n\n", red, reset)
	case err != nil:
		fmt.Printf("  %s✖ Other error: %v%s\n\n", red, err, reset)
	default:
		fmt.Printf("  %s✔ Success%s\n\n", green, reset)
	}

	// Case 2: Permission denied
	fmt.Printf("%s%s▸ Case 2: fetchResource(1, \"guest\") — Permission Denied%s\n", cyan, bold, reset)
	err = fetchResource(1, "guest")
	fmt.Printf("  err = %s%v%s\n", magenta, err, reset)
	switch {
	case errors.Is(err, ErrNotFound):
		fmt.Printf("  %s✖ Resource not found%s\n\n", red, reset)
	case errors.Is(err, ErrPermission):
		fmt.Printf("  errors.Is(err, ErrPermission) = %s%v%s\n", magenta, true, reset)
		fmt.Printf("  %s✖ Permission denied%s\n\n", red, reset)
	case err != nil:
		fmt.Printf("  %s✖ Other error: %v%s\n\n", red, err, reset)
	default:
		fmt.Printf("  %s✔ Success%s\n\n", green, reset)
	}

	// Case 3: Success
	fmt.Printf("%s%s▸ Case 3: fetchResource(1, \"alice\") — Success%s\n", cyan, bold, reset)
	err = fetchResource(1, "alice")
	fmt.Printf("  err = %s%v%s\n", magenta, err, reset)
	switch {
	case errors.Is(err, ErrNotFound):
		fmt.Printf("  %s✖ Resource not found%s\n\n", red, reset)
	case errors.Is(err, ErrPermission):
		fmt.Printf("  %s✖ Permission denied%s\n\n", red, reset)
	case err != nil:
		fmt.Printf("  %s✖ Other error: %v%s\n\n", red, err, reset)
	default:
		fmt.Printf("  %s✔ No error — resource fetched successfully%s\n\n", green, reset)
	}

	fmt.Printf("%s%s▸ Sentinel Errors Best Practices%s\n", cyan, bold, reset)
	fmt.Printf("  %s✔ Always compare sentinels with errors.Is() — it walks wrapped chains%s\n", green, reset)
	fmt.Printf("  %s⚠ Direct == comparison breaks when errors are wrapped with %%w%s\n", yellow, reset)
	fmt.Printf("  %s✔ Sentinels are great for expected failure modes (io.EOF, sql.ErrNoRows)%s\n", green, reset)
	fmt.Printf("  %s⚠ Don't create sentinels for errors that need structured data — use custom types%s\n", yellow, reset)
}
