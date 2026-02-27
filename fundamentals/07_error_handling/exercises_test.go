package error_handling

import (
	"errors"
	"testing"
)

func TestDivide(t *testing.T) {
	r, err := Divide(10, 2)
	if err != nil || r != 5 {
		t.Errorf("❌ Divide(10,2) = (%v,%v), want (5,nil)", r, err)
	} else {
		t.Logf("✅ Divide(10,2) = %.1f", r)
	}

	_, err = Divide(10, 0)
	if err == nil {
		t.Error("❌ Divide(10,0) should return an error  ← Hint: check b == 0")
	} else {
		t.Logf("✅ Divide(10,0) returned error: %v", err)
	}
}

func TestValidate(t *testing.T) {
	if err := Validate("Alice"); err != nil {
		t.Errorf("❌ Validate(\"Alice\") unexpected error: %v", err)
	} else {
		t.Logf("✅ Validate(\"Alice\") = nil (valid)")
	}

	err := Validate("")
	if err == nil {
		t.Fatal("❌ Validate(\"\") should return error  ← Hint: return *ValidationError when empty")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Errorf("❌ error should be *ValidationError, got %T", err)
	} else {
		t.Logf("✅ Validate(\"\") = %v", err)
	}
}

func TestSafeGet(t *testing.T) {
	m := map[string]int{"x": 10}
	v, err := SafeGet(m, "x")
	if err != nil || v != 10 {
		t.Errorf("❌ SafeGet(m, \"x\") = (%d,%v), want (10,nil)", v, err)
	} else {
		t.Logf("✅ SafeGet(m, \"x\") = %d", v)
	}

	_, err = SafeGet(m, "missing")
	if err == nil {
		t.Error("❌ SafeGet missing key should return error  ← Hint: use the ok idiom v, ok := m[key]")
	} else {
		t.Logf("✅ SafeGet(m, \"missing\") returned error: %v", err)
	}
}

func TestFindUser(t *testing.T) {
	u, err := FindUser(1)
	if err != nil || u != "User1" {
		t.Errorf("❌ FindUser(1) = (%q,%v), want (\"User1\",nil)", u, err)
	} else {
		t.Logf("✅ FindUser(1) = %q", u)
	}

	_, err = FindUser(0)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("❌ FindUser(0) should be ErrUserNotFound, got %v", err)
	} else {
		t.Logf("✅ FindUser(0) → ErrUserNotFound")
	}

	_, err = FindUser(999)
	if !errors.Is(err, ErrAccessDenied) {
		t.Errorf("❌ FindUser(999) should be ErrAccessDenied, got %v", err)
	} else {
		t.Logf("✅ FindUser(999) → ErrAccessDenied")
	}
}

func TestWrapError(t *testing.T) {
	orig := errors.New("original")
	wrapped := WrapError(orig, "context")
	if !errors.Is(wrapped, orig) {
		t.Error("❌ errors.Is should find original in wrapped  ← Hint: use fmt.Errorf(\"%s: %w\", ...)")
	} else {
		t.Logf("✅ WrapError wraps correctly: %v", wrapped)
	}
}
