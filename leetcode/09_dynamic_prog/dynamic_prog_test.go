package dynamic_prog

import "testing"

func TestClimbStairs(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 5},
		{5, 8},
		{10, 89},
	}
	for _, tt := range tests {
		got := ClimbStairs(tt.n)
		if got != tt.want {
			t.Errorf("ClimbStairs(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestCoinChange(t *testing.T) {
	tests := []struct {
		name   string
		coins  []int
		amount int
		want   int
	}{
		{"basic", []int{1, 2, 5}, 11, 3},
		{"impossible", []int{2}, 3, -1},
		{"zero amount", []int{1}, 0, 0},
		{"single coin exact", []int{5}, 5, 1},
		{"large amount", []int{1, 5, 10, 25}, 30, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CoinChange(tt.coins, tt.amount)
			if got != tt.want {
				t.Errorf("CoinChange(%v, %d) = %d, want %d", tt.coins, tt.amount, got, tt.want)
			}
		})
	}
}

func TestRob(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{1, 2, 3, 1}, 4},
		{"alternate better", []int{2, 7, 9, 3, 1}, 12},
		{"single", []int{5}, 5},
		{"two", []int{1, 2}, 2},
		{"empty", []int{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Rob(tt.nums)
			if got != tt.want {
				t.Errorf("Rob(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}

func TestUniquePaths(t *testing.T) {
	tests := []struct {
		m, n int
		want int
	}{
		{3, 7, 28},
		{3, 2, 3},
		{1, 1, 1},
		{7, 3, 28},
	}
	for _, tt := range tests {
		got := UniquePaths(tt.m, tt.n)
		if got != tt.want {
			t.Errorf("UniquePaths(%d, %d) = %d, want %d", tt.m, tt.n, got, tt.want)
		}
	}
}

func TestLongestCommonSubsequence(t *testing.T) {
	tests := []struct {
		name         string
		text1, text2 string
		want         int
	}{
		{"classic", "abcde", "ace", 3},
		{"same string", "abc", "abc", 3},
		{"no common", "abc", "def", 0},
		{"one empty", "", "abc", 0},
		{"partial", "oxcpqrsvwf", "shmtulqrypy", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestCommonSubsequence(tt.text1, tt.text2)
			if got != tt.want {
				t.Errorf("LCS(%q, %q) = %d, want %d", tt.text1, tt.text2, got, tt.want)
			}
		})
	}
}
