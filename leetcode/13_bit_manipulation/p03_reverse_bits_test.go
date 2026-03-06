package bit_manipulation

import "testing"

func TestReverseBits(t *testing.T) {
	if got := ReverseBits(43261596); got != 964176192 {
		t.Errorf("ReverseBits(43261596) = %d, want 964176192", got)
	}
}
