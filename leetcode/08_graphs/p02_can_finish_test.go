package graphs

import "testing"

func TestCanFinish(t *testing.T) {
	tests := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
		want          bool
	}{
		{"no prereqs", 2, [][]int{}, true},
		{"linear chain", 2, [][]int{{1, 0}}, true},
		{"simple cycle", 2, [][]int{{1, 0}, {0, 1}}, false},
		{"longer no cycle", 4, [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}}, true},
		{"longer cycle", 3, [][]int{{0, 1}, {1, 2}, {2, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanFinish(tt.numCourses, tt.prerequisites)
			if got != tt.want {
				t.Errorf("CanFinish(%d, %v) = %v, want %v", tt.numCourses, tt.prerequisites, got, tt.want)
			}
		})
	}
}
