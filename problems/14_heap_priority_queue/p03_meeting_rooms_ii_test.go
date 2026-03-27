package heap_priority_queue

import "testing"

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
