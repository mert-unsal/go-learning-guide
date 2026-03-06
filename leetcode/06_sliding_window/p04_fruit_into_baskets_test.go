package sliding_window

import "testing"

func TestTotalFruit(t *testing.T) {
	tests := []struct {
		name   string
		fruits []int
		want   int
	}{
		{"basic", []int{1, 2, 1}, 3},
		{"three types", []int{0, 1, 2, 2}, 3},
		{"longer", []int{1, 2, 3, 2, 2}, 4},
		{"single type", []int{1, 1, 1}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TotalFruit(tt.fruits)
			if got != tt.want {
				t.Errorf("TotalFruit(%v) = %d, want %d", tt.fruits, got, tt.want)
			}
		})
	}
}
