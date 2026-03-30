package graphs

import "testing"

func TestFindOrder(t *testing.T) {
	tests := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
		wantLen       int
	}{
		{"no prereqs", 2, [][]int{}, 2},
		{"linear", 4, [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}}, 4},
		{"cycle", 2, [][]int{{1, 0}, {0, 1}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindOrder(tt.numCourses, tt.prerequisites)
			if len(got) != tt.wantLen {
				t.Errorf("FindOrder(%d, %v) returned %d courses, want %d",
					tt.numCourses, tt.prerequisites, len(got), tt.wantLen)
			}
		})
	}
}
