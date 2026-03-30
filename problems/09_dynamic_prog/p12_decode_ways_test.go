package dynamic_prog

import "testing"

func TestNumDecodings(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"basic", "12", 2},
		{"triple", "226", 3},
		{"leading zero", "06", 0},
		{"single", "1", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NumDecodings(tt.s)
			if got != tt.want {
				t.Errorf("NumDecodings(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}
