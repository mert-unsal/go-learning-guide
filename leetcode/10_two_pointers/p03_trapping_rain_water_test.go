package two_pointers

import "testing"

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
