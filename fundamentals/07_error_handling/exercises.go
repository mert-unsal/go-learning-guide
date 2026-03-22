package error_handling

import (
	"errors"
	"fmt"
)

// ============================================================
// EXERCISES — 07 Error Handling
// ============================================================

// Exercise 1: Divide returns error if b == 0
// LESSON: The most basic pattern. Return (value, error). The caller MUST check the error.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

// Exercise 2: Custom error type
// LESSON: Attach structured data to errors when callers need to act on the details.
// The caller uses errors.As(err, &ve) to extract the *ValidationError.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s — %s", e.Field, e.Message)
}

func Validate(name string) error {
	if name == "" {
		// Return a pointer — the Error() method has a pointer receiver
		return &ValidationError{Field: "name", Message: "cannot be empty"}
	}
	return nil
}

// Exercise 3: Safe map access
// LESSON: Add context to errors so they're useful when they bubble up.
// fmt.Errorf("context: ...") is Level 1. No wrapping needed here since
// there is no original error to preserve.
func SafeGet(m map[string]int, key string) (int, error) {
	v, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("key %q not found in map", key)
	}
	return v, nil
}

// Exercise 4: Sentinel errors
// LESSON: Predefined errors let callers make decisions with errors.Is().
// Use var, not const — errors.New returns a pointer, each call is unique.
var ErrUserNotFound = errors.New("user not found")
var ErrAccessDenied = errors.New("access denied")

func FindUser(id int) (string, error) {
	if id <= 0 {
		return "", ErrUserNotFound
	}
	if id == 999 {
		return "", ErrAccessDenied
	}
	return fmt.Sprintf("User%d", id), nil
}

// Exercise 5: Wrap errors with context using %w
// LESSON: %w (not %v!) wraps the error so errors.Is/errors.As can still find
// the original inside the chain. This is how you add context without losing identity.
func WrapError(err error, context string) error {
	return fmt.Errorf("%s: %w", context, err)
}
