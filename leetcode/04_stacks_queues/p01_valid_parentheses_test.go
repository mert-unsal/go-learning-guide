package stacks_queues

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"simple pair", "()", true},
		{"multiple types", "()[]{}", true},
		{"wrong order", "(]", false},
		{"interleaved", "([)]", false},
		{"nested", "{[]}", true},
		{"empty", "", true},
		{"only open", "(((", false},
		{"only close", ")))", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValid(tt.s)
			if got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
