package trees

import "testing"

func TestMaxPathSum(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want int
	}{
		{"basic", []int{1, 2, 3}, 6},
		{"negative", []int{-10, 9, 20, 0, 0, 15, 7}, 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxPathSum(newTree(tt.vals))
			if got != tt.want {
				t.Errorf("MaxPathSum = %v, want %v", got, tt.want)
			}
		})
	}
}
