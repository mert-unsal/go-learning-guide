package arrays

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name  string
		nums1 []int
		m     int
		nums2 []int
		n     int
		want  []int
	}{
		// Basic cases
		{"basic", []int{1, 2, 3, 0, 0, 0}, 3, []int{2, 5, 6}, 3, []int{1, 2, 2, 3, 5, 6}},
		{"one empty", []int{1}, 1, []int{}, 0, []int{1}},
		{"first empty", []int{0}, 0, []int{1}, 1, []int{1}},
		{"all same", []int{2, 2, 0, 0}, 2, []int{2, 2}, 2, []int{2, 2, 2, 2}},

		// nums2 elements all smaller — forces full copy of nums2
		{"nums2 all smaller", []int{4, 5, 6, 0, 0, 0}, 3, []int{1, 2, 3}, 3, []int{1, 2, 3, 4, 5, 6}},

		// nums2 elements all larger — nums1 stays in place
		{"nums2 all larger", []int{1, 2, 3, 0, 0, 0}, 3, []int{7, 8, 9}, 3, []int{1, 2, 3, 7, 8, 9}},

		// Interleaved elements
		{"interleaved", []int{1, 3, 5, 7, 0, 0, 0, 0}, 4, []int{2, 4, 6, 8}, 4, []int{1, 2, 3, 4, 5, 6, 7, 8}},

		// Negative numbers
		{"negatives", []int{-5, -3, -1, 0, 0, 0}, 3, []int{-4, -2, 0}, 3, []int{-5, -4, -3, -2, -1, 0}},
		{"mixed negative positive", []int{-10, 0, 10, 0, 0, 0}, 3, []int{-5, 5, 15}, 3, []int{-10, -5, 0, 5, 10, 15}},

		// Single element arrays
		{"single into single", []int{2, 0}, 1, []int{1}, 1, []int{1, 2}},
		{"single larger into single", []int{1, 0}, 1, []int{2}, 1, []int{1, 2}},

		// Large gap between values
		{"large gap", []int{1, 1000000, 0, 0}, 2, []int{500, 999999}, 2, []int{1, 500, 999999, 1000000}},

		// Duplicates across both arrays
		{"cross duplicates", []int{1, 2, 3, 0, 0, 0}, 3, []int{1, 2, 3}, 3, []int{1, 1, 2, 2, 3, 3}},

		// nums1 has one element, nums2 has many
		{"one vs many", []int{5, 0, 0, 0, 0}, 1, []int{1, 2, 3, 4}, 4, []int{1, 2, 3, 4, 5}},

		// nums1 has many, nums2 has one
		{"many vs one", []int{1, 2, 3, 4, 0}, 4, []int{3}, 1, []int{1, 2, 3, 3, 4}},

		// All zeros
		{"all zeros", []int{0, 0, 0, 0}, 2, []int{0, 0}, 2, []int{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Merge(tt.nums1, tt.m, tt.nums2, tt.n)
			if !reflect.DeepEqual(tt.nums1, tt.want) {
				t.Errorf("Merge() = %v, want %v", tt.nums1, tt.want)
			}
		})
	}
}
