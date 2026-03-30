package graphs

import "testing"

func TestSurroundedRegions(t *testing.T) {
	tests := []struct {
		name  string
		board [][]byte
	}{
		{"basic", [][]byte{
			{'X', 'X', 'X', 'X'},
			{'X', 'O', 'O', 'X'},
			{'X', 'X', 'O', 'X'},
			{'X', 'O', 'X', 'X'},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: implement test validation
			SurroundedRegions(tt.board)
			t.Skip("not implemented")
		})
	}
}
