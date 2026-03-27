package bit_manipulation

import (
	"reflect"
	"testing"
)

func TestCountBits(t *testing.T) {
	got := CountBits(5)
	want := []int{0, 1, 1, 2, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CountBits(5) = %v, want %v", got, want)
	}
}
