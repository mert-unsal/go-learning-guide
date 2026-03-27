package heap_priority_queue

import "testing"

func TestFindKthLargest(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"basic", []int{3, 2, 1, 5, 6, 4}, 2, 5},
		{"with dups", []int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := make([]int, len(tt.nums))
			copy(nums, tt.nums)
			got := FindKthLargest(nums, tt.k)
			if got != tt.want {
				t.Errorf("FindKthLargest(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}
