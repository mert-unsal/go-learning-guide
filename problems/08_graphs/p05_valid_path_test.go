package graphs

import "testing"

func TestValidPath(t *testing.T) {
	tests := []struct {
		name        string
		n           int
		edges       [][]int
		source      int
		destination int
		want        bool
	}{
		{"direct edge", 3, [][]int{{0, 1}, {1, 2}}, 0, 2, true},
		{"no path", 6, [][]int{{0, 1}, {0, 2}, {3, 5}, {5, 4}, {4, 3}}, 0, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidPath(tt.n, tt.edges, tt.source, tt.destination)
			if got != tt.want {
				t.Errorf("ValidPath(%d, %v, %d, %d) = %v, want %v",
					tt.n, tt.edges, tt.source, tt.destination, got, tt.want)
			}
		})
	}
}
