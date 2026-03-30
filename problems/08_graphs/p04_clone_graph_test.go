package graphs

import "testing"

func TestCloneGraph(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"nil input"},
		{"single node"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: implement test cases
			t.Skip("not implemented")
		})
	}
}
