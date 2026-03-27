package binary_search

import "testing"

func TestSearch(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{"found in rotated", []int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
		{"not found", []int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
		{"single not found", []int{1}, 0, -1},
		{"single found", []int{1}, 1, 0},
		{"not rotated", []int{1, 2, 3, 4, 5}, 3, 2},
		{"target at end", []int{3, 1, 2}, 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Search(tt.nums, tt.target)
			if got != tt.want {
				t.Errorf("Search(%v, %d) = %d, want %d", tt.nums, tt.target, got, tt.want)
			}
		})
	}
}
