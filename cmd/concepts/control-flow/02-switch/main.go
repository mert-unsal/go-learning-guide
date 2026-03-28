// Switch in Go — demonstrates all switch variants including fallthrough.
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
	fmt.Printf("%s%s  Switch in Go — All Variants            %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// Basic switch — no 'break' needed, each case breaks automatically
	fmt.Printf("%s▸ Basic Switch (auto-break, multi-value cases)%s\n", cyan+bold, reset)
	fmt.Printf("  %sUnlike C/Java, each case breaks automatically — no fall-through by default%s\n", dim, reset)
	day := "Monday"
	fmt.Printf("  day = %s%q%s\n", magenta, day, reset)
	switch day {
	case "Saturday", "Sunday": // multiple values per case
		fmt.Printf("  %s✔ Matched \"Saturday\", \"Sunday\" → Weekend!%s\n", green, reset)
	case "Monday":
		fmt.Printf("  %s✔ Matched \"Monday\" → Start of the work week%s\n", green, reset)
	default:
		fmt.Printf("  Default → Weekday\n")
	}
	fmt.Printf("  %s✔ Multiple values per case with comma — replaces long || chains%s\n", green, reset)

	// Switch with no expression (acts like if/else chain)
	fmt.Printf("\n%s▸ Expressionless Switch (acts like if/else chain)%s\n", cyan+bold, reset)
	fmt.Printf("  %sSyntax: switch { case <bool>: } — each case is a boolean expression%s\n", dim, reset)
	x := 15
	fmt.Printf("  x = %s%d%s\n", magenta, x, reset)
	switch {
	case x < 0:
		fmt.Printf("  x < 0 → negative\n")
	case x == 0:
		fmt.Printf("  x == 0 → zero\n")
	case x > 0 && x < 10:
		fmt.Printf("  x > 0 && x < 10 → small positive\n")
	default:
		fmt.Printf("  %s✔ default branch → large positive (x=%d ≥ 10)%s\n", green, x, reset)
	}

	// Switch with initializer
	fmt.Printf("\n%s▸ Switch with Init Statement%s\n", cyan+bold, reset)
	fmt.Printf("  %sSyntax: switch <init>; { } — variable scoped to switch block%s\n", dim, reset)
	switch n := 42; {
	case n < 0:
		fmt.Printf("  n < 0 → negative\n")
	case n < 100:
		fmt.Printf("  %s✔ n := %s%d%s%s → n < 100 matched%s\n", green, magenta, n, reset, green, reset)
	default:
		fmt.Printf("  n >= 100\n")
	}
	fmt.Printf("  %s✔ Same scoping benefit as if-init — n not visible outside switch%s\n", green, reset)

	// fallthrough: explicitly continue to next case (rare in practice)
	fmt.Printf("\n%s▸ Fallthrough (explicit, unconditional)%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ fallthrough is UNCONDITIONAL — runs next case body regardless of condition%s\n", yellow, reset)
	fmt.Printf("  %sswitch 2 { case 1: ... case 2: fallthrough; case 3: ... case 4: ... }%s\n", dim, reset)
	switch 2 {
	case 1:
		fmt.Printf("  case 1: %s\"one\"%s + fallthrough\n", magenta, reset)
		fallthrough
	case 2:
		fmt.Printf("  case 2: %s\"two\"%s ← matched, + fallthrough\n", magenta, reset)
		fallthrough // executes next case even if condition doesn't match
	case 3:
		fmt.Printf("  case 3: %s\"three\"%s ← executed via fallthrough (condition NOT checked!)\n", magenta, reset)
	case 4:
		fmt.Printf("  case 4: %s\"four\"%s\n", magenta, reset) // this will NOT print (fallthrough stops)
	}
	fmt.Printf("  %s⚠ fallthrough only advances ONE case — it does not cascade like C%s\n", yellow, reset)
	fmt.Printf("  %s⚠ Rarely used in production Go — prefer explicit logic over fallthrough%s\n", yellow, reset)
}
