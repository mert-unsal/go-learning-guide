package trees

import "testing"

func TestGoodNodes(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want int
	}{
		{"basic", []int{3, 1, 4, 3, 0, 1, 5}, 4},
		{"all good", []int{3, 3, 0, 4, 2}, 3},
		{"single", []int{1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GoodNodes(newTree(tt.vals))
			if got != tt.want {
				t.Errorf("GoodNodes = %v, want %v", got, tt.want)
			}
		})
	}
}
