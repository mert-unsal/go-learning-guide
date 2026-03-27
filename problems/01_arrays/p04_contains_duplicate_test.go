package arrays

import "testing"

func TestContainsDuplicate(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want bool
	}{
		{"has duplicate", []int{1, 2, 3, 1}, true},
		{"all unique", []int{1, 2, 3, 4}, false},
		{"single", []int{1}, false},
		{"empty", []int{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsDuplicate(tt.nums)
			if got != tt.want {
				t.Errorf("ContainsDuplicate(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
