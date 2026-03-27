package arrays

import (
	"reflect"
	"testing"
)

func TestFindDisappearedNumbers(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{"basic", []int{4, 3, 2, 7, 8, 2, 3, 1}, []int{5, 6}},
		{"none missing", []int{1, 2}, []int(nil)},
		{"all missing except one", []int{2, 2}, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindDisappearedNumbers(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindDisappearedNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
