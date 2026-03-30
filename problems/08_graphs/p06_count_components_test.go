package graphs

import "testing"

func TestCountComponents(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		edges [][]int
		want  int
	}{
		{"two components", 5, [][]int{{0, 1}, {1, 2}, {3, 4}}, 2},
		{"one component", 5, [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountComponents(tt.n, tt.edges)
			if got != tt.want {
				t.Errorf("CountComponents(%d, %v) = %d, want %d", tt.n, tt.edges, got, tt.want)
			}
		})
	}
}
