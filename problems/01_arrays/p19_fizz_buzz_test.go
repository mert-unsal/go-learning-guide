package arrays

import (
	"reflect"
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	got := FizzBuzz(15)
	want := []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FizzBuzz(15) = %v, want %v", got, want)
	}
}
