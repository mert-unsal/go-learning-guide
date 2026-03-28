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
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Nil Guards — Defending Against the Trap %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	var u *User
	var s Stringer = u // trap: s holds (*User, nil)

	// ── Approach 1: reflect guard ──
	fmt.Printf("%s▸ Approach 1: reflect.ValueOf — last resort, checks inner value%s\n", cyan+bold, reset)
	fmt.Printf("  s == nil:       %s%v%s  ← %sfalse! iface has type *User%s\n", magenta, s == nil, reset, red, reset)
	fmt.Printf("  isTrulyNil(s):  %s%v%s  ← %scorrect! reflect sees nil pointer inside%s\n", magenta, isTrulyNil(s), reset, green, reset)
	fmt.Printf("  isTrulyNil(42): %s%v%s  ← int is not a nillable kind\n", magenta, isTrulyNil(42), reset)
	fmt.Printf("  %s⚠ reflect has runtime cost — use only for unknown external input%s\n\n", yellow, reset)

	// ── Approach 2: type assertion guard ──
	fmt.Printf("%s▸ Approach 2: Type assertion — zero reflection overhead%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Extract concrete *User, then check pointer for nil%s\n", green, reset)
	fmt.Printf("  safeCall(nil *User) → ")
	safeCall(s) // skipped safely
	fmt.Printf("  safeCall(Carol)     → ")
	safeCall(&User{Name: "Carol", Age: 28}) // printed — pointer matches *User assertion
	fmt.Println()

	// ── Approach 3: fix at the source ──
	fmt.Printf("%s▸ Approach 3: Fix at the source — never let typed nil escape%s\n", cyan+bold, reset)
	bad := findUserBad(false)
	good := findUserGood(false)
	fmt.Printf("  findUserBad(false):  returns %s(*User)(nil)%s in interface\n", magenta, reset)
	fmt.Printf("    bad  == nil: %s%v%s  ← %sBUG! typed nil leaks into iface%s\n", magenta, bad == nil, reset, red+bold, reset)
	fmt.Printf("  findUserGood(false): returns %suntyped nil%s\n", magenta, reset)
	fmt.Printf("    good == nil: %s%v%s  ← %scorrect! true (nil, nil) interface%s\n\n", magenta, good == nil, reset, green, reset)

	fmt.Printf("  %s✔ Best practice: always return bare nil, never a typed nil pointer%s\n", green, reset)
	fmt.Printf("  %s⚠ 'return u' when u is *User nil → BUG. 'return nil' → correct%s\n", yellow, reset)
}
