// Package interfaces covers Go interfaces the Go way:
// implicit satisfaction, behavior-based grouping, consumer-defined contracts,
// small interfaces, and structural typing — NOT OOP inheritance.
package interfaces

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

// ============================================================
// THE MENTAL MODEL — READ THIS FIRST
// ============================================================
//
// In OOP languages (Java, C#) you think:
//   "These types ARE a Shape, so I'll declare them as implementing Shape."
//   → The PRODUCER declares the relationship upfront.
//   → Types are grouped by what they ARE (taxonomy/hierarchy).
//
// In Go you think:
//   "My function needs something it can call Area() on. I'll define that contract."
//   → The CONSUMER defines the interface, right where it needs it.
//   → Types are grouped by what they CAN DO (behavior), not what they are.
//
// This is the fundamental difference. Go has no class, no extends,
// no implements. A type satisfies an interface just by having the methods —
// it doesn't know about the interface and doesn't need to.
//
// KEY RULES:
//   1. Keep interfaces small (1-3 methods). The smaller, the more useful.
//   2. Define interfaces where they are CONSUMED, not where types are defined.
//   3. Don't create interfaces speculatively — create them when you need abstraction.
//   4. Accept interfaces, return concrete types (Rob Pike's guideline).

// ============================================================
// 1. IMPLICIT SATISFACTION — NO 'implements' KEYWORD
// ============================================================
// The type and the interface have ZERO connection at the type definition site.
// The compiler checks at the point of USE, not at the point of definition.

// Stringer is defined in the fmt package. We define our own here to illustrate.
// Any type with a String() string method satisfies it — it never needs to say so.
type Stringer interface {
	String() string
}

type User struct {
	Name string
	Age  int
}

// User satisfies Stringer. User doesn't know Stringer exists.
// Stringer doesn't know User exists. The compiler connects them at use time.
func (u User) String() string {
	return fmt.Sprintf("%s (age %d)", u.Name, u.Age)
}

type Point struct {
	X, Y float64
}

// Point also satisfies Stringer — completely unrelated type, same behavior.
func (p Point) String() string {
	return fmt.Sprintf("(%.1f, %.1f)", p.X, p.Y)
}

// printIt was written knowing nothing about User or Point.
// It works with both because they share the behavior, not a type hierarchy.
func printIt(s Stringer) {
	fmt.Println(s.String())
}

func DemonstrateImplicitSatisfaction() {
	// User and Point are completely different types.
	// They were never declared to implement anything.
	// They work here purely because they have the right method.
	printIt(User{Name: "Alice", Age: 30})
	printIt(Point{X: 3.5, Y: -1.2})
}

// ============================================================
// 2. THE CONSUMER DEFINES THE INTERFACE
// ============================================================
// This is the key Go pattern. You do NOT go to the producer's package
// and make their type implement your interface. You define the interface
// in YOUR package describing only what YOU need.
//
// Imagine you're writing a function that saves data somewhere.
// You don't care if it's a file, a database, or an in-memory buffer.
// You only care that it can Write bytes.

// Saver is defined HERE, by the consumer. It describes the minimum
// behavior this package needs. It does NOT live next to File or DB.
type Saver interface {
	Save(data string) error
}

// FileSaver and DBSaver are defined independently — they know nothing about Saver.
type FileSaver struct{ Path string }

func (f FileSaver) Save(data string) error {
	// (pretend we write to a file)
	fmt.Printf("[FileSaver] writing %d bytes to %s\n", len(data), f.Path)
	return nil
}

type DBSaver struct{ Table string }

func (d DBSaver) Save(data string) error {
	// (pretend we write to a DB)
	fmt.Printf("[DBSaver] inserting into table %s: %q\n", d.Table, data)
	return nil
}

// persist only knows about Saver — it is fully decoupled from FileSaver and DBSaver.
// You can add a new storage backend without touching this function at all.
func persist(s Saver, data string) {
	if err := s.Save(data); err != nil {
		fmt.Println("save failed:", err)
	}
}

