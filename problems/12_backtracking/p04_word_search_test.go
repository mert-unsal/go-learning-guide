package backtracking

import "testing"

func TestExist(t *testing.T) {
	board := [][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}
	if !Exist(board, "ABCCED") {
		t.Error("Exist(board, ABCCED) = false, want true")
	}
}
