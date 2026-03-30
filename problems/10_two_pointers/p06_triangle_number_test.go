package two_pointers

import "testing"

func TestTriangleNumber(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"example 1", []int{2, 2, 3, 4}, 3},
		{"no triangles", []int{1, 1, 10}, 0},
		{"all same", []int{3, 3, 3, 3}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TriangleNumber(tt.nums)
			if got != tt.want {
				t.Errorf("TriangleNumber(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
