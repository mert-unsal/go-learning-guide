package stacks_queues

import "testing"

func TestCarFleet(t *testing.T) {
	tests := []struct {
		name     string
		target   int
		position []int
		speed    []int
		want     int
	}{
		{"basic", 12, []int{10, 8, 0, 5, 3}, []int{2, 4, 1, 1, 3}, 3},
		{"single car", 10, []int{3}, []int{3}, 1},
		{"all same speed", 100, []int{0, 2, 4}, []int{4, 2, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CarFleet(tt.target, tt.position, tt.speed)
			if got != tt.want {
				t.Errorf("CarFleet(%d, %v, %v) = %v, want %v", tt.target, tt.position, tt.speed, got, tt.want)
			}
		})
	}
}
