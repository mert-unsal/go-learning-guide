package trees

import "testing"

func TestIsSameTree(t *testing.T) {
	tests := []struct {
		name string
		p    []int
		q    []int
		want bool
	}{
		{"same", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"different", []int{1, 2}, []int{1, 0, 2}, false},
		{"both empty", []int{}, []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSameTree(newTree(tt.p), newTree(tt.q))
			if got != tt.want {
				t.Errorf("IsSameTree = %v, want %v", got, tt.want)
			}
		})
	}
}
