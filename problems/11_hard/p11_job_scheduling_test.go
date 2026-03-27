package hard

import "testing"

func TestJobScheduling(t *testing.T) {
	tests := []struct {
		name      string
		startTime []int
		endTime   []int
		profit    []int
		want      int
	}{
		{"example 1", []int{1, 2, 3, 3}, []int{3, 4, 5, 6}, []int{50, 10, 40, 70}, 120},
		{"example 2", []int{1, 2, 3, 4, 6}, []int{3, 5, 10, 6, 9}, []int{20, 20, 100, 70, 60}, 150},
		{"no overlap", []int{1, 4}, []int{3, 6}, []int{10, 20}, 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JobScheduling(tt.startTime, tt.endTime, tt.profit)
			if got != tt.want {
				t.Errorf("JobScheduling() = %d, want %d", got, tt.want)
			}
		})
	}
}
