package reflect_ex

// ============================================================
// EXERCISES -- 11 reflect: Runtime Type Introspection
// ============================================================
// 12 exercises covering Go's reflect package at production depth.
// Focus: TypeOf, ValueOf, struct tags, when reflect is justified.

import (
	"reflect"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: TypeName -- return the type name of any value
// ────────────────────────────────────────────────────────────
// Given any value, return its reflect.Type.Name().
// For pointers, return the Elem() type name prefixed with "*".
// Examples: "int", "string", "*MyStruct"

func TypeName(v any) string {
	_ = reflect.TypeOf(v)
	return ""
}

// ────────────────────────────────────────────────────────────
// Exercise 2: IsNilSafe -- nil-safe interface check
// ────────────────────────────────────────────────────────────
// Return true if v is nil OR if v is an interface holding a nil pointer.
// This solves the classic Go nil-interface trap.
// A non-nil interface with a nil concrete value: var err error = (*MyError)(nil)

func IsNilSafe(v any) bool {
	return false
}

// ────────────────────────────────────────────────────────────
// Exercise 3: FieldNames -- extract struct field names
// ────────────────────────────────────────────────────────────
// Given a struct (or pointer to struct), return a slice of field names.
// If not a struct, return nil.

func FieldNames(v any) []string {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 4: GetTag -- read a struct tag from a field
// ────────────────────────────────────────────────────────────
// Given a struct, a field name, and a tag key, return the tag value.
// Example: GetTag(User{}, "Name", "json") → "name"
// Return "" if field or tag not found.

func GetTag(v any, fieldName, tagKey string) string {
	return ""
}

// ────────────────────────────────────────────────────────────
// Exercise 5: SetField -- set a struct field by name
// ────────────────────────────────────────────────────────────
// Given a pointer to a struct, set the named field to the given value.
// Return an error if the field doesn't exist or isn't settable.
// v must be a pointer to a struct.

func SetField(v any, fieldName string, value any) error {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 6: StructToMap -- convert struct to map[string]any
// ────────────────────────────────────────────────────────────
// Convert a struct to a map using the json tag as keys.
// If no json tag, use the field name. Skip fields with json:"-".
// Unexported fields are skipped.

func StructToMap(v any) map[string]any {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 7: MapToStruct -- populate struct from map
// ────────────────────────────────────────────────────────────
// Given a pointer to struct and a map[string]any, set fields matching
// the map keys to the map values. Match by field name (case-sensitive).
// Skip keys that don't match any field.

func MapToStruct(v any, data map[string]any) error {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 8: CallFunc -- invoke a function by reflection
// ────────────────────────────────────────────────────────────
// Given a function and a slice of arguments, call it and return results.
// Panic if fn is not a function.

func CallFunc(fn any, args ...any) []any {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 9: DeepEqual -- implement a simplified equality check
// ────────────────────────────────────────────────────────────
// Compare two values for deep equality. Support: basic types, slices, maps.
// For slices: compare element by element. For maps: compare key-value pairs.
// You may use reflect.DeepEqual internally, but add nil-safe handling:
// nil slice == empty slice should return true (unlike reflect.DeepEqual).

func DeepEqualNilSafe(a, b any) bool {
	return false
}

// ────────────────────────────────────────────────────────────
// Exercise 10: MakeSlice -- create a typed slice dynamically
// ────────────────────────────────────────────────────────────
// Given a sample element value and a length, create a new slice of that type.
// Return the slice as any.
// Example: MakeSlice(0, 3) → []int{0, 0, 0}

func MakeSlice(sampleElem any, length int) any {
	_ = reflect.MakeSlice
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 11: Implements -- check if a type implements an interface
// ────────────────────────────────────────────────────────────
// Given a value, check if its type implements the error interface.
// Return true if the value's type (or pointer-to-type) has an Error() string method.

func ImplementsError(v any) bool {
	return false
}

// ────────────────────────────────────────────────────────────
// Exercise 12: ValidateRequired -- struct validation via tags
// ────────────────────────────────────────────────────────────
// Check all fields with `validate:"required"` tag.
// Return the name of the first field that is its zero value, or "" if all set.
// Only check exported fields.

func ValidateRequired(v any) string {
	return ""
}
