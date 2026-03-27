package binary_search

import (
	"reflect"
	"testing"
)

func TestClimbingLeaderboard(t *testing.T) {
	got1 := ClimbingLeaderboard([]int{100, 100, 50, 40, 40, 20, 10}, []int{5, 25, 50, 120})
	want1 := []int{6, 4, 2, 1}
	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("test1: got %v want %v", got1, want1)
	}
	got2 := ClimbingLeaderboard([]int{100, 90, 90, 80}, []int{70, 90, 95})
	want2 := []int{4, 2, 2}
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("test2: got %v want %v", got2, want2)
	}
}
