package arrays

import (
	"reflect"
	"testing"
)

// ============================================================
// TESTS — Arrays
// ============================================================

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{"basic", []int{2, 7, 11, 15}, 9, []int{0, 1}},
		{"middle pair", []int{3, 2, 4}, 6, []int{1, 2}},
		{"same element", []int{3, 3}, 6, []int{0, 1}},
		{"negative numbers", []int{-3, 4, 3, 90}, 0, []int{0, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TwoSum(tt.nums, tt.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TwoSum(%v, %d) = %v, want %v", tt.nums, tt.target, got, tt.want)
			}
		})
	}
}

func TestMaxProfit(t *testing.T) {
	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{"normal", []int{7, 1, 5, 3, 6, 4}, 5},
		{"no profit", []int{7, 6, 4, 3, 1}, 0},
		{"single day", []int{5}, 0},
		{"two days profit", []int{1, 5}, 4},
		{"empty", []int{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxProfit(tt.prices)
			if got != tt.want {
				t.Errorf("MaxProfit(%v) = %d, want %d", tt.prices, got, tt.want)
			}
		})
	}
}

func TestProductExceptSelf(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{"basic", []int{1, 2, 3, 4}, []int{24, 12, 8, 6}},
		{"with zero", []int{1, 0, 3, 4}, []int{0, 12, 0, 0}},
		{"two elements", []int{3, 4}, []int{4, 3}},
		{"negative", []int{-1, 1, 0, -3, 3}, []int{0, 0, 9, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProductExceptSelf(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductExceptSelf(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

func TestContainsDuplicate(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want bool
	}{
		{"has duplicate", []int{1, 2, 3, 1}, true},
		{"all unique", []int{1, 2, 3, 4}, false},
		{"single", []int{1}, false},
		{"empty", []int{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsDuplicate(tt.nums)
			if got != tt.want {
				t.Errorf("ContainsDuplicate(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

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

func TestMerge(t *testing.T) {
	tests := []struct {
		name  string
		nums1 []int
		m     int
		nums2 []int
		n     int
		want  []int
	}{
		// Basic cases
		{"basic", []int{1, 2, 3, 0, 0, 0}, 3, []int{2, 5, 6}, 3, []int{1, 2, 2, 3, 5, 6}},
		{"one empty", []int{1}, 1, []int{}, 0, []int{1}},
		{"first empty", []int{0}, 0, []int{1}, 1, []int{1}},
		{"all same", []int{2, 2, 0, 0}, 2, []int{2, 2}, 2, []int{2, 2, 2, 2}},

		// nums2 elements all smaller — forces full copy of nums2
		{"nums2 all smaller", []int{4, 5, 6, 0, 0, 0}, 3, []int{1, 2, 3}, 3, []int{1, 2, 3, 4, 5, 6}},

		// nums2 elements all larger — nums1 stays in place
		{"nums2 all larger", []int{1, 2, 3, 0, 0, 0}, 3, []int{7, 8, 9}, 3, []int{1, 2, 3, 7, 8, 9}},

		// Interleaved elements
		{"interleaved", []int{1, 3, 5, 7, 0, 0, 0, 0}, 4, []int{2, 4, 6, 8}, 4, []int{1, 2, 3, 4, 5, 6, 7, 8}},

		// Negative numbers
		{"negatives", []int{-5, -3, -1, 0, 0, 0}, 3, []int{-4, -2, 0}, 3, []int{-5, -4, -3, -2, -1, 0}},
		{"mixed negative positive", []int{-10, 0, 10, 0, 0, 0}, 3, []int{-5, 5, 15}, 3, []int{-10, -5, 0, 5, 10, 15}},

		// Single element arrays
		{"single into single", []int{2, 0}, 1, []int{1}, 1, []int{1, 2}},
		{"single larger into single", []int{1, 0}, 1, []int{2}, 1, []int{1, 2}},

		// Large gap between values
		{"large gap", []int{1, 1000000, 0, 0}, 2, []int{500, 999999}, 2, []int{1, 500, 999999, 1000000}},

		// Duplicates across both arrays
		{"cross duplicates", []int{1, 2, 3, 0, 0, 0}, 3, []int{1, 2, 3}, 3, []int{1, 1, 2, 2, 3, 3}},

		// nums1 has one element, nums2 has many
		{"one vs many", []int{5, 0, 0, 0, 0}, 1, []int{1, 2, 3, 4}, 4, []int{1, 2, 3, 4, 5}},

		// nums1 has many, nums2 has one
		{"many vs one", []int{1, 2, 3, 4, 0}, 4, []int{3}, 1, []int{1, 2, 3, 3, 4}},

		// All zeros
		{"all zeros", []int{0, 0, 0, 0}, 2, []int{0, 0}, 2, []int{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Merge(tt.nums1, tt.m, tt.nums2, tt.n)
			if !reflect.DeepEqual(tt.nums1, tt.want) {
				t.Errorf("Merge() = %v, want %v", tt.nums1, tt.want)
			}
		})
	}
}

func TestFindDisappearedNumbers(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{"basic", []int{4, 3, 2, 7, 8, 2, 3, 1}, []int{5, 6}},
		{"none missing", []int{1, 2}, []int(nil)},
		{"all missing except one", []int{2, 2}, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindDisappearedNumbers(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindDisappearedNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"k=3", []int{1, 2, 3, 4, 5, 6, 7}, 3, []int{5, 6, 7, 1, 2, 3, 4}},
		{"k=2", []int{-1, -100, 3, 99}, 2, []int{3, 99, -1, -100}},
		{"k=0", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"k=len", []int{1, 2, 3}, 3, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Rotate(tt.nums, tt.k)
			if !reflect.DeepEqual(tt.nums, tt.want) {
				t.Errorf("Rotate() = %v, want %v", tt.nums, tt.want)
			}
		})
	}
}

func TestFindMinRotated(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"rotated", []int{3, 4, 5, 1, 2}, 1},
		{"not rotated", []int{1, 2, 3, 4, 5}, 1},
		{"single", []int{5}, 5},
		{"two", []int{2, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMinRotated(tt.nums)
			if got != tt.want {
				t.Errorf("FindMinRotated(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}

func TestSubarraySum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"basic", []int{1, 1, 1}, 2, 2},
		{"single", []int{1, 2, 3}, 3, 2},
		{"negative", []int{1, -1, 1}, 1, 3},
		{"zero k", []int{0, 0, 0}, 0, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SubarraySum(tt.nums, tt.k)
			if got != tt.want {
				t.Errorf("SubarraySum(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}

func TestLongestConsecutive(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{100, 4, 200, 1, 3, 2}, 4},
		{"single", []int{1}, 1},
		{"empty", []int{}, 0},
		{"duplicates", []int{1, 2, 0, 1}, 3},
		{"all same", []int{0, 0, 0}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestConsecutive(tt.nums)
			if got != tt.want {
				t.Errorf("LongestConsecutive(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}

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

func TestIsValidSudoku(t *testing.T) {
	valid := [][]byte{
		{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
		{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
		{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
		{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
		{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
		{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
		{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
		{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
		{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
	}
	if !IsValidSudoku(valid) {
		t.Error("expected valid sudoku")
	}
}

func TestMajorityElement(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{3, 2, 3}, 3},
		{"longer", []int{2, 2, 1, 1, 1, 2, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MajorityElement(tt.nums)
			if got != tt.want {
				t.Errorf("MajorityElement(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}

func TestMergeIntervals(t *testing.T) {
	tests := []struct {
		name      string
		intervals [][]int
		want      [][]int
	}{
		{"basic", [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}, [][]int{{1, 6}, {8, 10}, {15, 18}}},
		{"overlap all", [][]int{{1, 4}, {4, 5}}, [][]int{{1, 5}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeIntervals(tt.intervals)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeIntervals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertInterval(t *testing.T) {
	tests := []struct {
		name        string
		intervals   [][]int
		newInterval []int
		want        [][]int
	}{
		{"basic", [][]int{{1, 3}, {6, 9}}, []int{2, 5}, [][]int{{1, 5}, {6, 9}}},
		{"merge multiple", [][]int{{1, 2}, {3, 5}, {6, 7}, {8, 10}, {12, 16}}, []int{4, 8}, [][]int{{1, 2}, {3, 10}, {12, 16}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InsertInterval(tt.intervals, tt.newInterval)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEraseOverlapIntervals(t *testing.T) {
	tests := []struct {
		name      string
		intervals [][]int
		want      int
	}{
		{"basic", [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}}, 1},
		{"all overlap", [][]int{{1, 2}, {1, 2}, {1, 2}}, 2},
		{"none overlap", [][]int{{1, 2}, {2, 3}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EraseOverlapIntervals(tt.intervals)
			if got != tt.want {
				t.Errorf("EraseOverlapIntervals() = %d, want %d", got, tt.want)
			}
		})
	}
}
