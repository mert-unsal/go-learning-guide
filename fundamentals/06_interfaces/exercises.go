package interfaces

import (
	"reflect"
)

// ExStringer ============================================================
// EXERCISES — 06 Interfaces
// ============================================================
// Exercise 1:
// Stringer — any type with String() string method.
// Book and Movie both implement it.
type ExStringer interface {
	String() string
}
type ExBook struct{ Title, Author string }

type ExMovie struct {
	Title string
	Year  int
}

func (b ExBook) String() string {
	return ""
}

func (m ExMovie) String() string {
	return ""
}

func PrintAll(items []ExStringer) {
}

// ExWriter Exercise 2:
// ExWriter interface — anything that can Write a string.
type ExWriter interface {
	Write(data string) error
}
type ExBufferWriter struct {
	Buffer []string
}

func (bw *ExBufferWriter) Write(data string) error {
	return nil
}

func WriteAll(w ExWriter, items []string) error {
	return nil
}

// Exercise 3:
// Type switch — describe what kind of value is passed.
func Describe(i interface{}) string {
	return ""
}

// ============================================================
// Exercise 4: The Nil Interface Trap & Guards
// ============================================================
// This exercise tests your understanding of iface internals.
// An interface is runtime.iface{tab *itab, data unsafe.Pointer}.
// It's nil ONLY when both tab and data are nil.
// A typed nil (e.g., *MyType(nil)) assigned to an interface fills
// the tab → the interface is NOT nil, but calling methods panics.
//
// You will implement three functions:
//
// 4a. SafeGetLabel — Guard 1 (fix at source)
//     Return a true nil interface when there's no Labeler, not a typed nil.
//
// 4b. SafeCallLabeler — Guard 2 (type assertion guard)
//     Given a Labeler interface, safely extract the concrete *Product
//     and return its label. Handle both true nil AND typed nil.
//
// 4c. IsTrulyNil — Guard 3 (reflect-based guard)
//     Given any interface{}, determine if it's truly nil — either a true
//     nil interface OR an interface holding a nil pointer/slice/map/chan/func.
//     Use reflect.ValueOf and v.IsNil(). Handle non-nillable kinds (int, struct, etc).

// Labeler is a simple interface for the nil guard exercises.
type Labeler interface {
	Label() string
}

// Product is a concrete type that implements Labeler.
type Product struct {
	Name string
}

func (p *Product) Label() string {
	return ""
}

// Exercise 4a: SafeGetLabel — Guard 1 (fix at source)
// Given a product name, return a Labeler.
// If the name is empty, return a TRUE nil interface (not a typed nil *Product).
// If the name is non-empty, return a &Product{Name: name}.
//
// Hint: the trap is doing `var p *Product; return p` — that returns
// iface{tab: *itab, data: nil} which is NOT nil.
func SafeGetLabel(name string) Labeler {
	return nil
}

// Exercise 4b: SafeCallLabeler — Guard 2 (type assertion guard)
// Given a Labeler interface value, safely return its label string.
// Return "no labeler" if:
//   - l is a true nil interface, OR
//   - l holds a typed nil (e.g., a nil *Product wrapped in the interface)
//
// Use a type assertion to extract *Product, then check if the pointer is nil.
func SafeCallLabeler(l Labeler) string {
	return ""
}

// Exercise 4c: IsTrulyNil — Guard 3 (reflect-based, for unknown types)
// Return true if i is:
//   - a true nil interface (both tab and data are nil), OR
//   - an interface holding a nil pointer, slice, map, chan, or func
//
// Return false for:
//   - non-nil concrete values (42, "hello", Product{}, etc.)
//   - non-nillable kinds (int, string, struct, bool, etc.) — these can NEVER be nil
//
// Use: reflect.ValueOf(i), v.Kind(), v.IsNil()
// Be careful: calling IsNil() on a non-nillable kind (e.g., reflect.Int) PANICS.
//
// Suppress the "reflect imported and not used" error by using reflect here.
var _ = reflect.ValueOf // keep import alive
func IsTrulyNil(i interface{}) bool {
	// TODO: implement
	return false
}

