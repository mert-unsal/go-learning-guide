package io_files
import (
"strings"
"testing"
)
func TestWriteAndReadFile(t *testing.T) {
content := "hello, world!\nline two"
got, err := WriteAndReadFileSolution("testfile_tmp.txt", content)
if err != nil {
t.Fatalf("WriteAndReadFile error: %v", err)
}
if got != content {
t.Errorf("WriteAndReadFile = %q, want %q", got, content)
}
}
func TestCountLines(t *testing.T) {
tests := []struct{ content string; want int }{
{"line1\nline2\nline3", 3},
{"single", 1},
{"", 0},
{"a\nb\nc\nd", 4},
}
for _, tt := range tests {
got := CountLinesSolution(tt.content)
if got != tt.want {
t.Errorf("CountLines(%q) = %d, want %d", tt.content, got, tt.want)
}
}
}
func TestWordFrequency(t *testing.T) {
r := strings.NewReader("Go is great\ngo is fun\nGO")
freq := WordFrequencySolution(r)
if freq["go"] != 3 {
t.Errorf("freq[go] = %d, want 3", freq["go"])
}
if freq["is"] != 2 {
t.Errorf("freq[is] = %d, want 2", freq["is"])
}
}
func TestCopyReader(t *testing.T) {
r := strings.NewReader("hello world")
got, err := CopyReaderSolution(r)
if err != nil || got != "hello world" {
t.Errorf("CopyReader = (%q, %v), want (hello world, nil)", got, err)
}
}