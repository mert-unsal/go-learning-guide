package strings_problems

import (
	"sort"
	"testing"
)

func TestGroupAnagrams(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want [][]string
	}{
		{
			name: "standard example",
			strs: []string{"eat", "tea", "tan", "ate", "nat", "bat"},
			want: [][]string{{"eat", "tea", "ate"}, {"tan", "nat"}, {"bat"}},
		},
		{
			name: "empty string",
			strs: []string{""},
			want: [][]string{{""}},
		},
		{
			name: "single element",
			strs: []string{"a"},
			want: [][]string{{"a"}},
		},
		{
			name: "no anagrams",
			strs: []string{"abc", "def", "ghi"},
			want: [][]string{{"abc"}, {"def"}, {"ghi"}},
		},
		{
			name: "all anagrams",
			strs: []string{"abc", "bca", "cab"},
			want: [][]string{{"abc", "bca", "cab"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupAnagrams(tt.strs)
			if !groupsEqual(got, tt.want) {
				t.Errorf("GroupAnagrams(%v) = %v, want %v", tt.strs, got, tt.want)
			}
		})
	}
}

// groupsEqual compares two [][]string ignoring group order and inner order.
func groupsEqual(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	normalize := func(groups [][]string) []string {
		var result []string
		for _, g := range groups {
			sorted := make([]string, len(g))
			copy(sorted, g)
			sort.Strings(sorted)
			result = append(result, join(sorted))
		}
		sort.Strings(result)
		return result
	}
	na, nb := normalize(a), normalize(b)
	for i := range na {
		if na[i] != nb[i] {
			return false
		}
	}
	return true
}

func join(ss []string) string {
	r := ""
	for i, s := range ss {
		if i > 0 {
			r += ","
		}
		r += s
	}
	return r
}
