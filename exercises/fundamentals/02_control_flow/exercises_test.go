package control_flow

import (
	"fmt"
	"testing"
)

func TestFizzBuzzSwitch(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{1, "1"}, {3, "Fizz"}, {5, "Buzz"}, {15, "FizzBuzz"},
		{9, "Fizz"}, {10, "Buzz"}, {30, "FizzBuzz"}, {7, "7"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("FizzBuzzSwitch(%d)", tt.n), func(t *testing.T) {
			got := FizzBuzzSwitch(tt.n)
			if got != tt.want {
				t.Errorf("❌ FizzBuzzSwitch(%d) = %q, want %q", tt.n, got, tt.want)
			} else {
				t.Logf("✅ FizzBuzzSwitch(%d) = %q", tt.n, got)
			}
		})
	}
}

func TestSumTo(t *testing.T) {
	tests := []struct{ n, want int }{
		{1, 1}, {5, 15}, {10, 55}, {100, 5050},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("SumTo(%d)", tt.n), func(t *testing.T) {
			got := SumTo(tt.n)
			if got != tt.want {
				t.Errorf("❌ SumTo(%d) = %d, want %d", tt.n, got, tt.want)
			} else {
				t.Logf("✅ SumTo(%d) = %d", tt.n, got)
			}
		})
	}
}

func TestCountVowels(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"hello", 2}, {"AEIOU", 5}, {"rhythm", 0}, {"Go is fun", 3},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("CountVowels(%q)", tt.s), func(t *testing.T) {
			got := CountVowels(tt.s)
			if got != tt.want {
				t.Errorf("❌ CountVowels(%q) = %d, want %d  ← Hint: handle uppercase too", tt.s, got, tt.want)
			} else {
				t.Logf("✅ CountVowels(%q) = %d", tt.s, got)
			}
		})
	}
}

func TestIsPrime(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{2, true}, {3, true}, {4, false}, {17, true},
		{1, false}, {0, false}, {97, true}, {100, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("IsPrime(%d)", tt.n), func(t *testing.T) {
			got := IsPrime(tt.n)
			if got != tt.want {
				t.Errorf("❌ IsPrime(%d) = %v, want %v  ← Hint: check divisors up to sqrt(n)", tt.n, got, tt.want)
			} else {
				t.Logf("✅ IsPrime(%d) = %v", tt.n, got)
			}
		})
	}
}

func TestDeferOrder(t *testing.T) {
	got := DeferOrder()
	want := []string{"third", "second", "first"}
	if len(got) != 3 {
		t.Fatalf("❌ DeferOrder() returned %d elements, want 3  ← Hint: defers run LIFO", len(got))
	}
	allPass := true
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("❌ DeferOrder()[%d] = %q, want %q", i, got[i], want[i])
			allPass = false
		}
	}
	if allPass {
		t.Logf("✅ DeferOrder() = %v", got)
	}
}

// ─── ADVANCED TESTS ──────────────────────────────────────────

func TestDeferModifyReturn(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{5, 20},   // 5*2 + 10
		{0, 10},   // 0*2 + 10
		{10, 30},  // 10*2 + 10
		{-3, 4},   // -3*2 + 10
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("DeferModifyReturn(%d)", tt.n), func(t *testing.T) {
			got := DeferModifyReturn(tt.n)
			if got != tt.want {
				t.Errorf("❌ DeferModifyReturn(%d) = %d, want %d\n"+
					"    Hint: use a named return (result int) and defer func() { result += 10 }()\n"+
					"    The defer runs AFTER 'result = n*2' but BEFORE the caller receives it.\n"+
					"    See: learnings/22, Section 3", tt.n, got, tt.want)
			} else {
				t.Logf("✅ DeferModifyReturn(%d) = %d", tt.n, got)
			}
		})
	}
}

func TestDeferArgCapture(t *testing.T) {
	got := DeferArgCapture()
	if got != 0 {
		t.Errorf("❌ DeferArgCapture() = %d, want 0\n"+
			"    Hint: defer captures the argument VALUE at defer-time, not execution-time.\n"+
			"    If counter is 0 when you defer, the deferred function received 0.\n"+
			"    This is different from a closure capturing &counter.\n"+
			"    See: learnings/22, Section 1 — Rule 1", got)
	} else {
		t.Logf("✅ DeferArgCapture() = 0 — defer captured the value at defer-time")
	}
}

