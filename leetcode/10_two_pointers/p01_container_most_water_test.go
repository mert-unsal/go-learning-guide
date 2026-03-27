package two_pointers

import "testing"

func TestMaxArea(t *testing.T) {
	tests := []struct {
		name   string
		height []int
		want   int
	}{
		{"basic", []int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{"two walls", []int{1, 1}, 1},
		{"increasing", []int{1, 2, 3, 4, 5}, 6},
		{"decreasing", []int{5, 4, 3, 2, 1}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxArea(tt.height)
			if got != tt.want {
				t.Errorf("MaxArea(%v) = %d, want %d", tt.height, got, tt.want)
			}
		})
	}
}
