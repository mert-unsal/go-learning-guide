package hard

import "testing"

func TestJump(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"example 1", []int{2, 3, 1, 1, 4}, 2},
		{"example 2", []int{2, 3, 0, 1, 4}, 2},
		{"single", []int{0}, 0},
		{"two", []int{1, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Jump(tt.nums)
			if got != tt.want {
				t.Errorf("Jump(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
