package linked_list
import "testing"
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
