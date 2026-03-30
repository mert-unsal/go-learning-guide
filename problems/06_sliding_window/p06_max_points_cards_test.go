package sliding_window

import "testing"

func TestMaxScore(t *testing.T) {
	tests := []struct {
		name       string
		cardPoints []int
		k          int
		want       int
	}{
		{"basic", []int{1, 2, 3, 4, 5, 6, 1}, 3, 12},
		{"all same", []int{2, 2, 2}, 2, 4},
		{"take all", []int{9, 7, 7, 9, 7, 7, 9}, 7, 55},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxScore(tt.cardPoints, tt.k)
			if got != tt.want {
				t.Errorf("MaxScore(%v, %d) = %v, want %v", tt.cardPoints, tt.k, got, tt.want)
			}
		})
	}
}
