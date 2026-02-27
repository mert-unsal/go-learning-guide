package maps

import (
	"reflect"
	"sort"
	"testing"
)

func TestCharFrequency(t *testing.T) {
	got := CharFrequency("hello")
	want := map[rune]int{'h': 1, 'e': 1, 'l': 2, 'o': 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("❌ CharFrequency(\"hello\") = %v, want %v  ← Hint: range over string gives runes", got, want)
	} else {
		t.Logf("✅ CharFrequency(\"hello\") = %v", got)
	}
	if len(CharFrequency("")) != 0 {
		t.Error("❌ CharFrequency(\"\") should return empty map")
	} else {
		t.Logf("✅ CharFrequency(\"\") = empty map")
	}
}

func TestGroupByFirstChar(t *testing.T) {
	got := GroupByFirstChar([]string{"ant", "bat", "bee", "ape"})
	for _, v := range got {
		sort.Strings(v)
	}
	if !reflect.DeepEqual(got['a'], []string{"ant", "ape"}) {
		t.Errorf("❌ group['a'] = %v, want [ant ape]", got['a'])
	} else {
		t.Logf("✅ group['a'] = %v", got['a'])
	}
	if !reflect.DeepEqual(got['b'], []string{"bat", "bee"}) {
		t.Errorf("❌ group['b'] = %v, want [bat bee]", got['b'])
	} else {
		t.Logf("✅ group['b'] = %v", got['b'])
	}
}

func TestIsAnagram(t *testing.T) {
	tests := []struct {
		s, tt string
		want  bool
	}{
		{"listen", "silent", true},
		{"hello", "world", false},
		{"anagram", "nagaram", true},
		{"rat", "car", false},
		{"a", "a", true},
		{"ab", "a", false},
	}
	for _, tt := range tests {
		got := IsAnagram(tt.s, tt.tt)
		if got != tt.want {
			t.Errorf("❌ IsAnagram(%q,%q) = %v, want %v  ← Hint: count chars in s, subtract for t", tt.s, tt.tt, got, tt.want)
		} else {
			t.Logf("✅ IsAnagram(%q,%q) = %v", tt.s, tt.tt, got)
		}
	}
}

func TestFirstDuplicate(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{4, 3, 2, 7, 8, 2, 3, 1}, 2},
		{[]int{1, 2, 3}, -1},
		{[]int{1, 1}, 1},
	}
	for _, tt := range tests {
		got := FirstDuplicate(tt.nums)
		if got != tt.want {
			t.Errorf("❌ FirstDuplicate(%v) = %d, want %d  ← Hint: use a map[int]bool as a seen set", tt.nums, got, tt.want)
		} else {
			t.Logf("✅ FirstDuplicate(%v) = %d", tt.nums, got)
		}
	}
}

func TestWordCount(t *testing.T) {
	got := WordCount("go is go")
	if got["go"] != 2 || got["is"] != 1 {
		t.Errorf("❌ WordCount(\"go is go\") = %v, want {go:2 is:1}", got)
	} else {
		t.Logf("✅ WordCount(\"go is go\") = %v", got)
	}
	if len(WordCount("")) != 0 {
		t.Error("❌ WordCount(\"\") should return empty map")
	} else {
		t.Logf("✅ WordCount(\"\") = empty map")
	}
}
