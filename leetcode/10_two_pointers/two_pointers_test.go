package two_pointers

import (
	"reflect"
	"testing"
)

func TestMaxArea(t *testing.T) {
	tests := []struct {
		name   string
		height []int
		want   int
	}{
		{"basic", []int{1, 8, 6, 2, 5, 4, 8, 3, 7}, 49},
		{"two walls", []int{1, 1}, 1},
		{"increasing", []int{1, 2, 3, 4, 5}, 6},
		{"decreasing", []int{5, 4, 3, 2, 1}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxArea(tt.height)
			if got != tt.want {
				t.Errorf("MaxArea(%v) = %d, want %d", tt.height, got, tt.want)
			}
		})
	}
}

func TestThreeSum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want [][]int
	}{
		{"classic", []int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{"all zeros", []int{0, 0, 0}, [][]int{{0, 0, 0}}},
		{"no solution", []int{1, 2, 3}, nil},
		{"with duplicates", []int{-2, 0, 0, 2, 2}, [][]int{{-2, 0, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ThreeSum(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ThreeSum(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

func TestTrap(t *testing.T) {
	tests := []struct {
		name   string
		height []int
		want   int
	}{
		{"classic", []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, 6},
		{"simple valley", []int{4, 2, 0, 3, 2, 5}, 9},
		{"no trap", []int{3, 2, 1}, 0},
		{"empty", []int{}, 0},
		{"flat", []int{3, 3, 3}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Trap(tt.height)
			if got != tt.want {
				t.Errorf("Trap(%v) = %d, want %d", tt.height, got, tt.want)
			}
		})
	}
}

func TestMoveZeroes(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"basic", []int{0, 1, 0, 3, 12}, []int{1, 3, 12, 0, 0}},
		{"no zeros", []int{1, 2, 3}, []int{1, 2, 3}},
		{"all zeros", []int{0, 0, 0}, []int{0, 0, 0}},
		{"single zero", []int{0}, []int{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MoveZeroes(tt.input)
			if !reflect.DeepEqual(tt.input, tt.want) {
				t.Errorf("MoveZeroes result = %v, want %v", tt.input, tt.want)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{1, 1, 2}, 2},
		{"multiple dups", []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}, 5},
		{"no dups", []int{1, 2, 3}, 3},
		{"single", []int{1}, 1},
		{"empty", []int{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveDuplicates(tt.nums)
			if got != tt.want {
				t.Errorf("RemoveDuplicates(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
