package trees

import (
	"reflect"
	"testing"
)

func TestLevelOrder(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want [][]int
	}{
		{"normal", []int{3, 9, 20, 0, 0, 15, 7}, [][]int{{3}, {9, 20}, {15, 7}}},
		{"single", []int{1}, [][]int{{1}}},
		{"empty", []int{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LevelOrder(newTree(tt.vals))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LevelOrder = %v, want %v", got, tt.want)
			}
		})
	}
}
