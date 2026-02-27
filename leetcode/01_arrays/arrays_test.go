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
