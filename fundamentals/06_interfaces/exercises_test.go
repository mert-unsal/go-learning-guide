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
