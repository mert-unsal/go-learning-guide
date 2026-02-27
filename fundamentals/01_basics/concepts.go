// Package basics covers Go fundamentals: variables, constants, types, and output.
// Run this file: go run fundamentals/01_basics/concepts.go
package basics

import "fmt"

// ============================================================
// 1. PACKAGE & IMPORTS
// ============================================================
// Every Go file starts with a package declaration.
// "main" package is special: it defines an executable program.
// Other packages (like this one) are libraries.
//
// Imports bring in other packages. Unused imports are a COMPILE ERROR in Go.

// ============================================================
// 2. VARIABLES
// ============================================================

func DemonstrateVariables() {
	// --- var keyword (explicit) ---
	var a int     // zero value: 0
	var b string  // zero value: ""
	var c bool    // zero value: false
	var d float64 // zero value: 0.0

	fmt.Println("Zero values:", a, b, c, d)

	// --- var with initializer ---
	var x int = 42
	var name string = "Gopher"
	fmt.Println(x, name)

	// --- Short variable declaration := (most common inside functions) ---
	// Type is INFERRED automatically
	age := 25
	pi := 3.14159
	isGoFun := true
	fmt.Println(age, pi, isGoFun)

	// --- Multiple assignment ---
	x1, y1 := 10, 20
	fmt.Println(x1, y1)

	// --- Swap values (Go's elegant way) ---
	x1, y1 = y1, x1
	fmt.Println("Swapped:", x1, y1)

	// --- Blank identifier _ (discard a value) ---
	_, second := 100, 200
	fmt.Println("Only second:", second)
}

// ============================================================
// 3. CONSTANTS
// ============================================================

// Constants are declared with 'const' and CANNOT be changed.
// They must be known at compile time.
const MaxSize = 100
const AppName = "LearnGo"
const Pi = 3.14159265358979

// Weekday --- iota: auto-incrementing constant generator ---
type Weekday int

const (
	Sunday    Weekday = iota // 0
	Monday                   // 1
	Tuesday                  // 2
	Wednesday                // 3
	Thursday                 // 4
	Friday                   // 5
	Saturday                 // 6
)

// iota with bit shifting — very common for flags/permissions
type Permission uint

const (
	Read    Permission = 1 << iota // 1  (001)
	Write                          // 2  (010)
	Execute                        // 4  (100)
)

func DemonstrateConstants() {
	fmt.Println("MaxSize:", MaxSize)
	fmt.Println("Monday:", Monday) // prints 1
	fmt.Println("Read:", Read, "Write:", Write, "Execute:", Execute)

	// Checking permissions with bitwise AND
	myPerms := Read | Write                        // 3 (011)
	fmt.Println("Can read?", myPerms&Read != 0)    // true
	fmt.Println("Can exec?", myPerms&Execute != 0) // false
}

// ============================================================
// 4. BASIC TYPES
// ============================================================
// Go is STATICALLY TYPED. Every variable has a fixed type.
//
// Integer types:
//   int8, int16, int32, int64  (signed)
//   uint8, uint16, uint32, uint64  (unsigned)
//   int, uint  (platform-dependent: 32 or 64 bit)
//   byte = uint8 (alias)
//   rune = int32 (alias, represents a Unicode code point)
//
// Float types:
//   float32, float64
//
// Complex types:
//   complex64, complex128
//
// String: immutable sequence of bytes (UTF-8 encoded)
// Bool: true or false

func DemonstrateTypes() {
	// Integer arithmetic
	var i int = 10
	var j int = 3
	fmt.Println(i+j, i-j, i*j, i/j, i%j) // 13 7 30 3 1

	// Type conversion (EXPLICIT — no implicit conversion in Go!)
	var f float64 = float64(i) / float64(j)
	fmt.Println(f) // 3.3333...

	// String operations
	s := "Hello, 世界"    // Go strings support Unicode
	fmt.Println(len(s)) // len() returns BYTES, not characters!

	// To iterate over characters (runes) use range:
	for i, ch := range s {
		if i < 5 {
			fmt.Printf("index=%d char=%c unicode=%d\n", i, ch, ch)
		}
	}

	// byte vs rune
	var myByte byte = 'A' // single quotes = rune/byte literal
	var myRune rune = '世'
	fmt.Printf("byte: %d, rune: %d\n", myByte, myRune)

	// String to/from byte slice
	str := "hello"
	bytes := []byte(str) // convert string to []byte
	bytes[0] = 'H'
	fmt.Println(string(bytes)) // "Hello"
}

// ============================================================
// 5. ZERO VALUES
// ============================================================
// In Go, variables are ALWAYS initialized to their zero value.
// This prevents undefined behavior common in C/C++.
//
//   int, float     → 0
//   bool           → false
//   string         → ""
//   pointer        → nil
//   slice          → nil
//   map            → nil
//   channel        → nil
//   function       → nil
//   interface      → nil
//   struct         → all fields set to their zero values

func DemonstrateZeroValues() {
	var i int
	var f float64
	var b bool
	var s string
	fmt.Printf("int: %d, float: %f, bool: %t, string: %q\n", i, f, b, s)
}

// ============================================================
// 6. fmt PACKAGE — FORMATTING VERBS
// ============================================================
func DemonstrateFmt() {
	name := "Gopher"
	age := 5
	pi := 3.14159

	// Printf: formatted printing
	fmt.Printf("Name: %s\n", name)       // string
	fmt.Printf("Age: %d\n", age)         // integer
	fmt.Printf("Pi: %.2f\n", pi)         // float with 2 decimal places
	fmt.Printf("Pi: %e\n", pi)           // scientific notation
	fmt.Printf("Bool: %t\n", true)       // boolean
	fmt.Printf("Type: %T\n", name)       // prints the TYPE of variable
	fmt.Printf("Value: %v\n", age)       // default format (works for anything)
	fmt.Printf("Go syntax: %#v\n", name) // Go-syntax representation

	// Sprintf: returns formatted string (doesn't print)
	msg := fmt.Sprintf("Hello, %s! You are %d years old.", name, age)
	fmt.Println(msg)

	// Stringer interface — any type with String() method works with %s and %v
}

// RunAll runs all demonstrations (useful for a main package)
func RunAll() {
	fmt.Println("\n=== Variables ===")
	DemonstrateVariables()
	fmt.Println("\n=== Constants ===")
	DemonstrateConstants()
	fmt.Println("\n=== Types ===")
	DemonstrateTypes()
	fmt.Println("\n=== Zero Values ===")
	DemonstrateZeroValues()
	fmt.Println("\n=== fmt ===")
	DemonstrateFmt()
}
