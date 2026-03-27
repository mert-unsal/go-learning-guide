package arrays

import "testing"

func TestCountingValleys(t *testing.T) {
	if got := CountingValleys("UDDDUDUU"); got != 1 {
		t.Errorf("got %d want 1", got)
	}
	if got := CountingValleys("DDUUUUDD"); got != 1 {
		t.Errorf("got %d want 1", got)
	}
}
