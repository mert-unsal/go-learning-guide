package binary_search

import "testing"

func TestFindPeakElement(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int // acceptable peak indices
	}{
		{"single peak", []int{1, 2, 3, 1}, []int{2}},
		{"two peaks", []int{1, 2, 1, 3, 5, 6, 4}, []int{1, 5}},
		{"single element", []int{1}, []int{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindPeakElement(tt.nums)
			valid := false
			for _, w := range tt.want {
				if got == w {
					valid = true
					break
				}
			}
			if !valid {
				t.Errorf("FindPeakElement(%v) = %v, want one of %v", tt.nums, got, tt.want)
			}
		})
	}
}
