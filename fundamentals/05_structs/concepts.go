// Package structs covers Go structs: definition, methods, embedding,
// struct tags, and composition patterns.
package structs

import "fmt"

// ============================================================
// 1. STRUCT DEFINITION
// ============================================================
// A struct groups related data. It's Go's primary way to create
// custom data types (no classes in Go!).

type Point struct {
	X float64
	Y float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

// Method on Rectangle (value receiver — no modification needed)
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// ============================================================
// 2. STRUCT INITIALIZATION
// ============================================================

func DemonstrateInitialization() {
	// Named fields (PREFERRED — order independent, self-documenting)
	p1 := Point{X: 1.0, Y: 2.0}

	// Positional (avoid — fragile if fields are added/reordered)
	p2 := Point{3.0, 4.0}

	// Zero value struct
	var p3 Point // X=0, Y=0

	// Pointer to struct
	p4 := &Point{X: 5.0, Y: 6.0}

	fmt.Println(p1, p2, p3, *p4)

	// Accessing fields
	fmt.Println("X:", p1.X, "Y:", p1.Y)

	// Through a pointer — Go auto-dereferences
	fmt.Println("p4.X:", p4.X) // same as (*p4).X
}

// ============================================================
// 3. STRUCT EMBEDDING (Composition over Inheritance)
// ============================================================
// Go doesn't have inheritance. Instead, use EMBEDDING to compose types.
// An embedded struct's fields and methods are PROMOTED to the outer struct.

type Animal struct {
	Name string
}

func (a Animal) Speak() string {
	return a.Name + " makes a sound"
}

type Dog struct {
	Animal // embedded (anonymous field) — NOT Animal Animal
	Breed  string
}

func (d Dog) Speak() string { // Dog can OVERRIDE Animal's method
	return d.Name + " says: Woof!"
}

type Cat struct {
	Animal
	Indoor bool
}

func DemonstrateEmbedding() {
	d := Dog{
		Animal: Animal{Name: "Rex"},
		Breed:  "Labrador",
	}

	// Access embedded fields directly (promoted)
	fmt.Println(d.Name)    // "Rex" — from Animal
	fmt.Println(d.Breed)   // "Labrador"
	fmt.Println(d.Speak()) // "Rex says: Woof!" — Dog's override

	// Still accessible via explicit path
	fmt.Println(d.Animal.Speak()) // "Rex makes a sound"

	c := Cat{Animal: Animal{Name: "Whiskers"}, Indoor: true}
	fmt.Println(c.Speak()) // "Whiskers makes a sound" — from Animal (no override)
	fmt.Println(c.Indoor)
}

// ============================================================
// 4. STRUCT TAGS
// ============================================================
// Tags add metadata to struct fields. Used by JSON, database, validation libs.
// Access via reflect package.

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"` // omit if empty
	Password string `json:"-"`               // never include in JSON
	Age      int    `json:"age,omitempty"`
}

// ============================================================
// 5. CONSTRUCTOR FUNCTIONS (Go pattern)
// ============================================================
// Go has no constructors. Convention: use a New* function.

type Stack struct {
	items []int
}

func NewStack() *Stack {
	return &Stack{items: make([]int, 0)}
}

func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	n := len(s.items)
	item := s.items[n-1]
	s.items = s.items[:n-1]
	return item, true
}

func (s *Stack) Peek() (int, bool) {
	if len(s.items) == 0 {
		return 0, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack) Len() int {
	return len(s.items)
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Stringer interface: implement String() for custom printing
func (s *Stack) String() string {
	return fmt.Sprintf("Stack%v", s.items)
}

func DemonstrateStack() {
	s := NewStack()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	fmt.Println(s) // Stack[1 2 3]

	if top, ok := s.Pop(); ok {
		fmt.Println("Popped:", top) // 3
	}
	fmt.Println("Length:", s.Len()) // 2
}

// ============================================================
// 6. ANONYMOUS STRUCTS
// ============================================================
// Useful for one-off data grouping, test cases, etc.

func DemonstrateAnonymousStructs() {
	// Define and initialize in one step
	point := struct {
		X, Y int
	}{X: 10, Y: 20}
	fmt.Println(point)

	// Very common in tests:
	tests := []struct {
		input    int
		expected int
	}{
		{1, 1},
		{2, 4},
		{3, 9},
	}
	for _, tt := range tests {
		result := tt.input * tt.input
		if result == tt.expected {
			fmt.Printf("PASS: %d^2 = %d\n", tt.input, result)
		}
	}
}

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Initialization ===")
	DemonstrateInitialization()
	fmt.Println("\n=== Embedding ===")
	DemonstrateEmbedding()
	fmt.Println("\n=== Stack ===")
	DemonstrateStack()
	fmt.Println("\n=== Anonymous Structs ===")
	DemonstrateAnonymousStructs()
}
