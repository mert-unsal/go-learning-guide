package error_handling

import (
	"errors"
	"fmt"
	"strconv"
)

// SOLUTIONS — 07 Error Handling
func DivideSolution(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

type ValidationErrorSolution struct {
	Field, Message string
}

func (e *ValidationErrorSolution) Error() string {
	return fmt.Sprintf("validation error: %s — %s", e.Field, e.Message)
}
func ValidateSolution(name string) error {
	if name == "" {
		return &ValidationErrorSolution{Field: "name", Message: "cannot be empty"}
	}
	return nil
}
func SafeGetSolution(m map[string]int, key string) (int, error) {
	val, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("key %q not found", key)
	}
	return val, nil
}
func FindUserSolution(id int) (string, error) {
	if id <= 0 {
		return "", ErrUserNotFound
	}
	if id == 999 {
		return "", ErrAccessDenied
	}
	return "User" + strconv.Itoa(id), nil
}
func WrapErrorSolution(err error, context string) error {
	return fmt.Errorf("%s: %w", context, err)
}

// ============================================================
// SOLUTIONS — Part B Advanced Error Patterns
// ============================================================

func ClassifyErrorSolution(err error) (field string, message string, ok bool) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve.Field, ve.Message, true
	}
	return "", "", false
}

func RetryableErrorSolution(err error, retryable bool) string {
	if retryable {
		return fmt.Sprintf("retryable: %s", err.Error())
	}
	return fmt.Sprintf("permanent: %s", err.Error())
}

func CollectErrorsSolution(inputs []string, validate func(string) error) error {
	var errs []error
	for _, input := range inputs {
		if err := validate(input); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func SafeDivideSolution(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	return a / b, nil
}

func RepositorySolution(id int) (string, error) {
	if id <= 0 {
		return "", ErrUserNotFound
	}
	return fmt.Sprintf("User%d", id), nil
}

func ServiceSolution(id int) (string, error) {
	user, err := RepositorySolution(id)
	if err != nil {
		return "", fmt.Errorf("service: get user %d: %w", id, err)
	}
	return user, nil
}

func HandlerSolution(id int) (string, error) {
	user, err := ServiceSolution(id)
	if err != nil {
		return "", fmt.Errorf("handler: process request: %w", err)
	}
	return user, nil
}

func StatusCodeErrorSolution(code int, msg string) string {
	return fmt.Sprintf("%d: %s", code, msg)
}

func StatusCodeIsMatchSolution(code1, code2 int) bool {
	return code1 == code2
}
