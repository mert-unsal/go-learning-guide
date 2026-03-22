package error_handling

import (
	"errors"
	"fmt"
	"strconv"
)

// ============================================================
// EXERCISES — 07 Error Handling
// ============================================================

// Divide Exercise 1: Divide returns error if b == 0
// LESSON: The most basic pattern. Return (value, error). The caller MUST check the error.
func Divide(a, b float64) (float64, error) {
	// TODO: implement
	if b == 0 {
		return b, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

// ValidationError Exercise 2: Custom error type
// LESSON: Attach structured data to errors when callers need to act on the details.
// The caller uses errors.As(err, &ve) to extract the *ValidationError.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	// TODO: return formatted "validation error: <Field> — <Message>"
	return fmt.Sprintf("validation error: %s — %s", e.Field, e.Message)
}

func Validate(name string) error {
	// TODO: return &ValidationError if name is empty, nil otherwise
	if name == "" {
		return &ValidationError{Field: "name", Message: "cannot be empty"}
	}
	return nil
}

// Exercise 3: Safe map access
// LESSON: Add context to errors so they're useful when they bubble up.
// Use fmt.Errorf to create an error with the key name embedded.
func SafeGet(m map[string]int, key string) (int, error) {
	// TODO: return value if key exists, descriptive error if not
	if value, found := m[key]; found {
		return value, nil
	}

	return 0, fmt.Errorf("key %q not found", key)
}

// ErrUserNotFound Exercise 4: Sentinel errors
// LESSON: Predefined errors let callers make decisions with errors.Is().
// Use var, not const — errors.New returns a pointer, each call is unique.
var ErrUserNotFound = errors.New("user not found")
var ErrAccessDenied = errors.New("access denied")

func FindUser(id int) (string, error) {
	// TODO: return ErrUserNotFound if id <= 0, ErrAccessDenied if id == 999,
	//       otherwise return "User<id>" with nil error
	if id <= 0 {
		return "", ErrUserNotFound
	}
	if id == 999 {
		return "", ErrAccessDenied
	}
	return "User" + strconv.Itoa(id), nil
}

// WrapError Exercise 5: Wrap errors with context using %w
// LESSON: %w (not %v!) wraps the error so errors.Is/errors.As can still find
// the original inside the chain. This is how you add context without losing identity.
func WrapError(err error, context string) error {
	// TODO: wrap err with context string using fmt.Errorf and %w
	return fmt.Errorf("%s: %w", context, err)
}
