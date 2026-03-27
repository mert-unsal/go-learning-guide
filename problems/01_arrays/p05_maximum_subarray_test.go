package arrays

import "testing"

func TestMaxSubArray(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		// Original cases
		{"mixed", []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}, 6},
		{"all positive", []int{1, 2, 3}, 6},
		{"all negative", []int{-3, -1, -2}, -1},
		{"single", []int{5}, 5},

		// Edge cases
		{"single negative", []int{-5}, -5},
		{"single zero", []int{0}, 0},
		{"two elements positive", []int{1, 2}, 3},
		{"two elements one negative", []int{-1, 2}, 2},
		{"two elements both negative", []int{-1, -2}, -1},

		// Subarray at different positions
		{"max at start", []int{5, 4, -10, 1, 2}, 9},
		{"max at end", []int{-1, -2, 3, 4, 5}, 12},
		{"max in middle", []int{-5, 4, 6, -3, -10}, 10},

		// Zeros
		{"zeros and positives", []int{0, 0, 3, 0, 0}, 3},
		{"zeros and negatives", []int{0, -1, 0, -2, 0}, 0},
		{"all zeros", []int{0, 0, 0}, 0},

		// Large negative then recover
		{"dip and recover", []int{2, -1, 2, 3, -9, 1}, 6},
		{"deep dip separates subarrays", []int{5, -100, 6}, 6},
		{"shallow dip keeps subarray", []int{5, -1, 6}, 10},

		// Entire array is the answer
		{"whole array", []int{1, 2, 3, 4, 5}, 15},

		// Alternating
		{"alternating signs", []int{-1, 2, -1, 2, -1, 2}, 4},
		{"alternating large", []int{-10, 20, -10, 20, -10}, 30},

		//CUSTOM
		{"alternating signs", []int{10, 20, 30, -30, -25, -5, 100, 10}, 110},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxSubArray(tt.nums)
			if got != tt.want {
				t.Errorf("MaxSubArray(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
