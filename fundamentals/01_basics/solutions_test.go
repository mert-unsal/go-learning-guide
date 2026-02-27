package basics

import (
	"testing"
)

// ============================================================
// TESTS — 01 Basics
// ============================================================
// These tests validate YOUR implementations in exercises.go.
// Run: go test ./fundamentals/01_basics/... -v

// Exercise 1
func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		celsius    float64
		fahrenheit float64
	}{
		{0, 32},
		{100, 212},
		{-40, -40}, // interesting: -40C == -40F
		{37, 98.6},
	}
	for _, tt := range tests {
		got := CelsiusToFahrenheit(tt.celsius)
		if got != tt.fahrenheit {
			t.Errorf("❌ CelsiusToFahrenheit(%v) = %v, want %v", tt.celsius, got, tt.fahrenheit)
		} else {
			t.Logf("✅ CelsiusToFahrenheit(%v) = %v", tt.celsius, got)
		}
	}
}

// Exercise 2
func TestSwapInts(t *testing.T) {
	a, b := SwapInts(3, 7)
	if a != 7 || b != 3 {
		t.Errorf("❌ SwapInts(3,7) = (%v,%v), want (7,3)", a, b)
	} else {
		t.Logf("✅ SwapInts(3,7) = (%v,%v)", a, b)
	}
}

// Exercise 3
func TestCharacterCount(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello", 5},
		{"世界", 2}, // 2 Chinese characters but 6 bytes — must use rune!
		{"Hello世界", 7},
		{"", 0},
	}
	for _, tt := range tests {
		got := CharacterCount(tt.input)
		if got != tt.want {
			t.Errorf("❌ CharacterCount(%q) = %v, want %v  ← Hint: use []rune(s) not len(s)", tt.input, got, tt.want)
		} else {
			t.Logf("✅ CharacterCount(%q) = %v", tt.input, got)
		}
	}
}

// Exercise 4
func TestMinMax(t *testing.T) {
	cases := []struct {
		a, b, c          int
		wantMin, wantMax int
	}{
		{3, 1, 2, 1, 3},
		{5, 5, 5, 5, 5},
		{-1, -5, 0, -5, 0},
	}
	for _, tt := range cases {
		gotMin, gotMax := MinMax(tt.a, tt.b, tt.c)
		if gotMin != tt.wantMin || gotMax != tt.wantMax {
			t.Errorf("❌ MinMax(%v,%v,%v) = (%v,%v), want (%v,%v)", tt.a, tt.b, tt.c, gotMin, gotMax, tt.wantMin, tt.wantMax)
		} else {
			t.Logf("✅ MinMax(%v,%v,%v) = (%v,%v)", tt.a, tt.b, tt.c, gotMin, gotMax)
		}
	}
}

// Exercise 5
func TestDirectionName(t *testing.T) {
	cases := []struct {
		dir  Direction
		want string
	}{
		{North, "North"},
		{East, "East"},
		{South, "South"},
		{West, "West"},
	}
	for _, tt := range cases {
		got := DirectionName(tt.dir)
		if got != tt.want {
			t.Errorf("❌ DirectionName(%v) = %q, want %q  ← Hint: use a switch statement", tt.dir, got, tt.want)
		} else {
			t.Logf("✅ DirectionName(%v) = %q", tt.dir, got)
		}
	}
}
