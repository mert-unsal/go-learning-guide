package bit_manipulation

import (
	"fmt"
	"testing"
)

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
		t.Run(fmt.Sprintf("IsPowerOfTwo(%d)", tt.n), func(t *testing.T) {
			if got := IsPowerOfTwo(tt.n); got != tt.want {
				t.Errorf("IsPowerOfTwo(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}
