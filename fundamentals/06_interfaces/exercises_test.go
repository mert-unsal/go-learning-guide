package interfaces

import "testing"

func TestExStringers(t *testing.T) {
	b := ExBook{Title: "The Go Programming Language", Author: "Donovan"}
	want := `"The Go Programming Language" by Donovan`
	got := b.String()
	if got != want {
		t.Errorf("❌ Book.String() = %q, want %q", got, want)
	} else {
		t.Logf("✅ Book.String() = %q", got)
	}

	m := ExMovie{Title: "Inception", Year: 2010}
	gotM := m.String()
	if gotM != "Inception (2010)" {
		t.Errorf("❌ Movie.String() = %q, want \"Inception (2010)\"", gotM)
	} else {
		t.Logf("✅ Movie.String() = %q", gotM)
	}
}

func TestExBufferWriter(t *testing.T) {
	bw := &ExBufferWriter{}
	bw.Write("hello")
	bw.Write("world")
	if len(bw.Buffer) != 2 {
		t.Errorf("❌ buffer has %d items, want 2  ← Hint: append data to Buffer", len(bw.Buffer))
	} else {
		t.Logf("✅ ExBufferWriter.Buffer = %v", bw.Buffer)
	}
	if len(bw.Buffer) >= 1 && bw.Buffer[0] != "hello" {
		t.Errorf("❌ buffer[0] = %q, want \"hello\"", bw.Buffer[0])
	}
}

func TestDescribe(t *testing.T) {
	tests := []struct {
		input interface{}
		want  string
	}{
		{42, "int: 42"},
		{"hello", "string: hello"},
		{true, "bool: true"},
		{3.14, "unknown"},
	}
	for _, tt := range tests {
		got := Describe(tt.input)
		if got != tt.want {
			t.Errorf("❌ Describe(%v) = %q, want %q  ← Hint: use a type switch", tt.input, got, tt.want)
		} else {
			t.Logf("✅ Describe(%v) = %q", tt.input, got)
		}
	}
}

// ─── Exercise 4a: SafeGetLabel (Guard 1 — fix at source) ───

func TestSafeGetLabel_NonEmpty(t *testing.T) {
	l := SafeGetLabel("Widget")
	if l == nil {
		t.Error("❌ SafeGetLabel(\"Widget\") returned nil, want a Labeler")
		return
	}
	got := l.Label()
	want := "Product: Widget"
	if got != want {
		t.Errorf("❌ SafeGetLabel(\"Widget\").Label() = %q, want %q", got, want)
	} else {
		t.Logf("✅ SafeGetLabel(\"Widget\").Label() = %q", got)
	}
}

func TestSafeGetLabel_Empty(t *testing.T) {
	l := SafeGetLabel("")
	if l != nil {
		t.Errorf("❌ SafeGetLabel(\"\") != nil — you returned a typed nil (*Product)(nil) instead of a true nil interface\n"+
			"   Hint: return nil directly, don't return a *Product variable that happens to be nil\n"+
			"   Under the hood: iface{tab: *itab(Product), data: nil} — tab is NOT nil!")
	} else {
		t.Log("✅ SafeGetLabel(\"\") == nil (true nil interface)")
	}
}

// ─── Exercise 4b: SafeCallLabeler (Guard 2 — type assertion) ───

func TestSafeCallLabeler_Valid(t *testing.T) {
	l := &Product{Name: "Gadget"}
	got := SafeCallLabeler(l)
	want := "Product: Gadget"
	if got != want {
		t.Errorf("❌ SafeCallLabeler(valid) = %q, want %q", got, want)
	} else {
		t.Logf("✅ SafeCallLabeler(valid) = %q", got)
	}
}

func TestSafeCallLabeler_TrueNil(t *testing.T) {
	got := SafeCallLabeler(nil)
	if got != "no labeler" {
		t.Errorf("❌ SafeCallLabeler(nil) = %q, want \"no labeler\"\n"+
			"   Hint: check l == nil first for the true nil case", got)
	} else {
		t.Log("✅ SafeCallLabeler(nil) = \"no labeler\"")
	}
}

func TestSafeCallLabeler_TypedNil(t *testing.T) {
	// This is the trap: a nil *Product wrapped in a Labeler interface.
	// The interface is NOT nil (tab is populated), but calling Label() would panic.
	var p *Product            // typed nil
	var l Labeler = p         // iface{tab: *itab, data: nil} — NOT nil!
	got := SafeCallLabeler(l)
	if got != "no labeler" {
		t.Errorf("❌ SafeCallLabeler(typed-nil) = %q, want \"no labeler\"\n"+
			"   This is the nil interface trap!\n"+
			"   l == nil is FALSE because iface.tab is populated with *Product type info.\n"+
			"   You must type-assert to *Product, then check if the pointer is nil.\n"+
			"   Hint: p, ok := l.(*Product); if !ok || p == nil { ... }", got)
	} else {
		t.Log("✅ SafeCallLabeler(typed-nil) = \"no labeler\" — you caught the trap!")
	}
}

// ─── Exercise 4c: IsTrulyNil (Guard 3 — reflect) ───

func TestIsTrulyNil(t *testing.T) {
	var nilPtr *Product
	var nilSlice []int
	var nilMap map[string]int
	var nilChan chan int
	var nilFunc func()

	tests := []struct {
		name  string
		input interface{}
		want  bool
		hint  string
	}{
		{"true nil interface", nil, true,
			"i == nil should catch this — both tab and data are nil"},
		{"nil *Product in interface", nilPtr, true,
			"eface has _type set but data is nil. reflect.ValueOf(i).IsNil() returns true"},
		{"nil slice in interface", nilSlice, true,
			"a nil slice wrapped in interface{} — IsNil() should catch it"},
		{"nil map in interface", nilMap, true,
			"a nil map wrapped in interface{} — IsNil() should catch it"},
		{"nil chan in interface", nilChan, true,
			"a nil chan wrapped in interface{} — IsNil() should catch it"},
		{"nil func in interface", nilFunc, true,
			"a nil func wrapped in interface{} — IsNil() should catch it"},
		{"valid *Product", &Product{Name: "X"}, false,
			"non-nil pointer — data word points to a real value"},
		{"int value", 42, false,
			"int is a non-nillable kind — calling IsNil() on it would PANIC"},
		{"string value", "hello", false,
			"string is non-nillable — guard against calling IsNil() on it"},
		{"struct value", Product{Name: "Y"}, false,
			"structs are non-nillable — reflect.Kind is reflect.Struct"},
		{"bool value", true, false,
			"bool is non-nillable"},
		{"non-nil slice", []int{1, 2}, false,
			"non-nil slice — data word has a valid pointer"},
		{"non-nil map", map[string]int{"a": 1}, false,
			"non-nil map — data word has a valid pointer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTrulyNil(tt.input)
			if got != tt.want {
				t.Errorf("❌ IsTrulyNil(%s) = %v, want %v\n   Hint: %s",
					tt.name, got, tt.want, tt.hint)
			} else {
				t.Logf("✅ IsTrulyNil(%s) = %v", tt.name, got)
			}
		})
	}
}
