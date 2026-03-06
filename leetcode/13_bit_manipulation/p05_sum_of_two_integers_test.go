package bit_manipulation

import "testing"

func TestGetSum(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{1, 2, 3},
		{-1, 1, 0},
		{0, 0, 0},
	}
	for _, tt := range tests {
		if got := GetSum(tt.a, tt.b); got != tt.want {
			t.Errorf("GetSum(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
