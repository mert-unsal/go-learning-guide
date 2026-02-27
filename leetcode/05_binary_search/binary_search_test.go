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

func TestFindMin(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"rotated 5", []int{3, 4, 5, 1, 2}, 1},
		{"rotated 7", []int{4, 5, 6, 7, 0, 1, 2}, 0},
		{"not rotated", []int{11, 13, 15, 17}, 11},
		{"single", []int{1}, 1},
		{"two elements rotated", []int{2, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMin(tt.nums)
			if got != tt.want {
				t.Errorf("FindMin(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{"found", []int{-1, 0, 3, 5, 9, 12}, 9, 4},
		{"not found", []int{-1, 0, 3, 5, 9, 12}, 2, -1},
		{"empty", []int{}, 5, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BinarySearch(tt.nums, tt.target)
			if got != tt.want {
				t.Errorf("BinarySearch(%v, %d) = %d, want %d", tt.nums, tt.target, got, tt.want)
			}
		})
	}
}

func TestSearchMatrix(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]int
		target int
		want   bool
	}{
		{"found", [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 3, true},
		{"not found", [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 13, false},
		{"single", [][]int{{1}}, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SearchMatrix(tt.matrix, tt.target)
			if got != tt.want {
				t.Errorf("SearchMatrix target=%d = %v, want %v", tt.target, got, tt.want)
			}
		})
	}
}