func TestDoubleScores(t *testing.T) {
	players := []Player{
		{"Alice", 10},
		{"Bob", 20},
		{"Charlie", 0},
	}
	DoubleScores(players)

	want := []int{20, 40, 0}
	for i, p := range players {
		if p.Score != want[i] {
			t.Errorf("❌ players[%d].Score = %d, want %d\n"+
				"    Hint: range gives you a COPY of each element.\n"+
				"    'for _, p := range players { p.Score *= 2 }' modifies the copy, not the original.\n"+
				"    Use the index: 'for i := range players { players[i].Score *= 2 }'",
				i, p.Score, want[i])
		}
	}
	if players[0].Score == 20 && players[1].Score == 40 && players[2].Score == 0 {
		t.Logf("✅ DoubleScores correctly modifies the original slice elements")
	}
}

func TestRuneValues(t *testing.T) {
	tests := []struct {
		s    string
		want []rune
	}{
		{"Go", []rune{'G', 'o'}},
		{"Go🚀", []rune{'G', 'o', '🚀'}},
		{"日本語", []rune{'日', '本', '語'}},
		{"", nil},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("RuneValues(%q)", tt.s), func(t *testing.T) {
			got := RuneValues(tt.s)
			if len(got) == 0 && len(tt.want) == 0 {
				t.Logf("✅ RuneValues(%q) = %v", tt.s, got)
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("❌ RuneValues(%q) returned %d runes, want %d\n"+
					"    Hint: range over string iterates RUNES, not bytes.\n"+
					"    '🚀' is one rune (4 bytes). 'for _, r := range s' gives you each rune.",
					tt.s, len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("❌ RuneValues(%q)[%d] = %c (U+%04X), want %c (U+%04X)",
						tt.s, i, got[i], got[i], tt.want[i], tt.want[i])
				}
			}
			t.Logf("✅ RuneValues(%q) = %v", tt.s, got)
		})
	}
}

func TestFindInMatrix(t *testing.T) {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	tests := []struct {
		target   int
		wantRow  int
		wantCol  int
	}{
		{5, 1, 1},
		{1, 0, 0},
		{9, 2, 2},
		{42, -1, -1},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("FindInMatrix(%d)", tt.target), func(t *testing.T) {
			r, c := FindInMatrix(matrix, tt.target)
			if r != tt.wantRow || c != tt.wantCol {
				t.Errorf("❌ FindInMatrix(matrix, %d) = (%d, %d), want (%d, %d)\n"+
					"    Hint: use a labeled break to exit both loops:\n"+
					"    outer: for i := range matrix { for j := range matrix[i] { break outer } }",
					tt.target, r, c, tt.wantRow, tt.wantCol)
			} else {
				t.Logf("✅ FindInMatrix(matrix, %d) = (%d, %d)", tt.target, r, c)
			}
		})
	}
}

func TestTypeDescribe(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want string
	}{
		{"int", 42, "int: 42"},
		{"string", "hello", "string: hello"},
		{"bool", true, "bool: true"},
		{"slice", []int{1, 2, 3}, "slice: len=3"},
		{"nil", nil, "nil"},
		{"float", 3.14, "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TypeDescribe(tt.v)
			if got != tt.want {
				t.Errorf("❌ TypeDescribe(%v) = %q, want %q\n"+
					"    Hint: use a type switch: switch v := v.(type) { case int: ... }",
					tt.v, got, tt.want)
			} else {
				t.Logf("✅ TypeDescribe(%v) = %q", tt.v, got)
			}
		})
	}
}

func TestSquares(t *testing.T) {
	tests := []struct {
		n    int
		want []int
	}{
		{0, nil},
		{1, []int{0}},
		{4, []int{0, 1, 4, 9}},
		{6, []int{0, 1, 4, 9, 16, 25}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Squares(%d)", tt.n), func(t *testing.T) {
			got := Squares(tt.n)
			if len(got) == 0 && len(tt.want) == 0 {
				t.Logf("✅ Squares(%d) = %v", tt.n, got)
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("❌ Squares(%d) returned %d elements, want %d\n"+
					"    Hint: use 'for i := range n' (Go 1.22+) to iterate 0..n-1",
					tt.n, len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("❌ Squares(%d)[%d] = %d, want %d", tt.n, i, got[i], tt.want[i])
				}
			}
			t.Logf("✅ Squares(%d) = %v", tt.n, got)
		})
	}
}