func DemonstrateConsumerDefinedInterface() {
	persist(FileSaver{Path: "/tmp/data.txt"}, "hello world")
	persist(DBSaver{Table: "events"}, "hello world")
}

// ============================================================
// 3. SMALL INTERFACES — THE io.Reader / io.Writer LESSON
// ============================================================
// Go's standard library is built on tiny, single-method interfaces.
// The smaller the interface, the more types can satisfy it,
// and the more places it can be used.
//
//   io.Reader  → Read(p []byte) (n int, err error)
//   io.Writer  → Write(p []byte) (n int, err error)
//   io.Closer  → Close() error
//   fmt.Stringer → String() string
//   error      → Error() string
//
// Because io.Reader is ONE method, strings, files, network connections,
// HTTP bodies, gzip readers, test buffers — ALL satisfy it.
// A 10-method interface would exclude most of them.
//
// Rule of thumb: if your interface has more than 3 methods, ask yourself
// if it's really one concern or several concerns bundled together.

// Measurer is a small, focused interface.
type Measurer interface {
	Measure() float64
}

// Notice: Circle and Rectangle are defined here NOT because they ARE shapes
// (that's OOP thinking), but because we need concrete types for the example.
// The insight is that ANYTHING with a Measure() method works — a distance,
// a weight, a duration, a file size.
type Circle struct{ Radius float64 }
type Rectangle struct{ Width, Height float64 }
type Segment struct{ Length float64 } // not a "shape" at all
type FileSize struct{ Bytes int64 }   // completely different domain

func (c Circle) Measure() float64    { return math.Pi * c.Radius * c.Radius }
func (r Rectangle) Measure() float64 { return r.Width * r.Height }
func (s Segment) Measure() float64   { return s.Length }
func (f FileSize) Measure() float64  { return float64(f.Bytes) }

func totalMeasure(items []Measurer) float64 {
	total := 0.0
	for _, m := range items {
		total += m.Measure()
	}
	return total
}

func DemonstrateSmallInterface() {
	// Completely unrelated types used together — because behavior, not type.
	items := []Measurer{
		Circle{Radius: 2},
		Rectangle{Width: 3, Height: 4},
		Segment{Length: 10},
		FileSize{Bytes: 1024},
	}
	fmt.Printf("Total: %.2f\n", totalMeasure(items))
}

// ============================================================
// 4. INTERFACE COMPOSITION — BUILD BIGGER FROM SMALL
// ============================================================
// Instead of one big interface, compose small ones.
// This mirrors how Go's stdlib works (io.ReadWriter = Reader + Writer).

type Loader interface {
	Load(key string) (string, error)
}

type Storer interface {
	Store(key, value string) error
}

// Cache composes Loader and Storer.
// A type must satisfy BOTH to be used as a Cache.
type Cache interface {
	Loader
	Storer
}

type MemCache struct {
	data map[string]string
}

func NewMemCache() *MemCache {
	return &MemCache{data: make(map[string]string)}
}

func (m *MemCache) Load(key string) (string, error) {
	v, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("key %q not found", key)
	}
	return v, nil
}

func (m *MemCache) Store(key, value string) error {
	m.data[key] = value
	return nil
}

// lookupOrStore only needs to load — it accepts the narrower Loader.
// This makes it easier to test and reuse.
func lookupOrStore(l Loader, key string) {
	v, err := l.Load(key)
	if err != nil {
		fmt.Println("miss:", err)
		return
	}
	fmt.Println("hit:", v)
}

func DemonstrateComposition() {
	c := NewMemCache()
	_ = c.Store("lang", "Go")
	lookupOrStore(c, "lang")    // hit: Go
	lookupOrStore(c, "missing") // miss: key not found
}

// ============================================================
// 5. INTERFACE VALUES — WHAT THEY ACTUALLY ARE AT RUNTIME
// ============================================================
// An interface value is a TWO-WORD pair internally: (type, value).
//
//   var s Stringer        → (nil,   nil  )   s == nil is TRUE
//   s = User{...}         → (User,  0xc0...)  s == nil is FALSE
//   s = (*User)(nil)      → (*User, nil  )   s == nil is FALSE  ← THE TRAP
//
// The trap: once a type is stored in the interface, the interface is no longer
// nil even if the pointer inside it is nil. Checking s == nil is NOT enough
// to know whether the underlying value is safe to use.

