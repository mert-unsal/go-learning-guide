package error_handling

import (
	"errors"
	"fmt"
	"strconv"
)

// ============================================================
// EXERCISES — 07 Error Handling
// ============================================================
//
// Exercises 1-5:  Producing errors (return, custom types, sentinels, wrapping)
// Exercises 6-11: Consuming errors (extracting, unwrap chains, join, recover, multi-layer, custom Is)

// Divide Exercise 1: Divide returns error if b == 0
// LESSON: The most basic pattern. Return (value, error). The caller MUST check the error.
func Divide(a, b float64) (float64, error) {
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
	if name == "" {
		return &ValidationError{Field: "name", Message: "cannot be empty"}
	}
	return nil
}

// SafeGet Exercise 3: Safe map access
// LESSON: Add context to errors so they're useful when they bubble up.
// Use fmt.Errorf to create an error with the key name embedded.
func SafeGet(m map[string]int, key string) (int, error) {
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
	return fmt.Errorf("%s: %w", context, err)
}

// ============================================================
// PART B — Advanced Error Patterns (Consuming Errors)
// ============================================================

// ClassifyError Exercise 6: Extract structured data from error chains using errors.As
//
// LESSON: In production, you receive errors from lower layers and need to
// extract structured data to decide what to do. errors.As walks the entire
// wrapping chain until it finds a matching type.
//
// The function receives an error that may (or may not) contain a *ValidationError
// somewhere in its wrapping chain. If found, return (field, message, true).
// If not found, return ("", "", false).
//
// Example:
//
//	inner := &ValidationError{Field: "email", Message: "invalid format"}
//	wrapped := fmt.Errorf("signup failed: %w", inner)
//	field, msg, ok := ClassifyError(wrapped)
//	// field="email", msg="invalid format", ok=true
func ClassifyError(err error) (field string, message string, ok bool) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve.Field, ve.Message, true
	}
	return "", "", false
}

// RetryableError Exercise 7: Custom error type with Unwrap method
//
// LESSON: When your custom error type wraps another error, you must implement
// the Unwrap() method so errors.Is/errors.As can see through it. Without
// Unwrap(), the chain breaks and callers can't find the original error.
//
// Implement:
//   - RetryableError struct with fields: Err error, Retryable bool
//   - Error() string → "retryable: <inner error>" or "permanent: <inner error>"
//   - Unwrap() error → returns the inner Err
type RetryableError struct {
	Err       error
	Retryable bool
}

func (e *RetryableError) Error() string {
	// TODO: return "retryable: <Err>" or "permanent: <Err>" based on Retryable flag
	if e.Retryable {
		return fmt.Sprintf("retryable: %s", e.Err)
	}
	return fmt.Sprintf("permanent: %s", e.Err)
}

func (e *RetryableError) Unwrap() error {
	// TODO: return the wrapped inner error
	return e.Err
}

// NewRetryableError wraps an error as retryable or permanent.
func NewRetryableError(err error, retryable bool) error {
	return &RetryableError{Err: err, Retryable: retryable}
}

