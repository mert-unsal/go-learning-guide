package bit_manipulation

import (
	"reflect"
	"testing"
)

func TestHammingWeight(t *testing.T) {
	tests := []struct {
		name string
		n    uint32
		want int
	}{
		{"11", 11, 3},
		{"128", 128, 1},
		{"max", 4294967293, 31},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HammingWeight(tt.n); got != tt.want {
				t.Errorf("HammingWeight(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}

func TestCountBits(t *testing.T) {
	got := CountBits(5)
	want := []int{0, 1, 1, 2, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CountBits(5) = %v, want %v", got, want)
	}
}

func TestReverseBits(t *testing.T) {
	if got := ReverseBits(43261596); got != 964176192 {
		t.Errorf("ReverseBits(43261596) = %d, want 964176192", got)
	}
}

func TestMissingNumber(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{3, 0, 1}, 2},
		{"zero missing", []int{1}, 0},
		{"last missing", []int{0, 1}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MissingNumber(tt.nums); got != tt.want {
				t.Errorf("MissingNumber(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}

func TestGetSum(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{1, 2, 3},
		{-1, 1, 0},
		{0, 0, 0},
	}
	for _, tt := range tests {
		if got := GetSum(tt.a, tt.b); got != tt.want {
			t.Errorf("GetSum(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestSingleNumber(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{2, 2, 1}, 1},
		{[]int{4, 1, 2, 1, 2}, 4},
	}
	for _, tt := range tests {
		if got := SingleNumber(tt.nums); got != tt.want {
			t.Errorf("SingleNumber(%v) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{1, true},
		{16, true},
		{3, false},
		{0, false},
		{-1, false},
	}
	for _, tt := range tests {
		if got := IsPowerOfTwo(tt.n); got != tt.want {
			t.Errorf("IsPowerOfTwo(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}
