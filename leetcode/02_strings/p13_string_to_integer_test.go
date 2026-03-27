package strings_problems

import "testing"

func TestMyAtoi(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"positive", "42", 42},
		{"negative", "   -42", -42},
		{"with words", "4193 with words", 4193},
		{"empty", "", 0},
		{"overflow", "2147483648", 2147483647},
		{"underflow", "-2147483649", -2147483648},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MyAtoi(tt.s)
			if got != tt.want {
				t.Errorf("MyAtoi(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}
