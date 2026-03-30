package heap_priority_queue

// PROBLEM 2: Find Median from Data Stream (LeetCode #295) — HARD
// Two-heap approach: maxHeap for smaller half, minHeap for larger half.

type MedianFinder struct {
	lo *IntMaxHeap
	hi *IntMinHeap
}

func NewMedianFinder() *MedianFinder {
	return nil
}
func (mf *MedianFinder) AddNum(num int) {
}
func (mf *MedianFinder) FindMedian() float64 {
	return 0
}
