package heap_priority_queue

import (
	"reflect"
	"testing"
)

func TestMaxSlidingWindow(t *testing.T) {
	got := MaxSlidingWindow([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3)
	want := []int{3, 3, 5, 5, 6, 7}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MaxSlidingWindow = %v, want %v", got, want)
	}
}
