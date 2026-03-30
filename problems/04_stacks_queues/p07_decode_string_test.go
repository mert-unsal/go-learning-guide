package stacks_queues

import "testing"

func TestDecodeString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"basic", "3[a]2[bc]", "aaabcbc"},
		{"nested", "3[a2[c]]", "accaccacc"},
		{"no encoding", "abc", "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecodeString(tt.s)
			if got != tt.want {
				t.Errorf("DecodeString(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}
