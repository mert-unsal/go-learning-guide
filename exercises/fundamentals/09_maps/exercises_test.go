package maps

import (
	"fmt"
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

func TestTopTwoFrequent(t *testing.T) {
	tests := []struct {
		nums []int
		want []int // sorted descending by frequency
	}{
		{[]int{1, 1, 1, 2, 2, 3}, []int{1, 2}},
		{[]int{5, 5, 3, 3, 3}, []int{3, 5}},
		{[]int{1, 2}, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("TopTwoFrequent(%v)", tt.nums), func(t *testing.T) {
			got := TopTwoFrequent(tt.nums)
			sort.Ints(got)
			sorted := make([]int, len(tt.want))
			copy(sorted, tt.want)
			sort.Ints(sorted)
			if !reflect.DeepEqual(got, sorted) {
				t.Errorf("❌ TopTwoFrequent(%v) = %v, want %v\n\t\t"+
					"Hint: Build a frequency map, then scan for top-2 by count",
					tt.nums, got, sorted)
			} else {
				t.Logf("✅ TopTwoFrequent(%v) = %v", tt.nums, got)
			}
		})
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
		t.Run(fmt.Sprintf("IsAnagram(%q,%q)", tt.s, tt.tt), func(t *testing.T) {
			got := IsAnagram(tt.s, tt.tt)
			if got != tt.want {
				t.Errorf("❌ IsAnagram(%q,%q) = %v, want %v  ← Hint: count chars in s, subtract for t", tt.s, tt.tt, got, tt.want)
			} else {
				t.Logf("✅ IsAnagram(%q,%q) = %v", tt.s, tt.tt, got)
			}
		})
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
		t.Run(fmt.Sprintf("FirstDuplicate(%v)", tt.nums), func(t *testing.T) {
			got := FirstDuplicate(tt.nums)
			if got != tt.want {
				t.Errorf("❌ FirstDuplicate(%v) = %d, want %d  ← Hint: use a map[int]bool as a seen set", tt.nums, got, tt.want)
			} else {
				t.Logf("✅ FirstDuplicate(%v) = %d", tt.nums, got)
			}
		})
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

// ── Tests for Go Map Internals (Exercises 7-12) ──

func TestNilMapRead(t *testing.T) {
	t.Run("nil_map", func(t *testing.T) {
		v, ok := NilMapRead(nil, "any")
		if v != 0 || ok != false {
			t.Errorf("❌ NilMapRead(nil, \"any\") = (%d, %v), want (0, false)\n\t\t"+
				"Hint: Reading from a nil map is safe — returns zero value. "+
				"Use comma-ok: v, ok := m[key]. See learnings/02 §1",
				v, ok)
		} else {
			t.Logf("✅ NilMapRead(nil, \"any\") = (0, false)")
		}
	})

	t.Run("key_exists", func(t *testing.T) {
		m := map[string]int{"score": 42}
		v, ok := NilMapRead(m, "score")
		if v != 42 || !ok {
			t.Errorf("❌ NilMapRead(m, \"score\") = (%d, %v), want (42, true)", v, ok)
		} else {
			t.Logf("✅ NilMapRead(m, \"score\") = (42, true)")
		}
	})

	t.Run("zero_value_vs_missing", func(t *testing.T) {
		m := map[string]int{"count": 0}
		_, okExists := NilMapRead(m, "count")
		_, okMissing := NilMapRead(m, "nope")
		if !okExists {
			t.Errorf("❌ key \"count\" exists with value 0, but ok=false\n\t\t"+
				"Hint: comma-ok distinguishes 'key exists with zero value' from 'key missing'")
		}
		if okMissing {
			t.Errorf("❌ key \"nope\" doesn't exist, but ok=true")
		}
		if okExists && !okMissing {
			t.Logf("✅ comma-ok correctly distinguishes zero value from missing key")
		}
	})
}

func TestInvertMap(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]int
		want  map[int][]string
	}{
		{
			"unique_values",
			map[string]int{"a": 1, "b": 2},
			map[int][]string{1: {"a"}, 2: {"b"}},
		},
		{
			"duplicate_values",
			map[string]int{"a": 1, "b": 2, "c": 1},
			map[int][]string{1: {"a", "c"}, 2: {"b"}},
		},
		{
			"empty",
			map[string]int{},
			map[int][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InvertMap(tt.input)
			if got == nil && len(tt.want) == 0 {
				return // both empty
			}
			// Sort slices for comparison (map iteration is random)
			for k := range got {
				sort.Strings(got[k])
			}
			for k := range tt.want {
				sort.Strings(tt.want[k])
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("❌ InvertMap(%v) = %v, want %v\n\t\t"+
					"Hint: For each k,v in input, append k to result[v]. "+
					"Multiple keys may share the same value → use []string",
					tt.input, got, tt.want)
			} else {
				t.Logf("✅ InvertMap = %v", got)
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	sum := func(a, b int) int { return a + b }
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	tests := []struct {
		name    string
		a, b    map[string]int
		resolve func(int, int) int
		want    map[string]int
	}{
		{
			"sum_collision",
			map[string]int{"x": 1, "y": 2},
			map[string]int{"x": 10, "z": 3},
			sum,
			map[string]int{"x": 11, "y": 2, "z": 3},
		},
		{
			"max_collision",
			map[string]int{"a": 5},
			map[string]int{"a": 3, "b": 7},
			max,
			map[string]int{"a": 5, "b": 7},
		},
		{
			"no_collision",
			map[string]int{"a": 1},
			map[string]int{"b": 2},
			sum,
			map[string]int{"a": 1, "b": 2},
		},
		{
			"empty_maps",
			map[string]int{},
			map[string]int{},
			sum,
			map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeMaps(tt.a, tt.b, tt.resolve)
			if got == nil && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("❌ MergeMaps(%v, %v) = %v, want %v\n\t\t"+
					"Hint: Copy all of a into result, then iterate b. "+
					"If key exists in result, use resolve(existing, new). "+
					"Otherwise just set it",
					tt.a, tt.b, got, tt.want)
			} else {
				t.Logf("✅ MergeMaps = %v", got)
			}
		})
	}
}

func TestSetDifference(t *testing.T) {
	tests := []struct {
		name string
		a, b map[string]bool
		want []string
	}{
		{
			"basic",
			map[string]bool{"go": true, "java": true, "python": true},
			map[string]bool{"java": true},
			[]string{"go", "python"},
		},
		{
			"no_overlap",
			map[string]bool{"a": true, "b": true},
			map[string]bool{"c": true},
			[]string{"a", "b"},
		},
		{
			"full_overlap",
			map[string]bool{"x": true},
			map[string]bool{"x": true},
			[]string{},
		},
		{
			"empty_a",
			map[string]bool{},
			map[string]bool{"x": true},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetDifference(tt.a, tt.b)
			if got == nil {
				got = []string{}
			}
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("❌ SetDifference(%v, %v) = %v, want %v\n\t\t"+
					"Hint: Iterate a, check if key is in b with: if b[key] { skip }. "+
					"map[T]bool used as set — checking membership is just m[key]. "+
					"Sort the result for deterministic output",
					tt.a, tt.b, got, tt.want)
			} else {
				t.Logf("✅ SetDifference = %v", got)
			}
		})
	}
}

func TestUniqueValues(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]int
		want  []int
	}{
		{
			"duplicates",
			map[string]int{"a": 3, "b": 1, "c": 3, "d": 2},
			[]int{1, 2, 3},
		},
		{
			"all_unique",
			map[string]int{"x": 10, "y": 20},
			[]int{10, 20},
		},
		{
			"empty",
			map[string]int{},
			[]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UniqueValues(tt.input)
			if got == nil {
				got = []int{}
			}
			sort.Ints(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("❌ UniqueValues(%v) = %v, want %v\n\t\t"+
					"Hint: Use a set (map[int]bool) to deduplicate values. "+
					"Map iteration order is RANDOM — sort for deterministic results. "+
					"See learnings/02 — map iteration randomness is deliberate",
					tt.input, got, tt.want)
			} else {
				t.Logf("✅ UniqueValues = %v", got)
			}
		})
	}
}

func TestMapEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b map[string]int
		want bool
	}{
		{"equal", map[string]int{"a": 1, "b": 2}, map[string]int{"b": 2, "a": 1}, true},
		{"different_value", map[string]int{"a": 1}, map[string]int{"a": 2}, false},
		{"extra_key", map[string]int{"a": 1}, map[string]int{"a": 1, "b": 2}, false},
		{"both_empty", map[string]int{}, map[string]int{}, true},
		{"both_nil", nil, nil, true},
		{"nil_vs_empty", nil, map[string]int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapEqual(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("❌ MapEqual(%v, %v) = %v, want %v\n\t\t"+
					"Hint: Maps can't use == (compile error). Check len first, "+
					"then iterate a and verify each key exists in b with same value. "+
					"Use comma-ok to distinguish missing from zero-value",
					tt.a, tt.b, got, tt.want)
			} else {
				t.Logf("✅ MapEqual(%v, %v) = %v", tt.a, tt.b, got)
			}
		})
	}
}

