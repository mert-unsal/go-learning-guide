// Package main demonstrates what interface values actually are at runtime.
//
// ============================================================
// 5. INTERFACE VALUES — WHAT THEY ACTUALLY ARE AT RUNTIME
// ============================================================
// An interface value is a TWO-WORD pair internally: (type, value).
//
// Under the hood, the Go runtime represents interfaces as:
//   - Empty interface (any/interface{}): runtime.eface = {_type *_type, data unsafe.Pointer}
//   - Non-empty interface:               runtime.iface = {tab *itab, data unsafe.Pointer}
//
// The itab (interface table) contains:
//   - the interface type descriptor
//   - the concrete type descriptor
//   - a table of method pointers for the concrete type's methods that match the interface
//   - itabs are cached globally — first call per (interface, concrete type) pair builds it
//
//   var s Stringer        → (nil,   nil  )   s == nil is TRUE
//   s = User{...}         → (User,  0xc0...)  s == nil is FALSE
//   s = (*User)(nil)      → (*User, nil  )   s == nil is FALSE  ← THE TRAP
//
// The trap: once a type is stored in the interface, the interface is no longer
// nil even if the pointer inside it is nil. Checking s == nil is NOT enough
// to know whether the underlying value is safe to use.
package main

import "fmt"

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

func main() {
	// Case 1: true nil interface — both slots are nil
	var s Stringer
	fmt.Println("nil interface:", s == nil) // true

	// Case 2: non-nil interface with a real value
	s = User{Name: "Bob", Age: 25}
	fmt.Println("non-nil interface:", s == nil) // false
	fmt.Println(s.String())

	// Case 3: THE TRAP — nil pointer stored inside a non-nil interface
	var u *User                                                  // u is nil
	s = u                                                        // s is now (*User, nil) — the type slot is filled
	fmt.Println("nil pointer in interface, s == nil:", s == nil) // FALSE!
	// s == nil is false, but calling s.String() would PANIC
	// because the method tries to dereference a nil pointer.
}
