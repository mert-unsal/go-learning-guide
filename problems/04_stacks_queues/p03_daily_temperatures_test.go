package stacks_queues

import (
	"reflect"
	"testing"
)

func TestDailyTemperatures(t *testing.T) {
	tests := []struct {
		name  string
		temps []int
		want  []int
	}{
		{"basic", []int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
		{"all decreasing", []int{5, 4, 3, 2, 1}, []int{0, 0, 0, 0, 0}},
		{"all increasing", []int{1, 2, 3, 4, 5}, []int{1, 1, 1, 1, 0}},
		{"single", []int{30}, []int{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DailyTemperatures(tt.temps)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DailyTemperatures(%v) = %v, want %v", tt.temps, got, tt.want)
			}
		})
	}
}
