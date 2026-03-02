package heap_priority_queue

import (
	"reflect"
	"testing"
)

func TestFindKthLargest(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"basic", []int{3, 2, 1, 5, 6, 4}, 2, 5},
		{"with dups", []int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Copy since quickselect modifies array
			nums := make([]int, len(tt.nums))
			copy(nums, tt.nums)
			got := FindKthLargest(nums, tt.k)
			if got != tt.want {
				t.Errorf("FindKthLargest(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}

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

func TestMinMeetingRooms(t *testing.T) {
	tests := []struct {
		name      string
		intervals [][]int
		want      int
	}{
		{"overlapping", [][]int{{0, 30}, {5, 10}, {15, 20}}, 2},
		{"no overlap", [][]int{{7, 10}, {2, 4}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinMeetingRooms(tt.intervals)
			if got != tt.want {
				t.Errorf("MinMeetingRooms(%v) = %d, want %d", tt.intervals, got, tt.want)
			}
		})
	}
}

func TestLeastInterval(t *testing.T) {
	tests := []struct {
		name  string
		tasks []byte
		n     int
		want  int
	}{
		{"basic", []byte{'A', 'A', 'A', 'B', 'B', 'B'}, 2, 8},
		{"no cooldown", []byte{'A', 'A', 'A', 'B', 'B', 'B'}, 0, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LeastInterval(tt.tasks, tt.n)
			if got != tt.want {
				t.Errorf("LeastInterval(%v, %d) = %d, want %d", tt.tasks, tt.n, got, tt.want)
			}
		})
	}
}

func TestMaxSlidingWindow(t *testing.T) {
	got := MaxSlidingWindow([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3)
	want := []int{3, 3, 5, 5, 6, 7}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MaxSlidingWindow = %v, want %v", got, want)
	}
}
