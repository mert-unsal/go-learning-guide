package arrays

import (
	"fmt"
	"testing"
)

func TestMiniMaxSum(t *testing.T) {
	tests := []struct {
		arr     []int
		wantMin int
		wantMax int
	}{
		{[]int{1, 2, 3, 4, 5}, 10, 14},
		{[]int{7, 69, 2, 221, 8974}, 299, 9271},
		{[]int{1, 1, 1, 1, 1}, 4, 4},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("MiniMaxSum(%v)", tt.arr), func(t *testing.T) {
			gotMin, gotMax := MiniMaxSum(tt.arr)
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("MiniMaxSum(%v) = (%d,%d), want (%d,%d)", tt.arr, gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}
