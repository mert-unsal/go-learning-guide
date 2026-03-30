package binary_search

import (
	"reflect"
	"testing"
)

func TestSearchRange(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{"found", []int{5, 7, 7, 8, 8, 10}, 8, []int{3, 4}},
		{"not found", []int{5, 7, 7, 8, 8, 10}, 6, []int{-1, -1}},
		{"empty", []int{}, 0, []int{-1, -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SearchRange(tt.nums, tt.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchRange(%v, %d) = %v, want %v", tt.nums, tt.target, got, tt.want)
			}
		})
	}
}
