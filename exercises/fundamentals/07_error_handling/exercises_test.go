package error_handling

import (
	"errors"
	"fmt"
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
		t.Errorf("❌ errors.Is should find original in wrapped  ← Hint: use fmt.Errorf(\"%%s: %%w\", ...)")
	} else {
		t.Logf("✅ WrapError wraps correctly: %v", wrapped)
	}
}

// ============================================================
// PART B — Advanced Error Pattern Tests
// ============================================================

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantField string
		wantMsg   string
		wantOK    bool
	}{
		{
			name:      "direct ValidationError",
			err:       &ValidationError{Field: "email", Message: "invalid"},
			wantField: "email", wantMsg: "invalid", wantOK: true,
		},
		{
			name:      "wrapped ValidationError",
			err:       fmt.Errorf("signup: %w", &ValidationError{Field: "age", Message: "too young"}),
			wantField: "age", wantMsg: "too young", wantOK: true,
		},
		{
			name:      "double wrapped ValidationError",
			err:       fmt.Errorf("handler: %w", fmt.Errorf("service: %w", &ValidationError{Field: "name", Message: "empty"})),
			wantField: "name", wantMsg: "empty", wantOK: true,
		},
		{
			name:   "non-ValidationError",
			err:    errors.New("something else"),
			wantOK: false,
		},
		{
			name:   "nil error",
			err:    nil,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, msg, ok := ClassifyError(tt.err)
			if ok != tt.wantOK {
				t.Errorf("❌ ClassifyError ok=%v, want %v  ← Hint: use errors.As(err, &ve)", ok, tt.wantOK)
			} else if ok && (field != tt.wantField || msg != tt.wantMsg) {
				t.Errorf("❌ ClassifyError = (%q,%q), want (%q,%q)", field, msg, tt.wantField, tt.wantMsg)
			} else {
				t.Logf("✅ ClassifyError(%v) → field=%q msg=%q ok=%v", tt.err, field, msg, ok)
			}
		})
	}
}

func TestRetryableError(t *testing.T) {
	inner := errors.New("connection refused")

	t.Run("retryable error message", func(t *testing.T) {
		re := NewRetryableError(inner, true)
		if re.Error() != "retryable: connection refused" {
			t.Errorf("❌ Error() = %q, want %q", re.Error(), "retryable: connection refused")
		} else {
			t.Logf("✅ retryable error: %v", re)
		}
	})

	t.Run("permanent error message", func(t *testing.T) {
		pe := NewRetryableError(inner, false)
		if pe.Error() != "permanent: connection refused" {
			t.Errorf("❌ Error() = %q, want %q", pe.Error(), "permanent: connection refused")
		} else {
			t.Logf("✅ permanent error: %v", pe)
		}
	})

	t.Run("Unwrap reveals inner error", func(t *testing.T) {
		re := NewRetryableError(inner, true)
		if !errors.Is(re, inner) {
			t.Errorf("❌ errors.Is should find inner error through Unwrap  ← Hint: implement Unwrap() error")
		} else {
			t.Logf("✅ errors.Is finds inner through Unwrap")
		}
	})

	t.Run("Unwrap through wrapping chain", func(t *testing.T) {
		re := NewRetryableError(inner, true)
		wrapped := fmt.Errorf("api call: %w", re)
		if !errors.Is(wrapped, inner) {
			t.Errorf("❌ errors.Is should find inner through double wrap")
		} else {
			t.Logf("✅ errors.Is finds inner through fmt.Errorf + Unwrap chain")
		}

		var target *RetryableError
		if !errors.As(wrapped, &target) {
			t.Errorf("❌ errors.As should find *RetryableError in chain")
		} else if !target.Retryable {
			t.Errorf("❌ should be retryable")
		} else {
			t.Logf("✅ errors.As extracted RetryableError, retryable=%v", target.Retryable)
		}
	})
}

