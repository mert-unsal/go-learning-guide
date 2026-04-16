package heap_priority_queue

// ============================================================
// PROBLEM 2: Find Median from Data Stream (LeetCode #295) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   The median is the middle value in an ordered integer list.
//   If the size is even, the median is the average of the two
//   middle values. Design a data structure that supports adding
//   integers from a data stream and finding the median of all
//   elements added so far.
//
// ── NewMedianFinder ──
// RETURN:
//   *MedianFinder — initialized MedianFinder instance
//
// ── AddNum ──
// PARAMETERS:
//   num int — integer to add from the data stream
//
// ── FindMedian ──
// RETURN:
//   float64 — the median of all elements added so far
//
// CONSTRAINTS:
//   • -10^5 <= num <= 10^5
//   • There will be at least one element before calling FindMedian
//   • At most 5 * 10^4 calls to AddNum and FindMedian
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   AddNum(1), AddNum(2) → FindMedian() = 1.5
//   AddNum(3)            → FindMedian() = 2.0
//   Why: After [1,2] median is (1+2)/2=1.5; after [1,2,3] median is 2.
//
// Example 2:
//   AddNum(6), AddNum(10), AddNum(2), AddNum(6)
//   FindMedian() = 6.0
//   Why: Sorted [2,6,6,10], median = (6+6)/2 = 6.0.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Two-heap approach: maxHeap for smaller half, minHeap for larger half
// • Balance heaps so their sizes differ by at most 1
// • Median is top of maxHeap, or average of both tops if sizes are equal
// • Target: O(log n) per AddNum, O(1) per FindMedian, O(n) space

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
