// If/Else in Go — demonstrates basic conditionals and init statements.
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
	fmt.Printf("%s%s  If/Else in Go — Conditionals & Init    %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// Basic if/else
	fmt.Printf("%s▸ Basic If / Else If / Else%s\n", cyan+bold, reset)
	fmt.Printf("  %sGo requires braces — no single-line if without them%s\n", dim, reset)
	x := 42
	fmt.Printf("  x = %s%d%s\n", magenta, x, reset)
	if x > 0 {
		fmt.Printf("  %s✔ x > 0 → \"positive\" branch taken%s\n", green, reset)
	} else if x < 0 {
		fmt.Printf("  x < 0 → \"negative\" branch taken\n")
	} else {
		fmt.Printf("  x == 0 → \"zero\" branch taken\n")
	}

	// KEY FEATURE: if with initialization statement
	fmt.Printf("\n%s▸ If with Init Statement (Go-specific feature)%s\n", cyan+bold, reset)
	fmt.Printf("  %sSyntax: if <init>; <condition> { }%s\n", dim, reset)
	fmt.Printf("  %s✔ The init variable is scoped ONLY to the if/else block%s\n", green, reset)
	if n := 10; n%2 == 0 {
		fmt.Printf("  n := %s%d%s → n%%%s2 == 0%s → %s\"even\"%s\n", magenta, n, reset, bold, reset, green, reset)
	} else {
		fmt.Printf("  n := %s%d%s → n%%%s2 != 0%s → %s\"odd\"%s\n", magenta, n, reset, bold, reset, green, reset)
	}
	fmt.Printf("  %s⚠ n is NOT accessible outside the if/else block — compile error if tried%s\n", yellow, reset)
	// fmt.Println(n) // COMPILE ERROR: n is not accessible here

	fmt.Printf("\n%s✔ Unlike C/Java, Go's if-init scoping prevents variable leakage%s\n", green, reset)
	fmt.Printf("%s✔ Use this pattern for error checks: if err := doThing(); err != nil { }%s\n", green, reset)
}
