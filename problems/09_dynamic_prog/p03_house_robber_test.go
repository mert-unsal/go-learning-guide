package dynamic_prog

import "testing"

func TestRob(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{1, 2, 3, 1}, 4},
		{"alternate better", []int{2, 7, 9, 3, 1}, 12},
		{"single", []int{5}, 5},
		{"two", []int{1, 2}, 2},
		{"empty", []int{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Rob(tt.nums)
			if got != tt.want {
				t.Errorf("Rob(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
