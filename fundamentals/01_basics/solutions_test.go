package basics

import (
	"testing"
)

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
		got := CelsiusToFahrenheitSolution(tt.celsius)
		if got != tt.fahrenheit {
			t.Errorf("CelsiusToFahrenheit(%v) = %v, want %v", tt.celsius, got, tt.fahrenheit)
		}
	}
}

func TestSwapInts(t *testing.T) {
	a, b := SwapIntsSolution(3, 7)
	if a != 7 || b != 3 {
		t.Errorf("SwapInts(3,7) = (%v,%v), want (7,3)", a, b)
	}
}

func TestCharacterCount(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello", 5},
		{"世界", 2}, // 2 Chinese characters but 6 bytes
		{"Hello世界", 7},
		{"", 0},
	}
	for _, tt := range tests {
		got := CharacterCountSolution(tt.input)
		if got != tt.want {
			t.Errorf("CharacterCount(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestMinMax(t *testing.T) {
	min, max := MinMaxSolution(3, 1, 2)
	if min != 1 || max != 3 {
		t.Errorf("MinMax(3,1,2) = (%v,%v), want (1,3)", min, max)
	}
}

func TestDirectionName(t *testing.T) {
	if DirectionNameSolution(North) != "North" {
		t.Error("Expected North")
	}
	if DirectionNameSolution(West) != "West" {
		t.Error("Expected West")
	}
}
