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

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Interface Values — Runtime Internals    %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Non-empty interface = runtime.iface{tab *itab, data unsafe.Pointer}%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ itab holds: interface type + concrete type + method pointers%s\n", green, reset)
	fmt.Printf("  %s✔ itabs are cached globally — first use builds it, subsequent calls reuse%s\n\n", green, reset)

	// Case 1: true nil interface — both slots are nil
	fmt.Printf("%s▸ Case 1: True nil interface — both (type, value) are nil%s\n", cyan+bold, reset)
	var s Stringer
	fmt.Printf("  var s Stringer          → iface = %s(nil, nil)%s\n", magenta, reset)
	fmt.Printf("  s == nil                → %s%v%s\n", magenta, s == nil, reset)
	fmt.Printf("  %s✔ Both type and data pointers are nil — this is a true nil%s\n\n", green, reset)

	// Case 2: non-nil interface with a real value
	fmt.Printf("%s▸ Case 2: Interface holds a concrete value%s\n", cyan+bold, reset)
	s = User{Name: "Bob", Age: 25}
	fmt.Printf("  s = User{\"Bob\", 25}     → iface = %s(type=User, data=0x...)%s\n", magenta, reset)
	fmt.Printf("  s == nil                → %s%v%s\n", magenta, s == nil, reset)
	fmt.Printf("  s.String()              → %s%q%s\n\n", magenta, s.String(), reset)

	// Case 3: THE TRAP — nil pointer stored inside a non-nil interface
	fmt.Printf("%s▸ Case 3: THE NIL INTERFACE TRAP%s\n", cyan+bold, reset)
	var u *User // u is nil
	s = u       // s is now (*User, nil) — the type slot is filled
	fmt.Printf("  var u *User             → u = %s%v%s\n", magenta, u, reset)
	fmt.Printf("  s = u                   → iface = %s(type=*User, data=nil)%s\n", magenta, reset)
	fmt.Printf("  s == nil                → %s%v%s  ← %sTHIS IS THE TRAP!%s\n", magenta, s == nil, reset, red+bold, reset)
	fmt.Printf("  %s⚠ The type slot (*User) is non-nil, so s != nil%s\n", yellow, reset)
	fmt.Printf("  %s⚠ But calling s.String() would PANIC — nil pointer dereference%s\n", yellow, reset)
	fmt.Printf("  %s⚠ Fix: never let a typed nil leak into an interface — return untyped nil%s\n", yellow, reset)
}