func DemonstrateInterfaceValues() {
	// Case 1: true nil interface — both slots are nil
	var s Stringer
	fmt.Println("nil interface:", s == nil) // true

	// Case 2: non-nil interface with a real value
	s = User{Name: "Bob", Age: 25}
	fmt.Println("non-nil interface:", s == nil) // false
	fmt.Println(s.String())

	// Case 3: THE TRAP — nil pointer stored inside a non-nil interface
	var u *User                                                  // u is nil
	s = u                                                        // s is now (*User, nil) — the type slot is filled
	fmt.Println("nil pointer in interface, s == nil:", s == nil) // FALSE!
	// s == nil is false, but calling s.String() would PANIC
	// because the method tries to dereference a nil pointer.
}

// ────────────────────────────────────────────────────────────
// HOW TO GUARD AGAINST THE NIL POINTER TRAP
// ────────────────────────────────────────────────────────────
//
// Three approaches, from last resort to best practice:
//
// APPROACH 1 — reflect.ValueOf
//   Use only when you receive an interface from outside code you don't control.
//   It checks whether the value stored inside the interface is itself nil.
//
// APPROACH 2 — type assertion
//   Clean, zero reflection overhead. Use when you know the possible concrete types.
//
// APPROACH 3 — fix it at the source (always the best option)
//   Never let a typed nil enter an interface. Return a true nil interface instead.

// Approach 1: reflect catches the nil pointer inside the interface.
func isTrulyNil(i any) bool {
	if i == nil {
		return true // fast path: interface itself is (nil, nil)
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface,
		reflect.Slice, reflect.Map,
		reflect.Chan, reflect.Func:
		return v.IsNil()
	}
	return false
}

// Approach 2: assert to the concrete type, then check the pointer.
func safeCall(s Stringer) {
	u, ok := s.(*User) // extract the concrete *User
	if !ok || u == nil {
		fmt.Println("  nil or wrong type — skipping")
		return
	}
	fmt.Println(" ", u.String())
}

// Approach 3 (BAD): returns (*User, nil) — caller's == nil check will FAIL.
func findUserBad(found bool) Stringer {
	var u *User
	if !found {
		return u // typed nil leaks into the interface → BUG
	}
	return &User{Name: "Alice", Age: 30}
}

// Approach 3 (GOOD): returns a true (nil, nil) interface — == nil works correctly.
func findUserGood(found bool) Stringer {
	if !found {
		return nil // untyped nil → true nil interface
	}
	return &User{Name: "Alice", Age: 30}
}

func DemonstrateNilGuards() {
	var u *User
	var s Stringer = u // trap: s holds (*User, nil)

	fmt.Println("── reflect guard ──")
	fmt.Println("s == nil:      ", s == nil)       // false — the trap
	fmt.Println("isTrulyNil(s): ", isTrulyNil(s))  // true  — correct
	fmt.Println("isTrulyNil(42):", isTrulyNil(42)) // false

	fmt.Println("── type assertion guard ──")
	safeCall(s)                            // skipped safely
	safeCall(User{Name: "Carol", Age: 28}) // printed

	fmt.Println("── fix at the source ──")
	bad := findUserBad(false)
	good := findUserGood(false)
	fmt.Println("bad  == nil:", bad == nil)  // false — the bug
	fmt.Println("good == nil:", good == nil) // true  — correct
}

// ============================================================
// 6. TYPE ASSERTIONS AND TYPE SWITCHES
// ============================================================
// When you have an interface value and need the concrete type back,
// use a type assertion (single type) or type switch (multiple types).
// Prefer type switches — they communicate intent more clearly.

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

func DemonstrateTypeAssertions() {
	var f Formatter
	f = JSONFormatter{}
	asJSON(f)
	describe(f, "event")

	f = TextFormatter{}
	asJSON(f) // not a JSONFormatter
	describe(f, "event")
}

