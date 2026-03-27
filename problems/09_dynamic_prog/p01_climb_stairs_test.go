package dynamic_prog

import (
	"fmt"
	"testing"
)

func TestClimbStairs(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 5},
		{5, 8},
		{10, 89},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("ClimbStairs(%d)", tt.n), func(t *testing.T) {
			got := ClimbStairs(tt.n)
			if got != tt.want {
				t.Errorf("ClimbStairs(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}
