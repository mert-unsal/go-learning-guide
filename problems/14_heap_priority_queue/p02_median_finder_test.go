package heap_priority_queue

import "testing"

func TestMedianFinder(t *testing.T) {
	mf := NewMedianFinder()
	mf.AddNum(1)
	mf.AddNum(2)
	if got := mf.FindMedian(); got != 1.5 {
		t.Errorf("FindMedian() = %f, want 1.5", got)
	}
	mf.AddNum(3)
	if got := mf.FindMedian(); got != 2.0 {
		t.Errorf("FindMedian() = %f, want 2.0", got)
	}
}
