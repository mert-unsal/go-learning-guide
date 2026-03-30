package dynamic_prog

import "testing"

func TestCanPartition(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want bool
	}{
		{"can split", []int{1, 5, 11, 5}, true},
		{"cannot split", []int{1, 2, 3, 5}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanPartition(tt.nums)
			if got != tt.want {
				t.Errorf("CanPartition(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
