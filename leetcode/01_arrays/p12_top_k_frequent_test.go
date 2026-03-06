package arrays

import "testing"

func TestTopKFrequent(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int // expected length
	}{
		{"basic", []int{1, 1, 1, 2, 2, 3}, 2, 2},
		{"single", []int{1}, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TopKFrequent(tt.nums, tt.k)
			if len(got) != tt.want {
				t.Errorf("TopKFrequent(%v, %d) returned %d elements, want %d", tt.nums, tt.k, len(got), tt.want)
			}
		})
	}
}
