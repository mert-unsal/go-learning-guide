package dynamic_prog

import "testing"

func TestCanJump(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want bool
	}{
		{"reachable", []int{2, 3, 1, 1, 4}, true},
		{"stuck", []int{3, 2, 1, 0, 4}, false},
		{"single", []int{0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanJump(tt.nums)
			if got != tt.want {
				t.Errorf("CanJump(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
