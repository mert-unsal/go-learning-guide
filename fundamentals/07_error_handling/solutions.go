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
if id <= 0 { return "", ErrUserNotFound }
if id == 999 { return "", ErrAccessDenied }
return "User" + strconv.Itoa(id), nil
}
func WrapErrorSolution(err error, context string) error {
return fmt.Errorf("%s: %w", context, err)
}