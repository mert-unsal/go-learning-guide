package trees

import "testing"

func TestKthSmallest(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		k    int
		want int
	}{
		{"basic", []int{3, 1, 4, 0, 2}, 1, 1},
		{"second", []int{5, 3, 6, 2, 4, 0, 0, 1}, 3, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KthSmallest(newTree(tt.vals), tt.k)
			if got != tt.want {
				t.Errorf("KthSmallest = %v, want %v", got, tt.want)
			}
		})
	}
}
