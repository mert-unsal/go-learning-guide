package error_handling
import (
"errors"
"fmt"
)
// ============================================================
// EXERCISES — 07 Error Handling
// ============================================================
// Exercise 1: Divide returns error if b == 0
func Divide(a, b float64) (float64, error) {
// TODO: return error "cannot divide by zero" when b == 0
return 0, nil
}
// Exercise 2: Custom error type
type ValidationError struct {
Field   string
Message string
}
func (e *ValidationError) Error() string {
return fmt.Sprintf("validation error: %s — %s", e.Field, e.Message)
}
func Validate(name string) error {
// TODO: return *ValidationError if name is empty
return nil
}
// Exercise 3: Safe map access
func SafeGet(m map[string]int, key string) (int, error) {
// TODO: return fmt.Errorf if key missing
return 0, nil
}
// Exercise 4: Sentinel errors (use Ex prefix to avoid redeclaring)
var ErrUserNotFound = errors.New("user not found")
var ErrAccessDenied = errors.New("access denied")
func FindUser(id int) (string, error) {
// TODO: id<=0 → ErrUserNotFound, id==999 → ErrAccessDenied, else "User{id}"
return "", nil
}
// Exercise 5: Wrap errors with context using %w
func WrapError(err error, context string) error {
// TODO: fmt.Errorf("%s: %w", context, err)
return nil
}