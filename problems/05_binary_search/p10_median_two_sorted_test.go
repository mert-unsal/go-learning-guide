package binary_search

import "testing"

func TestFindMedianSortedArrays(t *testing.T) {
	tests := []struct {
		name  string
		nums1 []int
		nums2 []int
		want  float64
	}{
		{"odd total", []int{1, 3}, []int{2}, 2.0},
		{"even total", []int{1, 2}, []int{3, 4}, 2.5},
		{"one empty", []int{}, []int{1}, 1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMedianSortedArrays(tt.nums1, tt.nums2)
			if got != tt.want {
				t.Errorf("FindMedianSortedArrays(%v, %v) = %v, want %v", tt.nums1, tt.nums2, got, tt.want)
			}
		})
	}
}
