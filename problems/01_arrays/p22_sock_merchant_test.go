package arrays

import "testing"

func TestSockMerchant(t *testing.T) {
	if got := SockMerchant([]int{10, 20, 20, 10, 10, 30, 50, 10, 20}); got != 3 {
		t.Errorf("got %d want 3", got)
	}
}
