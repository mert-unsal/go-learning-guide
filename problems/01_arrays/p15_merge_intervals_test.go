package arrays

import (
	"reflect"
	"testing"
)

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
