package bit_manipulation

import "testing"

func TestHammingWeight(t *testing.T) {
	tests := []struct {
		name string
		n    uint32
		want int
	}{
		{"11", 11, 3},
		{"128", 128, 1},
		{"max", 4294967293, 31},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HammingWeight(tt.n); got != tt.want {
				t.Errorf("HammingWeight(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}
