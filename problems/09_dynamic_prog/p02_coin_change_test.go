package dynamic_prog

import "testing"

func TestCoinChange(t *testing.T) {
	tests := []struct {
		name   string
		coins  []int
		amount int
		want   int
	}{
		{"basic", []int{1, 2, 5}, 11, 3},
		{"impossible", []int{2}, 3, -1},
		{"zero amount", []int{1}, 0, 0},
		{"single coin exact", []int{5}, 5, 1},
		{"large amount", []int{1, 5, 10, 25}, 30, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CoinChange(tt.coins, tt.amount)
			if got != tt.want {
				t.Errorf("CoinChange(%v, %d) = %d, want %d", tt.coins, tt.amount, got, tt.want)
			}
		})
	}
}
