package strings_problems

import "testing"

func TestCountAndSay(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"base case", 1, "1"},
		{"n=4", 4, "1211"},
		{"n=5", 5, "111221"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountAndSay(tt.n)
			if got != tt.want {
				t.Errorf("CountAndSay(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}
