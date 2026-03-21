package interfaces

import (
	"fmt"
	"reflect"
)
// ============================================================
// EXERCISES — 06 Interfaces
// ============================================================
// Exercise 1:
// Stringer — any type with String() string method.
// Book and Movie both implement it.
type ExStringer interface {
String() string
}
type ExBook struct{ Title, Author string }
type ExMovie struct{ Title string; Year int }
func (b ExBook) String() string {
	return fmt.Sprintf("%q by %s", b.Title, b.Author)
}

func (m ExMovie) String() string {
	return fmt.Sprintf("%s (%d)", m.Title, m.Year)
}
func PrintAll(items []ExStringer) {
// TODO: print each item.String()
}
// Exercise 2:
// ExWriter interface — anything that can Write a string.
type ExWriter interface {
Write(data string) error
}
type ExBufferWriter struct {
Buffer []string
}
func (bw *ExBufferWriter) Write(data string) error {
	bw.Buffer = append(bw.Buffer, data)
	return nil
}

func WriteAll(w ExWriter, items []string) error {
	for _, item := range items {
		if err := w.Write(item); err != nil {
			return err
		}
	}
	return nil
}
// Exercise 3:
// Type switch — describe what kind of value is passed.
func Describe(i interface{}) string {
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int: %d", v)
	case string:
		return fmt.Sprintf("string: %s", v)
	case bool:
		return fmt.Sprintf("bool: %t", v)
	default:
		return "unknown"
	}
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
	return "Product: " + p.Name
}

// Exercise 4a: SafeGetLabel — Guard 1 (fix at source)
// Given a product name, return a Labeler.
// If the name is empty, return a TRUE nil interface (not a typed nil *Product).
// If the name is non-empty, return a &Product{Name: name}.
//
// Hint: the trap is doing `var p *Product; return p` — that returns
// iface{tab: *itab, data: nil} which is NOT nil.
func SafeGetLabel(name string) Labeler {
	// TODO: implement
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
	// TODO: implement
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