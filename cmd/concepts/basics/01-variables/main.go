// Variables in Go — demonstrates declaration, initialization, and assignment.
//
// Every Go file starts with a package declaration.
// "main" package is special: it defines an executable program.
// Other packages are libraries.
//
// Imports bring in other packages. Unused imports are a COMPILE ERROR in Go.
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
	fmt.Printf("%s%s  Variables & Declaration Patterns        %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- var keyword (explicit) ---
	fmt.Printf("%s▸ var keyword — explicit declaration%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ 'var' declares a variable with an explicit type%s\n", green, reset)
	fmt.Printf("  %s✔ Uninitialized vars get their type's zero value — Go guarantees no garbage memory%s\n", green, reset)

	var a int     // zero value: 0
	var b string  // zero value: ""
	var c bool    // zero value: false
	var d float64 // zero value: 0.0

	fmt.Printf("  var a int     → %s%d%s\n", magenta, a, reset)
	fmt.Printf("  var b string  → %s%q%s\n", magenta, b, reset)
	fmt.Printf("  var c bool    → %s%t%s\n", magenta, c, reset)
	fmt.Printf("  var d float64 → %s%g%s\n\n", magenta, d, reset)

	// --- var with initializer ---
	fmt.Printf("%s▸ var with initializer%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ When you provide an initial value, the type can be inferred — but explicit is still valid%s\n", green, reset)

	var x int = 42
	var name string = "Gopher"
	fmt.Printf("  var x int = 42       → %s%d%s\n", magenta, x, reset)
	fmt.Printf("  var name string = .. → %s%s%s\n\n", magenta, name, reset)

	// --- Short variable declaration := (most common inside functions) ---
	// Type is INFERRED automatically
	fmt.Printf("%s▸ Short declaration := (idiomatic Go)%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ := declares AND initializes — type is inferred by the compiler%s\n", green, reset)
	fmt.Printf("  %s⚠ := only works inside functions, never at package level%s\n", yellow, reset)

	age := 25
	pi := 3.14159
	isGoFun := true
	fmt.Printf("  age := 25       → %s%d%s  (inferred %sint%s)\n", magenta, age, reset, dim, reset)
	fmt.Printf("  pi := 3.14159   → %s%g%s  (inferred %sfloat64%s, not float32!)\n", magenta, pi, reset, dim, reset)
	fmt.Printf("  isGoFun := true → %s%t%s  (inferred %sbool%s)\n\n", magenta, isGoFun, reset, dim, reset)

	// --- Multiple assignment ---
	fmt.Printf("%s▸ Multiple assignment%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Go evaluates ALL right-hand values before assigning to left-hand targets%s\n", green, reset)

	x1, y1 := 10, 20
	fmt.Printf("  x1, y1 := 10, 20 → x1=%s%d%s, y1=%s%d%s\n\n", magenta, x1, reset, magenta, y1, reset)

	// --- Swap values (Go's elegant way) ---
	fmt.Printf("%s▸ Swap values — Go's tuple assignment%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ No temp variable needed — right side is fully evaluated first%s\n", green, reset)

	x1, y1 = y1, x1
	fmt.Printf("  x1, y1 = y1, x1 → x1=%s%d%s, y1=%s%d%s\n\n", magenta, x1, reset, magenta, y1, reset)

	// --- Blank identifier _ (discard a value) ---
	fmt.Printf("%s▸ Blank identifier _%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ _ discards a value — commonly used to ignore one return from multi-return functions%s\n", green, reset)
	fmt.Printf("  %s⚠ Unused variables are a COMPILE ERROR in Go — _ is how you explicitly discard%s\n", yellow, reset)

	_, second := 100, 200
	fmt.Printf("  _, second := 100, 200 → second=%s%d%s  (100 is discarded)\n", magenta, second, reset)
}
