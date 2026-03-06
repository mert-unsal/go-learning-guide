package arrays

import "testing"

func TestMajorityElement(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{3, 2, 3}, 3},
		{"longer", []int{2, 2, 1, 1, 1, 2, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MajorityElement(tt.nums)
			if got != tt.want {
				t.Errorf("MajorityElement(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
