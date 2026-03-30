package binary_search

import "testing"

func TestMinEatingSpeed(t *testing.T) {
	tests := []struct {
		name  string
		piles []int
		h     int
		want  int
	}{
		{"basic", []int{3, 6, 7, 11}, 8, 4},
		{"tight", []int{30, 11, 23, 4, 20}, 5, 30},
		{"relaxed", []int{30, 11, 23, 4, 20}, 6, 23},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinEatingSpeed(tt.piles, tt.h)
			if got != tt.want {
				t.Errorf("MinEatingSpeed(%v, %d) = %v, want %v", tt.piles, tt.h, got, tt.want)
			}
		})
	}
}
