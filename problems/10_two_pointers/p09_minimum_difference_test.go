package two_pointers

import "testing"

func TestMinimumDifference(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"example 1", []int{90}, 1, 0},
		{"example 2", []int{9, 4, 1, 7}, 2, 2},
		{"sorted", []int{1, 2, 3, 4, 5}, 3, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinimumDifference(tt.nums, tt.k)
			if got != tt.want {
				t.Errorf("MinimumDifference(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}
