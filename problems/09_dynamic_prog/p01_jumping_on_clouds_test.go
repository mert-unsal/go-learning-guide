package dynamic_prog

import "testing"

func TestJumpingOnClouds(t *testing.T) {
	if got := JumpingOnClouds([]int{0, 0, 1, 0, 0, 1, 0}); got != 4 {
		t.Errorf("got %d want 4", got)
	}
}
