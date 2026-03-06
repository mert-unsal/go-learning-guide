package bit_manipulation

import "testing"

func TestIsPowerOfTwo(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{1, true},
		{16, true},
		{3, false},
		{0, false},
		{-1, false},
	}
	for _, tt := range tests {
		if got := IsPowerOfTwo(tt.n); got != tt.want {
			t.Errorf("IsPowerOfTwo(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}
