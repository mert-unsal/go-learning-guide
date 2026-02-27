package linked_list

import (
	"reflect"
	"testing"
)

func TestReverseList(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"normal", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"two nodes", []int{1, 2}, []int{2, 1}},
		{"single", []int{1}, []int{1}},
		{"empty", []int{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(ReverseList(newList(tt.input)))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReverseList(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestReverseListRecursive(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := []int{5, 4, 3, 2, 1}
	got := toSlice(ReverseListRecursive(newList(input)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReverseListRecursive(%v) = %v, want %v", input, got, want)
	}
}

func TestMergeTwoLists(t *testing.T) {
	tests := []struct {
		name  string
		list1 []int
		list2 []int
		want  []int
	}{
		{"normal", []int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
		{"both empty", []int{}, []int{}, nil},
		{"one empty", []int{}, []int{0}, []int{0}},
		{"first longer", []int{1, 3, 5, 7}, []int{2, 4}, []int{1, 2, 3, 4, 5, 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(MergeTwoLists(newList(tt.list1), newList(tt.list2)))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeTwoLists(%v, %v) = %v, want %v", tt.list1, tt.list2, got, tt.want)
			}
		})
	}
}

func TestHasCycle(t *testing.T) {
	t.Run("no cycle", func(t *testing.T) {
		head := newList([]int{1, 2, 3})
		if HasCycle(head) {
			t.Error("HasCycle on acyclic list returned true")
		}
	})
	t.Run("has cycle", func(t *testing.T) {
		// Build 3→2→0→-4 with -4 pointing back to 2
		head := &ListNode{Val: 3}
		n2 := &ListNode{Val: 2}
		n3 := &ListNode{Val: 0}
		n4 := &ListNode{Val: -4}
		head.Next = n2
		n2.Next = n3
		n3.Next = n4
		n4.Next = n2 // cycle!
		if !HasCycle(head) {
			t.Error("HasCycle on cyclic list returned false")
		}
	})
	t.Run("single no cycle", func(t *testing.T) {
		head := &ListNode{Val: 1}
		if HasCycle(head) {
			t.Error("expected false")
		}
	})
}

func TestRemoveNthFromEnd(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		n     int
		want  []int
	}{
		{"remove second from end", []int{1, 2, 3, 4, 5}, 2, []int{1, 2, 3, 5}},
		{"remove only node", []int{1}, 1, nil},
		{"remove last", []int{1, 2}, 1, []int{1}},
		{"remove first", []int{1, 2, 3}, 3, []int{2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(RemoveNthFromEnd(newList(tt.input), tt.n))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveNthFromEnd(%v, %d) = %v, want %v", tt.input, tt.n, got, tt.want)
			}
		})
	}
}

func TestMiddleNode(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		wantFirst int
	}{
		{"odd length", []int{1, 2, 3, 4, 5}, 3},
		{"even length", []int{1, 2, 3, 4}, 3},
		{"single", []int{1}, 1},
		{"two nodes", []int{1, 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MiddleNode(newList(tt.input))
			if got == nil || got.Val != tt.wantFirst {
				var gotVal int
				if got != nil {
					gotVal = got.Val
				}
				t.Errorf("MiddleNode(%v).Val = %d, want %d", tt.input, gotVal, tt.wantFirst)
			}
		})
	}
}
