package hard

import (
	"reflect"
	"testing"
)

// ============================================================
// TESTS — Hard Level Problems
// ============================================================

// ---- Helpers ----

func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func listToSlice(head *ListNode) []int {
	var res []int
	for cur := head; cur != nil; cur = cur.Next {
		res = append(res, cur.Val)
	}
	return res
}

func makeTree(val int, left, right *TreeNode) *TreeNode {
	return &TreeNode{Val: val, Left: left, Right: right}
}

// ---- Tests ----

func TestMergeKLists(t *testing.T) {
	tests := []struct {
		name  string
		lists [][]int
		want  []int
	}{
		{"three lists", [][]int{{1, 4, 5}, {1, 3, 4}, {2, 6}}, []int{1, 1, 2, 3, 4, 4, 5, 6}},
		{"empty lists", [][]int{{}, {}}, nil},
		{"single list", [][]int{{1, 2, 3}}, []int{1, 2, 3}},
		{"one empty", [][]int{{1}, {}, {2}}, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lists := make([]*ListNode, len(tt.lists))
			for i, l := range tt.lists {
				lists[i] = makeList(l)
			}
			got := listToSlice(MergeKLists(lists))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeKLists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrap(t *testing.T) {
	tests := []struct {
		name   string
		height []int
		want   int
	}{
		{"example 1", []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, 6},
		{"example 2", []int{4, 2, 0, 3, 2, 5}, 9},
		{"no water", []int{1, 2, 3}, 0},
		{"empty", []int{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Trap(tt.height)
			if got != tt.want {
				t.Errorf("Trap(%v) = %d, want %d", tt.height, got, tt.want)
			}
		})
	}
}

func TestLadderLength(t *testing.T) {
	tests := []struct {
		name      string
		beginWord string
		endWord   string
		wordList  []string
		want      int
	}{
		{"basic", "hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"}, 5},
		{"no path", "hit", "cog", []string{"hot", "dot", "dog", "lot", "log"}, 0},
		{"direct", "a", "c", []string{"a", "b", "c"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LadderLength(tt.beginWord, tt.endWord, tt.wordList)
			if got != tt.want {
				t.Errorf("LadderLength(%q, %q) = %d, want %d", tt.beginWord, tt.endWord, got, tt.want)
			}
		})
	}
}

func TestLongestValidParentheses(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"mixed", ")()())", 4},
		{"open left", "(()", 2},
		{"empty", "", 0},
		{"all valid", "()()", 4},
		{"nested", "((()))", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestValidParentheses(tt.s)
			if got != tt.want {
				t.Errorf("LongestValidParentheses(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

func TestJump(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"example 1", []int{2, 3, 1, 1, 4}, 2},
		{"example 2", []int{2, 3, 0, 1, 4}, 2},
		{"single", []int{0}, 0},
		{"two", []int{1, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Jump(tt.nums)
			if got != tt.want {
				t.Errorf("Jump(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}

func TestSolveNQueens(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		wantCount int
	}{
		{"n=1", 1, 1},
		{"n=4", 4, 2},
		{"n=5", 5, 10},
		{"n=6", 6, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveNQueens(tt.n)
			if len(got) != tt.wantCount {
				t.Errorf("SolveNQueens(%d) returned %d solutions, want %d", tt.n, len(got), tt.wantCount)
			}
		})
	}
}

func TestSerializeDeserialize(t *testing.T) {
	tests := []struct {
		name string
		root *TreeNode
	}{
		{"simple tree", makeTree(1, makeTree(2, nil, nil), makeTree(3, makeTree(4, nil, nil), makeTree(5, nil, nil)))},
		{"nil", nil},
		{"single node", makeTree(42, nil, nil)},
		{"left skewed", makeTree(1, makeTree(2, makeTree(3, nil, nil), nil), nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serialized := Serialize(tt.root)
			deserialized := Deserialize(serialized)
			// Re-serialize to compare
			if Serialize(tt.root) != Serialize(deserialized) {
				t.Errorf("Serialize/Deserialize round-trip failed for %q", tt.name)
			}
		})
	}
}

func TestMinWindow(t *testing.T) {
	tests := []struct {
		name string
		s, t string
		want string
	}{
		{"basic", "ADOBECODEBANC", "ABC", "BANC"},
		{"same", "a", "a", "a"},
		{"not found", "a", "aa", ""},
		{"exact match", "ABC", "ABC", "ABC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinWindow(tt.s, tt.t)
			if got != tt.want {
				t.Errorf("MinWindow(%q, %q) = %q, want %q", tt.s, tt.t, got, tt.want)
			}
		})
	}
}

func TestAlienOrder(t *testing.T) {
	tests := []struct {
		name      string
		words     []string
		wantEmpty bool // true if result must be ""
		wantLen   int  // expected length if non-empty
	}{
		{"basic", []string{"wrt", "wrf", "er", "ett", "rftt"}, false, 5},
		{"invalid prefix", []string{"abc", "ab"}, true, 0},
		{"simple", []string{"z", "x"}, false, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AlienOrder(tt.words)
			if tt.wantEmpty && got != "" {
				t.Errorf("AlienOrder(%v) = %q, want empty string", tt.words, got)
			}
			if !tt.wantEmpty && len(got) != tt.wantLen {
				t.Errorf("AlienOrder(%v) = %q (len %d), want len %d", tt.words, got, len(got), tt.wantLen)
			}
		})
	}
}

func TestIsMatch(t *testing.T) {
	tests := []struct {
		name string
		s, p string
		want bool
	}{
		{"no match", "aa", "a", false},
		{"star zero or more", "aa", "a*", true},
		{"dot star", "ab", ".*", true},
		{"complex", "aab", "c*a*b", true},
		{"exact", "mississippi", "mis*is*p*.", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsMatch(tt.s, tt.p)
			if got != tt.want {
				t.Errorf("IsMatch(%q, %q) = %v, want %v", tt.s, tt.p, got, tt.want)
			}
		})
	}
}

func TestMinDistance(t *testing.T) {
	tests := []struct {
		name         string
		word1, word2 string
		want         int
	}{
		{"horse to ros", "horse", "ros", 3},
		{"intention to execution", "intention", "execution", 5},
		{"empty to word", "", "abc", 3},
		{"word to empty", "abc", "", 3},
		{"same", "abc", "abc", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinDistance(tt.word1, tt.word2)
			if got != tt.want {
				t.Errorf("MinDistance(%q, %q) = %d, want %d", tt.word1, tt.word2, got, tt.want)
			}
		})
	}
}

func TestJobScheduling(t *testing.T) {
	tests := []struct {
		name      string
		startTime []int
		endTime   []int
		profit    []int
		want      int
	}{
		{"example 1", []int{1, 2, 3, 3}, []int{3, 4, 5, 6}, []int{50, 10, 40, 70}, 120},
		{"example 2", []int{1, 2, 3, 4, 6}, []int{3, 5, 10, 6, 9}, []int{20, 20, 100, 70, 60}, 150},
		{"no overlap", []int{1, 4}, []int{3, 6}, []int{10, 20}, 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JobScheduling(tt.startTime, tt.endTime, tt.profit)
			if got != tt.want {
				t.Errorf("JobScheduling() = %d, want %d", got, tt.want)
			}
		})
	}
}
