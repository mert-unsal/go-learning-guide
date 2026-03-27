package bit_manipulation

import (
	"fmt"
	"testing"
)

func TestSingleNumber(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{2, 2, 1}, 1},
		{[]int{4, 1, 2, 1, 2}, 4},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("SingleNumber(%v)", tt.nums), func(t *testing.T) {
			if got := SingleNumber(tt.nums); got != tt.want {
				t.Errorf("SingleNumber(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
