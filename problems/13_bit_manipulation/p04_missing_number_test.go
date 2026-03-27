package bit_manipulation

import "testing"

func TestMissingNumber(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{3, 0, 1}, 2},
		{"zero missing", []int{1}, 0},
		{"last missing", []int{0, 1}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MissingNumber(tt.nums); got != tt.want {
				t.Errorf("MissingNumber(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