// ============================================================
// Exercise 5: Method Sets & Interface Satisfaction
// ============================================================
// This exercise tests your understanding of Go's method set rules:
//
//   Type T  → method set = only value receivers
//   Type *T → method set = value receivers + pointer receivers
//
// A type satisfies an interface only if its method set includes
// ALL the interface's methods.
//
// This means:
//   - If interface requires value-receiver method  → both T and *T satisfy it
//   - If interface requires pointer-receiver method → only *T satisfies it
//
// WHY? Go can always dereference *T to get T, but it CANNOT always take
// the address of T (map values, return values, constants have no stable address).
//
// You will implement:
//
// 5a. Two types (Celsius, Kelvin) with specific receiver types,
//     and predict which can be assigned to an interface.
//
// 5b. A function that accepts a Sensor interface and calls both methods.
//     You must understand which concrete types can be passed.
//
// 5c. A function that collects multiple Stringer types into a slice,
//     choosing the correct way to add value vs pointer receiver types.

// Sensor is an interface requiring two methods:
//   - Reading() returns the sensor value (can be value receiver)
//   - Calibrate() adjusts the sensor (needs to mutate state → pointer receiver)
type Sensor interface {
	Reading() float64
	Calibrate(offset float64)
}

// Thermometer Exercise 5a: Implement Thermometer
// A Thermometer has a Temp field (float64).
// - Reading() should return t.Temp              → use VALUE receiver
// - Calibrate() should add offset to t.Temp     → use POINTER receiver (mutation!)
//
// Think: which type satisfies Sensor — Thermometer or *Thermometer?
type Thermometer struct {
	Temp float64
}

// TODO: implement Reading() with value receiver
func (t Thermometer) Reading() float64 {
	// TODO: return t.Temp
	return 0
}

// TODO: implement Calibrate() with pointer receiver
func (t *Thermometer) Calibrate(offset float64) {
	// TODO: add offset to t.Temp
}

// ReadAndCalibrate Exercise 5b: ReadAndCalibrate
// Given a Sensor, return the reading, then calibrate by the given offset,
// then return the new reading.
// Return (before, after).
//
// Think: can you pass Thermometer{Temp: 20.0} to this function?
//
//	Or must you pass &Thermometer{Temp: 20.0}?
func ReadAndCalibrate(s Sensor, offset float64) (before, after float64) {
	// TODO: implement
	return
}

// Displayer is an interface with only a value-receiver method.
type Displayer interface {
	Display() string
}

// Celsius uses a VALUE receiver for Display.
// Both Celsius and *Celsius satisfy Displayer.
type Celsius float64

// Display
// TODO: implement Display() on Celsius with VALUE receiver
// Return format: "XX.X°C" (e.g., "36.6°C")
// Hint: use fmt.Sprintf("%.1f°C", float64(c))
func (c Celsius) Display() string {
	// TODO: implement
	return ""
}

// Kelvin uses a POINTER receiver for Display.
// ONLY *Kelvin satisfies Displayer.
type Kelvin float64

// TODO: implement Display() on Kelvin with POINTER receiver
// Return format: "XX.XK" (e.g., "309.8K")
// Hint: use fmt.Sprintf("%.1fK", float64(*k))
func (k *Kelvin) Display() string {
	// TODO: implement
	return ""
}

// Exercise 5c: CollectDisplayers
// Return a slice of Displayer containing these values in order:
//  1. Celsius(36.6)
//  2. Kelvin(309.8)
//
// The challenge: one uses value receiver, one uses pointer receiver.
// You must figure out which needs & and which doesn't.
//
// Hint: If Kelvin.Display() has a pointer receiver, can you assign
//
//	Kelvin(309.8) directly to Displayer? Or do you need &Kelvin?
//	Remember: you can't take the address of a converted literal directly.
func CollectDisplayers() []Displayer {
	// TODO: implement — return a []Displayer with Celsius and Kelvin values
	return nil
}
