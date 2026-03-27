package hard

import (
	"reflect"
	"testing"
)

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
