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
