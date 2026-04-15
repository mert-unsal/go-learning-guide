package reflect_ex

import (
	"fmt"
	"testing"
)

// Test types used across exercises
type Person struct {
	Name    string `json:"name" validate:"required"`
	Age     int    `json:"age"`
	Email   string `json:"email" validate:"required"`
	private string //nolint
}

type MyError struct{ Msg string }

func (e *MyError) Error() string { return e.Msg }

// ────────────────────────────────────────────────────────────
// Exercise 1: TypeName
// ────────────────────────────────────────────────────────────

func TestTypeName(t *testing.T) {
	tests := []struct {
		input any
		want  string
	}{
		{42, "int"},
		{"hello", "string"},
		{3.14, "float64"},
		{Person{}, "Person"},
	}

	for _, tt := range tests {
		got := TypeName(tt.input)
		if got != tt.want {
			t.Errorf("❌ TypeName(%v) = %q, want %q\n\t\t"+
				"Hint: reflect.TypeOf(v).Name() for named types. "+
				"For pointers, check Kind() == reflect.Ptr, then Elem().Name()",
				tt.input, got, tt.want)
		} else {
			t.Logf("✅ TypeName(%T) = %q", tt.input, got)
		}
	}

	// Pointer case
	p := &Person{}
	got := TypeName(p)
	if got != "*Person" {
		t.Errorf("❌ TypeName(*Person) = %q, want \"*Person\"", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: IsNilSafe
// ────────────────────────────────────────────────────────────

func TestIsNilSafe(t *testing.T) {
	// Plain nil
	if !IsNilSafe(nil) {
		t.Error("❌ IsNilSafe(nil) = false\n\t\t" +
			"Hint: if v == nil { return true }. " +
			"Then use reflect.ValueOf(v) and check Kind for chan/func/interface/map/pointer/slice")
	} else {
		t.Log("✅ IsNilSafe(nil) = true")
	}

	// Non-nil interface with nil pointer (the classic trap)
	var err error = (*MyError)(nil) // err != nil but the pointer inside IS nil
	if !IsNilSafe(err) {
		t.Error("❌ IsNilSafe((*MyError)(nil) as error) = false, want true\n\t\t" +
			"Hint: rv := reflect.ValueOf(v); if rv.Kind() is a nillable kind, check rv.IsNil()")
	} else {
		t.Log("✅ IsNilSafe(nil-pointer-in-interface) = true")
	}

	// Actually non-nil
	if IsNilSafe(&Person{Name: "Go"}) {
		t.Error("❌ IsNilSafe(&Person{}) = true, want false")
	}

	// Non-pointer value
	if IsNilSafe(42) {
		t.Error("❌ IsNilSafe(42) = true, want false")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: FieldNames
// ────────────────────────────────────────────────────────────

func TestFieldNames(t *testing.T) {
	names := FieldNames(Person{})
	if names == nil || len(names) != 4 {
		t.Errorf("❌ FieldNames(Person) = %v (len %d), want 4 fields\n\t\t"+
			"Hint: t := reflect.TypeOf(v); if t.Kind() == reflect.Ptr { t = t.Elem() }; "+
			"loop t.NumField(), collect t.Field(i).Name",
			names, len(names))
		return
	}
	if names[0] != "Name" || names[1] != "Age" {
		t.Errorf("❌ FieldNames = %v, want [Name Age Email private]", names)
	} else {
		t.Logf("✅ FieldNames(Person) = %v", names)
	}

	// Pointer to struct
	ptrNames := FieldNames(&Person{})
	if len(ptrNames) != 4 {
		t.Error("❌ FieldNames(&Person{}) should work for pointer to struct too")
	}

	// Non-struct
	if FieldNames(42) != nil {
		t.Error("❌ FieldNames(42) should return nil")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: GetTag
// ────────────────────────────────────────────────────────────

func TestGetTag(t *testing.T) {
	got := GetTag(Person{}, "Name", "json")
	if got != "name" {
		t.Errorf("❌ GetTag(Person, \"Name\", \"json\") = %q, want \"name\"\n\t\t"+
			"Hint: field, ok := reflect.TypeOf(v).FieldByName(fieldName); "+
			"field.Tag.Get(tagKey)",
			got)
	} else {
		t.Logf("✅ GetTag(Person, Name, json) = %q", got)
	}

	// validate tag
	got = GetTag(Person{}, "Name", "validate")
	if got != "required" {
		t.Errorf("❌ GetTag(Person, Name, validate) = %q, want \"required\"", got)
	}

	// Missing field
	got = GetTag(Person{}, "Missing", "json")
	if got != "" {
		t.Errorf("❌ GetTag(Person, Missing, json) = %q, want \"\"", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: SetField
// ────────────────────────────────────────────────────────────

func TestSetField(t *testing.T) {
	p := &Person{Name: "old", Age: 20}
	err := SetField(p, "Name", "new")
	if err != nil {
		t.Errorf("❌ SetField error = %v\n\t\t"+
			"Hint: rv := reflect.ValueOf(v).Elem(); "+
			"field := rv.FieldByName(name); field.Set(reflect.ValueOf(value))",
			err)
		return
	}
	if p.Name != "new" {
		t.Errorf("❌ after SetField, Name = %q, want \"new\"", p.Name)
	} else {
		t.Logf("✅ SetField set Name to %q", p.Name)
	}

	// Set int field
	err = SetField(p, "Age", 99)
	if err != nil || p.Age != 99 {
		t.Errorf("❌ SetField Age: err=%v, Age=%d", err, p.Age)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: StructToMap
// ────────────────────────────────────────────────────────────

func TestStructToMap(t *testing.T) {
	p := Person{Name: "Go", Age: 15, Email: "go@example.com"}
	m := StructToMap(p)
	if m == nil {
		t.Fatal("❌ StructToMap returned nil\n\t\t" +
			"Hint: loop fields, skip unexported (field.PkgPath != \"\"), " +
			"use json tag as key if present, skip json:\"-\"")
	}
	if m["name"] != "Go" {
		t.Errorf("❌ m[\"name\"] = %v, want \"Go\" (uses json tag as key)", m["name"])
	}
	if m["age"] != 15 {
		t.Errorf("❌ m[\"age\"] = %v, want 15", m["age"])
	}
	if _, exists := m["private"]; exists {
		t.Error("❌ unexported field 'private' should be skipped")
	}
	t.Logf("✅ StructToMap = %v", m)
}

// ────────────────────────────────────────────────────────────
// Exercise 7: MapToStruct
// ────────────────────────────────────────────────────────────

func TestMapToStruct(t *testing.T) {
	p := &Person{}
	data := map[string]any{
		"Name":  "Go",
		"Age":   15,
		"Extra": "ignored",
	}
	err := MapToStruct(p, data)
	if err != nil {
		t.Errorf("❌ MapToStruct error = %v\n\t\t"+
			"Hint: rv := reflect.ValueOf(v).Elem(); for key, val := range data { "+
			"field := rv.FieldByName(key); if field.IsValid() && field.CanSet() { "+
			"field.Set(reflect.ValueOf(val)) } }",
			err)
		return
	}
	if p.Name != "Go" || p.Age != 15 {
		t.Errorf("❌ after MapToStruct: Name=%q Age=%d, want Go/15", p.Name, p.Age)
	} else {
		t.Logf("✅ MapToStruct set Name=%q Age=%d", p.Name, p.Age)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: CallFunc
// ────────────────────────────────────────────────────────────

func TestCallFunc(t *testing.T) {
	add := func(a, b int) int { return a + b }
	results := CallFunc(add, 3, 4)
	if results == nil || len(results) != 1 {
		t.Fatal("❌ CallFunc(add, 3, 4) returned nil or wrong count\n\t\t" +
			"Hint: rv := reflect.ValueOf(fn); args := make([]reflect.Value, len(args)); " +
			"out := rv.Call(args); convert back to []any")
	}
	if results[0] != 7 {
		t.Errorf("❌ CallFunc(add, 3, 4) = %v, want 7", results[0])
	} else {
		t.Logf("✅ CallFunc(add, 3, 4) = %v", results[0])
	}

	// Function with string args
	greet := func(name string) string { return fmt.Sprintf("Hello, %s!", name) }
	results = CallFunc(greet, "Go")
	if results == nil || results[0] != "Hello, Go!" {
		t.Errorf("❌ CallFunc(greet, Go) = %v, want \"Hello, Go!\"", results)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 9: DeepEqualNilSafe
// ────────────────────────────────────────────────────────────

func TestDeepEqualNilSafe(t *testing.T) {
	// Basic equality
	if !DeepEqualNilSafe(42, 42) {
		t.Error("❌ DeepEqualNilSafe(42, 42) = false\n\t\t" +
			"Hint: handle nil cases first, then use reflect.DeepEqual. " +
			"Special case: nil slice == empty slice should be true")
	}

	// Nil slice == empty slice
	var nilSlice []int
	emptySlice := []int{}
	if !DeepEqualNilSafe(nilSlice, emptySlice) {
		t.Error("❌ DeepEqualNilSafe(nil, []) = false, want true\n\t\t" +
			"Hint: check if both are slices with len 0")
	} else {
		t.Log("✅ nil slice == empty slice")
	}

	// Nil map == empty map
	var nilMap map[string]int
	emptyMap := map[string]int{}
	if !DeepEqualNilSafe(nilMap, emptyMap) {
		t.Error("❌ DeepEqualNilSafe(nil map, empty map) = false, want true")
	}

	// Unequal
	if DeepEqualNilSafe([]int{1, 2}, []int{1, 3}) {
		t.Error("❌ DeepEqualNilSafe([1,2], [1,3]) = true, want false")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 10: MakeSlice
// ────────────────────────────────────────────────────────────

func TestMakeSlice(t *testing.T) {
	result := MakeSlice(0, 3)
	if result == nil {
		t.Fatal("❌ MakeSlice(0, 3) returned nil\n\t\t" +
			"Hint: elemType := reflect.TypeOf(sampleElem); " +
			"sliceType := reflect.SliceOf(elemType); " +
			"slice := reflect.MakeSlice(sliceType, length, length); " +
			"return slice.Interface()")
	}

	intSlice, ok := result.([]int)
	if !ok {
		t.Fatalf("❌ MakeSlice(0, 3) type = %T, want []int", result)
	}
	if len(intSlice) != 3 {
		t.Errorf("❌ len = %d, want 3", len(intSlice))
	} else {
		t.Logf("✅ MakeSlice(0, 3) = %v ([]int)", intSlice)
	}

	// String slice
	strResult := MakeSlice("", 2)
	strSlice, ok := strResult.([]string)
	if !ok || len(strSlice) != 2 {
		t.Errorf("❌ MakeSlice(\"\", 2) = %v, want []string of len 2", strResult)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 11: ImplementsError
// ────────────────────────────────────────────────────────────

func TestImplementsError(t *testing.T) {
	if !ImplementsError(&MyError{Msg: "test"}) {
		t.Error("❌ ImplementsError(*MyError) = false, want true\n\t\t" +
			"Hint: errorType := reflect.TypeOf((*error)(nil)).Elem(); " +
			"reflect.TypeOf(v).Implements(errorType)")
	} else {
		t.Log("✅ *MyError implements error")
	}

	if ImplementsError("not an error") {
		t.Error("❌ ImplementsError(string) = true, want false")
	}

	if ImplementsError(42) {
		t.Error("❌ ImplementsError(42) = true, want false")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: ValidateRequired
// ────────────────────────────────────────────────────────────

func TestValidateRequired(t *testing.T) {
	// All required fields set
	p := Person{Name: "Go", Email: "go@go.dev"}
	if field := ValidateRequired(p); field != "" {
		t.Errorf("❌ ValidateRequired(valid) = %q, want \"\"\n\t\t"+
			"Hint: loop fields, check tag.Get(\"validate\") == \"required\", "+
			"then check field value == zero value using reflect.Value.IsZero()",
			field)
	} else {
		t.Log("✅ ValidateRequired(valid person) = \"\" (all set)")
	}

	// Missing Name
	p2 := Person{Email: "go@go.dev"}
	if field := ValidateRequired(p2); field != "Name" {
		t.Errorf("❌ ValidateRequired(missing Name) = %q, want \"Name\"", field)
	} else {
		t.Log("✅ ValidateRequired catches missing Name")
	}

	// Missing Email
	p3 := Person{Name: "Go"}
	if field := ValidateRequired(p3); field != "Email" {
		t.Errorf("❌ ValidateRequired(missing Email) = %q, want \"Email\"", field)
	} else {
		t.Log("✅ ValidateRequired catches missing Email")
	}
}
