package two_pointers

import "testing"

func TestBagOfTokensScore(t *testing.T) {
	tests := []struct {
		name   string
		tokens []int
		power  int
		want   int
	}{
		{"example 1", []int{100}, 50, 0},
		{"example 2", []int{200, 100}, 150, 1},
		{"example 3", []int{100, 200, 300, 400}, 200, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BagOfTokensScore(tt.tokens, tt.power)
			if got != tt.want {
				t.Errorf("BagOfTokensScore(%v, %d) = %d, want %d", tt.tokens, tt.power, got, tt.want)
			}
		})
	}
}
