// Package main demonstrates the error interface pattern.
//
// ============================================================
// 7. REAL-WORLD PATTERN: error IS an interface
// ============================================================
// The built-in error is just: type error interface { Error() string }
// Any type with Error() string is an error — no registration needed.
// This is Go's interfaces working exactly as designed.
//
// Key insights:
//   - error is the most ubiquitous interface in Go — one method, universally useful.
//   - Custom error types let you carry structured data (resource name, user, action).
//   - Callers can use type switches or errors.As() to handle specific error kinds.
//   - Functions return the error interface — callers don't need to import concrete error types.
package main

import "fmt"

type NotFoundError struct{ Resource string }
type PermissionError struct {
	User   string
	Action string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

func (e *PermissionError) Error() string {
	return fmt.Sprintf("user %q cannot %s", e.User, e.Action)
}

// findRecord returns the error interface — callers don't need to import
// NotFoundError or PermissionError to handle errors generically.
func findRecord(id int, user string) error {
	if id <= 0 {
		return &NotFoundError{Resource: fmt.Sprintf("record#%d", id)}
	}
	if user == "guest" {
		return &PermissionError{User: user, Action: "read records"}
	}
	return nil
}

func main() {
	for _, call := range []struct {
		id   int
		user string
	}{{-1, "alice"}, {1, "guest"}, {1, "alice"}} {
		err := findRecord(call.id, call.user)
		if err == nil {
			fmt.Printf("id=%d user=%s → ok\n", call.id, call.user)
			continue
		}
		// Use type switch to handle specific error kinds
		switch e := err.(type) {
		case *NotFoundError:
			fmt.Printf("not found: %s\n", e.Resource)
		case *PermissionError:
			fmt.Printf("denied: %s tried to %s\n", e.User, e.Action)
		default:
			fmt.Println("unexpected error:", err)
		}
	}
}
