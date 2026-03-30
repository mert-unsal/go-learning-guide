package dynamic_prog

import "testing"

func TestCoinChangeII(t *testing.T) {
	tests := []struct {
		name   string
		amount int
		coins  []int
		want   int
	}{
		{"basic", 5, []int{1, 2, 5}, 4},
		{"impossible", 3, []int{2}, 0},
		{"zero amount", 0, []int{7}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CoinChangeII(tt.amount, tt.coins)
			if got != tt.want {
				t.Errorf("CoinChangeII(%d, %v) = %d, want %d", tt.amount, tt.coins, got, tt.want)
			}
		})
	}
}