// CollectErrors Exercise 8: Aggregate multiple errors with errors.Join (Go 1.20+)
//
// LESSON: In production, batch operations produce multiple errors. You can't
// return just the first one — you need all of them. errors.Join combines
// multiple errors into one. The joined error's Error() concatenates all
// messages with newlines, and errors.Is/errors.As check each one.
//
// The function runs a validator function against each string in the input.
// Collect ALL errors (don't stop at the first), then return them joined.
// If no errors, return nil.
//
// Example:
//
//	validator := func(s string) error {
//	    if s == "" { return errors.New("empty string") }
//	    return nil
//	}
//	err := CollectErrors([]string{"a", "", "b", ""}, validator)
//	// err.Error() == "empty string\nempty string"
//	// errors.Is(err, errors.New("empty string")) — careful, this is false!
//	// (each errors.New creates a unique pointer)
func CollectErrors(inputs []string, validate func(string) error) error {
	// TODO: run validate on each input, collect errors, return errors.Join(...)
	var errs []error
	for _, input := range inputs {
		if err := validate(input); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

// SafeDivide Exercise 9: Recover from panic and convert to error
//
// LESSON: Panics are for programmer bugs, not expected errors. But at
// goroutine boundaries (HTTP handlers, worker loops), you MUST recover
// or the entire program crashes. The pattern: defer + recover + convert to error.
//
// SafeDivide performs integer division a/b. If b==0, Go panics with
// "runtime error: integer divide by zero". Catch this with recover and
// return it as an error instead.
//
// Rules:
//   - Use defer + anonymous function + recover
//   - Use named return values so the deferred function can set err
//   - Do NOT check b==0 yourself — let the panic happen and recover from it
func SafeDivide(a, b int) (result int, err error) {
	// TODO: defer a recovery function, then do a/b
	defer func() {
		if panicValue := recover(); panicValue != nil {
			err = fmt.Errorf("panic: %v", panicValue)
		}
	}()
	return a / b, nil
}

// MultiLayerOperation Exercise 10: Multi-layer error wrapping chain
//
// LESSON: In enterprise Go, errors pass through multiple layers:
// repository → service → handler. Each layer adds context with %w.
// The caller at the top must be able to find the original sentinel
// error through the entire chain using errors.Is.
//
// Simulate a 3-layer call stack:
//
//	Repository(id) → if id <= 0, return ErrUserNotFound
//	Service(id)    → calls Repository, wraps error with "service: get user %d: %w"
//	Handler(id)    → calls Service, wraps error with "handler: process request: %w"
//
// The test will verify: errors.Is(Handler(-1), ErrUserNotFound) == true
// even though the error has been wrapped twice.
func Repository(id int) (string, error) {
	// TODO: return ErrUserNotFound if id <= 0, return "User<id>" otherwise
	if id <= 0 {
		return "", fmt.Errorf("repository: UserID : %d, %w", id, ErrUserNotFound)
	}
	return "User" + strconv.Itoa(id), nil
}

func Service(id int) (string, error) {
	// TODO: call Repository, wrap error with service context using %w
	if userFoundInfo, err := Repository(id); err != nil {
		return "", fmt.Errorf("service: get user %d: %w", id, err)
	} else {
		return userFoundInfo, nil
	}
}

func Handler(id int) (string, error) {
	// TODO: call Service, wrap error with handler context using %w
	if serviceUser, err := Service(id); err != nil {
		return "", fmt.Errorf("handler: process request: %w", err)
	} else {
		return serviceUser, nil
	}
}

// StatusCodeError Exercise 11: Custom Is() method for flexible error matching
//
// LESSON: Sometimes you want errors.Is to match on VALUE, not identity.
// By default, errors.Is uses == comparison (pointer equality for sentinel errors).
// If you implement Is(target error) bool on your error type, errors.Is calls
// YOUR method instead, giving you control over matching logic.
//
// Implement StatusCodeError with:
//   - Fields: Code int, Msg string
//   - Error() string → "<Code>: <Msg>"
//   - Is(target error) bool → returns true if the target is also *StatusCodeError
//     with the same Code (ignore Msg). This lets you match by status code
//     regardless of the message.
//
// Example:
//
//	err1 := &StatusCodeError{Code: 404, Msg: "user not found"}
//	err2 := &StatusCodeError{Code: 404, Msg: "order not found"}
//	errors.Is(err1, err2) // true — same code
type StatusCodeError struct {
	Code int
	Msg  string
}

func (e *StatusCodeError) Error() string {
	// TODO: return "<Code>: <Msg>"
	return fmt.Sprintf("%d: %s", e.Code, e.Msg)
}

func (e *StatusCodeError) Is(target error) bool {
	// TODO: return true if target is *StatusCodeError with same Code
	if e.Code == target.(*StatusCodeError).Code {
		return true
	}
	return false
}
