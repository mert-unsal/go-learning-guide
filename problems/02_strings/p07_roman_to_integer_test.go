package strings_problems

import "testing"

func TestRomanToInt(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"simple", "III", 3},
		{"additive", "LVIII", 58},
		{"subtractive", "MCMXCIV", 1994},
		{"subtractive small", "IV", 4},
		{"nine", "IX", 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RomanToInt(tt.s)
			if got != tt.want {
				t.Errorf("RomanToInt(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
