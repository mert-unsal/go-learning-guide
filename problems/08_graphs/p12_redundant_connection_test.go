package graphs

import (
	"reflect"
	"testing"
)

func TestFindRedundantConnection(t *testing.T) {
	tests := []struct {
		name  string
		edges [][]int
		want  []int
	}{
		{"basic", [][]int{{1, 2}, {1, 3}, {2, 3}}, []int{2, 3}},
		{"longer", [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 4}, {1, 5}}, []int{1, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRedundantConnection(tt.edges)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindRedundantConnection(%v) = %v, want %v", tt.edges, got, tt.want)
			}
		})
	}
}
