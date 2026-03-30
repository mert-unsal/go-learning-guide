package two_pointers

import "testing"

func TestFourSum(t *testing.T) {
	tests := []struct {
		name      string
		nums      []int
		target    int
		wantCount int
	}{
		{"example 1", []int{1, 0, -1, 0, -2, 2}, 0, 3},
		{"example 2", []int{2, 2, 2, 2, 2}, 8, 1},
		{"no result", []int{1, 2, 3, 4}, 100, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FourSum(tt.nums, tt.target)
			if len(got) != tt.wantCount {
				t.Errorf("FourSum(%v, %d) returned %d quadruplets, want %d", tt.nums, tt.target, len(got), tt.wantCount)
			}
		})
	}
}