// ============================================================
// 7. REAL-WORLD PATTERN: error IS an interface
// ============================================================
// The built-in error is just: type error interface { Error() string }
// Any type with Error() string is an error — no registration needed.
// This is Go's interfaces working exactly as designed.

type NotFoundError struct{ Resource string }
type PermissionError struct {
	User   string
	Action string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

func (e *PermissionError) Error() string {
	return fmt.Sprintf("user %q cannot %s", e.User, e.Action)
}

// findRecord returns the error interface — callers don't need to import
// NotFoundError or PermissionError to handle errors generically.
func findRecord(id int, user string) error {
	if id <= 0 {
		return &NotFoundError{Resource: fmt.Sprintf("record#%d", id)}
	}
	if user == "guest" {
		return &PermissionError{User: user, Action: "read records"}
	}
	return nil
}

func DemonstrateErrorInterface() {
	for _, call := range []struct {
		id   int
		user string
	}{{-1, "alice"}, {1, "guest"}, {1, "alice"}} {
		err := findRecord(call.id, call.user)
		if err == nil {
			fmt.Printf("id=%d user=%s → ok\n", call.id, call.user)
			continue
		}
		// Use type switch to handle specific error kinds
		switch e := err.(type) {
		case *NotFoundError:
			fmt.Printf("not found: %s\n", e.Resource)
		case *PermissionError:
			fmt.Printf("denied: %s tried to %s\n", e.User, e.Action)
		default:
			fmt.Println("unexpected error:", err)
		}
	}
}

// ============================================================
// 8. THE EMPTY INTERFACE — any
// ============================================================
// interface{} (aliased as 'any' since Go 1.18) has zero methods.
// Every type satisfies it. Use it only when you truly cannot know
// the type at compile time (e.g., JSON decoding, fmt internals).
// It throws away compile-time type safety — use sparingly.

func printAnything(v any) {
	fmt.Printf("type: %-12T value: %v\n", v, v)
}

func DemonstrateEmptyInterface() {
	printAnything(42)
	printAnything("hello")
	printAnything(true)
	printAnything([]int{1, 2, 3})
	printAnything(nil)
}

// ============================================================
// SUMMARY: Go interfaces vs OOP interfaces
// ============================================================
//
//  OOP                            Go
//  ─────────────────────────────────────────────────────────
//  declared by the PRODUCER       defined by the CONSUMER
//  "class X implements I"         type just has the methods
//  big, upfront hierarchies       small, composed on demand
//  groups by WHAT THINGS ARE      groups by WHAT THINGS DO
//  coupling at definition time    coupling resolved at compile time, at use site
//  abstract classes, inheritance  no inheritance, only behavior contracts
//
// The practical result: in Go you can write an interface for a type
// defined in a package you don't own, and it will just work.
// No forking, no wrapping, no adapter boilerplate.

// RunAll runs all demonstrations
func RunAll() {
	sep := strings.Repeat("-", 40)
	fmt.Println("\n=== 1. Implicit Satisfaction ===")
	DemonstrateImplicitSatisfaction()
	fmt.Println(sep)
	fmt.Println("=== 2. Consumer-Defined Interface ===")
	DemonstrateConsumerDefinedInterface()
	fmt.Println(sep)
	fmt.Println("=== 3. Small Interface ===")
	DemonstrateSmallInterface()
	fmt.Println(sep)
	fmt.Println("=== 4. Interface Composition ===")
	DemonstrateComposition()
	fmt.Println(sep)
	fmt.Println("=== 5. Interface Values & Nil Trap ===")
	DemonstrateInterfaceValues()
	fmt.Println(sep)
	fmt.Println("=== 5b. Nil Guards ===")
	DemonstrateNilGuards()
	fmt.Println(sep)
	fmt.Println("=== 6. Type Assertions & Switches ===")
	DemonstrateTypeAssertions()
	fmt.Println(sep)
	fmt.Println("=== 7. Error Interface ===")
	DemonstrateErrorInterface()
	fmt.Println(sep)
	fmt.Println("=== 8. Empty Interface ===")
	DemonstrateEmptyInterface()
}
