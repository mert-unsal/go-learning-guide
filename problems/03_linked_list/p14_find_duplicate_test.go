package linked_list

import "testing"

func TestFindDuplicate(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{1, 3, 4, 2, 2}, 2},
		{"another", []int{3, 1, 3, 4, 2}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindDuplicate(tt.nums)
			if got != tt.want {
				t.Errorf("FindDuplicate(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
