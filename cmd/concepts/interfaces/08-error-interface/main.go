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

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

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
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  error Is an Interface                   %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Built-in error: type error interface { Error() string }%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Any type with Error() string satisfies error — no registration%s\n", green, reset)
	fmt.Printf("  %s✔ The most ubiquitous interface in Go — one method, universally useful%s\n", green, reset)
	fmt.Printf("  %s✔ Custom error types carry structured data (resource, user, action)%s\n\n", green, reset)

	fmt.Printf("%s▸ NotFoundError and PermissionError both satisfy error implicitly%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Functions return error interface — callers don't import concrete types%s\n\n", green, reset)

	for _, call := range []struct {
		id   int
		user string
	}{{-1, "alice"}, {1, "guest"}, {1, "alice"}} {
		err := findRecord(call.id, call.user)

		fmt.Printf("%s▸ findRecord(id=%s%d%s, user=%s%q%s)%s\n",
			cyan+bold, magenta, call.id, cyan+bold, magenta, call.user, cyan+bold, reset)

		if err == nil {
			fmt.Printf("  %s✔ No error — record found successfully%s\n\n", green, reset)
			continue
		}

		// Use type switch to handle specific error kinds
		switch e := err.(type) {
		case *NotFoundError:
			fmt.Printf("  %s⚠ NotFoundError: resource=%s%q%s%s\n", yellow, magenta, e.Resource, yellow, reset)
			fmt.Printf("  %sdim  Type switch matched *NotFoundError — the itab check%s\n", dim, reset)
		case *PermissionError:
			fmt.Printf("  %s⚠ PermissionError: user=%s%q%s action=%s%q%s%s\n",
				yellow, magenta, e.User, yellow, magenta, e.Action, yellow, reset)
			fmt.Printf("  %sdim  Type switch matched *PermissionError — structured data preserved%s\n", dim, reset)
		default:
			fmt.Printf("  %sunexpected error: %v%s\n", red, err, reset)
		}
		fmt.Println()
	}

	fmt.Printf("  %s⚠ In production: use errors.Is() for sentinel errors, errors.As() for typed errors%s\n", yellow, reset)
	fmt.Printf("  %s✔ Log errors ONCE at top level — don't log at every layer%s\n", green, reset)
}
