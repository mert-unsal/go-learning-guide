// Package interfaces covers Go interfaces: implicit implementation,
// the empty interface, type assertions, type switches, and common patterns.
package interfaces

import (
	"fmt"
	"math"
)

// ============================================================
// 1. DEFINING AND IMPLEMENTING INTERFACES
// ============================================================
// An interface defines a SET OF METHODS.
// A type implements an interface by implementing ALL its methods.
// There is NO 'implements' keyword — it's IMPLICIT (duck typing).

type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle implements Shape — it has both Area() and Perimeter()
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Rectangle also implements Shape
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// printShape accepts ANY type that implements Shape
func printShape(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func DemonstrateInterfaces() {
	c := Circle{Radius: 5}
	r := Rectangle{Width: 4, Height: 3}

	printShape(c)
	printShape(r)

	// Slice of interface — polymorphism!
	shapes := []Shape{c, r, Circle{Radius: 1}}
	totalArea := 0.0
	for _, s := range shapes {
		totalArea += s.Area()
	}
	fmt.Printf("Total area: %.2f\n", totalArea)
}

// ============================================================
// 2. THE EMPTY INTERFACE (interface{} or any)
// ============================================================
// interface{} is satisfied by ALL types (it has no methods).
// 'any' is an alias for interface{} (Go 1.18+).
// Use sparingly — prefer concrete types when possible.

func printAnything(v any) {
	fmt.Printf("Type: %T, Value: %v\n", v, v)
}

func DemonstrateEmptyInterface() {
	printAnything(42)
	printAnything("hello")
	printAnything(true)
	printAnything([]int{1, 2, 3})
	printAnything(nil)

	// Storing anything in a map
	m := map[string]any{
		"name":   "Alice",
		"age":    30,
		"active": true,
	}
	fmt.Println(m)
}

// ============================================================
// 3. TYPE ASSERTIONS
// ============================================================
// Extract the underlying concrete type from an interface value.
// Two forms:
//   t := i.(T)        // panics if i is not type T
//   t, ok := i.(T)    // safe form: ok is false if i is not T

func DemonstrateTypeAssertions() {
	var i any = "hello"

	// Safe type assertion
	s, ok := i.(string)
	if ok {
		fmt.Println("String:", s, "Length:", len(s))
	}

	// Wrong type — ok is false, no panic
	n, ok := i.(int)
	fmt.Println("Int:", n, "ok:", ok) // 0, false

	// Unsafe — would panic if wrong type
	// s2 := i.(int) // PANIC: interface conversion: string is not int
}

// ============================================================
// 4. TYPE SWITCH
// ============================================================
// A switch that checks the TYPE of an interface value.
// This is the idiomatic way to handle multiple types.

func describe(i any) string {
	switch v := i.(type) { // v holds the value as the specific type
	case int:
		return fmt.Sprintf("int: %d (doubled: %d)", v, v*2)
	case string:
		return fmt.Sprintf("string: %q (length: %d)", v, len(v))
	case bool:
		return fmt.Sprintf("bool: %t", v)
	case []int:
		return fmt.Sprintf("[]int with %d elements", len(v))
	case nil:
		return "nil value"
	default:
		return fmt.Sprintf("unknown type: %T", v)
	}
}

func DemonstrateTypeSwitch() {
	fmt.Println(describe(42))
	fmt.Println(describe("Go"))
	fmt.Println(describe(true))
	fmt.Println(describe([]int{1, 2, 3}))
	fmt.Println(describe(nil))
	fmt.Println(describe(3.14))
}

// ============================================================
// 5. INTERFACE COMPOSITION
// ============================================================
// Interfaces can embed other interfaces.

type Reader interface {
	Read() string
}

type Writer interface {
	Write(s string)
}

type ReadWriter interface {
	Reader // embeds Reader
	Writer // embeds Writer
}

type Buffer struct {
	data string
}

func (b *Buffer) Read() string   { return b.data }
func (b *Buffer) Write(s string) { b.data += s }

func DemonstrateComposition() {
	var rw ReadWriter = &Buffer{}
	rw.Write("Hello")
	rw.Write(", World!")
	fmt.Println(rw.Read()) // Hello, World!
}

// ============================================================
// 6. STRINGER INTERFACE (fmt.Stringer)
// ============================================================
// If your type implements String() string, fmt will use it automatically.

type Temperature struct {
	Celsius float64
}

func (t Temperature) String() string {
	return fmt.Sprintf("%.1f°C (%.1f°F)", t.Celsius, t.Celsius*9/5+32)
}

func DemonstrateStringer() {
	t := Temperature{Celsius: 37}
	fmt.Println(t)        // 37.0°C (98.6°F)
	fmt.Printf("%v\n", t) // same
	fmt.Printf("%s\n", t) // same
}

// ============================================================
// 7. ERROR INTERFACE
// ============================================================
// The built-in error interface: type error interface { Error() string }
// Any type with Error() string implements error.

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field=%s, msg=%s", e.Field, e.Message)
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Message: "must be non-negative"}
	}
	if age > 150 {
		return &ValidationError{Field: "age", Message: "unrealistic value"}
	}
	return nil
}

func DemonstrateErrorInterface() {
	err := validateAge(-5)
	if err != nil {
		fmt.Println("Error:", err)
		// Type assert to get the ValidationError
		if ve, ok := err.(*ValidationError); ok {
			fmt.Println("Field:", ve.Field)
		}
	}
}

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Interfaces ===")
	DemonstrateInterfaces()
	fmt.Println("\n=== Empty Interface ===")
	DemonstrateEmptyInterface()
	fmt.Println("\n=== Type Assertions ===")
	DemonstrateTypeAssertions()
	fmt.Println("\n=== Type Switch ===")
	DemonstrateTypeSwitch()
	fmt.Println("\n=== Interface Composition ===")
	DemonstrateComposition()
	fmt.Println("\n=== Stringer ===")
	DemonstrateStringer()
	fmt.Println("\n=== Error Interface ===")
	DemonstrateErrorInterface()
}
