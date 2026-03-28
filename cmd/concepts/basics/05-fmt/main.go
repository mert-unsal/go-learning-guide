// fmt Package — formatting verbs and output functions.
//
// Key formatting verbs:
//
//	%s  — string
//	%d  — integer (decimal)
//	%f  — float (default width)
//	%.2f — float with 2 decimal places
//	%e  — scientific notation
//	%t  — boolean
//	%T  — prints the TYPE of a variable
//	%v  — default format (works for anything)
//	%#v — Go-syntax representation
//	%q  — quoted string
//	%p  — pointer address
//	%b  — binary
//	%o  — octal
//	%x  — hexadecimal
//
// Stringer interface — any type with String() method works with %s and %v
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

// Animal demonstrates the Stringer interface.
type Animal struct {
	Name    string
	Sound   string
	Legs    int
}

func (a Animal) String() string {
	return fmt.Sprintf("%s (says %q, %d legs)", a.Name, a.Sound, a.Legs)
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  fmt Package — Formatting Verbs          %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	name := "Gopher"
	age := 5
	pi := 3.14159

	// Printf: formatted printing
	fmt.Printf("%s▸ Printf — formatted output verbs%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Printf never adds a newline — you must include \\n yourself%s\n", green, reset)
	fmt.Printf("  %s✔ Each verb starts with %% and specifies how to format the argument%s\n\n", green, reset)

	fmt.Printf("  %%s  string          → %s%s%s\n", magenta, name, reset)
	fmt.Printf("  %%d  integer         → %s%d%s\n", magenta, age, reset)
	fmt.Printf("  %%.2f float (2 dp)   → %s%.2f%s\n", magenta, pi, reset)
	fmt.Printf("  %%e  scientific      → %s%e%s\n", magenta, pi, reset)
	fmt.Printf("  %%t  boolean         → %s%t%s\n", magenta, true, reset)
	fmt.Printf("  %%T  TYPE of value   → %s%T%s  %s(reflection-based — useful for debugging)%s\n", magenta, name, reset, dim, reset)
	fmt.Printf("  %%v  default format  → %s%v%s  %s(the \"just print it\" verb)%s\n", magenta, age, reset, dim, reset)
	fmt.Printf("  %%#v Go syntax       → %s%#v%s  %s(shows type + value — great for debugging)%s\n", magenta, name, reset, dim, reset)
	fmt.Printf("  %%q  quoted string   → %s%q%s  %s(adds quotes and escapes special chars)%s\n", magenta, name, reset, dim, reset)
	fmt.Printf("  %%p  pointer address → %s%p%s\n", magenta, &name, reset)
	fmt.Printf("  %%b  binary          → %s%b%s  %s(age=5 in binary)%s\n", magenta, age, reset, dim, reset)
	fmt.Printf("  %%x  hexadecimal     → %s%x%s\n\n", magenta, age, reset)

	// Sprintf: returns formatted string (doesn't print)
	fmt.Printf("%s▸ Sprintf — format to string (no output)%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Sprintf returns the formatted string — use it to build values, not for output%s\n", green, reset)
	fmt.Printf("  %s⚠ In hot paths, prefer strconv.Itoa / strconv.FormatFloat — Sprintf allocates more%s\n", yellow, reset)

	msg := fmt.Sprintf("Hello, %s! You are %d years old.", name, age)
	fmt.Printf("  Sprintf result → %s%s%s\n\n", magenta, msg, reset)

	// Stringer interface — any type with String() method works with %s and %v
	fmt.Printf("%s▸ Stringer interface — custom formatting%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Implement String() string on your type and %%v/%%s uses it automatically%s\n", green, reset)
	fmt.Printf("  %s✔ This is Go's equivalent of Java's toString() or Python's __str__()%s\n", green, reset)

	dog := Animal{Name: "Dog", Sound: "woof", Legs: 4}
	fmt.Printf("  %%v with Stringer → %s%v%s\n", magenta, dog, reset)
	fmt.Printf("  %%#v ignores Stringer → %s%#v%s  %s(shows raw struct)%s\n", magenta, dog, reset, dim, reset)
}
