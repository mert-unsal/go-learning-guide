package arrays

import "testing"

func TestMaxProfit(t *testing.T) {
	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{"normal", []int{7, 1, 5, 3, 6, 4}, 5},
		{"no profit", []int{7, 6, 4, 3, 1}, 0},
		{"single day", []int{5}, 0},
		{"two days profit", []int{1, 5}, 4},
		{"empty", []int{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxProfit(tt.prices)
			if got != tt.want {
				t.Errorf("MaxProfit(%v) = %d, want %d", tt.prices, got, tt.want)
			}
		})
	}
}
