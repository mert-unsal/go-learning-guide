package sliding_window

import "testing"

func TestMinSubArrayLen(t *testing.T) {
	tests := []struct {
		name   string
		target int
		nums   []int
		want   int
	}{
		{"basic", 7, []int{2, 3, 1, 2, 4, 3}, 2},
		{"exact", 4, []int{1, 4, 4}, 1},
		{"not possible", 11, []int{1, 1, 1, 1, 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinSubArrayLen(tt.target, tt.nums)
			if got != tt.want {
				t.Errorf("MinSubArrayLen(%d, %v) = %v, want %v", tt.target, tt.nums, got, tt.want)
			}
		})
	}
}
