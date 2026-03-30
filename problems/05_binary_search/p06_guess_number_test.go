package binary_search

import "testing"

func TestGuessNumber(t *testing.T) {
	tests := []struct {
		name   string
		n      int
		picked int
	}{
		{"pick 6 from 10", 10, 6},
		{"pick 1 from 1", 1, 1},
		{"pick 2 from 2", 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guessFn := func(num int) int {
				if num > tt.picked {
					return -1
				} else if num < tt.picked {
					return 1
				}
				return 0
			}
			got := GuessNumber(tt.n, guessFn)
			if got != tt.picked {
				t.Errorf("GuessNumber(%d) = %v, want %v", tt.n, got, tt.picked)
			}
		})
	}
}