func TestCollectErrors(t *testing.T) {
	var ErrEmpty = errors.New("empty string")
	var ErrShort = errors.New("too short")

	t.Run("multiple errors joined", func(t *testing.T) {
		validator := func(s string) error {
			if s == "" {
				return ErrEmpty
			}
			return nil
		}
		err := CollectErrors([]string{"a", "", "b", ""}, validator)
		if err == nil {
			t.Fatal("❌ should return error when validation fails  ← Hint: use errors.Join()")
		}
		if !errors.Is(err, ErrEmpty) {
			t.Errorf("❌ joined error should match ErrEmpty via errors.Is")
		} else {
			t.Logf("✅ CollectErrors joined: %v", err)
		}
	})

	t.Run("no errors returns nil", func(t *testing.T) {
		validator := func(s string) error { return nil }
		err := CollectErrors([]string{"a", "b", "c"}, validator)
		if err != nil {
			t.Errorf("❌ should return nil when all valid, got: %v", err)
		} else {
			t.Logf("✅ CollectErrors all valid → nil")
		}
	})

	t.Run("mixed error types in join", func(t *testing.T) {
		validator := func(s string) error {
			if s == "" {
				return ErrEmpty
			}
			if len(s) < 3 {
				return ErrShort
			}
			return nil
		}
		err := CollectErrors([]string{"hello", "", "ab", "world"}, validator)
		if err == nil {
			t.Fatal("❌ should return error")
		}
		if !errors.Is(err, ErrEmpty) {
			t.Errorf("❌ should contain ErrEmpty")
		}
		if !errors.Is(err, ErrShort) {
			t.Errorf("❌ should contain ErrShort")
		}
		t.Logf("✅ CollectErrors mixed: %v", err)
	})

	t.Run("empty input", func(t *testing.T) {
		validator := func(s string) error { return ErrEmpty }
		err := CollectErrors([]string{}, validator)
		if err != nil {
			t.Errorf("❌ empty input should return nil, got: %v", err)
		} else {
			t.Logf("✅ CollectErrors empty input → nil")
		}
	})
}

func TestSafeDivide(t *testing.T) {
	tests := []struct {
		name       string
		a, b       int
		wantResult int
		wantErr    bool
	}{
		{"normal division", 10, 2, 5, false},
		{"divide by zero", 10, 0, 0, true},
		{"negative division", -10, 2, -5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeDivide(tt.a, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ SafeDivide(%d,%d) should return error  ← Hint: defer+recover, don't check b==0", tt.a, tt.b)
				} else {
					t.Logf("✅ SafeDivide(%d,%d) recovered: %v", tt.a, tt.b, err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ SafeDivide(%d,%d) unexpected error: %v", tt.a, tt.b, err)
				} else if result != tt.wantResult {
					t.Errorf("❌ SafeDivide(%d,%d) = %d, want %d", tt.a, tt.b, result, tt.wantResult)
				} else {
					t.Logf("✅ SafeDivide(%d,%d) = %d", tt.a, tt.b, result)
				}
			}
		})
	}
}

func TestMultiLayerOperation(t *testing.T) {
	t.Run("error propagates through layers", func(t *testing.T) {
		_, err := Handler(-1)
		if err == nil {
			t.Fatal("❌ Handler(-1) should return error")
		}
		if !errors.Is(err, ErrUserNotFound) {
			t.Errorf("❌ errors.Is(err, ErrUserNotFound) should be true through 2 layers of wrapping  ← Hint: use %%w at each layer")
		} else {
			t.Logf("✅ errors.Is finds ErrUserNotFound through chain: %v", err)
		}
	})

	t.Run("error message shows full chain", func(t *testing.T) {
		_, err := Handler(-1)
		msg := err.Error()
		if len(msg) < 20 {
			t.Errorf("❌ error message too short — each layer should add context: %q", msg)
		} else {
			t.Logf("✅ Full error chain: %v", err)
		}
	})

	t.Run("success case", func(t *testing.T) {
		user, err := Handler(42)
		if err != nil {
			t.Errorf("❌ Handler(42) unexpected error: %v", err)
		} else {
			t.Logf("✅ Handler(42) = %q", user)
		}
	})
}

func TestStatusCodeError(t *testing.T) {
	t.Run("same code different message matches", func(t *testing.T) {
		err1 := &StatusCodeError{Code: 404, Msg: "user not found"}
		err2 := &StatusCodeError{Code: 404, Msg: "order not found"}
		if !errors.Is(err1, err2) {
			t.Errorf("❌ same Code should match via Is()  ← Hint: implement Is(target error) bool")
		} else {
			t.Logf("✅ 404 matches 404 regardless of message")
		}
	})

	t.Run("different code does not match", func(t *testing.T) {
		err1 := &StatusCodeError{Code: 404, Msg: "not found"}
		err2 := &StatusCodeError{Code: 500, Msg: "server error"}
		if errors.Is(err1, err2) {
			t.Errorf("❌ different Code should NOT match")
		} else {
			t.Logf("✅ 404 ≠ 500")
		}
	})

	t.Run("works through wrapping chain", func(t *testing.T) {
		inner := &StatusCodeError{Code: 503, Msg: "service unavailable"}
		wrapped := fmt.Errorf("gateway: %w", inner)
		target := &StatusCodeError{Code: 503, Msg: "different message"}
		if !errors.Is(wrapped, target) {
			t.Errorf("❌ should find matching Code through wrapping chain")
		} else {
			t.Logf("✅ errors.Is finds Code=503 through wrapping")
		}
	})

	t.Run("error message format", func(t *testing.T) {
		e := &StatusCodeError{Code: 404, Msg: "not found"}
		want := "404: not found"
		if e.Error() != want {
			t.Errorf("❌ Error() = %q, want %q", e.Error(), want)
		} else {
			t.Logf("✅ Error() = %q", e.Error())
		}
	})
}
