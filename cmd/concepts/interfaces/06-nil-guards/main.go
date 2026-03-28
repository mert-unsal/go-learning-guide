// Package main demonstrates how to guard against the nil pointer trap
// in interface values.
//
// ────────────────────────────────────────────────────────────
// HOW TO GUARD AGAINST THE NIL POINTER TRAP
// ────────────────────────────────────────────────────────────
//
// Three approaches, from last resort to best practice:
//
// APPROACH 1 — reflect.ValueOf
//   Use only when you receive an interface from outside code you don't control.
//   It checks whether the value stored inside the interface is itself nil.
//
// APPROACH 2 — type assertion
//   Clean, zero reflection overhead. Use when you know the possible concrete types.
//
// APPROACH 3 — fix it at the source (always the best option)
//   Never let a typed nil enter an interface. Return a true nil interface instead.
package main

import (
	"fmt"
	"reflect"
)

type Stringer interface {
	String() string
}

type User struct {
	Name string
	Age  int
}

func (u User) String() string {
	return fmt.Sprintf("%s (age %d)", u.Name, u.Age)
}

// Approach 1: reflect catches the nil pointer inside the interface.
func isTrulyNil(i any) bool {
	if i == nil {
		return true // fast path: interface itself is (nil, nil)
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface,
		reflect.Slice, reflect.Map,
		reflect.Chan, reflect.Func:
		return v.IsNil()
	}
	return false
}

// Approach 2: assert to the concrete type, then check the pointer.
func safeCall(s Stringer) {
	u, ok := s.(*User) // extract the concrete *User
	if !ok || u == nil {
		fmt.Println("  nil or wrong type — skipping")
		return
	}
	fmt.Println(" ", u.String())
}

// Approach 3 (BAD): returns (*User, nil) — caller's == nil check will FAIL.
func findUserBad(found bool) Stringer {
	var u *User
	if !found {
		return u // typed nil leaks into the interface → BUG
	}
	return &User{Name: "Alice", Age: 30}
}

// Approach 3 (GOOD): returns a true (nil, nil) interface — == nil works correctly.
func findUserGood(found bool) Stringer {
	if !found {
		return nil // untyped nil → true nil interface
	}
	return &User{Name: "Alice", Age: 30}
}

func main() {
	var u *User
	var s Stringer = u // trap: s holds (*User, nil)

	fmt.Println("── reflect guard ──")
	fmt.Println("s == nil:      ", s == nil)       // false — the trap
	fmt.Println("isTrulyNil(s): ", isTrulyNil(s))  // true  — correct
	fmt.Println("isTrulyNil(42):", isTrulyNil(42)) // false

	fmt.Println("── type assertion guard ──")
	safeCall(s)                            // skipped safely
	safeCall(User{Name: "Carol", Age: 28}) // printed

	fmt.Println("── fix at the source ──")
	bad := findUserBad(false)
	good := findUserGood(false)
	fmt.Println("bad  == nil:", bad == nil)  // false — the bug
	fmt.Println("good == nil:", good == nil) // true  — correct
}
