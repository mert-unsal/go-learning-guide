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
	var f Formatter
	f = JSONFormatter{}
	asJSON(f)
	describe(f, "event")

	f = TextFormatter{}
	asJSON(f) // not a JSONFormatter
	describe(f, "event")
}
