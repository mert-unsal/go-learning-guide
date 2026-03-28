// Package main demonstrates type assertions and type switches.
//
// ============================================================
// 6. TYPE ASSERTIONS AND TYPE SWITCHES
// ============================================================
// When you have an interface value and need the concrete type back,
// use a type assertion (single type) or type switch (multiple types).
// Prefer type switches — they communicate intent more clearly.
//
// Under the hood:
//   - A type assertion checks the itab (or eface._type) to see if the
//     concrete type matches. It's a pointer comparison — very fast.
//   - A type switch is compiled into a series of these comparisons,
//     or a hash-based lookup for large switches.
//   - Neither involves reflection — it's all compile-time generated code.
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

type JSONFormatter struct{}
type TextFormatter struct{}

func (j JSONFormatter) Format(msg string) string {
	return fmt.Sprintf(`{"msg": %q}`, msg)
}

func (t TextFormatter) Format(msg string) string {
	return fmt.Sprintf("[TEXT] %s", msg)
}

type Formatter interface {
	Format(msg string) string
}

// Type assertion — when you need ONE specific type.
func asJSON(f Formatter) {
	jf, ok := f.(JSONFormatter)
	if !ok {
		fmt.Println("not a JSONFormatter")
		return
	}
	fmt.Println("JSON output:", jf.Format("hello"))
}

// Type switch — when you have several possibilities.
func describe(f Formatter, msg string) {
	switch v := f.(type) {
	case JSONFormatter:
		fmt.Println("using JSON:", v.Format(msg))
	case TextFormatter:
		fmt.Println("using text:", v.Format(msg))
	default:
		fmt.Printf("unknown formatter type: %T\n", v)
	}
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Type Assertions & Type Switches         %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ Under the hood: type assertion = itab pointer comparison (very fast)%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ No reflection involved — compile-time generated code%s\n", green, reset)
	fmt.Printf("  %s✔ Type switch compiles to a series of pointer comparisons or hash lookup%s\n\n", green, reset)

	// Type assertion — success path
	fmt.Printf("%s▸ Type Assertion — extracting ONE specific concrete type%s\n", cyan+bold, reset)
	var f Formatter = JSONFormatter{}
	fmt.Printf("  f = JSONFormatter{}  → iface = %s(type=JSONFormatter, data=...)%s\n", magenta, reset)
	fmt.Printf("  asJSON(f):  ")
	asJSON(f)
	fmt.Printf("  %s✔ Success: itab type matches JSONFormatter — assertion passes%s\n\n", green, reset)

	// Type assertion — failure path
	fmt.Printf("%s▸ Type Assertion — failure path (comma-ok pattern)%s\n", cyan+bold, reset)
	f = TextFormatter{}
	fmt.Printf("  f = TextFormatter{}  → iface = %s(type=TextFormatter, data=...)%s\n", magenta, reset)
	fmt.Printf("  asJSON(f):  ")
	asJSON(f)
	fmt.Printf("  %s⚠ Without comma-ok, a failed assertion PANICs: f.(JSONFormatter) → panic%s\n", yellow, reset)
	fmt.Printf("  %s✔ With comma-ok: jf, ok := f.(JSONFormatter) → ok=false, no panic%s\n\n", green, reset)

	// Type switch — multiple possibilities
	fmt.Printf("%s▸ Type Switch — matching against multiple concrete types%s\n", cyan+bold, reset)
	f = JSONFormatter{}
	fmt.Printf("  describe(JSONFormatter, \"event\"): ")
	describe(f, "event")
	f = TextFormatter{}
	fmt.Printf("  describe(TextFormatter, \"event\"): ")
	describe(f, "event")

	fmt.Printf("\n  %s✔ Prefer type switch over chained type assertions — clearer intent%s\n", green, reset)
	fmt.Printf("  %s⚠ Both are O(1)-ish via itab comparison, but switch handles exhaustiveness%s\n", yellow, reset)
}
