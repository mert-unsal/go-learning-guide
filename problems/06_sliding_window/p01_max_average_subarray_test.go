package sliding_window

import "testing"

func TestFindMaxAverage(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want float64
	}{
		{"basic", []int{1, 12, -5, -6, 50, 3}, 4, 12.75},
		{"single window", []int{5, 5, 5}, 3, 5.0},
		{"k=1", []int{3, 1, 4, 1, 5}, 1, 5.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMaxAverage(tt.nums, tt.k)
			if got != tt.want {
				t.Errorf("FindMaxAverage(%v, %d) = %v, want %v", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}
