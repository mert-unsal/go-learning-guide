package error_handling
import (
"errors"
"testing"
)
func TestDivide(t *testing.T) {
r, err := DivideSolution(10, 2)
if err != nil || r != 5 { t.Errorf("Divide(10,2)=(%v,%v) want (5,nil)", r, err) }
_, err = DivideSolution(10, 0)
if err == nil { t.Error("Divide by zero should error") }
}
func TestValidate(t *testing.T) {
if err := ValidateSolution("Alice"); err != nil {
t.Errorf("Validate(Alice) unexpected error: %v", err)
}
err := ValidateSolution("")
if err == nil { t.Fatal("Validate empty should error") }
var ve *ValidationErrorSolution
if !errors.As(err, &ve) { t.Error("should be *ValidationErrorSolution") }
}
func TestSafeGet(t *testing.T) {
m := map[string]int{"x": 10}
v, err := SafeGetSolution(m, "x")
if err != nil || v != 10 { t.Errorf("SafeGet(x)=(%d,%v) want (10,nil)", v, err) }
_, err = SafeGetSolution(m, "missing")
if err == nil { t.Error("missing key should error") }
}
func TestFindUser(t *testing.T) {
u, err := FindUserSolution(1)
if err != nil || u != "User1" { t.Errorf("FindUser(1)=(%q,%v) want (User1,nil)", u, err) }
_, err = FindUserSolution(0)
if !errors.Is(err, ErrUserNotFound) { t.Errorf("id=0 should be ErrUserNotFound, got %v", err) }
_, err = FindUserSolution(999)
if !errors.Is(err, ErrAccessDenied) { t.Errorf("id=999 should be ErrAccessDenied, got %v", err) }
}
func TestWrapError(t *testing.T) {
orig := errors.New("original")
wrapped := WrapErrorSolution(orig, "context")
if !errors.Is(wrapped, orig) { t.Error("errors.Is should find original in wrapped") }
}