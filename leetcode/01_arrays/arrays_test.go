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
		{"mixed", []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}, 6},
		{"all positive", []int{1, 2, 3}, 6},
		{"all negative", []int{-3, -1, -2}, -1},
		{"single", []int{5}, 5},
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
		{"basic", []int{1, 2, 3, 0, 0, 0}, 3, []int{2, 5, 6}, 3, []int{1, 2, 2, 3, 5, 6}},
		{"one empty", []int{1}, 1, []int{}, 0, []int{1}},
		{"first empty", []int{0}, 0, []int{1}, 1, []int{1}},
		{"all same", []int{2, 2, 0, 0}, 2, []int{2, 2}, 2, []int{2, 2, 2, 2}},
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
