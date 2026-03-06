package heap_priority_queue

import "testing"

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
