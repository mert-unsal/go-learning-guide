// Defer Argument Evaluation — a very common interview question!
//
// RULE: Arguments are evaluated at DEFER time, not when the deferred call executes.
// RULE: Defer CAN modify named return values.
// This is useful for wrapping errors or ensuring consistent cleanup.
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

// deferWithNamedReturn demonstrates that defer can modify named return values.
// The deferred closure overwrites whatever was returned.
func deferWithNamedReturn() (result string) {
	defer func() {
		fmt.Printf("    defer closure runs: overwriting result %s%q%s → %s%q%s\n",
			yellow, result, reset, green, "modified by defer", reset)
		result = "modified by defer" // overwrites whatever was returned
	}()
	fmt.Printf("    return %s%q%s is assigned to named return 'result'\n", magenta, "original", reset)
	return "original" // this value gets overwritten by defer above
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Defer Argument Evaluation (Interview!) %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Arguments Evaluated at Defer Time, NOT Execution Time%s\n", cyan+bold, reset)
	fmt.Printf("  %sThis is the #1 defer gotcha in interviews%s\n", dim, reset)

	x := 10
	fmt.Printf("  x = %s%d%s (initial value)\n", magenta, x, reset)
	fmt.Printf("  → %sdefer fmt.Println(x)%s — argument x=%s%d%s captured NOW\n", dim, reset, magenta, x, reset)
	defer fmt.Printf("  deferred x = %s%d%s ← %s✔ captured at defer time, NOT execution time%s\n", magenta, x, reset, green, reset)
	x = 99
	fmt.Printf("  x = %s%d%s (changed AFTER defer was registered)\n", magenta, x, reset)
	fmt.Printf("  %s✔ When defer runs, it will print x=%s10%s%s, not 99 — args were already evaluated%s\n", green, magenta, reset, green, reset)
	fmt.Printf("  %s⚠ To capture current value at execution time, use a closure instead:%s\n", yellow, reset)
	fmt.Printf("  %s  defer func() { fmt.Println(x) }() — closure reads x when it runs%s\n", dim, reset)

	fmt.Printf("\n%s▸ Defer Can Modify Named Return Values%s\n", cyan+bold, reset)
	fmt.Printf("  %sDeferred closures see named return vars — they can overwrite the result%s\n", dim, reset)
	fmt.Printf("  Calling deferWithNamedReturn():\n")
	result := deferWithNamedReturn()
	fmt.Printf("  Caller received: %s%q%s\n", magenta, result, reset)
	fmt.Printf("  %s✔ The deferred closure overwrote \"original\" → \"modified by defer\"%s\n", green, reset)
	fmt.Printf("  %s✔ Production use: defer func() { if err != nil { err = fmt.Errorf(\"wrap: %%w\", err) } }()%s\n", green, reset)

	fmt.Printf("\n%s── main() returning — deferred print fires below ──%s\n", blue+bold, reset)
}
