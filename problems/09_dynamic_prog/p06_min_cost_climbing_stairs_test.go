package dynamic_prog

import "testing"

func TestMinCostClimbingStairs(t *testing.T) {
	tests := []struct {
		name string
		cost []int
		want int
	}{
		{"basic", []int{10, 15, 20}, 15},
		{"longer", []int{1, 100, 1, 1, 1, 100, 1, 1, 100, 1}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinCostClimbingStairs(tt.cost)
			if got != tt.want {
				t.Errorf("MinCostClimbingStairs(%v) = %d, want %d", tt.cost, got, tt.want)
			}
		})
	}
}
