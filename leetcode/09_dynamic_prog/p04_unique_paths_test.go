package dynamic_prog

import "testing"

func TestUniquePaths(t *testing.T) {
	tests := []struct {
		m, n int
		want int
	}{
		{3, 7, 28},
		{3, 2, 3},
		{1, 1, 1},
		{7, 3, 28},
	}
	for _, tt := range tests {
		got := UniquePaths(tt.m, tt.n)
		if got != tt.want {
			t.Errorf("UniquePaths(%d, %d) = %d, want %d", tt.m, tt.n, got, tt.want)
		}
	}
}
