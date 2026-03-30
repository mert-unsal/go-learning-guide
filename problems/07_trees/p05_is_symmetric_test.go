package trees

import "testing"

func TestIsSymmetric(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want bool
	}{
		{"symmetric", []int{1, 2, 2, 3, 4, 4, 3}, true},
		{"not symmetric", []int{1, 2, 2, 0, 3, 0, 3}, false},
		{"single", []int{1}, true},
		{"empty", []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSymmetric(newTree(tt.vals))
			if got != tt.want {
				t.Errorf("IsSymmetric = %v, want %v", got, tt.want)
			}
		})
	}
}
