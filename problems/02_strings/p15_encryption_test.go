package strings_problems

import "testing"

func TestEncryption(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"haveaniceday", "hae and via ecy"},
		{"feedthedog", "fto ehg ee dd"},
		{"chillout", "clu hlt io"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Encryption(tt.input); got != tt.want {
				t.Errorf("Encryption(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
