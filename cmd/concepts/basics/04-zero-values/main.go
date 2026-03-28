// Zero Values in Go — every variable is ALWAYS initialized to its zero value.
// This prevents undefined behavior common in C/C++.
//
// Zero value reference table:
//
//	int, float     → 0
//	bool           → false
//	string         → ""
//	pointer        → nil
//	slice          → nil
//	map            → nil
//	channel        → nil
//	function       → nil
//	interface      → nil
//	struct         → all fields set to their zero values
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

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Zero Values — Go's Safety Guarantee     %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s✔ Every variable in Go is initialized to its zero value — no undefined memory%s\n", green, reset)
	fmt.Printf("%s✔ This eliminates an entire class of bugs common in C/C++ (reading uninitialized vars)%s\n", green, reset)
	fmt.Printf("%s✔ \"Make the zero value useful\" is a core Go proverb — design your types accordingly%s\n\n", green, reset)

	// --- Scalar zero values ---
	fmt.Printf("%s▸ Scalar types%s\n", cyan+bold, reset)

	var i int
	var f float64
	var b bool
	var s string
	fmt.Printf("  var i int     → %s%d%s\n", magenta, i, reset)
	fmt.Printf("  var f float64 → %s%g%s\n", magenta, f, reset)
	fmt.Printf("  var b bool    → %s%t%s\n", magenta, b, reset)
	fmt.Printf("  var s string  → %s%q%s  %s(empty string, not nil!)%s\n\n", magenta, s, reset, dim, reset)

	// --- Reference-type zero values ---
	fmt.Printf("%s▸ Reference types%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ nil slice/map/channel are valid zero values, but behave differently%s\n", yellow, reset)

	var sl []int
	var m map[string]int
	var ch chan int
	var p *int
	var fn func()
	var iface interface{}
	fmt.Printf("  var sl []int           → %s%v%s  (nil — append works, but len=0)\n", magenta, sl, reset)
	fmt.Printf("  var m map[string]int   → %s%v%s  (nil — read returns zero, but %swrite panics!%s)\n", magenta, m, reset, red+bold, reset)
	fmt.Printf("  var ch chan int        → %s%v%s  (nil — send/recv blocks forever)\n", magenta, ch, reset)
	fmt.Printf("  var p *int            → %s%v%s  (nil pointer)\n", magenta, p, reset)
	fmt.Printf("  var fn func()         → %s%v%s  (nil function)\n", magenta, fn, reset)
	fmt.Printf("  var iface interface{} → %s%v%s  (nil interface — both type and value are nil)\n\n", magenta, iface, reset)

	// --- Struct zero value ---
	fmt.Printf("%s▸ Struct zero values%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ All fields recursively set to their own zero values%s\n", green, reset)

	type Point struct {
		X, Y float64
		Label string
	}
	var pt Point
	fmt.Printf("  var pt Point → %s%+v%s\n", magenta, pt, reset)
	fmt.Printf("  %s✔ A zero-value Point{0,0,\"\"} is valid and usable — that's good design%s\n\n", green, reset)

	// --- Practical takeaways ---
	fmt.Printf("%s▸ Why this matters in production%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ sync.Mutex zero value is an unlocked mutex — ready to use, no constructor needed%s\n", green, reset)
	fmt.Printf("  %s✔ bytes.Buffer zero value is an empty buffer — just start writing%s\n", green, reset)
	fmt.Printf("  %s⚠ Always check: can your struct be used without an explicit constructor?%s\n", yellow, reset)
}
